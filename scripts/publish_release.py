#!/usr/bin/env python3
"""Publish the exact CI candidate after preflighting GitHub and npm together."""
import argparse
import base64
import hashlib
import json
import os
import pathlib
import re
import subprocess
import sys
import tempfile
import urllib.error
import urllib.request
from version import validate, SEMVER
from release_identity import resolve

REGISTRY = 'https://registry.npmjs.org/woobe-cli'


def run(*args, capture=False):
    result = subprocess.run(args, check=True, text=True, capture_output=capture)
    return result.stdout.strip() if capture else None


def fetch(url, token=None, missing=False):
    headers = {'Accept': 'application/json', 'User-Agent': 'woobe-cli-release'}
    if token:
        headers['Authorization'] = f'Bearer {token}'
    try:
        with urllib.request.urlopen(urllib.request.Request(url, headers=headers), timeout=60) as response:
            payload = response.read(80 * 1024 * 1024 + 1)
            if len(payload) > 80 * 1024 * 1024:
                raise ValueError('registry response exceeds distribution limit')
            return payload
    except urllib.error.HTTPError as error:
        if missing and error.code == 404:
            return None
        raise RuntimeError(f'Release lookup failed: HTTP {error.code} at {url}') from None


def github(path, missing=False):
    payload = fetch('https://api.github.com/' + path, os.environ.get('GH_TOKEN'), missing)
    return json.loads(payload) if payload is not None else None


def registry():
    payload = fetch(REGISTRY, missing=True)
    return json.loads(payload) if payload is not None else None


def semver_key(value):
    validate(value)
    core, separator, prerelease = value.partition('-')
    identifiers = tuple((0, int(part)) if part.isdigit() else (1, part)
                        for part in prerelease.split('.')) if separator else ()
    return (*map(int, core.split('.')), 0 if separator else 1, identifiers)


def assert_advances(version, published):
    known = [value for value in published if SEMVER.fullmatch(value)]
    if known and semver_key(version) <= max(map(semver_key, known)):
        raise ValueError('new releases must advance the highest published Semantic Version')


def verify_npm(packument, version, package):
    entry = packument.get('versions', {}).get(version) if packument else None
    if entry is None:
        return False
    integrity = 'sha512-' + base64.b64encode(hashlib.sha512(package.read_bytes()).digest()).decode()
    if entry.get('dist', {}).get('integrity') != integrity:
        raise ValueError('immutable npm version already exists with different bytes')
    url = entry['dist']['tarball']
    if not url.startswith('https://registry.npmjs.org/woobe-cli/-/'):
        raise ValueError('unexpected npm distribution URL')
    if fetch(url) != package.read_bytes():
        raise ValueError('published npm tarball differs from the validated candidate')
    return True


def verify_assets(repo, release, files):
    if release is None:
        return set()
    assets = release['assets']
    names = [asset['name'] for asset in assets]
    if len(names) != len(set(names)) or set(names) - set(files):
        raise ValueError('existing release has duplicate or unexpected assets')
    with tempfile.TemporaryDirectory() as directory:
        for asset in assets:
            expected = files[asset['name']]
            if asset['size'] != expected.stat().st_size:
                raise ValueError(f'immutable GitHub asset differs: {asset["name"]}')
            run('gh', 'release', 'download', release['tag_name'], '--repo', repo,
                '--pattern', asset['name'], '--dir', directory)
            if (pathlib.Path(directory) / asset['name']).read_bytes() != expected.read_bytes():
                raise ValueError(f'immutable GitHub asset differs: {asset["name"]}')
    return set(names)


def preflight(repo, version, commit, files, package):
    # Verify repository access first: an unauthenticated/private-repo 404 is
    # never interpreted as an absent release or tag.
    settings = github(f'repos/{repo}')
    if not settings.get('permissions', {}).get('push'):
        raise ValueError('GitHub credential lacks repository write permission')
    tag = github(f'repos/{repo}/git/ref/tags/v{version}', missing=True)
    if tag:
        obj = tag['object']
        for _ in range(5):
            if obj['type'] == 'commit':
                break
            if obj['type'] != 'tag':
                raise ValueError('release tag does not resolve to a commit')
            obj = github(f'repos/{repo}/git/tags/{obj["sha"]}')['object']
        if obj['type'] != 'commit' or obj['sha'] != commit:
            raise ValueError('immutable release tag points to a different commit')
    release = github(f'repos/{repo}/releases/tags/v{version}', missing=True)
    if release and not tag and (not release['draft'] or release['target_commitish'] != commit):
        raise ValueError('existing release has neither an immutable tag nor an exact draft revision')
    if release and release['prerelease'] != ('-' in version):
        raise ValueError('existing GitHub Release has a conflicting prerelease flag')
    assets = verify_assets(repo, release, files)
    packument = registry()
    npm_exists = verify_npm(packument, version, package)
    if not npm_exists:
        published = list((packument or {}).get('versions', {}))
        page = 1
        while True:
            releases = github(f'repos/{repo}/releases?per_page=100&page={page}')
            published.extend(item['tag_name'][1:] for item in releases
                             if not item['draft'] and item['tag_name'].startswith('v')
                             and item['tag_name'] != f'v{version}')
            if len(releases) < 100:
                break
            page += 1
        assert_advances(version, published)
        if packument is None and not os.environ.get('NODE_AUTH_TOKEN'):
            raise ValueError('first npm publication requires repository secret NPM_TOKEN; then configure Trusted Publisher for release.yml')
        if os.environ.get('NODE_AUTH_TOKEN'):
            run('npm', 'whoami', '--registry=https://registry.npmjs.org', capture=True)
        elif not os.environ.get('ACTIONS_ID_TOKEN_REQUEST_URL'):
            raise ValueError('npm publication requires NPM_TOKEN or GitHub Actions Trusted Publisher')
    return release, assets, npm_exists


def publish(repo, version, commit, files, package, preflight_only=False):
    release, assets, npm_exists = preflight(repo, version, commit, files, package)
    if preflight_only:
        print('Both destinations passed preflight; no publication requested')
        return
    tag = f'v{version}'
    npm_tag = 'next' if '-' in version else 'latest'
    if not npm_exists:
        run('npm', 'publish', str(package), '--ignore-scripts', '--provenance',
            '--access', 'public', '--tag', npm_tag, '--registry=https://registry.npmjs.org')
    # npm is immutable and verified before creating the permanent GitHub record.
    if not verify_npm(registry(), version, package):
        raise ValueError('npm publication is missing after publish')
    run(sys.executable, 'scripts/smoke_npm.py', '--commit', commit,
        '--report', 'dist/published-npm-smoke.json')
    if release is None:
        notes = pathlib.Path('dist/release-notes.md')
        notes.write_text(f'Woobe CLI {version}\n\nSource: {commit}\n\n'
                         f'Install: `npm install --global woobe-cli@{version}`\n\n'
                         f'Run: `npx --package=woobe-cli@{version} woobe version`\n\n'
                         'Native archives: Linux, macOS and Windows; amd64 and arm64.\n'
                         'Verify downloads with SHA256SUMS and release-manifest.json.\n')
        flags = ['--prerelease'] if '-' in version else []
        run('gh', 'release', 'create', tag, '--repo', repo, '--target', commit,
            '--draft', '--title', f'Woobe CLI {tag}', '--notes-file', str(notes), *flags)
    # Upload missing bytes only. Never --clobber an immutable release asset.
    for name, path in files.items():
        if name not in assets:
            run('gh', 'release', 'upload', tag, str(path), '--repo', repo)
    release = github(f'repos/{repo}/releases/tags/{tag}')
    if verify_assets(repo, release, files) != set(files):
        raise ValueError('GitHub publication is incomplete')
    if release['draft']:
        run('gh', 'release', 'edit', tag, '--repo', repo, '--draft=false')
    # Repeat all lookups after publication, including the tag and npm bytes.
    terminal, terminal_assets, terminal_npm = preflight(repo, version, commit, files, package)
    if terminal['draft'] or terminal_assets != set(files) or not terminal_npm:
        raise ValueError('release transaction did not complete')
    print(f'Published {tag}: all six native archives, npm and the permanent release manifest verified')


def candidate(version, commit, manifest_sha256):
    validate(version)
    if not re.fullmatch(r'[a-f0-9]{40}', commit) or not re.fullmatch(r'[a-f0-9]{64}', manifest_sha256):
        raise ValueError('invalid candidate identity')
    if run('git', 'rev-parse', 'HEAD', capture=True) != commit:
        raise ValueError('checkout differs from validated source')
    resolve(version, commit)
    manifest_path = pathlib.Path('dist/artifacts.json')
    if hashlib.sha256(manifest_path.read_bytes()).hexdigest() != manifest_sha256:
        raise ValueError('candidate manifest differs from successful CI output')
    run(sys.executable, 'scripts/verify_artifacts.py', '--commit', commit, '--version', version)
    manifest = json.loads(manifest_path.read_text())
    package = pathlib.Path('dist') / manifest['npm']['name']
    integrity = 'sha512-' + base64.b64encode(hashlib.sha512(package.read_bytes()).digest()).decode()
    record = dict(schema_version='1', version=version, commit=commit,
                  candidate_manifest_sha256=manifest_sha256, artifacts=manifest['artifacts'],
                  npm=dict(package='woobe-cli', version=version, integrity=integrity,
                           sha256=manifest['npm']['sha256']))
    permanent = pathlib.Path('dist/release-manifest.json')
    permanent.write_text(json.dumps(record, indent=2) + '\n')
    names = [item['name'] for item in manifest['artifacts']]
    names.extend([package.name, 'SHA256SUMS', 'artifacts.json', permanent.name])
    return {name: pathlib.Path('dist') / name for name in names}, package


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--version', required=True)
    parser.add_argument('--commit', required=True)
    parser.add_argument('--manifest-sha256', required=True)
    parser.add_argument('--repo', default=os.environ.get('GITHUB_REPOSITORY', 'A1b3rt0M3rcad0/woobe-cli'))
    parser.add_argument('--preflight-only', action='store_true')
    args = parser.parse_args()
    if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+', args.repo):
        parser.error('invalid GitHub repository')
    files, package = candidate(args.version, args.commit, args.manifest_sha256)
    publish(args.repo, args.version, args.commit, files, package, args.preflight_only)

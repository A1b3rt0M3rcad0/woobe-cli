#!/usr/bin/env python3
"""Publish the exact CI native candidate to GitHub, without npm or extra secrets."""
import argparse
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
from version import validate, SEMVER, semver_key
from release_identity import resolve


def run(*args, capture=False):
    result = subprocess.run(args, check=True, text=True, capture_output=capture)
    return result.stdout.strip() if capture else None


def fetch(url, token=None, missing=False):
    headers = {'Accept': 'application/json', 'User-Agent': 'woobe-cli-release'}
    if token:
        headers['Authorization'] = f'Bearer {token}'
    try:
        with urllib.request.urlopen(urllib.request.Request(url, headers=headers), timeout=60) as response:
            return response.read()
    except urllib.error.HTTPError as error:
        if missing and error.code == 404:
            return None
        raise RuntimeError(f'Release lookup failed: HTTP {error.code} at {url}') from None


def github(path, missing=False):
    payload = fetch('https://api.github.com/' + path, os.environ.get('GH_TOKEN'), missing)
    return json.loads(payload) if payload is not None else None


def published_versions(repo):
    versions = []
    page = 1
    while True:
        releases = github(f'repos/{repo}/releases?per_page=100&page={page}')
        versions.extend(item['tag_name'][1:] for item in releases
                        if not item['draft'] and item['tag_name'].startswith('v')
                        and SEMVER.fullmatch(item['tag_name'][1:]))
        if len(releases) < 100:
            return versions
        page += 1


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


def preflight(repo, version, commit, files):
    # Read repository access first: a private-repo/authentication 404 must not
    # be interpreted as a missing tag or release.
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
    if not tag and not release:
        versions = published_versions(repo)
        if versions and semver_key(version) <= max(map(semver_key, versions)):
            raise ValueError('new releases must advance the highest published Semantic Version')
    return release, assets, bool(tag)


def publish(repo, version, commit, files, preflight_only=False):
    release, assets, reserved = preflight(repo, version, commit, files)
    if preflight_only:
        print('GitHub native release passed preflight; no publication requested')
        return
    tag = f'v{version}'
    if not reserved:
        # Reserve the validated identity before uploading. A failed upload can
        # be retried; later pushes will not reuse this version for other bytes.
        run('gh', 'api', f'repos/{repo}/git/refs', '--method', 'POST',
            '-f', f'ref=refs/tags/{tag}', '-f', f'sha={commit}')
    if release is None:
        notes = pathlib.Path('dist/release-notes.md')
        links = '\n'.join(f'- [{name}](https://github.com/{repo}/releases/download/{tag}/{name})'
                          for name in files if name.endswith(('.tar.gz', '.zip')))
        notes.write_text(f'Woobe CLI {version}\n\nSource: {commit}\n\n'
                         'Download your platform archive, verify SHA256SUMS, extract\n'
                         'and put woobe (woobe.exe on Windows) on PATH. Run `woobe version`.\n'
                         'No Go, Node.js or npm installation is required.\n\n' + links + '\n')
        flags = ['--prerelease'] if '-' in version else []
        run('gh', 'release', 'create', tag, '--repo', repo, '--target', commit,
            '--verify-tag', '--draft', '--generate-notes', '--title', f'Woobe CLI {tag}', '--notes-file', str(notes), *flags)
    for name, path in files.items():
        if name not in assets:
            run('gh', 'release', 'upload', tag, str(path), '--repo', repo)
    release = github(f'repos/{repo}/releases/tags/{tag}')
    if verify_assets(repo, release, files) != set(files):
        raise ValueError('GitHub publication is incomplete')
    # Remote files have just been compared byte for byte with this candidate.
    run(sys.executable, 'scripts/smoke_artifacts.py', '--commit', commit,
        '--report', 'dist/published-native-smoke.json')
    if release['draft']:
        stable = [v for v in published_versions(repo) if '-' not in v]
        latest = '-' not in version and (not stable or semver_key(version) >= max(map(semver_key, stable)))
        run('gh', 'release', 'edit', tag, '--repo', repo, '--draft=false',
            f'--latest={str(latest).lower()}')
    terminal, terminal_assets, terminal_tag = preflight(repo, version, commit, files)
    if terminal['draft'] or terminal_assets != set(files) or not terminal_tag:
        raise ValueError('release transaction did not complete')
    print(f'Published {tag}: six native archives, immutable tag and permanent release manifest verified')


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
    if 'npm' in manifest:
        raise ValueError('this release accepts native-only candidates; npm distribution is deferred')
    record = dict(schema_version='1', version=version, commit=commit,
                  candidate_manifest_sha256=manifest_sha256, compiler=manifest['compiler'],
                  artifacts=manifest['artifacts'])
    permanent = pathlib.Path('dist/release-manifest.json')
    permanent.write_text(json.dumps(record, indent=2) + '\n')
    names = [item['name'] for item in manifest['artifacts']]
    names.extend(['SHA256SUMS', 'artifacts.json', permanent.name])
    return {name: pathlib.Path('dist') / name for name in names}


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
    files = candidate(args.version, args.commit, args.manifest_sha256)
    publish(args.repo, args.version, args.commit, files, args.preflight_only)

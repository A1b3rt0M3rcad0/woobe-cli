#!/usr/bin/env python3
"""Resolve a reviewed canonical revision; never merge or change VERSION."""
import argparse
import json
import pathlib
import re
import subprocess
from version import validate


def git(*args):
    return subprocess.check_output(['git', *args], text=True).strip()


def resolve(version, revision=''):
    validate(version)
    if revision and not re.fullmatch(r'[a-f0-9]{40}', revision):
        raise ValueError('revision must be a full lowercase commit SHA')
    tag = f'refs/tags/v{version}'
    exists = subprocess.run(['git', 'show-ref', '--verify', '--quiet', tag]).returncode == 0
    tagged = git('rev-parse', f'{tag}^{{commit}}') if exists else ''
    sha = revision or tagged or git('rev-parse', 'origin/master')
    if tagged and tagged != sha:
        raise ValueError('existing immutable tag points to a different commit')
    subprocess.run(['git', 'merge-base', '--is-ancestor', sha, 'origin/master'], check=True)
    canonical = git('show', f'{sha}:VERSION')
    manifest = json.loads(git('show', f'{sha}:packages/woobe-cli/package.json'))
    if canonical != version or manifest['version'] != version:
        raise ValueError('review VERSION and npm manifest in a PR first; release does not merge or bump versions')
    return sha


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--version', required=True)
    parser.add_argument('--revision', default='')
    parser.add_argument('--output', type=pathlib.Path, required=True)
    args = parser.parse_args()
    sha = resolve(args.version, args.revision)
    with args.output.open('a') as output:
        output.write(f'version={args.version}\nrevision={sha}\n')
    print(f'Release v{args.version} will validate exact master revision {sha}')

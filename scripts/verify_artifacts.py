#!/usr/bin/env python3
"""Verify released archive contents and the exact manifest/checksum contract."""
import argparse, hashlib, json, pathlib, re, tarfile, zipfile
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--commit")
options = parser.parse_args()
root = pathlib.Path('dist')
manifest = json.loads((root / 'artifacts.json').read_text())
assert len(manifest['artifacts']) == 6
assert re.fullmatch(r'[a-f0-9]{40}', manifest['commit'])
if options.commit:
    assert manifest['commit'] == options.commit, 'artifact source commit differs from expected checkout'
assert {(a['os'], a['arch']) for a in manifest['artifacts']} == {(s, a) for s in ('linux','darwin','windows') for a in ('amd64','arm64')}
assert len({a['name'] for a in manifest['artifacts']}) == 6
checksums = dict(line.split()[::-1] for line in (root / 'SHA256SUMS').read_text().splitlines())
assert set(checksums) == {a['name'] for a in manifest['artifacts']}
for artifact in manifest['artifacts']:
    suffix = 'zip' if artifact['os'] == 'windows' else 'tar.gz'
    assert artifact['name'] == f"woobe_{manifest['version']}_{artifact['os']}_{artifact['arch']}.{suffix}"
    path = root / artifact['name']
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    assert digest == artifact['sha256'] == checksums[path.name]
    assert path.stat().st_size == artifact['bytes']
    if path.suffix == '.zip':
        with zipfile.ZipFile(path) as archive:
            entries = archive.namelist()
            assert len(entries) == len(set(entries))
            names = set(entries)
            schema = json.loads(archive.read('manifest.schema.json'))
            resources = json.loads(archive.read('resources.schema.json'))
    else:
        with tarfile.open(path) as archive:
            entries = archive.getmembers()
            assert all(entry.isfile() for entry in entries)
            assert len(entries) == len({entry.name for entry in entries})
            names = {entry.name for entry in entries}
            schema = json.load(archive.extractfile('manifest.schema.json'))
            resources = json.load(archive.extractfile('resources.schema.json'))
    executable = 'woobe.exe' if artifact['os'] == 'windows' else 'woobe'
    assert names == {executable, 'README.md', 'USAGE.md', 'manifest.schema.json', 'resources.schema.json'}
    assert schema['$id'] == 'urn:woobe:manifest:steps:1'
    assert resources['$id'] == 'urn:woobe:manifest:resources:2'
print('Verified six archives, schemas and checksum metadata')

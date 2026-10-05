#!/usr/bin/env python3
"""Verify released archive contents and the exact manifest/checksum contract."""
import hashlib, json, pathlib, tarfile, zipfile
root = pathlib.Path('dist')
manifest = json.loads((root / 'artifacts.json').read_text())
assert len(manifest['artifacts']) == 6
checksums = dict(line.split()[::-1] for line in (root / 'SHA256SUMS').read_text().splitlines())
assert set(checksums) == {a['name'] for a in manifest['artifacts']}
for artifact in manifest['artifacts']:
    path = root / artifact['name']
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    assert digest == artifact['sha256'] == checksums[path.name]
    assert path.stat().st_size == artifact['bytes']
    if path.suffix == '.zip':
        with zipfile.ZipFile(path) as archive:
            names = set(archive.namelist())
            schema = json.loads(archive.read('manifest.schema.json'))
            resources = json.loads(archive.read('resources.schema.json'))
    else:
        with tarfile.open(path) as archive:
            names = set(archive.getnames())
            schema = json.load(archive.extractfile('manifest.schema.json'))
            resources = json.load(archive.extractfile('resources.schema.json'))
    executable = 'woobe.exe' if artifact['os'] == 'windows' else 'woobe'
    assert names == {executable, 'README.md', 'USAGE.md', 'manifest.schema.json', 'resources.schema.json'}
    assert schema['$id'] == 'urn:woobe:manifest:steps:1'
    assert resources['$id'] == 'urn:woobe:manifest:resources:2'
print('Verified six archives, schemas and checksum metadata')

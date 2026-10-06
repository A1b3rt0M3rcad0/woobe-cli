#!/usr/bin/env python3
"""Verify released archive contents and the exact manifest/checksum contract."""
import argparse, hashlib, json, pathlib, re, tarfile, zipfile
from version import validate
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--commit")
parser.add_argument("--version")
options = parser.parse_args()
root = pathlib.Path('dist')
manifest = json.loads((root / 'artifacts.json').read_text())
validate(manifest['version'])
if options.version:
    assert manifest['version'] == options.version, 'artifact version differs from expected release'
assert len(manifest['artifacts']) == 6
assert re.fullmatch(r'[a-f0-9]{40}', manifest['commit'])
if options.commit:
    assert manifest['commit'] == options.commit, 'artifact source commit differs from expected checkout'
assert {(a['os'], a['arch']) for a in manifest['artifacts']} == {(s, a) for s in ('linux','darwin','windows') for a in ('amd64','arm64')}
assert len({a['name'] for a in manifest['artifacts']}) == 6
checksums = dict(line.split()[::-1] for line in (root / 'SHA256SUMS').read_text().splitlines())
assert set(checksums) == {a['name'] for a in manifest['artifacts']} | {manifest['npm']['name']}
native_binaries = {}
native_notices = []
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
            native_binaries[(artifact['os'], artifact['arch'])] = archive.read('woobe.exe')
            native_notices.append(archive.read('THIRD_PARTY_NOTICES.txt'))
    else:
        with tarfile.open(path) as archive:
            entries = archive.getmembers()
            assert all(entry.isfile() for entry in entries)
            assert len(entries) == len({entry.name for entry in entries})
            names = {entry.name for entry in entries}
            schema = json.load(archive.extractfile('manifest.schema.json'))
            resources = json.load(archive.extractfile('resources.schema.json'))
            native_binaries[(artifact['os'], artifact['arch'])] = archive.extractfile('woobe').read()
            native_notices.append(archive.extractfile('THIRD_PARTY_NOTICES.txt').read())
    executable = 'woobe.exe' if artifact['os'] == 'windows' else 'woobe'
    assert names == {executable, 'README.md', 'USAGE.md', 'manifest.schema.json', 'resources.schema.json', 'THIRD_PARTY_NOTICES.txt'}
    assert schema['$id'] == 'urn:woobe:manifest:steps:1'
    assert resources['$id'] == 'urn:woobe:manifest:resources:2'
npm = manifest['npm']
assert npm['name'] == f"woobe-cli-{manifest['version']}.tgz"
npm_path = root / npm['name']
assert npm_path.stat().st_size == npm['bytes']
assert hashlib.sha256(npm_path.read_bytes()).hexdigest() == npm['sha256'] == checksums[npm['name']]
with tarfile.open(npm_path) as archive:
    entries = archive.getmembers()
    assert all(entry.isfile() for entry in entries)
    names = {entry.name for entry in entries}
    assert len(names) == len(entries)
    package = json.load(archive.extractfile('package/package.json'))
    assert package['name'] == 'woobe-cli' and package['version'] == manifest['version']
    assert package.get('private') is not True
    assert package['bin'] == {'woobe': 'bin/woobe.cjs'}
    assert not package.get('scripts') and not package.get('dependencies')
    expected = {'package/package.json', 'package/bin/woobe.cjs', 'package/README.md', 'package/USAGE.md', 'package/THIRD_PARTY_NOTICES.txt'}
    notice = archive.extractfile('package/THIRD_PARTY_NOTICES.txt').read()
    assert notice and all(value == notice for value in native_notices)
    for (system, arch), binary in native_binaries.items():
        executable = 'woobe.exe' if system == 'windows' else 'woobe'
        name = f'package/vendor/{system}-{arch}/{executable}'
        expected.add(name)
        assert archive.extractfile(name).read() == binary, 'npm and native archive executables differ'
        assert archive.getmember(name).mode & 0o111, 'packaged executable is not executable'
    assert names == expected
    assert archive.getmember('package/bin/woobe.cjs').mode & 0o111
print('Verified six archives, npm package, identical binaries, schemas and checksum metadata')

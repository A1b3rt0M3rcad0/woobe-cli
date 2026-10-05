#!/usr/bin/env python3
"""Exercise the downloaded executable for this host, without a backend or secrets."""
import argparse
import json
import os
import pathlib
import platform
import subprocess
import tarfile
import tempfile
import zipfile


def smoke(root, commit):
    manifest = json.loads((root / 'artifacts.json').read_text())
    if manifest['commit'] != commit:
        raise ValueError('artifact source commit differs from expected checkout')
    system = {'Linux': 'linux', 'Darwin': 'darwin', 'Windows': 'windows'}[platform.system()]
    arch = {'x86_64': 'amd64', 'AMD64': 'amd64', 'arm64': 'arm64', 'aarch64': 'arm64'}[platform.machine()]
    matches = [a for a in manifest['artifacts'] if a['os'] == system and a['arch'] == arch]
    if len(matches) != 1:
        raise ValueError('exactly one archive must match the native host')
    artifact = matches[0]
    executable = 'woobe.exe' if system == 'windows' else 'woobe'
    archive_path = root / artifact['name']
    if system == 'windows':
        with zipfile.ZipFile(archive_path) as archive:
            binary = archive.read(executable)
    else:
        with tarfile.open(archive_path) as archive:
            binary = archive.extractfile(executable).read()
    # Extract only the known executable; never trust archive paths for extraction.
    with tempfile.TemporaryDirectory() as directory:
        path = pathlib.Path(directory) / executable
        path.write_bytes(binary)
        path.chmod(0o700)
        env = {k: v for k, v in os.environ.items() if not k.startswith('WOOBE_')}

        def invoke(args, body=None, code=0):
            command = [str(path), *args, '--output', 'json', '--config', str(pathlib.Path(directory) / 'config.json')]
            result = subprocess.run(command, input=body, capture_output=True, text=True, env=env, cwd=directory, timeout=30)
            if result.returncode != code:
                raise ValueError(f'{args}: exit {result.returncode}, expected {code}: {result.stderr}')
            envelope = json.loads(result.stdout)
            if envelope['success'] != (code == 0):
                raise ValueError(f'{args}: envelope success differs from exit code')
            if code and envelope['error']['exit_code'] != code:
                raise ValueError(f'{args}: structured error code differs from exit code')
            return envelope

        version = invoke(['version'])['data']
        if any(version[k] != v for k, v in {'commit': commit, 'version': manifest['version'], 'os': system, 'arch': arch}.items()):
            raise ValueError('packaged binary identity differs from archive metadata')
        discovery = invoke(['help'])['data']
        commands = {row['command'] for row in discovery}
        if not {'manifest validate', 'manifest apply', 'manifest reconcile', 'runtime target run', 'version'} <= commands:
            raise ValueError('discovery is missing essential executable commands')
        for revision, identity in [('1', 'urn:woobe:manifest:steps:1'), ('2', 'urn:woobe:manifest:resources:2')]:
            schema = invoke(['schema', '--command', 'manifest validate', '--kind', 'document', '--manifest-version', revision])['data']
            if schema['$id'] != identity:
                raise ValueError('packaged manifest schema has incorrect identity')
        body = json.dumps({'schema_version': '2', 'project_id': 'p', 'resources': [{'key': 'a', 'kind': 'Agent', 'action': 'update', 'resource_id': 'a', 'spec': {'name': 'Smoke'}}]})
        validation = invoke(['manifest', 'validate', '--file', '-', '--project', 'p'], body)['data']
        if validation['valid'] is not True or validation['authorization'] != 'not_evaluated':
            raise ValueError('local manifest validation has incorrect semantics')
        invoke(['manifest', 'validate', '--file', '-'], '{', code=2)
    return {'commit': commit, 'version': manifest['version'], 'os': system, 'arch': arch,
            'archive': artifact['name'], 'success': True, 'checks': ['identity', 'discovery', 'schemas', 'manifest', 'invalid-input'],
            'backend_acceptance': 'not_evaluated', 'credential_provider_acceptance': 'not_evaluated'}


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--commit', required=True)
    parser.add_argument('--report', required=True)
    options = parser.parse_args()
    report = smoke(pathlib.Path('dist'), options.commit)
    pathlib.Path(options.report).write_text(json.dumps(report, indent=2) + '\n')
    print(json.dumps(report))

#!/usr/bin/env python3
"""Exercise the downloaded executable for this host, without a backend or secrets."""
import argparse
import json
import os
import pathlib
import platform
import http.server
import threading
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
        if validation['valid'] is not True or validation['source_schema_version'] != '2' or len(validation['manifest_hash']) != 64:
            raise ValueError('local manifest validation has incorrect semantics')
        invoke(['manifest', 'validate', '--file', '-'], '{', code=2)
        # A loopback fixture verifies the packaged client's body pagination on
        # each native OS; real Woobe policy remains the separate backend gate.
        class Pages(http.server.BaseHTTPRequestHandler):
            requests = []
            validation_reads = 0
            writes = 0

            def do_GET(self):
                from urllib.parse import urlparse, parse_qs
                route = urlparse(self.path)
                query = parse_qs(route.query)
                if route.path == '/openapi.json':
                    type(self).validation_reads += 1
                    schema = {'openapi': '3.1.0', 'paths': {'/ai/agents': {'post': {'requestBody': {'content': {'application/json': {'schema': {
                        'type': 'object', 'required': ['id', 'model_id', 'created_at'], 'properties': {
                            'id': {'type': 'string', 'readOnly': True},
                            'model_id': {'type': 'string', 'format': 'uuid'},
                            'created_at': {'type': 'string', 'format': 'date-time'}}}}}}}}}}
                    schema['paths']['/runtime/agents/{agent_id}/sessions'] = {'get': {'parameters': [
                        {'name': 'agent_id', 'in': 'path', 'required': True, 'schema': {'type': 'string'}},
                        {'name': 'project_id', 'in': 'query', 'required': True, 'schema': {'type': 'string'}},
                        {'name': 'limit', 'in': 'query', 'schema': {'type': 'integer', 'minimum': 1, 'maximum': 100}},
                        {'name': 'cursor', 'in': 'query', 'schema': {'type': 'string'}}]}}
                    self.send_response(200)
                    self.send_header('Content-Type', 'application/json')
                    self.end_headers()
                    self.wfile.write(json.dumps(schema).encode())
                    return
                self.requests.append(query)
                if route.path != '/runtime/agents/a/sessions' or query.get('project_id') != ['p'] or query.get('limit') != ['1']:
                    self.send_error(400)
                    return
                terminal = query.get('cursor') == ['first']
                payload = {'success': True, 'data': {'items': [{'id': 'second' if terminal else 'first'}], 'has_next': not terminal, 'next_cursor': None if terminal else 'first'}}
                self.send_response(200)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(payload).encode())

            def do_POST(self):
                type(self).writes += 1
                if self.path != '/ai/agents':
                    self.send_error(400)
                    return
                self.rfile.read(int(self.headers.get('Content-Length', '0')))
                self.send_response(200)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(b'{"id":"a"}')

            def log_message(self, *_):
                pass

        server = http.server.HTTPServer(('127.0.0.1', 0), Pages)
        thread = threading.Thread(target=server.serve_forever, daemon=True)
        thread.start()
        try:
            origin = f'http://127.0.0.1:{server.server_port}'
            pages = invoke(['runtime', 'agent', 'sessions', 'a', '--project', 'p', '--api-url', origin, '--limit', '1', '--all'])
            if pages['data']['page_count'] != 2 or pages['meta']['collection_complete'] != 'verified' or pages['meta']['complete'] is not True or len(Pages.requests) != 2:
                raise ValueError('packaged pagination failed to consume both pages')
            partial = invoke(['runtime', 'agent', 'sessions', 'a', '--project', 'p', '--api-url', origin, '--limit', '1'])
            if partial['meta']['complete'] is not False or partial['meta']['collection_complete'] != 'partial':
                raise ValueError('single page incorrectly claims collection completeness')
            valid_body = {'model_id': '6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'created_at': '2026-10-05T20:40:46-03:00'}
            valid = invoke(['validate-input', '--command', 'project agent create', '--api-url', origin, '--file', '-'], json.dumps(valid_body))
            if valid['data']['validation_direction'] != 'request' or valid['data']['executed'] is not False:
                raise ValueError('packaged input validation did not report request direction')
            for invalid in [dict(valid_body, model_id='invalid'), dict(valid_body, id='server'), dict(valid_body, created_at='2026-02-29T00:00:00Z')]:
                denied = invoke(['project', 'agent', 'create', '--validate-body', '--api-url', origin, '--file', '-'], json.dumps(invalid), code=2)
                if denied['error']['write_outcome'] != 'not_attempted':
                    raise ValueError('invalid body was not refused before mutation')
            if Pages.writes != 0:
                raise ValueError('invalid request body caused a mutation')
            invoke(['project', 'agent', 'create', '--validate-body', '--api-url', origin, '--file', '-'], json.dumps(valid_body))
            if Pages.writes != 1 or Pages.validation_reads != 5:
                raise ValueError('packaged request validation/write counts differ')
            parameter_preflight = invoke(['validate-input', '--command', 'runtime agent sessions', '--validate-parameters', '--path-param', 'agent_id=a', '--query', 'project_id=p', '--query', 'limit=1', '--api-url', origin])
            if parameter_preflight['data']['path_query_validation'] != 'supported_schema_subset' or parameter_preflight['data']['body_validation'] != 'not_evaluated':
                raise ValueError('parameter-only preflight incorrectly reports its scope')
            pages_before = len(Pages.requests)
            invoke(['runtime', 'agent', 'sessions', 'a', '--project', 'p', '--api-url', origin, '--query', 'limit=101', '--validate-parameters'], code=2)
            if len(Pages.requests) != pages_before:
                raise ValueError('invalid query caused a resource read')
            invoke(['runtime', 'agent', 'sessions', 'a', '--project', 'p', '--api-url', origin, '--limit', '1', '--all', '--validate-parameters'])
            if len(Pages.requests) != pages_before + 2:
                raise ValueError('validated packaged pagination did not consume both pages')
        finally:
            server.shutdown()
            server.server_close()
            thread.join(timeout=5)
    return {'commit': commit, 'version': manifest['version'], 'os': system, 'arch': arch,
            'archive': artifact['name'], 'success': True, 'checks': ['identity', 'discovery', 'schemas', 'manifest', 'invalid-input', 'body-pagination', 'partial-collection', 'request-direction', 'uuid-date-time', 'validation-before-write', 'path-query-validation', 'validated-pagination'],
            'backend_acceptance': 'not_evaluated', 'credential_provider_acceptance': 'not_evaluated'}


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--commit', required=True)
    parser.add_argument('--report', required=True)
    options = parser.parse_args()
    report = smoke(pathlib.Path('dist'), options.commit)
    pathlib.Path(options.report).write_text(json.dumps(report, indent=2) + '\n')
    print(json.dumps(report))

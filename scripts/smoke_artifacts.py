#!/usr/bin/env python3
"""Exercise the downloaded executable for this host, without a backend or secrets."""
import argparse
import hashlib
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
    arch = {'x86_64': 'amd64', 'AMD64': 'amd64', 'arm64': 'arm64', 'ARM64': 'arm64', 'aarch64': 'arm64'}[platform.machine()]
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
        invoke(['init'])
        author = pathlib.Path(directory) / 'author.json'
        author.write_text(json.dumps({'kind': 'Agent', 'metadata': {'name': 'Support'}, 'spec': {'instructions': 'First line\nSecond line\n', 'limit': 9007199254740993}}), encoding='utf-8')
        invoke(['resources', 'create', 'agent', 'support', '--file', str(author)])
        rendered = (pathlib.Path(directory) / '.woobe/agents/support/agent.yaml').read_text(encoding='utf-8')
        if not rendered.startswith('kind: Agent\nmetadata:\n  key: support\n') or '  instructions: |\n    First line\n    Second line\n' not in rendered or '  limit: 9007199254740993\n' not in rendered:
            raise ValueError('packaged CLI did not write readable, lossless block YAML')
        invoke(['resources', 'move', 'agent', '@support', '--path', 'agents/support.json'])
        moved = json.loads((pathlib.Path(directory) / '.woobe/agents/support.json').read_text(encoding='utf-8'))
        if moved['spec'] != json.loads(author.read_text(encoding='utf-8'))['spec']:
            raise ValueError('packaged CLI changed author values when moving YAML to JSON')
        invoke(['resources', 'clone', 'agent', '@support', '--alias', 'copy', '--path', 'agents/copy.json'])
        invoke(['resources', 'used-by', 'agent', '@copy'])
        discovery = invoke(['help'])['data']
        commands = {row['command'] for row in discovery}
        if not {'manifest validate', 'manifest apply', 'manifest reconcile', 'runtime target run', 'version', 'workspace authority category diff'} <= commands:
            raise ValueError('discovery is missing essential executable commands')
        for revision, identity in [('1', 'urn:woobe:manifest:steps:1'), ('2', 'urn:woobe:manifest:resources:2')]:
            schema = invoke(['schema', '--command', 'manifest validate', '--kind', 'document', '--manifest-version', revision])['data']
            if schema['$id'] != identity:
                raise ValueError('packaged manifest schema has incorrect identity')
        fixture = json.loads((pathlib.Path(__file__).resolve().parent.parent / 'testdata/package/shared/complete.json').read_text())['files']
        bundle = pathlib.Path(directory) / 'complete-package'
        bundle.mkdir()
        for relative, content in fixture.items():
            destination = bundle / relative
            destination.parent.mkdir(parents=True, exist_ok=True)
            destination.write_bytes(content.encode('utf-8'))
        validated = invoke(['package', 'validate', str(bundle)])['data']
        if not validated['valid'] or validated['entrypoint']['kind'] != 'Network':
            raise ValueError('complete Package fixture was not validated offline')
        inventory = [{'path': relative, 'size_bytes': len(content.encode('utf-8')), 'sha256': hashlib.sha256(content.encode('utf-8')).hexdigest()} for relative, content in sorted(fixture.items())]
        tuples = [[item['path'], item['size_bytes'], item['sha256']] for item in inventory]
        digest = hashlib.sha256(b'woobe-package@1.0\n' + json.dumps(tuples, separators=(',', ':'), ensure_ascii=False).encode('utf-8')).hexdigest()
        if validated['artifact_digest'] != digest or validated['inventory'] != inventory:
            raise ValueError('native Package inventory differs from Python canonical digest')
        lock = {'format': 'woobe-package', 'schema_version': '1.0', 'kind': 'PackageLock', 'algorithm': 'sha256', 'inventory': inventory, 'artifact_digest': digest}
        lock_bytes = json.dumps(lock).encode()
        (bundle / 'woobe.lock.json').write_bytes(lock_bytes)
        invoke(['package', 'validate', str(bundle), '--locked'])
        (bundle / 'agent.yaml').write_bytes(fixture['agent.yaml'].encode() + b'\n')
        invoke(['package', 'validate', str(bundle), '--locked'], code=2)
        if (bundle / 'woobe.lock.json').read_bytes() != lock_bytes:
            raise ValueError('Package validation rewrote an inventory lock')
        body = json.dumps({'schema_version': '2', 'project_id': 'p', 'resources': [{'key': 'a', 'kind': 'Agent', 'action': 'update', 'resource_id': 'a', 'spec': {'name': 'Smoke'}}]})
        validation = invoke(['manifest', 'validate', '--file', '-', '--project', 'p'], body)['data']
        if validation['valid'] is not True or validation['source_schema_version'] != '2' or len(validation['manifest_hash']) != 64:
            raise ValueError('local manifest validation has incorrect semantics')
        category_body = json.dumps({'schema_version': '2', 'workspace_id': 'w', 'resources': [{'key': 'reader', 'kind': 'AuthorityCategory', 'action': 'create', 'spec': {'name': 'reader', 'permissions': ['agent:read']}}]})
        category = invoke(['manifest', 'compile', '--file', '-'], category_body)['data']
        if category['document']['workspace_id'] != 'w' or category['executed'] is not False:
            raise ValueError('packaged category compilation lost its Workspace')
        invoke(['manifest', 'validate', '--file', '-'], '{', code=2)
        protected_plan = json.dumps({'schema_version': '2', 'project_id': 'p', 'resources': [{'key': 'tool', 'kind': 'Tool', 'action': 'create', 'spec': {'config': {'headers': [{'name': 'Authorization', 'value': {'$secret_ref': 'tool-auth'}}]}}}]})
        protected = invoke(['manifest', 'compile', '--file', '-'], protected_plan)['data']['document']
        if protected['steps'][0]['body']['config']['headers'][0]['value'] != {'$secret_ref': 'tool-auth'}:
            raise ValueError('packaged compilation lost protected credential reference')
        invoke(['project', 'tool', 'mcp', 'bulk', '--project', 'p', '--file', '-', '--dry-run'], json.dumps({'tool_id': '11111111-1111-4111-8111-111111111111', 'mode': 'invalid'}), code=2)
        for permission_mode in ['allow', 'deny', 'review']:
            invoke(['project', 'tool', 'mcp', 'bulk', '--project', 'p', '--file', '-', '--dry-run'], json.dumps({'tool_id': '11111111-1111-4111-8111-111111111111', 'mode': permission_mode}))

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
                    schema['paths']['/identity/workspaces/{workspace_id}/authority-categories/{category_id}'] = {'get': {'operationId': 'get_category'}}
                    self.send_response(200)
                    self.send_header('Content-Type', 'application/json')
                    self.end_headers()
                    self.wfile.write(json.dumps(schema).encode())
                    return
                if route.path == '/identity/workspaces/w/authority-categories/c':
                    revision = int(query['revision'][0])
                    payload = {'data': {'id': 'c', 'workspace_id': 'w', 'revision': revision, 'name': 'reader', 'description': None, 'scope': 'project', 'permissions': ['agent:read'] if revision == 1 else ['agent:write', 'agent:read'], 'conditions': {}, 'catalog_revision': 'fixture'}}
                    self.send_response(200)
                    self.send_header('Content-Type', 'application/json')
                    self.end_headers()
                    self.wfile.write(json.dumps(payload).encode())
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
            category_diff = invoke(['workspace', 'authority', 'category', 'diff', 'c', '--workspace', 'w', '--from-revision', '1', '--to-revision', '2', '--api-url', origin])['data']
            if category_diff['grants_changed'] is not False or category_diff['executed'] is not False or [row['field'] for row in category_diff['changes']] != ['permissions']:
                raise ValueError('packaged category diff has incorrect authority semantics')
        finally:
            server.shutdown()
            server.server_close()
            thread.join(timeout=5)
    return {'commit': commit, 'version': manifest['version'], 'os': system, 'arch': arch,
            'archive': artifact['name'], 'success': True, 'checks': ['identity', 'discovery', 'schemas', 'complete-package-offline', 'package-python-inventory-parity', 'package-lock-tamper', 'manifest', 'invalid-input', 'body-pagination', 'partial-collection', 'request-direction', 'uuid-date-time', 'validation-before-write', 'path-query-validation', 'validated-pagination', 'workspace-category-manifest', 'category-revision-diff', 'protected-reference-compile', 'canonical-mcp-permissions'],
            'backend_acceptance': 'not_evaluated', 'credential_provider_acceptance': 'not_evaluated'}


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--commit', required=True)
    parser.add_argument('--report', required=True)
    options = parser.parse_args()
    report = smoke(pathlib.Path('dist'), options.commit)
    pathlib.Path(options.report).write_text(json.dumps(report, indent=2) + '\n')
    print(json.dumps(report))

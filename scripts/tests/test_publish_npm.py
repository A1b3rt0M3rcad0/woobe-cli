import base64
import hashlib
import json
import os
import pathlib
import sys
import tempfile
import unittest
from unittest.mock import patch

sys.path.insert(0, str(pathlib.Path(__file__).resolve().parents[1]))
import publish_npm


class NpmPublicationTests(unittest.TestCase):
    def setUp(self):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        old = os.getcwd()
        os.chdir(directory.name)
        self.addCleanup(os.chdir, old)
        pathlib.Path('dist').mkdir()
        self.version = '1.2.3'
        self.commit = 'a' * 40
        archive = pathlib.Path('dist/woobe-cli-1.2.3.tgz')
        archive.write_bytes(b'tested npm candidate')
        metadata = pathlib.Path('dist/artifacts.json')
        metadata.write_text(json.dumps(dict(version=self.version, commit=self.commit,
                                            npm=dict(name=archive.name))))
        self.digest = hashlib.sha256(metadata.read_bytes()).hexdigest()
        self.integrity = 'sha512-' + base64.b64encode(hashlib.sha512(archive.read_bytes()).digest()).decode()

    def test_identical_retry_never_republishes(self):
        with patch.object(publish_npm, 'view', return_value=self.integrity), patch.object(publish_npm.subprocess, 'run') as run:
            publish_npm.publish(self.commit, self.version, self.digest)
            self.assertEqual(run.call_count, 1)

    def test_conflicting_existing_version_is_refused(self):
        with patch.object(publish_npm, 'view', return_value='sha512-other'), patch.object(publish_npm.subprocess, 'run') as run:
            with self.assertRaisesRegex(ValueError, 'different bytes'):
                publish_npm.publish(self.commit, self.version, self.digest)
            self.assertEqual(run.call_count, 1)

    def test_manifest_tampering_is_refused_before_registry_or_publish(self):
        with patch.object(publish_npm, 'view') as view, patch.object(publish_npm.subprocess, 'run') as run:
            with self.assertRaisesRegex(ValueError, 'manifest'):
                publish_npm.publish(self.commit, self.version, '0' * 64)
            view.assert_not_called()
            run.assert_not_called()

    def test_stable_prerelease_and_recovery_tags(self):
        for version, latest, expected in [('1.2.3', '1.2.2', 'latest'), ('1.2.3', '1.3.0', 'historical'), ('1.2.3-rc.1', '1.2.2', 'next')]:
            with self.subTest(version=version, latest=latest):
                metadata = pathlib.Path('dist/artifacts.json')
                data = json.loads(metadata.read_text())
                data['version'] = version
                metadata.write_text(json.dumps(data))
                digest = hashlib.sha256(metadata.read_bytes()).hexdigest()
                with patch.object(publish_npm, 'view', side_effect=[None, latest, self.integrity]), patch.object(publish_npm.subprocess, 'run') as run:
                    publish_npm.publish(self.commit, version, digest)
                    command = run.call_args_list[1].args[0]
                    self.assertTrue(pathlib.Path(command[2]).is_absolute())
                    self.assertEqual(command[command.index('--tag') + 1], expected)
                    self.assertIn('--provenance', command)

    def test_registry_auth_failure_is_not_missing_package(self):
        result = type('Result', (), dict(returncode=1, stdout='{"error":{"code":"E401"}}', stderr=''))()
        with patch.object(publish_npm.subprocess, 'run', return_value=result):
            with self.assertRaises(RuntimeError):
                publish_npm.view('woobe-cli@1.2.3', 'dist.integrity')

    def test_registry_404_allows_first_publication(self):
        result = type('Result', (), dict(returncode=1, stdout='{"error":{"code":"E404"}}', stderr=''))()
        with patch.object(publish_npm.subprocess, 'run', return_value=result):
            self.assertIsNone(publish_npm.view('woobe-cli@1.2.3', 'dist.integrity'))

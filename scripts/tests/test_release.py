import base64
import hashlib
import os
import pathlib
import sys
import tempfile
import unittest
from unittest.mock import patch

SCRIPTS = pathlib.Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SCRIPTS))
import publish_release as release
import release_identity as identity


class ReleaseTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.package = pathlib.Path(self.directory.name) / 'woobe-cli-0.1.0.tgz'
        self.package.write_bytes(b'exact validated npm bytes')
        self.files = {self.package.name: self.package}
        self.commit = 'a' * 40
        self.repo = 'owner/woobe-cli'
        self.env = patch.dict(os.environ, {'NODE_AUTH_TOKEN': '', 'ACTIONS_ID_TOKEN_REQUEST_URL': ''})
        self.env.start()
        self.addCleanup(self.env.stop)

    def packument(self, payload=None):
        integrity = 'sha512-' + base64.b64encode(hashlib.sha512(payload or self.package.read_bytes()).digest()).decode()
        return {'versions': {'0.1.0': {'dist': {'integrity': integrity,
                'tarball': 'https://registry.npmjs.org/woobe-cli/-/woobe-cli-0.1.0.tgz'}}}}

    def github(self, path, missing=False):
        if path == f'repos/{self.repo}':
            return {'permissions': {'push': True}}
        if '/releases?' in path:
            return []
        return None

    def test_semver_order_and_advancement(self):
        versions = ['0.9.0', '1.0.0-alpha', '1.0.0-alpha.2', '1.0.0-alpha.10',
                    '1.0.0-beta', '1.0.0-rc.1', '1.0.0', '1.0.1']
        self.assertEqual(sorted(reversed(versions), key=release.semver_key), versions)
        release.assert_advances('1.0.1', ['invalid', '1.0.0'])
        with self.assertRaises(ValueError):
            release.assert_advances('1.0.0-rc.1', ['1.0.0'])

    def test_unknown_http_error_is_never_absence(self):
        from urllib.error import HTTPError
        for code in [401, 403, 429, 500]:
            with self.subTest(code=code), patch.object(release.urllib.request, 'urlopen',
                    side_effect=HTTPError('https://example.com', code, 'error', {}, None)):
                with self.assertRaises(RuntimeError):
                    release.fetch('https://example.com', missing=True)

    def test_registry_404_is_absence(self):
        from urllib.error import HTTPError
        with patch.object(release.urllib.request, 'urlopen',
                side_effect=HTTPError('https://example.com', 404, 'not found', {}, None)):
            self.assertIsNone(release.fetch('https://example.com', missing=True))

    def test_first_publish_without_token_fails_before_writes(self):
        with patch.object(release, 'github', side_effect=self.github), \
             patch.object(release, 'registry', return_value=None), patch.object(release, 'run') as command:
            with self.assertRaisesRegex(ValueError, 'first npm publication'):
                release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            command.assert_not_called()

    def test_conflicting_npm_bytes_block_all_publication(self):
        with patch.object(release, 'github', side_effect=self.github), \
             patch.object(release, 'registry', return_value=self.packument(b'different')), \
             patch.object(release, 'run') as command:
            with self.assertRaisesRegex(ValueError, 'npm version'):
                release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            command.assert_not_called()

    def test_remote_tarball_is_checked_even_when_integrity_matches(self):
        with patch.object(release, 'fetch', return_value=b'different'):
            with self.assertRaisesRegex(ValueError, 'tarball differs'):
                release.verify_npm(self.packument(), '0.1.0', self.package)

    def test_conflicting_tag_blocks_all_publication(self):
        def github(path, missing=False):
            if '/git/ref/' in path:
                return {'object': {'type': 'commit', 'sha': 'b' * 40}}
            return self.github(path, missing)
        with patch.object(release, 'github', side_effect=github), patch.object(release, 'run') as command:
            with self.assertRaisesRegex(ValueError, 'different commit'):
                release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            command.assert_not_called()

    def test_asset_conflict_blocks_npm_publication(self):
        def github(path, missing=False):
            if '/git/ref/' in path:
                return {'object': {'type': 'commit', 'sha': self.commit}}
            if '/releases/tags/' in path:
                return {'tag_name': 'v0.1.0', 'draft': True, 'prerelease': False,
                        'assets': [{'name': self.package.name, 'size': 1}]}
            return self.github(path, missing)
        with patch.object(release, 'github', side_effect=github), patch.object(release, 'run') as command:
            with self.assertRaisesRegex(ValueError, 'asset differs'):
                release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            command.assert_not_called()

    def test_failed_preflight_never_writes(self):
        with patch.object(release, 'github', side_effect=RuntimeError('HTTP 403')), \
             patch.object(release, 'run') as command:
            with self.assertRaises(RuntimeError):
                release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            command.assert_not_called()

    def test_resume_matching_npm_does_not_republish_or_move_dist_tag(self):
        complete = {'tag_name': 'v0.1.0', 'draft': False, 'prerelease': False}
        state = (complete, set(self.files), True)
        with patch.object(release, 'preflight', return_value=state), \
             patch.object(release, 'verify_npm', return_value=True), \
             patch.object(release, 'verify_assets', return_value=set(self.files)), \
             patch.object(release, 'github', return_value=complete), \
             patch.object(release, 'registry', return_value=self.packument()), \
             patch.object(release, 'run') as command:
            release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            commands = [call.args for call in command.call_args_list]
            self.assertEqual(len(commands), 1)  # installed binary smoke only
            self.assertIn('scripts/smoke_npm.py', commands[0])

    def test_resume_partial_draft_uploads_only_missing_assets(self):
        draft = {'tag_name': 'v0.1.0', 'draft': True}
        published = dict(draft, draft=False)
        with patch.object(release, 'preflight', side_effect=[(draft, set(), True), (published, set(self.files), True)]), \
             patch.object(release, 'verify_npm', return_value=True), \
             patch.object(release, 'verify_assets', return_value=set(self.files)), \
             patch.object(release, 'github', return_value=draft), \
             patch.object(release, 'registry', return_value=self.packument()), \
             patch.object(release, 'run') as command:
            release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            commands = [call.args for call in command.call_args_list]
            self.assertTrue(any(args[:3] == ('gh', 'release', 'upload') for args in commands))
            self.assertFalse(any('--clobber' in args for args in commands))
            self.assertFalse(any(args[:2] == ('npm', 'publish') for args in commands))

    def test_fresh_release_publishes_npm_before_creating_github_record(self):
        draft = {'tag_name': 'v0.1.0', 'draft': True}
        published = dict(draft, draft=False)
        with patch.object(release, 'preflight', side_effect=[(None, set(), False), (published, set(self.files), True)]), \
             patch.object(release, 'verify_npm', return_value=True), \
             patch.object(release, 'verify_assets', return_value=set(self.files)), \
             patch.object(release, 'github', return_value=draft), \
             patch.object(release, 'registry', return_value=self.packument()), \
             patch.object(release, 'run') as command, \
             patch.object(release.pathlib.Path, 'write_text'):
            release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            commands = [call.args for call in command.call_args_list]
            self.assertEqual(commands[0][:2], ('npm', 'publish'))
            self.assertIn('--provenance', commands[0])
            create = next(args for args in commands if args[:3] == ('gh', 'release', 'create'))
            self.assertIn('--draft', create)
            self.assertEqual(create[create.index('--target') + 1], self.commit)

    def test_npm_publish_failure_cannot_create_github_tag_or_release(self):
        with patch.object(release, 'preflight', return_value=(None, set(), False)), \
             patch.object(release, 'run', side_effect=RuntimeError('npm permission denied')) as command:
            with self.assertRaises(RuntimeError):
                release.publish(self.repo, '0.1.0', self.commit, self.files, self.package)
            self.assertEqual(command.call_count, 1)
            self.assertEqual(command.call_args.args[:2], ('npm', 'publish'))

    def test_preflight_only_never_publishes(self):
        with patch.object(release, 'preflight', return_value=(None, set(), False)), \
             patch.object(release, 'run') as command:
            release.publish(self.repo, '0.1.0', self.commit, self.files, self.package, preflight_only=True)
            command.assert_not_called()

    def test_draft_without_tag_recovers_only_the_exact_commit(self):
        def github(path, missing=False):
            if '/releases/tags/' in path:
                return {'tag_name': 'v0.1.0', 'draft': True, 'prerelease': False,
                        'target_commitish': self.commit, 'assets': []}
            return self.github(path, missing)
        with patch.object(release, 'github', side_effect=github), \
             patch.object(release, 'registry', return_value=self.packument()), \
             patch.object(release, 'fetch', return_value=self.package.read_bytes()):
            draft, assets, npm_exists = release.preflight(self.repo, '0.1.0', self.commit, self.files, self.package)
            self.assertTrue(draft['draft'])
            self.assertTrue(npm_exists)
            with self.assertRaisesRegex(ValueError, 'exact draft revision'):
                release.preflight(self.repo, '0.1.0', 'b' * 40, self.files, self.package)


class IdentityTests(unittest.TestCase):
    def test_recovery_uses_old_tag_when_master_has_advanced(self):
        tagged = 'a' * 40
        def git(*args):
            return {('rev-parse', 'refs/tags/v0.1.0^{commit}'): tagged,
                    ('show', tagged + ':VERSION'): '0.1.0',
                    ('show', tagged + ':packages/woobe-cli/package.json'): '{"version":"0.1.0"}'}[args]
        with patch.object(identity, 'git', side_effect=git), \
             patch.object(identity.subprocess, 'run') as command:
            command.return_value.returncode = 0
            self.assertEqual(identity.resolve('0.1.0'), tagged)

    def test_full_revision_and_tag_must_match(self):
        with self.assertRaises(ValueError):
            identity.resolve('0.1.0', 'master')
        with patch.object(identity, 'git', return_value='a' * 40), \
             patch.object(identity.subprocess, 'run') as command:
            command.return_value.returncode = 0
            with self.assertRaisesRegex(ValueError, 'different commit'):
                identity.resolve('0.1.0', 'b' * 40)

    def test_unreviewed_version_is_rejected(self):
        def git(*args):
            if args[0] == 'rev-parse':
                return 'a' * 40
            if args[1].endswith(':VERSION'):
                return '0.1.0'
            return '{"version":"0.1.0"}'
        with patch.object(identity, 'git', side_effect=git), \
             patch.object(identity.subprocess, 'run') as command:
            command.return_value.returncode = 1
            with self.assertRaisesRegex(ValueError, 'review VERSION'):
                identity.resolve('0.2.0')


if __name__ == '__main__':
    unittest.main()

# Woobe CLI skill

The native CLI already provides offline `woobe skill install/status/uninstall`.
This optional separate npm package provides the same payload and compatible receipts.

Separate npm package `woobe-cli-skill`, executable `woobe-skill`. Node.js 22+;
Woobe CLI 0.23.0+. No dependencies, lifecycle scripts, keys or API requests.

```sh
npx --yes --package=woobe-cli-skill woobe-skill install --agent codex
npx --yes --package=woobe-cli-skill woobe-skill install --agent claude
```

See [the complete installation, update and publisher guide](../../docs/AGENT_SKILL.md).
The distribution README uses that complete guide. This directory is a private
source template; release packaging injects the calculated version/source commit
and creates the public tarball. Do not publish this directory directly.

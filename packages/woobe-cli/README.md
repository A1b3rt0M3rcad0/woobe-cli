# Woobe CLI

The native Woobe client, packaged for Linux, macOS and Windows on x64 and arm64.
Node.js 22 or newer is required for the npm launcher. Go is not required.

After the first npm publication, run without a global installation:

```sh
npx --package=woobe-cli woobe version
npx --package=woobe-cli woobe help
```

Or install once:

```sh
npm install --global woobe-cli
woobe version
```

Pin an exact version for repeatable automation:

```sh
npx --package=woobe-cli@0.1.0 woobe version
```

The package contains all six native executables. Installation needs no lifecycle
scripts or external binary downloads; `npm install --ignore-scripts` works.
Arguments, stdin, output and exit codes pass directly to the native client.

## Connect your Workspace

In Woobe, open **Workspace settings → CLI Keys**, create a Control Key and save
the secret before closing the view. Import it from stdin, then select your context:

```sh
woobe auth credential import --name workspace-cli --stdin < workspace.key
woobe context create woobe --api-url https://YOUR_WOOBE_API
woobe context credential attach woobe --credential workspace-cli
woobe context use woobe
woobe context set --workspace WORKSPACE_ID --project PROJECT_ID
woobe project agent list
```

The API URL is the backend URL. Project access comes from the key's explicit
grants. Runtime execution uses a separate Runtime Key.

For a terminal without Node, download the native archive from
[GitHub Releases](https://github.com/A1b3rt0M3rcad0/woobe-cli/releases).
See [installation and versioning](https://github.com/A1b3rt0M3rcad0/woobe-cli/blob/master/docs/INSTALLATION.md)
and the included `USAGE.md` for commands and credential handling.

The client is under active development. See
[implementation status](https://github.com/A1b3rt0M3rcad0/woobe-cli/blob/master/docs/STATUS.md)
for functional coverage and known limitations. Preparing this package does not
publish it to npm or create a GitHub Release.

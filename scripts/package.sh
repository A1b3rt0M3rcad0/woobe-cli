#!/usr/bin/env bash
set -euo pipefail
version="${1:?version required}"
[[ "$version" =~ ^v?[0-9][A-Za-z0-9._-]*$ ]] || { echo "invalid version" >&2; exit 2; }
mkdir -p dist
task_schema="$(mktemp)"
trap 'rm -f "$task_schema"; if [[ -n "${task_dir:-}" ]]; then rm -rf "$task_dir"; fi' EXIT
go run ./cmd/woobe schema --command "manifest validate" --kind document > "$task_schema"
for target_os in linux darwin windows; do
  for target_arch in amd64 arm64; do
    task_dir="$(mktemp -d)"
    executable=woobe
    if [[ "$target_os" == windows ]]; then executable=woobe.exe; fi
    CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" go build -trimpath -ldflags="-s -w -X github.com/A1b3rt0M3rcad0/woobe-cli/internal/cli.Version=$version" -o "$task_dir/$executable" ./cmd/woobe
    cp README.md "$task_dir/README.md"
 cp docs/USAGE.md "$task_dir/USAGE.md"
 python3 - "$task_schema" "$task_dir/manifest.schema.json" <<'PYSCHEMA'
import json,sys
with open(sys.argv[1]) as f: schema=json.load(f)["data"]
with open(sys.argv[2],"w") as f: json.dump(schema,f,indent=2);f.write("\n")
PYSCHEMA
    archive="woobe_${version}_${target_os}_${target_arch}"
    if [[ "$target_os" == windows ]]; then
      task_archive="$(pwd)/dist/$archive.zip"
      (cd "$task_dir" && zip -q "$task_archive" "$executable" README.md USAGE.md manifest.schema.json)
    else
      tar -czf "dist/$archive.tar.gz" -C "$task_dir" "$executable" README.md USAGE.md manifest.schema.json
    fi
    rm -rf "$task_dir"
  done
done
(cd dist && sha256sum "woobe_${version}_"*.tar.gz "woobe_${version}_"*.zip > SHA256SUMS)
python3 scripts/artifact_manifest.py "$version" "$(git rev-parse HEAD)" "$(go version)"

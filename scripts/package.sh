#!/usr/bin/env bash
set -euo pipefail
version="${1:?version required}"
mkdir -p dist
for target_os in linux darwin windows; do
  for target_arch in amd64 arm64; do
    task_dir="$(mktemp -d)"
    executable=woobe
    if [[ "$target_os" == windows ]]; then executable=woobe.exe; fi
    CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" go build -trimpath -ldflags="-s -w -X github.com/A1b3rt0M3rcad0/woobe-cli/internal/cli.Version=$version" -o "$task_dir/$executable" ./cmd/woobe
    cp README.md "$task_dir/README.md"
    archive="woobe_${version}_${target_os}_${target_arch}"
    if [[ "$target_os" == windows ]]; then
      task_archive="$(pwd)/dist/$archive.zip"
      (cd "$task_dir" && zip -q "$task_archive" "$executable" README.md)
    else
      tar -czf "dist/$archive.tar.gz" -C "$task_dir" "$executable" README.md
    fi
    rm -rf "$task_dir"
  done
done
(cd dist && sha256sum *.tar.gz *.zip > SHA256SUMS)

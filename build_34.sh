#!/usr/bin/env bash
set -euo pipefail

repo_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$repo_dir"

go -C qmk_macro_mapper run . ../bare_layout_34.json ../macros.json ../34_built_layout.json
qmk flash -kb eli_3x5 ./34_built_layout.json

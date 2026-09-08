#!/usr/bin/env bash
set -euo pipefail

repo_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$repo_dir"

go -C qmk_macro_mapper run . ../bare_layout.json ../macros.json ../dactyl_built_layout.json
qmk flash ./dactyl_built_layout.json


#!/usr/bin/env bash

set -euo pipefail

GO=$(which go)
ID=$(id -u)

sudo "$GO" run . bootstrap -uid "$ID" -home-dir "$HOME"

#!/usr/bin/env bash
set -e

REPO="chameleon-db/chameleondb"
ARTIFACT_DIR=".artifacts"

echo "🔍 Checking dependencies..."
command -v curl >/dev/null || { echo "curl not found"; exit 1; }
command -v tar  >/dev/null || { echo "tar not found"; exit 1; }

mkdir -p "$ARTIFACT_DIR"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

echo "📦 Downloading artifacts ($OS/$ARCH)"

URL="https://github.com/$REPO/releases/${VERSION}/download/chameleondb-${OS}-${ARCH}.tar.gz"
curl -L https://github.com/$REPO/releases/latest/download/chameleon-$OS-$ARCH.tar.gz | tar -xz -C "$ARTIFACT_DIR"

echo "✅ Artifacts ready"

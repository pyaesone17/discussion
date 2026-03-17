#!/bin/bash
set -euo pipefail

REPO_URL="https://github.com/kamranahmedse/developer-roadmap.git"
CLONE_DIR=$(mktemp -d)
DEST_DIR="$(cd "$(dirname "$0")" && pwd)/references"

# Add or remove roadmap names here to control which ones get copied
ROADMAPS=(
  "backend"
  "frontend"
  "ios"
  "android"
  "golang"
  "javascript"
  "typescript"
  "react"
  "software-architect"
)

echo "Cloning developer-roadmap into $CLONE_DIR..."
git clone --depth 1 "$REPO_URL" "$CLONE_DIR"

SRC_DIR="$CLONE_DIR/src/data/roadmaps"

if [ ! -d "$SRC_DIR" ]; then
  echo "ERROR: $SRC_DIR not found in cloned repo"
  rm -rf "$CLONE_DIR"
  exit 1
fi

mkdir -p "$DEST_DIR"

count=0
for name in "${ROADMAPS[@]}"; do
  content_dir="$SRC_DIR/$name/content"
  if [ -d "$content_dir" ]; then
    echo "Copying $name/content -> references/$name/content"
    mkdir -p "$DEST_DIR/$name"
    cp -r "$content_dir" "$DEST_DIR/$name/"
    count=$((count + 1))
  else
    echo "WARNING: $name/content not found, skipping"
  fi
done

echo ""
echo "Done. Copied $count roadmaps to $DEST_DIR"

rm -rf "$CLONE_DIR"
echo "Cleaned up temp clone."

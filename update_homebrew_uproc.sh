#!/usr/bin/env bash

set -euo pipefail

RELEASE_REPO="uproc-io/uproc-cli"
TAP_REPO="uproc-io/homebrew-uproc"
FORMULA_PATH="Formula/uproc.rb"
DEFAULT_BRANCH="main"

TAG=""
PUSH_MODE="pr"
BRANCH=""

usage() {
  cat <<'EOF'
Usage: ./update_homebrew_uproc.sh [--tag vX.Y.Z] [--push-pr|--push-direct]

Updates uproc Homebrew formula with version + sha256 values from a GitHub release.

Options:
  --tag <tag>       Release tag to use (default: latest release)
  --push-pr         Push a branch and create PR (default)
  --push-direct     Commit and push directly to main
  --help            Show this help message
EOF
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tag)
      if [[ $# -lt 2 ]]; then
        echo "Missing value for --tag" >&2
        exit 1
      fi
      TAG="$2"
      shift 2
      ;;
    --push-pr)
      PUSH_MODE="pr"
      shift
      ;;
    --push-direct)
      PUSH_MODE="direct"
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

for cmd in gh git mktemp python3; do
  if ! need_cmd "$cmd"; then
    echo "[ERROR] Required command not found: $cmd" >&2
    exit 1
  fi
done

if [[ -z "$TAG" ]]; then
  TAG="$(gh release view --repo "$RELEASE_REPO" --json tagName --jq '.tagName')"
fi

if [[ -z "$TAG" ]]; then
  echo "[ERROR] Could not resolve release tag" >&2
  exit 1
fi

VERSION="${TAG#v}"
CHECKSUMS_URL="https://github.com/${RELEASE_REPO}/releases/download/${TAG}/checksums.txt"

TMP_DIR="$(mktemp -d)"
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

echo "[INFO] Using release tag: $TAG"
echo "[INFO] Downloading checksums from: $CHECKSUMS_URL"
gh api --method GET "repos/${RELEASE_REPO}/releases/tags/${TAG}" >/dev/null
curl -fsSL "$CHECKSUMS_URL" -o "$TMP_DIR/checksums.txt"

checksum_for() {
  local asset_name="$1"
  local checksum
  checksum="$(awk -v f="$asset_name" '$2 == f {print $1}' "$TMP_DIR/checksums.txt")"
  if [[ -z "$checksum" ]]; then
    echo "[ERROR] Missing checksum for asset: $asset_name" >&2
    exit 1
  fi
  printf '%s' "$checksum"
}

DARWIN_AMD64="uproc.cli_${VERSION}_darwin_amd64.tar.gz"
DARWIN_ARM64="uproc.cli_${VERSION}_darwin_arm64.tar.gz"
LINUX_AMD64="uproc.cli_${VERSION}_linux_amd64.tar.gz"
LINUX_ARM64="uproc.cli_${VERSION}_linux_arm64.tar.gz"

SHA_DARWIN_AMD64="$(checksum_for "$DARWIN_AMD64")"
SHA_DARWIN_ARM64="$(checksum_for "$DARWIN_ARM64")"
SHA_LINUX_AMD64="$(checksum_for "$LINUX_AMD64")"
SHA_LINUX_ARM64="$(checksum_for "$LINUX_ARM64")"

echo "[INFO] Cloning tap repo: $TAP_REPO"
git clone "git@github.com:${TAP_REPO}.git" "$TMP_DIR/tap"

FORMULA_FILE="$TMP_DIR/tap/$FORMULA_PATH"
if [[ ! -f "$FORMULA_FILE" ]]; then
  echo "[ERROR] Formula file not found: $FORMULA_PATH" >&2
  exit 1
fi

python3 - "$FORMULA_FILE" "$VERSION" "$SHA_DARWIN_ARM64" "$SHA_DARWIN_AMD64" "$SHA_LINUX_AMD64" "$SHA_LINUX_ARM64" <<'PY'
import re
import sys
from pathlib import Path

path = Path(sys.argv[1])
version = sys.argv[2]
sha_darwin_arm64 = sys.argv[3]
sha_darwin_amd64 = sys.argv[4]
sha_linux_amd64 = sys.argv[5]
sha_linux_arm64 = sys.argv[6]

content = path.read_text()

content = re.sub(r'(?m)^\s*version\s+"[^"]+"', f'  version "{version}"', content)

def replace_sha(block_pattern: str, new_sha: str, text: str) -> str:
  m = re.search(block_pattern, text, re.DOTALL)
  if not m:
    raise SystemExit(f"Missing formula block pattern: {block_pattern}")
  block = m.group(0)
  block_new = re.sub(r'(?m)^\s*sha256\s+"[^"]+"', f'    sha256 "{new_sha}"', block)
  return text.replace(block, block_new)

content = replace_sha(r'if\s+OS\.mac\?\s+&&\s+Hardware::CPU\.arm\?.*?sha256\s+"[^"]+"', sha_darwin_arm64, content)
content = replace_sha(r'elsif\s+OS\.mac\?\s+&&\s+Hardware::CPU\.intel\?.*?sha256\s+"[^"]+"', sha_darwin_amd64, content)
content = replace_sha(r'elsif\s+OS\.linux\?\s+&&\s+Hardware::CPU\.intel\?.*?sha256\s+"[^"]+"', sha_linux_amd64, content)
content = replace_sha(r'elsif\s+OS\.linux\?\s+&&\s+Hardware::CPU\.arm\?.*?sha256\s+"[^"]+"', sha_linux_arm64, content)

path.write_text(content)
PY

pushd "$TMP_DIR/tap" >/dev/null

if git diff --quiet -- "$FORMULA_PATH"; then
  echo "[INFO] Formula already up to date: $FORMULA_PATH"
  exit 0
fi

git add "$FORMULA_PATH"
git commit -m "chore(formula): update uproc to ${TAG}"

if [[ "$PUSH_MODE" == "direct" ]]; then
  echo "[INFO] Pushing directly to ${DEFAULT_BRANCH}"
  git push origin "HEAD:${DEFAULT_BRANCH}"
  echo "[OK] Formula updated directly on ${TAP_REPO}:${DEFAULT_BRANCH}"
else
  BRANCH="chore/update-uproc-${TAG}"
  git push origin "HEAD:${BRANCH}"
  gh pr create \
    --repo "$TAP_REPO" \
    --base "$DEFAULT_BRANCH" \
    --head "$BRANCH" \
    --title "chore(formula): update uproc to ${TAG}" \
    --body "Update Formula/uproc.rb to ${TAG} using checksums from ${RELEASE_REPO}."
  echo "[OK] PR created for ${TAP_REPO}"
fi

popd >/dev/null

#!/usr/bin/env bash
#
# Validates terrafmt against real provider repos (the corpus) in two ways:
#
#  1. golden comparison: runs `terrafmt diff` built from HEAD and from BASE_REF
#     (default origin/main) over the corpus and diffs their output. Any
#     difference is either an intended new capability (review it) or a regression.
#
#  2. idempotency: copies the corpus, runs `terrafmt fmt` over it, then
#     `terrafmt diff --check` — which must come back clean: formatting a second
#     time must never produce further changes.
#
# Usage:
#   CORPUS_REPOS="/path/to/terraform-provider-azurerm /path/to/terraform-provider-aws" ./scripts/check-against-providers.sh
set -euo pipefail

CORPUS_REPOS="${CORPUS_REPOS:-$HOME/hashi/azure/azurerm $HOME/hashi/aws/terraform-provider-aws}"
BASE_REF="${BASE_REF:-origin/main}"

tmp="$(mktemp -d)"
trap 'git worktree remove --force "$tmp/base" >/dev/null 2>&1 || true; rm -rf "$tmp"' EXIT

echo "==> building terrafmt from HEAD..."
go build -o "$tmp/terrafmt-head" .

echo "==> building terrafmt from $BASE_REF..."
git worktree add --force "$tmp/base" "$BASE_REF" >/dev/null
(cd "$tmp/base" && go build -o "$tmp/terrafmt-base" .)

fail=0

for repo in $CORPUS_REPOS; do
  if [ ! -d "$repo" ]; then
    echo "!!> skipping missing corpus repo $repo"
    continue
  fi
  name="$(basename "$repo")"

  # dir:pattern:flags — website markdown is plain, go test files need --fmtcompat
  for spec in "website:*.markdown:" "internal:*_test.go:-f"; do
    dir="${spec%%:*}"
    rest="${spec#*:}"
    pattern="${rest%%:*}"
    flags="${rest##*:}"
    [ -d "$repo/$dir" ] || continue

    echo "==> golden: $name/$dir ($pattern ${flags:--})"
    # shellcheck disable=SC2086
    (cd "$repo" && "$tmp/terrafmt-base" diff "$dir" -p "$pattern" -u -v $flags >"$tmp/$name-$dir-base.out" 2>"$tmp/$name-$dir-base.err") || true
    # shellcheck disable=SC2086
    (cd "$repo" && "$tmp/terrafmt-head" diff "$dir" -p "$pattern" -u -v $flags >"$tmp/$name-$dir-head.out" 2>"$tmp/$name-$dir-head.err") || true

    for stream in out err; do
      if ! diff -u "$tmp/$name-$dir-base.$stream" "$tmp/$name-$dir-head.$stream" >"$tmp/$name-$dir.$stream.diff"; then
        echo "!!> std$stream differs from $BASE_REF for $name/$dir:"
        head -60 "$tmp/$name-$dir.$stream.diff"
        fail=1
      fi
    done
  done

  # idempotency: fmt a copy of the corpus, then diff --check must be clean
  echo "==> idempotency: $name (fmt, then diff --check)"
  copy="$tmp/$name-copy"
  mkdir -p "$copy/website" "$copy/internal"
  [ -d "$repo/website" ] && rsync -a --prune-empty-dirs --include='*/' --include='*.markdown' --exclude='*' "$repo/website/" "$copy/website/"
  [ -d "$repo/internal" ] && rsync -a --prune-empty-dirs --include='*/' --include='*_test.go' --exclude='*' "$repo/internal/" "$copy/internal/"

  [ -d "$copy/website" ] && "$tmp/terrafmt-head" fmt "$copy/website" -p '*.markdown' >/dev/null
  [ -d "$copy/internal" ] && "$tmp/terrafmt-head" fmt "$copy/internal" -p '*_test.go' -f >/dev/null

  if [ -d "$copy/website" ] && ! "$tmp/terrafmt-head" diff "$copy/website" -p '*.markdown' -c -q -u >"$tmp/$name-idem-website.out" 2>&1; then
    echo "!!> fmt is not idempotent (or errored) for $name/website:"
    head -40 "$tmp/$name-idem-website.out"
    fail=1
  fi
  if [ -d "$copy/internal" ] && ! "$tmp/terrafmt-head" diff "$copy/internal" -p '*_test.go' -f -c -q -u >"$tmp/$name-idem-internal.out" 2>&1; then
    echo "!!> fmt is not idempotent (or errored) for $name/internal:"
    head -40 "$tmp/$name-idem-internal.out"
    fail=1
  fi
done

if [ "$fail" -ne 0 ]; then
  echo "!!> provider check FAILED"
  exit 1
fi
echo "==> provider check passed"

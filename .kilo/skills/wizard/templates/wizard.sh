#!/usr/bin/env bash
# Safe interactive wizard template. Replace only the STAGES section.

set -euo pipefail

TOTAL_STAGES=0
STAGE_INDEX=0
ENV_FILE="${ENV_FILE:-.env}"
WRITTEN_KEYS=()

say() { printf '  %s\n' "$1"; }
warn() { printf '  WARNING: %s\n' "$1" >&2; }

confirm() {
  local reply=''
  printf '  %s [y/N] ' "$1"
  read -r reply || true
  [[ "$reply" =~ ^[Yy]$ ]]
}

stage() {
  STAGE_INDEX=$((STAGE_INDEX + 1))
  printf '\n== Stage %s/%s: %s ==\n' "$STAGE_INDEX" "$TOTAL_STAGES" "$1"
}

pause() {
  printf '  %s ' "${1:-Press Enter to continue.}"
  read -r _ || true
}

ask() {
  local variable="$1" prompt="$2" value=''
  printf '  %s ' "$prompt"
  read -r value
  printf -v "$variable" '%s' "$value"
}

ask_secret() {
  local variable="$1" prompt="$2" value=''
  printf '  %s ' "$prompt"
  read -rs value
  printf '\n'
  printf -v "$variable" '%s' "$value"
}

open_url() {
  local url="$1"
  say "URL: $url"
  confirm 'Open this URL in your browser now?' || return 0
  if command -v wslview >/dev/null 2>&1; then wslview "$url"
  elif command -v explorer.exe >/dev/null 2>&1; then explorer.exe "$url"
  elif command -v xdg-open >/dev/null 2>&1; then xdg-open "$url"
  elif command -v open >/dev/null 2>&1; then open "$url"
  else warn 'No supported browser opener found; open the URL manually.'
  fi
}

write_env() {
  local key="$1" value="$2" directory temporary
  [[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || { warn "Invalid environment key: $key"; return 1; }
  [[ "$value" != *$'\n'* && "$value" != *$'\r'* ]] || { warn "Refusing a multiline value for $key"; return 1; }
  say "Target: $ENV_FILE"
  say "Key: $key (value is hidden)"
  confirm "Create or replace $key in $ENV_FILE?" || return 0
  directory=$(dirname "$ENV_FILE")
  [[ -d "$directory" ]] || { warn "Directory does not exist: $directory"; return 1; }
  umask 077
  temporary=$(mktemp "${ENV_FILE}.tmp.XXXXXX")
  if [[ -f "$ENV_FILE" ]]; then
    awk -v key="$key" 'index($0, key "=") != 1 { print }' "$ENV_FILE" > "$temporary"
  fi
  printf '%s=%s\n' "$key" "$value" >> "$temporary"
  chmod 600 "$temporary"
  mv "$temporary" "$ENV_FILE"
  WRITTEN_KEYS+=("$key")
  say "Updated $key in $ENV_FILE."
}

finish() {
  printf '\nWizard complete.\n'
  if (( ${#WRITTEN_KEYS[@]} )); then
    say "Updated keys in $ENV_FILE: ${WRITTEN_KEYS[*]}"
  fi
}

# STAGES: replace this example. Every consequential action needs confirm nearby.
TOTAL_STAGES=1

stage 'Review this generated wizard'
say 'Replace this example with the authorized procedure stages.'
pause 'Press Enter after reviewing the script.'

finish
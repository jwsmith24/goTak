#!/usr/bin/env bash
# User-run human-in-the-loop reproduction template.
# Copy only after preview and authorization, customize observations, then run in a terminal.
# Never capture passwords, tokens, cookies, authorization headers, personal data, or secrets.

set -euo pipefail

step() {
  printf '\n>>> %s\n' "$1"
  read -r -p '    [Press Enter when complete] ' _
}

capture() {
  local variable="$1" question="$2" answer
  printf '\n>>> %s\n' "$question"
  read -r -p '    > ' answer
  printf -v "$variable" '%s' "$answer"
}

# Replace this example block with the previewed reproduction steps.
step 'Perform the authorized action that triggers the symptom.'
capture REPRODUCED 'Was the exact symptom reproduced? (yes/no)'
capture OBSERVATION 'Describe the non-sensitive observation briefly:'

printf '\n--- Captured observations ---\n'
printf 'REPRODUCED=%s\n' "$REPRODUCED"
printf 'OBSERVATION=%s\n' "$OBSERVATION"
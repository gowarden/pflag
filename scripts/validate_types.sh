#!/usr/bin/env bash

set -euo pipefail

cmd_funcs=(
  "|Var"
  "|"
)

fs_funcs=(
  "Get|"
  "MustGet|"
  "|Var"
  "|"
)

function join_by {
  local d=${1-} f=${2-}
  if shift 2; then
    printf %s "$f" "${@/#/$d}"
  fi
}

get_types() {
  local -a types
  mapfile -t types < <(awk '/^type\s*[^ ]+\s*interface\s*\{/{print $2}' flag.go)
  join_by "|" "${types[@]}"
}

iface_assert_types="$(get_types)"

validate_interface_asserts() {
  local file="$1"

  local pattern
  pattern="^var _ ($iface_assert_types) = \(\*([a-zA-Z0-9_]+)\)\(nil\)$"

  local line
  local type=""
  local line_no=0
  rc=0
  while IFS="" read -r line || [[ -n "$line" ]]; do
    ((line_no++))
    if [[ $line =~ ^type\ ([a-zA-Z0-9_]+)\  ]]; then
      type="${BASH_REMATCH[1]}"
    fi

    if [[ "$line" =~ $pattern ]]; then
      if [[ -z $type ]]; then
        printf "Fatal: %s:%d: type declaration must come before interface assertion\n" "$file" "$line_no"
        rc=1
      fi

      if [[ $type != "${BASH_REMATCH[2]}" ]]; then
        printf "Error: %s:%d: interface assertion for '%q' was type '%q' and does not match type '%q'\n" "$file" "$line_no" "${BASH_REMATCH[1]}" "${BASH_REMATCH[2]}" "$type"
        rc=1
      fi
    fi
  done <"$file"
  return $rc
}

validate_funcs() {
  local file="$1"

  pattern="var _ Typed = \(\*([a-zA-Z0-9]+)Value\)\(nil\)"
  rc=0
  while read -r type; do
    if [[ $type =~ $pattern ]]; then
      fn_type="${BASH_REMATCH[1]^}"
      if [[ $fn_type == Ip* ]]; then
        fn_type="IP${fn_type:2}"
      fi

      for req_fn in "${fs_funcs[@]}"; do
        expected_fn="${req_fn//\|/$fn_type}"
        pattern="func\s*\([a-z]+\s+\*FlagSet\)\s+$expected_fn\("
        if ! grep -qE "$pattern" "$file"; then
          printf "%s: Could not find function definition: func (f *FlagSet) %s\n" "$file" "$expected_fn"
          rc=1
        fi
      done

      for req_fn in "${cmd_funcs[@]}"; do
        expected_fn="${req_fn//\|/$fn_type}"
        pattern="func\s+$expected_fn\("
        if ! grep -qE "$pattern" "$file"; then
          printf "%s: Could not find function definition: func %s\n" "$file" "$expected_fn"
          rc=1
        fi
      done
    fi
  done < <(grep -oE "$pattern" "$file")

  return $rc
}

main() {
  rc=0

  for file in "$@"; do
    validate_interface_asserts "$file" || rc=1
    validate_funcs "$file" || rc=1
  done

  return $rc
}

main "$@"

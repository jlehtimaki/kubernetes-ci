#! /usr/bin/env bash

set -Eeuo pipefail
set -x

# Convert actions to multiple actions
opts=""
IFS=","
while (( "$#" )); do
    if [[ $1 == "--actions" ]]; then
      read -a tmp_opts <<< $2
      for o in "${tmp_opts[@]}"; do
        opts="$opts --actions $o"
      done
      break
    else
      opts="$opts $1"
      shift
    fi
done

# Run the script
if [ "$0" = "${BASH_SOURCE[*]}" ] ; then
  >&2 echo -E "\nRunning Kubernetes-CI...\n"
fi
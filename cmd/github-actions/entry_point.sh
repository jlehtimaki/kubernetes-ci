#! /usr/bin/env bash

# Convert actions to multiple actions
opts=""
IFS=","
parse_args() {
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
  # set remaining positional arguments (if any) in their proper place
  eval set -- "$opts"

  echo "${opts/ /}"
  return 0
}

# Generates client.
kubernetes-ci() {
  /bin/kubernetes-ci $opts
}

args=("$@")

if [ "$0" = "${BASH_SOURCE[*]}" ] ; then
  echo $opts
  >&2 echo -E "\nRunning Kubernetes-CI...\n"
  kubernetes-ci "${args[@]}"
fi
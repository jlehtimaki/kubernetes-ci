#! /usr/bin/env bash

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
echo $opts

COMMAND="${opts}"

/bin/kubernetes-ci $COMMAND
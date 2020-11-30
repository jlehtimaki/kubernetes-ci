#! /usr/bin/env bash

set -Eeuo pipefail
set -x

# Filter out arguments that are not available to this action
# args:
#   $@: Arguments to be filtered
parse_args() {
  local opts=""
  while (( "$#" )); do
    case "$1" in
      --actions)
        opts="$opts --actions=$2"
        shift
        ;;
      --type)
        opts="$opts --type=$2"
        shift 2
        ;;
      --k8s_ca)
        opts="$opts --k8s_ca=$2"
        shift
        ;;
      --k8s_token)
        opts="$opts --k8s_token=$2"
        shift 2
        ;;
      --k8s_user)
        opts="$opts --k8s_user=$2"
        shift
        ;;
      --k8s_server)
        opts="$opts --k8s_server=$2"
        shift 2
        ;;
      --assume_role)
        opts="$opts --assume_role=$2"
        shift 2
        ;;
      --kubectl_version)
        opts="$opts --kubectl_version=$2"
        shift 2
        ;;
      --cluster_name)
        opts="$opts --cluster_name=$2"
        shift
        ;;
      --manifest_dir)
        opts="$opts --manifest_dir=$2"
        shift
        ;;
      --kubernetes_namespace)
        opts="$opts --kubernetes_namespace=$2"
        shift
        ;;
      --region)
        opts="$opts --region=$2"
        shift
        ;;
      --kustomize)
        opts="$opts --kustomize=$2"
        shift
        ;;
      --image.version)
        opts="$opts --image.version=$2"
        shift
        ;;
      --image.name)
        opts="$opts --image.name=$2"
        shift
        ;;
      --rolloutCheck)
        opts="$opts --rolloutCheck=$2"
        shift
        ;;
      --rolloutTimeout)
        opts="$opts --rolloutTimeout=$2"
        shift
        ;;
      --googleProjectID)
        opts="$opts --googleProjectID=$2"
        shift
        ;;
      --googleSA)
        opts="$opts --googleSA=$2"
        shift
        ;;
      --) # end argument parsing
        shift
        break
        ;;
      -*) # unsupported flags
        >&2 echo "ERROR: Unsupported flag: '$1'"
        exit 1
        ;;
      *) # positional arguments
        shift  # ignore
        ;;
    esac
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
  >&2 echo -E "\nRunning Kubernetes-CI...\n"
  kubernetes-ci "${args[@]}"
fi
#! /usr/bin/env bash

if [ "$0" = "${BASH_SOURCE[*]}" ] ; then
  >&2 echo -E "\nRunning Kubernetes-CI...\n"
  /bin/kubernetes-ci $*
fi
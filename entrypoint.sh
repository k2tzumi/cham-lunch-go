#!/bin/sh

set -e

if [ -d ${HOME}/envs/scripts ]; then
  for file in ${HOME}/envs/scripts/*.sh ; do
    /bin/sh ${file}
  done
fi

exec "$@"

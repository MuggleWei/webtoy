#!/bin/bash

origin_dir=$(readlink -f "$(dirname "$0")")
cd $origin_dir/docker_volume/cert

rm -rf rootCA.* server.*

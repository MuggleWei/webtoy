#!/bin/bash

origin_dir=$(readlink -f "$(dirname "$0")")
cd $origin_dir

echo "---------------------------"
echo "- generate CA"
echo "---------------------------"

$origin_dir/nginx/gen_ca.sh

echo "---------------------------"
echo "- run docker"
echo "---------------------------"

sudo docker compose up -d

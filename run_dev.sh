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

sudo docker compose -f docker-compose.dev.yml up -d

echo "---------------------------"
echo "- build services"
echo "---------------------------"

mkdir -p $origin_dir/build

backends=(
	webtoy_gate
	webtoy_auth
	webtoy_captcha
)
for serv in ${backends[@]}
do
	echo "build service ${serv}"
	rm -rf $origin_dir/build/${serv}
	mkdir -p $origin_dir/build/${serv}
	cd $origin_dir/backend/${serv}
	go build -o $origin_dir/build/${serv}/${serv}
	cp -r ./config $origin_dir/build/${serv}/
done

frontends=(
	webtoy-front
)
for serv in ${frontends[@]}
do
	echo "build service ${serv}"
	rm -rf $origin_dir/build/${serv}
	mkdir -p $origin_dir/build/${serv}
	cd $origin_dir/frontend/${serv}
	npm install
	npm run build
	mv build $origin_dir/build/${serv}/${serv}
done

echo "---------------------------"
echo "- run services"
echo "---------------------------"

host=10.0.2.15

echo "run webtoy_gate"
cd $origin_dir/build/webtoy_gate
nohup ./webtoy_gate \
	--host=$host \
	--port=8080 \
	--sr.service.host=$host \
	--sr.service.port=8180 \
	--sr.service.id=webtoy-gate-0 \
	&

echo "run webtoy_auth"
cd $origin_dir/build/webtoy_auth
nohup ./webtoy_auth \
	--host=$host \
	--port=8180 \
	--sr.service.host=$host \
	--sr.service.port=8180 \
	--sr.service.id=webtoy-auth-0 \
	&

echo "run webtoy_captcha"
cd $origin_dir/build/webtoy_captcha
nohup ./webtoy_captcha \
	--host=$host \
	--port=8280 \
	--sr.service.host=$host \
	--sr.service.port=8280 \
	--sr.service.id=webtoy-captcha-0 \
	&

echo "run webtoy-front"
cd $origin_dir/build/webtoy-front
npm install -S serve
nohup node node_modules/serve/build/main.js -s webtoy-front -l 3000 &

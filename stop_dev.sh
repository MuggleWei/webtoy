#!/bin/bash

origin_dir=$(readlink -f "$(dirname "$0")")
cd $origin_dir

echo "---------------------------"
echo "- stop services"
echo "---------------------------"

backends=(
	webtoy_gate
	webtoy_auth
	webtoy_captcha
)
for serv in ${backends[@]}
do
	echo "stop service ${serv}"
	killall ${serv}
done

frontends=(
	webtoy-front
)
for serv in ${frontends[@]}
do
	echo "stop service ${serv}"
	arr=$(ps aux | grep "${serv}" | grep -v "stop.sh" | awk '{print $2}')
	num=${#arr[@]}
	echo "found $num ${serv}"
	if [[ $num -gt 0 ]]; then
		for i in "${arr[@]}"
		do
			if test -z "$i"
			then
				echo "nothing to kill"
			else
				echo "kill -9 $i"
				kill -9 $i
			fi
		done
	fi
done

echo "---------------------------"
echo "- stop docker"
echo "---------------------------"

#sudo docker compose -f docker-compose.dev.yml stop
#sudo docker compose -f docker-compose.dev.yml rm -f
sudo docker compose -f docker-compose.dev.yml down

echo "---------------------------"
echo "- clean CA"
echo "---------------------------"

$origin_dir/nginx/clean_ca.sh

#! /usr/bin/env bash

if [ -z "$1" ]; then
  echo I need an api key as argument
  exit 1
fi

res=$(curl -X "POST" "https://panel.degrensrp.be/api/state/schedule/update" -H "X-Api-Key: $1")
if [ $res != "false" ]; then
  echo "Update already scheduled, ending here"
  exit 0
fi

echo "No other update scheduled. Going to wait 1min"
sleep 60s

cd ~/panel
git fetch
git pull --rebase

docker login 

docker image pull registry-gitlab.pieter557.dscloud.me/degrens-21/panel/webui && docker image pull registry-gitlab.pieter557.dscloud.me/degrens-21/panel/server

docker compose down --remove-orphans
docker compose up -d


#!/bin/sh

docker run -d -t -i --env-file .env_list --name githubforkrefresh docker.io/library/app

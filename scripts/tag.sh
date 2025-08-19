#!/bin/bash


#    - /^[0-9]{8}\.[0-9]{6}-[a-zA-Z0-9]+$/
NAME_TAG=`date +"%Y%m%d.%H%M%S"`-`git rev-parse HEAD | head -c 10`

echo "$NAME_TAG"

git tag  "$NAME_TAG"
git push origin "$NAME_TAG"
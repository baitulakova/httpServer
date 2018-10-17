#!/bin/bash
file=$1
if [ ! -f "${file}" ]; then echo "Error - ${file} does not exist"; exit 1; fi
url=http://127.0.0.1:8080/upload
curl -Ffile=@$file $url


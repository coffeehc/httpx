#!/usr/bin/env bash

echo "build commons"
go build  github.com/coffeehc/commons
echo "build commons/httpcommons/client"
go build  github.com/coffeehc/commons/httpcommons/httpclient

#!/bin/bash

wrk -t12 -c200 -d5s http://127.0.0.1:9999/api/runner/beiluo/debug/_getApiInfo?router=/hello&method=GET
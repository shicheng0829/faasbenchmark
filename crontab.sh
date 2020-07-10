#!/usr/bin/env bash
cd /root/faasbenchmark
./faasbenchmark run aliyun all > log
python3 result.py
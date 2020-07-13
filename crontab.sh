#!/usr/bin/env bash
source /root/.bashrc
cd /root/faasbenchmark
./faasbenchmark run aliyun all > log
python3 result.py
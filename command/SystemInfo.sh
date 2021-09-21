#! /usr/bin/env bash
echo "系统总体负载统计"
top -b -n 1 | head -n 5
echo ""
echo "应用负载统计"
docker stats --no-stream
echo ""
echo "磁盘使用量统计"
df -hl | head -n 4
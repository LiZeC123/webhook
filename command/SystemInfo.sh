#! /usr/bin/env bash
function showProcessStat() {
    name=$1
    pattern=$2
    stat=$(pgrep $pattern)
    
    if [ ${stat} ]; then
        stat="On"
    else
        stat="Off"
    fi
    printf "%8s:%8s\n" ${name} ${stat}
}



echo "系统总体负载统计"
top -b -n 1 | head -n 5

echo ""
echo "应用负载统计"
docker stats --no-stream

echo ""
echo "磁盘使用量统计"
df -hl | head -n 4

echo ""
echo "进程状态监控"
showProcessStat OneDrive onedrive
showProcessStat Docker dockerd

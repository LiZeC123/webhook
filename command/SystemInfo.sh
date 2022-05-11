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
echo "<pre style='word-wrap: break-word; white-space: pre-wrap;'>"
top -b -n 1 | head -n 5
echo "</pre>"

echo "应用负载统计"
echo "<pre style='word-wrap: break-word; white-space: pre-wrap;'>"
docker stats --no-stream
echo "<hr />"
docker ps
echo "</pre>"

echo "磁盘使用量统计"
echo "<pre style='word-wrap: break-word; white-space: pre-wrap;'>"
df -hl | head -n 4
echo "</pre>"

echo "进程状态监控"
echo "<pre style='word-wrap: break-word; white-space: pre-wrap;'>"
showProcessStat OneDrive onedrive
showProcessStat Docker dockerd
echo "</pre>"
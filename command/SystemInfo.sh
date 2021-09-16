#! /usr/bin/env bash

top -b -n 1 | head -n 5
echo ""
docker stats --no-stream
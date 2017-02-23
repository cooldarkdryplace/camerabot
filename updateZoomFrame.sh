#!/bin/bash

/usr/bin/raspistill -n -q 98 --ISO 100 -ex auto -sh 100 -roi 0.593,0.52,0.04,0.04 -hf -vf -t 1 -o /tmp/zoomedFrame.png

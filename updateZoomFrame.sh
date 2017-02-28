#!/bin/bash

/usr/bin/raspistill -w 320 -h 240 \
-mm spot -n -q 90 --ISO 100 -ex auto \
-sa -100 -sh 100 -co 100 \
-roi 0.593,0.52,0.04,0.04 \
-hf -vf -t 1 \
-o /tmp/zoomedFrame.jpg

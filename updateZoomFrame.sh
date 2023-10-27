#!/bin/bash

/usr/bin/libcamera-still -w 320 -h 240 \
-roi 0.593,0.52,0.04,0.04 \
-o /tmp/zoomedFrame.jpg

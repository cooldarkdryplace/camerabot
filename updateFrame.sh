#!/bin/bash

/usr/bin/raspistill -w 1024 -h 768 -n -q 80 -ex auto -hf -vf -t 1 -o /tmp/frame.jpg

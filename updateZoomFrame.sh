#!/bin/bash

/usr/bin/raspistill -n -q 98 -ex auto -roi 0.5,0.5,0.25,0.25 -hf -vf -t 1 -o /tmp/frame.png


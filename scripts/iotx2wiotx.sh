#!/bin/bash
if [ "$#" -ne 1 ]; then
    echo "usage: $0 amount (IOTX)"
    exit
fi
addr=io15qr5fzpxsnp7garl4m7k355rafzqn8grrm0grz
ioctl contract invoke bytecode $addr d0e30db0 ${1}

#!/bin/bash
if [ "$#" -ne 1 ]; then
    echo "usage: $0 amount (IOTX)"
    exit
fi
function decimal2hex() {
  echo "ibase=A;obase=16;$1" | bc | xargs -0 -I h printf "%65s" h | sed "s/ /0/g"
}
addr=io15qr5fzpxsnp7garl4m7k355rafzqn8grrm0grz
amount=`decimal2hex $(echo $1 '*' 1000000000000000000 | bc | sed "s/\..*//g")`
ioctl contract invoke bytecode $addr 2e1a7d4d${amount}

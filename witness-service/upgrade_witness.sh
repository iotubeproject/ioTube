#!/bin/bash

git pull || { echo 'git pull failed' ; exit 1; }
bash start_witness.sh

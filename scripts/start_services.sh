#!/bin/sh

# Start Nginx in the background
nginx

OPT=/opt/guardicore
export PYTHONPATH=$OPT

# Start gc-inventory in the background
$OPT/gc-inventory &

# Start gc-onboarding in the foreground
$OPT/gc-onboarding & 


sleep 30
python3 -m reveal 

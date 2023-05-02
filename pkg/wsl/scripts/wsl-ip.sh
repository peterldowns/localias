#!/usr/bin/env bash
# Prints the IP address being used by the WSL container.
# From testing some manual edits to the windows hosts file,
#
#   172.20.166.118 explicit.test
#   127.0.0.1 homev4.test
#   ::1 homev6.test
#
# both explicit.test and homev4.test work but homev6.test does not.
# Therefore, hostctl should be able to continue to use 127.0.0.1 instead
# of needing to pass in a special IP address. I'll leave this script here
# just in case.
#
# Copied from HRX_
# https://superuser.com/a/1749524
ip addr show eth0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}'

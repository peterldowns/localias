#!/usr/bin/env bash

# Looks for the microsoft tag in the hostname.  If it's present, the host is
# assumed to be WSL.
#
#   Linux downs-windows 5.4.72-microsoft-standard-WSL2 #1 SMP Wed Oct 28 23:40:43 UTC 2020 x86_64 GNU/Linux
#   Linux Eve 4.4.0-18362-Microsoft #476-Microsoft Fri Nov 01 16:53:00 PST 2019 x86_64 x86_64 x86_64 GNU/Linux
#
# see https://stackoverflow.com/a/59765344
uname -a | grep -i 'microsoft' || echo ""
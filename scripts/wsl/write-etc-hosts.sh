#!/usr/bin/env bash
#
# Reads from the file $1 (or stdin if not passed) and writes it to
# the windows /etc/hosts file.
etchosts='$env:windir\System32\drivers\etc\hosts'
etchosts='./example.hosts'
tmpfile=$(mktemp /tmp/localias-XXXXXXX)
cat "${1:-/dev/stdin}" > $tmpfile

powershell.exe "Set-Content -Path $etchosts -Value (Get-Content -Path $tmpfile -Raw) -Force"

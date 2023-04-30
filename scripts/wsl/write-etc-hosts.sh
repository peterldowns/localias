#!/usr/bin/env bash
#
# Reads file contents from stdin and writes it to the windows
# /etc/hosts file.
tmpfile=$(mktemp /tmp/localias-XXXXXXX)
# Send stdin from this script to that tmpfile
cat - > "$tmpfile"
cat "$tmpfile"
wintmpfile=$(wslpath -w "$tmpfile")
# Using '$env:windir' intentionally to get the windows path
# to the base windows installation. We don't want bash to
# expand this variable.
# shellcheck disable=SC2016
etchosts='$env:windir\System32\drivers\etc\hosts'
echo powershell.exe ./script.ps1 "$wintmpfile" "$etchosts" sudo
powershell.exe ./script.ps1 "$wintmpfile" "$etchosts" sudo
echo "rm $tmpfile"
rm "$tmpfile"

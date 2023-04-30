#!/usr/bin/env bash
tmpfile=$(mktemp /tmp/localias-XXXXXX)
wintmpfile=$(wslpath -w $tmpfile)
echo "   tmpfile = $tmpfile"
echo "wintmpfile = $wintmpfile"
# Send stdin from this script to that tmpfile
cat - > $tmpfile

outfile=$1
winoutfile=$(wslpath -w $1)
echo "   outfile = $outfile"
echo "winoutfile = $winoutfile"
powershell.exe ./script.ps1 $wintmpfile $winoutfile $2
rm $tmpfile

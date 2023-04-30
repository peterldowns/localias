#!/usr/bin/env bash
# Using '$env:windir' intentionally to get the windows path
# to the base windows installation. We don't want bash to
# expand this variable.
# shellcheck disable=SC2016
winetchosts=$(powershell.exe -c 'Write-Host (Resolve-Path $env:windir\System32\drivers\etc\hosts)')
cat "$(wslpath -u "$winetchosts")"
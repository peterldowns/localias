#!/usr/bin/env bash
#
# powershell.exe will consume STDIN, which it inherits from this script,
# automatically.  when this script is embedded inside the localias binary, it is
# executed by being passed via STDIN to the bash program. If we don't use
# </dev/null inside the powershell.exe subshell, powershell.exe will "eat" the
# rest of this script and bash will have nothing else to execute.
#
# This uses '$env:windir' intentionally to get the windows path to the base
# windows installation. We don't want bash to expand this variable, it has
# special meaning to Powershell.
# shellcheck disable=SC2016
winetchosts=$(powershell.exe -c 'Write-Host (Resolve-Path $env:windir\System32\drivers\etc\hosts)' </dev/null)
cat "$(wslpath -u "$winetchosts")"
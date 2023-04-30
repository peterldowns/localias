param(
  [Parameter()]
  [String]$infile,
  [String]$outfile,
  [String]$sudo
)
function Sudo {
    Start-Process powershell.exe -WindowStyle hidden -Wait -Verb RunAs -ArgumentList @args
}
function Run {
    Start-Process powershell.exe -WindowStyle hidden -Wait -ArgumentList @args
}
$command="Set-Content -Path $outfile -Value (Get-Content -Path $infile) -Force"
if ($sudo.ToString() -eq "sudo") {
  Write-Host "(Sudo) $command"
  Sudo $command
} else {
  Write-Host "(Run) $command"
  Run $command
}
  


#$args = "Set-Content -Path $outfile -Value (Get-Content -Path $infile) -Force"
#Write-Host $args
#Start-Process powershell.exe -WindowStyle hidden -Wait -ArgumentList $args


# if ($sudo.ToString() -eq "") {
#   Write-Host "using sudo"
#   Sudo powershell.exe 'Set-Content -Path $outfile -Value (Get-Content -Path $infile) -Force'
# } else {
# Run powershell.exe "Set-Content -Path $outfile -Value (Get-Content -Path $infile) -Force"


#sudo powershell.exe 'Add-Content -Path $env:windir\System32\drivers\etc\hosts -Value "`n172.29.224.1`tlocal.test" -Force'
#function Copy
#Write-Output
#Write-Output $filepath
#$y = Get-Content -Path $filepath
#Write-Output $y
### $nl = [Environment]::NewLine
##$x = Get-Content -Path ""
##Write-host $x
##Write-Host "done"

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
  Sudo $command
} else {
  Run $command
}
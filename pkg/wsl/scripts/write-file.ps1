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

$command="Copy-Item -Path $infile -Destination $outfile -Force"

if ($sudo.ToString() -eq "sudo") {
  Sudo $command
} else {
  Run $command
}
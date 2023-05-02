param(
  [Parameter()]
  [String]$infile,
  [String]$sudo
)
function Sudo {
    Start-Process powershell.exe -WindowStyle hidden -Wait -Verb RunAs -ArgumentList @args
}

function Run {
    Start-Process powershell.exe -WindowStyle hidden -Wait -ArgumentList @args
}

$command="Import-Certificate -FilePath '$infile' -CertStoreLocation Cert:\CurrentUser\Root"

if ($sudo.ToString() -eq "sudo") {
  Sudo $command
} else {
  Run $command
}
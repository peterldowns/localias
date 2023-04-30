function sudo {
    Start-Process @args -verb runas
}

sudo powershell.exe 'Add-Content -Path $env:windir\System32\drivers\etc\hosts -Value "`n172.29.224.1`tlocal.test" -Force'

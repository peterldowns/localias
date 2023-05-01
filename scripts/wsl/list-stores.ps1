# Get the current user's certificate stores
$stores = Get-ChildItem -Path "Cert:\CurrentUser" -Recurse -Force | Where-Object { $_.PSIsContainer -and $_.Name -ne "My" -and $_.Name -ne "CA" -and $_.Name -ne "Root" -and $_.Name -ne "Trust" }

# Display the store names
$stores | ForEach-Object { Write-Host $_.Name }

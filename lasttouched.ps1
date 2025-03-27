# Set the path to your folder share
$folderPath = "C:\Path\To\Your\SharedFolder"

# Define the age thresholds
$now = Get-Date
$threshold9Months = $now.AddMonths(-9)
$threshold12Months = $now.AddMonths(-12)

# Get all files recursively and filter by LastAccessTime
$files = Get-ChildItem -Path $folderPath -Recurse -File -ErrorAction SilentlyContinue

$old9Months = $files | Where-Object { $_.LastAccessTime -lt $threshold9Months -and $_.LastAccessTime -ge $threshold12Months }
$old12Months = $files | Where-Object { $_.LastAccessTime -lt $threshold12Months }

# Output to screen
Write-Host "`n=== Files not accessed in 9–12 months ===`n"
$old9Months | Select-Object FullName, LastAccessTime

Write-Host "`n=== Files not accessed in over 12 months ===`n"
$old12Months | Select-Object FullName, LastAccessTime

# Optional: export to CSV
$old9Months | Select-Object FullName, LastAccessTime | Export-Csv -Path ".\old_files_9_months.csv" -NoTypeInformation
$old12Months | Select-Object FullName, LastAccessTime | Export-Csv -Path ".\old_files_12_months.csv" -NoTypeInformation

$bucket = "aoscms"
$downloadRoot = "downloaded_files"
$dateFilter = "2025-02-26"

# Make sure the root download folder exists
if (!(Test-Path -Path $downloadRoot)) {
    New-Item -ItemType Directory -Path $downloadRoot | Out-Null
}

# Get all possible accountname prefixes (top-level folders)
$topLevel = aws s3api list-objects-v2 --bucket $bucket --delimiter '/' `
  | ConvertFrom-Json

foreach ($prefixObj in $topLevel.CommonPrefixes) {
    $accountPrefix = $prefixObj.Prefix  # e.g., "accountname1/"
    $accountName = $accountPrefix.TrimEnd('/')

    # Set sub-path to Letters/Mittera
    $subPrefix = "$accountName/Letters/Mittera/"

    Write-Host "`nSearching under: s3://$bucket/$subPrefix"

    # Paginate results under Letters/Mittera for each accountname
    $objects = aws s3api list-objects-v2 --bucket $bucket --prefix $subPrefix --output json `
        | ConvertFrom-Json

    foreach ($obj in $objects.Contents) {
        $key = $obj.Key
        $lastModified = $obj.LastModified

        if ($key -like "*SYFSDB*" -and $key.ToLower().EndsWith(".pdf")) {
            if ($lastModified.StartsWith($dateFilter)) {
                $relativePath = $key.Substring($subPrefix.Length)
                $localPath = Join-Path "$downloadRoot\$accountName" $relativePath
                $localDir = Split-Path $localPath -Parent

                if (!(Test-Path -Path $localDir)) {
                    New-Item -ItemType Directory -Path $localDir -Force | Out-Null
                }

                Write-Host "Downloading: $key --> $localPath"
                aws s3 cp "s3://$bucket/$key" $localPath
            }
        }
    }
}

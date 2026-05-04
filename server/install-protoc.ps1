$ProgressPreference = 'SilentlyContinue'
$zipPath = "$env:TEMP\protoc.zip"
$extractPath = "$env:TEMP\protoc-extract"

Write-Host "Downloading protoc..."
Invoke-WebRequest -Uri "https://github.com/protocolbuffers/protobuf/releases/download/v25.3/protoc-25.3-win64.zip" -OutFile $zipPath

Write-Host "Extracting..."
Remove-Item -Path $extractPath -Recurse -Force -ErrorAction SilentlyContinue
Expand-Archive -Path $zipPath -DestinationPath $extractPath

# Copy to GOPATH/bin
$goBin = "$env:USERPROFILE\go\bin"
Copy-Item "$extractPath\bin\protoc.exe" "$goBin\protoc.exe" -Force

Write-Host "protoc installed to $goBin\protoc.exe"
Write-Host "Cleaning up..."
Remove-Item $zipPath -Force
Remove-Item $extractPath -Recurse -Force

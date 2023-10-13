$GitHubRepo = "kevincobain2000/gobrew"
$GoBrewRoot = [IO.Path]::Combine($HOME, ".gobrew")
$GoRoot = [IO.Path]::Combine($GoBrewRoot, "current", "go")

# check if env GOBREW_ROOT is set
if ($null -ne [Environment]::GetEnvironmentVariable("GOBREW_ROOT")) {
  $GoBrewRoot = [IO.Path]::Combine($Env:GOBREW_ROOT, ".gobrew")
}
$GoBrewBin = [IO.Path]::Combine($GoBrewRoot, "bin")

$AddUserPaths = @(
  $GoBrewBin
  [IO.Path]::Combine($GoBrewRoot, "current", "bin")
  [IO.Path]::Combine($HOME, "go", "bin")
)

$AddEnvVars = @{
  'GOROOT' = $GoRoot
  'GOPATH' = [IO.Path]::Combine($HOME, "go")
}

# create bin dir if not exists
New-Item -ItemType Directory -Force -Path $GoBrewBin -ErrorAction Stop | Out-Null

##

if ($PSVersionTable.PSVersion.Major -lt 6) {
  $global:IsWindows = ([Environment]::OSVersion.Platform -eq "Win32NT")
}

$GoArch = {
  $archMap = @{
    "amd64" = "x86_64", "i386_64"
    "386"   = "i386", "i686", "x86"    
  }

  $platform = {
    if ($IsWindows) {
      if (![String]::IsNullOrEmpty($Env:PROCESSOR_ARCHITEW6432)) {
        return $Env:PROCESSOR_ARCHITEW6432
      }
      return $Env:PROCESSOR_ARCHITECTURE
    }
    return $(uname -p)
  }.Invoke().ToLower()

  foreach ($k in $archMap.Keys) {
    foreach ($v in $archMap[$k]) {
      if ($platform.Equals($v)) {        
        return $k
      }
    }
  }
  return $platform
}.Invoke()

if ([String]::IsNullOrEmpty($GoArch)) {
  Write-Host "Unsupported CPU architecture" -f Red
  Exit 1
}

$GoOS = {
  $osMap = @{
    "linux" = "mingw64_nt", "mingw32_nt"
  }

  $platform = {
    if ($IsWindows) {
      return "windows"
    }
    return $(uname -s)
  }.Invoke().ToLower()

  foreach ($k in $osMap.Keys) {
    foreach ($v in $osMap[$k]) {
      if ($platform.StartsWith($v)) {        
        return $k
      }
    }
  }
  return $platform  
}.Invoke()

if ([String]::IsNullOrEmpty($GoOS)) {
  Write-Host "Unsupported OS: $([Environment]::OSVersion.Platform)" -f Red
  Exit 1
}

##

$GoBrewBinName = "gobrew"
if ($IsWindows) {
  $GoBrewBinName += ".exe"
}

$fileMask = "gobrew-$($GoOS)-$($GoArch)"

$releases = "https://api.github.com/repos/$GitHubRepo/releases/latest"
$latest = (Invoke-WebRequest -Uri $releases -UseBasicParsing -ErrorAction Stop | ConvertFrom-Json)[0]
$tagName = $latest.tag_name
$latestAssets = $latest.assets
$latestBin = $latestAssets | Where-Object { $_.name -match "^$($fileMask).exe" }
$latestZip = $latestAssets | Where-Object { $_.name -match "^$($fileMask).zip" }
$matchedAsset = @($latestZip, $latestBin).Where({ ![String]::IsNullOrEmpty($_) }, "First")
$downloadUrl = $matchedAsset.browser_download_url
if ([String]::IsNullOrEmpty($downloadUrl)) {
  Write-Host "Unable to find a valid release URL!" -f Magenta
  Write-Host "Unable to match OS: $GoOS, CPU: $GoArch" -f Magenta
  Write-Host "Please open an issue at: https://github.com/$GitHubRepo/issues" -f Magenta 
  Write-Host "Remember to provide your OS, CPU details and powershell version." -f Magenta
  Exit 1
}

Write-Host "Downloading gobrew $tagName from: $downloadUrl ..." -f Cyan
$GoBrewDownloadPath = [IO.Path]::Combine($GoBrewBin, $matchedAsset.name)
Remove-Item -Path $GoBrewDownloadPath -Force -ErrorAction SilentlyContinue
Invoke-WebRequest -OutFile $GoBrewDownloadPath -Uri $downloadUrl -UseBasicParsing -ErrorAction Stop
if ($null -eq $?) {
  Write-Host "Failed to download gobrew from: $downloadUrl" -f Red
  Exit 1
}

Unblock-File -Path $GoBrewDownloadPath -ErrorAction SilentlyContinue

if (![IO.File]::Exists($GoBrewDownloadPath)) {
  Write-Host "Failed to download gobrew to: $GoBrewBinPath" -f Red
  Exit 1
}

if ($GoBrewDownloadPath.EndsWith('.zip') -and ([IO.Directory]::Exists($GoBrewDownloadPath))) {
  Write-Host "Extracting gobrew to: $GoBrewBin" -f Cyan
  Expand-Archive -Path $GoBrewDownloadPath -DestinationPath $GoBrewBin -Force
  if ($null -eq $?) {
    Write-Host "Failed to extract gobrew to: $GoBrewBin" -f Red
    Exit 1
  }
  Remove-Item -Path $GoBrewDownloadPath -Force
}
else {
  $GoBrewBinPath = [IO.Path]::Combine($GoBrewBin, $GoBrewBinName)
  Remove-Item -Path $GoBrewBinPath -Force -ErrorAction SilentlyContinue
  Rename-Item -Path $GoBrewDownloadPath -NewName $GoBrewBinName -Force -ErrorAction Stop
  if ([IO.File]::Exists($GoBrewBinPath)) {
    Write-Host "Gobrew successfully installed to: $GoBrewBinPath" -f Green
    Write-Host "Run 'gobrew' to see available commands" -f Green
  }
  else {
    Write-Host "Failed to install gobrew to: $GoBrewBinPath" -f Red
    Exit 1
  }
}

# Add paths to user PATH
$UserEnvVars = [Environment]::GetEnvironmentVariables('User')
$userPath = $UserEnvVars.Path.Split(';');

[System.Collections.ArrayList]$newPaths = @()
foreach ($path in $AddUserPaths) {
  if ($userPath.Contains($path)) {
    Write-Host "$path is already in user PATH" -f DarkGreen
    continue
  }
  Write-Host "Adding $path to user PATH" -f DarkCyan
  [void]$newPaths.Add($path)
}
$updatedPath = [String]::Join(';', $userPath + $newPaths.ToArray())
[Environment]::SetEnvironmentVariable('Path', $updatedPath, 'User')
$Env:Path = $updatedPath.TrimEnd(';') + ';' + [Environment]::GetEnvironmentVariables('Machine').Path

# Set extra environment variables
foreach ($key in $AddEnvVars.Keys) {
  $value = $AddEnvVars[$key]
  if ($UserEnvVars.Contains($key)) {
    Write-Host "$key is already set" -f DarkGreen
  }
  else {
    Write-Host "Setting $key to $value" -f DarkCyan
    [Environment]::SetEnvironmentVariable($key, $value, 'User')
  }
  Set-Item "Env:$key" -Value $value
}

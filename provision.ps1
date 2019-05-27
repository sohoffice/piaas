# Enable SSH
echo "Enable SSH"
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
echo "Starting SSH"
Start-Service sshd
# OPTIONAL but recommended:
Set-Service -Name sshd -StartupType 'Automatic'
# Confirm the Firewall rule is configured. It should be created automatically by setup.
Get-NetFirewallRule -Name *ssh*
# There should be a firewall rule named "OpenSSH-Server-In-TCP", which should be enabled

# Enable RDP
echo "Enable RDP"
Set-ItemProperty 'HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server\' -Name "fDenyTSConnections" -Value 0
Set-ItemProperty 'HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp\' -Name "UserAuthentication" -Value 1
Enable-NetFirewallRule -DisplayGroup "Remote Desktop"

# Common variables
$downloadDir = $env:TEMP

# Install 7zip
if (!(Test-Path "$downloadDir/7za920.zip")) {
  echo "Install 7zip"
  Invoke-WebRequest -Uri "https://www.7-zip.org/a/7za920.zip" -OutFile "$downloadDir/7za920.zip" -UseBasicParsing
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  [System.IO.Compression.ZipFile]::ExtractToDirectory("$downloadDir/7za920.zip", "C:\Users\IEUser\Bin")

  echo "Add ~/bin to path"
  $userenv = [System.Environment]::GetEnvironmentVariable("Path", "User")
  [System.Environment]::SetEnvironmentVariable("PATH", $userenv + ";C:\Users\IEUser\Bin", "User")
}

if (!(Test-Path "$downloadDir/PortableGit.7z.exe")) {
  echo "Install Git"
  Invoke-WebRequest -Uri "https://github.com/git-for-windows/git/releases/download/v2.21.0.windows.1/PortableGit-2.21.0-64-bit.7z.exe" -OutFile "$downloadDir/PortableGit.7z.exe" -UseBasicParsing
  mkdir C:\Git
  Set-Location -Path C:\Git
  & "C:\Users\IEUser\Bin\7za.exe" x "$downloadDir/PortableGit.7z.exe"

  echo "Add ~/Git/bin to path"
  $userenv = [System.Environment]::GetEnvironmentVariable("Path", "User")
  [System.Environment]::SetEnvironmentVariable("PATH", $userenv + ";C:\Git\Bin", "User")
}

if (!(Test-Path "$downloadDir/mingw-w64.7z")) {
  echo "Install mingw"
  Invoke-WebRequest -Uri "https://nchc.dl.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z" -OutFile "$downloadDir/mingw-w64.7z" -UseBasicParsing
  Set-Location -Path C:\
  & "C:\Users\IEUser\Bin\7za.exe" x "$downloadDir/mingw-w64.7z"

  if (Test-Path "C:\mingw64\bin\mingw32-make.exe") {
    echo "Mingw64 installed"
    cp C:\mingw64\bin\mingw32-make.exe C:\mingw64\bin\make.exe

    echo "Add C:/mingw64/bin to path"
    $userenv = [System.Environment]::GetEnvironmentVariable("Path", "User")
    [System.Environment]::SetEnvironmentVariable("PATH", $userenv + ";C:\mingw64\bin", "User")
  }
}

# Enable WSL
echo "Enable WSL"
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux
echo "WSL enabled, you may have to restart it to continue."

if (!(Test-Path "$downloadDir/Ubuntu.zip")) {
  echo "Download ubuntu 16.04"
  Invoke-WebRequest -Uri https://aka.ms/wsl-ubuntu-1604 -OutFile "$downloadDir/Ubuntu.appx" -UseBasicParsing
  mv "$downloadDir/Ubuntu.appx" "$downloadDir/Ubuntu.zip"
}

if (!(Test-Path "C:\Users\IEUser\Ubuntu")) {
  echo "Extracting Ubuntu"
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  [System.IO.Compression.ZipFile]::ExtractToDirectory("$downloadDir/Ubuntu.zip", "C:\Users\IEUser\Ubuntu")

  echo "Add Ubuntu to path"
  $userenv = [System.Environment]::GetEnvironmentVariable("Path", "User")
  [System.Environment]::SetEnvironmentVariable("PATH", $userenv + ";C:\Users\IEUser\Ubuntu", "User")

  echo "Please initialize ubuntu with the following steps:"
  echo "  cd C:\users\IEUser\Ubuntu"
  echo "  ./ubuntu.exe"
}

# Download and install go
# Extracted from https://gist.github.com/andrewkroh/2c93f8a5953f6093a505
$version = "1.11.10"
$packageName = 'golang'
$url = 'https://storage.googleapis.com/golang/go' + $version + '.windows-amd64.zip'
$goroot = "C:\go"

if (Test-Path "$goroot\bin\go.exe") {
  Write-Host "Go is already installed"
} else {
  echo "Downloading $url"
  $zip = "$downloadDir\golang-$version.zip"
  if (!(Test-Path "$zip")) {
    $downloader = new-object System.Net.WebClient
    $downloader.DownloadFile($url, $zip)
  }

  echo "Extracting $zip to $goroot"
  if (Test-Path "$downloadDir\go") {
    rm -Force -Recurse -Path "$downloadDir\go"
  }
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  [System.IO.Compression.ZipFile]::ExtractToDirectory("$zip", $downloadDir)
  mv "$downloadDir\go" $goroot

  if (!(Test-Path "$goroot\bin\go.exe")) {
    Write-Host "Go is not installed to $goroot !!"
    exit
  }

  echo "Setting GOROOT and PATH for Machine"
  [System.Environment]::SetEnvironmentVariable("GOROOT", "$goroot", "Machine")
  $p = [System.Environment]::GetEnvironmentVariable("PATH", "User")
  $p = "$goroot\bin;$p"
  [System.Environment]::SetEnvironmentVariable("PATH", "$p", "User")
}

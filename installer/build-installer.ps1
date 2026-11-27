# Builds GAN-Setup.exe.
#
#   powershell -ExecutionPolicy Bypass -File installer\build-installer.ps1
#
# It needs gan.exe to be built already (go build -o gan.exe .). It uses IExpress, the
# self-extracting-installer maker that comes with Windows, so nothing has to be installed.
#
# What goes into the installer:
#   gan.exe, libs\, examples\ (with modules\), tutorial\ (if it exists), icons\, README.txt
# Inside the .exe: install.ps1 (the setup window) and payload.zip (all of the above).

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$work = Join-Path $root 'installer\build'
$payload = Join-Path $work 'payload'
$out = Join-Path $root 'installer\GAN-Setup.exe'

if (-not (Test-Path (Join-Path $root 'gan.exe'))) { throw 'gan.exe not found. Build it first: go build -o gan.exe .' }

# start from an empty work folder (it only ever holds files this script made)
if (Test-Path $work) { Remove-Item -Recurse -Force $work }
New-Item -ItemType Directory -Force -Path $payload | Out-Null

Write-Host 'Collecting files...'
Copy-Item (Join-Path $root 'gan.exe') $payload
Copy-Item (Join-Path $root 'libs') (Join-Path $payload 'libs') -Recurse

# examples: the working programs (not the ones that are meant to fail)
$examples = Join-Path $payload 'examples'
New-Item -ItemType Directory -Force -Path $examples | Out-Null
foreach ($n in 'hello', 'var', 'comments', 'math', 'parens', 'decimals', 'strings', 'if', 'elseif', 'logic',
               'while', 'for', 'func', 'scope', 'final', 'lists', 'builtins', 'types', 'imports',
               'breaks', 'tensors', 'tensortools', 'autograd', 'gradients') {
    $f = Join-Path $root "$n.gan"
    if (Test-Path $f) { Copy-Item $f $examples }
}
# ask.gan needs Ollama running on the computer
Copy-Item (Join-Path $root 'live\ask.gan') $examples
Copy-Item (Join-Path $root 'modules') (Join-Path $examples 'modules') -Recurse

if (Test-Path (Join-Path $root 'tutorial')) {
    Copy-Item (Join-Path $root 'tutorial') (Join-Path $payload 'tutorial') -Recurse
    # the lessons only: leave out the checker and its scratch folder
    Remove-Item (Join-Path $payload 'tutorial\.run') -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item (Join-Path $payload 'tutorial\check.py') -Force -ErrorAction SilentlyContinue
}

$icons = Join-Path $payload 'icons'
New-Item -ItemType Directory -Force -Path $icons | Out-Null
Copy-Item (Join-Path $root 'gan-icon.ico') $icons
Copy-Item (Join-Path $root 'gan-icon.svg') $icons
Copy-Item (Join-Path $root 'vscode-gan\gan-seti.vsix') $icons

@'
GAN 0.3 (beta)
==============

Open a NEW terminal window, then:

    gan --version
    gan examples\hello.gan

Write your own program in a file that ends in .gan, for example hi.gan:

    let name = input("Your name? ")
    print("Hello, " + name)

and run it with:   gan hi.gan

Libraries are used with import:

    import "math"
    print(math.square(5))

Folders
    examples\   small programs, one idea each
    libs\       the standard libraries: math, text, lists (written in GAN)
    tutorial\   a step-by-step course (if it came with this version)

To remove GAN: Windows Settings > Apps > Installed apps > GAN programming language > Uninstall.
'@ | Set-Content -Path (Join-Path $payload 'README.txt') -Encoding UTF8

Write-Host 'Packing payload.zip...'
$pkg = Join-Path $work 'package'
New-Item -ItemType Directory -Force -Path $pkg | Out-Null
# a virus scanner or editor can briefly lock freshly copied files, so try a few times
$zipped = $false
for ($try = 1; $try -le 6 -and -not $zipped; $try++) {
    try {
        Compress-Archive -Path (Join-Path $payload '*') -DestinationPath (Join-Path $pkg 'payload.zip') -Force -ErrorAction Stop
        $zipped = $true
    } catch {
        Write-Host "  a file is busy, trying again ($try of 6)..."
        Start-Sleep -Seconds 2
    }
}
if (-not $zipped) { throw 'Could not create payload.zip (files stayed locked).' }
Copy-Item (Join-Path $PSScriptRoot 'install.ps1') $pkg
Copy-Item (Join-Path $root 'gan-icon.ico') $pkg

Write-Host 'Building GAN-Setup.exe with IExpress...'
if (Test-Path $out) { Remove-Item -Force $out }
$sed = Join-Path $work 'gan.sed'
@"
[Version]
Class=IEXPRESS
SEDVersion=3
[Options]
PackagePurpose=InstallApp
ShowInstallProgramWindow=0
HideExtractAnimation=1
UseLongFileName=1
InsideCompressed=0
CAB_FixedSize=0
CAB_ResvCodeSigning=0
RebootMode=N
InstallPrompt=%InstallPrompt%
DisplayLicense=%DisplayLicense%
FinishMessage=%FinishMessage%
TargetName=%TargetName%
FriendlyName=%FriendlyName%
AppLaunched=%AppLaunched%
PostInstallCmd=%PostInstallCmd%
AdminQuietInstCmd=%AdminQuietInstCmd%
UserQuietInstCmd=%UserQuietInstCmd%
SourceFiles=SourceFiles
[Strings]
InstallPrompt=
DisplayLicense=
FinishMessage=
TargetName=$out
FriendlyName=GAN Setup
AppLaunched=powershell.exe -NoProfile -ExecutionPolicy Bypass -STA -WindowStyle Hidden -File install.ps1
PostInstallCmd=<None>
AdminQuietInstCmd=
UserQuietInstCmd=
FILE0="install.ps1"
FILE1="payload.zip"
FILE2="gan-icon.ico"
[SourceFiles]
SourceFiles0=$pkg\
[SourceFiles0]
%FILE0%=
%FILE1%=
%FILE2%=
"@ | Set-Content -Path $sed -Encoding ASCII

$p = Start-Process -FilePath "$env:WINDIR\System32\iexpress.exe" -ArgumentList '/N', '/Q', $sed -Wait -PassThru
if (-not (Test-Path $out)) { throw "IExpress did not produce $out (exit code $($p.ExitCode))" }
$mb = [math]::Round((Get-Item $out).Length / 1MB, 2)
Write-Host "Built $out ($mb MB)"

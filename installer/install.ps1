# GAN installer / uninstaller
#
# Run by GAN-Setup.exe (and by the uninstall entry in Windows "Installed apps").
# It shows a small setup window. Everything is installed for the current user only,
# so no administrator rights are needed.
#
# Options for testing and scripting:
#   -Silent       no window, just do it
#   -Uninstall    remove GAN instead of installing it
#   -InstallDir   where to install (default: %LOCALAPPDATA%\Programs\GAN)
#   -NoPath / -NoAssoc / -NoVSCode   skip that part of the install
#   -RegBase      registry location (default HKCU:\Software). Used by tests so the real
#                 registry is not touched.
#   -SelfTest     build the window without showing it to the user, then stop
#   -Screenshot   with -SelfTest: save a picture of the window to this file

param(
    [string]$InstallDir = (Join-Path $env:LOCALAPPDATA 'Programs\GAN'),
    [switch]$Silent,
    [switch]$Uninstall,
    [switch]$NoPath,
    [switch]$NoAssoc,
    [switch]$NoVSCode,
    [string]$RegBase = 'HKCU:\Software',
    [switch]$SelfTest,
    [string]$Screenshot = '',
    [switch]$FromTemp
)

$ErrorActionPreference = 'Stop'
$AppName = 'GAN'
$AppVersion = '0.3 beta'
$ScriptPath = $MyInvocation.MyCommand.Path
$ScriptDir = Split-Path -Parent $ScriptPath
$UninstallKey = "$RegBase\Microsoft\Windows\CurrentVersion\Uninstall\GAN"

# ---------------------------------------------------------------- PATH helpers
# These two work on plain text, so they can be tested without touching the real PATH.

function Add-PathEntry([string]$path, [string]$dir) {
    $parts = @($path -split ';' | Where-Object { $_.Trim() -ne '' })
    $clean = $dir.TrimEnd('\')
    foreach ($p in $parts) {
        if ($p.Trim().TrimEnd('\') -ieq $clean) { return ($parts -join ';') }  # already there
    }
    return (($parts + $clean) -join ';')
}

function Remove-PathEntry([string]$path, [string]$dir) {
    $clean = $dir.TrimEnd('\')
    $kept = @($path -split ';' | Where-Object { $_.Trim() -ne '' -and $_.Trim().TrimEnd('\') -ine $clean })
    return ($kept -join ';')
}

# Read and write the user PATH straight in the registry. (The usual .NET call can turn a PATH that
# contains %VARIABLES% into plain text, which would break other programs, so we avoid it.)
# $name is "Path" for real use; the tests use a private name so your real PATH is never touched.
function Get-UserPathRaw([string]$name = 'Path') {
    $key = [Microsoft.Win32.Registry]::CurrentUser.OpenSubKey('Environment')
    if (-not $key) { return '' }
    $value = $key.GetValue($name, '', [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)
    $key.Close()
    return [string]$value
}

function Set-UserPathRaw([string]$value, [string]$name = 'Path') {
    $key = [Microsoft.Win32.Registry]::CurrentUser.OpenSubKey('Environment', $true)
    $key.SetValue($name, $value, [Microsoft.Win32.RegistryValueKind]::ExpandString)
    $key.Close()
}

# Tell running programs that the environment changed, so NEW terminals see the new PATH.
function Send-EnvironmentChange {
    try {
        Add-Type -Namespace GanSetup -Name Native -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", SetLastError = true, CharSet = System.Runtime.InteropServices.CharSet.Auto)]
public static extern System.IntPtr SendMessageTimeout(System.IntPtr hWnd, uint Msg, System.UIntPtr wParam, string lParam, uint flags, uint timeout, out System.UIntPtr result);
[System.Runtime.InteropServices.DllImport("shell32.dll")]
public static extern void SHChangeNotify(int eventId, uint flags, System.IntPtr item1, System.IntPtr item2);
'@ -ErrorAction SilentlyContinue
        $r = [UIntPtr]::Zero
        [GanSetup.Native]::SendMessageTimeout([IntPtr]0xffff, 0x1A, [UIntPtr]::Zero, 'Environment', 2, 5000, [ref]$r) | Out-Null
        [GanSetup.Native]::SHChangeNotify(0x08000000, 0, [IntPtr]::Zero, [IntPtr]::Zero)  # refresh file icons
    } catch { }
}

# ---------------------------------------------------------------- install
function Install-Gan([string]$dir, [bool]$addPath, [bool]$assoc, [bool]$vscode, [scriptblock]$report) {
    & $report 5 'Creating the install folder...'
    New-Item -ItemType Directory -Force -Path $dir | Out-Null

    & $report 15 'Copying GAN...'
    $zip = Join-Path $ScriptDir 'payload.zip'
    if (-not (Test-Path $zip)) { throw "payload.zip was not found next to the installer ($zip)" }
    Expand-Archive -Path $zip -DestinationPath $dir -Force
    Copy-Item $ScriptPath (Join-Path $dir 'installer.ps1') -Force   # needed later for uninstalling

    if ($addPath) {
        & $report 55 'Adding GAN to your PATH...'
        Set-UserPathRaw (Add-PathEntry (Get-UserPathRaw) $dir)
    }

    if ($assoc) {
        & $report 65 'Setting up .gan files (icon, double-click to run)...'
        $classes = "$RegBase\Classes"
        New-Item -Path "$classes\.gan" -Force | Out-Null
        Set-ItemProperty -Path "$classes\.gan" -Name '(default)' -Value 'GANFile'
        New-Item -Path "$classes\GANFile" -Force | Out-Null
        Set-ItemProperty -Path "$classes\GANFile" -Name '(default)' -Value 'GAN Program'
        New-Item -Path "$classes\GANFile\DefaultIcon" -Force | Out-Null
        Set-ItemProperty -Path "$classes\GANFile\DefaultIcon" -Name '(default)' -Value (Join-Path $dir 'icons\gan-icon.ico')
        New-Item -Path "$classes\GANFile\shell\open\command" -Force | Out-Null
        # cmd /k keeps the window open after the program ends, so you can read the output
        $cmd = 'cmd.exe /k ""' + (Join-Path $dir 'gan.exe') + '" "%1""'
        Set-ItemProperty -Path "$classes\GANFile\shell\open\command" -Name '(default)' -Value $cmd
    }

    & $report 75 'Registering with Windows (Installed apps)...'
    $uninstallCmd = 'powershell.exe -NoProfile -ExecutionPolicy Bypass -STA -WindowStyle Hidden -File "' + (Join-Path $dir 'installer.ps1') + '" -Uninstall'
    if ($RegBase -ne 'HKCU:\Software') { $uninstallCmd += ' -RegBase "' + $RegBase + '"' }
    New-Item -Path $UninstallKey -Force | Out-Null
    $values = @{
        DisplayName     = 'GAN programming language'
        DisplayVersion  = $AppVersion
        Publisher       = 'GAN'
        InstallLocation = $dir
        DisplayIcon     = (Join-Path $dir 'icons\gan-icon.ico')
        UninstallString = $uninstallCmd
    }
    foreach ($k in $values.Keys) { Set-ItemProperty -Path $UninstallKey -Name $k -Value $values[$k] }
    Set-ItemProperty -Path $UninstallKey -Name 'NoModify' -Value 1 -Type DWord
    Set-ItemProperty -Path $UninstallKey -Name 'NoRepair' -Value 1 -Type DWord

    $extra = ''
    if ($vscode) {
        & $report 85 'Installing the VS Code icon theme...'
        $vsix = Join-Path $dir 'icons\gan-seti.vsix'
        $code = Get-Command code -ErrorAction SilentlyContinue
        if ($code -and (Test-Path $vsix)) {
            try { & $code.Source --install-extension $vsix --force 2>&1 | Out-Null } catch { }
            $extra = ' The VS Code icon theme "Seti + GAN" was installed; choose it with Ctrl+Shift+P > Preferences: File Icon Theme.'
        } else {
            $extra = ' (VS Code was not found, so the icon theme was skipped.)'
        }
    }

    & $report 95 'Finishing...'
    Send-EnvironmentChange
    & $report 100 'Done.'
    return $extra
}

# ---------------------------------------------------------------- uninstall
function Uninstall-Gan([string]$dir, [scriptblock]$report) {
    & $report 10 'Removing GAN from your PATH...'
    $old = Get-UserPathRaw
    if ($old) { Set-UserPathRaw (Remove-PathEntry $old $dir) }

    & $report 40 'Removing .gan file settings...'
    $classes = "$RegBase\Classes"
    # only remove the .gan setting if it is the one we made
    $assoc = Get-ItemProperty -Path "$classes\.gan" -ErrorAction SilentlyContinue
    if ($assoc -and $assoc.'(default)' -eq 'GANFile') { Remove-Item -Path "$classes\.gan" -Recurse -Force }
    if (Test-Path "$classes\GANFile") { Remove-Item -Path "$classes\GANFile" -Recurse -Force }

    & $report 60 'Removing the Windows registration...'
    if (Test-Path $UninstallKey) { Remove-Item -Path $UninstallKey -Recurse -Force }

    & $report 80 'Deleting the GAN folder...'
    Set-Location $env:TEMP   # so the folder we delete is not the one we are standing in
    if (Test-Path $dir) { Remove-Item -Path $dir -Recurse -Force }

    Send-EnvironmentChange
    & $report 100 'Done.'
}

# ---------------------------------------------------------------- the window
function New-SetupWindow {
    Add-Type -AssemblyName System.Windows.Forms
    Add-Type -AssemblyName System.Drawing
    [System.Windows.Forms.Application]::EnableVisualStyles()

    $cyan = [System.Drawing.ColorTranslator]::FromHtml('#00ADD8')
    $form = New-Object System.Windows.Forms.Form
    $form.Text = 'GAN Setup'
    $form.ClientSize = New-Object System.Drawing.Size(540, 430)
    $form.StartPosition = 'CenterScreen'
    $form.FormBorderStyle = 'FixedDialog'
    $form.MaximizeBox = $false
    $form.Font = New-Object System.Drawing.Font('Segoe UI', 9.5)
    $ico = Join-Path $ScriptDir 'gan-icon.ico'
    if (Test-Path $ico) { try { $form.Icon = New-Object System.Drawing.Icon($ico) } catch { } }

    # header strip
    $header = New-Object System.Windows.Forms.Panel
    $header.Dock = 'Top'; $header.Height = 90; $header.BackColor = $cyan
    $title = New-Object System.Windows.Forms.Label
    $title.Text = 'GAN'
    $title.Font = New-Object System.Drawing.Font('Segoe UI', 30, ([System.Drawing.FontStyle]::Bold -bor [System.Drawing.FontStyle]::Italic))
    $title.ForeColor = [System.Drawing.Color]::White
    $title.AutoSize = $true; $title.Location = New-Object System.Drawing.Point(24, 10)
    $sub = New-Object System.Windows.Forms.Label
    $sub.Text = "Setup  -  version $AppVersion"
    $sub.ForeColor = [System.Drawing.Color]::White
    $sub.AutoSize = $true; $sub.Location = New-Object System.Drawing.Point(300, 52)
    $header.Controls.AddRange(@($title, $sub))

    $intro = New-Object System.Windows.Forms.Label
    $intro.Text = "This installs the GAN programming language on your computer:`r`nthe gan command, the standard libraries, examples and the file icon.`r`nNo administrator rights are needed."
    $intro.Location = New-Object System.Drawing.Point(24, 104); $intro.Size = New-Object System.Drawing.Size(495, 62)

    $folderLabel = New-Object System.Windows.Forms.Label
    $folderLabel.Text = 'Install GAN in this folder:'
    $folderLabel.AutoSize = $true; $folderLabel.Location = New-Object System.Drawing.Point(24, 176)
    $folderBox = New-Object System.Windows.Forms.TextBox
    $folderBox.Text = $InstallDir
    $folderBox.Location = New-Object System.Drawing.Point(24, 200); $folderBox.Size = New-Object System.Drawing.Size(395, 26)
    $browse = New-Object System.Windows.Forms.Button
    $browse.Text = 'Browse...'
    $browse.Location = New-Object System.Drawing.Point(428, 198); $browse.Size = New-Object System.Drawing.Size(92, 30)
    $browse.Add_Click({
        $dlg = New-Object System.Windows.Forms.FolderBrowserDialog
        $dlg.Description = 'Choose where to install GAN'
        if ($dlg.ShowDialog() -eq 'OK') { $folderBox.Text = Join-Path $dlg.SelectedPath 'GAN' }
    })

    $cbPath = New-Object System.Windows.Forms.CheckBox
    $cbPath.Text = 'Add GAN to my PATH (so "gan file.gan" works in any folder)'
    $cbPath.Checked = $true; $cbPath.AutoSize = $true; $cbPath.Location = New-Object System.Drawing.Point(24, 242)
    $cbAssoc = New-Object System.Windows.Forms.CheckBox
    $cbAssoc.Text = 'Give .gan files the GAN icon, and run them when double-clicked'
    $cbAssoc.Checked = $true; $cbAssoc.AutoSize = $true; $cbAssoc.Location = New-Object System.Drawing.Point(24, 270)
    $cbCode = New-Object System.Windows.Forms.CheckBox
    $cbCode.Text = 'Install the VS Code icon theme (if VS Code is found)'
    $cbCode.Checked = $true; $cbCode.AutoSize = $true; $cbCode.Location = New-Object System.Drawing.Point(24, 298)

    $status = New-Object System.Windows.Forms.Label
    $status.Text = 'Ready to install.'
    $status.AutoSize = $false; $status.Location = New-Object System.Drawing.Point(24, 338); $status.Size = New-Object System.Drawing.Size(495, 22)
    $bar = New-Object System.Windows.Forms.ProgressBar
    $bar.Location = New-Object System.Drawing.Point(24, 362); $bar.Size = New-Object System.Drawing.Size(495, 18)

    $install = New-Object System.Windows.Forms.Button
    $install.Text = 'Install'
    $install.Location = New-Object System.Drawing.Point(322, 390); $install.Size = New-Object System.Drawing.Size(92, 30)
    $cancel = New-Object System.Windows.Forms.Button
    $cancel.Text = 'Cancel'
    $cancel.Location = New-Object System.Drawing.Point(428, 390); $cancel.Size = New-Object System.Drawing.Size(92, 30)
    $cancel.Add_Click({ $form.Close() })
    $form.AcceptButton = $install

    $install.Add_Click({
        $install.Enabled = $false; $cancel.Enabled = $false; $folderBox.Enabled = $false; $browse.Enabled = $false
        $cbPath.Enabled = $false; $cbAssoc.Enabled = $false; $cbCode.Enabled = $false
        $progress = { param($pct, $text) $bar.Value = $pct; $status.Text = $text; [System.Windows.Forms.Application]::DoEvents() }
        try {
            $extra = Install-Gan $folderBox.Text.Trim() $cbPath.Checked $cbAssoc.Checked $cbCode.Checked $progress
            $status.Text = 'GAN is installed.'
            [System.Windows.Forms.MessageBox]::Show(
                ("GAN $AppVersion is installed.`r`n`r`nOpen a NEW terminal window and try:`r`n`r`n    gan --version`r`n    gan `"" + (Join-Path $folderBox.Text.Trim() 'examples\hello.gan') + "`"`r`n" + $extra),
                'GAN Setup', 'OK', 'Information') | Out-Null
            $form.Close()
        } catch {
            $status.Text = 'Setup failed.'
            [System.Windows.Forms.MessageBox]::Show("Setup could not finish:`r`n`r`n" + $_.Exception.Message, 'GAN Setup', 'OK', 'Error') | Out-Null
            $cancel.Enabled = $true
        }
    })

    $form.Controls.AddRange(@($header, $intro, $folderLabel, $folderBox, $browse, $cbPath, $cbAssoc, $cbCode, $status, $bar, $install, $cancel))
    return $form
}

# ---------------------------------------------------------------- main
$quiet = { param($pct, $text) if ($Silent) { Write-Host ("[{0,3}%] {1}" -f $pct, $text) } }

if ($SelfTest) {
    $form = New-SetupWindow
    if ($Screenshot) {
        $form.Show(); [System.Windows.Forms.Application]::DoEvents()
        $bmp = New-Object System.Drawing.Bitmap($form.Width, $form.Height)
        $form.DrawToBitmap($bmp, (New-Object System.Drawing.Rectangle(0, 0, $form.Width, $form.Height)))
        $bmp.Save($Screenshot, [System.Drawing.Imaging.ImageFormat]::Png)
    }
    $form.Dispose()
    Write-Host 'self-test: the setup window was built without errors'
    exit 0
}

if ($Uninstall) {
    # running from inside the folder we are about to delete: copy ourselves to TEMP first
    if (-not $FromTemp -and $ScriptDir.TrimEnd('\') -ieq $InstallDir.TrimEnd('\')) {
        $tmp = Join-Path $env:TEMP 'gan-uninstall.ps1'
        Copy-Item $ScriptPath $tmp -Force
        $argsList = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-STA', '-WindowStyle', 'Hidden', '-File', $tmp, '-Uninstall', '-FromTemp', '-InstallDir', $InstallDir, '-RegBase', $RegBase)
        if ($Silent) { $argsList += '-Silent' }
        Start-Process powershell.exe -ArgumentList $argsList
        exit 0
    }
    # find where GAN is installed, if we were not told
    $item = Get-ItemProperty -Path $UninstallKey -ErrorAction SilentlyContinue
    if ($item -and $item.InstallLocation -and -not $PSBoundParameters.ContainsKey('InstallDir')) { $InstallDir = $item.InstallLocation }
    if (-not $Silent) {
        Add-Type -AssemblyName System.Windows.Forms
        $answer = [System.Windows.Forms.MessageBox]::Show("Remove GAN from this computer?`r`n`r`n$InstallDir", 'GAN Setup', 'YesNo', 'Question')
        if ($answer -ne 'Yes') { exit 0 }
    }
    Uninstall-Gan $InstallDir $quiet
    if (-not $Silent) { [System.Windows.Forms.MessageBox]::Show('GAN was removed.', 'GAN Setup', 'OK', 'Information') | Out-Null }
    exit 0
}

if ($Silent) {
    $extra = Install-Gan $InstallDir (-not $NoPath) (-not $NoAssoc) (-not $NoVSCode) $quiet
    Write-Host "GAN $AppVersion installed in $InstallDir.$extra"
    exit 0
}

$form = New-SetupWindow
[void]$form.ShowDialog()

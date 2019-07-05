$env:VMNAME="pangolin"
$env:SWITCHNAME="pangolin"
$env:DOCKERNAME="pangolin"

# $SERVERIP="0.0.0.0"
# $SERVERPORT="12345"
# $TOKENS='["token01", "token02"]'
# $ROLE="CLIENT"

function installPangolin () {
    $Adapter=(Get-NetRoute | Where-Object -FilterScript {$_.NextHop -Ne "::"} | Where-Object -FilterScript { $_.NextHop -Ne "0.0.0.0" } | Where-Object -FilterScript { ($_.NextHop.SubString(0,6) -Ne "fe80::") } | Get-NetAdapter ).Name.ToString()
    New-VMSwitch -Name $env:SWITCHNAME -NetAdapterName $Adapter -AllowManagementOS $true
    docker-machine create -d hyperv --hyperv-virtual-switch $env:SWITCHNAME $env:VMNAME
    docker-machine env $env:VMNAME
    & docker-machine.exe env $env:VMNAME | Invoke-Expression
    docker build -t $env:DOCKERNAME .
}

function uninstallPangolin () {
    docker-machine stop $env:VMNAME
    docker-machine rm -f $env:VMNAME
    Remove-VMSwitch -Name $env:SWITCHNAME -Force
}

function startPangolin () {
    & docker-machine.exe env $env:VMNAME | Invoke-Expression
    docker run --cap-add NET_ADMIN --cap-add NET_RAW --device /dev/net/tun:/dev/net/tun --net host --env ROLE=$env:ROLE --env SERVERIP=$env:SERVERIP --env SERVERPORT=$env:SERVERPORT --env TOKENS=$env:TOKENS pangolin
    sleep 3
}

function stopPangolin () {
    & docker-machine.exe env $env:VMNAME | Invoke-Expression
    docker-machine stop $env:VMNAME
    sleep 3
}


#UI##################################################################
Add-Type -assembly System.Windows.Forms
$X0=10; $X1=110; 
$Y0=10; $Y1=40; $Y2=70; $Y3=100; $Y4=130; $Y5=160;
$W0 = 100; $W1 = 300;
$BW = 100
$BX0=10; $BX1=110; $BX2=210; $BX3=310;

$labelIp = New-Object System.Windows.Forms.Label
$labelIp.Text = "Server IP"
$labelIp.Location = New-Object System.Drawing.Point($X0, $Y0)
$labelIp.Width = $W0

$labelPort = New-Object System.Windows.Forms.Label
$labelPort.Text = "Server Port"
$labelPort.Location = New-Object System.Drawing.Point($X0, $Y1)
$labelPort.Width = $W0

$labelTokens = New-Object System.Windows.Forms.Label
$labelTokens.Text = "Tokens"
$labelTokens.Location = New-Object System.Drawing.Point($X0, $Y2)
$labelTokens.Width = $W0

$labelRole = New-Object System.Windows.Forms.Label
$labelRole.Text = "Role"
$labelRole.Location = New-Object System.Drawing.Point($X0, $Y3)
$labelRole.Width = $W0

$textIp = New-Object System.Windows.Forms.TextBox
$textIp.Text = "0.0.0.0"
$textIp.Location = New-Object System.Drawing.Point($X1, $Y0)
$textIp.Width = $W1

$textPort = New-Object System.Windows.Forms.TextBox
$textPort.Text = "12345"
$textPort.Location = New-Object System.Drawing.Point($X1, $Y1)
$textPort.Width = $W1

$textTokens = New-Object System.Windows.Forms.TextBox
$textTokens.Text = '["toke01", "token02"]'
$textTokens.Location = New-Object System.Drawing.Point($X1, $Y2)
$textTokens.Width = $W1

$textBoxOutput = New-Object System.Windows.Forms.TextBox
$textBoxOutput.Location = New-Object System.Drawing.Point($X0, $Y5)
$textBoxOutput.Width = $W1 + $W0
$textBoxOutput.Height = 300
$textBoxOutput.AutoSize = $false
$textBoxOutput.Enabled = $true
$textBoxOutput.ScrollBars = "Vertical"
$textBoxOutput.Multiline = $true

$comboRole = New-Object System.Windows.Forms.ComboBox
$comboRole.Items.Add("SERVER")
$comboRole.Items.Add("CLIENT")
$comboRole.Location = New-Object System.Drawing.Point($X1, $Y3)
$comboRole.Width = $W1

$buttonInstall = New-Object System.Windows.Forms.Button
$buttonInstall.Text = 'Install'
$buttonInstall.Location = New-Object System.Drawing.Point($BX0, $Y4)
$buttonInstall.Width = $BW
$buttonInstall.BackColor = [System.Drawing.Color]::Yellow

$buttonUninstall = New-Object System.Windows.Forms.Button
$buttonUninstall.Text = 'Uninstall'
$buttonUninstall.Location = New-Object System.Drawing.Point($BX1, $Y4)
$buttonUninstall.Width = $BW
$buttonUninstall.BackColor = [System.Drawing.Color]::Red

$buttonStart = New-Object System.Windows.Forms.Button
$buttonStart.Text = 'Start'
$buttonStart.Location = New-Object System.Drawing.Point($BX2, $Y4)
$buttonStart.Width = $BW

$buttonStop = New-Object System.Windows.Forms.Button
$buttonStop.Text = 'Stop'
$buttonStop.Location = New-Object System.Drawing.Point($BX3, $Y4)
$buttonStop.Width = $BW

function initEnv () {
    $env:SERVERIP = $textIp.Text
    $env:SERVERPORT = $textPort.Text
    $env:TOKENS = $textTokens.Text
    $env:ROLE = $comboRole.SelectedText

    $env:installPangolin=$script:installPangolin
    $env:uninstallPangolin=$script:uninstallPangolin
    $env:startPangolin=$script:startPangolin
    $env:stopPangolin=$script:stopPangolin
}

function outputToGUI {
    param([parameter(ValueFromPipeline=$true)]$a)
    Process {
        $a | Out-String -OutBuffer 1 -Stream | ForEach-Object {
            $textBoxOutput.Lines = $textBoxOutput.Lines + $_
            $textBoxOutput.Select($textBoxOutput.Text.Length, 0)
            $textBoxOutput.ScrollToCaret()
            $main.Update()
        }
    }
}

$buttonInstall.Add_Click(
    {
        initEnv
        [System.Windows.Forms.MessageBox]::Show("Install will start, please wait","Install")
        installPangolin | outputToGUI
        [System.Windows.Forms.MessageBox]::Show("Done","Install")
    }
)

$buttonUninstall.Add_Click(
    {
        initEnv
        [System.Windows.Forms.MessageBox]::Show("Uninstall will start, please wait","Uninstall")
        uninstallPangolin | outputToGUI
        [System.Windows.Forms.MessageBox]::Show("Done","Uninstall")
    }
)

$buttonStart.Add_Click(
    {
        initEnv
        Start-Process powershell {
            $env:startPangolin
        } | outputToGUI
    }
)

$buttonStop.Add_Click(
    {
        initEnv
        Start-Process powershell {
            $env:stopPangolin
        } | outputToGUI
    }
)

$main = New-Object System.Windows.Forms.Form
$main.Text = "Pangolin"
$main.Width = 435
$main.Height = 510
$main.FormBorderStyle = "FixedDialog"
$main.MaximizeBox = $false
$main.Controls.Add($labelIp)
$main.Controls.Add($labelPort)
$main.Controls.Add($textIp)
$main.Controls.Add($textPort)
$main.Controls.Add($labelTokens)
$main.Controls.Add($labelRole)
$main.Controls.Add($comboRole)
$main.Controls.Add($textTokens)
$main.Controls.Add($buttonInstall)
$main.Controls.Add($buttonUninstall)
$main.Controls.Add($buttonStart)
$main.Controls.Add($buttonStop)
$main.Controls.Add($textBoxOutput)
$main.ShowDialog()

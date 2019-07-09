$env:VMNAME="pangolin"
$env:SWITCHNAME="pangolin"
$env:DOCKERNAME="pangolin"

$SERVERIP="0.0.0.0"
$SERVERPORT="12345"
$TOKENS='["token01", "token02"]'
$ROLE="CLIENT"

function installPangolin () {
    echo "Install"
    $Adapter=(Get-NetRoute | Where-Object -FilterScript {$_.NextHop -Ne "::"} | Where-Object -FilterScript { $_.NextHop -Ne "0.0.0.0" } | Where-Object -FilterScript { ($_.NextHop.SubString(0,6) -Ne "fe80::") } | Get-NetAdapter ).Name.ToString()
    New-VMSwitch -Name $env:SWITCHNAME -NetAdapterName $Adapter -AllowManagementOS $true
    docker-machine create -d hyperv --hyperv-virtual-switch $env:SWITCHNAME $env:VMNAME
    docker-machine env $env:VMNAME
    docker-machine.exe env $env:VMNAME | Invoke-Expression
    docker build -t $env:DOCKERNAME .
}

function uninstallPangolin () {
    echo "Uninstall"
    docker-machine stop $env:VMNAME
    docker-machine rm -f $env:VMNAME
    Remove-VMSwitch -Name $env:SWITCHNAME -Force
}

function startPangolin () {
    echo "Start Pangolin"
    & docker-machine.exe env $env:VMNAME | Invoke-Expression
    $pangolinIp=docker-machine.exe ip $env:VMNAME
    Get-NetRoute | where { $_.DestinationPrefix -eq '0.0.0.0/0' } | select { $_.NextHop } | route delete 0.0.0.0
    route add 0.0.0.0 mask 0.0.0.0 $pangolinIp
    $Adapter=((Get-NetRoute | Where-Object -FilterScript {$_.NextHop -Ne "::"} | Where-Object -FilterScript { $_.NextHop -Ne "0.0.0.0" } | Where-Object -FilterScript { ($_.NextHop.SubString(0,6) -Ne "fe80::") } | Get-NetAdapter ).Name.ToString())
    Set-DnsClientServerAddress -InterfaceAlias $Adapter -ServerAddresses("8.8.8.8")
    docker run -d --cap-add NET_ADMIN --cap-add NET_RAW --device /dev/net/tun:/dev/net/tun --net host --env ROLE=$env:ROLE --env SERVERIP=$env:SERVERIP --env SERVERPORT=$env:SERVERPORT --env TOKENS=$env:TOKENS pangolin
    sleep 3
}

function stopPangolin () {
    echo "Stop Pangolin"
    initEnv
    & docker-machine.exe env $env:VMNAME | Invoke-Expression
    docker ps | foreach { $s=$_ -split '\s+'; if($s[1] -eq $env:DOCKERNAME){docker kill $s[0];}}
    $Adapter=(Get-NetRoute | Where-Object -FilterScript {$_.NextHop -Ne "::"} | Where-Object -FilterScript { $_.NextHop -Ne "0.0.0.0" } | Where-Object -FilterScript { ($_.NextHop.SubString(0,6) -Ne "fe80::") } | Get-NetAdapter ).Name.ToString()
    Restart-NetAdapter $Adapter
    Set-DnsClientServerAddress -InterfaceAlias $Adapter -ResetServerAddresses
}

function restartVM () {
    echo "Restart VM"
    initEnv
    & docker-machine.exe env $env:VMNAME | Invoke-Expression
    docker-machine.exe restart $env:VMNAME
}

If (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
{
  # Relaunch as an elevated process:
  Start-Process powershell.exe "-File",('"{0}"' -f $MyInvocation.MyCommand.Path) -Verb RunAs
  exit
}

#UI##################################################################
Add-Type -assembly System.Windows.Forms
$X0=10; $X1=110; 
$Y0=40; $Y1=70; $Y2=100; $Y3=130; $Y4=160; $Y5=190;
$W0 = 100; $W1 = 300;

#managemenu
$installMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$installMenu.Text = "Install"
$uninstallMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$uninstallMenu.Text = "Uninstall"
$manageMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$manageMenu.Text = "Manage"
$manageMenu.DropDownItems.Add($installMenu)
$manageMenu.DropDownItems.Add($uninstallMenu)
#controlMenu
$restartVMMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$restartVMMenu.Text = "Restart VM"
$startMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$startMenu.text = 'Start'
$stopMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$stopMenu.text = 'Sop'
$controlMenu = New-Object System.Windows.Forms.ToolStripMenuItem
$controlMenu.Text = 'Control'
$controlMenu.DropDownItems.Add($startMenu)
$controlMenu.DropDownItems.Add($stopMenu)
$controlMenu.DropDownItems.Add($restartVMMenu)
#menu bar
$menuBar = New-Object System.Windows.Forms.MenuStrip
$menuBar.Items.Add($manageMenu)
$menuBar.Items.Add($controlMenu)


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
$textIp.Text = $SERVERIP
$textIp.Location = New-Object System.Drawing.Point($X1, $Y0)
$textIp.Width = $W1

$textPort = New-Object System.Windows.Forms.TextBox
$textPort.Text = $SERVERPORT
$textPort.Location = New-Object System.Drawing.Point($X1, $Y1)
$textPort.Width = $W1

$textTokens = New-Object System.Windows.Forms.TextBox
$textTokens.Text = $TOKENS
$textTokens.Location = New-Object System.Drawing.Point($X1, $Y2)
$textTokens.Width = $W1

$textBoxOutput = New-Object System.Windows.Forms.TextBox
$textBoxOutput.Location = New-Object System.Drawing.Point($X0, $Y4)
$textBoxOutput.Width = $W1 + $W0
$textBoxOutput.Height = 200
$textBoxOutput.AutoSize = $false
$textBoxOutput.Enabled = $true
$textBoxOutput.ScrollBars = "Vertical"
$textBoxOutput.Multiline = $true

$comboRole = New-Object System.Windows.Forms.ComboBox
$comboRole.Items.Add("SERVER")
$comboRole.Items.Add("CLIENT")
$comboRole.Location = New-Object System.Drawing.Point($X1, $Y3)
$comboRole.Width = $W1
if($ROLE -eq "CLIENT"){
    $comboRole.SelectedIndex = 1
} else {
    $comboRole.SelectedIndex = 0
}

function initEnv () {
    $env:SERVERIP = $textIp.Text
    $env:SERVERPORT = $textPort.Text
    $env:TOKENS = $textTokens.Text.Replace('"','\"')
    $env:ROLE = $comboRole.Text
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

$installMenu.Add_Click(
    {
        initEnv
        [System.Windows.Forms.MessageBox]::Show("Install will start, please wait","Install")
        installPangolin | outputToGUI
        [System.Windows.Forms.MessageBox]::Show("Done","Install")
    }
)

$uninstallMenu.Add_Click(
    {
        initEnv
        [System.Windows.Forms.MessageBox]::Show("Uninstall will start, please wait","Uninstall")
        uninstallPangolin | outputToGUI
        [System.Windows.Forms.MessageBox]::Show("Done","Uninstall")
    }
)

$startMenu.Add_Click(
    {
        initEnv
        startPangolin | outputToGUI
        #Start-Process powershell -ArgumentList ({. $env:FUNCFILE; startPangolin}) -NoNewWindow
        [System.Windows.Forms.MessageBox]::Show("Done","Start")
    }
)

$stopMenu.Add_Click(
    {
        initEnv
        stopPangolin | outputToGUI
        [System.Windows.Forms.MessageBox]::Show("Done","Stop")
    }
)

$restartVMMenu.Add_Click(
    {
        initEnv
        restartVM | outputToGUI
        [System.Windows.Forms.MessageBox]::Show("Done","Restart VM")
    }
)

$main = New-Object System.Windows.Forms.Form
$main.Text = "Pangolin"
$main.Width = 435
$main.Height = 410
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
$main.Controls.Add($textBoxOutput)
$main.Controls.Add($menuBar)
$main.ShowDialog()

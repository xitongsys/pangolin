$TOKENS='["token01", "token02"]'
$SERVERIP="47.240.40.78"
$SERVERPORT="12345"
$VMNAME="pangolin01"
$SWITCHNAME="pangolin"
$DOCKERNAME="pangolin"

function MacroReplace ($file) {
    (Get-Content $file).Replace('{SERVERIP}', $SERVERIP).Replace('{SERVERPORT}', $SERVERPORT).Replace('{TOKENS}', $TOKENS) | Set-Content $file
}

MacroReplace .\Dockerfile
MacroReplace .\pangolin\configs\cfg_client.json
MacroReplace .\pangolin\configs\cfg_server.json

$Adapter=(Get-NetRoute | Where-Object -FilterScript {$_.NextHop -Ne "::"} | Where-Object -FilterScript { $_.NextHop -Ne "0.0.0.0" } | Where-Object -FilterScript { ($_.NextHop.SubString(0,6) -Ne "fe80::") } | Get-NetAdapter ).Name.ToString()
New-VMSwitch -Name $SWITCHNAME -NetAdapterName $Adapter -AllowManagementOS $true
docker-machine create -d hyperv --hyperv-virtual-switch $SWITCHNAME $VMNAME
docker-machine restart $VMNAME
docker-machine env $VMNAME
& docker-machine.exe env $VMNAME | Invoke-Expression
docker build -t $DOCKERNAME .

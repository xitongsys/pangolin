param($cmd="install")

$VMNAME="pangolin01"
$SWITCHNAME="pangolin"
$DOCKERNAME="pangolin"

function install () {
    $Adapter=(Get-NetRoute | Where-Object -FilterScript {$_.NextHop -Ne "::"} | Where-Object -FilterScript { $_.NextHop -Ne "0.0.0.0" } | Where-Object -FilterScript { ($_.NextHop.SubString(0,6) -Ne "fe80::") } | Get-NetAdapter ).Name.ToString()
    New-VMSwitch -Name $SWITCHNAME -NetAdapterName $Adapter -AllowManagementOS $true
    docker-machine create -d hyperv --hyperv-virtual-switch $SWITCHNAME $VMNAME
    docker-machine restart $VMNAME
    docker-machine env $VMNAME
    & docker-machine.exe env $VMNAME | Invoke-Expression
    docker build -t $DOCKERNAME .
}

function uninstall () {
    docker-machine stop $VMNAME
    docker-machine rm -f $VMNAME
    Remove-VMSwitch -Name $SWITCHNAME -Force
}

if ($cmd -eq "install") {
    install
}

if ($cmd -eq "uninstall") {
    uninstall
}


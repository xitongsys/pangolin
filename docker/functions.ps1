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
    docker-machine stop $env:VMNAME
    sleep 3
}

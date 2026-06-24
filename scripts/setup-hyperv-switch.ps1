#Requires -RunAsAdministrator
<#
.SYNOPSIS
    Creates a Hyper-V Internal switch with a static host IP and a Windows NAT
    so that minikube gets a predictable, stable IP across restarts.

.DESCRIPTION
    The Hyper-V "Default Switch" assigns a dynamic 172.x.x.x address to the
    minikube VM on every restart, which breaks docker-env, helm proxy bypass,
    and any tooling that expects a fixed endpoint.

    This script creates an Internal switch (host<->VM only, no physical NIC
    bridge), assigns the host side a static IP, and enables Windows NAT so the
    VM can reach the internet through the host.

    After running this script, initialise gokube with:

        gokube init --driver=hyperv `
                    --hyperv-virtual-switch=<SwitchName> `
                    --static-ip=<MinikubeIP>

    Then add the minikube IP to NO_PROXY once (or make it permanent):

        $env:NO_PROXY = "$env:NO_PROXY,<MinikubeIP>"

.PARAMETER SwitchName
    Name for the new Hyper-V Internal switch. Defaults to "minikube-switch".
    Pass the same value to gokube init --hyperv-virtual-switch.

.PARAMETER HostIP
    IP address assigned to the host-side virtual NIC (the gateway the VM uses).
    Defaults to "192.168.200.1".

.PARAMETER PrefixLength
    Subnet prefix length. Defaults to 24 (i.e. /24, giving 192.168.200.0/24).

.PARAMETER MinikubeIP
    The static IP you intend to give the minikube VM. Only used to print the
    recommended gokube init command at the end - the script does not configure
    the VM directly. Defaults to "192.168.200.100".

.PARAMETER NatName
    Name for the Windows NetNat rule. Defaults to "minikube-nat".

.EXAMPLE
    # Default settings - creates minikube-switch on 192.168.200.0/24
    .\setup-hyperv-switch.ps1

.EXAMPLE
    # Custom subnet
    .\setup-hyperv-switch.ps1 -SwitchName "k8s-switch" -HostIP "10.10.20.1" -MinikubeIP "10.10.20.100"
#>

param(
    [string] $SwitchName   = "minikube-switch",
    [string] $HostIP       = "192.168.200.1",
    [int]    $PrefixLength = 24,
    [string] $MinikubeIP   = "192.168.200.100",
    [string] $NatName      = "minikube-nat"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# Derive subnet prefix from host IP + prefix length for the NAT rule.
$ipBytes   = [System.Net.IPAddress]::Parse($HostIP).GetAddressBytes()
$maskInt   = [uint32]([System.Math]::Pow(2, 32) - [System.Math]::Pow(2, 32 - $PrefixLength))
$maskBytes = [System.BitConverter]::GetBytes($maskInt)
if ([System.BitConverter]::IsLittleEndian) { [System.Array]::Reverse($maskBytes) }
$netBytes  = for ($i = 0; $i -lt 4; $i++) { $ipBytes[$i] -band $maskBytes[$i] }
$SubnetPrefix = ($netBytes -join ".") + "/$PrefixLength"

Write-Host ""
Write-Host "=== Hyper-V minikube switch setup ===" -ForegroundColor Cyan
Write-Host "  Switch      : $SwitchName"
Write-Host "  Subnet      : $SubnetPrefix"
Write-Host "  Host IP     : $HostIP/$PrefixLength"
Write-Host "  Minikube IP : $MinikubeIP  (passed to gokube init --static-ip)"
Write-Host "  NAT name    : $NatName"
Write-Host ""

# -- 1. Hyper-V Internal switch ----------------------------------------------
$existingSwitch = Get-VMSwitch -Name $SwitchName -ErrorAction SilentlyContinue
if ($existingSwitch) {
    Write-Host "[1/3] Virtual switch '$SwitchName' already exists - skipping." -ForegroundColor Yellow
} else {
    Write-Host "[1/3] Creating Internal Hyper-V switch '$SwitchName'..." -ForegroundColor Green
    New-VMSwitch -Name $SwitchName -SwitchType Internal | Out-Null
    Write-Host "      Done."
}

# -- 2. Host adapter static IP ------------------------------------------------
# The adapter for an Internal switch is named "vEthernet (<SwitchName>)".
$adapterName = "vEthernet ($SwitchName)"

# Wait up to 10 s for the adapter to appear (it is created asynchronously).
$adapter = $null
for ($i = 0; $i -lt 10; $i++) {
    $adapter = Get-NetAdapter -Name $adapterName -ErrorAction SilentlyContinue
    if ($adapter) { break }
    Start-Sleep -Seconds 1
}
if (-not $adapter) {
    Write-Error "Adapter '$adapterName' did not appear after 10 s. Aborting."
    exit 1
}

$existingIP = Get-NetIPAddress -InterfaceIndex $adapter.ifIndex -IPAddress $HostIP -ErrorAction SilentlyContinue
if ($existingIP) {
    Write-Host "[2/3] Host IP $HostIP already assigned to '$adapterName' - skipping." -ForegroundColor Yellow
} else {
    Write-Host "[2/3] Assigning $HostIP/$PrefixLength to '$adapterName'..." -ForegroundColor Green
    New-NetIPAddress -IPAddress $HostIP -PrefixLength $PrefixLength -InterfaceIndex $adapter.ifIndex | Out-Null
    Write-Host "      Done."
}

# -- 3. Windows NAT (outbound internet for the VM) ---------------------------
$existingNat = Get-NetNat -Name $NatName -ErrorAction SilentlyContinue
if ($existingNat) {
    Write-Host "[3/3] NetNat '$NatName' already exists - skipping." -ForegroundColor Yellow
} else {
    Write-Host "[3/3] Creating NetNat '$NatName' for $SubnetPrefix..." -ForegroundColor Green
    New-NetNat -Name $NatName -InternalIPInterfaceAddressPrefix $SubnetPrefix | Out-Null
    Write-Host "      Done."
}

# -- Summary ------------------------------------------------------------------

# Build display strings once to keep Write-Host calls simple.
$initCmd = "  gokube init --driver=hyperv --hyperv-virtual-switch=`"$SwitchName`" --static-ip=`"$MinikubeIP`""
$noProxySession   = '  $env:NO_PROXY = "$env:NO_PROXY,' + $MinikubeIP + '"'
$noProxyPermanent = '  [System.Environment]::SetEnvironmentVariable("NO_PROXY", $env:NO_PROXY + ",' + $MinikubeIP + '", "Machine")'
$dockerHost  = '  $env:DOCKER_HOST      = "tcp://' + $MinikubeIP + ':2376"'
$dockerTls   = '  $env:DOCKER_TLS_VERIFY = "1"'
$dockerCerts = '  $env:DOCKER_CERT_PATH  = "$env:USERPROFILE\.minikube\certs"'

Write-Host ""
Write-Host "=== Setup complete ===" -ForegroundColor Green
Write-Host ""
Write-Host "1. Run gokube init:" -ForegroundColor Cyan
Write-Host $initCmd -ForegroundColor White
Write-Host ""
Write-Host "2. Add the minikube IP to NO_PROXY (required to bypass corporate proxy):" -ForegroundColor Cyan
Write-Host $noProxySession -ForegroundColor White
Write-Host ""
Write-Host "   To make it permanent (system-wide, requires admin):" -ForegroundColor Cyan
Write-Host $noProxyPermanent -ForegroundColor White
Write-Host ""
Write-Host "3. Point the Docker client at the fixed address:" -ForegroundColor Cyan
Write-Host $dockerHost  -ForegroundColor White
Write-Host $dockerTls   -ForegroundColor White
Write-Host $dockerCerts -ForegroundColor White
Write-Host ""
Write-Host "   Or auto-configure each session via your PowerShell profile:" -ForegroundColor Cyan
Write-Host '  minikube docker-env | Invoke-Expression' -ForegroundColor White
Write-Host ""

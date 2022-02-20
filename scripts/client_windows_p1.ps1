$y = $args[0]
$x =  netsh interface ip show interfaces | Select-String -Pattern "^\s*(\d*).+$y$"
$interface_num = $x.Matches.groups[1].value
$ip = [System.Net.Dns]::GetHostAddresses($args[0])[0].IPAddressToString
route add 0.0.0.0 mask 128.0.0.0 $ip if $interface_num
route add 128.0.0.0 mask 128.0.0.0 $ip if $interface_num
netsh interface ipv4 set subinterface "$y" mtu=1400

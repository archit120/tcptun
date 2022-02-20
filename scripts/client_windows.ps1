$x = route PRINT | Select-String -Pattern "\s*0\.0\.0\.0\s*0\.0\.0\.0\s*([\d\.]*)"
$x -match "\s*0\.0\.0\.0\s*0\.0\.0\.0\s*([\d\.]*)"
$a = $Matches.1
$ip = [System.Net.Dns]::GetHostAddresses($args[0])[0].IPAddressToString
route add $ip mask 255.255.255.255 $a

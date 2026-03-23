Write-Host "Starting 3 nodes..."

# Node 1
$env:PORT="8080"; $env:PEERS=""
Start-Process "go" -ArgumentList "run","." -WindowStyle Hidden
Start-Sleep 3

# Node 2
$env:PORT="8081"; $env:PEERS="localhost:8080"
Start-Process "go" -ArgumentList "run","." -WindowStyle Hidden
Start-Sleep 3

# Node 3
$env:PORT="8082"; $env:PEERS="localhost:8080,localhost:8081"
Start-Process "go" -ArgumentList "run","." -WindowStyle Hidden

Write-Host "All nodes started!"
Write-Host "Check: http://localhost:8080/blockchain"
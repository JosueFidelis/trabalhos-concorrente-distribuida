set /p numberOfClients="Choose the number of clients"

start cmd /k "go run .\server\server.go"

timeout /t 2 /nobreak

set /A numberOfClients=numberOfClients-1
for /l %%x in (1, 1, %numberOfClients%) do (
   start cmd /k "go run .\client\client.go -1"
)
set /A numberOfClients=numberOfClients+1

start cmd /k "go run .\client\client.go %numberOfClients%"

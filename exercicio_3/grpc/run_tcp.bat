set /p numberOfClients="Choose the number of clients"

start cmd /k "go run .\server\sortarray_server.go"

timeout /t 2 /nobreak

set /A numberOfClients=numberOfClients-1
for /l %%x in (1, 1, %numberOfClients%) do (
   start cmd /k "go run .\client\sortarray_client.go -1"
)
set /A numberOfClients=numberOfClients+1

start cmd /k "go run .\client\sortarray_client.go %numberOfClients%"

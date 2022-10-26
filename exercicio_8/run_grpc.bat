setlocal enabledelayedexpansion
set /p numberOfClients="Choose the number of clients"
::set /a numberOfClients=6

start cmd /k "go run .\server\server.go"

timeout /t 2 /nobreak

for /l %%x in (1, 1, %numberOfClients%) do (
   echo %%x
   start cmd /k "go run .\client\client.go %%x"
)

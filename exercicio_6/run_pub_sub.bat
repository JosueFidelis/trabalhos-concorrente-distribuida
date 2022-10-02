set /p messageSize="Choose message size"

start cmd /k "go run .\publisher\publisher.go"

timeout /t 2 /nobreak

for /l %%x in (1, 1, 40) do (
   start cmd /k "go run .\subscriber\subscriber.go %messageSize%"
)

start cmd /k "go run .\client\main.go %messageSize% 1"

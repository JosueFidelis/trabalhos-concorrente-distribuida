set /p pc="Choose PC"

start cmd /k "go run .\subscriber\subscriber.go %pc%"

timeout /t 2 /nobreak

start cmd /k "go run .\publisher\publisher.go"


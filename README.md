# poc-sse
Server sent event poc

# Presquision
- Setup redis with config 
Addr:     "localhost:16379"
Username: "default",
Password: "moneyforward123",
- Install golang, npm, nodejs

## run sse-fe
- cd sse-fe
- npm i
- npm run dev
- Open browser : localhost:3000/sse

## run sse-server
- go run main.go 8080 // change arg[1] to run multi port for testing load balance

gooseup: 
		goose -dir ./db/migration/ -v mysql "edtech:E61@ch123@/edtech?parseTime=true" up

sqlc:
		sqlc generate
	
genproto:
		protoc -I ./api/rpc/ ./api/rpc/data.porto --go_out=:./rpc/  --go-grpc_out=./rpc/

goosedown: 
		goose -dir ./db/migration/ -v mysql "edtech:E61@ch123@/edtech?parseTime=true" down

mock:
	mockgen -package mockdb --destination db/mock/store.go Edtech_Golang/db/sqlc Store
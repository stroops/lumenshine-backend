
goose -dir ./db-files/migrations_src/ postgres "user=icop password=jw8s0F4 dbname=icop port=5433 sslmode=disable" down-to 0
goose -dir ./db-files/migrations_src/ postgres "user=icop password=jw8s0F4 dbname=icop port=5433 sslmode=disable" up
go generate
rice embed-go
go build

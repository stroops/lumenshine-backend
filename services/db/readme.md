# Local Migration

Running the migration, without embedding can be done manual with goose:
`goose -dir db-files/migrations_src postgres "user=icop password=xxx dbname=icop sslmode=disable port=5433" up`
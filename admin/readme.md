h2 Update db-models

I use this steps to regenerate the models from the db

* Edit the sql in the migration_src
* Run _rice embed-go_ to create the new ricefile with the embedded SQLs
* Run _go build_ to embed the new sql
* Run _.\db.exe migrate down-to 0_ to get an empty db
* Run _.\db.exe migrate up_ to get an current updated db
* Run _go generate_ to get the newest models

One can also run the migration, without building the code (e.g. if it does not compile)

* Install the goose command _go get -u github.com/pressly/goose/cmd/goose_
  You need gcc installed for this to work
* you can then run \_

goose.exe -dir .\db-files\migrations_src\ postgres "user=icop password=jw8s0F4 dbname=icop sslmode=disable" up

\_ in order to run the migration scripts

* after this you can run _go generate_ to let sqlboiler recreate the models

h1 For Christi and Theo

Please take care for the following "things"

* If there is a time in the db, please use
  `timestamp with time zone NOT NULL default current_timestamp`
  as a fieldtype in the DB. The timezone is important, if you convert the value to UTC to pass it to a proto-buffer call
* We DON'T use null's in the DB. Not even for string. Every field has a default value. This makes the servercode much more predictable and easyer to maintain
* Please create your tables with migrations only. If you do it manualy, you will have problems/difficulties, writing the migrations later, as the tables depend on each other
* when creating an index, please don't use `USING btree`. Just let the DB select the correct index type, based on the datatype. Exception to this is, if you know that you need another index type
* in the migrations, try to make your sqls as simple as possible (no joins in inserts e.g.)
* end every line with a semicolon in the sql files
* Please take care to export only functions, which are realy needed in other packages
* Please use the functions only for what they are designed for. E.G. a function named _LogAndReturnError_ is used to log and error. Dont use it without an error ...
* function should alsways return only error, not IcopError. IcopError is used only to the frontend. Stick to the golang style ... if no error is present (e.g. user not found), create one with errors.New()
* please DON'T set the http return status in a gin context, from inside a nested function. This makes maintainace hard, because we don't know where the status was set. An exception to this is, if you know exactly that you handle the state on your own (JWT/Authentication etc. but not e.g. inside a multiple used function like getUser).
* take care to log EVERY (go)error. Use info and warning for noncritical levels. If you don't have a logger in the function ... pass one to the function. We have the logger on every request set by default (with requestID, language etc.). Another possibility is, to just return the error (goish-way) and log it in the calling function. A little verbose, but this way everyone knows what's happening.
* please take care to not check for nil, we don't have null values in the db ... this will result in business logic runtime errors
* REST-Request will NEVER have pointer (*string *int) as parameters. please specify valuetypes only
* please use the validators from the icop portal. Till now we have icop_email, icop_phone
* Please NEVER EVER EVER EVER construct SQLs by yourself. NEVER. If you have to use string functions to handle the SQLs, then there is some error in your thinking ... NEVER DO THIS. This will introduce sql injection possibilites. Please NEVER do this. Using queries.Raw is an exception, that you only use to READ data, but never to write data. And even then, one should not need to use in most cases. @Cristi: please rewrite getUserGroupInsertQuery ...
* please read this look (in your freetime :) )
  https://github.com/bjut-hz/E-Books/blob/master/program%20language/Effective%20Go%20-%20The%20Go%20Programming%20Language.pdf
  This is important if you want to learn go ...
* @Cristi: Please don't implement functions like updateLastLogin, getUserByEmail inside the api package. If you want to update/query the DB, then move those functions in a db packeg or similar. General spoken, please divide the unit of consern inside the code


How to update the backend:
- git pull
- dep ensure
- go build
- icop build
- icop migrate --down
- icop migrate
- icop run
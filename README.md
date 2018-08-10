# lumenshine-backend


This is the backend for the lumenshine wallet.

Currently the project is alpha. You can run the project on your local machine or on a test server.

# Prerequirements

- We run and develop the backend on linux
- We use go for the whole backend
- You need the following go packages installed localy
- `make`
- `go get -u github.com/golang/dep/cmd/dep`
- `go get -u github.com/GeertJohan/go.rice/rice`
- `go get -u -t github.com/volatiletech/sqlboiler`
- We run some docker images with docker-compose for the DB and ETH/BTC clients. Example files are provided. You should install the latest docker/compose version

## SQLBoiler setup

We use go generate to create the models.
The SQLBoiler files are configured inside the main.go files.
We call SQLBoiler like this:

`//go:generate sqlboiler --wipe -b goose_db_version --no-tests --tinyint-as-bool=true --config $HOME/.config/sqlboiler/sqlboiler_db.toml postgres`

Therefore you need to create the `$HOME/.config/sqlboiler` folder and copy/configure the approriated file.


## Application/Service configuration

We use `https://github.com/spf13/viper` for reading the configuration.

Please define a global enviroment variable named `ICOP_CONFIG_DIR`
Inside this directory you have to place all configurations for the services/APIs.
Every service/API has its own configuration file.
Every configuration can be overriden by a local configuration file, that includes only the different values. This files are always named `xxx-local.toml`
Example configuration files are provided in the example-data directory.

# Local developing

If you start debugging the project(s) with VS-Code e.g., you have to pass in some argument to the services/APIs.

Example:
~~~~
{
    "name": "SRV-DB",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceRoot}/services/db",            
    "env": {"icop_debug": "1"}
}
~~~~

The `"env": {"icop_debug": "1"}` is needed, because we use live-reload of the configuration from viper. There is a problem with this feature and debugging, so that the debugger will not stop the application on exit. Therefore you can specify this argument and live-reload will not be enabled.

# Hint

The project is also refered as `icop` in some places. This will be changed in future.


# Directroy description

The backend is build as a monorepo.
Here is some information on the directories.

## addons

Includes some tools we use to extend the basic functionality

## admin

This is the backend for the admin-frontend.

## api

This includes all exposed APIs

## constants

This includes some commonly used constants

## db

This includes some commonly used db functionality

## example-data

Examples for the service configuration, docker, etc.

## helpers

Some helper package with shared functionality

## icop_errors

Defines some commonly used error constants

## pb

Includes the definitions and buildoutput for the gRPC services

## services

This includes all services, which are exposed internaly as gRPC.



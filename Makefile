BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse HEAD)
DATE := $(shell date '+%Y-%m-%d_%H:%M:%S')
CURRENT_DIR = $(shell pwd)
TARGET_DIR = $(CURRENT_DIR)/dist
CONFIG_DIR = $(TARGET_DIR)/config

define clear
	rm -rf $(TARGET_DIR)	
	mkdir $(TARGET_DIR)	
	mkdir $(CONFIG_DIR)
endef

define copy_services
    cp services/db/db $(TARGET_DIR)
	#cp services/db/data/db-config.toml $(CONFIG_DIR)

	cp services/2fa/2fa $(TARGET_DIR)
	#cp services/2fa/data/2fa-config.toml $(CONFIG_DIR)

	cp services/jwt/jwt $(TARGET_DIR)
	#cp services/jwt/data/jwt-config.toml $(CONFIG_DIR)

	cp services/mail/mail $(TARGET_DIR)
	#cp services/mail/data/mail-config.toml $(CONFIG_DIR)

	cp api/userapi/userapi $(TARGET_DIR)
	#cp api/userapi/data/userapi-config.toml $(CONFIG_DIR)

	cp api/payapi/payapi $(TARGET_DIR)
	#cp api/payapi/data/payapi-config.toml $(CONFIG_DIR)

	cp admin/admin $(TARGET_DIR)
	#cp admin/data/admin-config.toml $(CONFIG_DIR)
endef

.PHONY : all service-db service-2fa service-jwt service-mail api-userapi admin-api api-payapi service-pay

all : service-db service-2fa service-jwt service-mail api-userapi admin-api api-payapi service-pay
	$(call clear)
	$(call copy_services)
	
all-go: service-db service-2fa service-jwt service-mail api-userapi admin-api api-payapi service-pay
	$(call clear)
	$(call copy_services)

service-db: 
	cd services/db; rice embed-go; go build

service-2fa: 
	cd services/2fa; go build	

service-jwt: 
	cd services/jwt; go build		

service-mail: 
	cd services/mail; go build			

api-userapi: 	
	cd api/userapi; rice embed-go; go build -ldflags "-X main.buildDate=$(DATE) -X main.gitVersion=$(HASH) -X main.gitRemote=$(BRANCH)"

api-payapi: 	
	cd api/payapi; rice embed-go; go build -ldflags "-X main.buildDate=$(DATE) -X main.gitVersion=$(HASH) -X main.gitRemote=$(BRANCH)"

service-pay: 
	cd services/pay; go build	

admin-api: 
	cd admin; rice embed-go; go build

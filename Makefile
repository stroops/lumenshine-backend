BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse HEAD)
DATE := $(shell date '+%Y-%m-%d_%H:%M:%S')
CURRENT_DIR = $(shell pwd)
TARGET_DIR = $(CURRENT_DIR)/dist

define clear
	rm -rf $(TARGET_DIR)
	mkdir $(TARGET_DIR)
endef

define copy_services
    	cp services/db/db $(TARGET_DIR)

	cp services/2fa/2fa $(TARGET_DIR)

	cp services/jwt/jwt $(TARGET_DIR)

	cp services/mail/mail $(TARGET_DIR)

	cp api/userapi/userapi $(TARGET_DIR)

	cp api/payapi/payapi $(TARGET_DIR)

	cp services/pay/pay $(TARGET_DIR)

	cp admin/admin $(TARGET_DIR)

	cp addons/charts/charts $(TARGET_DIR)

endef

.PHONY : all service-db service-2fa service-jwt service-mail api-userapi admin-api api-payapi service-pay charts-addon docs

all: service-db service-2fa service-jwt service-mail api-userapi admin-api api-payapi service-pay charts-addon docs
	$(call clear)
	$(call copy_services)

docs:
	cd api/userapi; swagger generate spec -o ./userapi_swagger.yml -m
	cd api/payapi; swagger generate spec -o ./pay_api_swagger.yml -m
	cd admin; swagger generate spec -o ./adminapi_swagger.yml -m

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
	/bin/cp -rf "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"
	cd services/pay; go build

admin-api:
	cd admin; rice embed-go; go build

charts-addon:
	cd addons/charts; rice embed-go; go build

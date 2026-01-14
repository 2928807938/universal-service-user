module github.com/2928807938/universal-service-user/sdk

go 1.24.11

require (
	github.com/redis/go-redis/v9 v9.7.0
	gorm.io/gorm v1.25.12
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace (
	github.com/2928807938/universal-service-user/api/user-api => ../api/user-api
	github.com/2928807938/universal-service-user/auth/domain => ../auth/domain
	github.com/2928807938/universal-service-user/auth/infrastructure => ../auth/infrastructure
	github.com/2928807938/universal-service-user/hook => ../hook
	github.com/2928807938/universal-service-user/notification/domain => ../notification/domain
	github.com/2928807938/universal-service-user/notification/infrastructure => ../notification/infrastructure
	github.com/2928807938/universal-service-user/oauth/infrastructure => ../oauth/infrastructure
	github.com/2928807938/universal-service-user/rules => ../rules
	github.com/2928807938/universal-service-user/share => ../share
	github.com/2928807938/universal-service-user/user/domain => ../user/domain
	github.com/2928807938/universal-service-user/user/infrastructure => ../user/infrastructure
	github.com/2928807938/universal-service-user/verification/domain => ../verification/domain
	github.com/2928807938/universal-service-user/verification/infrastructure => ../verification/infrastructure
)

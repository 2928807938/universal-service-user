module github.com/2928807938/universal-service-user/auth/infrastructure

go 1.24.11

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
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
	github.com/2928807938/universal-service-user/auth/domain => ../domain
	github.com/2928807938/universal-service-user/bom => ../../bom
	github.com/2928807938/universal-service-user/share => ../../share
	github.com/2928807938/universal-service-user/user/domain => ../../user/domain
)

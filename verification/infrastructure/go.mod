module github.com/2928807938/universal-service-user/verification/infrastructure

go 1.24.11

require github.com/redis/go-redis/v9 v9.7.0

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace (
	github.com/2928807938/universal-service-user/share => ../../share
	github.com/2928807938/universal-service-user/verification/domain => ../domain
)

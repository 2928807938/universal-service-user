module universal-service-user/user/domain

go 1.24.11

require (
	// 通用工具
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.7.0
	golang.org/x/crypto v0.22.0
	universal-service-user/rules v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace universal-service-user/bom => ../../bom

replace universal-service-user/rules => ../../rules

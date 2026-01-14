module universal-service-user/sdk

go 1.24.11

require (
	github.com/redis/go-redis/v9 v9.7.0
	gorm.io/gorm v1.25.12
	universal-service-user/api/user-api v0.0.0
	universal-service-user/auth/domain v0.0.0
	universal-service-user/auth/infrastructure v0.0.0
	universal-service-user/hook v0.0.0
	universal-service-user/notification/domain v0.0.0
	universal-service-user/notification/infrastructure v0.0.0
	universal-service-user/oauth/infrastructure v0.0.0
	universal-service-user/share v0.0.0
	universal-service-user/user/domain v0.0.0
	universal-service-user/user/infrastructure v0.0.0
	universal-service-user/verification/domain v0.0.0
	universal-service-user/verification/infrastructure v0.0.0
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.63.107 // indirect
	github.com/bytedance/go-tagexpr/v2 v2.9.2 // indirect
	github.com/bytedance/gopkg v0.1.1 // indirect
	github.com/bytedance/sonic v1.12.6 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/hertz v0.9.3 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/henrylee2cn/ameda v1.4.10 // indirect
	github.com/henrylee2cn/goutil v0.0.0-20210127050712-89660552f6f8 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.1.49 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.1.49 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.2.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/driver/postgres v1.5.11 // indirect
	gorm.io/driver/sqlite v1.5.7 // indirect
	universal-service-user/rules v0.0.0 // indirect
)

replace (
	universal-service-user/api/user-api => ../api/user-api
	universal-service-user/auth/domain => ../auth/domain
	universal-service-user/auth/infrastructure => ../auth/infrastructure
	universal-service-user/hook => ../hook
	universal-service-user/notification/domain => ../notification/domain
	universal-service-user/notification/infrastructure => ../notification/infrastructure
	universal-service-user/oauth/infrastructure => ../oauth/infrastructure
	universal-service-user/rules => ../rules
	universal-service-user/share => ../share
	universal-service-user/user/domain => ../user/domain
	universal-service-user/user/infrastructure => ../user/infrastructure
	universal-service-user/verification/domain => ../verification/domain
	universal-service-user/verification/infrastructure => ../verification/infrastructure
)

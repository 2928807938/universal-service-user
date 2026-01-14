module universal-service-user/api/config-api

go 1.24.11

require (
	github.com/cloudwego/hertz v0.9.3
	github.com/google/uuid v1.6.0
	universal-service-user/configcenter/domain v0.0.0
	universal-service-user/share v0.0.0
)

replace (
	universal-service-user/configcenter/domain => ../../configcenter/domain
	universal-service-user/share => ../../share
)

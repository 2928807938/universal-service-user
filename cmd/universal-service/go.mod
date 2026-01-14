module universal-service-user/cmd/universal-service

go 1.24.11

require (
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.5.7
	gorm.io/driver/postgres v1.5.7
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
	universal-service-user/configcenter/domain v0.0.0
	universal-service-user/configcenter/infrastructure v0.0.0
)

replace (
	universal-service-user/configcenter/domain => ../../configcenter/domain
	universal-service-user/configcenter/infrastructure => ../../configcenter/infrastructure
)

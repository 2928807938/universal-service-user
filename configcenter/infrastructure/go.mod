module universal-service-user/configcenter/infrastructure

go 1.24.11

require (
	gorm.io/datatypes v1.2.6
	gorm.io/gorm v1.25.12
	universal-service-user/configcenter/domain v0.0.0
	universal-service-user/share v0.0.0
)

replace (
	universal-service-user/configcenter/domain => ../domain
	universal-service-user/share => ../../share
)

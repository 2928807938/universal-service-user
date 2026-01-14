module github.com/2928807938/universal-service-user/configcenter/infrastructure

go 1.24.11

require (
	gorm.io/datatypes v1.2.6
	gorm.io/gorm v1.25.12
	github.com/2928807938/universal-service-user/configcenter/domain v0.0.0
	github.com/2928807938/universal-service-user/share v0.0.0
)

replace (
	github.com/2928807938/universal-service-user/configcenter/domain => ../domain
	github.com/2928807938/universal-service-user/share => ../../share
)

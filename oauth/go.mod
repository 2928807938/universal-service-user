module github.com/2928807938/universal-service-user/oauth

go 1.24.11

require (
	github.com/2928807938/universal-service-user/oauth/domain v0.0.0
	github.com/2928807938/universal-service-user/oauth/infrastructure v0.0.0
)

replace (
	github.com/2928807938/universal-service-user/oauth/domain => ./domain
	github.com/2928807938/universal-service-user/oauth/infrastructure => ./infrastructure
)

module github.com/2928807938/universal-service-user/user

go 1.24.11

replace (
	github.com/2928807938/universal-service-user/user/domain => ./domain
	github.com/2928807938/universal-service-user/user/infrastructure => ./infrastructure
)

module universal-service-user/auth

go 1.24.11

replace (
	universal-service-user/auth/domain => ./domain
	universal-service-user/auth/infrastructure => ./infrastructure
)

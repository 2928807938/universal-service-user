module universal-service-user/user

go 1.24.11

replace (
	universal-service-user/user/domain => ./domain
	universal-service-user/user/infrastructure => ./infrastructure
)

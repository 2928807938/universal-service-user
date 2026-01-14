module universal-service-user/oauth

go 1.24.11

require (
	universal-service-user/oauth/domain v0.0.0
	universal-service-user/oauth/infrastructure v0.0.0
)

replace (
	universal-service-user/oauth/domain => ./domain
	universal-service-user/oauth/infrastructure => ./infrastructure
)

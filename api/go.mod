module universal-service-user/api

go 1.24.11

require (
	universal-service-user/api/config-api v0.0.0
	universal-service-user/api/user-api v0.0.0
)

replace universal-service-user/api/user-api => ./user-api
replace universal-service-user/api/config-api => ./config-api

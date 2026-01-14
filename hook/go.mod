module universal-service-user/hook

go 1.24.11

require universal-service-user/user/domain v0.0.0

require golang.org/x/crypto v0.22.0 // indirect

replace universal-service-user/user/domain => ../user/domain

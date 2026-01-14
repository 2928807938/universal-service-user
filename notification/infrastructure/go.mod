module github.com/2928807938/universal-service-user/notification/infrastructure

go 1.24.11

replace github.com/2928807938/universal-service-user/notification/domain => ../domain

replace github.com/2928807938/universal-service-user/verification => ../../verification

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.63.107
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.1.49
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.1.49
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

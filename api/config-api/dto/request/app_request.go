package request

// RegisterAppRequest app registration request.
type RegisterAppRequest struct {
	AppName     string `json:"app_name" vd:"len($)>0 && len($)<=128"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

package model

import "time"

// UserSSO - UserSSO
type UserSSO struct {
	//ID           string `bson:"id"`
	/*Username    string `json:"username" validate:"required"`
	Password    string `json:"password,omitempty" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Roles       string `json:"roles"`*/
	//ID            string `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Lastname      string `json:"lastname"`
	Useremail     string `json:"useremail"`
	Usertel       string `json:"usertel"`
	Userdate      string `json:"userdate"`
	Userstatus    string `json:"userstatus"`
	Userparentid  string `json:"userparentid"`
	Comid         string `json:"comid"`
	Comdepartment string `json:"comdepartment"`
	Usercode      string `json:"usercode"`
	Usertype      string `json:"usertype"`
	Userid        string `json:"userid"`
}

// UserSSOInfoResp - UserSSOInfoResp
type UserSSOInfoResp struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Lastname      string `json:"lastname"`
	Useremail     string `json:"useremail"`
	Usertel       string `json:"usertel"`
	Userdate      string `json:"userdate"`
	Userstatus    string `json:"userstatus"`
	Userparentid  string `json:"userparentid"`
	Comid         string `json:"comid"`
	Comdepartment string `json:"comdepartment"`
	Usercode      string `json:"usercode"`
	Usertype      string `json:"usertype"`
	Userid        string `json:"userid"`
}

// UsersSSO - UsersSSO
type UsersSSO []UserSSO

//Check Parent Have Child
type CheckChild struct {
	Loginusername string `json:"loginusername"`
	Checkusername string `json:"checkusername"`
}

//CompanyUserProductRole
type CompanyUserProductRole struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Productid string `json:"productid"`
}

//product_integrated
type ProductIntegrated struct {
	Productid   string `json:"productid"`
	Productname string `json:"productname"`
	Productdesc string `json:"productdesc"`
	Productdate string `json:"productdate"`
}
type ProductIntegrateds []ProductIntegrated

//company_product
type CompanyProduct struct {
	Comid          string `json:"comid"`
	Productid      string `json:"productid"`
	Contractcode   string `json:"contractcode"`
	Contractdate   string `json:"contractdate"`
	Contractstatus string `json:"contractstatus"`
	Datecreate     string `json:"datecreate"`
}
type CompanyProducts []CompanyProduct

//Company Info
type CompanyInfo struct {
	Comid            string `json:"comid"`
	Comshortname     string `json:"comshortname"`
	Comfullname      string `json:"comfullname"`
	Comaddress       string `json:"comaddress"`
	Comtel           string `json:"comtel"`
	Compersoncontact string `json:"compersoncontact"`
	Compersontel     string `json:"compersontel"`
	Compersonemail   string `json:"compersonemail"`
	Comstatus        string `json:"comstatus"`
	Comdate          string `json:"comdate"`
}
type CompanyInfos []CompanyInfo

//Company Info
type CompanyFunction struct {
	Functioncode   string `json:"functioncode"`
	Functionname   string `json:"functionname"`
	Functiondesc   string `json:"functiondesc"`
	Functionstatus string `json:"functionstatus"`
	Functiondate   string `json:"functiondate"`
	Comid          string `json:"comid"`
}
type CompanyFunctions []CompanyFunction

//Company Role
type CompanyRole struct {
	Rolecode   string `json:"rolecode"`
	Rolename   string `json:"rolename"`
	Roledesc   string `json:"roledesc"`
	Rolestatus string `json:"rolestatus"`
	Roledate   string `json:"roledate"`
	Comid      string `json:"comid"`
}
type CompanyRoles []CompanyRole

//Company Role Function
type CompanyRoleFunction struct {
	Rolecode         string   `json:"rolecode"`
	Productid        string   `json:"productid"`
	Functioncodelist []string `json:"functioncodelist"`
	Comid            string   `json:"comid"`
	Contractstatus   string   `json:"contractstatus"`
}
type CompanyRoleFunctions []CompanyRoleFunction

//company user role
type CompanyUserRole struct {
	Username   string   `json:"username"`
	Rolecode   []string `json:"rolecode"`
	Datecreate string   `json:"datecreate"`
	Comid      string   `json:"comid"`
}
type CompanyUserRoles []CompanyUserRole

//action log
type LogELK struct {
	Datetime time.Time `json:"datetime"`
	Log      string    `json:"log"`
}

//Request Token
//TokenBody
type TokenRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenRequestBodyv2 struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Clientid     string `json:"clientid"`
	Clientsecret string `json:"clientsecret"`
}
type AuthenRequestBody struct {
	Token string `json:"token"`
}
type AuthenRequestBodyV2 struct {
	Token     string `json:"token"`
	Productid string `json:"productid"`
}
type IoTRequestBody struct {
	Ssojwtoken string `json:"ssojwtoken"`
}

//Password management
//ForgotUsername
type ForgotUsername struct {
	InputUserInfo string `json:"inputuserinfo"`
}

//ResetPassword
type RequestResetPassword struct {
	Username string `json:"username"`
}
type ResetPassword struct {
	Otp         string `json:"otp"`
	Newpassword string `json:"newpassword"`
}

//refresh token
type RefreshTokenBody struct {
	Refreshtoken string `json:"refreshtoken"`
}

//paging
type Paging struct {
	Start int `json:"start"`
	Stop  int `json:"stop"`
}

//paging and filter
type PagingandFilter struct {
	Filterkey   string `json:"filterkey"`
	Filtervalue string `json:"filtervalue"`
	Start       int    `json:"start"`
	Stop        int    `json:"stop"`
}

//Kafka
type KafkaStruct struct {
	Action     string `json:"action"`
	Collection string `json:"collection"`
	Data       string `json:"data"`
}

//User service
type CompanyUserService struct {
	Username        string `json:"username"`
	Servicename     string `json:"servicename"`
	Servicetype     string `json:"servicetype"`
	Servicemetadata ServiceMetaData
	Comid           string `json:"comid"`
	Status          string `json:"status"`
}
type CompanyUserServices []CompanyUserService
type ServiceMetaData struct {
	Clientid     string `json:"clientId"`
	Clientsecret string `json:"clientSecret"`
	Serviceid    string `json:"serviceId"`
}

package service

import (
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"

	"github.com/labstack/echo"
)

//
var SysLog = new(model.LogELK)

// CreateUser - CreateUser
func CreateUser(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}

	if err := c.Validate(userSSO); err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}

	if !isStringAlphabetic(userSSO.Username) {
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	/*if !isPhoneNumber(userSSO.PhoneNumber) || (len(userSSO.PhoneNumber) < 10 || len(userSSO.PhoneNumber) > 11) {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "PhoneNumber is Invalid"})
	}*/
	//check length password
	if len(userSSO.Password) < 8 {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "Password less 8 character"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	if userSSO.Useremail == "" && userSSO.Usertel == "" {
		SysLog.Datetime = time.Now()
		SysLog.Log = "Email empty & Phone empty"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Email empty or Phone empty", "data": map[string]interface{}{"info": nil}})
	}
	countEmail, err2 := CheckUserEmailExist(userSSO.Useremail)
	countPhone, err3 := CheckUserPhoneExist(userSSO.Usertel)
	if err2 != nil || err3 != nil {
		SysLog.Datetime = time.Now()
		SysLog.Log = "MongoDB connection refused"
		WriteLog(SysLog)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	if countEmail > 0 || countPhone > 0 {
		SysLog.Datetime = time.Now()
		SysLog.Log = "Email or Phone number exist"
		WriteLog(SysLog)
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Email or Phone number exist", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	_, err = profileRepository.FindByUser(userSSO.Username)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "User is Exists", "data": map[string]interface{}{"info": nil}})
	}
	userSSO.Userid = uuid.New().String()
	err = profileRepository.Save(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": nil}})
}

// Update End User
func UpdateEndUser(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)

	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//check length password
	if len(userSSO.Password) < 8 {
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	err = profileRepository.Update(userSSO.Username, userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

// Update Company User
func UpdateCompanyUser(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//check length password
	if len(userSSO.Password) < 8 {
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	err = profileRepository.Update(userSSO.Username, userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	err = UpdateRelationshipNode(userSSO.Username, userSSO.Userstatus)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED UPDATE CHILD USER STATUS"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

// DeleteUser - DeleteUser
func DeleteEndUser(c echo.Context) error {
	// Get Parameter
	username := c.Param("username")

	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//delete user in collection users
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	err = profileRepository.Delete(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//delete user in collection company_user_role (cascade relationship)
	profileRepository2 := repository.NewProfileRepositoryMongo(db, "company_user_role")
	_, err = profileRepository2.FindCompanyUserRoleUsername(username)
	if err != mgo.ErrNotFound {
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Username exist"})
		err = profileRepository2.DeleteCompanyUserRole(username)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
		}
	}
	//Delete cas service
	err3 := RemoveServiceCasByUserName(username)
	if err3 != nil {
		fmt.Println(err3)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED REMOVE CAS SERVICE"})
	}
	//Delete user service
	err2 := DeleteAllUserServiceByUsername(username)
	if err2 != nil {
		fmt.Println(err2)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED REMOVE USER SERVICE"})
	}
	/*err = profileRepository2.DeleteCompanyUserRole(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}*/
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

// Get ONE UserInfo By Username
func GetUserInfo(c echo.Context) error {
	// Get info
	username := c.Param("username")

	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	userSSO, err := profileRepository.FindByUser(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "User is NonExists"})
	}
	userSSO.Password = ""
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": userSSO})
}

//Find all by Username
func GetAllUserByComId(c echo.Context) error {
	// Get info
	comid := c.Param("comid")

	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	usersSSO, err := profileRepository.FindAll(comid)
	//Che password
	for i := 0; i < len(usersSSO); i++ {
		usersSSO[i].Password = ""
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "User is NonExists"})
	}
	//usersSSO.Password = ""
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": usersSSO})
}

//
func GetUserChildInfo(c echo.Context) error {
	// Get info
	username := c.Param("username")
	err, result := GetInfomationChildOfNode(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//return c.JSON(http.StatusOK, result)
	return c.JSON(200, map[string]interface{}{"code": "0", "message": result})
}

//
type ArrayAuth struct {
	Rolecode         string   `json:"rolecode"`
	Functioncodelist []string `json:"functioncodelist"`
}
type ArrayAuths []ArrayAuth

func Auth(c echo.Context) error {
	//Sys log
	//sysLog := new(model.LogELK)
	//Sys log
	//GET INFO
	companyUserProductRole := new(model.CompanyUserProductRole)
	err := c.Bind(companyUserProductRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	username := companyUserProductRole.Username
	password := companyUserProductRole.Password
	productid := companyUserProductRole.Productid
	token, _ := RequestToken(username, password)
	//fmt.Println(token)
	if token == "" {
		//fmt.Println(token)
		//Sys log
		SysLog.Datetime = time.Now()
		SysLog.Log = "SSO SYSTEM AUTHEN USER FAILED"
		WriteLog(SysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "AUTHEN FAILED"})
	}
	authResponse, err := BasicAuth(token)
	//fmt.Println(authResponse)
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "VALID TOKEN FAILED"})
	}
	//authResponse := `{"service":"exampleOauthClient","attributes":{"address":"42 Phạm Ngọc Thạch P6 Q3","email":"naphaluan211092@gmail.com","first_name":"Lam","last_name":"Nguyen","phone_number":"0907888511","roles":"ROLE_ADMIN","user_id":"2bd2b378-6a4a-11ea-9686-beaadoooceb7"},"id":"thanhlam","client_id":"exampleOauthClient"}`
	//return c.JSON(200, map[string]interface{}{"code": "1", "message": authResponse})
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)
	attributes := result["attributes"].(map[string]interface{})

	//
	userstatus := attributes["userstatus"]

	if userstatus == "DISABLE" {
		//Sys log
		SysLog.Datetime = time.Now()
		SysLog.Log = "USER DISABLE"
		WriteLog(SysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "USER DISABLE"})
	}
	comid := attributes["comid"].(string)
	rolesUserInProductID, err := GetUserRoleWithProductID(username, comid, productid)
	//fmt.Println(rolesUserInProductID)
	//fmt.Println(err)
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Not found data"})
	}
	var arrayAuth ArrayAuth
	var arrayAuths ArrayAuths
	for i := 0; i < len(rolesUserInProductID); i++ {
		err, responseArr := GetCompanyUserRoleByRoleCode(rolesUserInProductID[i])
		if err != nil {
			panic(err)
		}
		arrayAuth.Rolecode = rolesUserInProductID[i]
		arrayAuth.Functioncodelist = responseArr
		arrayAuths = append(arrayAuths, arrayAuth)
		//fmt.Println(responseArr)
	}
	//fmt.Println(rolesUserInProductID)
	//return c.JSON(http.StatusOK, rolesUserInProductID)
	//return c.JSON(200, map[string]interface{}{"code": "1", "message": rolesUserInProductID})
	return c.JSON(200, map[string]interface{}{"code": "0", "message": arrayAuths})
}
func DeleteRelationshipUser(c echo.Context) error {
	// Get Parameter
	username := c.Param("username")
	//Check user
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	_, err = profileRepository.FindByUser(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "User is NonExists"})
	}
	//Delete cascase company_user_role node and sub-child node
	err = GetAllChildOfNodeAndDeleteCompanyUserRole(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP"})
	}
	//Delete relationship all child node
	err = DeleteNode(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//Delete cas service
	err3 := RemoveServiceCasByUserName(username)
	if err3 != nil {
		fmt.Println(err3)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED REMOVE CAS SERVICE"})
	}
	//Delete user service
	err2 := DeleteAllUserServiceByUsername(username)
	if err2 != nil {
		fmt.Println(err2)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED REMOVE USER SERVICE"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}
func CheckParentHaveChild(c echo.Context) error {
	checkChild := new(model.CheckChild)
	err := c.Bind(checkChild)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	err, exist := CheckParentAndChild(checkChild.Loginusername, checkChild.Checkusername)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": exist})
}
func RequestSSOToken(c echo.Context) error {
	tokenRequestBody := new(model.TokenRequestBody)
	err := c.Bind(tokenRequestBody)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	username := tokenRequestBody.Username
	password := tokenRequestBody.Password
	token, refreshtoken := RequestToken(username, password)
	if token == "" {
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "AUTHEN FAILED"})
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Authen Failed", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	//return c.JSON(200, map[string]interface{}{"code": "0", "message": token})
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"token": token, "refresh_token": refreshtoken}})
}
func RequestSSOTokenv2(c echo.Context) error {
	tokenRequestBodyv2 := new(model.TokenRequestBodyv2)
	err := c.Bind(tokenRequestBodyv2)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	username := tokenRequestBodyv2.Username
	password := tokenRequestBodyv2.Password
	clientID := tokenRequestBodyv2.Clientid
	clientSecret := tokenRequestBodyv2.Clientsecret
	serviceStatus, err := CheckServiceStatus(clientID)
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Failed Check Service Status", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	if serviceStatus != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Service Disable", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	token, refreshtoken := RequestTokenv2(clientID, clientSecret, username, password)
	if token == "" {
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "AUTHEN FAILED"})
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Authen Failed", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	//return c.JSON(200, map[string]interface{}{"code": "0", "message": token})
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"token": token, "refresh_token": refreshtoken}})
}
func AuthenSSOToken(c echo.Context) error {
	authenRequestBody := new(model.AuthenRequestBody)
	err := c.Bind(authenRequestBody)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	token := authenRequestBody.Token
	authResponse, err := BasicAuth(token)
	//
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		//return c.JSON(200, map[string]interface{}{"code": "-4", "message": err.Error()})
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)

	if result["error"] != nil {
		//errorToken := result["error"]
		if result["message"] != nil {
			//syslog
			SysLog.Datetime = time.Now()
			SysLog.Log = (result["message"]).(string)
			WriteLog(SysLog)
			//syslog
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = (result["error"]).(string)
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	//test
	attributes := result["attributes"].(map[string]interface{})
	attributes["username"] = result["id"]
	//fmt.Println(username)
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"info": attributes}})
}
func RefreshSSOToken(c echo.Context) error {
	refreshTokenBody := new(model.RefreshTokenBody)
	err := c.Bind(refreshTokenBody)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"token": nil}})
	}
	refreshtoken := refreshTokenBody.Refreshtoken

	token := RequestTokenByRT(refreshtoken)
	if token == "" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Invalid Request", "data": map[string]interface{}{"token": nil}})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"token": token}})
}
func AuthenSSOTokenV2(c echo.Context) error {
	authenRequestBody := new(model.AuthenRequestBody)
	err := c.Bind(authenRequestBody)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	token := authenRequestBody.Token
	authResponse, err := BasicAuth(token)
	//
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)

	if result["error"] != nil {
		//errorToken := result["error"]
		if result["message"] != nil {
			//syslog
			SysLog.Datetime = time.Now()
			SysLog.Log = (result["message"]).(string)
			WriteLog(SysLog)
			//syslog
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = (result["error"]).(string)
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	//test
	attributes := result["attributes"].(map[string]interface{})
	attributes["username"] = result["id"]
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"info": attributes}})
}

//paging
func PagingUserInfo(c echo.Context) error {
	paging := new(model.Paging)
	err := c.Bind(paging)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	usersSSO := new(model.UsersSSO)
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	usersSSO, err = profileRepository.PagingUser(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": usersSSO})
}

//paging with filter
func PagingAndFilterUserInfo(c echo.Context) error {
	pagingandFilter := new(model.PagingandFilter)
	err := c.Bind(pagingandFilter)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	usersSSO := new(model.UsersSSO)
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	usersSSO, err = profileRepository.PagingUserAndFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": usersSSO})
}

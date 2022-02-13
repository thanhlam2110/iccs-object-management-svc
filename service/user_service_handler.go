package service

import (
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"github.com/labstack/echo"
)

//CREATE COMPANY FUNCTION
func CreateUserServiceOauth(c echo.Context) error {
	companyUserService := new(model.CompanyUserService)
	serviceMetaData := new(model.ServiceMetaData)
	err := c.Bind(companyUserService)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	metaData := companyUserService.Servicemetadata
	serviceMetaData.Clientid = metaData.Clientid
	serviceMetaData.Clientsecret = metaData.Clientsecret
	serviceMetaData.Serviceid = metaData.Serviceid
	//db, err := config.GetMongoDB()
	//fmt.Println(companyUserService.Servicename)
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CONNECT DATABASE"})
	}
	//CHECK USERNAME EXIST
	countUsername, err1 := CheckUserNameExist(companyUserService.Username)
	if err1 != nil {
		fmt.Println(err1)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK COMID EXIST"})
	}
	if countUsername < 1 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Username non-exist"})
	}
	//CHECK SERVICE NAME
	count, err2 := CheckUserServiceNameExist(companyUserService.Servicename)
	//fmt.Println(count)
	if err2 != nil {
		fmt.Println(err2)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK SERVICE NAME EXIST"})
	}
	if count >= 1 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Service name exist"})
	}
	//CHECK OAUTH2 CLIENT ID
	countClientId, err3 := CheckUserServiceOauthClientIdExist(metaData.Clientid)
	if err3 != nil {
		fmt.Println(err3)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK COMID EXIST"})
	}
	if countClientId >= 1 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Oauth2.0 ClientID exist"})
	}
	//CLONE SERVICE
	err = CloneAndUpdateServiceOauth(companyUserService.Servicename, serviceMetaData)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CLONE SERVICE"})
	}
	//ADD THÔNG TIN QUẢN LÝ
	err = AddUserService(companyUserService)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CREATE COMPANY USER SERVICE"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}
func UpdateUserServiceOauth(c echo.Context) error {
	companyUserService := new(model.CompanyUserService)
	err := c.Bind(companyUserService)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	err = UpdateOauthUserService(companyUserService)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Fail Update Service"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//GET ALL BY USERNAME
func GetAllUserServiceByUsername(c echo.Context) error {
	// Get info
	username := c.Param("username")
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_service")
	companyUserServices, err := profileRepository.FindAllUserSeriveByUsername(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "User Service is NonExists"})
	}
	//usersSSO.Password = ""
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyUserServices})
}
func GetAllUserService(c echo.Context) error {
	// Get info
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_service")
	companyUserServices, err := profileRepository.FindAllUserSerive()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "User Service is NonExists"})
	}
	//usersSSO.Password = ""
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyUserServices})
}
func GetAllUserServiceByComid(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_service")
	companyUserServices, err := profileRepository.FindAllUserSeriveByComid(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "User Service is NonExists"})
	}
	//usersSSO.Password = ""
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyUserServices})
}

//
func PagingCompanyUserService(c echo.Context) error {
	paging := new(model.Paging)
	err := c.Bind(paging)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	companyUserServices := new(model.CompanyUserServices)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_service")
	companyUserServices, err = profileRepository.PagingUserServiceDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyUserServices})
}

func PagingCompanyUserServiceAndFilter(c echo.Context) error {
	pagingandFilter := new(model.PagingandFilter)
	err := c.Bind(pagingandFilter)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	companyUserServices := new(model.CompanyUserServices)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_service")
	companyUserServices, err = profileRepository.PagingUserServiceAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyUserServices})
}
func DeleteUserServiceOauthByServiceName(c echo.Context) error {
	serviceName := c.Param("servicename")
	//
	//fmt.Println(serviceName)
	err := RemoveServiceCas("name", string(serviceName))
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED REMOVE CAS SERVICE"})
	}
	err1 := RemoveUserService("servicename", string(serviceName))
	if err1 != nil {
		fmt.Println(err1)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED REMOVE USER SERVICE"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

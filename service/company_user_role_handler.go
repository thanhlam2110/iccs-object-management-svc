package service

import (
	"fmt"
	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

//FIND ONE
func GetCompanyUserRoleByUsername(c echo.Context) error {
	// Get info
	username := c.Param("username")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	companyUserRole, err := profileRepository.FindCompanyUserRoleUsername(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Username is NonExists"})
	}
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyUserRole})
}

//FIND ALL
func GetCompanyUserRoleByComId(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	companyUserRoles, err := profileRepository.FindCompanyUserRoleComId(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Comid is NonExists"})
	}
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyUserRoles})
}
func CreateCompanyUserRole(c echo.Context) error {
	companyUserRole := new(model.CompanyUserRole)
	err := c.Bind(companyUserRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	_, err = profileRepository.FindCompanyUserRoleUsername(companyUserRole.Username)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Username exist"})
	}
	//companyRoleFunction.Contractstatus="ACTIVE"
	err = profileRepository.SaveCompanyUserRole(companyUserRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//DELETE BY ROLECODE
func DeleteCompanyUserRoleByUsername(c echo.Context) error {
	// Get Parameter
	username := c.Param("username")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//delete company role
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	err = profileRepository.DeleteCompanyUserRole(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//UPDATE
func UpdateCompanyUserRole(c echo.Context) error {
	companyUserRole := new(model.CompanyUserRole)
	err := c.Bind(companyUserRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}

	//connect Mongo
	//db, err := config.GetMongoDB()
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	/*profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	err = profileRepository.UpdateCompanyUserRole(companyUserRole.Username, companyUserRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}*/
	err = UpdateCompanyUserRoleRelationship(companyUserRole.Username, companyUserRole.Datecreate, companyUserRole.Rolecode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//paging
func PagingCompanyUseRole(c echo.Context) error {
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
	companyUserRoles := new(model.CompanyUserRoles)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	companyUserRoles, err = profileRepository.PagingCompanyUseRoleDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyUserRoles})
}

//paging with filter
func PagingCompanyUseRoleAndFilter(c echo.Context) error {
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
	companyUserRoles := new(model.CompanyUserRoles)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	companyUserRoles, err = profileRepository.PagingCompanyUseRoleAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyUserRoles})
}

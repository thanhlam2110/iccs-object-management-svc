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
func GetCompanyRoleFunctionByRoleCode(c echo.Context) error {
	// Get info
	rolecode := c.Param("rolecode")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	companyRole, err := profileRepository.FindCompanyRoleFunctionRoleCode(rolecode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "RoleCode is NonExists"})
	}
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyRole})
}
func GetCompanyRoleFunctionByComId(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	companyRole, err := profileRepository.FindCompanyRoleFunctionComId(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Comid is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyRole})
}
func CreateCompanyRoleFunction(c echo.Context) error {
	companyRoleFunction := new(model.CompanyRoleFunction)
	err := c.Bind(companyRoleFunction)
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

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	_, err = profileRepository.FindCompanyRoleFunctionRoleCode(companyRoleFunction.Rolecode)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Role code exist"})
	}
	//companyRoleFunction.Contractstatus="ACTIVE"
	err = profileRepository.SaveCompanyRoleFunction(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}
func UpdateFunctionCodeCompanyRoleFunction(c echo.Context) error {
	companyRoleFunction := new(model.CompanyRoleFunction)
	err := c.Bind(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	err = UpdateArrayFunctionCodeList(companyRoleFunction.Rolecode, companyRoleFunction.Functioncodelist)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Update success")
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//DELETE BY ROLECODE
func DeleteCompanyRoleFunctionByRoleCode(c echo.Context) error {
	// Get Parameter
	rolecode := c.Param("rolecode")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//delete company role
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	err = profileRepository.DeleteCompanyRoleFunction(rolecode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//paging
func PagingCompanyRoleFunction(c echo.Context) error {
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
	companyRoleFunctions := new(model.CompanyRoleFunctions)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	companyRoleFunctions, err = profileRepository.PagingCompanyRoleFunctionDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyRoleFunctions})
}

//paging with filter
func PagingCompanyRoleFunctionAndFilter(c echo.Context) error {
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
	companyRoleFunctions := new(model.CompanyRoleFunctions)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	companyRoleFunctions, err = profileRepository.PagingCompanyRoleFunctionAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyRoleFunctions})
}

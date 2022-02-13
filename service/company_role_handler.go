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

//GET INFO COMPANY INFO BY CODE
func GetCompanyRoleByCode(c echo.Context) error {
	// Get info
	rolecode := c.Param("rolecode")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	companyRole, err := profileRepository.FindCompanyRoleByCode(rolecode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "RoleCode is NonExists"})
	}
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyRole})
}
func GetCompanyRoleByComID(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	companyRoles, err := profileRepository.FindCompanyRole(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "ComID is NonExists"})
	}
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyRoles})
}
func CreateCompanyRole(c echo.Context) error {
	companyRole := new(model.CompanyRole)
	err := c.Bind(companyRole)
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

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	_, err = profileRepository.FindCompanyRoleByCode(companyRole.Rolecode)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Role code exist"})
	}
	err = profileRepository.SaveCompanyRole(companyRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//UPDATE COMPANY FUNCTION BY FUNCTION CODE
func UpdateCompanyRoleByRoleCode(c echo.Context) error {
	companyRole := new(model.CompanyRole)
	err := c.Bind(companyRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	err = profileRepository.UpdateCompanyRoleCode(companyRole.Rolecode, companyRole)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//update collection company_role_function when update Rolestatus
	/*err = UpdateCompanyRoleFunctionRoleStatusByRoleCode(companyRole.Rolecode, companyRole.Rolestatus)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED UPDATE RELATIONSHIP ROLE STATUS"})
	}*/
	if companyRole.Rolestatus == "DISABLE" {
		var roleCodeArray []string
		roleCodeArray = append(roleCodeArray, companyRole.Rolecode)
		err = DeleteArrayRoleCodeCompanyUserRole(roleCodeArray)
		roleCodeArray = nil
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP ROLE CODE"})
		}
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//DELETE COMPANY PRODUCT BY  CODE
func DeleteCompanyRoleCode(c echo.Context) error {
	// Get Parameter
	rolecode := c.Param("rolecode")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED 1"})
	}
	//delete company role
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	err = profileRepository.DeleteCompanyRole(rolecode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED 2"})
	}
	////delete collection company_user_role when delete company_role
	var roleCodeArray []string
	roleCodeArray = append(roleCodeArray, rolecode)
	err = DeleteArrayRoleCodeCompanyUserRole(roleCodeArray)
	roleCodeArray = nil
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP ROLE CODE"})
	}
	//
	////delete collection company_role_function when delete company_role
	profileRepository2 := repository.NewProfileRepositoryMongo(db, "company_role_function")
	_, err = profileRepository2.FindCompanyRoleFunctionRoleCode(rolecode)
	if err != mgo.ErrNotFound {
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Role code exist"})
		err = profileRepository2.DeleteCompanyRoleFunction(rolecode)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP COMPANY ROLE FUNCTION"})
		}
	}
	/*err = profileRepository2.DeleteCompanyRoleFunction(rolecode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP COMPANY ROLE FUNCTION"})
	}*/
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//paging
func PagingCompanyRole(c echo.Context) error {
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
	companyRoles := new(model.CompanyRoles)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	companyRoles, err = profileRepository.PagingCompanyRoleDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyRoles})
}

//paging with filter
func PagingCompanyRoleAndFilter(c echo.Context) error {
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
	companyRoles := new(model.CompanyRoles)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	companyRoles, err = profileRepository.PagingCompanyRoleAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyRoles})
}

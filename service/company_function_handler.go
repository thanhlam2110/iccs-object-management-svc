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
func GetCompanyFunctionByCode(c echo.Context) error {
	// Get info
	functioncode := c.Param("functioncode")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	companyFunction, err := profileRepository.FindCompanyFunctionByCode(functioncode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Function code is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyFunction})
}

//GET INFO COMPANY FUNCTION BY COMID
func GetCompanyFunctionInfoByComID(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	companyFunctions, err := profileRepository.FindCompanyFunction(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "ComID is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyFunctions})
}

//CREATE COMPANY FUNCTION
func CreateCompanyFunction(c echo.Context) error {
	companyFunction := new(model.CompanyFunction)
	err := c.Bind(companyFunction)
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

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	_, err = profileRepository.FindCompanyFunctionByCode(companyFunction.Functioncode)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Function code exist"})
	}
	err = profileRepository.SaveCompanyFunction(companyFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//DELETE COMPANY PRODUCT BY  CODE
func DeleteCompanyFunctionCode(c echo.Context) error {
	// Get Parameter
	functioncode := c.Param("functioncode")

	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	////delete collection company_role_function when delete company_role_function
	var functionCodeArray []string
	functionCodeArray = append(functionCodeArray, functioncode)
	err = DeleteArrayFunctionCodeCompanyRoleFunction(functionCodeArray)
	functionCodeArray = nil
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP FUNCTION CODE"})
	}
	//
	err = profileRepository.DeleteCompanyFunctionByFunctionCode(functioncode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//DELETE COMPANY PRODUCT BY COMID
/*func DeleteCompanyFunctionComId(c echo.Context) error {
	// Get Parameter
	comid := c.Param("comid")

	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	err = profileRepository.DeleteCompanyFunction(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}*/

//UPDATE COMPANY FUNCTION BY FUNCTION CODE
func UpdateCompanyFunctionByCode(c echo.Context) error {
	companyFunction := new(model.CompanyFunction)
	err := c.Bind(companyFunction)
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
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	err = profileRepository.UpdateCompanyFunction(companyFunction.Functioncode, companyFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	//delete collection company_role_function when update Functionstatus=="DISABLE"
	if companyFunction.Functionstatus == "DISABLE" {
		var functionCodeArray []string
		functionCodeArray = append(functionCodeArray, companyFunction.Functioncode)
		err = DeleteArrayFunctionCodeCompanyRoleFunction(functionCodeArray)
		functionCodeArray = nil
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP FUNCTION CODE"})
		}
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//paging
func PagingCompanyFunction(c echo.Context) error {
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
	companyFunctions := new(model.CompanyFunctions)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	companyFunctions, err = profileRepository.PagingCompanyFunctionDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyFunctions})
}

//paging with filter
func PagingCompanyFunctionAndFilter(c echo.Context) error {
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
	companyFunctions := new(model.CompanyFunctions)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	companyFunctions, err = profileRepository.PagingCompanyFunctionAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyFunctions})
}

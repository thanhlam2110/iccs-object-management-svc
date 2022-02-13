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

func CreateCompanyInfo(c echo.Context) error {
	companyInfo := new(model.CompanyInfo)
	err := c.Bind(companyInfo)
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

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	_, err = profileRepository.FindCompanyInfo(companyInfo.Comid)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Comid code exist"})
	}
	err = profileRepository.SaveCompanyInfo(companyInfo)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//GET INFO COMPANY INFO BY COMID
func GetCompanyInfoByComID(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	//connect Mongo
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	companyInfo, err := profileRepository.FindCompanyInfo(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "ComID is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyInfo})
}

//DELETE PRODUCT INTEGRATED
func DeleteCompanyInfo(c echo.Context) error {
	// Get Parameter
	comid := c.Param("comid")
	err := RemoveAllRelatedCompanyComid(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//UPDATE PRODUCT INTEGRATED
func UpdateCompanyInfo(c echo.Context) error {
	companyInfo := new(model.CompanyInfo)
	err := c.Bind(companyInfo)
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

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	err = profileRepository.UpdateCompanyInfoByComId(companyInfo.Comid, companyInfo)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	if companyInfo.Comstatus == "DISABLE" {
		err = UpdateCascadeUserStatusCompanyInfo(companyInfo.Comid)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
		}
		err = UpdateCascadeContractStatusCompanyInfo(companyInfo.Comid)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
		}
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//GET all company info
func GetAllCompanyInfo(c echo.Context) error {
	// Get info
	//all := c.Param("all")
	//fmt.Println(all)
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
		return c.JSON(400, map[string]interface{}{"code": "10", "message": "Error connect MongoDB", "data": map[string]interface{}{"info": nil}})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	companyInfos, err := profileRepository.FindCompanyInfoAll()
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Bla bla"})
		return c.JSON(200, map[string]interface{}{"code": "4", "message": "No data", "data": map[string]interface{}{"info": nil}})
	}
	session.Close()
	//return c.JSON(200, map[string]interface{}{"code": "1", "message": companyInfos})
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": companyInfos}})
}

//paging
func PagingCompanyInfo(c echo.Context) error {
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
	companyInfos := new(model.CompanyInfos)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	companyInfos, err = profileRepository.PagingCompanyInfoDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyInfos})
}

//paging with filter
func PagingAndFilterCompanyInfo(c echo.Context) error {
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
	companyInfos := new(model.CompanyInfos)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	companyInfos, err = profileRepository.PagingCompanyInfoAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyInfos})
}

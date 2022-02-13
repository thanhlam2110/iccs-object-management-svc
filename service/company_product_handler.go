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

func CreateCompanyProduct(c echo.Context) error {
	companyProduct := new(model.CompanyProduct)
	err := c.Bind(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	//Check product integrated
	count, err2 := CheckProductIntegratedElementExist("productid", companyProduct.Productid)
	if err2 != nil {
		fmt.Println(err2)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK PRODUCTID"})
	}
	//Check product integrated
	_, err = profileRepository.FindCompanyProductContract(companyProduct.Contractcode)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Contract code exist"})
	}
	if count < 1 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Product Integrated not found"})
	}
	//Check twice comid + productid
	checkTwice, err3 := CheckTwiceComIDandProductIDinCompanyProduct(companyProduct.Comid, companyProduct.Productid)
	if err3 != nil {
		fmt.Println(err3)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK DOUBLE COMID AND PRODUCTID"})
	}
	if checkTwice > 0 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "DON'T ALLOW CREATE DUPLICATE SAME COMID AND PRODUCTID"})
	}
	//Check twice comid + productid
	err = profileRepository.SaveCompanyProduct(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//GET INFO COMPANY PRODUCT BY COMID
func GetCompanyProductInfoByComID(c echo.Context) error {
	// Get info
	comid := c.Param("comid")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	companyProducts, err := profileRepository.FindCompanyProduct(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "ComID is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyProducts})
}

//GET INFO COMPANY PRODUCT BY CONTRACT CODE
func GetCompanyProductInfoByContractCode(c echo.Context) error {
	// Get info
	contractcode := c.Param("contractcode")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	companyProduct, err := profileRepository.FindCompanyProductContract(contractcode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Contract code is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": companyProduct})
}

//DELETE ALL COMPANY PRODUCT BY COMID
func DeleteCompanyProductByComId(c echo.Context) error {
	// Get Parameter
	comid := c.Param("comid")

	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	err = profileRepository.DeleteCompanyProduct(comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//DELETE COMPANY PRODUCT BY CONTRACT CODE
func DeleteCompanyProductByContractCode(c echo.Context) error {
	// Get Parameter
	contractcode := c.Param("contractcode")
	//get product_id and com_id by contract code
	comid, productid, _ := FindDeleteInfoCompanyProduct(contractcode)
	fmt.Println(comid)
	fmt.Println(productid)
	//
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	err = profileRepository.DeleteCompanyProductContractCode(contractcode)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//Xoá document collection company_role_function khi xoá productid
	err = DeleteCompanyRoleFunctionRelationshipCompanyProduct(productid, comid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP CONTRACT STATUS"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//UPDATE COMPANY PRODUCT BY CONTRACT CODE
func UpdateCompanyProductByContractCode(c echo.Context) error {
	companyProduct := new(model.CompanyProduct)
	err := c.Bind(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	//err = profileRepository.UpdateCompanyProduct(companyProduct.Contractcode, companyProduct)
	//Check product integrated
	count, err2 := CheckProductIntegratedElementExist("productid", companyProduct.Productid)
	fmt.Println(count)
	if err2 != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK PRODUCTID"})
	}
	//Check product integrated
	if count < 1 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Product Integrated not found"})
	}
	err = profileRepository.UpdateCompanyProduct(companyProduct.Contractcode, companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	//update company_role_function
	err = UpdateCompanyRoleFunctionContractStatus(companyProduct.Productid, companyProduct.Contractstatus)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED UPDATE RELATIONSHIP CONTRACT STATUS"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//UPDATE COMPANY PRODUCT STATUS BY COMID
func UpdateCompanyProductStatusByComId(c echo.Context) error {
	companyProduct := new(model.CompanyProduct)
	err := c.Bind(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	err = profileRepository.UpdateCompanyProductStatus(companyProduct.Comid, companyProduct.Contractstatus)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//paging
func PagingCompanyProduct(c echo.Context) error {
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
	companyProducts := new(model.CompanyProducts)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	companyProducts, err = profileRepository.PagingCompanyProductDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": companyProducts})
}

//paging with filter
func PagingAndFilterCompanyProduct(c echo.Context) error {
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
	companyProducts := new(model.CompanyProducts)
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	companyProducts, err = profileRepository.PagingCompanyProductAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": companyProducts})
}

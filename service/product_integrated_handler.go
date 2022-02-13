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

func CreateProductIntegrated(c echo.Context) error {
	productIntegrated := new(model.ProductIntegrated)
	err := c.Bind(productIntegrated)
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
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	_, err = profileRepository.FindProductIntegrated(productIntegrated.Productid)
	if err != mgo.ErrNotFound {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Productid integrated is Exists"})
	}
	err = profileRepository.SaveProductIntegrated(productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//UPDATE PRODUCT INTEGRATED
func UpdateProductIntegratedByProductId(c echo.Context) error {
	productIntegrated := new(model.ProductIntegrated)
	err := c.Bind(productIntegrated)
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
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	err = profileRepository.UpdateProductIntegrated(productIntegrated.Productid, productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//GET INFO PRODUCT_INTEGRATED
func GetProductIntegratedInfo(c echo.Context) error {
	//Sys log
	sysLog := new(model.LogELK)
	//Sys log
	// Get info
	productid := c.Param("productid")
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//Sys log
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		//Sys log
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	productIntegrated, err := profileRepository.FindProductIntegrated(productid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Productid integrated is NonExists"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "1", "message": productIntegrated})
}

//DELETE PRODUCT INTEGRATED OLD
/*func DeleteProductIntegratedByProductId(c echo.Context) error {
	// Get Parameter
	productid := c.Param("productid")
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
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	err = profileRepository.DeleteProductIntegrated(productid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}*/

//DELETE PRODUCT INTEGRATED NEW RELATIONSHIP
func DeleteProductIntegratedByProductId(c echo.Context) error {
	// Get Parameter
	productid := c.Param("productid")
	err := RemoveAllRelatedProductId(productid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})
}

//GET ALL INTEGRATED PRODUCT
func GetAllProductIntegrated(c echo.Context) error {
	// Get info
	//all := c.Param("all")
	//fmt.Println(all)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "10", "message": "Error connect MongoDB", "data": map[string]interface{}{"info": nil}})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	productIntegrateds, err := profileRepository.FindProductIntegratedAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "4", "message": "No data", "data": map[string]interface{}{"info": nil}})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": productIntegrateds}})
}

//paging
func PagingProductIntegrated(c echo.Context) error {
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
	productIntegrateds := new(model.ProductIntegrateds)
	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	productIntegrateds, err = profileRepository.PagingProductIntegratedDriver(paging.Start, paging.Stop)
	countDocument, err := profileRepository.CountDocument()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocument, "data": productIntegrateds})
}

//paging with filter
func PagingAndFilterProductIntegrated(c echo.Context) error {
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
	productIntegrateds := new(model.ProductIntegrateds)
	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	productIntegrateds, err = profileRepository.PagingProductIntegratedAndFilterDriver(pagingandFilter.Filterkey, pagingandFilter.Filtervalue, pagingandFilter.Start, pagingandFilter.Stop)
	countDocumentWithFilter, err := profileRepository.CountDocumentWithFilter(pagingandFilter.Filterkey, pagingandFilter.Filtervalue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "numDocument": countDocumentWithFilter, "data": productIntegrateds})
}

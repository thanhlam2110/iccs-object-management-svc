package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

func StoreCreateCompanyProduct(message string) {
	var status string
	companyProduct := new(model.CompanyProduct)
	json.Unmarshal([]byte(message), &companyProduct)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	err = profileRepository.SaveCompanyProduct(companyProduct)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-COMPANY_PRODUCT-" + string(companyProduct.Comid) + "-" + string(companyProduct.Productid) + "-" + string(companyProduct.Contractcode))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-COMPANY_PRODUCT-" + string(companyProduct.Comid) + "-" + string(companyProduct.Productid) + "-" + string(companyProduct.Contractcode)
	if err != nil {
		status = err.Error()
	} else {
		status = "SAVED COMPANY_PRODUCT " + companyProduct.Comid + "-" + companyProduct.Productid + " MongoDB Success"
	}
	fmt.Println(status)
	//status := "Saved MongoDB success"
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Saved status success")
	session.Close()
}
func StoreUpdateCompanyProduct(message string) {
	var status string
	companyProduct := new(model.CompanyProduct)
	json.Unmarshal([]byte(message), &companyProduct)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	err = profileRepository.UpdateCompanyProduct(companyProduct.Contractcode, companyProduct)
	if err != nil {
		fmt.Println(err)
	}

	//update company_role_function
	err = service.UpdateCompanyRoleFunctionContractStatus(companyProduct.Productid, companyProduct.Contractstatus)
	if err != nil {
		fmt.Println(err.Error() + " FAILED UPDATE RELATIONSHIP CONTRACT STATUS")
	}
	session.Close()
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-COMPANY_PRODUCT-" + string(companyProduct.Contractcode))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-COMPANY_PRODUCT-" + string(companyProduct.Contractcode)
	if err != nil {
		status = err.Error()
	} else {
		status = "UPDATE COMPANY_PRODUCT- " + companyProduct.Contractcode + " MongoDB Success"
	}
	fmt.Println(status)
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Update status success")
	session.Close()
}
func StoreDeleteCompanyProduct(message string) {
	var status string
	//get product_id and com_id by contract code
	comid, productid, _ := service.FindDeleteInfoCompanyProduct(message)
	fmt.Println(comid)
	fmt.Println(productid)
	//
	//connect Mongo
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	err = profileRepository.DeleteCompanyProductContractCode(message)
	if err != nil {
		fmt.Println(err)
	}
	//Xoá document collection company_role_function khi xoá productid
	err = service.DeleteCompanyRoleFunctionRelationshipCompanyProduct(productid, comid)
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED DELETE RELATIONSHIP CONTRACT STATUS")
	}
	session.Close()
	///
	//// Delete status cũ
	err = KafkaDeleteStatus("DELETE-COMPANY_PRODUCT-" + message)
	if err != nil {
		fmt.Println(err)
	}
	////Insert status mới
	key := "DELETE-COMPANY_PRODUCT-" + message
	if err != nil {
		status = err.Error()
	} else {
		status = "Delete Company Product " + message + " MongoDB Success"
	}
	fmt.Println(status)
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Updated status success")
}

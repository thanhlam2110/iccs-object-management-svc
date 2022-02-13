package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

func StoreCreateCompanyRoleFunction(message string) {
	var status string
	companyRoleFunction := new(model.CompanyRoleFunction)
	json.Unmarshal([]byte(message), &companyRoleFunction)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	err = profileRepository.SaveCompanyRoleFunction(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-COMPANY_ROLE_FUNCTION-" + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-COMPANY_ROLE_FUNCTION-" + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "SAVED COMPANY_ROLE_FUNCTION " + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid) + " MongoDB Success"
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
func StoreUpdateCompanyRoleFunction(message string) {
	var status string
	companyRoleFunction := new(model.CompanyRoleFunction)
	json.Unmarshal([]byte(message), &companyRoleFunction)
	/*_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}*/
	//
	err := service.UpdateArrayFunctionCodeList(companyRoleFunction.Rolecode, companyRoleFunction.Functioncodelist)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Update success")
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-COMPANY_ROLE_FUNCTION-" + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-COMPANY_ROLE_FUNCTION-" + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "UPDATE COMPANY_ROLE_FUNCTION " + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid) + " MongoDB Success"
	}
	fmt.Println(status)
	//status := "Saved MongoDB success"
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated status success")
	//session.Close()
}
func StoreDeleteCompanyRoleFunction(message string) {
	var status string
	//connect Mongo
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED CONNECT MONGODB")
	}
	////
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	err = profileRepository.DeleteCompanyRoleFunction(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED DELETE COMPANY ROLE FUNCTION")
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("DELETE-COMPANY-ROLE-FUNCTION-" + message)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "DELETE-COMPANY-ROLE-FUNCTION-" + message
	if err != nil {
		status = err.Error()
	} else {
		status = "DELETE-COMPANY-ROLE-FUNCTION- " + message + " MongoDB Success"
	}
	fmt.Println(status)
	//status := "Saved MongoDB success"
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated status success")
	session.Close()

}

package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

func StoreCreateCompanyFunction(message string) {
	var status string
	companyFunction := new(model.CompanyFunction)
	json.Unmarshal([]byte(message), &companyFunction)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	err = profileRepository.SaveCompanyFunction(companyFunction)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-COMPANY-FUNCTION-" + companyFunction.Functioncode)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-COMPANY-FUNCTION-" + string(companyFunction.Functioncode)
	if err != nil {
		status = err.Error()
	} else {
		status = "Saved COMPANY-FUNCTION " + companyFunction.Functioncode + " MongoDB Success"
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
func StoreUpdateCompanyFunction(message string) {
	var status string
	companyFunction := new(model.CompanyFunction)
	json.Unmarshal([]byte(message), &companyFunction)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//Update
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	err = profileRepository.UpdateCompanyFunction(companyFunction.Functioncode, companyFunction)
	if err != nil {
		fmt.Println(err)
	}
	if companyFunction.Functionstatus == "DISABLE" {
		var functionCodeArray []string
		functionCodeArray = append(functionCodeArray, companyFunction.Functioncode)
		err = service.DeleteArrayFunctionCodeCompanyRoleFunction(functionCodeArray)
		functionCodeArray = nil
		if err != nil {
			fmt.Println(err)
		}
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-COMPANY-FUNCTION-" + companyFunction.Functioncode)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-COMPANY-FUNCTION-" + string(companyFunction.Functioncode)
	if err != nil {
		status = err.Error()
	} else {
		status = "Update COMPANY-FUNCTION " + companyFunction.Functioncode + " MongoDB Success"
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
func StoreDeleteCompanyFunction(message string) {
	var status string
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_function")
	////delete collection company_role_function when delete company_role_function
	var functionCodeArray []string
	functionCodeArray = append(functionCodeArray, message)
	err = service.DeleteArrayFunctionCodeCompanyRoleFunction(functionCodeArray)
	functionCodeArray = nil
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED DELETE RELATIONSHIP FUNCTION CODE"})
	}
	//
	err = profileRepository.DeleteCompanyFunctionByFunctionCode(message)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("DELETE-COMPANY-FUNCTION-" + message)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "DELETE-COMPANY-FUNCTION-" + message
	if err != nil {
		status = err.Error()
	} else {
		status = "DELETE COMPANY-FUNCTION " + message + " MongoDB Success"
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

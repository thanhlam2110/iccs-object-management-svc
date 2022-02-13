package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
	"gopkg.in/mgo.v2"
)

func StoreCreateCompanyRole(message string) {
	var status string
	companyRole := new(model.CompanyRole)
	json.Unmarshal([]byte(message), &companyRole)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	err = profileRepository.SaveCompanyRole(companyRole)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Saved success")
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-COMPANY_ROLE-" + string(companyRole.Rolecode) + "-" + string(companyRole.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-COMPANY_ROLE-" + string(companyRole.Rolecode) + "-" + string(companyRole.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "SAVED COMPANY_ROLE " + string(companyRole.Rolecode) + "-" + string(companyRole.Comid) + " MongoDB Success"
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
func StoreUpdateCompanyRole(message string) {
	var status string
	companyRole := new(model.CompanyRole)
	json.Unmarshal([]byte(message), &companyRole)
	////
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	err = profileRepository.UpdateCompanyRoleCode(companyRole.Rolecode, companyRole)
	if err != nil {
		fmt.Println(err)
	}
	if companyRole.Rolestatus == "DISABLE" {
		var roleCodeArray []string
		roleCodeArray = append(roleCodeArray, companyRole.Rolecode)
		err = service.DeleteArrayRoleCodeCompanyUserRole(roleCodeArray)
		roleCodeArray = nil
		if err != nil {
			fmt.Println(err)
		}
	}
	session.Close()
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-COMPANY_ROLE-" + string(companyRole.Rolecode) + "-" + string(companyRole.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	////
	key := "UPDATE-COMPANY_ROLE-" + string(companyRole.Rolecode) + "-" + string(companyRole.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "UPDATE COMPANY_ROLE " + string(companyRole.Rolecode) + "-" + string(companyRole.Comid) + " MongoDB Success"
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
func StoreDeleteCompanyRole(message string) {
	var status string
	//connect Mongo
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED CONNECT MONGODB")
	}
	//delete company role
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role")
	err = profileRepository.DeleteCompanyRole(message)
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED DELETE COMPANY ROLE  MONGODB")
	}
	////delete collection company_user_role when delete company_role
	var roleCodeArray []string
	roleCodeArray = append(roleCodeArray, message)
	err = service.DeleteArrayRoleCodeCompanyUserRole(roleCodeArray)
	roleCodeArray = nil
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED DELETE RELATIONSHIP ROLE CODE")
	}
	//
	////delete collection company_role_function when delete company_role
	profileRepository2 := repository.NewProfileRepositoryMongo(db, "company_role_function")
	_, err = profileRepository2.FindCompanyRoleFunctionRoleCode(message)
	if err != mgo.ErrNotFound {
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Role code exist"})
		err = profileRepository2.DeleteCompanyRoleFunction(message)
		if err != nil {
			fmt.Println(err)
			fmt.Println("FAILED DELETE RELATIONSHIP COMPANY ROLE FUNCTION")
		}
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("DELETE-COMPANY-FUNCTION-" + message)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "DELETE-COMPANY-ROLE-" + message
	if err != nil {
		status = err.Error()
	} else {
		status = "DELETE-COMPANY-ROLE " + message + " MongoDB Success"
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

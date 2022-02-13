package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

func StoreCreateCompanyUserRole(message string) {
	var status string
	companyUserRole := new(model.CompanyUserRole)
	json.Unmarshal([]byte(message), &companyUserRole)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	err = profileRepository.SaveCompanyUserRole(companyUserRole)
	if err != nil {
		fmt.Println(err)
	}
	session.Close()
	fmt.Println("Saved success")
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-COMPANY_USER_ROLE-" + string(companyUserRole.Username) + "-" + string(companyUserRole.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-COMPANY_USER_ROLE-" + string(companyUserRole.Username) + "-" + string(companyUserRole.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "SAVED COMPANY_USER_ROLE " + string(companyUserRole.Username) + "-" + string(companyUserRole.Comid) + " MongoDB Success"
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
func StoreUpdateCompanyUserRole(message string) {
	var status string
	companyUserRole := new(model.CompanyUserRole)
	json.Unmarshal([]byte(message), &companyUserRole)
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	err = service.UpdateCompanyUserRoleRelationship(companyUserRole.Username, companyUserRole.Datecreate, companyUserRole.Rolecode)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Update success")
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-COMPANY_USER_ROLE-" + string(companyUserRole.Username) + "-" + string(companyUserRole.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-COMPANY_USER_ROLE-" + string(companyUserRole.Username) + "-" + string(companyUserRole.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "UPDATE COMPANY_USER_ROLE " + string(companyUserRole.Username) + "-" + string(companyUserRole.Comid) + " MongoDB Success"
	}
	fmt.Println(status)
	//status := "Saved MongoDB success"
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated status success")
	session.Close()
	//
}

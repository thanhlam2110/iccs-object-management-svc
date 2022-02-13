package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

//

func StoreCreateCompanyInfo(message string) {
	var status string
	companyInfo := new(model.CompanyInfo)
	json.Unmarshal([]byte(message), &companyInfo)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	err = profileRepository.SaveCompanyInfo(companyInfo)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-COMPANY-INFO-" + companyInfo.Comid)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-COMPANY-INFO-" + string(companyInfo.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "Saved COMPANY INFO " + companyInfo.Comid + " MongoDB Success"
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
func StoreUpdateCompanyInfo(message string) {
	var status string
	companyInfo := new(model.CompanyInfo)
	json.Unmarshal([]byte(message), &companyInfo)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_info")
	err = profileRepository.UpdateCompanyInfoByComId(companyInfo.Comid, companyInfo)
	if err != nil {
		fmt.Println(err)
	}
	if companyInfo.Comstatus == "DISABLE" {
		fmt.Println(companyInfo.Comid)
		err = service.UpdateCascadeUserStatusCompanyInfo(companyInfo.Comid)
		if err != nil {
			fmt.Println(err)
		}
		err = service.UpdateCascadeContractStatusCompanyInfo(companyInfo.Comid)
		if err != nil {
			fmt.Println(err)
		}
		err = service.UpdateCascadeUserServiceStatusCompanyInfo(companyInfo.Comid)
		if err != nil {
			fmt.Println(err)
		}
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-COMPANY-INFO-" + companyInfo.Comid)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-COMPANY-INFO-" + string(companyInfo.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "Updated COMPANY INFO " + companyInfo.Comid + " MongoDB Success"
	}
	//status = "Updated CompanyInfo MongoDB Success"
	err2 := KafkaCreateStatus(key, status)
	if err != nil {
		fmt.Println(err2)
	}
	fmt.Println("Update status success")
	session.Close()
}
func StoreDeleteCompanyInfo(message string) {
	var status string
	err := service.RemoveAllRelatedCompanyComid(message)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("DELETE-COMPANY_INFO" + message)
	if err != nil {
		fmt.Println(err)
	}
	////Insert status mới
	key := "DELETE-COMPANY_INFO" + message
	if err != nil {
		status = err.Error()
	} else {
		status = "Delete Company Info " + message + " MongoDB Success"
	}
	fmt.Println(status)
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Updated status success")
}

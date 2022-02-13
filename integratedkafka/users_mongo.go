package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

func StoreCreateUsers(message string) {
	var status string
	userSSO := new(model.UserSSO)
	json.Unmarshal([]byte(message), &userSSO)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	err = profileRepository.Save(userSSO)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	session.Close()
	fmt.Println("Saved success")
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "SAVED USERS " + string(userSSO.Username) + "-" + string(userSSO.Comid) + " MongoDB Success"
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
func StoreUpdateCompanyUsers(message string) {
	var status string
	userSSO := new(model.UserSSO)
	json.Unmarshal([]byte(message), &userSSO)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	err = profileRepository.Update(userSSO.Username, userSSO)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	err = service.UpdateRelationshipNode(userSSO.Username, userSSO.Userstatus)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED UPDATE CHILD USER STATUS"})
	}
	//
	fmt.Println("Update success")
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "UPDATE USERS " + string(userSSO.Username) + "-" + string(userSSO.Comid) + " MongoDB Success"
	}
	fmt.Println(status)
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Saved status success")
	session.Close()
}
func StoreUpdateEndUsers(message string) {
	var status string
	userSSO := new(model.UserSSO)
	json.Unmarshal([]byte(message), &userSSO)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	//
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	err = profileRepository.Update(userSSO.Username, userSSO)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//
	fmt.Println("Update success")
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid))
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid)
	if err != nil {
		status = err.Error()
	} else {
		status = "UPDATE USERS " + string(userSSO.Username) + "-" + string(userSSO.Comid) + " MongoDB Success"
	}
	fmt.Println(status)
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println("Saved status success")
	session.Close()
}

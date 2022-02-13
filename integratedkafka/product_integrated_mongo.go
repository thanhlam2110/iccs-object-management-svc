package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
)

//var status string

func StoreCreateProductIntegrated(message string) {
	var status string
	productIntegrated := new(model.ProductIntegrated)
	json.Unmarshal([]byte(message), &productIntegrated)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	err = profileRepository.SaveProductIntegrated(productIntegrated)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("INSERT-" + productIntegrated.Productid)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "INSERT-" + string(productIntegrated.Productid)
	if err != nil {
		status = err.Error()
	} else {
		status = "Saved " + productIntegrated.Productid + " MongoDB Success"
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
func StoreUpdateProductIntegrated(message string) {
	var status string
	productIntegrated := new(model.ProductIntegrated)
	json.Unmarshal([]byte(message), &productIntegrated)
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	err = profileRepository.UpdateProductIntegrated(productIntegrated.Productid, productIntegrated)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("UPDATE-" + productIntegrated.Productid)
	if err != nil {
		fmt.Println(err)
	}
	////
	key := "UPDATE-" + string(productIntegrated.Productid)
	if err != nil {
		status = err.Error()
	} else {
		//status = "Saved ProductIntegrated MongoDB Success"
		status = "Updated " + productIntegrated.Productid + " MongoDB Success"
	}
	fmt.Println(status)
	//status = "Updated " + productIntegrated.Productid + " MongoDB Success"
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Updated status success")
	session.Close()
}
func StoreDeleteProductIntegrated(message string) {
	var status string
	err := service.RemoveAllRelatedProductId(message)
	if err != nil {
		fmt.Println(err)
	}
	//// Delete status cũ
	err = KafkaDeleteStatus("DELETE-" + message)
	if err != nil {
		fmt.Println(err)
	}
	////Insert status mới
	key := "DELETE-" + message
	if err != nil {
		status = err.Error()
	} else {
		status = "Delete " + message + " MongoDB Success"
	}
	fmt.Println(status)
	err2 := KafkaCreateStatus(key, status)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Updated status success")
}

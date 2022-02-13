package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/transport"
	"github.com/labstack/echo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/mgo.v2"
)

func KafkaCreateCompanyProduct(c echo.Context) error {
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
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_product")
	//Check product integrated
	count, err2 := service.CheckProductIntegratedElementExist("productid", companyProduct.Productid)
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
	checkTwice, err3 := service.CheckTwiceComIDandProductIDinCompanyProduct(companyProduct.Comid, companyProduct.Productid)
	if err3 != nil {
		fmt.Println(err3)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK DOUBLE COMID AND PRODUCTID"})
	}
	if checkTwice > 0 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "DON'T ALLOW CREATE DUPLICATE SAME COMID AND PRODUCTID"})
	}
	//Check twice comid + productid

	/*err = profileRepository.SaveCompanyProduct(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")*/
	session.Close()
	//Parse struct to string
	out, err := json.Marshal(companyProduct)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "CREATE"
	kafkaStruct.Collection = "company_product"
	kafkaStruct.Data = string(out)
	message, err := json.Marshal(kafkaStruct)
	if err != nil {
		panic(err)
	}
	//Kafka
	go func() {
		for e := range kafkaProducerCreate.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	topic := "SSO_CREATE_TOPIC"
	err = transport.PublishMessage(kafkaProducerCreate, topic, kafka.PartitionAny, string(message))
	//
	//session.Close()
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Processing"})
}
func GetKafkaCreateCompanyProductStatus(c echo.Context) error {
	companyProduct := new(model.CompanyProduct)
	err := c.Bind(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "INSERT-COMPANY_PRODUCT-" + string(companyProduct.Comid) + "-" + string(companyProduct.Productid) + "-" + string(companyProduct.Contractcode)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}

//update
func KafkaUpdateCompanyProductByContractCode(c echo.Context) error {
	companyProduct := new(model.CompanyProduct)
	err := c.Bind(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//check connect
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	count, err2 := service.CheckProductIntegratedElementExist("productid", companyProduct.Productid)
	fmt.Println(count)
	if err2 != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CHECK PRODUCTID"})
	}
	//Check product integrated
	if count < 1 {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Product Integrated not found"})
	}
	//
	session.Close()
	//Parse struct to string
	out, err := json.Marshal(companyProduct)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "UPDATE"
	kafkaStruct.Collection = "company_product"
	kafkaStruct.Data = string(out)
	message, err := json.Marshal(kafkaStruct)
	if err != nil {
		panic(err)
	}
	//Kafka
	go func() {
		for e := range kafkaProducerUpdate.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	topic := "SSO_UPDATE_TOPIC"
	err = transport.PublishMessage(kafkaProducerUpdate, topic, kafka.PartitionAny, string(message))
	//
	//session.Close()
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Processing"})
}
func GetKafkaUpdateCompanyProductStatus(c echo.Context) error {
	companyProduct := new(model.CompanyProduct)
	err := c.Bind(companyProduct)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "UPDATE-COMPANY_PRODUCT-" + string(companyProduct.Contractcode)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}
func KafkaDeleteCompanyProduct(c echo.Context) error {
	contractcode := c.Param("contractcode")
	///////Check CONNECT before UPDATE
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CONNECT DATABASE"})
	}
	session.Close()
	//////
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "DELETE"
	kafkaStruct.Collection = "company_product"
	kafkaStruct.Data = string(contractcode)
	message, err := json.Marshal(kafkaStruct)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(message))
	//Kafka
	go func() {
		for e := range kafkaProducerDelete.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	topic := "SSO_DELETE_TOPIC"
	err = transport.PublishMessage(kafkaProducerDelete, topic, kafka.PartitionAny, string(message))
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Processing"})
}
func GetKafkaDeleteCompanyProductStatus(c echo.Context) error {
	contractcode := c.Param("contractcode")
	key := "DELETE-COMPANY_PRODUCT-" + contractcode
	//fmt.Println(key)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}

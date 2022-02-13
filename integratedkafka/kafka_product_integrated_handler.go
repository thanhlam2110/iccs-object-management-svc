package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/transport"

	"github.com/labstack/echo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/mgo.v2"
)

func KafkaCreateProductIntegrated(c echo.Context) error {
	productIntegrated := new(model.ProductIntegrated)
	err := c.Bind(productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	//db, err := config.GetMongoDB()
	//check before connect
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED CONNECT MONGODB"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "product_integrated")
	_, err = profileRepository.FindProductIntegrated(productIntegrated.Productid)
	if err != mgo.ErrNotFound {
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Productid integrated is Exists"})
	}
	//
	/*err = profileRepository.SaveProductIntegrated(productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Saved success")*/
	session.Close()
	//Parse struct to string
	out, err := json.Marshal(productIntegrated)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "CREATE"
	kafkaStruct.Collection = "product_integrated"
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
func GetKafkaCreateProductIntegratedStatus(c echo.Context) error {
	productIntegrated := new(model.ProductIntegrated)
	err := c.Bind(productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "INSERT-" + string(productIntegrated.Productid)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}
func KafkaUpdateProductIntegrated(c echo.Context) error {
	productIntegrated := new(model.ProductIntegrated)
	err := c.Bind(productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	///////Check CONNECT before UPDATE
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED CONNECT DATABASE"})
	}
	session.Close()
	//////
	//Parse struct to string
	out, err := json.Marshal(productIntegrated)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "UPDATE"
	kafkaStruct.Collection = "product_integrated"
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
func GetKafkaUpdateProductIntegratedStatus(c echo.Context) error {
	productIntegrated := new(model.ProductIntegrated)
	err := c.Bind(productIntegrated)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "UPDATE-" + string(productIntegrated.Productid)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}
func KafkaDeleteProductIntegrated(c echo.Context) error {
	productid := c.Param("productid")
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
	kafkaStruct.Collection = "product_integrated"
	kafkaStruct.Data = string(productid)
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
func GetKafkaDeleteProductIntegratedStatus(c echo.Context) error {
	productid := c.Param("productid")
	key := "DELETE-" + productid
	//fmt.Println(key)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}

package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/transport"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/mgo.v2"
)

func KafkaCreateUser(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}

	/*if err := c.Validate(userSSO); err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}*/

	/*if !service.isStringAlphabetic(userSSO.Username) {
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}*/
	//check length password
	if len(userSSO.Password) < 8 {
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	if userSSO.Useremail == "" && userSSO.Usertel == "" {
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Email empty or Phone empty", "data": map[string]interface{}{"info": nil}})
	}
	countEmail, err2 := service.CheckUserEmailExist(userSSO.Useremail)
	countPhone, err3 := service.CheckUserPhoneExist(userSSO.Usertel)
	if err2 != nil || err3 != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	if countEmail > 0 || countPhone > 0 {
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Email or Phone number exist", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	//connect Mongo
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	_, err = profileRepository.FindByUser(userSSO.Username)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "User is Exists", "data": map[string]interface{}{"info": nil}})
	}
	userSSO.Userid = uuid.New().String()
	//insert
	/*err = profileRepository.Save(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}*/
	//insert
	//kafka
	//Parse struct to string
	out, err := json.Marshal(userSSO)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "CREATE"
	kafkaStruct.Collection = "users"
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
	session.Close()
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Processing"})
}
func GetKafkaCreateUserStatus(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "INSERT-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid)
	val, err := KafkaReadStatus(key)
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "9", "status": err.Error()})
	}
	if val == "" {
		return c.JSON(200, map[string]interface{}{"code": "4", "status": "key does not exist"})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}
func KafkaUpdateUser(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	if len(userSSO.Password) < 8 {
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	//connect Mongo
	//db, err := config.GetMongoDB()
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	//kafka
	//Parse struct to string
	out, err := json.Marshal(userSSO)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "UPDATE"
	kafkaStruct.Collection = "users"
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
	topic := "SSO_UPDATE_TOPIC"
	err = transport.PublishMessage(kafkaProducerCreate, topic, kafka.PartitionAny, string(message))
	//
	session.Close()
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Processing"})
}
func GetKafkaUpdateUserStatus(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "UPDATE-USERS-" + string(userSSO.Username) + "-" + string(userSSO.Comid)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}

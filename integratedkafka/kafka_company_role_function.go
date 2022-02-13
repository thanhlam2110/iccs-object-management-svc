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

func KafkaCreateCompanyRoleFunction(c echo.Context) error {
	companyRoleFunction := new(model.CompanyRoleFunction)
	err := c.Bind(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	/////////
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "company_role_function")
	_, err = profileRepository.FindCompanyRoleFunctionRoleCode(companyRoleFunction.Rolecode)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "Role code exist"})
	}

	session.Close()
	//Parse struct to string
	out, err := json.Marshal(companyRoleFunction)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "CREATE"
	kafkaStruct.Collection = "company_role_function"
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
func GetKafkaCreateCompanyRoleFunctionStatus(c echo.Context) error {
	companyRoleFunction := new(model.CompanyRoleFunction)
	err := c.Bind(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "INSERT-COMPANY_ROLE_FUNCTION-" + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}
func KafkaUpdateCompanyRoleFunction(c echo.Context) error {
	////
	companyRoleFunction := new(model.CompanyRoleFunction)
	err := c.Bind(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	_, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	/*err = UpdateArrayFunctionCodeList(companyRoleFunction.Rolecode, companyRoleFunction.Functioncodelist)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "-1", "message": "FAILED"})
	}
	fmt.Println("Update success")
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "SUCCESS"})*/
	//////
	session.Close()
	//Parse struct to string
	out, err := json.Marshal(companyRoleFunction)
	if err != nil {
		panic(err)
	}
	kafkaStruct := new(model.KafkaStruct)
	kafkaStruct.Action = "UPDATE"
	kafkaStruct.Collection = "company_role_function"
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
func GetKafkaUpdateCompanyRoleFunctionStatus(c echo.Context) error {
	companyRoleFunction := new(model.CompanyRoleFunction)
	err := c.Bind(companyRoleFunction)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "-1", "message": "Body is Invalid"})
	}
	key := "UPDATE-COMPANY_ROLE_FUNCTION-" + string(companyRoleFunction.Rolecode) + "-" + string(companyRoleFunction.Productid) + "-" + string(companyRoleFunction.Comid)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}
func KafkaDeleteCompanyRoleFunction(c echo.Context) error {
	rolecode := c.Param("rolecode")
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
	kafkaStruct.Collection = "company_role_function"
	kafkaStruct.Data = string(rolecode)
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
func GetKafkaDeleteCompanyRoleFunctionStatus(c echo.Context) error {
	rolecode := c.Param("rolecode")
	key := "DELETE-COMPANY-ROLE-FUNCTION-" + rolecode
	//fmt.Println(key)
	val, err := KafkaReadStatus(key)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "status": val})
}

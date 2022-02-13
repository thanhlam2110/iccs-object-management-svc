package integratedkafka

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"

	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var collection string

func HandlerConsumer() {
	fmt.Println("-------------------Init Consumer---------------------")
	config.ReadConfig()
	url := viper.GetString(`kafka.url`)
	groupID := viper.GetString(`kafka.groupID`)
	autoOffsetReset := viper.GetString(`kafka.autoOffsetReset`)
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": url,
		"group.id":          groupID,
		"auto.offset.reset": autoOffsetReset,
	})

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"SSO_CREATE_TOPIC", "SSO_UPDATE_TOPIC", "SSO_DELETE_TOPIC"}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			//fmt.Printf("Nhận Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			ReveiveMessageFromKafka(string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
func ReveiveMessageFromKafka(message string) {
	fmt.Println("Nhận message " + message)
	kafkaStruct := new(model.KafkaStruct)
	json.Unmarshal([]byte(message), &kafkaStruct)
	fmt.Println("=========================================================")
	if kafkaStruct.Action == "CREATE" {
		fmt.Println("GỌI HÀM CREATE")
		fmt.Println(kafkaStruct.Collection)
		collection := string(kafkaStruct.Collection)
		switch collection {
		case "company_info":
			fmt.Println("create_company_info")
			StoreCreateCompanyInfo(kafkaStruct.Data)
		case "product_integrated":
			fmt.Println("product_integrated")
			StoreCreateProductIntegrated(kafkaStruct.Data)
		case "company_product":
			fmt.Println("company_product")
			StoreCreateCompanyProduct(kafkaStruct.Data)
		case "company_function":
			fmt.Println("company_function")
			StoreCreateCompanyFunction(kafkaStruct.Data)
		case "company_role":
			fmt.Println("company_role")
			StoreCreateCompanyRole(kafkaStruct.Data)
		case "company_role_function":
			fmt.Println("company_role_function")
			StoreCreateCompanyRoleFunction(kafkaStruct.Data)
		case "company_user_role":
			fmt.Println("company_user_role")
			StoreCreateCompanyUserRole(kafkaStruct.Data)
		case "users":
			fmt.Println("users")
			StoreCreateUsers(kafkaStruct.Data)
		}
		fmt.Println(kafkaStruct.Data)
	} else if kafkaStruct.Action == "UPDATE" {
		fmt.Println("GỌI HÀM UPDATE")
		fmt.Println(kafkaStruct.Collection)
		collection := string(kafkaStruct.Collection)
		switch collection {
		case "company_info":
			fmt.Println("update_company_info")
			StoreUpdateCompanyInfo(kafkaStruct.Data)
		case "product_integrated":
			fmt.Println("product_integrated")
			StoreUpdateProductIntegrated(kafkaStruct.Data)
		case "company_product":
			fmt.Println("company_product")
			StoreUpdateCompanyProduct(kafkaStruct.Data)
		case "company_function":
			fmt.Println("company_function")
			StoreUpdateCompanyFunction(kafkaStruct.Data)
		case "company_role":
			fmt.Println("company_role")
			StoreUpdateCompanyRole(kafkaStruct.Data)
		case "company_role_function":
			fmt.Println("company_role_function")
			StoreUpdateCompanyRoleFunction(kafkaStruct.Data)
		case "company_user_role":
			fmt.Println("company_user_role")
			StoreUpdateCompanyUserRole(kafkaStruct.Data)
		case "users":
			fmt.Println("users")
			data := kafkaStruct.Data
			userSSO := new(model.UserSSO)
			json.Unmarshal([]byte(data), &userSSO)
			userType := userSSO.Usertype
			fmt.Println(userType)
			if userType == "ADMIN_ENDUSER" {
				StoreUpdateEndUsers(kafkaStruct.Data)
			} else {
				StoreUpdateCompanyUsers(kafkaStruct.Data)
			}
		}
		fmt.Println(kafkaStruct.Data)
	} else if kafkaStruct.Action == "DELETE" {
		fmt.Println("GỌI HÀM DELETE")
		fmt.Println(kafkaStruct.Collection)
		collection := string(kafkaStruct.Collection)
		switch collection {
		case "product_integrated":
			fmt.Println("product_integrated")
			StoreDeleteProductIntegrated(kafkaStruct.Data)
		case "company_info":
			fmt.Println("company_info")
			StoreDeleteCompanyInfo(kafkaStruct.Data)
		case "company_product":
			fmt.Println("company_product")
			StoreDeleteCompanyProduct(kafkaStruct.Data)
		case "company_function":
			fmt.Println("company_function")
			StoreDeleteCompanyFunction(kafkaStruct.Data)
		case "company_role":
			fmt.Println("company_role")
			StoreDeleteCompanyRole(kafkaStruct.Data)
		case "company_role_function":
			fmt.Println("company_role_function")
			StoreDeleteCompanyRoleFunction(kafkaStruct.Data)
		}
	} else {
		fmt.Println("KHÁC")
	}
	fmt.Println("=========================================================")
}

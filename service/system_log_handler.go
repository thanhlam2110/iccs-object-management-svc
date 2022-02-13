package service

import (
	"context"
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"github.com/spf13/viper"
	elastic "gopkg.in/olivere/elastic.v7"
)

func WriteLog(log *model.LogELK) {
	//Read config
	config.ReadConfig()
	index := viper.GetString(`elk.index`)
	//fmt.Println(index)
	//Read config
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	dataJSON, err := json.Marshal(log)
	fmt.Println(string(dataJSON))
	js := string(dataJSON)
	_, err = esclient.Index().
		Index(index).
		BodyJson(js).
		Do(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println("[Elastic][InsertProduct]Insertion Successful")
}
func GetESClient() (*elastic.Client, error) {
	//Read config
	config.ReadConfig()
	link := viper.GetString(`elk.url`)
	//fmt.Println(link)
	//fmt.Printf("%T\n", link)
	//Read config
	client, err := elastic.NewClient(elastic.SetURL(link),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	fmt.Println("ES initialized...")
	return client, err

}

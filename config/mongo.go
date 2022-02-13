package config

import (
	"fmt"

	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

func GetMongoDB() (*mgo.Database, error) {
	//hàm cũ không release connection
	ReadConfig()
	//standalone
	link := viper.GetString(`mongo.link`)
	fmt.Println(link)
	//cluster
	//link := viper.GetString(`mongocluster.link`)
	//
	//session, err := mgo.Dial("mongodb://casuser:Mellon@203.162.141.35:27017/users")
	session, err := mgo.Dial(link)
	if err != nil {
		return nil, err
	}
	db := session.DB("users")
	return db, nil
}
func ReadConfig() {
	// Get file config
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func GetMongoDataBase() (*mgo.Database, *mgo.Session, error) {
	//
	ReadConfig()
	//link := viper.GetString(`mongo.url`)
	link := viper.GetString(`mongocluster.link`)
	session, err := mgo.Dial(link)
	if err != nil {
		return nil, nil, err
	}
	db := session.DB("users")
	return db, session, nil
}
func GetMongoDataBaseCas() (*mgo.Database, *mgo.Session, error) {
	//
	ReadConfig()
	link := viper.GetString(`mongo.url`)
	//link := viper.GetString(`mongocluster.link`)
	session, err := mgo.Dial(link)
	if err != nil {
		return nil, nil, err
	}
	db := session.DB("cas")
	return db, session, nil
}

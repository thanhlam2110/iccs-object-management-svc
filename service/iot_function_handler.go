package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

//
//redis
var ctxRedis = context.Background()

func RequestSSOIoToken(c echo.Context) error {
	//read body
	tokenRequestBody := new(model.TokenRequestBody)
	err := c.Bind(tokenRequestBody)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"ssotoken": nil}})
	}
	username := tokenRequestBody.Username
	password := tokenRequestBody.Password
	//request SSO JWT token
	token, refreshtoken := RequestToken(username, password)
	if token == "" {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "Authen Failed"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Authen Failed", "data": map[string]interface{}{"ssotoken": nil, "ssorefreshtoken": nil}})
	}
	//parse token ---> check user status
	authResponse, err := BasicAuth(token)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"ssotoken": nil, "ssorefreshtoken": nil}})
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)
	attributes := result["attributes"].(map[string]interface{})
	attributes["username"] = result["id"]
	userstatus := attributes["userstatus"]
	if userstatus == "DISABLE" {
		//Sys log
		SysLog.Datetime = time.Now()
		SysLog.Log = "USER DISABLE"
		WriteLog(SysLog)
		//Sys log
		return c.JSON(200, map[string]interface{}{"code": "9", "message": "USER DISABLE", "data": map[string]interface{}{"ssotoken": nil, "ssorefreshtoken": nil}})
	}
	//check IOT user
	CreateIoTUser((attributes["useremail"]).(string), password)
	//fmt.Println(createStatus)
	//Create IOT user
	//request iot authen key
	ioTAuthKey, status := RequestIoToken((attributes["useremail"]).(string), password)
	//fmt.Println(ioTAuthKey)
	if status != "201 Created" {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = status
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "IOT " + status, "data": map[string]interface{}{"ssotoken": nil, "ssorefreshtoken": nil}})
	}
	//INSERT REDIS
	//connect redis
	config.ReadConfig()
	url := viper.GetString(`redis.url`)
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "ReDis0rimx", // no password set
		DB:       0,            // use default DB
	})
	_, err = rdb.Ping(ctxRedis).Result()
	//insert to redis
	// Output: PONG <nil>
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "Error connect redis"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Error connect redis", "data": map[string]interface{}{"ssotoken": nil, "ssorefreshtoken": nil}})
	}
	err = rdb.SetNX(ctxRedis, username, ioTAuthKey, 7200*time.Second).Err()
	//fmt.Println(err)
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "Error insert data redis"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Error insert data redis", "data": map[string]interface{}{"ssotoken": nil, "ssorefreshtoken": nil}})
	}
	//
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"ssotoken": token, "ssorefreshtoken": refreshtoken}})
}

func GetIoTAuthKeyFromRedis(c echo.Context) error {
	authenRequestBody := new(model.AuthenRequestBody)
	err := c.Bind(authenRequestBody)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	token := authenRequestBody.Token
	authResponse, err := BasicAuth(token)
	//
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)

	if result["error"] != nil {
		if result["message"] != nil {
			//syslog
			SysLog.Datetime = time.Now()
			SysLog.Log = (result["message"]).(string)
			WriteLog(SysLog)
			//syslog
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = (result["error"]).(string)
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	//test
	username := result["id"]
	//connect redis
	config.ReadConfig()
	url := viper.GetString(`redis.url`)
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "ReDis0rimx", // no password set
		DB:       0,            // use default DB
	})
	val, err := rdb.Get(ctxRedis, username.(string)).Result()
	if err == redis.Nil {
		//fmt.Println("key does not exist")
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "Redis key does not exist"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "1", "message": "AuthenKey không tồn tại", "data": map[string]interface{}{"info": nil}})
	} else if err != nil {
		//panic(err)
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "System error", "data": map[string]interface{}{"info": nil}})
	} else {
		//fmt.Println(username, val)
		return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"username": username, "IoTAuthenKey": val}})
	}
}

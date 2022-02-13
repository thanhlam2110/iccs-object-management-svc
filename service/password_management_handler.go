package service

import (
	"fmt"
	"regexp"
	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/repository"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func RequestForgotUserName(c echo.Context) error {
	inputUserInfo := new(model.ForgotUsername)
	err := c.Bind(inputUserInfo)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	info := inputUserInfo.InputUserInfo
	//connect database
	db, err := config.GetMongoDB()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//
	if isEmailValid(info) {
		//lay thong tin username
		profileRepository := repository.NewProfileRepositoryMongo(db, "users")
		userSSO, err := profileRepository.FindUsernameByEmail(info)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "4", "message": "Can't find username", "data": map[string]interface{}{"info": nil}})
		}
		//lay thong tin username
		fmt.Println("Send username via email " + info)
		response := "Send username via email " + info
		SendEmail(info, "Thông tin username", "Username: "+userSSO.Username)
		return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": response}})
	} else {
		//lay thong tin username
		profileRepository := repository.NewProfileRepositoryMongo(db, "users")
		userSSO, err := profileRepository.FindUsernameByPhone(info)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "4", "message": "Can't find username", "data": map[string]interface{}{"info": nil}})
		}
		//lay thong tin username
		fmt.Println("Send username via phone number " + info)
		response := "Send username via phone number " + info
		SendSms(info, "Username: "+userSSO.Username)
		return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": response}})
	}
	//return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": nil}})
}

//reset password
func RequestResetPassword(c echo.Context) error {
	requestResetPassword := new(model.RequestResetPassword)
	err := c.Bind(requestResetPassword)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"otp": nil}})
	}
	username := requestResetPassword.Username
	//create random otp
	otp := RandomOtpString(8)
	//
	//Store {"otp":"username"} --> redis
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
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Error connect redis", "data": map[string]interface{}{"otp": nil}})
	}
	err = rdb.SetNX(ctxRedis, otp, username, 300*time.Second).Err()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "Error insert data redis"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Error insert data redis", "data": map[string]interface{}{"otp": nil}})
	}
	//
	// Find Email or Phone by username
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Error connect MongoDB", "data": map[string]interface{}{"otp": nil}})
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	userSSO, err := profileRepository.FindByUser(username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "4", "message": "Username Non-Exist", "data": map[string]interface{}{"otp": nil}})
	}
	phone := userSSO.Usertel
	email := userSSO.Useremail
	//fmt.Println(phone)
	//fmt.Println(email)
	if phone != "" {
		SendSms(phone, "Your OTP reset password is : "+otp)
	} else if email != "" {
		fmt.Println("Send OTP via email")
		SendEmail(email, "Thông tin OTP reset password", "Your OTP reset password is: "+otp)
	} else {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "System error", "data": map[string]interface{}{"otp": nil}})
	}
	//
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"otp": otp}})
}
func ResetPassword(c echo.Context) error {
	resetPassword := new(model.ResetPassword)
	err := c.Bind(resetPassword)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	otp := resetPassword.Otp
	newPassword := resetPassword.Newpassword
	//check length password
	if len(newPassword) < 8 {
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	//check length password
	//connect redis get username
	//connect redis
	config.ReadConfig()
	url := viper.GetString(`redis.url`)
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "ReDis0rimx", // no password set
		DB:       0,            // use default DB
	})
	username, err := rdb.Get(ctxRedis, otp).Result()
	if err == redis.Nil {
		//fmt.Println("key does not exist")
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = "OTP Reset Password không hợp lệ"
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "1", "message": "OTP Reset Password không hợp lệ", "data": map[string]interface{}{"info": nil}})
	} else if err != nil {
		//panic(err)
		//syslog
		SysLog.Datetime = time.Now()
		SysLog.Log = err.Error()
		WriteLog(SysLog)
		//syslog
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "System error", "data": map[string]interface{}{"info": nil}})
	} else {
		//update new password
		err = UpdatePasswordByUsername(username, newPassword)
		if err != nil {
			fmt.Println(err)
			return c.JSON(200, map[string]interface{}{"code": "10", "message": "FAIL UPDATE PASSWORD", "data": map[string]interface{}{"info": nil}})
		}
		return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": nil}})
	}

	//
	//return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": nil}})
}

//check valid email
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

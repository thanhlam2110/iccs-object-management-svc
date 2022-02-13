package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

//

type Attributes struct {
	ComDepartment string `json:"comdepartment"`
	ComId         string `json:"comid"`
	Lastname      string `json:"lastname"`
	UserCode      string `json:"usercode"`
	UserDate      string `json:"userdate"`
	UserEmail     string `json:"useremail"`
	UserStatus    string `json:"userstatus"`
	UserTel       string `json:"usertel"`
	UserType      string `json:"usertype"`
	UserParentid  string `json:"userparentid"`
}
type Tokens struct {
	Accesstoken  string `json:"access_token"`
	Refreshtoken string `json:"refresh_token"`
	Tokentype    string `json:"token_type"`
	Expiresin    string `json:"expires_in"`
	Scope        string `json:"scope"`
}
type RTokens struct {
	Accesstoken string `json:"access_token"`
	Tokentype   string `json:"token_type"`
	Expiresin   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

func BasicAuth(token string) (s string, err error) {
	config.ReadConfig()
	link := viper.GetString(`sso.url`)
	var username string = "exampleOauthClient"
	var passwd string = "exampleOauthClientSecret"
	client := &http.Client{}
	url := link + "/profile?access_token=" + token
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s = string(bodyText)
	return s, nil
}
func RequestToken(username string, password string) (string, string) {
	//readConfig()
	config.ReadConfig()
	link := viper.GetString(`sso.url`)
	url := link + "/token?grant_type=password&client_id=exampleOauthClient&client_secret=exampleOauthClientSecret&username=" + username + "&password=" + password
	req, _ := http.NewRequest("GET", url, nil)
	//res, _ := http.DefaultClient.Do(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		fmt.Println(err)
		return "", ""
		//Sys log
	}
	body, _ := ioutil.ReadAll(res.Body)
	stringBody := string(body)
	if stringBody != "" {
		var tokens Tokens
		json.Unmarshal([]byte(stringBody), &tokens)
		return tokens.Accesstoken, tokens.Refreshtoken
	} else {
		return "", ""
	}
}
func RequestTokenv2(client_id string, client_secret string, username string, password string) (string, string) {
	//readConfig()
	config.ReadConfig()
	link := viper.GetString(`sso.url`)
	url := link + "/token?grant_type=password&client_id=" + client_id + "&client_secret=" + client_secret + "&username=" + username + "&password=" + password
	req, _ := http.NewRequest("GET", url, nil)
	//res, _ := http.DefaultClient.Do(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		fmt.Println(err)
		return "", ""
		//Sys log
	}
	body, _ := ioutil.ReadAll(res.Body)
	stringBody := string(body)
	if stringBody != "" {
		var tokens Tokens
		json.Unmarshal([]byte(stringBody), &tokens)
		return tokens.Accesstoken, tokens.Refreshtoken
	} else {
		return "", ""
	}
}

//GetCompanyUserRole
//ham kien truc cu
/*func GetCompanyUserRole(username string, comid string) (arrayRole []string) {
	//connect Mongo
	db, err := config.GetMongoDB()
	if err != nil {
		fmt.Println(err)
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "company_user_role")
	companyUserRole, err := profileRepository.FindRoleByUser(username)
	if err != nil {
		fmt.Println(err)
	}
	arrayRole = strings.Split(companyUserRole.Rolecode, ",")
	return arrayRole
}*/

func GetUserRoleWithProductID(username string, comid string, productid string) (results []string, err error) {
	//db, err := config.GetMongoDB()
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println(err)
	}
	c := db.C("company_role_function")
	//var results []string
	//roleArray := []string{"DHBK_ROLE_01", "DHBK_ROLE_03"}-->[DHBK_ROLE_01 DHBK_ROLE_03]
	err, roleArray := GetCompanyUserRole(username, comid)
	if err != nil {
		//panic(err)
		return results, err
	}
	err = c.Find(bson.M{"rolecode": bson.M{"$in": roleArray}, "productid": productid, "contractstatus": "ACTIVE"}).Distinct("rolecode", &results)
	if err != nil {
		//panic(err)
		return results, err
	}
	//fmt.Println(results[0])
	/*err, responseArr := GetCompanyUserRoleByRoleCode(results[0])
	fmt.Println(responseArr)*/
	/*for i := 0; i < len(results); i++ {
		err, responseArr := GetCompanyUserRoleByRoleCode(results[i])
		if err != nil {
			panic(err)
		}
		//fmt.Println(responseArr)
	}*/
	session.Close()
	return results, err
}

//IOT
type IoToken struct {
	Token string `json:"token"`
}

func RequestIoToken(email string, password string) (token string, status string) {
	config.ReadConfig()
	//url := viper.GetString(`iot.url`) + "tokens"
	url := viper.GetString(`iot.urltokens`)
	stringAuth := `{"email":"` + email + `", "password":"` + password + `"}`
	var jsonStr = []byte(stringAuth)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	var ioToken IoToken
	json.Unmarshal(body, &ioToken)
	/*if resp.Status != "201 Created" {
		return resp.Status
	}*/
	//fmt.Println(ioToken.Token)
	return ioToken.Token, resp.Status
}
func CreateIoTUser(email string, password string) string {
	config.ReadConfig()
	//url := viper.GetString(`iot.url`) + "users"
	url := viper.GetString(`iot.urlusers`)
	//url := "https://iotnew.vdc2.com.vn/users"
	stringAuth := `{"email":"` + email + `", "password":"` + password + `"}`
	var jsonStr = []byte(stringAuth)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status
}
func SendSms(phone string, message string) string {
	config.ReadConfig()
	urlSmS := viper.GetString(`oneshop.sms`)
	stringSendSmS := `{"sender_id":"L33QuyBiqcsnGP57yfytqD", "phone_number":"` + phone + `","message":"` + message + `"}`
	var jsonStr = []byte(stringSendSmS)
	req, err := http.NewRequest("POST", urlSmS, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status
}
func SendEmail(sendto string, subject string, body string) string {
	config.ReadConfig()
	urlSmS := viper.GetString(`oneshop.email`)
	stringSendEmail := `{"sender":"sE9tikaouLFpsCMjPYNdpj", "to": ["` + sendto + `"],"subject":"` + subject + `","body":"` + body + `"}`
	var jsonStr = []byte(stringSendEmail)
	req, err := http.NewRequest("POST", urlSmS, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status
}

//create Random OTP
func RandomOtpString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	//length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	otp := b.String()
	//fmt.Println(otp)
	return otp
}
func RequestTokenByRT(rtoken string) string {
	config.ReadConfig()
	link := viper.GetString(`sso.url`)
	url := link + "/accessToken?grant_type=refresh_token&client_id=exampleOauthClient&client_secret=exampleOauthClientSecret&refresh_token=" + rtoken
	//url := "https://ssostandalone.vdc2.com.vn:8443/cas/oauth2.0/accessToken?grant_type=refresh_token&client_id=exampleOauthClient&client_secret=exampleOauthClientSecret&refresh_token=" + rtoken
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//Sys log
		sysLog := new(model.LogELK)
		sysLog.Datetime = time.Now()
		sysLog.Log = err.Error()
		WriteLog(sysLog)
		fmt.Println(err)
		return ""
		//Sys log
	}
	//fmt.Println(err)
	stringBody := string(body)
	//fmt.Println(stringBody)
	if stringBody != "" {
		var rtokens RTokens
		json.Unmarshal([]byte(stringBody), &rtokens)
		//fmt.Println(rtokens.Accesstoken)
		return rtokens.Accesstoken
	} else {
		return ""
	}
}

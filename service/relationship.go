package service

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"

	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

//Hàm delete user cha thì delete user con
func DeleteNode(deteleNode string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://casuser:Mellon@222.255.102.145:27017/users"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		//panic(err)
		return err
	}
	defer client.Disconnect(ctx)

	/*database := client.Database("users")
	users := database.Collection("users")*/
	collection := client.Database("users").Collection("users")

	matchStage := bson.D{{"$match", bson.D{{"$or", bson.A{bson.D{{"username", deteleNode}}, bson.D{{"hierarchy.username", deteleNode}}}}}}}
	graphStage := bson.D{{"$graphLookup", bson.D{{"from", "users"}, {"startWith", "$userparentid"}, {"connectFromField", "userparentid"}, {"connectToField", "username"}, {"as", "hierarchy"}}}}

	//showInfoCursor, err := users.Aggregate(ctx, mongo.Pipeline{graphStage, matchStage})
	showInfoCursor, err := collection.Aggregate(ctx, mongo.Pipeline{graphStage, matchStage})
	if err != nil {
		//panic(err)
		return err
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		//panic(err)
		return err
	}

	for _, doc := range showsWithInfo {
		//_, err := users.DeleteOne(ctx, bson.M{"username": doc["username"]})
		_, err := collection.DeleteOne(ctx, bson.M{"username": doc["username"]})
		if err != nil {
			//panic(err)
			return err
		}
	}
	return err
}

//Update status user cha thì update status user con
func UpdateRelationshipNode(updateNode string, status string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("users")
	matchStage := bson.D{{"$match", bson.D{{"$or", bson.A{bson.D{{"username", updateNode}}, bson.D{{"hierarchy.username", updateNode}}}}}}}
	graphStage := bson.D{{"$graphLookup", bson.D{{"from", "users"}, {"startWith", "$userparentid"}, {"connectFromField", "userparentid"}, {"connectToField", "username"}, {"as", "hierarchy"}}}}
	showInfoCursor, err := collection.Aggregate(ctx, mongo.Pipeline{graphStage, matchStage})
	if err != nil {
		return err
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		return err
	}
	for _, doc := range showsWithInfo {
		_, err := collection.UpdateOne(ctx, bson.M{"username": doc["username"]}, bson.M{"$set": bson.M{"userstatus": status}})
		if err != nil {
			return err
		}
	}
	return err
}

//Remove 1 function code trong function code list (dùng khi xoá 1 function_code or DISABLE) trong collection company_role_function
func DeleteArrayFunctionCodeCompanyRoleFunction(functionCode []string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_role_function")
	_, err = collection.UpdateMany(ctx, bson.D{{}}, bson.M{"$pull": bson.M{"functioncodelist": bson.M{"$in": functionCode}}})
	if err != nil {
		return err
	}
	return nil
}

//Update đồng thời contractstatus của collection company_role_function khi update contractstatus của collection company_product
func UpdateCompanyRoleFunctionContractStatus(productid string, contractstatus string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_role_function")
	_, err = collection.UpdateOne(ctx, bson.M{"productid": productid}, bson.M{"$set": bson.M{"contractstatus": contractstatus}})
	if err != nil {
		return err
	}
	return nil
}

//Xoá nhiều colletion company_role_function khi xoá 1 product_id của company_product
func DeleteCompanyRoleFunctionComProduct(productid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_role_function")
	_, err = collection.DeleteMany(ctx, bson.M{"productid": productid})
	if err != nil {
		return err
	}
	return nil
}

//Tìm thông tin comid và productid trong collection company_product tương ứng vs contract_code
type Response struct {
	Comid     string
	Productid string
}

func FindDeleteInfoCompanyProduct(contractcode string) (comid string, productid string, err error) {
	var response Response
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return "", "", err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_product")
	err = collection.FindOne(ctx, bson.M{"contractcode": contractcode}).Decode(&response)
	if err != nil {
		return "", "", err
	}
	comid = response.Comid
	productid = response.Productid
	return comid, productid, nil
}

//Xoá document trong collection company_role_function khi xoá document cua collection company_product
func DeleteCompanyRoleFunctionRelationshipCompanyProduct(productid string, comid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_role_function")
	_, err = collection.DeleteOne(ctx, bson.M{"productid": productid, "comid": comid})
	if err != nil {
		return err
	}
	return nil
}

//
//GetCompanyUserRole
type ArrayRole struct {
	Rolecode []string
}

func GetCompanyUserRole(username string, comid string) (err error, arrayUserRole []string) {
	//
	var arrayRole ArrayRole
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err, arrayUserRole
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_role")
	err = collection.FindOne(ctx, bson.M{"username": username, "comid": comid}).Decode(&arrayRole)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(arrayRole.Rolecode)
	return err, arrayRole.Rolecode
}

//Remove 1 role code trong array role code collection company_user_role (dùng khi xoá 1 role_code or DISABLE)
func DeleteArrayRoleCodeCompanyUserRole(roleCode []string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_role")
	_, err = collection.UpdateMany(ctx, bson.D{{}}, bson.M{"$pull": bson.M{"rolecode": bson.M{"$in": roleCode}}})
	if err != nil {
		return err
	}
	return nil
}

//Update Array Function Code List
func UpdateArrayFunctionCodeList(roleCode string, functionCodeList []string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_role_function")
	_, err = collection.UpdateOne(ctx, bson.M{"rolecode": roleCode}, bson.M{"$set": bson.M{"functioncodelist": functionCodeList}})
	if err != nil {
		return err
	}
	return nil
}

//Remove all related company comid
func RemoveAllRelatedCompanyComid(comid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	arrayCollection := []string{"company_info", "company_product", "company_function", "company_role", "company_role_function", "company_user_role", "users", "company_user_service"}
	for i := 0; i < len(arrayCollection); i++ {
		collection := client.Database("users").Collection(arrayCollection[i])
		_, err = collection.DeleteMany(ctx, bson.M{"comid": comid})
		if err != nil {
			return err
		}
	}
	return nil
}

//Update collection company_info also update collection users
func UpdateCascadeUserStatusCompanyInfo(comid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("users")
	_, err = collection.UpdateMany(ctx, bson.M{"comid": comid}, bson.M{"$set": bson.M{"userstatus": "DISABLE"}})
	if err != nil {
		return err
	}
	return nil
}

//Update collection company_info also update collection company_product, company_role_function
func UpdateCascadeContractStatusCompanyInfo(comid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	arrayCollection := []string{"company_product", "company_role_function"}
	for i := 0; i < len(arrayCollection); i++ {
		collection := client.Database("users").Collection(arrayCollection[i])
		_, err = collection.UpdateMany(ctx, bson.M{"comid": comid}, bson.M{"$set": bson.M{"contractstatus": "DISABLE"}})
		if err != nil {
			return err
		}
	}
	return nil
}

//Delete cascade company_user_role when delete users
func GetAllChildOfNodeAndDeleteCompanyUserRole(node string) error {
	//
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	//
	database := client.Database("users")
	users := database.Collection("users")
	matchStage := bson.D{{"$match", bson.D{{"username", node}}}}
	graphStage := bson.D{{"$graphLookup", bson.D{{"from", "users"}, {"startWith", "$username"}, {"connectFromField", "username"}, {"connectToField", "userparentid"}, {"as", "descendants"}}}}
	unWind := bson.D{{"$unwind", "$descendants"}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$descendants"}}}}
	proJect := bson.D{{"$project", bson.D{{"descendants", 0}}}}
	showInfoCursor, err := users.Aggregate(ctx, mongo.Pipeline{matchStage, graphStage, unWind, replaceRoot, proJect})
	if err != nil {
		panic(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		panic(err)
	}
	userole := database.Collection("company_user_role")
	for _, s := range showsWithInfo {
		//fmt.Println(s["username"])
		_, err = userole.DeleteOne(ctx, bson.M{"username": s["username"].(string)})
	}
	_, err = userole.DeleteOne(ctx, bson.M{"username": node})
	return nil
}
func GetInfomationChildOfNode(node string) (err error, usersSSO *model.UsersSSO) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err, usersSSO
	}
	defer client.Disconnect(ctx)
	//
	database := client.Database("users")
	users := database.Collection("users")
	matchStage := bson.D{{"$match", bson.D{{"username", node}}}}
	graphStage := bson.D{{"$graphLookup", bson.D{{"from", "users"}, {"startWith", "$username"}, {"connectFromField", "username"}, {"connectToField", "userparentid"}, {"as", "descendants"}}}}
	unWind := bson.D{{"$unwind", "$descendants"}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$descendants"}}}}
	proJect := bson.D{{"$project", bson.D{{"descendants", 0}}}}
	showInfoCursor, err := users.Aggregate(ctx, mongo.Pipeline{matchStage, graphStage, unWind, replaceRoot, proJect})
	if err != nil {
		return err, usersSSO
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		return err, usersSSO
	}
	//Che password
	for _, s := range showsWithInfo {
		s["password"] = ""
	}
	//convert bson.M to struct
	bs, err := json.Marshal(showsWithInfo)
	if err != nil {
		return err, usersSSO
	}
	if err := json.Unmarshal(bs, &usersSSO); err != nil {
		return err, usersSSO
	}
	//
	return nil, usersSSO
}
func CheckParentAndChild(loginusername string, checkusername string) (err error, exist bool) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err, false
	}
	defer client.Disconnect(ctx)
	//
	database := client.Database("users")
	users := database.Collection("users")
	matchStage := bson.D{{"$match", bson.D{{"username", loginusername}}}}
	graphStage := bson.D{{"$graphLookup", bson.D{{"from", "users"}, {"startWith", "$username"}, {"connectFromField", "username"}, {"connectToField", "userparentid"}, {"as", "descendants"}}}}
	unWind := bson.D{{"$unwind", "$descendants"}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$descendants"}}}}
	proJect := bson.D{{"$project", bson.D{{"descendants", 0}}}}
	showInfoCursor, err := users.Aggregate(ctx, mongo.Pipeline{matchStage, graphStage, unWind, replaceRoot, proJect})
	if err != nil {
		return err, false
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		return err, false
	}
	var arrayChild []string
	//fmt.Println(showsWithInfo)
	for i := 0; i < len(showsWithInfo); i++ {
		arrayChild = append(arrayChild, showsWithInfo[i]["username"].(string))
	}
	//fmt.Println(arrayChild)
	exist = contains(arrayChild, checkusername)
	return err, exist
}

//Check if element contain in array
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//Response RoleCode of Thing. Dang Array
type ArrayFunction struct {
	Functioncodelist []string
}

func GetCompanyUserRoleByRoleCode(rolecode string) (err error, responseArr []string) {
	//
	var arrayFunction ArrayFunction
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err, responseArr
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_role_function")
	err = collection.FindOne(ctx, bson.M{"rolecode": rolecode}).Decode(&arrayFunction)
	if err != nil {
		fmt.Println(err)
	}
	return err, arrayFunction.Functioncodelist
}

//Password Management
//Update Password with Username
func UpdatePasswordByUsername(username string, newpassword string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("users")
	_, err = collection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"password": newpassword}})
	if err != nil {
		return err
	}
	return nil
}

//Check Product Integrated Exist or not?
func CheckProductIntegratedElementExist(fileter_key string, fileter_value string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("product_integrated")
	count, err = collection.CountDocuments(ctx, bson.M{fileter_key: fileter_value})
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}

//Check trùng comid + productid
func CheckTwiceComIDandProductIDinCompanyProduct(comid string, productid string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_product")
	count, err = collection.CountDocuments(ctx, bson.M{"comid": comid, "productid": productid})
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}

//Remove all related product id
func RemoveAllRelatedProductId(productid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	arrayCollection := []string{"product_integrated", "company_product", "company_role_function"}
	for i := 0; i < len(arrayCollection); i++ {
		collection := client.Database("users").Collection(arrayCollection[i])
		_, err = collection.DeleteMany(ctx, bson.M{"productid": productid})
		if err != nil {
			return err
		}
	}
	return nil
}

//Check email exist
func CheckUserEmailExist(email string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("users")
	count, err = collection.CountDocuments(ctx, bson.M{"useremail": email})
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}
func CheckUserPhoneExist(usertel string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("users")
	count, err = collection.CountDocuments(ctx, bson.M{"usertel": usertel})
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}

//Check ComId Exist or Not?
func CheckComIdExist(fileter_key string, fileter_value string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_info")
	count, err = collection.CountDocuments(ctx, bson.M{fileter_key: fileter_value})
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}

//Clone and Update Service
func CloneAndUpdateServiceOauth(name string, serviceMetaData *model.ServiceMetaData) error {
	var m bson.M
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	templateOauth := viper.GetString(`servicetemplate.oauth2`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	database := client.Database("cas")
	service := database.Collection("casserviceregistry")
	err = service.FindOne(ctx, bson.M{"clientId": templateOauth}).Decode(&m)
	if err != nil {
		return err
	}
	//generate random
	rand.Seed(int64(time.Now().Nanosecond()))
	value := rand.Int()
	//
	//m["_id"] = rand.Int31() + rand.Int31()
	m["_id"] = value
	m["clientSecret"] = serviceMetaData.Clientsecret
	m["clientId"] = serviceMetaData.Clientid
	m["serviceId"] = serviceMetaData.Serviceid
	m["name"] = name
	_, err = service.InsertOne(context.TODO(), m)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}

//check service name exist
func CheckUserServiceNameExist(servicename string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_service")
	count, err = collection.CountDocuments(ctx, bson.M{"servicename": servicename})
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}

//check clienrID exist
func CheckUserServiceOauthClientIdExist(clientid string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_service")
	count, err = collection.CountDocuments(ctx, bson.M{"servicemetadata.clientid": clientid})
	//count, err = collection.FindOne(ctx, bson.M{"servicemetadata.clientid": clientid}).CountDocuments()
	if err != nil {
		return -1, err
	}
	//fmt.Println(count)
	return count, nil
}
func CheckUserNameExist(username string) (count int64, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return -1, err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("users")
	count, err = collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return -1, err
	}
	return count, nil
}
func AddUserService(companyUserService *model.CompanyUserService) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_service")
	_, err = collection.InsertOne(ctx, companyUserService)
	if err != nil {
		return err
	}
	return nil
}
func UpdateOauthUserService(companyUserService *model.CompanyUserService) error {
	var encodeStringSecret string
	serviceMetaData := companyUserService.Servicemetadata
	serviceName := companyUserService.Servicename
	status := companyUserService.Status
	clientSecret := serviceMetaData.Clientsecret
	serviceId := serviceMetaData.Serviceid

	//encode base64
	encodeStringSecret = b64.StdEncoding.EncodeToString([]byte(clientSecret))
	//Connect database
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	service := client.Database("cas").Collection("casserviceregistry")
	collection := client.Database("users").Collection("company_user_service")
	//_, err = collection.InsertOne(ctx, companyUserService)
	//update user service
	updateUserService := bson.M{"$set": bson.M{"status": status, "servicemetadata.serviceid": serviceId, "servicemetadata.clientsecret": clientSecret}}
	_, err1 := collection.UpdateOne(ctx, bson.M{"servicename": serviceName}, updateUserService)
	if err1 != nil {
		return err1
	}
	//DISABLE thì chuyển thành base64
	if status == "DISABLE" {
		updateServiceString := bson.M{"$set": bson.M{"clientSecret": encodeStringSecret, "serviceId": serviceId}}
		_, err2 := service.UpdateOne(ctx, bson.M{"name": serviceName}, updateServiceString)
		if err2 != nil {
			return err2
		}
	} else if status == "ACTIVE" {
		updateServiceString := bson.M{"$set": bson.M{"clientSecret": clientSecret, "serviceId": serviceId}}
		_, err3 := service.UpdateOne(ctx, bson.M{"name": serviceName}, updateServiceString)
		if err3 != nil {
			return err3
		}
	}
	return nil
}

/*func SearchAllUserService() error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_service")
	var companyUserServices CompanyUserServices
	_, err = collection.Find(ctx, bson.M{}).All(&companyUserServices)
	if err != nil {
		return err
	}
	return nil
}*/
func RemoveServiceCas(filter, condition string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	database := client.Database("cas")
	service := database.Collection("casserviceregistry")
	_, err = service.DeleteOne(ctx, bson.M{filter: condition})
	if err != nil {
		return err
	}
	return nil
}
func RemoveUserService(filter, condition string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	database := client.Database("users")
	service := database.Collection("company_user_service")
	_, err = service.DeleteOne(ctx, bson.M{filter: condition})
	if err != nil {
		return err
	}
	return nil
}
func GetAllServiceNameByUserName(username string) (serviceName []string, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return serviceName, err
	}
	defer client.Disconnect(ctx)
	database := client.Database("users")
	service := database.Collection("company_user_service")
	filterCursor, err := service.Find(ctx, bson.M{"username": username})
	if err != nil {
		return serviceName, err
	}
	var showsWithInfo []bson.M
	if err = filterCursor.All(ctx, &showsWithInfo); err != nil {
		return serviceName, err
	}
	for i := 0; i < len(showsWithInfo); i++ {
		a := showsWithInfo[i]["servicename"]
		serviceName = append(serviceName, a.(string))
	}
	//fmt.Println(serviceName)
	return serviceName, nil
}
func RemoveServiceCasByUserName(username string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	database := client.Database("cas")
	service := database.Collection("casserviceregistry")
	arrayServiceName, err := GetAllServiceNameByUserName(username)
	if err != nil {
		return err
	}
	//fmt.Println(arrayServiceName)
	for i := 0; i < len(arrayServiceName); i++ {
		_, err = service.DeleteOne(ctx, bson.M{"name": arrayServiceName[i]})
		//fmt.Println("Delete " + arrayServiceName[i])
	}
	if err != nil {
		return err
	}
	return nil
}
func DeleteAllUserServiceByUsername(username string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_service")
	_, err = collection.DeleteMany(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	return nil
}

//update user role
func UpdateCompanyUserRoleRelationship(username, datecreate string, roleCode []string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	database := client.Database("users")
	service := database.Collection("company_user_role")
	update := bson.M{"$set": bson.M{"datecreate": datecreate, "rolecode": roleCode}}
	_, err = service.UpdateOne(ctx, bson.M{"username": username}, update)
	if err != nil {
		return err
	}
	return nil
}

//Update status company_info also update collection company_user_service
func UpdateCascadeUserServiceStatusCompanyInfo(comid string) error {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("users").Collection("company_user_service")
	_, err = collection.UpdateMany(ctx, bson.M{"comid": comid}, bson.M{"$set": bson.M{"status": "DISABLE"}})
	if err != nil {
		return err
	}
	return nil
}
func CheckServiceStatus(clientid string) (status string, err error) {
	config.ReadConfig()
	link := viper.GetString(`mongo.link`)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(link))
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)
	database := client.Database("users")
	service := database.Collection("company_user_service")
	var result bson.M
	err = service.FindOne(ctx, bson.M{"servicemetadata.clientid": clientid}).Decode(&result)
	if err != nil {
		return "", nil
	}
	//fmt.Println(result)
	status = (result["status"]).(string)
	return status, nil
}

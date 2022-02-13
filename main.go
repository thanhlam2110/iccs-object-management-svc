package main

import (
	"net/http"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/integratedkafka"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/service"
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/transport"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	go func() {
		//init producer
		config.ReadConfig()
		url := viper.GetString(`kafka.url`)
		integratedkafka.KafkaInit(url)
		producer, _ := transport.InitProducer(url)
		transport.PublishMessage(producer, "SSO_CREATE_TOPIC", kafka.PartitionAny, "INIT SSO_CREATE_TOPIC")
		transport.PublishMessage(producer, "SSO_UPDATE_TOPIC", kafka.PartitionAny, "INIT SSO_UPDATE_TOPIC")
		transport.PublishMessage(producer, "SSO_DELETE_TOPIC", kafka.PartitionAny, "INIT SSO_DELETE_TOPIC")
		//int consumer
		integratedkafka.HandlerConsumer()
	}()
}
func main() {
	e := echo.New()

	e.Validator = &service.CustomValidator{Validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "SSO API")
	})

	//USER MANAGEMENT
	e.POST("/api/sso/v1/user/register", service.CreateUser)
	e.PATCH("/api/sso/v1/user/updatenduser", service.UpdateEndUser)
	e.PATCH("/api/sso/v1/user/updatecompanyuser", service.UpdateCompanyUser)
	e.DELETE("/api/sso/v1/user/deletenduser/:username", service.DeleteEndUser)
	e.GET("/api/sso/v1/user/info/:username", service.GetUserInfo)
	e.GET("/api/sso/v1/user/childinfo/:username", service.GetUserChildInfo)
	e.GET("/api/sso/v1/user/infobycomid/:comid", service.GetAllUserByComId)
	e.DELETE("/api/sso/v1/user/deleteUserRelationship/:username", service.DeleteRelationshipUser)
	e.PATCH("/api/sso/v1/user/checkparenthavechild", service.CheckParentHaveChild)
	e.POST("/api/sso/v1/user/paging", service.PagingUserInfo)
	e.POST("/api/sso/v1/user/pagingandfilter", service.PagingAndFilterUserInfo)
	//SSO
	e.POST("/api/sso/v1/user/auth", service.Auth)
	e.POST("/api/sso/v1/user/requestToken", service.RequestSSOToken)
	e.POST("/api/sso/v2/user/requestToken", service.RequestSSOTokenv2)
	e.POST("/api/sso/v1/user/authenToken", service.AuthenSSOToken)
	e.POST("/api/sso/v1/user/ioToken", service.RequestSSOIoToken)
	e.POST("/api/sso/v1/user/ioTAuthenKey", service.GetIoTAuthKeyFromRedis)
	e.POST("/api/sso/v1/user/getTokenByRFToken", service.RefreshSSOToken)

	//service.GetUserRoleWithProductID("KHOA_DIEN", "DHBK", "IOT_Platform")
	//fmt.Println("rolecode:" + rolecode)

	//PRODUCT_INTEGRATED
	e.POST("/api/sso/v1/productintegrated/create", service.CreateProductIntegrated)
	e.GET("/api/sso/v1/user/productintegrated/:productid", service.GetProductIntegratedInfo)
	e.PATCH("/api/sso/v1/user/productintegrated/update", service.UpdateProductIntegratedByProductId)
	e.DELETE("/api/sso/v1/user/productintegrated/:productid", service.DeleteProductIntegratedByProductId)
	e.GET("/api/sso/v1/user/productintegrated/all", service.GetAllProductIntegrated)
	e.POST("/api/sso/v1/productintegrated/paging", service.PagingProductIntegrated)
	e.POST("/api/sso/v1/productintegrated/pagingandfilter", service.PagingAndFilterProductIntegrated)
	//COMPANY PRODUCT
	e.POST("/api/sso/v1/companyproduct/create", service.CreateCompanyProduct)
	e.GET("/api/sso/v1/companyproductbycomid/:comid", service.GetCompanyProductInfoByComID)
	e.GET("/api/sso/v1/companyproductbycontractcode/:contractcode", service.GetCompanyProductInfoByContractCode)
	//e.DELETE("/api/sso/v1/user/companyproductbycomid/:comid", service.DeleteCompanyProductByComId)
	e.DELETE("/api/sso/v1/user/companyproductbycontractcode/:contractcode", service.DeleteCompanyProductByContractCode)
	e.PATCH("/api/sso/v1/user/companyproduct/update", service.UpdateCompanyProductByContractCode)
	e.POST("/api/sso/v1/companyproduct/paging", service.PagingCompanyProduct)
	e.POST("/api/sso/v1/companyproduct/pagingandfilter", service.PagingAndFilterCompanyProduct)
	//e.PATCH("/api/sso/v1/user/companyproduct/updatestatus", service.UpdateCompanyProductStatusByComId)
	//COMPANY INFO
	e.POST("/api/sso/v1/companyinfo/create", service.CreateCompanyInfo)
	e.GET("/api/sso/v1/companyinfo/:comid", service.GetCompanyInfoByComID)
	e.DELETE("/api/sso/v1/user/companyinfo/:comid", service.DeleteCompanyInfo)
	e.PATCH("/api/sso/v1/user/companyinfo/update", service.UpdateCompanyInfo)
	e.GET("/api/sso/v1/companyinfo/all", service.GetAllCompanyInfo)
	e.POST("/api/sso/v1/companyinfo/paging", service.PagingCompanyInfo)
	e.POST("/api/sso/v1/companyinfo/pagingandfilter", service.PagingAndFilterCompanyInfo)
	//COMPANY FUNCTION
	e.GET("/api/sso/v1/companyfunctionbycode/:functioncode", service.GetCompanyFunctionByCode)
	e.GET("/api/sso/v1/companyfunctionbycomid/:comid", service.GetCompanyFunctionInfoByComID)
	e.POST("/api/sso/v1/companyfunction/create", service.CreateCompanyFunction)
	e.DELETE("/api/sso/v1/user/companyfunctionbycode/:functioncode", service.DeleteCompanyFunctionCode)
	//e.DELETE("/api/sso/v1/user/companyfunctionbycomid/:comid", service.DeleteCompanyFunctionComId)
	e.PATCH("/api/sso/v1/user/companyfunctionbycode/update", service.UpdateCompanyFunctionByCode)
	e.POST("/api/sso/v1/companyfunction/paging", service.PagingCompanyFunction)
	e.POST("/api/sso/v1/companyfunction/pagingandfilter", service.PagingCompanyFunctionAndFilter)
	//COMPANY ROLE
	e.GET("/api/sso/v1/companyrolebycode/:rolecode", service.GetCompanyRoleByCode)
	e.GET("/api/sso/v1/companyrolebycomid/:comid", service.GetCompanyRoleByComID)
	e.POST("/api/sso/v1/companyrole/create", service.CreateCompanyRole)
	e.PATCH("/api/sso/v1/user/companyrolebycode/update", service.UpdateCompanyRoleByRoleCode)
	e.DELETE("/api/sso/v1/user/companyrolebycode/:rolecode", service.DeleteCompanyRoleCode)
	e.POST("/api/sso/v1/companyrole/paging", service.PagingCompanyRole)
	e.POST("/api/sso/v1/companyrole/pagingandfilter", service.PagingCompanyRoleAndFilter)
	//COMPANY ROLE FUNCTION
	e.GET("/api/sso/v1/companyrolefunctionbycode/:rolecode", service.GetCompanyRoleFunctionByRoleCode)
	e.GET("/api/sso/v1/companyrolefunctionbycomid/:comid", service.GetCompanyRoleFunctionByComId)
	e.POST("/api/sso/v1/companyrolefunction/create", service.CreateCompanyRoleFunction)
	e.PATCH("/api/sso/v1/user/companyrolefunction/update", service.UpdateFunctionCodeCompanyRoleFunction)
	e.DELETE("/api/sso/v1/user/companyrolefunctionrolecode/:rolecode", service.DeleteCompanyRoleFunctionByRoleCode)
	e.POST("/api/sso/v1/companyrolefunction/paging", service.PagingCompanyRoleFunction)
	e.POST("/api/sso/v1/companyrolefunction/pagingandfilter", service.PagingCompanyRoleFunctionAndFilter)
	//COMPANY USER ROLE
	e.GET("/api/sso/v1/companyuserrolebyusername/:username", service.GetCompanyUserRoleByUsername)
	e.GET("/api/sso/v1/companyuserrolebycomid/:comid", service.GetCompanyUserRoleByComId)
	e.POST("/api/sso/v1/companyuserrole/create", service.CreateCompanyUserRole)
	e.DELETE("/api/sso/v1/user/companyuserole/:username", service.DeleteCompanyUserRoleByUsername)
	e.PATCH("/api/sso/v1/user/companyuserrole/update", service.UpdateCompanyUserRole)
	e.POST("/api/sso/v1/companyuserole/paging", service.PagingCompanyUseRole)
	e.POST("/api/sso/v1/companyuserole/pagingandfilter", service.PagingCompanyUseRoleAndFilter)
	//PASSWORD MANAGEMENT
	e.POST("/api/sso/v1/passwordmanagement/forgotusername", service.RequestForgotUserName)
	e.POST("/api/sso/v1/passwordmanagement/requestresetpassword", service.RequestResetPassword)
	e.POST("/api/sso/v1/passwordmanagement/resetpassword", service.ResetPassword)
	////////////////////////////////KAFKA INTEGRATED/////////////////////
	//-------------------Company Info-------------------//
	e.POST("/api/sso/v1/companyinfokafka/create", integratedkafka.KafkaCreateCompanyInfo)
	e.POST("/api/sso/v1/companyinfokafkacreate/status", integratedkafka.GetKafkaCreateCompanyInfoStatus)
	e.PATCH("/api/sso/v1/companyinfokafka/update", integratedkafka.KafkaUpdateCompanyInfo)
	e.POST("/api/sso/v1/companyinfokafkaupdate/status", integratedkafka.GetKafkaUpdateCompanyInfoStatus)
	e.DELETE("/api/sso/v1/companyinfokafkadelete/:comid", integratedkafka.KafkaDeleteCompanyInfo)
	e.POST("/api/sso/v1/companyinfokafkadelete/status/:comid", integratedkafka.GetKafkaDeleteCompanyInfoStatus)
	//-------------------Company Info-------------------//
	//-------------------Product Integrated-------------------//
	e.POST("/api/sso/v1/productintegratedkafka/create", integratedkafka.KafkaCreateProductIntegrated)
	e.POST("/api/sso/v1/productintegratedkafkacreate/status", integratedkafka.GetKafkaCreateProductIntegratedStatus)
	e.PATCH("/api/sso/v1/productintegratedkafka/update", integratedkafka.KafkaUpdateProductIntegrated)
	e.POST("/api/sso/v1/productintegratedkafkaupdate/status", integratedkafka.GetKafkaUpdateProductIntegratedStatus)
	e.DELETE("/api/sso/v1/productintegratedkafkadelete/:productid", integratedkafka.KafkaDeleteProductIntegrated)
	e.POST("/api/sso/v1/productintegratedkafkadelete/status/:productid", integratedkafka.GetKafkaDeleteProductIntegratedStatus)
	//-------------------Product Integrated-------------------//
	//-------------------Company Product-------------------//
	e.POST("/api/sso/v1/companyproductkafka/create", integratedkafka.KafkaCreateCompanyProduct)
	e.POST("/api/sso/v1/companyproductkafkacreate/status", integratedkafka.GetKafkaCreateCompanyProductStatus)
	e.PATCH("/api/sso/v1/companyproductkafka/updatebycontractcode", integratedkafka.KafkaUpdateCompanyProductByContractCode)
	e.POST("/api/sso/v1/companyproductkafkaupdatebycontractcode/status", integratedkafka.GetKafkaUpdateCompanyProductStatus)
	e.DELETE("/api/sso/v1/companyproductkafkadeletebycontractcode/:contractcode", integratedkafka.KafkaDeleteCompanyProduct)
	e.POST("/api/sso/v1/companyproductkafkadeletebycontractcode/status/:contractcode", integratedkafka.GetKafkaDeleteCompanyProductStatus)
	//-------------------Company Product-------------------//
	//-------------------Company Function-------------------//
	e.POST("/api/sso/v1/companyfunctionkafka/create", integratedkafka.KafkaCreateCompanyFunction)
	e.POST("/api/sso/v1/companyfunctionkafkacreate/status", integratedkafka.GetKafkaCreateCompanyFunctionStatus)
	e.PATCH("/api/sso/v1/companyfunctionkafka/update", integratedkafka.KafkaUpdateCompanyFunction)
	e.POST("/api/sso/v1/companyfunctionkafkaupdate/status", integratedkafka.GetKafkaUpdateCompanyFunctionStatus)
	e.DELETE("/api/sso/v1/companyfunctionkafkadelete/:functioncode", integratedkafka.KafkaDeleteCompanyFunction)
	e.POST("/api/sso/v1/companyfunctionkafkadelete/status/:functioncode", integratedkafka.GetKafkaDeleteCompanyFunctionStatus)
	//-------------------Company Function-------------------//
	//-------------------Company Role-------------------//
	e.POST("/api/sso/v1/companyrolekafka/create", integratedkafka.KafkaCreateCompanyRole)
	e.POST("/api/sso/v1/companyrolekafkacreate/status", integratedkafka.GetKafkaCreateCompanyRoleStatus)
	e.PATCH("/api/sso/v1/companyrolekafka/update", integratedkafka.KafkaUpdateCompanyRole)
	e.POST("/api/sso/v1/companyrolekafkaupdate/status", integratedkafka.GetKafkaUpdateCompanyRoleStatus)
	e.DELETE("/api/sso/v1/companyrolekafkadelete/:rolecode", integratedkafka.KafkaDeleteCompanyRole)
	e.POST("/api/sso/v1/companyrolekafkadelete/status/:rolecode", integratedkafka.GetKafkaDeleteCompanyRoleStatus)
	//-------------------Company Role-------------------//
	//-------------------Company Role Function-------------------//
	e.POST("/api/sso/v1/companyrolefunctionkafka/create", integratedkafka.KafkaCreateCompanyRoleFunction)
	e.PATCH("/api/sso/v1/companyrolefunctionkafka/update", integratedkafka.KafkaUpdateCompanyRoleFunction)
	e.POST("/api/sso/v1/companyrolefunctionkafkacreate/status", integratedkafka.GetKafkaCreateCompanyRoleFunctionStatus)
	e.POST("/api/sso/v1/companyrolefunctionkafkaupdate/status", integratedkafka.GetKafkaUpdateCompanyRoleFunctionStatus)
	//-------------------Company Role Function-------------------//
	//-------------------Company User Role-------------------//
	e.POST("/api/sso/v1/companyuserolekafka/create", integratedkafka.KafkaCreateCompanyUseRole)
	e.PATCH("/api/sso/v1/companyuserolekafka/update", integratedkafka.KafkaUpdateCompanyUseRole)
	e.POST("/api/sso/v1/companyuserolekafkacreate/status", integratedkafka.GetKafkaCreateCompanyUseRoleStatus)
	e.POST("/api/sso/v1/companyuserolekafkaupdate/status", integratedkafka.GetKafkaUpdateCompanyUseRoleStatus)
	//-------------------Company User Role-------------------//
	//-------------------User-------------------//
	e.POST("/api/sso/v1/userkafka/create", integratedkafka.KafkaCreateUser)
	e.POST("/api/sso/v1/userkafkacreate/status", integratedkafka.GetKafkaCreateUserStatus)
	e.PATCH("/api/sso/v1/userkafka/update", integratedkafka.KafkaUpdateUser)
	//e.PATCH("/api/sso/v1/userkafka/updatenduser", integratedkafka.KafkaUpdateEndUser)
	e.POST("/api/sso/v1/userkafkaupdate/status", integratedkafka.GetKafkaUpdateUserStatus)
	//-------------------User-------------------//
	////////////////////////////////KAFKA INTEGRATED/////////////////////
	//-------------------User Service-------------------//
	e.POST("/api/sso/v1/userservice/create", service.CreateUserServiceOauth)
	e.PATCH("/api/sso/v1/userservice/update", service.UpdateUserServiceOauth)
	e.GET("/api/sso/v1/userservicegetallbyusername/:username", service.GetAllUserServiceByUsername)
	e.GET("/api/sso/v1/userservicegetall/", service.GetAllUserService)
	e.GET("/api/sso/v1/userservicegetallbycomid/:comid", service.GetAllUserServiceByComid)
	e.POST("/api/sso/v1/userservice/paging", service.PagingCompanyUserService)
	e.POST("/api/sso/v1/userservice/pagingandfilter", service.PagingCompanyUserServiceAndFilter)
	e.DELETE("/api/sso/v1/userservice/:servicename", service.DeleteUserServiceOauthByServiceName)
	//-------------------User Service-------------------//
	//TEST
	//var functionCode []string
	//functionCode = append(functionCode, "DHBK1_FUNC_1")
	//service.DeleteArrayFunctionCodeCompanyRoleFunction(functionCode)
	//
	//service.GetUserRoleWithProductID("KHOA_DIEN", "DHBK", "IOT_Platform")
	//fmt.Println("rolecode:" + rolecode)
	//
	e.Logger.Fatal(e.Start(":1323"))

}

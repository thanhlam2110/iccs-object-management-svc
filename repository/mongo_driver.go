package repository

import (
	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
)

// ProfileRepositoryMongo - ProfileRepositoryMongo
type ProfileRepositoryMongo struct {
	db         *mgo.Database
	collection string
}

//NewProfileRepositoryMongo - NewProfileRepositoryMongo
func NewProfileRepositoryMongo(db *mgo.Database, collection string) *ProfileRepositoryMongo {
	return &ProfileRepositoryMongo{
		db:         db,
		collection: collection,
	}
}

//<------------------CLOSE CONNECTON------------------>

//<------------------CLOSE CONNECTON------------------>
//<------------------USER------------------>
//Save - Save
func (r *ProfileRepositoryMongo) Save(userSSO *model.UserSSO) error {
	err := r.db.C(r.collection).Insert(userSSO)
	return err
}

//Update - Update
func (r *ProfileRepositoryMongo) Update(username string, userSSO *model.UserSSO) error {
	err := r.db.C(r.collection).Update(bson.M{"username": username}, userSSO)
	return err
}

//Delete - Delete End User
func (r *ProfileRepositoryMongo) Delete(username string) error {
	err := r.db.C(r.collection).Remove(bson.M{"username": username})
	return err
}

//FindByUser - FindByUser
func (r *ProfileRepositoryMongo) FindByUser(username string) (*model.UserSSO, error) {
	var userSSO model.UserSSO
	err := r.db.C(r.collection).Find(bson.M{"username": username}).One(&userSSO)
	if err != nil {
		return nil, err
	}
	return &userSSO, nil
}

//FindAll - FindAll
func (r *ProfileRepositoryMongo) FindAll(comid string) (model.UsersSSO, error) {
	var usersSSO model.UsersSSO

	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&usersSSO)
	if err != nil {
		return nil, err
	}
	return usersSSO, nil
}

//Find Rolecode in Company_user_role collection
/*func (r *ProfileRepositoryMongo) FindRoleByUser(username string) (*model.CompanyUserRole, error) {
	var companyUserRole model.CompanyUserRole
	err := r.db.C(r.collection).Find(bson.M{"username": username}).One(&companyUserRole)
	if err != nil {
		return nil, err
	}
	return &companyUserRole, nil
}*/
//<------------------PRODUCT INTEGRATED------------------>
//Create product_integrated
func (r *ProfileRepositoryMongo) SaveProductIntegrated(productIntegrated *model.ProductIntegrated) error {
	err := r.db.C(r.collection).Insert(productIntegrated)
	return err
}

//Find  product_integrated by product_id
func (r *ProfileRepositoryMongo) FindProductIntegrated(productid string) (*model.ProductIntegrated, error) {
	var productIntegrated model.ProductIntegrated
	err := r.db.C(r.collection).Find(bson.M{"productid": productid}).One(&productIntegrated)
	if err != nil {
		return nil, err
	}
	return &productIntegrated, nil
}

//Update product_integrated by product_id
func (r *ProfileRepositoryMongo) UpdateProductIntegrated(productid string, productIntegrated *model.ProductIntegrated) error {
	err := r.db.C(r.collection).Update(bson.M{"productid": productid}, productIntegrated)
	return err
}

//Delete product_integrated by product_id
func (r *ProfileRepositoryMongo) DeleteProductIntegrated(productid string) error {
	err := r.db.C(r.collection).Remove(bson.M{"productid": productid})
	return err
}

//Find all product_integrated
func (r *ProfileRepositoryMongo) FindProductIntegratedAll() (*model.ProductIntegrateds, error) {
	var productIntegrateds model.ProductIntegrateds
	err := r.db.C(r.collection).Find(nil).All(&productIntegrateds)
	if err != nil {
		return nil, err
	}
	return &productIntegrateds, nil
}

//<------------------COMPANY PRODUCT------------------>
//Create company_product
func (r *ProfileRepositoryMongo) SaveCompanyProduct(companyProduct *model.CompanyProduct) error {
	err := r.db.C(r.collection).Insert(companyProduct)
	return err
}

//Find  company_product by com_id
func (r *ProfileRepositoryMongo) FindCompanyProduct(comid string) (*model.CompanyProducts, error) {
	var companyProducts model.CompanyProducts
	//err := r.db.C(r.collection).Find(bson.M{"comid": comid}).One(&companyProduct)
	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&companyProducts)
	if err != nil {
		return nil, err
	}
	return &companyProducts, nil
}

//Find  company_product by contract code
func (r *ProfileRepositoryMongo) FindCompanyProductContract(contractcode string) (*model.CompanyProduct, error) {
	var companyProduct model.CompanyProduct
	err := r.db.C(r.collection).Find(bson.M{"contractcode": contractcode}).One(&companyProduct)
	if err != nil {
		return nil, err
	}
	return &companyProduct, nil
}

//Delete company_product by com_id
func (r *ProfileRepositoryMongo) DeleteCompanyProduct(comid string) error {
	//err := r.db.C(r.collection).Remove(bson.M{"comid": comid})
	_, err := r.db.C(r.collection).RemoveAll(bson.M{"comid": comid})
	return err
}

//Delete company_product by contract code
func (r *ProfileRepositoryMongo) DeleteCompanyProductContractCode(contractcode string) error {
	//err := r.db.C(r.collection).Remove(bson.M{"comid": comid})
	err := r.db.C(r.collection).Remove(bson.M{"contractcode": contractcode})
	return err
}

//Update company product by contract code
func (r *ProfileRepositoryMongo) UpdateCompanyProduct(contractcode string, companyProduct *model.CompanyProduct) error {
	err := r.db.C(r.collection).Update(bson.M{"contractcode": contractcode}, companyProduct)
	return err
}

//Update company product status by comid
func (r *ProfileRepositoryMongo) UpdateCompanyProductStatus(comid string, contractstatus string) error {
	_, err := r.db.C(r.collection).UpdateAll(bson.M{"comid": comid}, bson.M{"$set": bson.M{"contractstatus": contractstatus}})
	return err
}

//<------------------COMPANY INFO------------------>
//Create company info
func (r *ProfileRepositoryMongo) SaveCompanyInfo(companyInfo *model.CompanyInfo) error {
	err := r.db.C(r.collection).Insert(companyInfo)
	return err
}

//Find company info by com_id
func (r *ProfileRepositoryMongo) FindCompanyInfo(comid string) (*model.CompanyInfo, error) {
	var companyInfo model.CompanyInfo
	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).One(&companyInfo)
	if err != nil {
		return nil, err
	}
	return &companyInfo, nil
}

//Find all company info
func (r *ProfileRepositoryMongo) FindCompanyInfoAll() (*model.CompanyInfos, error) {
	var companyInfos model.CompanyInfos
	err := r.db.C(r.collection).Find(nil).All(&companyInfos)
	if err != nil {
		return nil, err
	}
	return &companyInfos, nil
}

//Update company info by comid
func (r *ProfileRepositoryMongo) UpdateCompanyInfoByComId(comid string, companyInfo *model.CompanyInfo) error {
	err := r.db.C(r.collection).Update(bson.M{"comid": comid}, companyInfo)
	return err
}

//<------------------COMPANY FUNCTION------------------>
//Create company function
func (r *ProfileRepositoryMongo) SaveCompanyFunction(companyFunction *model.CompanyFunction) error {
	err := r.db.C(r.collection).Insert(companyFunction)
	return err
}

//Find company function by function code
func (r *ProfileRepositoryMongo) FindCompanyFunctionByCode(functioncode string) (*model.CompanyFunction, error) {
	var companyFunction model.CompanyFunction
	err := r.db.C(r.collection).Find(bson.M{"functioncode": functioncode}).One(&companyFunction)
	if err != nil {
		return nil, err
	}
	return &companyFunction, nil
}

//Find company_function by com_id
func (r *ProfileRepositoryMongo) FindCompanyFunction(comid string) (*model.CompanyFunctions, error) {
	var companyFunctions model.CompanyFunctions
	//err := r.db.C(r.collection).Find(bson.M{"comid": comid}).One(&companyProduct)
	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&companyFunctions)
	if err != nil {
		return nil, err
	}
	return &companyFunctions, nil
}

//Delete company_function by code
func (r *ProfileRepositoryMongo) DeleteCompanyFunctionByFunctionCode(functioncode string) error {
	err := r.db.C(r.collection).Remove(bson.M{"functioncode": functioncode})
	return err
}

//Delete company_function by com_id
/*func (r *ProfileRepositoryMongo) DeleteCompanyFunction(comid string) error {
	_, err := r.db.C(r.collection).RemoveAll(bson.M{"comid": comid})
	return err
}*/

//Update company function by function code
func (r *ProfileRepositoryMongo) UpdateCompanyFunction(functioncode string, companyFunction *model.CompanyFunction) error {
	err := r.db.C(r.collection).Update(bson.M{"functioncode": functioncode}, companyFunction)
	return err
}

//<------------------COMPANY ROLE------------------>
func (r *ProfileRepositoryMongo) FindCompanyRoleByCode(rolecode string) (*model.CompanyRole, error) {
	var companyRole model.CompanyRole
	err := r.db.C(r.collection).Find(bson.M{"rolecode": rolecode}).One(&companyRole)
	if err != nil {
		return nil, err
	}
	return &companyRole, nil
}

//Find  company_role by com_id
func (r *ProfileRepositoryMongo) FindCompanyRole(comid string) (*model.CompanyRoles, error) {
	var companyRoles model.CompanyRoles
	//err := r.db.C(r.collection).Find(bson.M{"comid": comid}).One(&companyProduct)
	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&companyRoles)
	if err != nil {
		return nil, err
	}
	return &companyRoles, nil
}

//Create company role
func (r *ProfileRepositoryMongo) SaveCompanyRole(companyRole *model.CompanyRole) error {
	err := r.db.C(r.collection).Insert(companyRole)
	return err
}

//DELETE BY ROLE CODE
func (r *ProfileRepositoryMongo) DeleteCompanyRole(rolecode string) error {
	err := r.db.C(r.collection).Remove(bson.M{"rolecode": rolecode})
	return err
}

//Update company role by role code
func (r *ProfileRepositoryMongo) UpdateCompanyRoleCode(rolecode string, companyRole *model.CompanyRole) error {
	err := r.db.C(r.collection).Update(bson.M{"rolecode": rolecode}, companyRole)
	return err
}

//<------------------COMPANY ROLE FUNCTION------------------>
//FIND BY ROLE CODE
func (r *ProfileRepositoryMongo) FindCompanyRoleFunctionRoleCode(rolecode string) (*model.CompanyRoleFunction, error) {
	var companyRoleFunction model.CompanyRoleFunction
	err := r.db.C(r.collection).Find(bson.M{"rolecode": rolecode}).One(&companyRoleFunction)
	if err != nil {
		return nil, err
	}
	return &companyRoleFunction, nil
}

//FIND ALL BY COMID
func (r *ProfileRepositoryMongo) FindCompanyRoleFunctionComId(comid string) (*model.CompanyRoleFunctions, error) {
	var companyRoleFunctions model.CompanyRoleFunctions
	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&companyRoleFunctions)
	if err != nil {
		return nil, err
	}
	return &companyRoleFunctions, nil
}

//DELETE BY ROLE CODE
func (r *ProfileRepositoryMongo) DeleteCompanyRoleFunction(rolecode string) error {
	err := r.db.C(r.collection).Remove(bson.M{"rolecode": rolecode})
	return err
}

//CREATE COMPANY_ROLE_FUCNTION
func (r *ProfileRepositoryMongo) SaveCompanyRoleFunction(companyRoleFunction *model.CompanyRoleFunction) error {
	err := r.db.C(r.collection).Insert(companyRoleFunction)
	return err
}

//<------------------COMPANY USER ROLE------------------>
//FIND ONE BY Username
func (r *ProfileRepositoryMongo) FindCompanyUserRoleUsername(username string) (*model.CompanyUserRole, error) {
	var companyUserRole model.CompanyUserRole
	err := r.db.C(r.collection).Find(bson.M{"username": username}).One(&companyUserRole)
	if err != nil {
		return nil, err
	}
	return &companyUserRole, nil
}

//FIND ALL
func (r *ProfileRepositoryMongo) FindCompanyUserRoleComId(comid string) (*model.CompanyUserRoles, error) {
	var companyUserRoles model.CompanyUserRoles
	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&companyUserRoles)
	if err != nil {
		return nil, err
	}
	return &companyUserRoles, nil
}

//CREATE COMPANY_USER_ROLE
func (r *ProfileRepositoryMongo) SaveCompanyUserRole(companyUserRole *model.CompanyUserRole) error {
	err := r.db.C(r.collection).Insert(companyUserRole)
	return err
}

//DELETE BY USERNAME
func (r *ProfileRepositoryMongo) DeleteCompanyUserRole(username string) error {
	err := r.db.C(r.collection).Remove(bson.M{"username": username})
	return err
}
func (r *ProfileRepositoryMongo) UpdateCompanyUserRole(username string, companyUserRole *model.CompanyUserRole) error {
	err := r.db.C(r.collection).Update(bson.M{"username": username}, companyUserRole)
	return err
}

//<----------------Password Management----------------->
//FindUsernameByPhone
func (r *ProfileRepositoryMongo) FindUsernameByPhone(usertel string) (*model.UserSSO, error) {
	var userSSO model.UserSSO
	err := r.db.C(r.collection).Find(bson.M{"usertel": usertel}).One(&userSSO)
	if err != nil {
		return nil, err
	}
	return &userSSO, nil
}

//<----------------Password Management----------------->
//FindUsernameEmail
func (r *ProfileRepositoryMongo) FindUsernameByEmail(useremail string) (*model.UserSSO, error) {
	var userSSO model.UserSSO
	err := r.db.C(r.collection).Find(bson.M{"useremail": useremail}).One(&userSSO)
	if err != nil {
		return nil, err
	}
	return &userSSO, nil
}

//<----------------Paging----------------->
//count document
func (r *ProfileRepositoryMongo) CountDocument() (int, error) {
	numDoc, err := r.db.C(r.collection).Count()
	if err != nil {
		return -1, err
	}
	return numDoc, nil
}

//count document with filter
func (r *ProfileRepositoryMongo) CountDocumentWithFilter(fileter_key string, filter_value string) (int, error) {
	numDoc, err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Count()
	if err != nil {
		return -1, err
	}
	return numDoc, nil
}

//For User
//Paging without filter
func (r *ProfileRepositoryMongo) PagingUser(start int, stop int) (*model.UsersSSO, error) {
	var usersSSO model.UsersSSO
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&usersSSO)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &usersSSO, nil
}

//Paging without filter
func (r *ProfileRepositoryMongo) PagingUserAndFilter(fileter_key string, filter_value string, start int, stop int) (*model.UsersSSO, error) {
	var usersSSO model.UsersSSO
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&usersSSO)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &usersSSO, nil
}

//For company_info
//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyInfoDriver(start int, stop int) (*model.CompanyInfos, error) {
	var companyInfos model.CompanyInfos
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyInfos)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyInfos, nil
}

//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyInfoAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyInfos, error) {
	var companyInfos model.CompanyInfos
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyInfos)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyInfos, nil
}

//Company Product
//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyProductDriver(start int, stop int) (*model.CompanyProducts, error) {
	var companyProducts model.CompanyProducts
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyProducts)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyProducts, nil
}

//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyProductAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyProducts, error) {
	var companyProducts model.CompanyProducts
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyProducts)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyProducts, nil
}

//Product Integrated
//Paging without filter
func (r *ProfileRepositoryMongo) PagingProductIntegratedDriver(start int, stop int) (*model.ProductIntegrateds, error) {
	var productIntegrateds model.ProductIntegrateds
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&productIntegrateds)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &productIntegrateds, nil
}

//Paging with filter
func (r *ProfileRepositoryMongo) PagingProductIntegratedAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.ProductIntegrateds, error) {
	var productIntegrateds model.ProductIntegrateds
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&productIntegrateds)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &productIntegrateds, nil
}

//
//Company_Role
//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyRoleDriver(start int, stop int) (*model.CompanyRoles, error) {
	var companyRoles model.CompanyRoles
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyRoles)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyRoles, nil
}

//Paging with filter
func (r *ProfileRepositoryMongo) PagingCompanyRoleAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyRoles, error) {
	var companyRoles model.CompanyRoles
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyRoles)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyRoles, nil
}

//Company_Function
//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyFunctionDriver(start int, stop int) (*model.CompanyFunctions, error) {
	var companyFunctions model.CompanyFunctions
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyFunctions)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyFunctions, nil
}

//Paging with filter
func (r *ProfileRepositoryMongo) PagingCompanyFunctionAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyFunctions, error) {
	var companyFunctions model.CompanyFunctions
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyFunctions)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyFunctions, nil
}

//Company_Role_function
//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyRoleFunctionDriver(start int, stop int) (*model.CompanyRoleFunctions, error) {
	var companyRoleFunctions model.CompanyRoleFunctions
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyRoleFunctions)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyRoleFunctions, nil
}

//Paging with filter
func (r *ProfileRepositoryMongo) PagingCompanyRoleFunctionAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyRoleFunctions, error) {
	var companyRoleFunctions model.CompanyRoleFunctions
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyRoleFunctions)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyRoleFunctions, nil
}

//Company_User_Role
//Paging without filter
func (r *ProfileRepositoryMongo) PagingCompanyUseRoleDriver(start int, stop int) (*model.CompanyUserRoles, error) {
	var companyUserRoles model.CompanyUserRoles
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyUserRoles)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyUserRoles, nil
}

//Paging with filter
func (r *ProfileRepositoryMongo) PagingCompanyUseRoleAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyUserRoles, error) {
	var companyUserRoles model.CompanyUserRoles
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyUserRoles)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyUserRoles, nil
}

//FindAll User Service by Username
func (r *ProfileRepositoryMongo) FindAllUserSeriveByUsername(username string) (model.CompanyUserServices, error) {
	var companyUserServices model.CompanyUserServices

	err := r.db.C(r.collection).Find(bson.M{"username": username}).All(&companyUserServices)
	if err != nil {
		return nil, err
	}
	return companyUserServices, nil
}

//FindAll User Service by ComID
func (r *ProfileRepositoryMongo) FindAllUserSeriveByComid(comid string) (model.CompanyUserServices, error) {
	var companyUserServices model.CompanyUserServices

	err := r.db.C(r.collection).Find(bson.M{"comid": comid}).All(&companyUserServices)
	if err != nil {
		return nil, err
	}
	return companyUserServices, nil
}

//FindAll User Service
func (r *ProfileRepositoryMongo) FindAllUserSerive() (model.CompanyUserServices, error) {
	var companyUserServices model.CompanyUserServices

	err := r.db.C(r.collection).Find(bson.M{}).All(&companyUserServices)
	if err != nil {
		return nil, err
	}
	return companyUserServices, nil
}

//Paging without filter
func (r *ProfileRepositoryMongo) PagingUserServiceDriver(start int, stop int) (*model.CompanyUserServices, error) {
	var companyUserServices model.CompanyUserServices
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(nil).Limit(limit).Skip(skip).All(&companyUserServices)
	//All(&userSSO)
	if err != nil {
		return nil, err
	}
	return &companyUserServices, nil
}

//Paging with filter
func (r *ProfileRepositoryMongo) PagingUserServiceAndFilterDriver(fileter_key string, filter_value string, start int, stop int) (*model.CompanyUserServices, error) {
	var companyUserServices model.CompanyUserServices
	limit := (stop - start) + 1
	var skip int
	if start > 0 {
		skip = start - 1
	} else {
		skip = 0
	}
	err := r.db.C(r.collection).Find(bson.M{fileter_key: bson.M{"$regex": filter_value}}).Limit(limit).Skip(skip).All(&companyUserServices)
	if err != nil {
		return nil, err
	}
	return &companyUserServices, nil
}

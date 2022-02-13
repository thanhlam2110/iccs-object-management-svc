package repository

import "bitbucket.org/cloud-platform/vnpt-sso-usermgnt/model"

//IProfileRepository - IProfileRepository
type IProfileRepository interface {
	Save(*model.UserSSO) error
	Update(string, *model.UserSSO) error
	Delete(string) error
	FindByUser(string) (*model.UserSSO, error)
	FindAll(comid string) (model.UsersSSO, error)
	//ROLE
	//FindRoleByUser(string) (*model.CompanyUserRole, error)
	//PRODUCT_INTEGRATED
	SaveProductIntegrated(*model.ProductIntegrated) error
	FindProductIntegrated(string) (*model.ProductIntegrated, error)
	UpdateProductIntegrated(string, *model.ProductIntegrated) error
	DeleteProductIntegrated(string) error
	//COMPANY PRODUCT
	SaveCompanyProduct(*model.CompanyProduct) error
	FindCompanyProduct(comid string) (*model.CompanyProducts, error)
	FindCompanyProductContract(contractcode string) (*model.CompanyProduct, error)
	//DeleteCompanyProduct(comid string) error
	DeleteCompanyProductProductId(productid string) error
	UpdateCompanyProduct(string, *model.CompanyProduct) error
	//UpdateCompanyProductStatus(comid string, contractstatus string) error
	//COMPANY INFO
	SaveCompanyInfo(companyInfo *model.CompanyInfo) error
	FindCompanyInfo(comid string) (*model.CompanyInfo, error)
	UpdateCompanyInfoByComId(comid string, companyInfo *model.CompanyInfo) error
	//COMPANY FUNCTION
	FindCompanyFunctionByCode(functioncode string) (*model.CompanyFunction, error)
	FindCompanyFunction(comid string) (*model.CompanyFunctions, error)
	SaveCompanyFunction(companyFunction *model.CompanyFunction) error
	DeleteCompanyFunctionByFunctionCode(functioncode string) error
	UpdateCompanyFunction(functioncode string, companyFunction *model.CompanyFunction) error
	//COMPANY ROLE
	FindCompanyRoleByCode(rolecode string) (*model.CompanyRole, error)
	FindCompanyRole(comid string) (*model.CompanyRoles, error)
	SaveCompanyRole(companyRole *model.CompanyRole) error
	UpdateCompanyRoleCode(rolecode string, companyRole *model.CompanyRole) error
	DeleteCompanyRole(rolecode string) error
	//COMPANY ROLE FUNCTION
	//DeleteCompanyFunction(comid string) error
	DeleteCompanyRoleFunction(rolecode string) error
	FindCompanyRoleFunctionFunID(functioncode string) (*model.CompanyRoleFunctions, error)
	FindCompanyRoleFunctionComId(comid string) (*model.CompanyRoleFunctions, error)
	SaveCompanyRoleFunction(companyRoleFunction *model.CompanyRoleFunction) error
	//COMPANY USER ROLE
	FindCompanyUserRoleUsername(username string) (*model.CompanyUserRole, error)
	FindCompanyUserRoleComId(comid string) (*model.CompanyUserRoles, error)
	SaveCompanyUserRole(companyUserRole *model.CompanyUserRole) error
	DeleteCompanyUserRole(username string) error
	UpdateCompanyUserRole(username string, companyUserRole *model.CompanyUserRole) error
}

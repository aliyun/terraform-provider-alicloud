package servicecatalog

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Servicecatalog service.
// The resource and data source factories stay package-private; this package is
// the only caller that assembles them into the provider maps via conns.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "servicecatalog",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_service_catalog_portfolio", Factory: resourceAliCloudServiceCatalogPortfolio},
			{TypeName: "alicloud_service_catalog_principal_portfolio_association", Factory: resourceAliCloudServiceCatalogPrincipalPortfolioAssociation},
			{TypeName: "alicloud_service_catalog_product", Factory: resourceAliCloudServiceCatalogProduct},
			{TypeName: "alicloud_service_catalog_product_portfolio_association", Factory: resourceAliCloudServiceCatalogProductPortfolioAssociation},
			{TypeName: "alicloud_service_catalog_product_version", Factory: resourceAliCloudServiceCatalogProductVersion},
			{TypeName: "alicloud_service_catalog_provisioned_product", Factory: resourceAlicloudServiceCatalogProvisionedProduct},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_service_catalog_end_user_products", Factory: dataSourceAlicloudServiceCatalogEndUserProducts},
			{TypeName: "alicloud_service_catalog_launch_options", Factory: dataSourceAlicloudServiceCatalogLaunchOptions},
			{TypeName: "alicloud_service_catalog_portfolios", Factory: dataSourceAlicloudServiceCatalogPortfolios},
			{TypeName: "alicloud_service_catalog_product_as_end_users", Factory: dataSourceAlicloudServiceCatalogProductAsEndUsers},
			{TypeName: "alicloud_service_catalog_product_versions", Factory: dataSourceAlicloudServiceCatalogProductVersions},
			{TypeName: "alicloud_service_catalog_provisioned_products", Factory: dataSourceAlicloudServiceCatalogProvisionedProducts},
		},
	}
}

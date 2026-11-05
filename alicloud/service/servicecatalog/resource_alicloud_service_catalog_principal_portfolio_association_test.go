package servicecatalog_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/servicecatalog"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ServiceCatalog PrincipalPortfolioAssociation. >>> Resource test cases, automatically generated.
// Case 恬裕测试_副本1724752307297 7656
func TestAccAliCloudServiceCatalogPrincipalPortfolioAssociation_basic7656(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_principal_portfolio_association.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudServiceCatalogPrincipalPortfolioAssociationMap7656)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return servicecatalog.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceCatalogPrincipalPortfolioAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalog%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudServiceCatalogPrincipalPortfolioAssociationBasicDependence7656)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"principal_id":   "${alicloud_ram_role.default48JHf4.id}",
					"portfolio_id":   "${alicloud_service_catalog_portfolio.defaultDaXVxI.id}",
					"principal_type": "RamRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"principal_id":   acctest.CHECKSET,
						"portfolio_id":   acctest.CHECKSET,
						"principal_type": "RamRole",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudServiceCatalogPrincipalPortfolioAssociationMap7656 = map[string]string{}

func AlicloudServiceCatalogPrincipalPortfolioAssociationBasicDependence7656(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_service_catalog_portfolio" "defaultDaXVxI" {
  provider_name  = "tianyu"
  description    = "desc"
  portfolio_name = var.name
}

resource "alicloud_ram_role" "default48JHf4" {
  name        = var.name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}


`, name)
}

// Test ServiceCatalog PrincipalPortfolioAssociation. <<< Resource test cases, automatically generated.

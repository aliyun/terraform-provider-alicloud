package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApiGatewayDataset_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_dataset.default"
	checkoutSupportedRegions(t, true, connectivity.ApiGatewaySupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayDatasetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewaydataset%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayDatasetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dataset_name": name,
					"dataset_type": "JWT_BLOCKING",
					"description":  "tf test acc dataset 0",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "tf testacc 0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dataset_name": name,
						"dataset_type": "JWT_BLOCKING",
						"description":  "tf test acc dataset 0",
						"tags.%":       "2",
						"tags.Created": "tfTestAcc0",
						"tags.For":     "tf testacc 0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dataset_name": name + "u",
					"dataset_type": "JWT_BLOCKING",
					"description":  "tf test acc dataset 0 updated",
					"tags": map[string]string{
						"Created": "tfTestAcc0u",
						"Env":     "tf testacc 0 updated",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dataset_name": name + "u",
						"dataset_type": "JWT_BLOCKING",
						"description":  "tf test acc dataset 0 updated",
						"tags.%":       "2",
						"tags.Created": "tfTestAcc0u",
						"tags.For":     REMOVEKEY,
						"tags.Env":     "tf testacc 0 updated",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudApiGatewayDataset_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_dataset.default"
	checkoutSupportedRegions(t, true, connectivity.ApiGatewaySupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayDatasetMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayDataset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapigatewaydataset%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayDatasetBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dataset_name": name,
					"dataset_type": "IP_WHITELIST_CIDR",
					"description":  "tf test acc dataset cidr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dataset_name": name,
						"dataset_type": "IP_WHITELIST_CIDR",
						"description":  "tf test acc dataset cidr",
						"tags.%":       REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dataset_name": name + "u",
					"dataset_type": "IP_WHITELIST_CIDR",
					"description":  "tf test acc dataset cidr updated",
					"tags": map[string]string{
						"Owner": "tfTestAcc1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dataset_name": name + "u",
						"dataset_type": "IP_WHITELIST_CIDR",
						"description":  "tf test acc dataset cidr updated",
						"tags.%":       "1",
						"tags.Owner":   "tfTestAcc1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudApiGatewayDatasetMap0 = map[string]string{
	"tags.%":       CHECKSET,
	"dataset_type": CHECKSET,
}

var AlicloudApiGatewayDatasetMap1 = map[string]string{
	"tags.%":       CHECKSET,
	"dataset_type": CHECKSET,
}

func AlicloudApiGatewayDatasetBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func AlicloudApiGatewayDatasetBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

package sddp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sddp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSDDPConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSDDPConfigMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpconfig%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSDDPConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "description",
					"lang":        "zh",
					"code":        "access_failed_cnt",
					"value":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code":  "access_failed_cnt",
						"value": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "zh",
					"value":       "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "40",
					}),
				),
			},
			//todo : SDDP's Bug. reopen after fixing
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"description": "description",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"description": "description",
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudSDDPConfig_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSDDPConfigMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpconfig%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSDDPConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"code":  "access_permission_exprie_max_days",
					"value": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code":  "access_permission_exprie_max_days",
						"value": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "40",
					}),
				),
			},
			//todo : SDDP's Bug. reopen after fixing
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"description": "description",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"description": "description",
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudSDDPConfig_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSDDPConfigMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpconfig%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSDDPConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"code":  "log_datasize_avg_days",
					"value": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"code":  "log_datasize_avg_days",
						"value": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "40",
					}),
				),
			},
			//todo : SDDP's Bug. reopen after fixing
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"description": "description",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"description": "description",
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudSDDPConfigMap0 = map[string]string{}

func AlicloudSDDPConfigBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

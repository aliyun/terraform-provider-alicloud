package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ens Eip. >>> Resource test cases, automatically generated.
// Case 5131
func TestAccAliCloudEnsEip_basic5131(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsEipMap5131)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsEip")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseip%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsEipBasicDependence5131)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":         "PayAsYouGo",
					"ens_region_id":        "cn-wuxi-telecom_unicom_cmcc-2",
					"internet_charge_type": "95BandwidthByMonth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":         "PayAsYouGo",
						"ens_region_id":        "cn-wuxi-telecom_unicom_cmcc-2",
						"internet_charge_type": "95BandwidthByMonth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "EipDescription_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "EipDescription_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_name": name,
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

var AliCloudEnsEipMap5131 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"bandwidth":   acctest.CHECKSET,
	"isp":         acctest.CHECKSET,
}

func AliCloudEnsEipBasicDependence5131(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5131  twin
func TestAccAliCloudEnsEip_basic5131_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_eip.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsEipMap5131)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsEip")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senseip%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsEipBasicDependence5131)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "EipDescription_UPDATE_autost",
					"bandwidth":            "6",
					"isp":                  "cmcc",
					"payment_type":         "PayAsYouGo",
					"ens_region_id":        "cn-wuxi-telecom_unicom_cmcc-2",
					"eip_name":             name,
					"internet_charge_type": "95BandwidthByMonth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "EipDescription_UPDATE_autost",
						"bandwidth":            "6",
						"isp":                  "cmcc",
						"payment_type":         "PayAsYouGo",
						"ens_region_id":        "cn-wuxi-telecom_unicom_cmcc-2",
						"eip_name":             name,
						"internet_charge_type": "95BandwidthByMonth",
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

// Test Ens Eip. <<< Resource test cases, automatically generated.

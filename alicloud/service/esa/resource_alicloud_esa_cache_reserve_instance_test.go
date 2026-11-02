package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA CacheReserveInstance. >>> Resource test cases, automatically generated.
// Test ESA CacheReserveInstance. <<< Resource test cases, automatically generated.// Test Esa CacheReserveInstance. >>> Resource test cases, automatically generated.
// Case CacheReserveInstance_test 10493
func TestAccAliCloudEsaCacheReserveInstance_basic10493(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_cache_reserve_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEsaCacheReserveInstanceMap10493)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCacheReserveInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEsaCacheReserveInstanceBasicDependence10493)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_gb":     "10240",
					"cr_region":    "CN-beijing",
					"auto_renew":   "true",
					"period":       "1",
					"payment_type": "Subscription",
					"auto_pay":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_gb":     "10240",
						"cr_region":    "CN-beijing",
						"auto_renew":   "true",
						"period":       "1",
						"payment_type": "Subscription",
						"auto_pay":     "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_gb": "51200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_gb": "51200",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew"},
			},
		},
	})
}

var AlicloudEsaCacheReserveInstanceMap10493 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudEsaCacheReserveInstanceBasicDependence10493(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Esa CacheReserveInstance. <<< Resource test cases, automatically generated.

package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudOnsStaticAccount_basic1775(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_static_account.default"
	ra := resourceAttrInit(resourceId, AliCloudOnsStaticAccountMap1775)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpStaticAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sAmqpStaticAccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudOnsStaticAccountBasicDependence1775)
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
					"instance_id": "${data.alicloud_amqp_instances.default.ids.0}",
					"access_key":  os.Getenv("ALICLOUD_ACCESS_KEY"),
					"secret_key":  os.Getenv("ALICLOUD_SECRET_KEY"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

var AliCloudOnsStaticAccountMap1775 = map[string]string{
	"user_name":   CHECKSET,
	"password":    CHECKSET,
	"master_uid":  CHECKSET,
	"create_time": CHECKSET,
}

func AliCloudOnsStaticAccountBasicDependence1775(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_amqp_instances" "default" {
  		status = "SERVING"
	}
`, name)
}

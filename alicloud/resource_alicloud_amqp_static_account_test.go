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
func TestAccAlicloudOnsStaticAccount_basic1775(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_static_account.default"
	ra := resourceAttrInit(resourceId, AlicloudOnsStaticAccountMap1775)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAmqpStaticAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sAmqpStaticAccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOnsStaticAccountBasicDependence1775)
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
						"user_name":   CHECKSET,
						"instance_id": CHECKSET,
						"create_time": CHECKSET,
						"master_uid":  CHECKSET,
						"password":    CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

var AlicloudOnsStaticAccountMap1775 = map[string]string{}

func AlicloudOnsStaticAccountBasicDependence1775(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_amqp_instances" "default" {
	status = "SERVING"
}
`, name)
}

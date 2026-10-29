package amqpopen_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/amqpopen"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudOnsStaticAccount_basic1775(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_amqp_static_account.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudOnsStaticAccountMap1775)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return amqpopen.NewAmqpOpenService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAmqpStaticAccount")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sAmqpStaticAccount%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudOnsStaticAccountBasicDependence1775)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_amqp_instances.default.ids.0}",
					"access_key":  os.Getenv("ALICLOUD_ACCESS_KEY"),
					"secret_key":  os.Getenv("ALICLOUD_SECRET_KEY"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
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
	"user_name":   acctest.CHECKSET,
	"password":    acctest.CHECKSET,
	"master_uid":  acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
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

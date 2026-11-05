package smartag_test

import (
	"fmt"
	"os"
	"testing"

	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/smartag"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudConnectNetworkAttachment_basic(t *testing.T) {
	var sag smartagsdk.SmartAccessGateway
	resourceId := "alicloud_cloud_connect_network_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, ccnAttachmentMap)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &sag, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetworkAttachment-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCcnAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheckWithEnvVariable(t, "SAG_INSTANCE_ID")
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ccn_id":     "${alicloud_cloud_connect_network.ccn.id}",
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"depends_on": []string{"alicloud_cloud_connect_network.ccn"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ccn_id": acctest.CHECKSET,
						"sag_id": os.Getenv("SAG_INSTANCE_ID"),
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

var ccnAttachmentMap = map[string]string{
	"ccn_id": acctest.CHECKSET,
}

func resourceCcnAttachmentDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	resource "alicloud_cloud_connect_network" "ccn" {
	  	name = "${var.name}"
	  	is_default = "true"
	}
`, name)
}

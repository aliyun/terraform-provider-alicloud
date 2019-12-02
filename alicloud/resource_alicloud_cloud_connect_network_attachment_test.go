package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCloudConnectNetworkAttachment_basic(t *testing.T) {
	var sag smartag.SmartAccessGateway
	resourceId := "alicloud_cloud_connect_network_attachment.default"
	ra := resourceAttrInit(resourceId, ccnAttachmentMap)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &sag, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCloudConnectNetworkAttachment-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCcnAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithSmartAccessGatewaySetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ccn_id":     "${alicloud_cloud_connect_network.ccn.id}",
					"sag_id":     os.Getenv("SAG_INSTANCE_ID"),
					"depends_on": []string{"alicloud_cloud_connect_network.ccn"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ccn_id": CHECKSET,
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
	"ccn_id": CHECKSET,
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

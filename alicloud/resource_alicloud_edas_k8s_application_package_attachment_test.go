package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudEdasK8sApplicationPackageAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_k8s_application_deployment.default"
	ra := resourceAttrInit(resourceId, edasK8sAPAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-edask8sdeploymentbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sAPAttachmentDependence)
	region := os.Getenv("ALICLOUD_REGION")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckK8sDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":    "0204f66d-9eef-4370-8e29-cb364462a6bc",
					"replicas":  "1",
					"image_url": fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/consumer:1.0", region),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testEdasCheckK8sDeploymentDestroy(s *terraform.State) error {
	return nil
}

var edasK8sAPAttachmentBasicMap = map[string]string{
	"app_id":    CHECKSET,
	"replicas":  CHECKSET,
	"image_url": CHECKSET,
}

func resourceEdasK8sAPAttachmentDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		`, name)
}

package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSImageImport(t *testing.T) {
	var v ecs.Image
	resourceId := "alicloud_image_import.default"
	ra := resourceAttrInit(resourceId, testAccImageImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rand := acctest.RandIntRange(1000, 9999)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageImageBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOSSForImageImport(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
					"image_name":   name,
					"architecture": "x86_64",
					"license_type": "Auto",
					"platform":     "Ubuntu",
					"os_type":      "linux",
					"disk_device_mapping": []map[string]interface{}{
						{
							"oss_bucket": os.Getenv("ALICLOUD_OSS_BUCKET_FOR_IMAGE"),
							"oss_object": os.Getenv("ALICLOUD_OSS_OBJECT_FOR_IMAGE"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
						"image_name":                       name,
						"architecture":                     "x86_64",
						"license_type":                     "Auto",
						"platform":                         "Ubuntu",
						"os_type":                          "linux",
						"disk_device_mapping.#":            "1",
						"disk_device_mapping.0.oss_bucket": CHECKSET,
						"disk_device_mapping.0.oss_object": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescriptionchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescriptionchange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%dchange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
					"image_name":  fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
						"image_name":  fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%d", rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"license_type"},
			},
		},
	})

}

var testAccImageImageCheckMap = map[string]string{}

func resourceImageImageBasicConfigDependence(name string) string {
	return ""
}

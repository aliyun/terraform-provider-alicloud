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

func TestAccAlicloudSagClientUser_basic(t *testing.T) {
	var user smartag.User
	resourceId := "alicloud_sag_client_user.default"
	ra := resourceAttrInit(resourceId, sagClientUserMap)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &user, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-username-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagClientUserDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithSmartAccessGatewayAppSetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":    os.Getenv("SAG_APP_INSTANCE_ID"),
					"bandwidth": "20",
					"user_mail": "${var.name}@test.com",
					"user_name": "${var.name}",
					"password":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sag_id":    os.Getenv("SAG_APP_INSTANCE_ID"),
						"bandwidth": "20",
						"user_mail": fmt.Sprintf("%s@test.com", name),
						"user_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "20",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSagClientUser_multi(t *testing.T) {
	var user smartag.User
	resourceId := "alicloud_sag_client_user.default.9"
	ra := resourceAttrInit(resourceId, sagClientUserMap)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &user, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-username-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagClientUserDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			testAccPreCheckWithSmartAccessGatewayAppSetting(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sag_id":    os.Getenv("SAG_APP_INSTANCE_ID"),
					"count":     "10",
					"bandwidth": "20",
					"user_mail": "${var.name}-${count.index}@test.com",
					"user_name": "${var.name}-${count.index}",
					"password":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var sagClientUserMap = map[string]string{
	"sag_id": CHECKSET,
}

func resourceSagClientUserDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

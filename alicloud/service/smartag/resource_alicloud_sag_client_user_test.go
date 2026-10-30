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

func TestAccAliCloudSagClientUser_basic(t *testing.T) {
	var user smartagsdk.User
	resourceId := "alicloud_sag_client_user.default"
	ra := acctest.ResourceAttrInit(resourceId, sagClientUserMap)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &user, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-username-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSagClientUserDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheckWithEnvVariable(t, "SAG_APP_INSTANCE_ID")
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

func TestAccAliCloudSagClientUser_multi(t *testing.T) {
	var user smartagsdk.User
	resourceId := "alicloud_sag_client_user.default.9"
	ra := acctest.ResourceAttrInit(resourceId, sagClientUserMap)
	serviceFunc := func() interface{} {
		return smartag.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &user, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-username-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceSagClientUserDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
			acctest.PreCheckWithEnvVariable(t, "SAG_APP_INSTANCE_ID")
		},

		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"sag_id": acctest.CHECKSET,
}

func resourceSagClientUserDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

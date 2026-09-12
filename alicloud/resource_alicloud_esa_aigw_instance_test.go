package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEsaAigwInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_aigw_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaAigwInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaAigwInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccaigw%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaAigwInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aigw_instance_name": name,
					"comment":            "test-comment-basic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aigw_instance_name": name,
						"comment":            "test-comment-basic",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aigw_instance_name": name,
					"comment":            "test-comment-updated",
					"enable_auth":        "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aigw_instance_name": name,
						"comment":            "test-comment-updated",
						"enable_auth":        "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aigw_instance_name": name,
					"comment":            "test-comment-updated",
					"enable_auth":        "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aigw_instance_name": name,
						"comment":            "test-comment-updated",
						"enable_auth":        "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEsaAigwInstance_datasource(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_aigw_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaAigwInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaAigwInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccaigw%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaAigwInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aigw_instance_name": name,
					"comment":            "test-comment-ds",
				}) + testAccAlicloudEsaAigwInstanceDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aigw_instance_name": name,
						"comment":            "test-comment-ds",
					}),
					resource.TestCheckResourceAttr("data.alicloud_esa_aigw_instances.default", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccAlicloudEsaAigwInstanceDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "alicloud_esa_aigw_instances" "default" {
  ids = [alicloud_esa_aigw_instance.default.aigw_instance_id]
}
`)
}

var AliCloudEsaAigwInstanceMap = map[string]string{
	"aigw_instance_id": CHECKSET,
	"auth_key":         CHECKSET,
	"create_time":      CHECKSET,
	"record_count":     CHECKSET,
	"status":           CHECKSET,
	"update_time":      CHECKSET,
}

func AliCloudEsaAigwInstanceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

`, name)
}

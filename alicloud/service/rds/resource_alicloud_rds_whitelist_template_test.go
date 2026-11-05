package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudWhitelistTemplate(t *testing.T) {
	var WhitelistTemplate map[string]interface{}
	resourceId := "alicloud_rds_whitelist_template.default"
	var WhitelistTemplateMap = map[string]string{
		"ip_white_list": acctest.CHECKSET,
		"template_name": acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, WhitelistTemplateMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &WhitelistTemplate, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeWhitelistTemplate")
	rand := tfacctest.RandString(5)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-%s", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, WhitelistTemplateBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": "127.0.0.1",
					"template_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list": acctest.CHECKSET,
						"template_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_name": fmt.Sprintf("%s-updated", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white_list": "127.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list": acctest.CHECKSET,
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

func WhitelistTemplateBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}

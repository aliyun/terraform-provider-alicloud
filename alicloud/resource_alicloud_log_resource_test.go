package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogResource_basic(t *testing.T) {
	var v *sls.Resource
	resourceId := "alicloud_log_resource.default"
	ra := resourceAttrInit(resourceId, logResourceMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclog-resource-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogResourceDependence)

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
					"name":             "user.user_whitelist_ip",
					"schema":           `{"schema": [{"desc": "whitelist ip range", "type": "string", "column": "ip", "indexed": false, "ext_info": {}, "required": true}, {"desc": "0:ip, 1:ip range", "type": "string", "column": "type", "indexed": false, "ext_info": "optional", "required": true}]}`,
					"ext_info":         `{}`,
					"create_time":      "1646314020",
					"last_modify_time": "1646314020",
					"description":      "all whitelist",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             "user.user_whitelist_ip",
						"schema":           `{"schema": [{"desc": "whitelist ip range", "type": "string", "column": "ip", "indexed": false, "ext_info": {}, "required": true}, {"desc": "0:ip, 1:ip range", "type": "string", "column": "type", "indexed": false, "ext_info": "optional", "required": true}]}`,
						"ext_info":         `{}`,
						"create_time":      "1646314020",
						"last_modify_time": "1646314020",
						"description":      "all whitelist",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     "user.user_whitelist_ip",
					"ext_info": `{"ext1": "ext1_val"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     "user.user_whitelist_ip",
						"ext_info": `{"ext1": "ext1_val"}`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "user.user_whitelist_ip",
					"create_time": "1646314520",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        "user.user_whitelist_ip",
						"create_time": "1646314520",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "user.user_whitelist_ip",
					"description": `all whitelist new`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        "user.user_whitelist_ip",
						"description": `all whitelist new`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   "user.user_whitelist_ip",
					"schema": `{"schema": [{"desc": "whitelist ip range new", "type": "string", "column": "ip", "indexed": false, "ext_info": {}, "required": true}, {"desc": "0:ip, 1:ip range", "type": "string", "column": "type", "indexed": false, "ext_info": "optional", "required": true}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   "user.user_whitelist_ip",
						"schema": `{"schema": [{"desc": "whitelist ip range new", "type": "string", "column": "ip", "indexed": false, "ext_info": {}, "required": true}, {"desc": "0:ip, 1:ip range", "type": "string", "column": "type", "indexed": false, "ext_info": "optional", "required": true}]}`,
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

var logResourceMap = map[string]string{}

func resourceLogResourceDependence(name string) string {
	return ""
}

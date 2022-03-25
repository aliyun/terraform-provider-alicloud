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
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	logResourceType := "userdefine"

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testacclog_resource_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogResourceDependence)

	resource.Test(t, resource.TestCase{

		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.LogResourceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":        logResourceType,
					"name":        name,
					"schema":      `{\"schema\":[{\"column\":\"ip\",\"required\":true,\"indexed\":false,\"desc\":\"whitelist ip range\",\"type\":\"string\",\"ext_info\":{}},{\"column\":\"type\",\"required\":true,\"indexed\":false,\"desc\":\"0:ip, 1:ip range\",\"type\":\"string\",\"ext_info\":\"optional\"}]}`,
					"ext_info":    "{}",
					"description": "all whitelist",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        logResourceType,
						"name":        name,
						"schema":      CHECKSET,
						"ext_info":    "{}",
						"description": "all whitelist",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name,
					"ext_info": `{\"ext1\": \"ext1_val\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"ext_info": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "all whitelist new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "all whitelist new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":   name,
					"schema": `{\"schema\": [{\"desc\": \"whitelist ip range new\", \"type\": \"string\", \"column\": \"ip\", \"indexed\": false, \"ext_info\": {}, \"required\": true}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   name,
						"schema": CHECKSET,
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

var logResourceMap = map[string]string{
	"name": CHECKSET,
}

func resourceLogResourceDependence(name string) string {
	return `
	`
}

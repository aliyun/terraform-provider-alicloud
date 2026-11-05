package sls_test

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	slssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogResource_basic(t *testing.T) {
	var v *sls.Resource
	resourceId := "alicloud_log_resource.default"
	ra := acctest.ResourceAttrInit(resourceId, logResourceMap)
	serviceFunc := func() interface{} {
		return slssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	logResourceType := "userdefine"

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testacclog_resource_%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogResourceDependence)

	resource.Test(t, resource.TestCase{

		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.LogResourceSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"schema":      acctest.CHECKSET,
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
						"ext_info": acctest.CHECKSET,
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
						"schema": acctest.CHECKSET,
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
	"name": acctest.CHECKSET,
}

func resourceLogResourceDependence(name string) string {
	return `
	`
}

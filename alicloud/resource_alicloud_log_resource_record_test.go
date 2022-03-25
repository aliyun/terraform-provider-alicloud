package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogResourceRecord_basic(t *testing.T) {
	var v *sls.ResourceRecord
	resourceId := "alicloud_log_resource_record.default"
	ra := resourceAttrInit(resourceId, logResourceRecordMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclog-resource-record-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogResourceRecordDependence)

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
					"record_id":     name,
					"tag":           "tag name",
					"value":         `{\"email\": [\"aaa@xxx.com\"], \"phone\": \"18888888888\", \"enabled\": true, \"user_id\": \"test_214958\", \"user_name\": \"test_name\", \"sms_enabled\": true, \"country_code\": \"86\", \"voice_enabled\": false}`,
					"resource_name": "sls.common.user",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_id":     name,
						"tag":           "tag name",
						"value":         CHECKSET,
						"resource_name": "sls.common.user",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"record_id": name,
					"tag":       "tag name new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_id": name,
						"tag":       "tag name new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"record_id": name,
					"value":     `{\"email\": [\"bbb@xxx.com\"], \"phone\": \"18888888888\", \"enabled\": true, \"user_id\": \"test_214958\", \"user_name\": \"test_name\", \"sms_enabled\": true, \"country_code\": \"86\", \"voice_enabled\": false}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_id": name,
						"value":     CHECKSET,
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

var logResourceRecordMap = map[string]string{
	"record_id": CHECKSET,
	"tag":       CHECKSET,
}

func resourceLogResourceRecordDependence(name string) string {
	return `
	`
}

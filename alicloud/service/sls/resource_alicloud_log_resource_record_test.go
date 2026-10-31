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

func TestAccAliCloudLogResourceRecord_basic(t *testing.T) {
	var v *sls.ResourceRecord
	resourceId := "alicloud_log_resource_record.default"
	ra := acctest.ResourceAttrInit(resourceId, logResourceRecordMap)
	serviceFunc := func() interface{} {
		return slssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclog-resource-record-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogResourceRecordDependence)

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
					"record_id":     name,
					"tag":           "tag name",
					"value":         `{\"email\": [\"aaa@xxx.com\"], \"phone\": \"18888888888\", \"enabled\": true, \"user_id\": \"test_214958\", \"user_name\": \"test_name\", \"sms_enabled\": true, \"country_code\": \"86\", \"voice_enabled\": false}`,
					"resource_name": "sls.common.user",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_id":     name,
						"tag":           "tag name",
						"value":         acctest.CHECKSET,
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
						"value":     acctest.CHECKSET,
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

func TestAccAliCloudLogResourceRecord_resourceName_update(t *testing.T) {
	var v *sls.ResourceRecord
	resourceId := "alicloud_log_resource_record.default"
	ra := acctest.ResourceAttrInit(resourceId, logResourceRecordMap)
	serviceFunc := func() interface{} {
		return slssvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclog-resource-record-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceLogResourceRecordDependence)

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
					"record_id":     name,
					"tag":           "tag name",
					"value":         `{\"email\": [\"aaa@xxx.com\"], \"phone\": \"18888888888\", \"enabled\": true, \"user_id\": \"test_214958\", \"user_name\": \"test_name\", \"sms_enabled\": true, \"country_code\": \"86\", \"voice_enabled\": false}`,
					"resource_name": "sls.common.user",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_id":     name,
						"tag":           "tag name",
						"resource_name": "sls.common.user",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"record_id":     name,
					"tag":           "tag name",
					"value":         `{\"email\": [\"aaa@xxx.com\"], \"phone\": \"18888888888\", \"enabled\": true, \"user_id\": \"test_214958\", \"user_name\": \"test_name\", \"sms_enabled\": true, \"country_code\": \"86\", \"voice_enabled\": false}`,
					"resource_name": "sls.common.alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"record_id":     name,
						"tag":           "tag name",
						"resource_name": "sls.common.alert",
					}),
				),
			},
		},
	})
}

var logResourceRecordMap = map[string]string{
	"record_id": acctest.CHECKSET,
	"tag":       acctest.CHECKSET,
}

func resourceLogResourceRecordDependence(name string) string {
	return `
	`
}

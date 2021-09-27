package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudIMPAppTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_imp_app_template.default"
	ra := resourceAttrInit(resourceId, AlicloudIMPAppTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImpAppTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%simpapptemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudIMPAppTemplateBasicDependence0)
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
					"scene":             "business",
					"integration_mode":  "paasSDK",
					"component_list":    []string{"component.live"},
					"app_template_name": "tf_testAcc_GWcpq51dSi5td18Qd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scene":             "business",
						"integration_mode":  "paasSDK",
						"component_list.#":  "1",
						"app_template_name": "tf_testAcc_GWcpq51dSi5td18Qd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_template_name": "tf_testAcc_IN1u0gHPAo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_template_name": "tf_testAcc_IN1u0gHPAo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_list": []map[string]interface{}{
						{
							"key":   "config.appCallbackAuthKey",
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV",
						},
						{
							"key":   "config.appCallbackUrl",
							"value": "http://aliyun.com/tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV",
						},
						{
							"key":   "config.livePullDomain",
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV.com",
						},
						{
							"key":   "config.livePushDomain",
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTD.com",
						},
						{
							"key":   "config.regionId",
							"value": "cn-hangzhou",
						},
						{
							"key":   "config.streamChangeCallbackUrl",
							"value": "https://aliyun.com/tf-testAcc-jdD4qhGOujVlYcCU",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_list.#": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_template_name": "tf_testAcc_tqPHQU5xU",
					"config_list": []map[string]interface{}{
						{
							"key":   "config.appCallbackAuthKey",
							"value": "tf-testAcc-jdD4qhGOxxxxxxxx",
						},
						{
							"key":   "config.appCallbackUrl",
							"value": "http://aliyun.com/tf-testAcc-jdD4qhGOxxxxxxxx",
						},
						{
							"key":   "config.livePullDomain",
							"value": "tf-testAcc-jdD4qhGOxxxxxxxx.com",
						},
						{
							"key":   "config.livePushDomain",
							"value": "tf-testAcc-jdD4qhGOxxxxxxxx.com",
						},
						{
							"key":   "config.regionId",
							"value": "cn-shanghai",
						},
						{
							"key":   "config.streamChangeCallbackUrl",
							"value": "https://aliyun.com/tf-testAcc-jdD4qhGOxxxxxxxx",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_template_name": "tf_testAcc_tqPHQU5xU",
						"config_list.#":     "6",
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

var AlicloudIMPAppTemplateMap0 = map[string]string{
	"component_list.#": CHECKSET,
	"config_list.#":    CHECKSET,
	"status":           CHECKSET,
}

func AlicloudIMPAppTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

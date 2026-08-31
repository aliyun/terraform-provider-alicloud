package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApiGatewayStageModel_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_stage_model.default"
	ra := resourceAttrInit(resourceId, AlicloudApiGatewayStageModelMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayStageModel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("TF%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApiGatewayStageModelBasicDependence0)
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
					"stage_model_name":  name,
					"stage_model_alias": "tf-testAcc-alias",
					"description":       "tf-testAcc-desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stage_model_name":  name,
						"stage_model_alias": "tf-testAcc-alias",
						"description":       "tf-testAcc-desc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stage_model_alias": "tf-testAcc-alias-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stage_model_alias": "tf-testAcc-alias-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAcc-desc-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc-desc-update",
					}),
				),
			},
		},
	})
}

var AlicloudApiGatewayStageModelMap0 = map[string]string{
	"stage_model_id": CHECKSET,
	"type":           CHECKSET,
	"created_time":   CHECKSET,
	"modified_time":  CHECKSET,
}

func AlicloudApiGatewayStageModelBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

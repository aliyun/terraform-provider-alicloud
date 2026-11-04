package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudApiGatewayStageModel_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_stage_model.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayStageModelMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewCloudApiService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayStageModel")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("TF%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayStageModelBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
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
	"stage_model_id": acctest.CHECKSET,
	"type":           acctest.CHECKSET,
	"created_time":   acctest.CHECKSET,
	"modified_time":  acctest.CHECKSET,
}

func AlicloudApiGatewayStageModelBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

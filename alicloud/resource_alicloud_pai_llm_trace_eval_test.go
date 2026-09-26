package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPaiLlmTraceEval_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_llm_trace_eval.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiLlmTraceEvalMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaillmtraceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaillmtraceEval")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-paillmtrace-eval-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiLlmTraceEvalBasicDependence0)
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
					"app_name":        "tf-test-app",
					"data_source":     "tf-test-data-source",
					"description":     "tf-test-eval-description",
					"eval_name":       name,
					"evaluation_data": "{}",
					"metadata":        "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":    "tf-test-app",
						"data_source": "tf-test-data-source",
						"description": "tf-test-eval-description",
						"eval_name":   name,
						"region_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":        "tf-test-app-updated",
					"data_source":     "tf-test-data-source-updated",
					"description":     "tf-test-eval-description-updated",
					"eval_name":       fmt.Sprintf("%s-updated", name),
					"evaluation_data": "{}",
					"metadata":        "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":    "tf-test-app-updated",
						"data_source": "tf-test-data-source-updated",
						"description": "tf-test-eval-description-updated",
						"eval_name":   fmt.Sprintf("%s-updated", name),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_name", "data_source", "evaluation_data", "metadata"},
			},
		},
	})
}

var AlicloudPaiLlmTraceEvalMap0 = map[string]string{
	"region_id":       CHECKSET,
	"gmt_create_time": CHECKSET,
	"record_count":    CHECKSET,
	"eval_id":         CHECKSET,
}

func AlicloudPaiLlmTraceEvalBasicDependence0(name string) string {
	return ""
}

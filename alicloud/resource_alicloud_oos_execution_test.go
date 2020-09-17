package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/oos"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_oos_execution", &resource.Sweeper{
		Name: "alicloud_oos_execution",
		F:    testSweepOosExecution,
	})
}

func testSweepOosExecution(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := oos.CreateListExecutionsRequest()
	raw, err := client.WithOosClient(func(OosClient *oos.Client) (interface{}, error) {
		return OosClient.ListExecutions(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Oos Executions: %s", WrapError(err))
	}
	response, _ := raw.(*oos.ListExecutionsResponse)

	sweeped := false
	for _, v := range response.Executions {
		id := v.ExecutionId
		name := v.TemplateName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Oos Executions: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting Oos Executions: %s (%s)", name, id)
		req := oos.CreateDeleteExecutionsRequest()
		req.ExecutionIds = convertListToJsonString(convertListStringToListInterface([]string{id}))
		_, err := client.WithOosClient(func(OosClient *oos.Client) (interface{}, error) {
			return OosClient.DeleteExecutions(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Oos Executions (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these Oos Executions have been deleted.
		time.Sleep(10 * time.Second)
	}
	return nil
}

func TestAccAlicloudOOSExecution_basic(t *testing.T) {
	var v oos.Execution
	resourceId := "alicloud_oos_execution.default"
	ra := resourceAttrInit(resourceId, OosExecutionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosExecution%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OosExecutionBasicdependence)
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
					"template_name": "${alicloud_oos_template.default.template_name}",
					"description":   "From TF Test",
					"parameters":    `{\"Status\":\"Running\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_name": CHECKSET,
						"description":   "From TF Test",
						"parameters":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "loop_mode", "safety_check", "parameters"},
			},
		},
	})
}

var OosExecutionMap = map[string]string{
	"create_date": CHECKSET,
	"executed_by": CHECKSET,
	"is_parent":   CHECKSET,
	"start_date":  CHECKSET,
	"status":      CHECKSET,
	"template_id": CHECKSET,
	"update_date": CHECKSET,
}

func OosExecutionBasicdependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Describe instances of given status",
			"Parameters":{
			  "Status":{
				"Type": "String",
				"Description": "(Required) The status of the Ecs instance."
			  }
			},
			"Tasks": [
			  {
				"Properties" :{
				  "Parameters":{
					"Status": "{{ Status }}"
				  },
				  "API": "DescribeInstances",
				  "Service": "Ecs"
				},
				"Name": "foo",
				"Action": "ACS::ExecuteApi"
			  }]
		  }
		  EOF
		  template_name = "%s"
		  version_name = "test"
		  tags = {
			"Created" = "TF",
			"For" = "template Test"
		  }
		}
	`, name)
}

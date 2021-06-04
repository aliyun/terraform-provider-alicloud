package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_oos_template", &resource.Sweeper{
		Name: "alicloud_oos_template",
		F:    testSweepOosTemplate,
	})
}

func testSweepOosTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	action := "ListTemplates"
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oos_template", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Templates", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Templates", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["TemplateName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping OOS Template: %s", item["TemplateName"].(string))
				continue
			}
			sweeped = true
			action = "DeleteTemplate"
			request := map[string]interface{}{
				"TemplateName": item["TemplateName"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete OOS Template (%s): %s", item["TemplateName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure OOS Template have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete OOS Template success: %s ", item["TemplateName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudOosTemplate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_template.default"
	ra := resourceAttrInit(resourceId, OosTemplateMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosTemplate%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, OosTemplateBasicdependence)
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
					"content":       `{\"FormatVersion\":\"OOS-2019-06-01\",\"Description\":\"Describe instances of given status\",\"Parameters\":{\"Status\":{\"Type\":\"String\",\"Description\":\"(Required) The status of the Ecs instance.\"}},\"Tasks\":[{\"Properties\":{\"Parameters\":{\"Status\":\"{{ Status }}\"},\"API\":\"DescribeInstances\",\"Service\":\"Ecs\"},\"Name\":\"foo\",\"Action\":\"ACS::ExecuteApi\"}]}`,
					"template_name": name,
					"version_name":  "test1",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":      CHECKSET,
						"version_name": "test1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_delete_executions", "content", "version_name"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": `{\"FormatVersion\":\"OOS-2019-06-01\",\"Description\":\"Update Describe instances of given status\",\"Parameters\":{\"Status\":{\"Type\":\"String\",\"Description\":\"(Required) The status of the Ecs instance.\"}},\"Tasks\":[{\"Properties\":{\"Parameters\":{\"Status\":\"{{ Status }}\"},\"API\":\"DescribeInstances\",\"Service\":\"Ecs\"},\"Name\":\"foo\",\"Action\":\"ACS::ExecuteApi\"}]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content":      `{\"FormatVersion\":\"OOS-2019-06-01\",\"Description\":\"Update Describe instances of given status\",\"Parameters\":{\"Status\":{\"Type\":\"String\",\"Description\":\"(Required) The status of the Ecs instance.\"}},\"Tasks\":[{\"Properties\":{\"Parameters\":{\"Status\":\"{{ Status }}\"},\"API\":\"DescribeInstances\",\"Service\":\"Ecs\"},\"Name\":\"foo\",\"Action\":\"ACS::ExecuteApi\"}]}`,
					"version_name": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":      CHECKSET,
						"version_name": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": `{\"FormatVersion\":\"OOS-2019-06-01\",\"Description\":\"Update Describe instances of given status\",\"Parameters\":{\"Status\":{\"Type\":\"String\",\"Description\":\"(Required) The status of the Ecs instance.\"}},\"Tasks\":[{\"Properties\":{\"Parameters\":{\"Status\":\"{{ Status }}\"},\"API\":\"DescribeInstances\",\"Service\":\"Ecs\"},\"Name\":\"foo\",\"Action\":\"ACS::ExecuteApi\"}]}`,
					"tags": map[string]string{
						"Created": "TF-Test",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":      CHECKSET,
						"tags.%":       "2",
						"tags.Created": "TF-Test",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": `{\"FormatVersion\":\"OOS-2019-06-01\",\"Description\":\"Describe instances of given status\",\"Parameters\":{\"Status\":{\"Type\":\"String\",\"Description\":\"(Required) The status of the Ecs instance.\"}},\"Tasks\":[{\"Properties\":{\"Parameters\":{\"Status\":\"{{ Status }}\"},\"API\":\"DescribeInstances\",\"Service\":\"Ecs\"},\"Name\":\"foo\",\"Action\":\"ACS::ExecuteApi\"}]}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
					"version_name": "test3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":      CHECKSET,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
						"version_name": "test3",
					}),
				),
			},
		},
	})
}

var OosTemplateMap = map[string]string{
	"created_by":      CHECKSET,
	"created_date":    CHECKSET,
	"description":     CHECKSET,
	"has_trigger":     CHECKSET,
	"share_type":      CHECKSET,
	"template_format": CHECKSET,
	"template_id":     CHECKSET,
	"template_type":   CHECKSET,
	"updated_by":      CHECKSET,
	"updated_date":    CHECKSET,
}

func OosTemplateBasicdependence(name string) string {
	return ""
}

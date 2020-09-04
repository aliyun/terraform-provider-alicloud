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
	resource.AddTestSweepers("alicloud_oos_template", &resource.Sweeper{
		Name: "alicloud_oos_template",
		F:    testSweepOosTemplate,
	})
}

func testSweepOosTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := oos.CreateListTemplatesRequest()
	raw, err := client.WithOosClient(func(OosClient *oos.Client) (interface{}, error) {
		return OosClient.ListTemplates(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Oos Templates: %s", WrapError(err))
	}
	response, _ := raw.(*oos.ListTemplatesResponse)

	sweeped := false
	for _, v := range response.Templates {
		id := v.TemplateId
		name := v.TemplateName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Oos Templates: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting Oos Templates: %s (%s)", name, id)
		req := oos.CreateDeleteTemplateRequest()
		req.TemplateName = name
		_, err := client.WithOosClient(func(OosClient *oos.Client) (interface{}, error) {
			return OosClient.DeleteTemplate(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Oos Templates (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these Oos Templates have been deleted.
		time.Sleep(10 * time.Second)
	}
	return nil
}

func TestAccAlicloudOOSTemplate_basic(t *testing.T) {
	var v oos.Template
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

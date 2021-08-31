package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alb_health_check_template",
		&resource.Sweeper{
			Name: "alicloud_alb_health_check_template",
			F:    testSweepAlbHealthCheckTemplate,
		})
}

func testSweepAlbHealthCheckTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListHealthCheckTemplates"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.HealthCheckTemplates", response)

		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.HealthCheckTemplates", action, err)
			return nil
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["HealthCheckTemplateName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ALB HealthCheckTemplate: %s", item["HealthCheckTemplateName"].(string))
				continue
			}

			sweeped = true
			action := "DeleteHealthCheckTemplates"
			request := map[string]interface{}{
				"HealthCheckTemplateIds.1": item["HealthCheckTemplateId"],
			}
			request["ClientToken"] = buildClientToken("DeleteHealthCheckTemplate")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ALB HealthCheckTemplate (%s): %s", item["HealthCheckTemplateName"].(string), err)
			}
			if sweeped {
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete ALB HealthCheckTemplate success: %s ", item["HealthCheckTemplateName"].(string))
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}

func TestAccAlicloudALBHealthCheckTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_health_check_template.default"
	ra := resourceAttrInit(resourceId, AlicloudALBHealthCheckTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbHealthCheckTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbhealthchecktemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBHealthCheckTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name,
					"health_check_protocol":      "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name,
						"health_check_protocol":      "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_host": "www.test.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_host": "www.test.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_http_version": "HTTP1.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_version": "HTTP1.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_method": "GET",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_method": "GET",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_codes": []string{"http_3xx", "http_4xx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_codes.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "8080",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "8080",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "TCP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "TCP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_name": name,
					"health_check_protocol":      "HTTP",
					"health_check_connect_port":  "0",
					"health_check_host":          "www.test2.com",
					"health_check_http_version":  "HTTP1.1",
					"health_check_interval":      "2",
					"health_check_method":        "HEAD",
					"health_check_path":          "/test2",
					"health_check_timeout":       "5",
					"healthy_threshold":          "3",
					"unhealthy_threshold":        "3",
					"health_check_codes":         []string{"http_2xx", "http_5xx"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_name": name,
						"health_check_protocol":      "HTTP",
						"health_check_connect_port":  "0",
						"health_check_host":          "www.test2.com",
						"health_check_http_version":  "HTTP1.1",
						"health_check_interval":      "2",
						"health_check_method":        "HEAD",
						"health_check_path":          "/test2",
						"health_check_timeout":       "5",
						"healthy_threshold":          "3",
						"unhealthy_threshold":        "3",
						"health_check_codes.#":       "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBHealthCheckTemplateMap0 = map[string]string{}

func AlicloudALBHealthCheckTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

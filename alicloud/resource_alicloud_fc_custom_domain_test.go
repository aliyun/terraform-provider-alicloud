package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_fc_custom_domain", &resource.Sweeper{
		Name: "alicloud_fc_custom_domain",
		F:    testSweepFCCustomDomain,
		Dependencies: []string{
			"alicloud_fc_custom_domain",
		},
	})
}

func testSweepFCCustomDomain(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.ListServices(fc.NewListServicesInput())
	})
	if err != nil {
		return fmt.Errorf("Error retrieving FC services: %s", err)
	}

	// Delete FC function and services.
	services, _ := raw.(*fc.ListServicesOutput)
	for _, v := range services.Services {
		name := *v.ServiceName
		id := *v.ServiceID
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping FC services: %s (%s)", name, id)
			continue
		}
		// Remove functions
		nextToken := ""
		for {
			args := fc.NewListFunctionsInput(name)
			if nextToken != "" {
				args.NextToken = &nextToken
			}

			raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				return fcClient.ListFunctions(args)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to list functions of service (%s (%s)): %s", name, id, err)
			}
			resp, _ := raw.(*fc.ListFunctionsOutput)

			if resp.Functions == nil || len(resp.Functions) < 1 {
				break
			}

			for _, function := range resp.Functions {
				// Remove triggers
				triggerDeleted := false
				triggerNext := ""
				for {
					req := fc.NewListTriggersInput(name, *function.FunctionName)
					if triggerNext != "" {
						req.NextToken = &triggerNext
					}

					raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
						return fcClient.ListTriggers(req)
					})
					if err != nil {
						log.Printf("[ERROR] Failed to list triggers of functiion (%s): %s", name, err)
					}
					resp, _ := raw.(*fc.ListTriggersOutput)

					if resp == nil || resp.Triggers == nil || len(resp.Triggers) < 1 {
						break
					}
					for _, trigger := range resp.Triggers {
						triggerDeleted = true
						if _, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
							return fcClient.DeleteTrigger(&fc.DeleteTriggerInput{
								ServiceName:  StringPointer(name),
								FunctionName: function.FunctionName,
								TriggerName:  trigger.TriggerName,
							})
						}); err != nil {
							log.Printf("[ERROR] Failed to delete trigger %s of function: %s.", *trigger.TriggerName, *function.FunctionName)
						}
					}

					triggerNext = ""
					if resp.NextToken != nil {
						triggerNext = *resp.NextToken
					}
					if triggerNext == "" {
						break
					}
				}
				//remove function
				if triggerDeleted {
					time.Sleep(5 * time.Second)
				}
				if _, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
					return fcClient.DeleteFunction(&fc.DeleteFunctionInput{
						ServiceName:  StringPointer(name),
						FunctionName: function.FunctionName,
					})
				}); err != nil {
					log.Printf("[ERROR] Failed to delete function %s of services: %s (%s)", *function.FunctionName, name, id)
				}
			}

			nextToken = ""
			if resp.NextToken != nil {
				nextToken = *resp.NextToken
			}
			if nextToken == "" {
				break
			}
		}

		log.Printf("[INFO] Deleting FC services: %s (%s)", name, id)
		_, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.DeleteService(&fc.DeleteServiceInput{
				ServiceName: StringPointer(name),
			})
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete FC services (%s (%s)): %s", name, id, err)
		}
	}

	// Delete FC custom domains.
	nextToken := ""
	for {
		raw, err = client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListCustomDomains(fc.NewListCustomDomainsInput().WithNextToken(nextToken))
		})
		if err != nil {
			return fmt.Errorf("Error retrieving FC custom domains: %s", err)
		}
		response := raw.(*fc.ListCustomDomainsOutput)
		nextToken = *response.NextToken
		for _, domain := range response.CustomDomains {
			log.Printf("[INFO] Deleting FC custom domain: %s", *domain.DomainName)
			_, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				return fcClient.DeleteCustomDomain(fc.NewDeleteCustomDomainInput(*domain.DomainName))
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete FC custom domains(%s): %s", *domain.DomainName, err)
			}
		}
		if nextToken == "" {
			break
		}
	}

	return nil
}

func TestAccAlicloudFCCustomDomainUpdate(t *testing.T) {
	var v *fc.GetCustomDomainOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc-%s-alicloudfccustomdomain-%d-cd", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":               CHECKSET,
		"created_time":       CHECKSET,
		"last_modified_time": CHECKSET,
	}
	resourceId := "alicloud_fc_custom_domain.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcCustomDomainConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     "terraform.functioncompute.com",
					"protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTP",
					"route_config": []map[string]interface{}{
						{
							"path":          "/*",
							"service_name":  "${alicloud_fc_service.default.name}",
							"function_name": "${alicloud_fc_function.default.name}",
							"qualifier":     "?query",
							"methods":       []string{"GET", "POST"},
						},
						{
							"path":          "/test",
							"service_name":  "${alicloud_fc_service.default.name}",
							"function_name": "${alicloud_fc_function.default.name}",
							"qualifier":     "?region",
							"methods":       []string{"HEAD", "PATCH"},
						},
					},
					/*"cert_config": []map[string]interface{}{
						{
							"name":        "fake",
							"private_key": "${var.private_key}",
							"certificate": "${var.certificate}",
						},
					},*/
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":                     "HTTP",
						"route_config.0.path":          "/*",
						"route_config.0.service_name":  name,
						"route_config.0.function_name": name,
						"route_config.0.qualifier":     "?query",
						"route_config.0.methods.0":     "GET",
						"route_config.0.methods.1":     "POST",
						"route_config.1.path":          "/test",
						"route_config.1.qualifier":     "?region",
						"route_config.1.methods.0":     "HEAD",
						"route_config.1.methods.1":     "PATCH",
					}),
				),
			},
		},
	})
}

func resourceFcCustomDomainConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
resource "alicloud_fc_service" "default" {
  name = "${var.name}"
  internet_access = false
}
resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}
resource "alicloud_oss_bucket_object" "default" {
  bucket = "${alicloud_oss_bucket.default.id}"
  key = "fc/hello.zip"
  content = <<EOF
    # -*- coding: utf-8 -*-
  def handler(event, context):
      print "hello world"
      return 'hello world'
  EOF
}
resource "alicloud_fc_function" "default" {
  service = "${alicloud_fc_service.default.name}"
  name = "${var.name}"
  oss_bucket = "${alicloud_oss_bucket.default.id}"
  oss_key = "${alicloud_oss_bucket_object.default.key}"
  memory_size = 512
  runtime = "python2.7"
  handler = "hello.handler"
}
`, name)
}

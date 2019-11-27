package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_fc_service", &resource.Sweeper{
		Name: "alicloud_fc_service",
		F:    testSweepFCServices,
		Dependencies: []string{
			"alicloud_fc_function",
		},
	})
}

func testSweepFCServices(region string) error {
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
	return nil
}

func TestAccAlicloudFCServiceUpdate(t *testing.T) {
	var v *fc.GetServiceOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%salicloudfcservice-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":          name,
		"last_modified": CHECKSET,
	}
	resourceId := "alicloud_fc_service.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcServiceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
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
					"description": "tf unit test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_access": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role":       "${alicloud_ram_role.default.arn}",
					"depends_on": []string{"alicloud_ram_role_policy_attachment.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_config": []map[string]string{
						{
							"project":  "${alicloud_log_store.default.project}",
							"logstore": "${alicloud_log_store.default.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_config.0.project":  name,
						"log_config.0.logstore": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_config":      REMOVEKEY,
					"role":            REMOVEKEY,
					"internet_access": REMOVEKEY,
					"description":     REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_config.0.project":  REMOVEKEY,
						"log_config.0.logstore": REMOVEKEY,
						"role":                  REMOVEKEY,
						"internet_access":       REMOVEKEY,
						"description":           REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudFCServiceVpcUpdate(t *testing.T) {
	var v *fc.GetServiceOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%salicloudfcservice-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":          name,
		"role":          CHECKSET,
		"vpc_config.#":  "1",
		"last_modified": CHECKSET,
	}
	resourceId := "alicloud_fc_service.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcServiceConfigVpcDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"role": "${alicloud_ram_role.default.arn}",
					"vpc_config": []map[string]interface{}{
						{
							"vswitch_ids":       "${alicloud_vswitch.default.*.id}",
							"security_group_id": "${alicloud_security_group.default.id}",
						},
					},
					"depends_on": []string{"alicloud_ram_role_policy_attachment.default", "alicloud_ram_role_policy_attachment.default1"},
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
					"description": "tf unit test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_access": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_config": []map[string]string{
						{
							"project":  "${alicloud_log_store.default.project}",
							"logstore": "${alicloud_log_store.default.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_config.0.project":  name,
						"log_config.0.logstore": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_config":      REMOVEKEY,
					"internet_access": REMOVEKEY,
					"description":     REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_config.0.project":  REMOVEKEY,
						"log_config.0.logstore": REMOVEKEY,
						"internet_access":       REMOVEKEY,
						"description":           REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudFCServiceMulti(t *testing.T) {
	var v *fc.GetServiceOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%salicloudfcservice-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":          name + "_9",
		"role":          CHECKSET,
		"last_modified": CHECKSET,
	}
	resourceId := "alicloud_fc_service.default.9"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcServiceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":      "10",
					"name":       "${var.name}_${count.index}",
					"role":       "${alicloud_ram_role.default.arn}",
					"depends_on": []string{"alicloud_ram_role_policy_attachment.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceFcServiceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    name = "${var.name}"
}

resource "alicloud_log_store" "default" {
    project = "${alicloud_log_project.default.name}"
    name = "${var.name}"
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<DEFINITION
  %s
  DEFINITION
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}`, name, testFCRoleTemplate)
}

func resourceFcServiceConfigVpcDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    name = "${var.name}"
}

resource "alicloud_log_store" "default" {
    project = "${alicloud_log_project.default.name}"
    name = "${var.name}"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_zones" "default" {
    available_resource_creation = "FunctionCompute"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
  name = "${var.name}"
  document = <<DEFINITION
  %s
  DEFINITION
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_ram_policy" "default" {
  name = "${var.name}"
  document = <<DEFINITION
  %s
  DEFINITION
}
resource "alicloud_ram_role_policy_attachment" "default1" {
  role_name = "${alicloud_ram_role.default.name}"
  policy_name = "${alicloud_ram_policy.default.name}"
  policy_type = "Custom"
}`, name, testFCRoleTemplate, testFCVpcPolicyTemplate)
}

var testFCRoleTemplate = `
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "fc.aliyuncs.com"
        ]
      }
    }
  ],
  "Version": "1"
}
`

var testFCVpcPolicyTemplate = `
{
  "Version": "1",
  "Statement": [
    {
      "Action": "vpc:*",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action": [
        "ecs:*NetworkInterface*"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
`

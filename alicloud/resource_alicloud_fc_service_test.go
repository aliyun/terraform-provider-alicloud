package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		log.Printf("Error retrieving FC services: %s", err)
		return nil
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

		// Remove eni
		log.Printf("[INFO] Prepare to delete eni which FC created...")
		if *v.VPCConfig.VPCID != "" || len(v.VPCConfig.VSwitchIDs) > 0 {
			action := "DescribeNetworkInterfaces"
			request := make(map[string]interface{})
			request["VpcId"] = *v.VPCConfig.VPCID
			request["VSwitchId"] = v.VPCConfig.VSwitchIDs[0]
			request["RegionId"] = client.RegionId
			request["PageSize"] = PageSizeLarge
			request["PageNumber"] = 1
			conn, err := client.NewEcsClient()
			if err != nil {
				return WrapError(err)
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				return WrapError(err)
			}
			addDebug(action, response, request)
			resp, err := jsonpath.Get("$.NetworkInterfaceSets.NetworkInterfaceSet", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkInterfaceSets.NetworkInterfaceSet", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				if fmt.Sprint(item["NetworkInterfaceName"]) != "fc-eni" {
					continue
				}
				log.Printf("[INFO] Deleting FC eni: %s (%s)", item["NetworkInterfaceName"], item["NetworkInterfaceId"])
				action := "DeleteNetworkInterface"
				request := make(map[string]interface{})
				request["RegionId"] = client.RegionId
				request["NetworkInterfaceId"] = fmt.Sprint(item["NetworkInterfaceId"])
				_, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					return WrapError(err)
				}
			}
		}
		// Delete the service versions.
		log.Printf("[INFO] Prepare to delete FC versions...")
		input := &fc.ListServiceVersionsInput{
			ServiceName: v.ServiceName,
			Limit:       Int32Pointer(100),
		}
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListServiceVersions(input)
		})
		if err != nil {
			return WrapError(err)
		}

		output := raw.(*fc.ListServiceVersionsOutput)
		for _, v := range output.Versions {
			log.Printf("[INFO] Deleting FC service %s version: %s", *input.ServiceName, *v.VersionID)
			input := &fc.DeleteServiceVersionInput{
				ServiceName: input.ServiceName,
				VersionID:   v.VersionID,
			}
			_, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				return fcClient.DeleteServiceVersion(input)
			})
			if err != nil {
				return WrapError(err)
			}
		}
		log.Printf("[INFO] Deleting FC services: %s (%s)", name, id)
		_, err = client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
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
					"name":    "${var.name}",
					"publish": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "publish"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf unit test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf unit test",
						"version":     "2",
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
						"version":         "3",
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
						"role":    CHECKSET,
						"version": "4",
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
						"version":               "5",
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
						"version":               "6",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudFCServiceVpcAndNasUpdate(t *testing.T) {
	var v *fc.GetServiceOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testacc%salicloudfcservice-%d", defaultRegionToTest, rand)
	var basicMap = map[string]string{
		"name":          name,
		"role":          CHECKSET,
		"vpc_config.#":  "1",
		"nas_config.#":  "1",
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
							"vswitch_ids":       []string{"${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"},
							"security_group_id": "${alicloud_security_group.default.id}",
						},
					},
					"nas_config": []map[string]interface{}{
						{
							"user_id":  "9527",
							"group_id": "9528",
							"mount_points": []map[string]interface{}{
								{
									"server_addr": "${local.mount_target_domain}",
									"mount_dir":   "/mnt/nas",
								},
							},
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
				ImportStateVerifyIgnore: []string{"name_prefix", "publish"},
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
							"project":  "${alicloud_log_project.default.name}",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"nas_config": []map[string]interface{}{
						{
							"user_id":  "9527",
							"group_id": "9528",
							"mount_points": []map[string]interface{}{
								{
									"server_addr": "${local.mount_target_domain1}",
									"mount_dir":   "/mnt/nas",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nas_config.0.user_id":                    "9527",
						"nas_config.0.group_id":                   "9528",
						"nas_config.0.mount_points.0.server_addr": CHECKSET,
						"nas_config.0.mount_points.0.mount_dir":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nas_config": []map[string]interface{}{
						{
							"user_id":  "9627",
							"group_id": "9628",
							"mount_points": []map[string]interface{}{
								{
									"server_addr": "${local.mount_target_domain}",
									"mount_dir":   "/mnt/nas",
								},
								{
									"server_addr": "${local.mount_target_domain1}",
									"mount_dir":   "/home/nas",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nas_config.0.user_id":                    "9627",
						"nas_config.0.group_id":                   "9628",
						"nas_config.0.mount_points.0.server_addr": CHECKSET,
						"nas_config.0.mount_points.0.mount_dir":   CHECKSET,
						"nas_config.0.mount_points.1.server_addr": CHECKSET,
						"nas_config.0.mount_points.1.mount_dir":   CHECKSET,
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

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}

resource "alicloud_nas_file_system" "this" {
  protocol_type = "NFS"
  storage_type = "Performance"
}

resource "alicloud_nas_access_group" "this" {
  access_group_name = "${var.name}"
  access_group_type = "Vpc"
}

resource "alicloud_nas_mount_target" "this" {
  count = 2
  access_group_name = alicloud_nas_access_group.this.access_group_name
  file_system_id = alicloud_nas_file_system.this.id
  vswitch_id = data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0
}

locals {
  mount_target_domain = format("%%s://mnt",split(":",alicloud_nas_mount_target.this[0].id)[1])
  mount_target_domain1 = format("%%s://mnt",split(":",alicloud_nas_mount_target.this[1].id)[1])
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

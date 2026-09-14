package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			response, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
				_, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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

func TestAccAliCloudFCServiceUpdate(t *testing.T) {
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
					"log_config": []map[string]string{
						{
							"project":                 "${alicloud_log_store.default.project}",
							"logstore":                "${alicloud_log_store.default.name}",
							"enable_request_metrics":  "true",
							"enable_instance_metrics": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_config.0.project":                 name,
						"log_config.0.logstore":                name,
						"log_config.0.enable_request_metrics":  CHECKSET,
						"log_config.0.enable_instance_metrics": CHECKSET,
						"version":                              "6",
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
						"log_config.0.project":                 REMOVEKEY,
						"log_config.0.logstore":                REMOVEKEY,
						"log_config.0.enable_request_metrics":  REMOVEKEY,
						"log_config.0.enable_instance_metrics": REMOVEKEY,
						"role":                                 REMOVEKEY,
						"internet_access":                      REMOVEKEY,
						"description":                          REMOVEKEY,
						"version":                              "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"TestKey1": "test-value1",
						"TestKey2": "test-value2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.TestKey1": "test-value1",
						"tags.TestKey2": "test-value2",
						"version":       "8",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudFCServiceVpcAndNasUpdate(t *testing.T) {
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
							"vswitch_ids":       []string{"${alicloud_vswitch.default.id}"},
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
							"project":                 "${alicloud_log_project.default.name}",
							"logstore":                "${alicloud_log_store.default.name}",
							"enable_request_metrics":  "true",
							"enable_instance_metrics": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_config.0.project":                 name,
						"log_config.0.logstore":                name,
						"log_config.0.enable_request_metrics":  CHECKSET,
						"log_config.0.enable_instance_metrics": CHECKSET,
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
						"log_config.0.project":                 REMOVEKEY,
						"log_config.0.logstore":                REMOVEKEY,
						"log_config.0.enable_request_metrics":  REMOVEKEY,
						"log_config.0.enable_instance_metrics": REMOVEKEY,
						"internet_access":                      REMOVEKEY,
						"description":                          REMOVEKEY,
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

func TestAccAliCloudFCServiceMulti(t *testing.T) {
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

data "alicloud_fc_zones" "default" {}

resource "alicloud_log_project" "default" {
    name = "${var.name}"
}

resource "alicloud_log_store" "default" {
    project = "${alicloud_log_project.default.name}"
    name = "${var.name}"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/21"
  zone_id      = "${data.alicloud_fc_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
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
  vswitch_id = alicloud_vswitch.default.id
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

// TestUnitAliCloudFCServiceParseVpcConfig guards parseVpcConfig against a regression where a
// vpc_config block with an empty vswitch_ids set caused an index-out-of-range panic. It also
// pins the documented behavior that a vpc_config block whose vswitch_ids and security_group_id
// are both empty is treated as unset (see website/docs/r/fc_service.html.markdown), including
// the requirement that this is decided before the 'role' guard. Every case returns before
// parseVpcConfig reaches DescribeVSwitch, so the test performs no network call and needs no
// credentials.
func TestUnitAliCloudFCServiceParseVpcConfig(t *testing.T) {
	meta := &connectivity.AliyunClient{}
	const testRole = "acs:ram::1234567890123456:role/fc-service-unit-test"

	// Case 1 (regression anchor): vpc_config present but vswitch_ids and security_group_id both
	// empty, with role set. On the unfixed code this panics at vswitchIds[0]; the fix returns
	// (nil, nil).
	t.Run("both empty is treated as unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlicloudFCService().Schema, map[string]interface{}{
			"role": testRole,
			"vpc_config": []interface{}{
				map[string]interface{}{
					"vswitch_ids":       []interface{}{},
					"security_group_id": "",
				},
			},
		})
		config, err := parseVpcConfig(d, meta)
		if err != nil {
			t.Fatalf("expected no error for an empty vpc_config block, got: %v", err)
		}
		if config != nil {
			t.Fatalf("expected nil VPCConfig for an empty vpc_config block, got: %#v", config)
		}
	})

	// Case 2: vpc_config entirely absent must also produce no VPCConfig.
	t.Run("missing vpc_config block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlicloudFCService().Schema, map[string]interface{}{
			"role": testRole,
		})
		config, err := parseVpcConfig(d, meta)
		if err != nil {
			t.Fatalf("expected no error when vpc_config is absent, got: %v", err)
		}
		if config != nil {
			t.Fatalf("expected nil VPCConfig when vpc_config is absent, got: %#v", config)
		}
	})

	// Case 3: an empty vswitch_ids with a non-empty security_group_id cannot build a valid
	// VPCConfig, so parseVpcConfig must return a clear error instead of panicking.
	t.Run("empty vswitch_ids with security_group_id errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlicloudFCService().Schema, map[string]interface{}{
			"role": testRole,
			"vpc_config": []interface{}{
				map[string]interface{}{
					"vswitch_ids":       []interface{}{},
					"security_group_id": "test-security-group",
				},
			},
		})
		config, err := parseVpcConfig(d, meta)
		if err == nil {
			t.Fatalf("expected an error when vswitch_ids is empty but security_group_id is set, got config: %#v", config)
		}
	})

	// Case 4 (ordering anchor): an empty vpc_config block must be treated as unset before the
	// 'role' requirement is enforced, so omitting role must not surface the "'role' is required"
	// error.
	t.Run("empty vpc_config does not require role", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlicloudFCService().Schema, map[string]interface{}{
			"vpc_config": []interface{}{
				map[string]interface{}{
					"vswitch_ids":       []interface{}{},
					"security_group_id": "",
				},
			},
		})
		config, err := parseVpcConfig(d, meta)
		if err != nil {
			t.Fatalf("expected no error (role not required for an empty vpc_config), got: %v", err)
		}
		if config != nil {
			t.Fatalf("expected nil VPCConfig for an empty vpc_config block, got: %#v", config)
		}
	})
}

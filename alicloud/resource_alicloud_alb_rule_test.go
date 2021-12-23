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
	resource.AddTestSweepers(
		"alicloud_alb_rule",
		&resource.Sweeper{
			Name: "alicloud_alb_rule",
			F:    testSweepAlbRule,
		})
}

func testSweepAlbRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListRules"
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

		resp, err := jsonpath.Get("$.Rules", response)

		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Rules", action, err)
			return nil
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["RuleName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ALB Rule: %s", item["RuleName"].(string))
				continue
			}

			sweeped = true
			action := "DeleteRule"
			request := map[string]interface{}{
				"RuleId": item["RuleId"],
			}
			request["ClientToken"] = buildClientToken("DeleteRule")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ALB Rule (%s): %s", item["RuleName"].(string), err)
			}
			if sweeped {
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete ALB Rule success: %s ", item["RuleName"].(string))
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}

func TestAccAlicloudALBRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudALBRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBRuleBasicDependence0)
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
					"listener_id": "${alicloud_alb_listener.default.id}",
					"rule_name":   "${var.name}",
					"priority":    "666",
					"rule_conditions": []map[string]interface{}{
						{
							"cookie_config": []map[string]interface{}{
								{
									"values": []map[string]interface{}{
										{
											"key":   "created",
											"value": "tf",
										},
									},
								},
							},
							"type": "Cookie",
						},
					},
					"rule_actions": []map[string]interface{}{
						{
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${alicloud_alb_server_group.default.id}",
										},
									},
								},
							},
							"order": "9",
							"type":  "ForwardGroup",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_id":       CHECKSET,
						"rule_name":         name,
						"priority":          "666",
						"rule_actions.#":    "1",
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"redirect_config": []map[string]interface{}{
								{
									"host":      "ww.ali.com",
									"http_code": "301",
									"path":      "/test",
									"port":      "10",
									"protocol":  "HTTP",
									"query":     "query",
								},
							},
							"order": "4",
							"type":  "Redirect",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"fixed_response_config": []map[string]interface{}{
								{
									"content":      "tf-testAcc",
									"content_type": "application/json",
									"http_code":    "200",
								},
							},
							"order": "2",
							"type":  "FixedResponse",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"insert_header_config": []map[string]interface{}{
								{
									"key":        "tf-insert-header",
									"value":      "SLBId",
									"value_type": "SystemDefined",
								},
							},
							"order": "3",
							"type":  "InsertHeader",
						},
						{
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${alicloud_alb_server_group.default.id}",
										},
									},
								},
							},
							"order": "9",
							"type":  "ForwardGroup",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_actions": []map[string]interface{}{
						{
							"rewrite_config": []map[string]interface{}{
								{
									"host":  "www.test.com",
									"path":  "/test",
									"query": "test",
								},
							},
							"order": "5",
							"type":  "Rewrite",
						},
						{
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${alicloud_alb_server_group.default.id}",
										},
									},
								},
							},
							"order": "9",
							"type":  "ForwardGroup",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_actions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"header_config": []map[string]interface{}{
								{
									"key":    "Port",
									"values": []string{"5006"},
								},
							},
							"type": "Header",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"host_config": []map[string]interface{}{
								{
									"values": []string{"www.test.com"},
								},
							},
							"type": "Host",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"method_config": []map[string]interface{}{
								{
									"values": []string{"PUT"},
								},
							},
							"type": "Method",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"path_config": []map[string]interface{}{
								{
									"values": []string{"/test"},
								},
							},
							"type": "Path",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_conditions": []map[string]interface{}{
						{
							"query_string_config": []map[string]interface{}{
								{
									"values": []map[string]interface{}{
										{
											"key":   "test",
											"value": "test",
										},
									},
								},
							},
							"type": "QueryString",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name,
					"priority":  "777",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name,
						"priority":  "777",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBRuleMap0 = map[string]string{}

func AlicloudALBRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 

variable "name" {	
	default = "%s"
}

data "alicloud_alb_zones" "default"{}

data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count             = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id =  data.alicloud_alb_zones.default.zones.0.id
  vswitch_name              = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count             = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name              = var.name
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id =              data.alicloud_vpcs.default.ids.0
  address_type =        "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name =    var.name
  load_balancer_edition = "Standard"
  load_balancer_billing_config {
    pay_type = 	"PayAsYouGo"
  }
  zone_mappings{
		vswitch_id =  length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
		zone_id =  data.alicloud_alb_zones.default.zones.0.id
	}
  zone_mappings{
		vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
		zone_id =   data.alicloud_alb_zones.default.zones.1.id
	}
}

resource "alicloud_alb_server_group" "default" {
	protocol = "HTTP"
	vpc_id = data.alicloud_vpcs.default.vpcs.0.id
	server_group_name = var.name
	health_check_config {
       health_check_enabled = "false"
	}
	sticky_session_config {
       sticky_session_enabled = "false"
	}
}

resource "alicloud_alb_listener" "default" {
	load_balancer_id = alicloud_alb_load_balancer.default.id
	listener_protocol =  "HTTP"
	listener_port = 8080
	listener_description = var.name
	default_actions{
		type = "ForwardGroup"
		forward_group_config{
			server_group_tuples{
				server_group_id = alicloud_alb_server_group.default.id
			}
		}
	}
}

`, name)
}

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
		"alicloud_alb_server_group",
		&resource.Sweeper{
			Name: "alicloud_alb_server_group",
			F:    testSweepAlbServerGroup,
		})
}

func testSweepAlbServerGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListServerGroups"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeXLarge

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

		resp, err := jsonpath.Get("$.ServerGroups", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ServerGroups", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["ServerGroupName"]; !ok {
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["ServerGroupName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Alb Server Group: %s", item["ServerGroupName"].(string))
					continue
				}
			}
			action := "DeleteServerGroup"
			request := map[string]interface{}{
				"ServerGroupId": item["ServerGroupId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alb Server Group (%s): %s", item["ServerGroupName"].(string), err)
			}
			log.Printf("[INFO] Delete Alb Server Group success: %s ", item["ServerGroupName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudALBServerGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBServerGroupBasicDependence0)
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
					"server_group_name": name,
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "true",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "true",
						},
					},
					"slow_start_config": []map[string]interface{}{
						{
							"slow_start_enabled":  "true",
							"slow_start_duration": "30",
						},
					},
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":       name,
						"vpc_id":                  CHECKSET,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
						"cross_zone_enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "false",
					"slow_start_config": []map[string]interface{}{
						{
							"slow_start_enabled":  "true",
							"slow_start_duration": "40",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"slow_start_config": []map[string]interface{}{
						{
							"slow_start_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "300",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "false",
							"connection_drain_timeout": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upstream_keepalive_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upstream_keepalive_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upstream_keepalive_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upstream_keepalive_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wlc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wlc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Uch",
					"uch_config": []map[string]interface{}{
						{
							"type":  "QueryString",
							"value": "abc",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Uch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_template_id": "${alicloud_alb_health_check_template.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_template_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Server",
							"cookie":                 name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled":      "true",
							"health_check_connect_port": "46325",
							"health_check_host":         "tf-testAcc.com",
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/tf-testAcc",
							"health_check_protocol":     "HTTP",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
							"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_id":         "${alicloud_instance.default.id}",
							"server_type":       "Ecs",
							"server_ip":         "${alicloud_instance.default.private_ip}",
							"port":              "80",
							"remote_ip_enabled": "false",
							"weight":            "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_id":         "${alicloud_instance.default.id}",
							"server_type":       "Ecs",
							"server_ip":         "${alicloud_instance.default.private_ip}",
							"port":              "80",
							"remote_ip_enabled": "false",
							"weight":            "10",
							"description":       name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ServerGroup",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"health_check_template_id"},
			},
		},
	})
}

func TestAccAliCloudALBServerGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBServerGroupBasicDependence0)
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
					"server_group_name": name,
					"server_group_type": "Instance",
					"protocol":          "HTTP",
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"scheduler":         "Wlc",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Server",
							"cookie":                 name,
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled":      "true",
							"health_check_connect_port": "46325",
							"health_check_host":         "tf-testAcc.com",
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/tf-testAcc",
							"health_check_protocol":     "HTTP",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
							"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
					"servers": []map[string]interface{}{
						{
							"server_id":         "${alicloud_instance.default.id}",
							"server_type":       "Ecs",
							"server_ip":         "${alicloud_instance.default.private_ip}",
							"port":              "80",
							"remote_ip_enabled": "false",
							"weight":            "10",
							"description":       name,
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":       name,
						"server_group_type":       "Instance",
						"protocol":                "HTTP",
						"vpc_id":                  CHECKSET,
						"scheduler":               "Wlc",
						"resource_group_id":       CHECKSET,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
						"servers.#":               "1",
						"tags.%":                  "2",
						"tags.Created":            "TF",
						"tags.For":                "ServerGroup",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudALBServerGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBServerGroupBasicDependence0)
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
					"server_group_name": name,
					"server_group_type": "Ip",
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "false",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":       name,
						"server_group_type":       "Ip",
						"vpc_id":                  CHECKSET,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "Wlc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": "Wlc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Insert",
							"cookie_timeout":         "1000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled":      "true",
							"health_check_connect_port": "46325",
							"health_check_host":         "tf-testAcc.com",
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/tf-testAcc",
							"health_check_protocol":     "HTTP",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
							"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_id":         "${cidrhost(data.alicloud_vpcs.default.vpcs.0.cidr_block, 2)}",
							"server_type":       "Ip",
							"server_ip":         "${cidrhost(data.alicloud_vpcs.default.vpcs.0.cidr_block, 2)}",
							"port":              "80",
							"remote_ip_enabled": "false",
							"weight":            "10",
							"description":       name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ServerGroup",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudALBServerGroup_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBServerGroupBasicDependence0)
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
					"server_group_name": name,
					"server_group_type": "Ip",
					"protocol":          "HTTP",
					"vpc_id":            "${data.alicloud_vpcs.default.vpcs.0.id}",
					"scheduler":         "Wlc",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"sticky_session_config": []map[string]interface{}{
						{
							"sticky_session_enabled": "true",
							"sticky_session_type":    "Insert",
							"cookie_timeout":         "1000",
						},
					},
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled":      "true",
							"health_check_connect_port": "46325",
							"health_check_host":         "tf-testAcc.com",
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/tf-testAcc",
							"health_check_protocol":     "HTTP",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
							"health_check_codes":        []string{"http_2xx", "http_3xx", "http_4xx"},
						},
					},
					"servers": []map[string]interface{}{
						{
							"server_id":         "${cidrhost(data.alicloud_vpcs.default.vpcs.0.cidr_block, 2)}",
							"server_type":       "Ip",
							"server_ip":         "${cidrhost(data.alicloud_vpcs.default.vpcs.0.cidr_block, 2)}",
							"port":              "80",
							"remote_ip_enabled": "false",
							"weight":            "10",
							"description":       name,
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":       name,
						"server_group_type":       "Ip",
						"protocol":                "HTTP",
						"vpc_id":                  CHECKSET,
						"scheduler":               "Wlc",
						"resource_group_id":       CHECKSET,
						"sticky_session_config.#": "1",
						"health_check_config.#":   "1",
						"servers.#":               "1",
						"tags.%":                  "2",
						"tags.Created":            "TF",
						"tags.For":                "ServerGroup",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudALBServerGroup_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBServerGroupBasicDependence0)
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
					"server_group_name": name,
					"server_group_type": "Fc",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":     name,
						"server_group_type":     "Fc",
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled":      "true",
							"health_check_connect_port": "80",
							"health_check_host":         "${data.alicloud_account.default.id}." + defaultRegionToTest + ".fc.aliyuncs.com",
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ServerGroup",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudALBServerGroup_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudALBServerGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBServerGroupBasicDependence0)
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
					"server_group_name": name,
					"server_group_type": "Fc",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled":      "true",
							"health_check_connect_port": "80",
							"health_check_host":         "${data.alicloud_account.default.id}." + defaultRegionToTest + ".fc.aliyuncs.com",
							"health_check_http_version": "HTTP1.1",
							"health_check_interval":     "2",
							"health_check_method":       "HEAD",
							"health_check_path":         "/",
							"health_check_timeout":      "5",
							"healthy_threshold":         "3",
							"unhealthy_threshold":       "3",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_name":     name,
						"server_group_type":     "Fc",
						"resource_group_id":     CHECKSET,
						"health_check_config.#": "1",
						"tags.%":                "2",
						"tags.Created":          "TF",
						"tags.For":              "ServerGroup",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudALBServerGroupMap0 = map[string]string{
	"server_group_type": CHECKSET,
	"scheduler":         CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
}

func AliCloudALBServerGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		instance_name              = var.name
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_alb_health_check_template" "default" {
	  health_check_template_name = var.name
	}
`, name)
}

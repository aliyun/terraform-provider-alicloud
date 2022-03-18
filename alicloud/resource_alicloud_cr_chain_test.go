package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_chain", &resource.Sweeper{
		Name: "alicloud_cr_chain",
		F:    testSweepCRChain,
	})
}

func testSweepCRChain(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftestAcc",
	}
	conn, err := client.NewAcrClient()
	if err != nil {
		log.Printf("[ERROR] Failed to fetch client.NewAcrClient: %s ", err)
		return nil
	}

	crService := &CrService{client}
	pageNo := 1
	pageSize := 50

	var instances []cr_ee.InstancesItem
	for {
		resp, err := crService.ListCrEEInstances(pageNo, pageSize)
		if err != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s ", err)
			return nil
		}
		instances = append(instances, resp.Instances...)
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}

	instanceIds := make([]string, 0)
	for _, instance := range instances {
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(instance.InstanceName), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping cr ee instance: %s ", instance.InstanceName)
			continue
		}
		instanceIds = append(instanceIds, instance.InstanceId)
	}

	for _, instanceId := range instanceIds {
		action := "ListChain"
		request := map[string]interface{}{
			"InstanceId": instanceId,
		}
		var response map[string]interface{}
		chainIds := make([]string, 0)

		pageNo, pageSize := 1, PageSizeLarge
		for {
			request["PageNo"] = pageNo
			request["PageSize"] = pageSize
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Failed To List CR Chains : %s", err)
			}
			resp, err := jsonpath.Get("$.Chains", response)
			if err != nil {
				log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", action, "$.Chains", response)
				return nil
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				skip := true
				item := v.(map[string]interface{})
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["Name"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping CR Chain: %v (%v)", item["Name"], item["ChainId"])
					continue
				}
				chainIds = append(chainIds, fmt.Sprint(item["InstanceId"], ":", item["ChainId"]))
			}
			if len(result) < pageSize {
				break
			}
			pageNo++
		}

		if len(chainIds) > 0 {
			log.Printf("[INFO] Deleting CR Chains: (%s)", strings.Join(chainIds, ","))
			action = "DeleteChain"

			for _, chainId := range chainIds {
				parts, err := ParseResourceId(chainId, 2)
				if err != nil {
					log.Printf("[ERROR] Failed to parse ResourceId %s, %s", parts, err)
					return nil
				}
				deleteRequest := map[string]interface{}{
					"ChainId":    parts[1],
					"InstanceId": parts[0],
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(time.Minute*10, func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-01"), StringPointer("AK"), nil, deleteRequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Failed To Delete CR Chain (%s): %v", chainId, err)
					continue
				}
			}
		}
	}
	return nil
}

func TestAccAlicloudCRChain_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_chain.default"
	checkoutSupportedRegions(t, true, connectivity.CRSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCRChainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrChain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacccrchain%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCRChainBasicDependence0)
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
					"instance_id":         "${alicloud_cr_ee_namespace.default.instance_id}",
					"chain_name":          "tf-testacc-1",
					"repo_namespace_name": "${alicloud_cr_ee_namespace.default.name}",
					"repo_name":           "${alicloud_cr_ee_repo.default.name}",
					"chain_config": []map[string]interface{}{
						{
							"routers": []map[string]interface{}{
								{
									"from": []map[string]interface{}{
										{
											"node_name": "DOCKER_IMAGE_BUILD",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "DOCKER_IMAGE_PUSH",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "DOCKER_IMAGE_PUSH",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "VULNERABILITY_SCANNING",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "VULNERABILITY_SCANNING",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "ACTIVATE_REPLICATION",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "ACTIVATE_REPLICATION",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "TRIGGER",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "VULNERABILITY_SCANNING",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "SNAPSHOT",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "SNAPSHOT",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "TRIGGER_SNAPSHOT",
										},
									},
								},
							},
							"nodes": []map[string]interface{}{
								{
									"enable":    "true",
									"node_name": "DOCKER_IMAGE_BUILD",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "DOCKER_IMAGE_PUSH",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "VULNERABILITY_SCANNING",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{
													"issue_level": "MEDIUM",
													"issue_count": "1",
													"action":      "BLOCK_DELETE_TAG",
													"logic":       "AND",
												},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "ACTIVATE_REPLICATION",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "TRIGGER",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "false",
									"node_name": "SNAPSHOT",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "false",
									"node_name": "TRIGGER_SNAPSHOT",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"chain_name":  "tf-testacc-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"chain_name": "tf-testacc-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"chain_name": "tf-testacc-2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "describe",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "describe",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"chain_config": []map[string]interface{}{
						{
							"routers": []map[string]interface{}{
								{
									"from": []map[string]interface{}{
										{
											"node_name": "DOCKER_IMAGE_PUSH",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "VULNERABILITY_SCANNING",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "DOCKER_IMAGE_BUILD",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "DOCKER_IMAGE_PUSH",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "VULNERABILITY_SCANNING",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "ACTIVATE_REPLICATION",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "ACTIVATE_REPLICATION",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "TRIGGER",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "VULNERABILITY_SCANNING",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "SNAPSHOT",
										},
									},
								},
								{
									"from": []map[string]interface{}{
										{
											"node_name": "SNAPSHOT",
										},
									},
									"to": []map[string]interface{}{
										{
											"node_name": "TRIGGER_SNAPSHOT",
										},
									},
								},
							},
							"nodes": []map[string]interface{}{
								{
									"enable":    "true",
									"node_name": "DOCKER_IMAGE_PUSH",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "DOCKER_IMAGE_BUILD",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "VULNERABILITY_SCANNING",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{
													"issue_level": "MEDIUM",
													"issue_count": "1",
													"action":      "BLOCK_DELETE_TAG",
													"logic":       "AND",
												},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "ACTIVATE_REPLICATION",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "true",
									"node_name": "TRIGGER",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "false",
									"node_name": "SNAPSHOT",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
								{
									"enable":    "false",
									"node_name": "TRIGGER_SNAPSHOT",
									"node_config": []map[string]interface{}{
										{
											"deny_policy": []map[string]interface{}{
												{},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"repo_namespace_name", "repo_name"},
			},
		},
	})
}

var AlicloudCRChainMap0 = map[string]string{
	"chain_id":    CHECKSET,
	"instance_id": CHECKSET,
	"chain_name":  CHECKSET,
}

func AlicloudCRChainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_cr_ee_instances" "default" {
  name_regex  = "tf-testacc"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = data.alicloud_cr_ee_instances.default.ids[0]
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_namespace.default.instance_id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = var.name
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}
`, name)
}

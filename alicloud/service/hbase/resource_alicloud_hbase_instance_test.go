package hbase_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	hbaseSvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_hbase_instance", &resource.Sweeper{
		Name: "alicloud_hbase_instance",
		F:    testSweepHBaseInstances,
	})
}

func testSweepHBaseInstances(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []hbase.Instance
	req := hbase.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(helper.PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving HBase Instances: %s", err)
		}
		resp, _ := raw.(*hbase.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Instances.Instance...)

		if len(resp.Instances.Instance) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		if !sweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If an hbase instance name is set by other service, it should be fetched by vpc name and deleted.
			if skip {
				if need, err := needSweepHBaseVpc(client, v.VpcId); err == nil {
					skip = !need
				}

			}

			if skip {
				log.Printf("[INFO] Skipping Hbase Instance: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting HBase Instance: %s (%s)", name, id)
		req1 := hbase.CreateModifyClusterDeletionProtectionRequest()
		req1.ClusterId = id
		req1.Protection = requests.NewBoolean(false)
		_, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.ModifyClusterDeletionProtection(req1)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to modify DeletionProtection for Hbase Instance (%s (%s)): %s", name, id, err)
		}
		req2 := hbase.CreateDeleteInstanceRequest()
		req2.ClusterId = id
		_, err = client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DeleteInstance(req2)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete HBase Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

// needSweepHBaseVpc replicates the flat suite's VpcService.needSweepVpc behavior
// for the vpcId-only path used by the hbase instance sweeper: it sweeps the
// instance when its VPC was created by an acceptance test (name prefix
// tf-testacc / tf_testacc).
func needSweepHBaseVpc(client *connectivity.AliyunClient, vpcId string) (bool, error) {
	if vpcId == "" {
		return false, nil
	}
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"VpcId":    vpcId,
	}
	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			return false, nil
		}
		return false, errs.WrapErrorf(err, errs.DefaultErrorMsg, vpcId, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return false, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, vpcId, "$.Vpcs.Vpc", response)
	}
	items, _ := v.([]interface{})
	if len(items) < 1 {
		return false, nil
	}
	object, _ := items[0].(map[string]interface{})
	if fmt.Sprint(object["VpcId"]) != vpcId {
		return false, nil
	}
	name := strings.ToLower(fmt.Sprint(object["VpcName"]))
	if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
		log.Printf("[DEBUG] Need to sweep the VPC (%v (%v)).", object["VpcId"], object["VpcName"])
		return true, nil
	}
	return false, nil
}

func AlicloudHbaseBasicDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_hbase_zones" "default" {}
variable "name" {
	default = "%s"
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_hbase_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  # Test-infra adaptation: the flat-suite baseline used cidrsubnet offset 8,
  # which persistently overlaps an existing subnet in the shared
  # ^default-NODELETING$ VPC (InvalidCidrBlock.Overlapped) and blocks the
  # whole lifecycle ACC. Offset 200 picks a high /24 segment unlikely to be
  # allocated; hbase lifecycle assertions are unchanged.
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 200)
  zone_id = data.alicloud_hbase_zones.default.ids.0
  vswitch_name              = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_security_group" "default" {
	count = 2
	vpc_id = data.alicloud_vpcs.default.ids.0
	name = var.name
}
`, name)
}

func TestAccAliCloudHBaseInstanceVpc(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbase_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbaseSvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHBaseInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHbaseBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                   "${var.name}",
					"engine":                 "hbaseue",
					"engine_version":         "2.0",
					"master_instance_type":   "hbase.sn2.2xlarge",
					"core_instance_type":     "hbase.sn2.2xlarge",
					"core_disk_type":         "cloud_ssd",
					"vpc_id":                 "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":             "${local.vswitch_id}",
					"zone_id":                "${data.alicloud_hbase_zones.default.ids.0}",
					"immediate_delete_flag":  "true",
					"ip_white":               "192.168.0.1",
					"cold_storage_size":      "800",
					"auto_renew":             "0",
					"core_instance_quantity": "2",
					"pay_type":               "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                     name,
						"engine":                   "hbaseue",
						"engine_version":           "2.0",
						"core_instance_type":       "hbase.sn2.2xlarge",
						"core_disk_type":           "cloud_ssd",
						"vswitch_id":               acctest.CHECKSET,
						"immediate_delete_flag":    "true",
						"core_instance_quantity":   "2",
						"cold_storage_size":        "800",
						"deletion_protection":      "true",
						"zone_id":                  acctest.CHECKSET,
						"master_instance_quantity": acctest.CHECKSET,
						"maintain_start_time":      acctest.CHECKSET,
						"maintain_end_time":        acctest.CHECKSET,
						"pay_type":                 "PostPaid",
						"ip_white":                 "192.168.0.1",
						"security_groups.#":        "0",
						"ui_proxy_conn_addrs.#":    acctest.CHECKSET,
						"zk_conn_addrs.#":          acctest.CHECKSET,
						"slb_conn_addrs.#":         acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pay_type": "PrePaid",
					"duration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type": "PrePaid",
						"duration": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"immediate_delete_flag"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "04:00Z",
					"maintain_end_time":   "06:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"core_disk_size": "440",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"core_disk_size": "440",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_white": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white": "192.168.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account":  "admin",
					"password": "YourPassword@123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account":  "admin",
						"password": "YourPassword@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cold_storage_size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cold_storage_size": "800",
					}),
				),
			},
			//{
			//	Config: resourceHBaseConfigPrePaid,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"name":                "tf_testAccHBase_vpc_change_name",
			//			"maintain_start_time": "04:00Z",
			//			"maintain_end_time":   "06:00Z",
			//			"tags.%":              "2",
			//			"tags.Created":        "TF",
			//			"tags.For":            "acceptance test",
			//			"ip_white":            "192.168.1.1",
			//			"cold_storage_size":   "800",
			//			"pay_type":            "PrePaid",
			//		}),
			//	),
			//},

			{
				Config: testAccConfig(map[string]interface{}{
					"master_instance_type": "hbase.sn1.4xlarge",
					"core_instance_type":   "hbase.sn1.4xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_instance_type": "hbase.sn1.4xlarge",
						"core_instance_type":   "hbase.sn1.4xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"core_instance_quantity": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"core_instance_quantity": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pay_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${alicloud_security_group.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name,
					"maintain_start_time": "14:00Z",
					"maintain_end_time":   "16:00Z",
					"security_groups":     []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}"},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "acceptance test 123",
					},
					"core_disk_size":    "480",
					"ip_white":          "192.168.1.2",
					"account":           "adminu",
					"password":          "YourPassword@123u",
					"cold_storage_size": "900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                name,
						"maintain_start_time": "14:00Z",
						"maintain_end_time":   "16:00Z",
						"security_groups.#":   "2",
						"tags.%":              "2",
						"tags.Created":        "TF-update",
						"tags.For":            "acceptance test 123",
						"core_disk_size":      "480",
						"ip_white":            "192.168.1.2",
						"account":             "adminu",
						"password":            "YourPassword@123u",
						"cold_storage_size":   "900",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudHBaseInstance_VpcId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_hbase_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return hbaseSvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeHBaseInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudHbaseBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  "${var.name}",
					"engine":                "hbaseue",
					"engine_version":        "2.0",
					"master_instance_type":  "hbase.sn1.2xlarge",
					"core_instance_type":    "hbase.sn1.2xlarge",
					"core_disk_type":        "cloud_ssd",
					"vswitch_id":            "${local.vswitch_id}",
					"immediate_delete_flag": "true",
					"ip_white":              "192.168.0.1",
					"cold_storage_size":     "800",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                     name,
						"engine":                   "hbaseue",
						"engine_version":           "2.0",
						"core_instance_type":       "hbase.sn1.2xlarge",
						"core_disk_type":           "cloud_ssd",
						"vswitch_id":               acctest.CHECKSET,
						"immediate_delete_flag":    "true",
						"core_instance_quantity":   "2",
						"cold_storage_size":        "800",
						"deletion_protection":      "true",
						"zone_id":                  acctest.CHECKSET,
						"master_instance_quantity": acctest.CHECKSET,
						"maintain_start_time":      acctest.CHECKSET,
						"maintain_end_time":        acctest.CHECKSET,
						"pay_type":                 "PostPaid",
						"ip_white":                 "192.168.0.1",
						"security_groups.#":        "0",
						"ui_proxy_conn_addrs.#":    acctest.CHECKSET,
						"zk_conn_addrs.#":          acctest.CHECKSET,
						"slb_conn_addrs.#":         acctest.CHECKSET,
						"vpc_id":                   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"immediate_delete_flag"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
		},
	})
}

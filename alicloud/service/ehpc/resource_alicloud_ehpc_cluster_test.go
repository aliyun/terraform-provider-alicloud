package ehpc_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ehpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass ram_node_types -- pre-existing stock TypeList in a pure directory migration (no schema/CRUD/behavior change); the V1 ehpc_cluster ACC is unavailable in the EHPC 1.0 stop-sale environment, so a convergent two-member reorder apply cannot be run

func testSweepEhpcCluster(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
	}

	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = region
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	for {
		action := "ListClusters"
		response, err := client.RpcGet("EHPC", "2018-04-12", action, request, nil)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Clusters.ClusterInfoSimple", response)
		if err != nil {
			log.Println(err)
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			itemName := fmt.Sprint(item["Name"])
			itemId := fmt.Sprint(item["Id"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(itemName, prefix) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ehpc cluster: %s(%s) ", itemId, itemName)
				continue
			}
			log.Printf("[Info] Delete Ehpc cluster: %s(%s)", itemId, itemName)

			{
				action := "DeleteCluster"
				request := map[string]interface{}{
					"ClusterId":       itemId,
					"RegionId":        client.RegionId,
					"ReleaseInstance": "true",
				}
				request["ClientToken"] = helper.BuildClientToken("DeleteCluster")
				response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
				logs.AddDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] Deleting ephc cluster %s got an error: %s", itemId, err)
				}
			}

		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudEhpcCluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.EhpcSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEhpcClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ehpc.NewEhpcService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEhpcCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sehpccluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEhpcClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name":          "${var.name}",
					"deploy_mode":           "Simple",
					"os_tag":                "CentOS_7.6_64",
					"manager_count":         "1",
					"manager_instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
					"compute_count":         "1",
					"compute_instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
					"login_count":           "1",
					"login_instance_type":   "${data.alicloud_instance_types.default.instance_types.0.id}",
					"volume_protocol":       "nfs",
					"volume_id":             "${alicloud_nas_file_system.default.id}",
					"volume_mountpoint":     "${alicloud_nas_mount_target.default.mount_target_domain}",
					"password":              "Your-password123",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":               "${data.alicloud_zones.default.zones.0.id}",
					"release_instance":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":          name,
						"deploy_mode":           "Simple",
						"os_tag":                "CentOS_7.6_64",
						"manager_count":         "1",
						"manager_instance_type": acctest.CHECKSET,
						"compute_count":         "1",
						"compute_instance_type": acctest.CHECKSET,
						"login_count":           "1",
						"login_instance_type":   acctest.CHECKSET,
						"volume_protocol":       "nfs4",
						"volume_id":             acctest.CHECKSET,
						"volume_mountpoint":     acctest.CHECKSET,
						"vswitch_id":            acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":          "${data.alicloud_images.default.images.0.id}",
					"image_owner_alias": "system",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id":          acctest.CHECKSET,
						"image_owner_alias": "system",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "ram_role_name", "system_disk_level", "system_disk_type", "ram_node_types", "plugin", "resource_group_id", "domain", "volume_mount_option", "zone_id", "compute_enable_ht", "ecs_charge_type", "release_instance", "cluster_version", "input_file_url", "system_disk_size", "compute_spot_strategy", "without_elastic_ip", "additional_volumes", "security_group_name", "period", "compute_spot_price_limit", "manager_count", "job_queue", "without_agent", "auto_renew", "is_compute_ess", "ehpc_version", "remote_vis_enable", "auto_renew_period", "period_unit", "volume_protocol"},
			},
		},
	})
}

func TestAccAlicloudEhpcCluster_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.EhpcSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEhpcClusterMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ehpc.NewEhpcService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEhpcCluster")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sehpccluster%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEhpcClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":           "${var.name}",
					"cluster_name":          "${var.name}",
					"volume_protocol":       "nfs",
					"volume_id":             "${alicloud_nas_file_system.default.id}",
					"volume_mountpoint":     "${alicloud_nas_mount_target.default.mount_target_domain}",
					"deploy_mode":           "Simple",
					"image_id":              "${data.alicloud_images.default.images.0.id}",
					"image_owner_alias":     "system",
					"cluster_version":       "1.0",
					"volume_type":           "nas",
					"remote_directory":      "/",
					"scheduler_type":        "pbs",
					"account_type":          "nis",
					"ha_enable":             "false",
					"os_tag":                "CentOS_7.6_64",
					"manager_count":         "1",
					"manager_instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
					"compute_count":         "1",
					"compute_instance_type": "${data.alicloud_instance_types.default.instance_types.0.id}",
					"login_count":           "1",
					"login_instance_type":   "${data.alicloud_instance_types.default.instance_types.0.id}",
					"password":              "Your-password123",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":               "${data.alicloud_zones.default.zones.0.id}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"system_disk_level":     "PL0",
					"system_disk_size":      "40",
					"system_disk_type":      "cloud_essd",
					"ecs_charge_type":       "PostPaid",
					"client_version":        "1.0.1",
					"remote_vis_enable":     "false",
					"additional_volumes": []map[string]interface{}{
						{
							"job_queue":        "high",
							"local_directory":  "/ff",
							"location":         "PublicCloud",
							"remote_directory": "/test",
							"roles": []map[string]interface{}{
								{
									"name": "Compute",
								},
							},
							"volume_id":           "${alicloud_nas_file_system.default1.id}",
							"volume_mountpoint":   "${alicloud_nas_mount_target.default1.mount_target_domain}",
							"volume_mount_option": "-t nfs -o vers=4.0",
							"volume_protocol":     "nfs",
							"volume_type":         "nas",
						},
					},
					"application": []map[string]interface{}{
						{
							"tag": "singularity_3.8.3",
						},
					},
					"post_install_script": []map[string]interface{}{
						{
							"args": "bashfile.sh",
							"url":  "/opt/job.sh",
						},
					},
					"ram_node_types":      []string{"manager"},
					"security_group_name": "${var.name}",
					"volume_mount_option": "-t nfs -o vers=4.0",
					"without_agent":       "false",
					"without_elastic_ip":  "false",
					"ehpc_version":        "1.0.0",
					"is_compute_ess":      "false",
					"release_instance":    "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"cluster_name":          name,
						"volume_protocol":       "nfs",
						"volume_id":             acctest.CHECKSET,
						"deploy_mode":           "Simple",
						"image_id":              acctest.CHECKSET,
						"volume_mountpoint":     acctest.CHECKSET,
						"image_owner_alias":     "system",
						"cluster_version":       "1.0",
						"volume_type":           "nas",
						"remote_directory":      "/",
						"scheduler_type":        "pbs",
						"account_type":          "nis",
						"ha_enable":             "false",
						"os_tag":                "CentOS_7.6_64",
						"manager_count":         "1",
						"manager_instance_type": acctest.CHECKSET,
						"compute_count":         "1",
						"compute_instance_type": acctest.CHECKSET,
						"login_count":           "1",
						"login_instance_type":   acctest.CHECKSET,
						"vswitch_id":            acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
						"system_disk_level":     "PL0",
						"system_disk_size":      "40",
						"system_disk_type":      "cloud_essd",
						"client_version":        "1.0.1",
						"additional_volumes.#":  "1",
						"application.#":         "1",
						"post_install_script.#": "1",
						"remote_vis_enable":     "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "ram_role_name", "system_disk_level", "system_disk_type", "ram_node_types", "plugin", "resource_group_id", "domain", "volume_mount_option", "zone_id", "compute_enable_ht", "ecs_charge_type", "release_instance", "cluster_version", "input_file_url", "system_disk_size", "compute_spot_strategy", "without_elastic_ip", "additional_volumes", "security_group_name", "period", "compute_spot_price_limit", "manager_count", "job_queue", "without_agent", "auto_renew", "is_compute_ess", "ehpc_version", "remote_vis_enable", "auto_renew_period", "period_unit"},
			},
		},
	})
}

var AlicloudEhpcClusterMap0 = map[string]string{}

func AlicloudEhpcClusterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
data "alicloud_zones" default {
  available_resource_creation  = "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}
data "alicloud_instance_types" "default" {
 availability_zone = data.alicloud_zones.default.zones.0.id
}
variable "storage_type" {
  default = "Capacity"
}
resource "alicloud_nas_file_system" "default" {
  storage_type = var.storage_type
  protocol_type = "NFS"
}
resource "alicloud_nas_mount_target" "default" {
	file_system_id = alicloud_nas_file_system.default.id
	access_group_name = "DEFAULT_VPC_GROUP_NAME"
	vswitch_id = data.alicloud_vswitches.default.ids.0
}
resource "alicloud_nas_file_system" "default1" {
  storage_type = var.storage_type
  protocol_type = "NFS"
}
resource "alicloud_nas_mount_target" "default1" {
	file_system_id = alicloud_nas_file_system.default1.id
	access_group_name = "DEFAULT_VPC_GROUP_NAME"
	vswitch_id = data.alicloud_vswitches.default.ids.0
}
data "alicloud_images" "default" {
  name_regex  = "^centos_7_6_x64*"
  owners      = "system"
}

`, name)
}

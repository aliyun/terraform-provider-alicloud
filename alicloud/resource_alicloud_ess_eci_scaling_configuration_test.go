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
		"alicloud_ess_eci_scaling_configuration",
		&resource.Sweeper{
			Name: "alicloud_ess_eci_scaling_configuration",
			F:    testSweepEciScalingConfiguration,
		})
}

func testSweepEciScalingConfiguration(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc`",
	}
	var response map[string]interface{}
	conn, err := client.NewEssClient()
	if err != nil {
		return WrapError(err)
	}
	action := "DescribeEciScalingConfigurations"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_eci_scaling_configuration", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	resp, err := jsonpath.Get("$.ScalingConfigurations", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ScalingConfigurations", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := item["ScalingConfigurationName"].(string)
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Eci configuration: %s ", name)
			continue
		}
		log.Printf("[INFO] Delete Eci configuration: %s ", name)
		action := "DeleteEciScalingConfiguration"
		request := map[string]interface{}{
			"ScalingConfigurationId": item["ScalingConfigurationId"],
		}
		request["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(time.Minute*10, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Eci Scaling Configuration (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlicloudEssEciScalingConfigurationBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ess_eci_scaling_configuration.default"
	ra := resourceAttrInit(resourceId, AlicloudEssEciScalingConfigurationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEssEciScalingConfiguration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudEciContainerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssEciScalingConfiguration)
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
					"scaling_group_id":           "${alicloud_ess_scaling_group.default.id}",
					"scaling_configuration_name": name,
					"description":                "desc",
					"security_group_id":          "sg-bp1hi5tpb5c3e51a15pf",
					"container_group_name":       "containerGroupName",
					"restart_policy":             "restartPolicy",
					"cpu":                        "2",
					"memory":                     "4",
					"resource_group_id":          "resourceGroupId",
					"dns_policy":                 "dnsPolicy",
					"enable_sls":                 "true",
					"ram_role_name":              "ramRoleName",
					"spot_strategy":              "SpotWithPriceLimit",
					"spot_price_limit":           "1.1",
					"auto_create_eip":            "true",
					"eip_bandwidth":              "1",
					"host_name":                  "hostname",
					"ingress_bandwidth":          "1",
					"egress_bandwidth":           "1",
					"tags": map[string]string{
						"name": "tf-test",
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "password",
							"server":   "server",
							"username": "username",
						},
					},
					"containers": []map[string]interface{}{
						{
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "key",
									"value": "value",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "path",
									"name":       "name",
									"read_only":  "true",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"ports": []map[string]interface{}{
								{
									"protocol": "protocol",
									"port":     "1",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "key",
									"value": "value",
								},
							},
							"working_dir":       "workingDir",
							"args":              []string{"arg"},
							"cpu":               "1",
							"gpu":               "1",
							"memory":            "1",
							"name":              "name",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "policy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "path",
									"name":       "name",
									"read_only":  "true",
								},
							},
							"commands": []string{"cmd"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content",
									"path":    "path",
								},
							},
							"disk_volume_disk_id":   "disk_volume_disk_id",
							"disk_volume_fs_type":   "disk_volume_fs_type",
							"disk_volume_disk_size": "1",
							"flex_volume_driver":    "flex_volume_driver",
							"flex_volume_fs_type":   "flex_volume_fs_type",
							"flex_volume_options":   "flex_volume_options",
							"nfs_volume_path":       "nfs_volume_path",
							"nfs_volume_read_only":  "true",
							"nfs_volume_server":     "nfs_volume_server",
							"name":                  "name",
							"type":                  "type",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames"},
							"ip":        "ip",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":             CHECKSET,
						"scaling_configuration_name":   name,
						"description":                  "desc",
						"security_group_id":            "sg-bp1hi5tpb5c3e51a15pf",
						"container_group_name":         "containerGroupName",
						"restart_policy":               "restartPolicy",
						"cpu":                          "2",
						"memory":                       "4",
						"resource_group_id":            "resourceGroupId",
						"dns_policy":                   "dnsPolicy",
						"enable_sls":                   "true",
						"ram_role_name":                "ramRoleName",
						"spot_strategy":                "SpotWithPriceLimit",
						"spot_price_limit":             "1.1",
						"auto_create_eip":              "true",
						"host_name":                    "hostname",
						"ingress_bandwidth":            "1",
						"egress_bandwidth":             "1",
						"tags.name":                    "tf-test",
						"image_registry_credentials.#": "1",
						"containers.#":                 "1",
						"init_containers.#":            "1",
						"volumes.#":                    "1",
						"host_aliases.#":               "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaling_configuration_name": "newName",
					"description":                "newDesc",
					"container_group_name":       "newName",
					"restart_policy":             "newPolicy",
					"cpu":                        "2",
					"memory":                     "2",
					"resource_group_id":          "newGroupId",
					"dns_policy":                 "newDnsPolicy",
					"enable_sls":                 "false",
					"ram_role_name":              "newRoleName",
					"spot_strategy":              "SpotAsPriceGo",
					"spot_price_limit":           "1.2",
					"auto_create_eip":            "false",
					"eip_bandwidth":              "3",
					"host_name":                  "newHostName",
					"ingress_bandwidth":          "2",
					"egress_bandwidth":           "2",
					"tags": map[string]string{
						"name": "tf-test2",
					},
					"image_registry_credentials": []map[string]interface{}{
						{
							"password": "newPassword",
							"server":   "newServer",
							"username": "newUserName",
						},
					},
					"containers": []map[string]interface{}{
						{
							"ports": []map[string]interface{}{
								{
									"protocol": "newProtocol",
									"port":     "2",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "newKey",
									"value": "newValue",
								},
							},
							"working_dir":       "newWorkingDir",
							"args":              []string{"arg2"},
							"cpu":               "2",
							"gpu":               "2",
							"memory":            "2",
							"name":              "newName",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "newPolicy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "newPath",
									"name":       "newName",
									"read_only":  "false",
								},
							},
							"commands": []string{"cmd2"},
						},
					},
					"init_containers": []map[string]interface{}{
						{
							"ports": []map[string]interface{}{
								{
									"protocol": "newProtocol",
									"port":     "2",
								},
							},
							"environment_vars": []map[string]interface{}{
								{
									"key":   "newKey",
									"value": "newValue",
								},
							},
							"working_dir":       "newWorkingDir",
							"args":              []string{"arg2"},
							"cpu":               "2",
							"gpu":               "2",
							"memory":            "2",
							"name":              "newName",
							"image":             "registry-vpc.aliyuncs.com/eci_open/alpine:3.5",
							"image_pull_policy": "newPolicy",
							"volume_mounts": []map[string]interface{}{
								{
									"mount_path": "newPath",
									"name":       "newName",
									"read_only":  "false",
								},
							},
							"commands": []string{"cmd2"},
						},
					},
					"volumes": []map[string]interface{}{
						{
							"config_file_volume_config_file_to_paths": []map[string]interface{}{
								{
									"content": "content2",
									"path":    "path2",
								},
							},
							"disk_volume_disk_id":   "disk_volume_disk_id2",
							"disk_volume_fs_type":   "disk_volume_fs_type2",
							"disk_volume_disk_size": "2",
							"flex_volume_driver":    "flex_volume_driver2",
							"flex_volume_fs_type":   "flex_volume_fs_type2",
							"flex_volume_options":   "flex_volume_options2",
							"nfs_volume_path":       "nfs_volume_path2",
							"nfs_volume_read_only":  "false",
							"nfs_volume_server":     "nfs_volume_server2",
							"name":                  "name2",
							"type":                  "type2",
						},
					},
					"host_aliases": []map[string]interface{}{
						{
							"hostnames": []string{"hostnames2"},
							"ip":        "ip2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_configuration_name":   "newName",
						"description":                  "newDesc",
						"container_group_name":         "newName",
						"restart_policy":               "newPolicy",
						"cpu":                          "2",
						"memory":                       "2",
						"resource_group_id":            "newGroupId",
						"dns_policy":                   "newDnsPolicy",
						"enable_sls":                   "false",
						"ram_role_name":                "newRoleName",
						"spot_strategy":                "SpotAsPriceGo",
						"spot_price_limit":             "1.2",
						"auto_create_eip":              "false",
						"eip_bandwidth":                "3",
						"host_name":                    "newHostName",
						"ingress_bandwidth":            "2",
						"egress_bandwidth":             "2",
						"tags.name":                    "tf-test2",
						"image_registry_credentials.#": "1",
						"containers.#":                 "1",
						"init_containers.#":            "1",
						"volumes.#":                    "1",
						"host_aliases.#":               "1",
					}),
				),
			},
		},
	})
}

var AlicloudEssEciScalingConfigurationMap = map[string]string{}

func resourceEssEciScalingConfiguration(name string) string {
	return fmt.Sprintf(`
	%s

	variable "name" {
		default = "%s"
	}

	resource "alicloud_security_group" "default1" {
	  name   = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	}

	resource "alicloud_ess_scaling_group" "default" {
		min_size = 0
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alicloud_vswitch.default.id}"]
		group_type = "ECI"
	}`, EcsInstanceCommonTestCase, name)
}

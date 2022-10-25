package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEhpcCluster_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster.default"
	checkoutSupportedRegions(t, true, connectivity.EhpcSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEhpcClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sehpccluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcClusterBasicDependence0)
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
					"password":              "your-password123",
					"vswitch_id":            "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":               "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":          name,
						"deploy_mode":           "Simple",
						"os_tag":                "CentOS_7.6_64",
						"manager_count":         "1",
						"manager_instance_type": CHECKSET,
						"compute_count":         "1",
						"compute_instance_type": CHECKSET,
						"login_count":           "1",
						"login_instance_type":   CHECKSET,
						"volume_protocol":       "nfs",
						"volume_id":             CHECKSET,
						"volume_mountpoint":     CHECKSET,
						"vswitch_id":            CHECKSET,
						"vpc_id":                CHECKSET,
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
						"image_id":          CHECKSET,
						"image_owner_alias": "system",
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

func TestAccAlicloudEhpcCluster_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_cluster.default"
	checkoutSupportedRegions(t, true, connectivity.EhpcSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEhpcClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sehpccluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcClusterBasicDependence0)
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
					"password":              "your-password123",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"cluster_name":          name,
						"volume_protocol":       "nfs",
						"volume_id":             CHECKSET,
						"deploy_mode":           "Simple",
						"image_id":              CHECKSET,
						"volume_mountpoint":     CHECKSET,
						"image_owner_alias":     "system",
						"cluster_version":       "1.0",
						"volume_type":           "nas",
						"remote_directory":      "/",
						"scheduler_type":        "pbs",
						"account_type":          "nis",
						"ha_enable":             "false",
						"os_tag":                "CentOS_7.6_64",
						"manager_count":         "1",
						"manager_instance_type": CHECKSET,
						"compute_count":         "1",
						"compute_instance_type": CHECKSET,
						"login_count":           "1",
						"login_instance_type":   CHECKSET,
						"vswitch_id":            CHECKSET,
						"vpc_id":                CHECKSET,
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
	name_regex = "default-NODELETING"
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

func TestUnitAlicloudEhpcCluster(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description":              "CreateEhpcClusterValue",
		"cluster_name":             "CreateEhpcClusterValue",
		"volume_protocol":          "CreateEhpcClusterValue",
		"volume_id":                "CreateEhpcClusterValue",
		"volume_mountpoint":        "CreateEhpcClusterValue",
		"deploy_mode":              "CreateEhpcClusterValue",
		"image_id":                 "CreateEhpcClusterValue",
		"image_owner_alias":        "CreateEhpcClusterValue",
		"cluster_version":          "CreateEhpcClusterValue",
		"volume_type":              "CreateEhpcClusterValue",
		"remote_directory":         "CreateEhpcClusterValue",
		"scheduler_type":           "CreateEhpcClusterValue",
		"account_type":             "CreateEhpcClusterValue",
		"ha_enable":                false,
		"os_tag":                   "CreateEhpcClusterValue",
		"manager_count":            2,
		"manager_instance_type":    "CreateEhpcClusterValue",
		"compute_count":            1,
		"compute_instance_type":    "CreateEhpcClusterValue",
		"login_count":              1,
		"login_instance_type":      "CreateEhpcClusterValue",
		"password":                 "CreateEhpcClusterValue",
		"vswitch_id":               "CreateEhpcClusterValue",
		"vpc_id":                   "CreateEhpcClusterValue",
		"zone_id":                  "CreateEhpcClusterValue",
		"auto_renew":               true,
		"auto_renew_period":        1,
		"compute_enable_ht":        true,
		"compute_spot_price_limit": "CreateEhpcClusterValue",
		"compute_spot_strategy":    "CreateEhpcClusterValue",
		"domain":                   "CreateEhpcClusterValue",
		"ecs_charge_type":          "CreateEhpcClusterValue",
		"ehpc_version":             "CreateEhpcClusterValue",
		"input_file_url":           "CreateEhpcClusterValue",
		"is_compute_ess":           true,
		"job_queue":                "CreateEhpcClusterValue",
		"period":                   1,
		"period_unit":              "CreateEhpcClusterValue",
		"plugin":                   "CreateEhpcClusterValue",
		"key_pair_name":            "CreateEhpcClusterValue",
		"ram_role_name":            "CreateEhpcClusterValue",
		"remote_vis_enable":        true,
		"resource_group_id":        "CreateEhpcClusterValue",
		"scc_cluster_id":           "CreateEhpcClusterValue",
		"security_group_id":        "CreateEhpcClusterValue",
		"security_group_name":      "CreateEhpcClusterValue",
		"system_disk_level":        "CreateEhpcClusterValue",
		"system_disk_size":         40,
		"without_agent":            true,
		"without_elastic_ip":       true,
		"system_disk_type":         "CreateEhpcClusterValue",
		"volume_mount_option":      "CreateEhpcClusterValue",
		"client_version":           "CreateEhpcClusterValue",
		"release_instance":         true,
		"additional_volumes": []map[string]interface{}{
			{
				"job_queue":        "CreateEhpcClusterValue",
				"local_directory":  "CreateEhpcClusterValue",
				"location":         "CreateEhpcClusterValue",
				"remote_directory": "CreateEhpcClusterValue",
				"roles": []map[string]interface{}{
					{
						"name": "CreateEhpcClusterValue",
					},
				},
				"volume_id":           "CreateEhpcClusterValue",
				"volume_mount_option": "CreateEhpcClusterValue",
				"volume_mountpoint":   "CreateEhpcClusterValue",
				"volume_protocol":     "CreateEhpcClusterValue",
				"volume_type":         "CreateEhpcClusterValue",
			},
		},
		"application": []map[string]interface{}{
			{
				"tag": "CreateEhpcClusterValue",
			},
		},
		"post_install_script": []map[string]interface{}{
			{
				"args": "CreateEhpcClusterValue",
				"url":  "CreateEhpcClusterValue",
			},
		},
		"ram_node_types": []string{"CreateEhpcClusterValue"},
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"ClusterInfo": map[string]interface{}{
			"Status":           "running",
			"VpcId":            "CreateEhpcClusterValue",
			"KeyPairName":      "CreateEhpcClusterValue",
			"EcsChargeType":    "CreateEhpcClusterValue",
			"SecurityGroupId":  "CreateEhpcClusterValue",
			"SccClusterId":     "CreateEhpcClusterValue",
			"CreateTime":       "CreateEhpcClusterValue",
			"AccountType":      "CreateEhpcClusterValue",
			"VolumeProtocol":   "CreateEhpcClusterValue",
			"Description":      "CreateEhpcClusterValue",
			"VolumeId":         "CreateEhpcClusterValue",
			"HaEnable":         false,
			"BaseOsTag":        "CreateEhpcClusterValue",
			"Name":             "CreateEhpcClusterValue",
			"ImageId":          "CreateEhpcClusterValue",
			"SchedulerType":    "CreateEhpcClusterValue",
			"DeployMode":       "CreateEhpcClusterValue",
			"ImageOwnerAlias":  "CreateEhpcClusterValue",
			"OsTag":            "CreateEhpcClusterValue",
			"VolumeMountpoint": "CreateEhpcClusterValue",
			"RemoteDirectory":  "CreateEhpcClusterValue",
			"RegionId":         "CreateEhpcClusterValue",
			"VSwitchId":        "CreateEhpcClusterValue",
			"ImageName":        "CreateEhpcClusterValue",
			"VolumeType":       "CreateEhpcClusterValue",
			"Location":         "CreateEhpcClusterValue",
			"Id":               "CreateEhpcClusterValue",
			"ClientVersion":    "CreateEhpcClusterValue",
			"Applications": map[string]interface{}{
				"ApplicationInfo": []interface{}{
					map[string]interface{}{
						"Tag":     "CreateEhpcClusterValue",
						"Name":    "CreateEhpcClusterValue",
						"Version": "CreateEhpcClusterValue",
					}},
			},
			"PostInstallScripts": map[string]interface{}{
				"PostInstallScriptInfo": []interface{}{
					map[string]interface{}{
						"Url":  "CreateEhpcClusterValue",
						"Args": "CreateEhpcClusterValue",
					},
				},
			},
			"EcsInfo": map[string]interface{}{
				"Manager": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        2,
				},
				"Compute": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        1,
				},
				"Login": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        1,
				},
				"ProxyMgr": map[string]interface{}{
					"InstanceType": "CreateEhpcClusterValue",
					"Count":        1,
				},
			},
			"OnPremiseInfo": []interface{}{
				map[string]interface{}{
					"Type":     "CreateEhpcClusterValue",
					"HostName": "CreateEhpcClusterValue",
					"IP":       "CreateEhpcClusterValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"ClusterId": "CreateEhpcClusterValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ehpc_cluster", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcClusterCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEhpcClusterCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcClusterUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"description":       "UpdateEhpcClusterValue",
		"cluster_name":      "UpdateEhpcClusterValue",
		"image_id":          "UpdateEhpcClusterValue",
		"image_owner_alias": "UpdateEhpcClusterValue",
	}
	diff, err := newInstanceDiff("alicloud_ehpc_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ClusterInfo": map[string]interface{}{
			"Name":            "UpdateEhpcClusterValue",
			"Description":     "UpdateEhpcClusterValue",
			"ImageId":         "UpdateEhpcClusterValue",
			"ImageOwnerAlias": "UpdateEhpcClusterValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyClusterAttributes" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEhpcClusterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_ehpc_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeCluster" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEhpcClusterRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcClusterDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ehpc_cluster", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_cluster"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ClusterNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteCluster" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					if errorCodes[retryIndex] == "ClusterNotFound" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DescribeCluster" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEhpcClusterDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

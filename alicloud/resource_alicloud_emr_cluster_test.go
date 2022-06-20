package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_emr_cluster", &resource.Sweeper{
		Name: "alicloud_emr_cluster",
		F:    testSweepEmrCluster,
	})
}

func testSweepEmrCluster(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := emr.CreateListClustersRequest()
	statusList := []string{"CREATING", "RUNNING", "IDLE"}
	req.StatusList = &statusList
	req.PageNumber = requests.Integer(strconv.Itoa(1))
	req.PageSize = requests.Integer(strconv.Itoa(PageSizeMedium))
	req.DefaultStatus = requests.Boolean(strconv.FormatBool(true))
	//vpcService := VpcService{client}
	for {
		raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
			return emrClient.ListClusters(req)
		})

		if err != nil {
			return fmt.Errorf("Error retrieving EMR Cluster: %s", err)
		}

		resp, _ := raw.(*emr.ListClustersResponse)
		if resp == nil || len(resp.Clusters.ClusterInfo) == 0 {
			break
		}
		for _, v := range resp.Clusters.ClusterInfo {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(v.Id, prefix) {
					skip = false
				}
			}
			// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
			//if skip {
			//	if need, err := vpcService.needSweepVpc(v., v.VSwitchId); err == nil {
			//		skip = !need
			//	}
			//
			//}
			if skip {
				log.Printf("[INFO] Skipping emr: %s (%s)", v.Name, v.Id)
				continue
			}
			request := emr.CreateReleaseClusterRequest()
			request.Id = v.Id
			request.ForceRelease = requests.NewBoolean(true)

			raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
				return emrClient.ReleaseCluster(request)
			})

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, v.Id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}
	return nil
}

func TestAccAlicloudEmrCluster_basic(t *testing.T) {
	var v *emr.DescribeClusterV2Response
	resourceId := "alicloud_emr_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrClusterCommonConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  testAccAlicloudEmrClusterDestroy,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name,
					"emr_ver":                   "${data.alicloud_emr_main_versions.default.main_versions.0.emr_version}",
					"cluster_type":              "${data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0}",
					"deposit_type":              "HALF_MANAGED",
					"high_availability_enable":  "false",
					"zone_id":                   "${data.alicloud_emr_instance_types.default.types.0.zone_id}",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"is_open_public_ip":         "true",
					"charge_type":               "PostPaid",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"user_defined_emr_ecs_role": "${alicloud_ram_role.default.name}",
					"ssh_enable":                "true",
					"master_pwd":                "ABCtest1234!",
					"tags": map[string]interface{}{
						"Created": "TF",
						"For":     "acceptance test",
					},

					"host_group": []map[string]interface{}{
						{
							"host_group_type":   "MASTER",
							"node_count":        "1",
							"instance_type":     "${data.alicloud_emr_instance_types.default.types.0.id}",
							"disk_type":         "${data.alicloud_emr_disk_types.data_disk.types.0.value}",
							"disk_capacity":     "${data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160}",
							"disk_count":        "1",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160}",
						},
						{
							"host_group_type":   "CORE",
							"node_count":        "2",
							"instance_type":     "${data.alicloud_emr_instance_types.default.types.0.id}",
							"disk_type":         "${data.alicloud_emr_disk_types.data_disk.types.0.value}",
							"disk_capacity":     "${data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160}",
							"disk_count":        "4",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"emr_ver":      CHECKSET,
						"cluster_type": CHECKSET,
						"charge_type":  "PostPaid",
						"zone_id":      CHECKSET,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssh_enable", "master_pwd", "is_open_public_ip", "host_group"},
			},
		},
	})
}

func TestAccAlicloudEmrCluster_local_storage(t *testing.T) {
	var v *emr.DescribeClusterV2Response
	resourceId := "alicloud_emr_cluster.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrClusterLocalStorageConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  testAccAlicloudEmrClusterDestroy,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name,
					"emr_ver":                   "${data.alicloud_emr_main_versions.default.main_versions.0.emr_version}",
					"cluster_type":              "${data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0}",
					"deposit_type":              "HALF_MANAGED",
					"high_availability_enable":  "false",
					"zone_id":                   "${data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id}",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"is_open_public_ip":         "true",
					"charge_type":               "PostPaid",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"user_defined_emr_ecs_role": "${alicloud_ram_role.default.name}",
					"ssh_enable":                "true",
					"master_pwd":                "ABCtest1234!",

					"bootstrap_action": []map[string]interface{}{
						{
							"name":                    "bootstrap_test",
							"path":                    "oss://emrtf/tf.sh",
							"arg":                     "--a=b",
							"execution_target":        "core_group",
							"execution_moment":        "BEFORE_INSTALL",
							"execution_fail_strategy": "FAILED_CONTINUE",
						},
					},

					"host_group": []map[string]interface{}{
						{
							"host_group_name":   "master_group",
							"host_group_type":   "MASTER",
							"node_count":        "1",
							"instance_type":     "${data.alicloud_emr_instance_types.cloud_disk.types.0.id}",
							"disk_type":         "${data.alicloud_emr_disk_types.data_disk.types.0.value}",
							"disk_capacity":     "${data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160}",
							"disk_count":        "1",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160}",
						},
						{
							"host_group_name":   "core_group",
							"host_group_type":   "CORE",
							"node_count":        "3",
							"instance_type":     "${data.alicloud_emr_instance_types.local_disk.types.0.id}",
							"disk_type":         "local_disk",
							"disk_capacity":     "${data.alicloud_emr_instance_types.local_disk.types.0.local_storage_capacity}",
							"disk_count":        "1",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"emr_ver":      CHECKSET,
						"cluster_type": CHECKSET,
						"charge_type":  "PostPaid",
						"zone_id":      CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEmrCluster_gateway(t *testing.T) {
	var v *emr.DescribeClusterV2Response
	resourceId := "alicloud_emr_cluster.default_gateway"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrClusterGatewayConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  testAccAlicloudEmrClusterDestroy,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                      name,
					"emr_ver":                   "${data.alicloud_emr_main_versions.default.main_versions.0.emr_version}",
					"cluster_type":              "GATEWAY",
					"zone_id":                   "${data.alicloud_emr_instance_types.default.types.0.zone_id}",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"is_open_public_ip":         "true",
					"charge_type":               "PostPaid",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"user_defined_emr_ecs_role": "${alicloud_ram_role.default.name}",
					"ssh_enable":                "true",
					"master_pwd":                "ABCtest1234!",
					"related_cluster_id":        "${alicloud_emr_cluster.default.id}",

					"host_group": []map[string]interface{}{
						{
							"host_group_type":   "GATEWAY",
							"node_count":        "1",
							"instance_type":     "${data.alicloud_emr_instance_types.gateway.types.0.id}",
							"disk_type":         "${data.alicloud_emr_disk_types.gateway_data_disk.types.0.value}",
							"disk_capacity":     "${data.alicloud_emr_disk_types.gateway_data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.gateway_data_disk.types.0.min : 160}",
							"disk_count":        "1",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.gateway_system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.gateway_system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.gateway_system_disk.types.0.min : 160}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"emr_ver":      CHECKSET,
						"cluster_type": "GATEWAY",
						"charge_type":  "PostPaid",
						"zone_id":      CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEmrCluster_multicluster(t *testing.T) {
	var v *emr.DescribeClusterV2Response
	resourceId := "alicloud_emr_cluster.default.4"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EmrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc%sEmrClusterConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEmrClusterCommonConfigDependence)

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
					"count":                     "5",
					"name":                      name,
					"emr_ver":                   "${data.alicloud_emr_main_versions.default.main_versions.0.emr_version}",
					"cluster_type":              "${data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0}",
					"deposit_type":              "HALF_MANAGED",
					"high_availability_enable":  "true",
					"zone_id":                   "${data.alicloud_emr_instance_types.default.types.0.zone_id}",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"is_open_public_ip":         "true",
					"charge_type":               "PostPaid",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"user_defined_emr_ecs_role": "${alicloud_ram_role.default.name}",
					"ssh_enable":                "true",
					"master_pwd":                "ABCtest1234!",

					"host_group": []map[string]interface{}{
						{
							"host_group_type":   "MASTER",
							"node_count":        "2",
							"instance_type":     "${data.alicloud_emr_instance_types.default.types.0.id}",
							"disk_type":         "${data.alicloud_emr_disk_types.data_disk.types.0.value}",
							"disk_capacity":     "${data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160}",
							"disk_count":        "1",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160}",
						},
						{
							"host_group_type":   "CORE",
							"node_count":        "3",
							"instance_type":     "${data.alicloud_emr_instance_types.default.types.0.id}",
							"disk_type":         "${data.alicloud_emr_disk_types.data_disk.types.0.value}",
							"disk_capacity":     "${data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160}",
							"disk_count":        "4",
							"sys_disk_type":     "${data.alicloud_emr_disk_types.system_disk.types.0.value}",
							"sys_disk_capacity": "${data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"emr_ver":      CHECKSET,
						"cluster_type": CHECKSET,
						"charge_type":  "PostPaid",
						"zone_id":      CHECKSET,
						"tags.%":       NOSET,
						"tags.Created": NOSET,
						"tags.For":     NOSET,
					}),
				),
			},
		},
	})
}

func resourceEmrClusterCommonConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, EmrCommonTestCase, name)
}

func resourceEmrClusterGatewayConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, EmrGatewayTestCase, name)
}

func resourceEmrClusterLocalStorageConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, EmrLocalStorageTestCase, name)
}

package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			flag := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(v.ClusterId, prefix) {
					flag = true
				}
			}
			if flag {
				request := emr.CreateReleaseClusterRequest()
				request.Id = v.ClusterId
				request.ForceRelease = requests.NewBoolean(true)

				raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
					return emrClient.ReleaseCluster(request)
				})

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, v.ClusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
				}

				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
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
	name := fmt.Sprintf("tf-testAccEmrClusterConfig%d", rand)
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
					"emr_ver":                   "EMR-3.22.0",
					"cluster_type":              "HADOOP",
					"deposit_type":              "HALF_MANAGED",
					"high_availability_enable":  "true",
					"option_software_list":      []string{"HBASE", "PRESTO"},
					"zone_id":                   "cn-huhehaote-a",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"is_open_public_ip":         "true",
					"charge_type":               "PostPaid",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"user_defined_emr_ecs_role": "EMRUserDefineRole-Role1",
					"ssh_enable":                "true",
					"master_pwd":                "ABCtest1234!",

					"host_group": []map[string]interface{}{
						map[string]interface{}{
							"host_group_type":   "MASTER",
							"node_count":        "2",
							"instance_type":     "ecs.g5.xlarge",
							"disk_type":         "cloud_ssd",
							"disk_capacity":     "80",
							"disk_count":        "1",
							"sys_disk_type":     "cloud_ssd",
							"sys_disk_capacity": "80",
						},
						map[string]interface{}{
							"host_group_type":   "CORE",
							"node_count":        "3",
							"instance_type":     "ecs.g5.xlarge",
							"disk_type":         "cloud_ssd",
							"disk_capacity":     "80",
							"disk_count":        "4",
							"sys_disk_type":     "cloud_ssd",
							"sys_disk_capacity": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"emr_ver":      "EMR-3.22.0",
						"cluster_type": "HADOOP",
						"zone_id":      "cn-huhehaote-a",
						"charge_type":  "PostPaid",
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
	name := fmt.Sprintf("tf-testAccEmrClusterConfig%d", rand)
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
					"emr_ver":                   "EMR-3.22.0",
					"cluster_type":              "HADOOP",
					"deposit_type":              "HALF_MANAGED",
					"high_availability_enable":  "true",
					"option_software_list":      []string{"HBASE", "PRESTO"},
					"zone_id":                   "cn-huhehaote-a",
					"security_group_id":         "${alicloud_security_group.default.id}",
					"is_open_public_ip":         "true",
					"charge_type":               "PostPaid",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"user_defined_emr_ecs_role": "EMRUserDefineRole-Role1",
					"ssh_enable":                "true",
					"master_pwd":                "ABCtest1234!",

					"host_group": []map[string]interface{}{
						map[string]interface{}{
							"host_group_type":   "MASTER",
							"node_count":        "2",
							"instance_type":     "ecs.g5.xlarge",
							"disk_type":         "cloud_ssd",
							"disk_capacity":     "80",
							"disk_count":        "1",
							"sys_disk_type":     "cloud_ssd",
							"sys_disk_capacity": "80",
						},
						map[string]interface{}{
							"host_group_type":   "CORE",
							"node_count":        "3",
							"instance_type":     "ecs.g5.xlarge",
							"disk_type":         "cloud_ssd",
							"disk_capacity":     "80",
							"disk_count":        "4",
							"sys_disk_type":     "cloud_ssd",
							"sys_disk_capacity": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"emr_ver":      "EMR-3.22.0",
						"cluster_type": "HADOOP",
						"zone_id":      "cn-huhehaote-a",
						"charge_type":  "PostPaid",
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

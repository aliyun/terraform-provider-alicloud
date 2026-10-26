package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudElasticsearchSnapshotPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudElasticsearchSnapshotPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSnapshotSetting")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESSnapshotPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudElasticsearchSnapshotPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_elasticsearch_instance.default.id}",
					"enable":       true,
					"quartz_regex": "0 0 01 ? * * *",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"enable":       "true",
						"quartz_regex": "0 0 01 ? * * *",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quartz_regex": "0 0 02 ? * * *",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quartz_regex": "0 0 02 ? * * *",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable":       false,
					"quartz_regex": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable":       "false",
						"quartz_regex": REMOVEKEY,
					}),
				),
			},
			{
				// Destroy only the snapshot policy while the Elasticsearch instance
				// still exists, then verify the policy is disabled by reading it back.
				Config: AliCloudElasticsearchSnapshotPolicyBasicDependence0(name),
				Check: resource.ComposeTestCheckFunc(
					checkElasticsearchSnapshotSettingDisabled(),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchSnapshotPolicy_nonExistentInstance(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_snapshot_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudElasticsearchSnapshotPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSnapshotSetting")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sESSnapshotPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudElasticsearchSnapshotPolicyNonExistentInstanceDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "es-cn-tfnonexistent0001",
					"enable":       true,
					"quartz_regex": "0 0 01 ? * * *",
				}),
				ExpectError: regexp.MustCompile(`NotFound`),
			},
		},
	})
}

func checkElasticsearchSnapshotSettingDisabled() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_elasticsearch_instance.default"]
		if !ok {
			return fmt.Errorf("not found: alicloud_elasticsearch_instance.default")
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		elasticsearchService := &ElasticsearchService{client}
		object, err := elasticsearchService.DescribeSnapshotSetting(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("checking the snapshot setting of instance %s exists error: %s ", rs.Primary.ID, err)
		}
		enable, _ := snapshotSettingAttribute(object, "Enable", "enable").(bool)
		if enable {
			return fmt.Errorf("the snapshot policy of instance %s has not been disabled, enable: %v", rs.Primary.ID, enable)
		}
		return nil
	}
}

var AliCloudElasticsearchSnapshotPolicyMap0 = map[string]string{}

func AliCloudElasticsearchSnapshotPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	# The ACC test account has exhausted its VPC quota in cn-beijing, so the
	# instance is created in the pre-provisioned vswitch of the shared
	# default-NODELETING VPC. The region is pinned to cn-beijing in PreCheck
	# to match the fixed zone below.
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
		cidr_block = "172.16.0.0/12"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = "cn-beijing-l"
	}

	resource "alicloud_elasticsearch_instance" "default" {
		description       = var.name
		instance_category = "x-pack"
		version           = "7.10.0_with_X-Pack"
		payment_type      = "PayAsYouGo"
		zone_count        = 1
		vswitch_id        = data.alicloud_vswitches.default.ids.0
		password          = "YourPassword1234!"

		data_node_configuration {
			spec              = "elasticsearch.sn1ne.large.new"
			amount            = 2
			disk              = 20
			disk_type         = "cloud_essd"
			performance_level = "PL1"
		}

		kibana_configuration {
			spec   = "elasticsearch.n4.small"
			amount = 1
			disk   = 0
		}
	}
	`, name)
}

func AliCloudElasticsearchSnapshotPolicyNonExistentInstanceDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
`, name)
}

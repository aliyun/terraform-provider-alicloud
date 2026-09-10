package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudThreatDetectionBackupPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_backup_policy.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudThreatDetectionBackupPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccThreatDetectionBackupPolicy-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudThreatDetectionBackupPolicyBasicDependence)
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
					"backup_policy_name": name,
					"policy":             `{\"Exclude\":[\"/bin/\",\"/usr/bin/\",\"/sbin/\",\"/boot/\",\"/proc/\",\"/sys/\",\"/srv/\",\"/lib/\",\"/selinux/\",\"/usr/sbin/\",\"/run/\",\"/lib32/\",\"/lib64/\",\"/lost+found/\",\"/var/lib/kubelet/\",\"/var/lib/ntp/proc\",\"/var/lib/container\"],\"ExcludeSystemPath\":true,\"Include\":[],\"IsDefault\":1,\"Retention\":7,\"Schedule\":\"I|1668703620|PT24H\",\"Source\":[],\"SpeedLimiter\":\"\",\"UseVss\":true}`,
					"policy_version":     "2.0.0",
					"uuid_list":          []string{"${data.alicloud_threat_detection_assets.default.ids.0}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_policy_name": name,
						"policy":             "{\"Exclude\":[\"/bin/\",\"/usr/bin/\",\"/sbin/\",\"/boot/\",\"/proc/\",\"/sys/\",\"/srv/\",\"/lib/\",\"/selinux/\",\"/usr/sbin/\",\"/run/\",\"/lib32/\",\"/lib64/\",\"/lost+found/\",\"/var/lib/kubelet/\",\"/var/lib/ntp/proc\",\"/var/lib/container\"],\"ExcludeSystemPath\":true,\"Include\":[],\"IsDefault\":1,\"Retention\":7,\"Schedule\":\"I|1668703620|PT24H\",\"Source\":[],\"SpeedLimiter\":\"\",\"UseVss\":true}",
						"policy_version":     "2.0.0",
						"uuid_list.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_policy_name": name + "_update",
					"policy":             `{\"Exclude\":[\"/bin/\",\"/usr/bin/\",\"/sbin/\",\"/boot/\",\"/proc/\",\"/sys/\",\"/srv/\",\"/lib/\",\"/selinux/\",\"/usr/sbin/\",\"/run/\",\"/lib32/\",\"/lib64/\",\"/lost+found/\",\"/var/lib/kubelet/\",\"/var/lib/ntp/proc\",\"/var/lib/container\"],\"ExcludeSystemPath\":true,\"Include\":[],\"IsDefault\":1,\"Retention\":6,\"Schedule\":\"I|1668703620|PT24H\",\"Source\":[],\"SpeedLimiter\":\"\",\"UseVss\":true}`,
					"uuid_list":          []string{"${data.alicloud_threat_detection_assets.default.ids.0}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_policy_name": name + "_update",
						"policy":             "{\"Exclude\":[\"/bin/\",\"/usr/bin/\",\"/sbin/\",\"/boot/\",\"/proc/\",\"/sys/\",\"/srv/\",\"/lib/\",\"/selinux/\",\"/usr/sbin/\",\"/run/\",\"/lib32/\",\"/lib64/\",\"/lost+found/\",\"/var/lib/kubelet/\",\"/var/lib/ntp/proc\",\"/var/lib/container\"],\"ExcludeSystemPath\":true,\"Include\":[],\"IsDefault\":1,\"Retention\":6,\"Schedule\":\"I|1668703620|PT24H\",\"Source\":[],\"SpeedLimiter\":\"\",\"UseVss\":true}",
						"uuid_list.#":        "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy_region_id"},
			},
		},
	})
}

var resourceAlicloudThreatDetectionBackupPolicyMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudThreatDetectionBackupPolicyBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}
	data "alicloud_threat_detection_assets" "default" {
	  machine_types = "ecs"
	}
`, name)
}

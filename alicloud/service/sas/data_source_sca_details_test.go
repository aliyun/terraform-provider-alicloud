package sas_test

import (
	"fmt"
	"os"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudThreatDetectionScaDetailsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// search_item, search_info filter, matching
			{
				Config: testAccScaDetailsConfig(name, fmt.Sprintf(`
			search_item = "name"
			search_info = "liblzma"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.type"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type_dispaly"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_started"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.config_path", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.pid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.port", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.cmdline"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_ip", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.version"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_name"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_status", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.web_path", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.ppid", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.intranet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.uuid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_protocol", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.image_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.path"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.container_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.proof"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.create_timestamp"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.runtime_env_version", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pod_name", ""),
				),
			},
			{
				Config: testAccScaDetailsConfig(name, fmt.Sprintf(`
			search_item = "type"
			search_info = "software_library"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.#"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.type", "library"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type_dispaly"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_started"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.config_path", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.pid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.port", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.cmdline"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type", "software_library"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_ip", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.version"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_name"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_status", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.web_path", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.ppid", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.intranet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.uuid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_protocol", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.image_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.path"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.container_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.proof"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.create_timestamp"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.runtime_env_version", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pod_name", ""),
				),
			},
			// search_item_sub, search_info_sub filter, matching
			{
				Config: testAccScaDetailsConfig(name, fmt.Sprintf(`
			search_item_sub  = "pid"
			search_info_sub = "689"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.type"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type_dispaly"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_started"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.internet_ip"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.config_path", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pid", "689"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.port"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.cmdline"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type", "software_library"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.listen_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.version"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.listen_status"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.web_path", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ppid"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.intranet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.uuid"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.listen_protocol"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.image_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.path"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.container_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.proof"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.create_timestamp"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.runtime_env_version", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pod_name", ""),
				),
			},
			// biz filter, matching
			{
				Config: testAccScaDetailsConfig(name, fmt.Sprintf(`
			biz = "sca"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.type"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type_dispaly", "软件库"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_started"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.config_path", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.pid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.port", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.cmdline"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type", "software_library"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_ip", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.version"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_name"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_status", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.web_path", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.ppid", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.intranet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.uuid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_protocol", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.image_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.path"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.container_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.proof"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.create_timestamp"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.runtime_env_version", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pod_name", ""),
				),
			},
			// both filters, matching
			{
				Config: testAccScaDetailsConfig(name, fmt.Sprintf(`
				search_item     = "name"
				search_info     = "rsyslogd"
				search_item_sub = "pid"
				search_info_sub = "770"
				biz             = "sca"
				lang            = "en"`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.type"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type_dispaly", "DevOps Application"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_started"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.internet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.config_path"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pid", "770"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.port", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.cmdline"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.biz_type", "devops_application"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.version", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_name"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_status", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.name", "rsyslogd"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.process_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.web_path", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.ppid"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.intranet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.uuid"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.listen_protocol", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.image_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.path"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.container_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.proof"),
					resource.TestCheckResourceAttrSet("data.alicloud_threat_detection_sca_details.default", "details.0.create_timestamp"),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.runtime_env_version", ""),
					resource.TestCheckResourceAttr("data.alicloud_threat_detection_sca_details.default", "details.0.pod_name", ""),
				),
			},
		},
	})
}

func testAccScaDetailsConfig(name string, filter string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_threat_detection_sca_details" "default" {
  %s
}
`, name, filter)
}

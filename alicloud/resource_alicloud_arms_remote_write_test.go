package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudArmsRemoteWrite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_remote_write.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsRemoteWriteMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsRemoteWrite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10, 99)
	name := fmt.Sprintf("tf-testacc-ArmsRW%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsRemoteWriteBasicDependence0)
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
					"cluster_id":        "${alicloud_arms_prometheus.default.id}",
					"remote_write_yaml": `remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 666, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":        CHECKSET,
						"remote_write_yaml": "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 666, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_write_yaml": `remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 888, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_write_yaml": "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 888, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudArmsRemoteWriteMap = map[string]string{
	"remote_write_name": CHECKSET,
}

func AlicloudArmsRemoteWriteBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_arms_prometheus" "default" {
  		cluster_type        = "ecs"
  		grafana_instance_id = "free"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		security_group_id   = alicloud_security_group.default.id
  		cluster_name        = "${var.name}-${data.alicloud_vpcs.default.ids.0}"
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}
`, name)
}

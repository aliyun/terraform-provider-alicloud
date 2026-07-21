package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// resource alicloud_arms_remote_writes has been deprecated from version 1.228.0
func SkipTestAccAliCloudArmsRemoteWrite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_remote_write.default"
	ra := resourceAttrInit(resourceId, AliCloudArmsRemoteWriteMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsRemoteWrite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10, 99)
	name := fmt.Sprintf("tf-testacc-ArmsRW%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudArmsRemoteWriteBasicDependence0)
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
						"remote_write_yaml": "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth:\n    username: 666\n    password: '******'\n  write_relabel_configs:\n  - source_labels:\n    - instance_id\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_write_yaml": `remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 888, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_write_yaml": "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth:\n    username: 888\n    password: '******'\n  write_relabel_configs:\n  - source_labels:\n    - instance_id\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n",
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

var AliCloudArmsRemoteWriteMap = map[string]string{
	"remote_write_name": CHECKSET,
}

func AliCloudArmsRemoteWriteBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	
	data "alicloud_resource_manager_resource_groups" "default" {}
	
	resource "alicloud_vpc" "default" {
	  vpc_name = var.name
	}
	
	resource "alicloud_vswitch" "vswitch" {
	  vpc_id       = alicloud_vpc.default.id
	  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
	  zone_id      = data.alicloud_zones.default.zones.1.id
	  vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_arms_prometheus" "default" {
  		cluster_type        = "ecs"
  		grafana_instance_id = "free"
  		vpc_id              = alicloud_vpc.default.id
  		vswitch_id          = alicloud_vswitch.vswitch.id
  		security_group_id   = alicloud_security_group.default.id
  		cluster_name        = "${var.name}-${alicloud_vpc.default.id}"
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}
`, name)
}

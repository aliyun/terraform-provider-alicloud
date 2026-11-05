package slb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudSLBListenerAclAttachment_basic0 covers the HTTP protocol path:
// create, in-place ACL replacement, import and removal of the attachment while
// the listener keeps running.
func TestAccAliCloudSLBListenerAclAttachment_basic0(t *testing.T) {
	var (
		name         = fmt.Sprintf("tf-testaccslblisteneracl%d", tfacctest.RandIntRange(10000, 99999))
		resourceId   = "alicloud_slb_listener_acl_attachment.default"
		listenerId   = "alicloud_slb_listener.default"
		dataSourceId = "data.alicloud_slb_listeners.default"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSLBListenerAclAttachmentConfigHttp(name, "first", "white", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "listener_id", listenerId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "acl_id", "alicloud_slb_acl.first", "id"),
					resource.TestCheckResourceAttr(resourceId, "acl_type", "white"),
					resource.TestCheckResourceAttr(resourceId, "acl_status", "on"),
					// Remote reads through the data source prove that binding the
					// ACL does not reset other listener settings.
					resource.TestCheckResourceAttrPair(dataSourceId, "slb_listeners.0.id", listenerId, "id"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.0.bandwidth", "10"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.0.scheduler", "wrr"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.0.health_check", "on"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.0.health_check_uri", "/cons"),
				),
			},
			{
				Config: testAccSLBListenerAclAttachmentConfigHttp(name, "second", "black", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "acl_id", "alicloud_slb_acl.second", "id"),
					resource.TestCheckResourceAttr(resourceId, "acl_type", "black"),
					resource.TestCheckResourceAttr(resourceId, "acl_status", "on"),
					resource.TestCheckResourceAttrPair(dataSourceId, "slb_listeners.0.id", listenerId, "id"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.0.bandwidth", "10"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.0.scheduler", "wrr"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					listener := s.RootModule().Resources[listenerId]
					acl := s.RootModule().Resources["alicloud_slb_acl.second"]
					if listener == nil || acl == nil {
						return "", fmt.Errorf("cannot find listener or ACL state for import")
					}
					return listener.Primary.ID + ":" + acl.Primary.ID, nil
				},
			},
			{
				// Removing the attachment disables access control but must keep
				// the listener itself alive.
				Config: testAccSLBListenerAclAttachmentConfigHttp(name, "second", "black", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceId, "slb_listeners.0.id", listenerId, "id"),
					resource.TestCheckResourceAttr(dataSourceId, "slb_listeners.#", "1"),
				),
			},
		},
	})
}

// TestAccAliCloudSLBListenerAclAttachment_tcp covers the TCP protocol dispatch
// path, which maps to a different Set/Describe action pair.
func TestAccAliCloudSLBListenerAclAttachment_tcp(t *testing.T) {
	var (
		name       = fmt.Sprintf("tf-testaccslblisteneracltcp%d", tfacctest.RandIntRange(10000, 99999))
		resourceId = "alicloud_slb_listener_acl_attachment.default"
		listenerId = "alicloud_slb_listener.default"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSLBListenerAclAttachmentConfigTcp(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "listener_id", listenerId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "acl_id", "alicloud_slb_acl.first", "id"),
					resource.TestCheckResourceAttr(resourceId, "acl_type", "white"),
					resource.TestCheckResourceAttr(resourceId, "acl_status", "on"),
				),
			},
		},
	})
}

func testAccSLBListenerAclAttachmentConfigHttp(name, aclName, aclType string, withAttachment bool) string {
	attachment := ""
	if withAttachment {
		attachment = fmt.Sprintf(`
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.%s.id
  acl_type    = "%s"
}
`, aclName, aclType)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name   = var.name
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  load_balancer_spec   = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  bandwidth                 = 10
  scheduler                 = "wrr"
  health_check              = "on"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx,http_3xx"

  # The ACL binding is managed by alicloud_slb_listener_acl_attachment, so this
  # listener (which is only used to provision the listener itself) ignores the
  # ACL drift on purpose. This mirrors the documented management boundary.
  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}

resource "alicloud_slb_acl" "first" {
  name       = "${var.name}-first"
  ip_version = "ipv4"
  entry_list {
    entry = "10.10.10.0/24"
  }
}

resource "alicloud_slb_acl" "second" {
  name       = "${var.name}-second"
  ip_version = "ipv4"
  entry_list {
    entry = "10.10.11.0/24"
  }
}
%s
data "alicloud_slb_listeners" "default" {
  load_balancer_id = alicloud_slb_listener.default.load_balancer_id
}
`, name, attachment)
}

func testAccSLBListenerAclAttachmentConfigTcp(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name   = var.name
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  load_balancer_spec   = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 22
  frontend_port             = 22
  protocol                  = "tcp"
  bandwidth                 = 10
  health_check_type         = "tcp"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}

resource "alicloud_slb_acl" "first" {
  name       = "${var.name}-first"
  ip_version = "ipv4"
  entry_list {
    entry = "10.10.10.0/24"
  }
}

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.first.id
  acl_type    = "white"
}
`, name)
}

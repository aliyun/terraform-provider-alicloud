package ens_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test Ens LoadBalancerHttpsListener. >>> Resource test cases.
//
// The HTTPS listener requires an ENS server certificate. ENS does not expose a
// server certificate upload API, so the certificate id must be supplied via the
// ALICLOUD_ENS_SERVER_CERTIFICATE_ID environment variable; the test is skipped
// when it is unset, and is documented as a regression exemption in the PR body.

const ensLoadBalancerHttpsListenerResourceId = "alicloud_ens_load_balancer_https_listener.default"

// Case: basic - create with health check off, update to on + running, import.
func TestAccAliCloudEnsLoadBalancerHttpsListener_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	certId := os.Getenv("ALICLOUD_ENS_SERVER_CERTIFICATE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "ALICLOUD_ENS_SERVER_CERTIFICATE_ID")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEnsLoadBalancerHttpsListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnsLoadBalancerHttpsListenerBasicConfig(name, certId, map[string]string{
					"listener_port":             "443",
					"health_check":              "off",
					"health_check_uri":          "/health",
					"health_check_connect_port": "443",
					"health_check_http_code":    "http_2xx",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "listener_port", "443"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "server_certificate_id", certId),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check", "off"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "scheduler", "wrr"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "idle_timeout", "15"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "request_timeout", "60"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "listener_forward", "off"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "status", "stopped"),
				),
			},
			{
				Config: testAccEnsLoadBalancerHttpsListenerBasicConfig(name, certId, map[string]string{
					"description":               name,
					"scheduler":                 "wlc",
					"health_check":              "on",
					"health_check_uri":          "/status",
					"health_check_domain":       "127.0.0.1",
					"healthy_threshold":         "5",
					"unhealthy_threshold":       "5",
					"health_check_timeout":      "10",
					"health_check_connect_port": "81",
					"health_check_interval":     "5",
					"health_check_http_code":    "http_3xx",
					"health_check_method":       "get",
					"idle_timeout":              "30",
					"request_timeout":           "90",
					"status":                    "running",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "description", name),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "scheduler", "wlc"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check", "on"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check_uri", "/status"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check_domain", "127.0.0.1"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "healthy_threshold", "5"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "unhealthy_threshold", "5"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check_timeout", "10"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check_connect_port", "81"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check_interval", "5"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "health_check_method", "get"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "idle_timeout", "30"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "request_timeout", "90"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "status", "running"),
				),
			},
			{
				ResourceName:      ensLoadBalancerHttpsListenerResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Write-only sticky session attributes: the ENS Describe API
					// does not return them, so they cannot be re-imported.
					"sticky_session_type", "cookie_timeout", "cookie",
				},
			},
		},
	})
}

// Case: sticky - sticky_session_type is ForceNew; recreate on its change.
func TestAccAliCloudEnsLoadBalancerHttpsListener_sticky(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	certId := os.Getenv("ALICLOUD_ENS_SERVER_CERTIFICATE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "ALICLOUD_ENS_SERVER_CERTIFICATE_ID")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEnsLoadBalancerHttpsListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnsLoadBalancerHttpsListenerBasicConfig(name, certId, map[string]string{
					"listener_port":       "8443",
					"backend_server_port": "80",
					"sticky_session_type": "insert",
					"cookie_timeout":      "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "listener_port", "8443"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "backend_server_port", "80"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "sticky_session_type", "insert"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "cookie_timeout", "100"),
				),
			},
			{
				Config: testAccEnsLoadBalancerHttpsListenerBasicConfig(name, certId, map[string]string{
					"listener_port":       "8443",
					"backend_server_port": "80",
					"sticky_session_type": "server",
					"cookie":              "tf-cookie",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "listener_port", "8443"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "backend_server_port", "80"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "sticky_session_type", "server"),
					resource.TestCheckResourceAttr(ensLoadBalancerHttpsListenerResourceId, "cookie", "tf-cookie"),
				),
			},
			{
				ResourceName:      ensLoadBalancerHttpsListenerResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"sticky_session_type", "cookie_timeout", "cookie",
				},
			},
		},
	})
}

func testAccCheckEnsLoadBalancerHttpsListenerDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	ensService := ens.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ens_load_balancer_https_listener" {
			continue
		}
		_, err := ensService.DescribeEnsLoadBalancerHttpsListener(rs.Primary.ID)
		if err != nil {
			if errs.IsProviderNotFound(err) || errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing ENS LoadBalancerHttpsListener %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("ENS LoadBalancerHttpsListener %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccEnsLoadBalancerHttpsListenerConfig renders the HTTPS listener config.
// attributes carries the per-step attribute overrides; only the supplied keys
// are emitted, so a step can omit optional attributes to exercise defaults.
func testAccEnsLoadBalancerHttpsListenerBasicConfig(name, certId string, attributes map[string]string) string {
	attrs := ""
	for k, v := range attributes {
		attrs += fmt.Sprintf("  %s = %q\n", k, v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_network" "network" {
  network_name   = "HttpsListenerNetwork_autotest"
  description    = var.name
  cidr_block     = "192.168.0.0/16"
  ens_region_id  = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description    = "HttpsListenerVSwitch_autotest"
  cidr_block     = "192.168.2.0/24"
  vswitch_name   = var.name
  ens_region_id  = "cn-chenzhou-telecom_unicom_cmcc"
  network_id     = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer" "default" {
  payment_type       = "PayAsYouGo"
  ens_region_id      = "cn-chenzhou-telecom_unicom_cmcc"
  load_balancer_spec = "elb.s1.small"
  load_balancer_name = var.name
  vswitch_id         = alicloud_ens_vswitch.switch.id
  network_id         = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer_https_listener" "default" {
  load_balancer_id     = alicloud_ens_load_balancer.default.id
  server_certificate_id = "%s"
%s
}
`, name, certId, attrs)
}

// Test Ens LoadBalancerHttpsListener. <<< Resource test cases.

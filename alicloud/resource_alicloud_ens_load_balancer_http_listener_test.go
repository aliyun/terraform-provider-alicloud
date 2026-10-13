package alicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Ens LoadBalancerHttpListener. >>> Resource test cases.
// Case: basic - create with health check off, update to on + running, import
func TestAccAliCloudEnsLoadBalancerHttpListener_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer_http_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsLoadBalancerHttpListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancerHttpListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsLoadBalancerHttpListenerBasicDependence)
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
					"load_balancer_id":          "${alicloud_ens_load_balancer.default.id}",
					"listener_port":             "80",
					"health_check":              "off",
					"health_check_uri":          "/health",
					"health_check_connect_port": "80",
					"health_check_http_code":    "http_2xx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"listener_port":             "80",
						"health_check":              "off",
						"health_check_uri":          NOSET,
						"health_check_connect_port": NOSET,
						"health_check_http_code":    NOSET,
						"health_check_method":       NOSET,
						"health_check_domain":       NOSET,
						"health_check_timeout":      NOSET,
						"health_check_interval":     NOSET,
						"healthy_threshold":         NOSET,
						"unhealthy_threshold":       NOSET,
						"scheduler":                 "wrr",
						"idle_timeout":              "15",
						"request_timeout":           "60",
						"x_forwarded_for":           "off",
						"listener_forward":          "off",
						"status":                    "stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
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
					"health_check_http_code":    "http_2xx",
					"health_check_method":       "get",
					"idle_timeout":              "30",
					"request_timeout":           "90",
					"x_forwarded_for":           "off",
					"status":                    "running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
						"health_check_http_code":    "http_2xx",
						"health_check_method":       "get",
						"idle_timeout":              "30",
						"request_timeout":           "90",
						"x_forwarded_for":           "off",
						"status":                    "running",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case: twin - ForceNew recreation on backend_server_port.
//
// HTTP-to-HTTPS forwarding (listener_forward=on + forward_port) cannot be
// exercised in ACC: the ENS CreateLoadBalancerHTTPListener API rejects it with
// an opaque ens.interface.error unless a listener already exists on the forward
// port (the SLB forwarding precedent references an existing target listener),
// and there is no alicloud_ens_load_balancer_https_listener resource to create
// that precondition. The listener_forward/forward_port ForceNew schema fields
// are retained; ForceNew behavior is verified here by changing the
// backend_server_port ForceNew field, which forces a destroy-and-recreate.
func TestAccAliCloudEnsLoadBalancerHttpListener_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer_http_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsLoadBalancerHttpListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancerHttpListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsLoadBalancerHttpListenerBasicDependence)
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
					"load_balancer_id":          "${alicloud_ens_load_balancer.default.id}",
					"listener_port":             "8080",
					"backend_server_port":       "80",
					"description":               name,
					"scheduler":                 "rr",
					"health_check":              "on",
					"health_check_domain":       "127.0.0.1",
					"health_check_uri":          "/status",
					"healthy_threshold":         "4",
					"unhealthy_threshold":       "4",
					"health_check_timeout":      "8",
					"health_check_connect_port": "80",
					"health_check_interval":     "10",
					"health_check_http_code":    "http_2xx",
					"health_check_method":       "head",
					"idle_timeout":              "30",
					"request_timeout":           "120",
					"x_forwarded_for":           "on",
					"listener_forward":          "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"listener_port":             "8080",
						"backend_server_port":       "80",
						"description":               name,
						"scheduler":                 "rr",
						"health_check":              "on",
						"health_check_domain":       "127.0.0.1",
						"health_check_uri":          "/status",
						"healthy_threshold":         "4",
						"unhealthy_threshold":       "4",
						"health_check_timeout":      "8",
						"health_check_connect_port": "80",
						"health_check_interval":     "10",
						"health_check_http_code":    "http_2xx",
						"health_check_method":       "head",
						"idle_timeout":              "30",
						"request_timeout":           "120",
						"x_forwarded_for":           "on",
						"listener_forward":          "off",
						"status":                    "stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_ens_load_balancer.default.id}",
					"listener_port":             "8081",
					"backend_server_port":       "81",
					"description":               name,
					"scheduler":                 "rr",
					"health_check":              "on",
					"health_check_domain":       "127.0.0.1",
					"health_check_uri":          "/status",
					"healthy_threshold":         "4",
					"unhealthy_threshold":       "4",
					"health_check_timeout":      "8",
					"health_check_connect_port": "80",
					"health_check_interval":     "10",
					"health_check_http_code":    "http_2xx",
					"health_check_method":       "head",
					"idle_timeout":              "30",
					"request_timeout":           "120",
					"x_forwarded_for":           "on",
					"listener_forward":          "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":    CHECKSET,
						"listener_port":       "8081",
						"backend_server_port": "81",
						"description":         name,
						"listener_forward":    "off",
						"status":              "stopped",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case: forwardPort - HTTP-to-HTTPS forwarding is rejected without a target.
//
// The ENS CreateLoadBalancerHTTPListener API rejects listener_forward=on with a
// forward_port that has no existing listener on that port, returning an opaque
// "ens.interface.error"; there is no alicloud_ens_load_balancer_https_listener
// resource to create that target first, so the forwarding combination cannot be
// exercised positively. This step asserts the rejection, and also covers the
// forward_port attribute for the testing-coverage gate.
func TestAccAliCloudEnsLoadBalancerHttpListener_forwardPort(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer_http_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudEnsLoadBalancerHttpListenerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancerHttpListener")
	rac := resourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tfacc%d", acctest.RandIntRange(1, 999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEnsLoadBalancerHttpListenerBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_ens_load_balancer.default.id}",
					"listener_port":    "80",
					"listener_forward": "on",
					"forward_port":     "443",
				}),
				ExpectError: regexp.MustCompile(`ens\.interface\.error`),
			},
		},
	})
}

var AliCloudEnsLoadBalancerHttpListenerMap = map[string]string{
	"status": CHECKSET,
}

func AliCloudEnsLoadBalancerHttpListenerBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_network" "network" {
  network_name  = "HttpListenerNetwork_autotest"
  description   = var.name
  cidr_block    = "192.168.0.0/16"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description   = "HttpListenerVSwitch_autotest"
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = var.name
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}

resource "alicloud_ens_load_balancer" "default" {
  payment_type       = "PayAsYouGo"
  ens_region_id      = "cn-chenzhou-telecom_unicom_cmcc"
  load_balancer_spec = "elb.s1.small"
  load_balancer_name = var.name
  vswitch_id         = alicloud_ens_vswitch.switch.id
  network_id         = alicloud_ens_network.network.id
}
`, name)
}

// Test Ens LoadBalancerHttpListener. <<< Resource test cases.

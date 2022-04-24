package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSLBListener_http_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbListenerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"backend_port":              "80",
					"frontend_port":             "80",
					"protocol":                  "http",
					"bandwidth":                 "10",
					"sticky_session":            "on",
					"sticky_session_type":       "insert",
					"cookie_timeout":            "86400",
					"cookie":                    "testslblistenercookie",
					"health_check":              "on",
					"health_check_domain":       "ali.com",
					"health_check_uri":          "/cons",
					"health_check_connect_port": "20",
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "5",
					"health_check_http_code":    "http_2xx,http_3xx",
					"x_forwarded_for": []map[string]interface{}{
						{
							"retrive_slb_ip": "true",
							"retrive_slb_id": "true",
						},
					},
					"acl_status":      "on",
					"acl_type":        "white",
					"acl_id":          "${alicloud_slb_acl.default.id}",
					"request_timeout": "80",
					"idle_timeout":    "30",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":                    CHECKSET,
						"backend_port":                        "80",
						"frontend_port":                       "80",
						"protocol":                            "http",
						"bandwidth":                           "10",
						"scheduler":                           string(WRRScheduler),
						"sticky_session":                      string(OnFlag),
						"sticky_session_type":                 string(InsertStickySessionType),
						"cookie_timeout":                      "86400",
						"health_check":                        "on",
						"health_check_uri":                    "/cons",
						"health_check_domain":                 "ali.com",
						"health_check_connect_port":           "20",
						"healthy_threshold":                   "8",
						"unhealthy_threshold":                 "8",
						"health_check_timeout":                "8",
						"health_check_interval":               "5",
						"health_check_http_code":              string(HTTP_2XX) + "," + string(HTTP_3XX),
						"x_forwarded_for.0.retrive_client_ip": "true",
						"x_forwarded_for.0.retrive_slb_ip":    "true",
						"x_forwarded_for.0.retrive_slb_id":    "true",
						"x_forwarded_for.0.retrive_slb_proto": "false",
						"acl_status":                          "on",
						"acl_type":                            string(AclTypeWhite),
						"acl_id":                              CHECKSET,
						"gzip":                                "true",
						"idle_timeout":                        "30",
						"request_timeout":                     "80",
						"description":                         name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": string(WLCScheduler),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": string(WLCScheduler),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cookie_timeout": "80000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie_timeout": "80000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_uri": "/con",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/con",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_domain": "al.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_domain": "al.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gzip": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"idle_timeout": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_timeout": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_status": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_status": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sticky_session_type": string(ServerStickySessionType),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_type": string(ServerStickySessionType),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": string(RRScheduler),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": string(RRScheduler),
					}),
				),
			},
		},
	})
}

func TestAccCheckSlbListenerForward(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerForwardConfigSpot%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbListenerHttpForward(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"frontend_port":             "80",
						"protocol":                  "http",
						"listener_forward":          "on",
						"forward_port":              "443",
						"bandwidth":                 NOSET,
						"scheduler":                 string(WRRScheduler),
						"sticky_session":            NOSET,
						"sticky_session_type":       NOSET,
						"cookie_timeout":            NOSET,
						"health_check":              NOSET,
						"health_check_uri":          NOSET,
						"health_check_domain":       NOSET,
						"health_check_connect_port": NOSET,
						"healthy_threshold":         "3",
						"unhealthy_threshold":       "3",
						"health_check_timeout":      "5",
						"health_check_interval":     "2",
						"acl_status":                string(OffFlag),
						"acl_type":                  NOSET,
						"acl_id":                    NOSET,
						"gzip":                      NOSET,
						"idle_timeout":              NOSET,
						"request_timeout":           NOSET,
						"health_check_http_code":    NOSET,
					}),
				),
			},
		},
	})
}
func TestAccAlicloudSLBListener_same_port(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerSamePort%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbListenerSamePort(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"frontend_port":    "80",
						"protocol":         "tcp",
						"bandwidth":        "10",
						"backend_port":     "80",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudSLBListener_https_update(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbHTTPSListenerConfigDependence)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"backend_port":              "80",
					"frontend_port":             "80",
					"protocol":                  "https",
					"bandwidth":                 "10",
					"sticky_session":            "on",
					"sticky_session_type":       "insert",
					"cookie_timeout":            "86400",
					"cookie":                    "testslblistenercookie",
					"health_check":              "on",
					"health_check_uri":          "/cons",
					"health_check_domain":       "internal-health-check",
					"health_check_connect_port": "20",
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "5",
					"health_check_http_code":    "http_2xx,http_3xx",
					"x_forwarded_for": []map[string]interface{}{
						{
							"retrive_slb_ip": "true",
							"retrive_slb_id": "true",
						},
					},
					"acl_status":            "on",
					"acl_type":              "white",
					"acl_id":                "${alicloud_slb_acl.default.id}",
					"request_timeout":       "80",
					"idle_timeout":          "30",
					"enable_http2":          "on",
					"tls_cipher_policy":     "tls_cipher_policy_1_2",
					"server_certificate_id": "${alicloud_slb_server_certificate.default.id}",
					"ca_certificate_id":     "${alicloud_slb_ca_certificate.default.id}",
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"frontend_port":             "80",
						"backend_port":              "80",
						"protocol":                  "https",
						"bandwidth":                 "10",
						"scheduler":                 string(WRRScheduler),
						"sticky_session":            string(OnFlag),
						"sticky_session_type":       string(InsertStickySessionType),
						"cookie_timeout":            "86400",
						"health_check":              "on",
						"health_check_connect_port": "20",
						"health_check_domain":       "internal-health-check",
						"healthy_threshold":         "8",
						"unhealthy_threshold":       "8",
						"health_check_timeout":      "8",
						"health_check_interval":     "5",
						"acl_status":                "on",
						"acl_type":                  string(AclTypeWhite),
						"acl_id":                    CHECKSET,
						"gzip":                      "true",
						"idle_timeout":              "30",
						"request_timeout":           "80",
						"health_check_http_code":    string(HTTP_2XX) + "," + string(HTTP_3XX),
						"server_certificate_id":     CHECKSET,
						"ca_certificate_id":         CHECKSET,
						"enable_http2":              "on",
						"x_forwarded_for.#":         "1",
						"tls_cipher_policy":         "tls_cipher_policy_1_2",
						"description":               name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_cipher_policy": "tls_cipher_policy_1_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_cipher_policy": "tls_cipher_policy_1_1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": string(WLCScheduler),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": string(WLCScheduler),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cookie_timeout": "80000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie_timeout": "80000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_uri": "/con",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/con",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gzip": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"idle_timeout": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_timeout": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_status": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_status":          "off",
						"health_check_domain": "internal-health-check",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBListener_tcp_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbListenerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"frontend_port":             "22",
					"backend_port":              "22",
					"protocol":                  "tcp",
					"scheduler":                 string(WRRScheduler),
					"bandwidth":                 "10",
					"acl_status":                "on",
					"acl_type":                  string(AclTypeBlack),
					"acl_id":                    "${alicloud_slb_acl.default.id}",
					"persistence_timeout":       "3600",
					"health_check_type":         string(HTTPHealthCheckType),
					"health_check_domain":       "",
					"health_check_uri":          "/console",
					"health_check_connect_port": "20",
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "5",
					"health_check_http_code":    string(HTTP_2XX),
					"established_timeout":       "600",
					"description":               name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"frontend_port":             "22",
						"backend_port":              "22",
						"protocol":                  "tcp",
						"scheduler":                 string(WRRScheduler),
						"bandwidth":                 "10",
						"acl_status":                "on",
						"acl_type":                  string(AclTypeBlack),
						"acl_id":                    CHECKSET,
						"persistence_timeout":       "3600",
						"health_check_type":         string(HTTPHealthCheckType),
						"health_check_domain":       "",
						"health_check_uri":          "/console",
						"health_check_connect_port": "20",
						"healthy_threshold":         "8",
						"unhealthy_threshold":       "8",
						"health_check_timeout":      "8",
						"health_check_interval":     "5",
						"health_check_http_code":    string(HTTP_2XX),
						"established_timeout":       "600",
						"description":               name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"persistence_timeout": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistence_timeout": "3000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_uri": "/cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_http_code": string(HTTP_2XX) + "," + string(HTTP_3XX),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_code": string(HTTP_2XX) + "," + string(HTTP_3XX),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_type": string(TCPHealthCheckType),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_type":      string(TCPHealthCheckType),
						"health_check_http_code": "",
						"health_check_uri":       "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"established_timeout": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"established_timeout": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "4",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBListener_tcp_server_group(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbListenerServerGroupConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"frontend_port":             "22",
					"protocol":                  "tcp",
					"backend_port":              "22",
					"bandwidth":                 "10",
					"persistence_timeout":       "3600",
					"health_check_type":         "http",
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "5",
					"health_check_http_code":    "http_2xx",
					"health_check_connect_port": "20",
					"health_check_uri":          "/console",
					"established_timeout":       "600",
					"server_group_id":           "${alicloud_slb_server_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"frontend_port":             "22",
						"backend_port":              "22",
						"protocol":                  "tcp",
						"bandwidth":                 "10",
						"persistence_timeout":       "3600",
						"health_check_type":         "http",
						"healthy_threshold":         "8",
						"unhealthy_threshold":       "8",
						"health_check_timeout":      "8",
						"health_check_interval":     "5",
						"health_check_http_code":    "http_2xx",
						"health_check_connect_port": "20",
						"health_check_uri":          "/console",
						"established_timeout":       "600",
						"server_group_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_group_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_slave_server_group_id": "${alicloud_slb_master_slave_server_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_slave_server_group_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBListener_udp_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbListenerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"backend_port":              "2001",
					"frontend_port":             "2001",
					"protocol":                  "udp",
					"bandwidth":                 "10",
					"scheduler":                 string(WRRScheduler),
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "4",
					"persistence_timeout":       "3600",
					"health_check_connect_port": "20",
					"acl_status":                "on",
					"acl_type":                  string(AclTypeBlack),
					"acl_id":                    "${alicloud_slb_acl.default.id}",
					"description":               name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(
						map[string]string{
							"load_balancer_id":          CHECKSET,
							"backend_port":              "2001",
							"frontend_port":             "2001",
							"protocol":                  "udp",
							"bandwidth":                 "10",
							"scheduler":                 string(WRRScheduler),
							"healthy_threshold":         "8",
							"unhealthy_threshold":       "8",
							"health_check_timeout":      "8",
							"health_check_interval":     "4",
							"persistence_timeout":       "3600",
							"health_check_connect_port": "20",
							"acl_status":                "on",
							"acl_type":                  string(AclTypeBlack),
							"acl_id":                    CHECKSET,
							"description":               name,
						}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_connect_port": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"unhealthy_threshold": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_timeout": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"persistence_timeout": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistence_timeout": "3000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "15",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBListener_http_healcheckmethod(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbListenerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.HttpHttpsHealthCheckMehtodSupportedRegions)
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"backend_port":              "80",
					"frontend_port":             "80",
					"protocol":                  "http",
					"bandwidth":                 "10",
					"sticky_session":            "on",
					"sticky_session_type":       "insert",
					"cookie_timeout":            "86400",
					"cookie":                    "testslblistenercookie",
					"health_check":              "on",
					"health_check_domain":       "ali.com",
					"health_check_method":       "head",
					"health_check_uri":          "/cons",
					"health_check_connect_port": "20",
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "5",
					"health_check_http_code":    "http_2xx,http_3xx",
					"x_forwarded_for": []map[string]interface{}{
						{
							"retrive_slb_ip": "true",
							"retrive_slb_id": "true",
						},
					},
					"acl_status":      "on",
					"acl_type":        "white",
					"acl_id":          "${alicloud_slb_acl.default.id}",
					"request_timeout": "80",
					"idle_timeout":    "30",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":                    CHECKSET,
						"backend_port":                        "80",
						"frontend_port":                       "80",
						"protocol":                            "http",
						"bandwidth":                           "10",
						"scheduler":                           string(WRRScheduler),
						"sticky_session":                      string(OnFlag),
						"sticky_session_type":                 string(InsertStickySessionType),
						"cookie_timeout":                      "86400",
						"health_check":                        "on",
						"health_check_uri":                    "/cons",
						"health_check_domain":                 "ali.com",
						"health_check_method":                 "head",
						"health_check_connect_port":           "20",
						"healthy_threshold":                   "8",
						"unhealthy_threshold":                 "8",
						"health_check_timeout":                "8",
						"health_check_interval":               "5",
						"health_check_http_code":              string(HTTP_2XX) + "," + string(HTTP_3XX),
						"x_forwarded_for.0.retrive_client_ip": "true",
						"x_forwarded_for.0.retrive_slb_ip":    "true",
						"x_forwarded_for.0.retrive_slb_id":    "true",
						"x_forwarded_for.0.retrive_slb_proto": "false",
						"acl_status":                          "on",
						"acl_type":                            string(AclTypeWhite),
						"acl_id":                              CHECKSET,
						"gzip":                                "true",
						"idle_timeout":                        "30",
						"request_timeout":                     "80",
						"description":                         name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_method": "get",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_method": "get",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudSLBListener_https_healcheckmethod(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccSlbListenerConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbHTTPSListenerConfigDependence)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.HttpHttpsHealthCheckMehtodSupportedRegions)
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id":          "${alicloud_slb_load_balancer.default.id}",
					"backend_port":              "80",
					"frontend_port":             "80",
					"protocol":                  "https",
					"bandwidth":                 "10",
					"sticky_session":            "on",
					"sticky_session_type":       "insert",
					"cookie_timeout":            "86400",
					"cookie":                    "testslblistenercookie",
					"health_check":              "on",
					"health_check_uri":          "/cons",
					"health_check_method":       "head",
					"health_check_connect_port": "20",
					"healthy_threshold":         "8",
					"unhealthy_threshold":       "8",
					"health_check_timeout":      "8",
					"health_check_interval":     "5",
					"health_check_http_code":    "http_2xx,http_3xx",
					"x_forwarded_for": []map[string]interface{}{
						{
							"retrive_slb_ip": "true",
							"retrive_slb_id": "true",
						},
					},
					"acl_status":            "on",
					"acl_type":              "white",
					"acl_id":                "${alicloud_slb_acl.default.id}",
					"request_timeout":       "80",
					"idle_timeout":          "30",
					"enable_http2":          "on",
					"tls_cipher_policy":     "tls_cipher_policy_1_2",
					"server_certificate_id": "${alicloud_slb_server_certificate.default.id}",
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":          CHECKSET,
						"frontend_port":             "80",
						"backend_port":              "80",
						"protocol":                  "https",
						"bandwidth":                 "10",
						"scheduler":                 string(WRRScheduler),
						"sticky_session":            string(OnFlag),
						"sticky_session_type":       string(InsertStickySessionType),
						"cookie_timeout":            "86400",
						"health_check":              "on",
						"health_check_method":       "head",
						"health_check_connect_port": "20",
						"healthy_threshold":         "8",
						"unhealthy_threshold":       "8",
						"health_check_timeout":      "8",
						"health_check_interval":     "5",
						"acl_status":                "on",
						"acl_type":                  string(AclTypeWhite),
						"acl_id":                    CHECKSET,
						"gzip":                      "true",
						"idle_timeout":              "30",
						"request_timeout":           "80",
						"health_check_http_code":    string(HTTP_2XX) + "," + string(HTTP_3XX),
						"server_certificate_id":     CHECKSET,
						"enable_http2":              "on",
						"x_forwarded_for.#":         "1",
						"tls_cipher_policy":         "tls_cipher_policy_1_2",
						"description":               name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_method": "get",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_method": "get",
					}),
				),
			},
		},
	})
}

func testAccSlbListenerHttpForward(name string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_slb_listener" "default"{
  		load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  		frontend_port = 80
  		protocol = "http"
  		listener_forward = "on"
  		forward_port = "${alicloud_slb_listener.default-1.frontend_port}"
	}
	resource "alicloud_slb_listener" "default-1" {
  		load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  		backend_port = 80
  		frontend_port = 443
  		protocol = "https"
  		sticky_session = "off"
  		health_check = "on"
  		health_check_uri = "/cons"
  		health_check_connect_port = 20
  		healthy_threshold = 8
  		unhealthy_threshold = 8
  		health_check_timeout = 8
  		health_check_interval = 5
  		health_check_http_code = "http_2xx,http_3xx"
  		bandwidth = 10
  		server_certificate_id= "${alicloud_slb_server_certificate.default.id}"
	}
	`, resourceSlbHTTPSListenerConfigDependence(name))
}

func testAccSlbListenerSamePort(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
  		default = "%s"
	}
	resource "alicloud_slb_listener" "default"{
  		load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  		frontend_port = 80
  		protocol = "tcp"
		bandwidth = "10"
		backend_port = 80
	}
	resource "alicloud_slb_listener" "default-1" {
  		load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  		frontend_port = 80
  		protocol = "udp"
		bandwidth = "10"
		backend_port = 80
	}`, SlbListenerCommonTestCase, name)
}

func resourceSlbListenerConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, SlbListenerCommonTestCase, name)
}

func resourceSlbListenerServerGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	`, SlbListenerVserverCommonTestCase, name)
}

func resourceSlbHTTPSListenerConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
    variable "name" {
		default = "%s"
	}

    resource "alicloud_slb_ca_certificate" "default" {
	  ca_certificate_name           = "${var.name}"
	  ca_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
    }
    resource "alicloud_slb_server_certificate" "default" {
  		name = "${var.name}"
  		server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  		private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
	}
	`, SlbListenerCommonTestCase, name)
}

func testAccCheckSlbListenerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_listener" {
			continue
		}
		protocol := ""
		port := 0
		var err error
		// Try to find the Slb
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) == 3 {
			protocol = parts[1]
			port, err = strconv.Atoi(parts[2])
		} else {
			port, err = strconv.Atoi(parts[1])
		}
		if err != nil {
			return fmt.Errorf("Parsing SlbListener's id got an error: %#v", err)
		}
		loadBalancer, err := slbService.DescribeSlb(parts[0])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
		}
		if len(parts) == 3 {
			for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
				if portAndProtocol.ListenerPort == port && portAndProtocol.ListenerProtocol == protocol {
					return fmt.Errorf("SLB listener still exist")
				}
			}
		} else {
			for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
				if portAndProtocol.ListenerPort == port {
					return fmt.Errorf("SLB listener still exist")
				}
			}
		}

	}

	return nil
}

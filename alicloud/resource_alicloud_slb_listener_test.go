package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSlbListener_http_basic(t *testing.T) {
	var v map[string]interface{}
	rand := acctest.RandInt()
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbHttpListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerHttpConfig(rand),
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
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateBandwidth(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "15",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateScheduler(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": string(WLCScheduler),
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateCookieTimeout(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie_timeout": "80000",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthCheckUri(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/con",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthCheckDomain(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_domain": "al.com",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthCheckConnectPort(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthThreshold(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateUnHealthThreshold(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthCheckTimeout(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthCheckInterval(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "4",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateGzip(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "false",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateIdleTimeout(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "40",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateRequestTimeout(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "90",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateHealthCheck(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check": "off",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateAclStatus(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_status": "off",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigUpdateStickySessionType(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sticky_session_type": string(ServerStickySessionType),
					}),
				),
			},
			{
				Config: testAccSlbListenerHttpConfigRRScheduler(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":                    CHECKSET,
						"backend_port":                        "80",
						"frontend_port":                       "80",
						"protocol":                            "http",
						"bandwidth":                           "10",
						"scheduler":                           string(RRScheduler),
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
	}, "DescribeSlbHttpListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerHttpForward,
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

func TestAccCheckSlbListener_http_multi(t *testing.T) {
	var v map[string]interface{}
	rand := acctest.RandInt()
	resourceId := "alicloud_slb_listener.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbHttpListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerHttpConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":                    CHECKSET,
						"backend_port":                        "10",
						"frontend_port":                       "10",
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
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_https_update(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbHttpsListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerHttps,
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
						"ssl_certificate_id":        CHECKSET,
						"enable_http2":              "on",
						"x_forwarded_for.#":         "1",
						"tls_cipher_policy":         "tls_cipher_policy_1_2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlbListenerHttps_tls_cipher_policy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_cipher_policy": "tls_cipher_policy_1_1",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_scheduler,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler": string(WLCScheduler),
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_cookie_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie_timeout": "80000",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_health_check_uri,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/con",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_health_check_connect_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_healthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_unhealthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_health_check_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_health_check_interval,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "4",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_gzip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip": "false",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_idle_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "40",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_request_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "90",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_health_check,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check": "off",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps_bandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "15",
					}),
				),
			},
			{
				Config: testAccSlbListenerHttps,
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
						"health_check_uri":          "/cons",
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
						"ssl_certificate_id":        CHECKSET,
						"enable_http2":              "on",
						"x_forwarded_for.#":         "1",
						"tls_cipher_policy":         "tls_cipher_policy_1_2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_tcp_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbTcpListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerTcp,
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
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlbListenerTcp_persistence_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistence_timeout": "3000",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_health_check_uri,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_uri": "/cn",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_health_check_http_code,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_http_code": string(HTTP_2XX) + "," + string(HTTP_3XX),
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_health_check_type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_type":      string(TCPHealthCheckType),
						"health_check_http_code": "",
						"health_check_uri":       "",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_established_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"established_timeout": "500",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_health_check_connect_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_healthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_unhealthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_health_check_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_health_check_interval,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "4",
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp,
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
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_tcp_server_group(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbTcpListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerTcp_server_group,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
						"server_group_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccSlbListenerTcp_server_group_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": "",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbListener_udp_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_listener.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbUdpListener")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccSlbListenerUdp,
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
						}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSlbListenerUdp_health_check_connect_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_connect_port": "30",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp_healthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp_unhealthy_threshold,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"unhealthy_threshold": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp_health_check_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_timeout": "9",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp_health_check_interval,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "5",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp_persistence_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistence_timeout": "3000",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp_bandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "15",
					}),
				),
			},
			{
				Config: testAccSlbListenerUdp,
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
						}),
				),
			},
		},
	})
}

func testAccCheckSlbListenerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_listener" {
			continue
		}

		// Try to find the Slb
		parts := strings.Split(rs.Primary.ID, ":")
		port, err := strconv.Atoi(parts[1])
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
		for _, portAndProtocol := range loadBalancer.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if portAndProtocol.ListenerPort == port {
				return fmt.Errorf("SLB listener still exist")
			}
		}

	}

	return nil
}

func testAccSlbListenerHttpConfig(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 10
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 86400
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
  }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
  }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateBandwidth(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 86400
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
  }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateScheduler(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 86400
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
  }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
  }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateCookieTimeout(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthCheckUri(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/con"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthCheckDomain(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthCheckConnectPort(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthThreshold(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateUnHealthThreshold(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthCheckTimeout(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthCheckInterval(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateGzip(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateIdleTimeout(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 40
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateRequestTimeout(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateHealthCheck(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "off"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateAclStatus(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "off"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  gzip = false
  acl_status = "off"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigUpdateStickySessionType(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 15
  scheduler = "wlc"
  sticky_session = "on"
  sticky_session_type = "server"
  cookie_timeout = 80000
  cookie = "testslblistenercookie"
  health_check = "off"
  health_check_domain = "al.com"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  gzip = false
  acl_status = "off"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

func testAccSlbListenerHttpConfigRRScheduler(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  bandwidth = 10
  scheduler = "rr"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 86400
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

const testAccSlbListenerHttpForward = `
variable "name" {
  default = "tf-testAcc-https-forward"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default"{
  load_balancer_id = "${alicloud_slb.default.id}"
  frontend_port = 80
  protocol = "http"
  listener_forward = "on"
  forward_port = "${alicloud_slb_listener.default-1.frontend_port}"
}
resource "alicloud_slb_listener" "default-1" {
  load_balancer_id = "${alicloud_slb.default.id}"
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
  ssl_certificate_id= "${alicloud_slb_server_certificate.default.id}"
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_2"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_tls_cipher_policy = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_scheduler = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_cookie_timeout = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_health_check_uri = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_health_check_connect_port = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_domain = "al.com"
  health_check_connect_port = 30
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_healthy_threshold = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_unhealthy_threshold = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_health_check_timeout = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_health_check_interval = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_gzip = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_idle_timeout = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 80
  idle_timeout              = 40
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_request_timeout = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "on"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_health_check = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "off"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerHttps_bandwidth = `
variable "name" {
  default = "tf-testAcc-https-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttps"
  internet_charge_type = "PayByTraffic"
  internet = true
  specification = "slb.s1.small"
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "https"
  sticky_session = "on"
  scheduler = "wlc"
  sticky_session_type = "insert"
  cookie = "testslblistenercookie"
  cookie_timeout = 80000
  health_check = "off"
  health_check_uri = "/con"
  health_check_connect_port = 30
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 15
  x_forwarded_for  {
    retrive_slb_ip = false
    retrive_slb_id = false
  }
  gzip = false
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  ssl_certificate_id        = "${alicloud_slb_server_certificate.default.id}"
  request_timeout           = 90
  idle_timeout              = 40
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_1"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
resource "alicloud_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP\nMA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0\nZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow\ndjELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE\nChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG\n9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ\nAoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB\ncoG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook\nKOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw\nHQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy\n+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC\nQkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN\nMAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ\nAJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT\ncQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1\nOfi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd\nDUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV\nkg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM\nywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB\nAoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd\n6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP\nhwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4\nMdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz\n71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm\nEv9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE\nqygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8\n9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM\nzWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe\nDrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}
`

const testAccSlbListenerTcp_server_group = `
data "alicloud_zones" "default" {
  available_disk_category = "cloud_efficiency"
  available_resource_creation= "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count = 1
  memory_size = 2
}

data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
  most_recent = true
  owners = "system"
}

variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    }
   servers {
      server_ids = "${alicloud_instance.default.*.id}"
      port = 80
      weight = 100
    }
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = "${alicloud_slb.default.id}"
  backend_port              = "22"
  frontend_port             = "22"
  protocol                  = "tcp"
  bandwidth                 = "10"
  persistence_timeout       = 3600
  health_check_type         = "http"
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx"
  health_check_connect_port = 20
  health_check_uri          = "/console"
  established_timeout       = 600
  server_group_id           = "${alicloud_slb_server_group.default.id}"
}
`

const testAccSlbListenerTcp = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "http"
  persistence_timeout = 3600
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx"
  health_check_connect_port = 20
  health_check_uri = "/console"
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 600
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_persistence_timeout = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "http"
  persistence_timeout = 3000
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx"
  health_check_timeout = 8
  health_check_connect_port = 20
  health_check_uri = "/console"
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 600
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_health_check_type = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx"
  health_check_connect_port = 20
  health_check_uri = "/console"
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 600
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_health_check_uri = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "http"
  persistence_timeout = 3000
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx"
  health_check_uri = "/cn"
  health_check_connect_port = 20
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 600
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_health_check_http_code = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "http"
  persistence_timeout = 3000
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 20
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 600
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_established_timeout = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 20
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 500
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_health_check_connect_port = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 500
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_healthy_threshold = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 9
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 500
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_unhealthy_threshold = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 500
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_health_check_timeout = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 500
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_health_check_interval = `
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerTcp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "22"
  frontend_port = "22"
  protocol = "tcp"
  bandwidth = "10"
  health_check_type = "tcp"
  persistence_timeout = 3000
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_http_code = "http_2xx,http_3xx"
  health_check_uri = "/cn"
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
  established_timeout = 500
}
variable "name" {
  default = "tf-testAcc-tcp-listener-acl-5"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerTcp_server_group_update = `
data "alicloud_zones" "default" {
  available_disk_category = "cloud_efficiency"
  available_resource_creation= "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count = 1
  memory_size = 2
}

data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
  most_recent = true
  owners = "system"
}

variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "group" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alicloud_security_group.group.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    }
   servers {
      server_ids = "${alicloud_instance.default.*.id}"
      port = 80
      weight = 100
    }
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id          = "${alicloud_slb.default.id}"
  backend_port              = "22"
  frontend_port             = "22"
  protocol                  = "tcp"
  bandwidth                 = "10"
  persistence_timeout       = 3600
  health_check_type         = "http"
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx"
  health_check_connect_port = 20
  health_check_uri          = "/console"
  established_timeout       = 600
}
`
const testAccSlbListenerUdp = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 4
  health_check_connect_port = 20
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_health_check_connect_port = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 4
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_healthy_threshold = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 9
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 4
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_unhealthy_threshold = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 8
  health_check_interval = 4
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_health_check_timeout = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 4
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_health_check_interval = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3600
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 5
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_persistence_timeout = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 10
  persistence_timeout = 3000
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 5
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

const testAccSlbListenerUdp_bandwidth = `
variable "name" {
  default = "tf-testAcc-udp-listener-acl"
}
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerUdp"
  internet_charge_type = "PayByTraffic"
  internet = true
  bandwidth = 20
}
resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = 2001
  frontend_port = 2001
  protocol = "udp"
  bandwidth = 15
  persistence_timeout = 3000
  healthy_threshold = 9
  unhealthy_threshold = 9
  health_check_timeout = 9
  health_check_interval = 5
  health_check_connect_port = 30
  acl_status = "on"
  acl_type   = "black"
  acl_id     = "${alicloud_slb_acl.default.id}"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`

func testAccSlbListenerHttpConfig_multi(rand int) string {
	return fmt.Sprintf(`
variable "number" {
  default = 10
}
variable "name" {
  default = "tf-testAcc-http-listener-acl-%d"
}
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "tf-testAccSlbListenerHttp"
  internet_charge_type = "PayByTraffic"
  internet = true
}
resource "alicloud_slb_listener" "default" {
  count = "${var.number}"
  load_balancer_id = "${alicloud_slb.default.id}"
  backend_port = "${count.index+1}"
  frontend_port = "${count.index+1}"
  protocol = "http"
  bandwidth = 10
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie_timeout = 86400
  cookie = "testslblistenercookie"
  health_check = "on"
  health_check_domain = "ali.com"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  acl_status = "on"
  acl_type   = "white"
  acl_id     = "${alicloud_slb_acl.default.id}"
  request_timeout           = 80
  idle_timeout              = 30
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
    }
   entry_list {
      entry="168.10.10.0/24"
      comment="second"
    }
}
`, rand)
}

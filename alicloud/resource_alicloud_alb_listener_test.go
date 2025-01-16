package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alb_listener",
		&resource.Sweeper{
			Name: "alicloud_alb_listener",
			F:    testSweepAlbListener,
		})
}

func testSweepAlbListener(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListListeners"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeXLarge

	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Listeners", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Listeners", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["ListenerDescription"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["ListenerDescription"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Alb Listener: %s", item["ListenerDescription"].(string))
				continue
			}
			action := "DeleteListener"
			request := map[string]interface{}{
				"ListenerId": item["ListenerId"],
			}
			request["ClientToken"] = buildClientToken("DeleteListener")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alb Listener (%s): %s", item["ListenerId"].(string), err)
			}
			log.Printf("[INFO] Delete Alb Listener success: %s ", item["ListenerId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudALBListener_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBListenerBasicDependence0)
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
					"load_balancer_id":  "${local.load_balancer_id}",
					"listener_protocol": "HTTPS",
					"listener_port":     port,
					"default_actions": []map[string]interface{}{
						{
							"type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id}",
										},
									},
								},
							},
						},
					},
					"certificates": []map[string]interface{}{
						{
							"certificate_id": "${local.certificate_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":  CHECKSET,
						"listener_protocol": "HTTPS",
						"listener_port":     port,
						"default_actions.#": "1",
						"certificates.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"xforwarded_for_config": []map[string]interface{}{
						{
							"xforwardedforclientcertclientverifyalias":   "test_client-verify-alias_123451",
							"xforwardedforclientcertclientverifyenabled": "true",
							"xforwardedforclientcertfingerprintalias":    "test_client-verify-alias_123452",
							"xforwardedforclientcertfingerprintenabled":  "true",
							"xforwardedforclientcert_issuerdnalias":      "test_client-verify-alias_123453",
							"xforwardedforclientcert_issuerdnenabled":    "true",
							"xforwardedforclientcertsubjectdnalias":      "test_client-verify-alias_123454",
							"xforwardedforclientcertsubjectdnenabled":    "true",
							"xforwardedforclientsrcportenabled":          "true",
							"xforwardedforenabled":                       "true",
							"xforwardedforprotoenabled":                  "true",
							"xforwardedforslbidenabled":                  "true",
							"xforwardedforslbportenabled":                "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"xforwarded_for_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http2_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http2_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http2_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http2_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_id": "tls_cipher_policy_1_0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_id": "tls_cipher_policy_1_0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_description": "tf-testAccListener_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_description": "tf-testAccListener_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_timeout": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"request_timeout": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_timeout": "70",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gzip_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gzip_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gzip_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_actions": []map[string]interface{}{
						{
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id_update}",
										},
									},
								},
							},
							"type": "ForwardGroup",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_actions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"idle_timeout": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"idle_timeout": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_log_record_customized_headers_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_log_record_customized_headers_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_log_tracing_config": []map[string]interface{}{
						{
							"tracing_enabled": "true",
							"tracing_sample":  "800",
							"tracing_type":    "Zipkin",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_log_tracing_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http2_enabled": "true",
					"quic_config": []map[string]interface{}{
						{
							"quic_upgrade_enabled": "false",
						},
					},
					"security_policy_id":   "tls_cipher_policy_1_0",
					"listener_description": "tf-testAccListener",
					"request_timeout":      "60",
					"gzip_enabled":         "true",
					"certificates": []map[string]interface{}{
						{
							"certificate_id": "${local.certificate_id_update}",
						},
					},
					"default_actions": []map[string]interface{}{
						{
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id}",
										},
									},
								},
							},
							"type": "ForwardGroup",
						},
					},
					"idle_timeout": "15",
					"xforwarded_for_config": []map[string]interface{}{
						{
							"xforwardedforclientcertclientverifyalias":   "test_client-verify-alias_223451",
							"xforwardedforclientcertclientverifyenabled": "true",
							"xforwardedforclientcertfingerprintalias":    "test_client-verify-alias_223452",
							"xforwardedforclientcertfingerprintenabled":  "true",
							"xforwardedforclientcert_issuerdnalias":      "test_client-verify-alias_223453",
							"xforwardedforclientcert_issuerdnenabled":    "true",
							"xforwardedforclientcertsubjectdnalias":      "test_client-verify-alias_223454",
							"xforwardedforclientcertsubjectdnenabled":    "true",
							"xforwardedforclientsrcportenabled":          "true",
							"xforwardedforenabled":                       "true",
							"xforwardedforprotoenabled":                  "true",
							"xforwardedforslbidenabled":                  "true",
							"xforwardedforslbportenabled":                "true",
						},
					},
					"access_log_record_customized_headers_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http2_enabled":        "true",
						"security_policy_id":   "tls_cipher_policy_1_0",
						"listener_description": "tf-testAccListener",
						"request_timeout":      "60",
						"gzip_enabled":         "true",
						"certificates.#":       "1",
						"default_actions.#":    "1",
						"idle_timeout":         "15",
						"access_log_record_customized_headers_enabled": "false",
						"xforwarded_for_config.#":                      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Listener",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Listener",
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

func TestAccAliCloudALBListener_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBListenerBasicDependence0)
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
					"load_balancer_id":     "${local.load_balancer_id}",
					"listener_protocol":    "HTTPS",
					"listener_port":        port,
					"listener_description": "tf-testAccListener_new",
					"default_actions": []map[string]interface{}{
						{
							"type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id}",
										},
									},
								},
							},
						},
					},
					"certificates": []map[string]interface{}{
						{
							"certificate_id": "${local.certificate_id}",
						},
					},
					"acl_config": []map[string]interface{}{
						{
							"acl_type": "White",
							"acl_relations": []map[string]interface{}{
								{
									"acl_id": "${alicloud_alb_acl.default.0.id}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":     CHECKSET,
						"listener_protocol":    "HTTPS",
						"listener_port":        port,
						"listener_description": "tf-testAccListener_new",
						"default_actions.#":    "1",
						"certificates.#":       "1",
						"acl_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_config": []map[string]interface{}{
						{
							"acl_type": "Black",
							"acl_relations": []map[string]interface{}{
								{
									"acl_id": "${alicloud_alb_acl.default.0.id}",
								},
								{
									"acl_id": "${alicloud_alb_acl.default.1.id}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_config": []map[string]interface{}{
						{
							"acl_type": "White",
							"acl_relations": []map[string]interface{}{
								{
									"acl_id": "${alicloud_alb_acl.default.0.id}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_config.#": "1",
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

func TestAccAliCloudALBListener_basic2(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBListenerBasicDependence0)
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
					"load_balancer_id":     "${local.load_balancer_id}",
					"listener_protocol":    "HTTPS",
					"listener_port":        port,
					"listener_description": "tf-testAccListener_new",
					"default_actions": []map[string]interface{}{
						{
							"type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id}",
										},
									},
								},
							},
						},
					},
					"certificates": []map[string]interface{}{
						{
							"certificate_id": "${local.certificate_id}",
						},
					},
					"acl_config": []map[string]interface{}{
						{
							"acl_type": "White",
							"acl_relations": []map[string]interface{}{
								{
									"acl_id": "${alicloud_alb_acl.default.0.id}",
								},
							},
						},
					},
					"gzip_enabled":       "true",
					"http2_enabled":      "true",
					"idle_timeout":       "20",
					"request_timeout":    "60",
					"security_policy_id": "tls_cipher_policy_1_0",
					"xforwarded_for_config": []map[string]interface{}{
						{
							"xforwardedforclientcertclientverifyalias":   "test_client-verify-alias_223451",
							"xforwardedforclientcertclientverifyenabled": "true",
							"xforwardedforclientcertfingerprintalias":    "test_client-verify-alias_223452",
							"xforwardedforclientcertfingerprintenabled":  "true",
							"xforwardedforclientcert_issuerdnalias":      "test_client-verify-alias_223453",
							"xforwardedforclientcert_issuerdnenabled":    "true",
							"xforwardedforclientcertsubjectdnalias":      "test_client-verify-alias_223454",
							"xforwardedforclientcertsubjectdnenabled":    "true",
							"xforwardedforclientsrcportenabled":          "true",
							"xforwardedforenabled":                       "true",
							"xforwardedforprotoenabled":                  "true",
							"xforwardedforslbidenabled":                  "true",
							"xforwardedforslbportenabled":                "true",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Listener",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":        CHECKSET,
						"listener_protocol":       "HTTPS",
						"listener_port":           port,
						"listener_description":    "tf-testAccListener_new",
						"default_actions.#":       "1",
						"certificates.#":          "1",
						"acl_config.#":            "1",
						"gzip_enabled":            "true",
						"http2_enabled":           "true",
						"idle_timeout":            "20",
						"request_timeout":         "60",
						"security_policy_id":      "tls_cipher_policy_1_0",
						"xforwarded_for_config.#": "1",
						"tags.%":                  "2",
						"tags.Created":            "TF",
						"tags.For":                "Listener",
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

func TestAccAliCloudALBListener_basic3(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 1000)
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBListenerBasicDependence0)
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
					"load_balancer_id":     "${local.load_balancer_id}",
					"listener_protocol":    "HTTPS",
					"listener_port":        port,
					"listener_description": "tf-testAccListener_new",
					"default_actions": []map[string]interface{}{
						{
							"type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id}",
										},
									},
								},
							},
						},
					},
					"certificates": []map[string]interface{}{
						{
							"certificate_id": "${local.certificate_id}",
						},
					},
					"acl_config": []map[string]interface{}{
						{
							"acl_type": "White",
							"acl_relations": []map[string]interface{}{
								{
									"acl_id": "${alicloud_alb_acl.default.0.id}",
								},
							},
						},
					},
					"ca_enabled": "true",
					"ca_certificates": []map[string]interface{}{
						{
							"certificate_id": "${var.ca_certificate_id}",
						},
					},
					"gzip_enabled":       "true",
					"http2_enabled":      "true",
					"idle_timeout":       "20",
					"request_timeout":    "60",
					"security_policy_id": "tls_cipher_policy_1_0",
					"x_forwarded_for_config": []map[string]interface{}{
						{
							"x_forwarded_for_client_cert_client_verify_alias":   "test_client-verify-alias_223451",
							"x_forwarded_for_client_cert_client_verify_enabled": "true",
							"x_forwarded_for_client_cert_finger_print_alias":    "test_client-verify-alias_223452",
							"x_forwarded_for_client_cert_finger_print_enabled":  "true",
							"x_forwarded_for_client_cert_issuer_dn_alias":       "test_client-verify-alias_223453",
							"x_forwarded_for_client_cert_issuer_dn_enabled":     "true",
							"x_forwarded_for_client_cert_subject_dn_alias":      "test_client-verify-alias_223454",
							"x_forwarded_for_client_cert_subject_dn_enabled":    "true",
							"x_forwarded_for_client_src_port_enabled":           "true",
							"x_forwarded_for_enabled":                           "true",
							"x_forwarded_for_proto_enabled":                     "true",
							"x_forwarded_for_slb_id_enabled":                    "true",
							"x_forwarded_for_slb_port_enabled":                  "true",
							"x_forwarded_for_client_source_ips_enabled":         "true",
							"x_forwarded_for_client_source_ips_trusted":         "192.168.1.0/24",
							"x_forwarded_for_host_enabled":                      "true",
							"x_forwarded_for_processing_mode":                   "append",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":         CHECKSET,
						"listener_protocol":        "HTTPS",
						"listener_port":            port,
						"listener_description":     "tf-testAccListener_new",
						"default_actions.#":        "1",
						"certificates.#":           "1",
						"acl_config.#":             "1",
						"gzip_enabled":             "true",
						"http2_enabled":            "true",
						"idle_timeout":             "20",
						"request_timeout":          "60",
						"security_policy_id":       "tls_cipher_policy_1_0",
						"x_forwarded_for_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_enabled": "false",
					"ca_certificates": []map[string]interface{}{
						{
							"certificate_id": "${var.ca_certificate_id}",
						},
					},
					"x_forwarded_for_config": []map[string]interface{}{
						{
							"x_forwarded_for_client_cert_client_verify_alias":   "test_client-verify-alias_123451",
							"x_forwarded_for_client_cert_client_verify_enabled": "true",
							"x_forwarded_for_client_cert_finger_print_alias":    "test_client-verify-alias_123452",
							"x_forwarded_for_client_cert_finger_print_enabled":  "true",
							"x_forwarded_for_client_cert_issuer_dn_alias":       "test_client-verify-alias_123453",
							"x_forwarded_for_client_cert_issuer_dn_enabled":     "true",
							"x_forwarded_for_client_cert_subject_dn_alias":      "test_client-verify-alias_123454",
							"x_forwarded_for_client_cert_subject_dn_enabled":    "true",
							"x_forwarded_for_client_src_port_enabled":           "true",
							"x_forwarded_for_enabled":                           "true",
							"x_forwarded_for_proto_enabled":                     "true",
							"x_forwarded_for_slb_id_enabled":                    "true",
							"x_forwarded_for_slb_port_enabled":                  "true",
							"x_forwarded_for_client_source_ips_enabled":         "false",
							"x_forwarded_for_client_source_ips_trusted":         "192.168.1.0/24",
							"x_forwarded_for_host_enabled":                      "true",
							"x_forwarded_for_processing_mode":                   "remove",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"x_forwarded_for_config.#": "1",
					}),
				),
			},
			//x_forwarded_for_enabled cannot be set to false
			{
				Config: testAccConfig(map[string]interface{}{
					"ca_enabled": "true",
					"ca_certificates": []map[string]interface{}{
						{
							"certificate_id": "${var.ca_certificate_id_update}",
						},
					},
					"x_forwarded_for_config": []map[string]interface{}{
						{
							"x_forwarded_for_client_cert_client_verify_alias":   "test_client-verify-alias_123451",
							"x_forwarded_for_client_cert_client_verify_enabled": "false",
							"x_forwarded_for_client_cert_finger_print_alias":    "test_client-verify-alias_123452",
							"x_forwarded_for_client_cert_finger_print_enabled":  "false",
							"x_forwarded_for_client_cert_issuer_dn_alias":       "test_client-verify-alias_123453",
							"x_forwarded_for_client_cert_issuer_dn_enabled":     "false",
							"x_forwarded_for_client_cert_subject_dn_alias":      "test_client-verify-alias_123454",
							"x_forwarded_for_client_cert_subject_dn_enabled":    "false",
							"x_forwarded_for_client_src_port_enabled":           "false",
							"x_forwarded_for_enabled":                           "false",
							"x_forwarded_for_proto_enabled":                     "false",
							"x_forwarded_for_slb_id_enabled":                    "false",
							"x_forwarded_for_slb_port_enabled":                  "false",
							"x_forwarded_for_client_source_ips_enabled":         "true",
							"x_forwarded_for_client_source_ips_trusted":         "192.168.2.0/24",
							"x_forwarded_for_host_enabled":                      "false",
							"x_forwarded_for_processing_mode":                   "remove",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"x_forwarded_for_config.#": "1",
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

func TestAccAliCloudALBListener_basic4(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AliCloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 1000)
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudALBListenerBasicDependence0)
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
					"load_balancer_id":     "${local.load_balancer_id}",
					"listener_protocol":    "HTTPS",
					"listener_port":        port,
					"listener_description": "tf-testAccListener_new",
					"default_actions": []map[string]interface{}{
						{
							"type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"server_group_id": "${local.server_group_id}",
										},
									},
								},
							},
						},
					},
					"certificates": []map[string]interface{}{
						{
							"certificate_id": "${local.certificate_id}",
						},
					},
					"acl_config": []map[string]interface{}{
						{
							"acl_type": "White",
							"acl_relations": []map[string]interface{}{
								{
									"acl_id": "${alicloud_alb_acl.default.0.id}",
								},
							},
						},
					},
					"gzip_enabled":       "true",
					"http2_enabled":      "true",
					"idle_timeout":       "20",
					"request_timeout":    "60",
					"security_policy_id": "tls_cipher_policy_1_0",
					"x_forwarded_for_config": []map[string]interface{}{
						{
							"x_forwarded_for_client_cert_client_verify_alias":   "test_client-verify-alias_223451",
							"x_forwarded_for_client_cert_client_verify_enabled": "true",
							"x_forwarded_for_client_cert_finger_print_alias":    "test_client-verify-alias_223452",
							"x_forwarded_for_client_cert_finger_print_enabled":  "true",
							"x_forwarded_for_client_cert_issuer_dn_alias":       "test_client-verify-alias_223453",
							"x_forwarded_for_client_cert_issuer_dn_enabled":     "true",
							"x_forwarded_for_client_cert_subject_dn_alias":      "test_client-verify-alias_223454",
							"x_forwarded_for_client_cert_subject_dn_enabled":    "true",
							"x_forwarded_for_client_src_port_enabled":           "true",
							"x_forwarded_for_enabled":                           "true",
							"x_forwarded_for_proto_enabled":                     "true",
							"x_forwarded_for_slb_id_enabled":                    "true",
							"x_forwarded_for_slb_port_enabled":                  "true",
							"x_forwarded_for_client_source_ips_enabled":         "true",
							"x_forwarded_for_client_source_ips_trusted":         "192.168.1.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":         CHECKSET,
						"listener_protocol":        "HTTPS",
						"listener_port":            port,
						"listener_description":     "tf-testAccListener_new",
						"default_actions.#":        "1",
						"certificates.#":           "1",
						"acl_config.#":             "1",
						"gzip_enabled":             "true",
						"http2_enabled":            "true",
						"idle_timeout":             "20",
						"request_timeout":          "60",
						"security_policy_id":       "tls_cipher_policy_1_0",
						"x_forwarded_for_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"x_forwarded_for_config": []map[string]interface{}{
						{
							"x_forwarded_for_client_cert_client_verify_alias":   "test_client-verify-alias_123451",
							"x_forwarded_for_client_cert_client_verify_enabled": "true",
							"x_forwarded_for_client_cert_finger_print_alias":    "test_client-verify-alias_123452",
							"x_forwarded_for_client_cert_finger_print_enabled":  "true",
							"x_forwarded_for_client_cert_issuer_dn_alias":       "test_client-verify-alias_123453",
							"x_forwarded_for_client_cert_issuer_dn_enabled":     "true",
							"x_forwarded_for_client_cert_subject_dn_alias":      "test_client-verify-alias_123454",
							"x_forwarded_for_client_cert_subject_dn_enabled":    "true",
							"x_forwarded_for_client_src_port_enabled":           "true",
							"x_forwarded_for_enabled":                           "true",
							"x_forwarded_for_proto_enabled":                     "true",
							"x_forwarded_for_slb_id_enabled":                    "true",
							"x_forwarded_for_slb_port_enabled":                  "true",
							"x_forwarded_for_client_source_ips_enabled":         "false",
							"x_forwarded_for_client_source_ips_trusted":         "192.168.1.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"x_forwarded_for_config.#": "1",
					}),
				),
			},
			//x_forwarded_for_enabled cannot be set to false
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"x_forwarded_for_config": []map[string]interface{}{
			//			{
			//				"x_forwarded_for_client_cert_client_verify_alias":   "test_client-verify-alias_123451",
			//				"x_forwarded_for_client_cert_client_verify_enabled": "false",
			//				"x_forwarded_for_client_cert_finger_print_alias":    "test_client-verify-alias_123452",
			//				"x_forwarded_for_client_cert_finger_print_enabled":  "false",
			//				"x_forwarded_for_client_cert_issuer_dn_alias":       "test_client-verify-alias_123453",
			//				"x_forwarded_for_client_cert_issuer_dn_enabled":     "false",
			//				"x_forwarded_for_client_cert_subject_dn_alias":      "test_client-verify-alias_123454",
			//				"x_forwarded_for_client_cert_subject_dn_enabled":    "false",
			//				"x_forwarded_for_client_src_port_enabled":           "false",
			//				"x_forwarded_for_enabled":                           "false",
			//				"x_forwarded_for_proto_enabled":                     "false",
			//				"x_forwarded_for_slb_id_enabled":                    "false",
			//				"x_forwarded_for_slb_port_enabled":                  "false",
			//				"x_forwarded_for_client_source_ips_enabled":         "true",
			//				"x_forwarded_for_client_source_ips_trusted":         "192.168.2.0/24",
			//			},
			//		},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"x_forwarded_for_config.#": "1",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"quic_config": []map[string]interface{}{
						{
							"quic_upgrade_enabled": "true",
							"quic_listener_id":     "${alicloud_alb_listener.quic.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quic_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quic_config": []map[string]interface{}{
						{
							"quic_upgrade_enabled": "false",
							"quic_listener_id":     "",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quic_config.#": "1",
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

var AliCloudALBListenerMap0 = map[string]string{
	"access_log_record_customized_headers_enabled": CHECKSET,
	"gzip_enabled":       CHECKSET,
	"http2_enabled":      CHECKSET,
	"idle_timeout":       CHECKSET,
	"request_timeout":    CHECKSET,
	"security_policy_id": CHECKSET,
	"status":             CHECKSET,
}

func AliCloudALBListenerBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	variable "ca_certificate_id" {
  		default = "1efd3cac-e8c1-611d-b971-f96e5e432aba"
	}

	variable "ca_certificate_id_update" {
  		default = "1efd3caf-3cc7-6f42-be33-29ca9cecdfd3"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_alb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "vswitch_1" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  		zone_id      = data.alicloud_alb_zones.default.zones.0.id
  		vswitch_name = var.name
	}

	resource "alicloud_vswitch" "vswitch_2" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
  		zone_id      = data.alicloud_alb_zones.default.zones.1.id
  		vswitch_name = var.name
	}

	resource "alicloud_log_project" "default" {
  		name        = var.name
  		description = "created by terraform"
	}

	resource "alicloud_log_store" "default" {
  		project               = alicloud_log_project.default.name
  		name                  = var.name
  		shard_count           = 3
  		auto_split            = true
  		max_split_shard_count = 60
  		append_meta           = true
	}

	resource "alicloud_alb_load_balancer" "default_3" {
  		vpc_id                 = alicloud_vpc.default.id
  		address_type           = "Internet"
  		address_allocated_mode = "Fixed"
  		load_balancer_name     = var.name
  		load_balancer_edition  = "Standard"
  		resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  		load_balancer_billing_config {
    		pay_type = "PayAsYouGo"
  		}
  		tags = {
    		Created = "TF"
  		}
  		zone_mappings {
    		vswitch_id = alicloud_vswitch.vswitch_1.id
    		zone_id    = data.alicloud_alb_zones.default.zones.0.id
  		}
  		zone_mappings {
    		vswitch_id = alicloud_vswitch.vswitch_2.id
    		zone_id    = data.alicloud_alb_zones.default.zones.1.id
  		}
  		modification_protection_config {
    		status = "NonProtection"
  		}
  		access_log_config {
    		log_project = alicloud_log_project.default.name
    		log_store   = alicloud_log_store.default.name
  		}
	}

	resource "alicloud_alb_server_group" "default_4" {
  		count             = 2
  		protocol          = "HTTP"
  		vpc_id            = alicloud_vpc.default.id
  		server_group_name = var.name
  		resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  		health_check_config {
    		health_check_enabled = "false"
  		}
  		sticky_session_config {
    		sticky_session_enabled = "false"
  		}
  		tags = {
    		Created = "TF"
  		}
	}

	resource "alicloud_ssl_certificates_service_certificate" "default" {
  		certificate_name = var.name
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID1zCCAr+gAwIBAgIRAOrWWz1qmkcSg90JDHjuzFwwDQYJKoZIhvcNAQELBQAw
XjELMAkGA1UEBhMCQ04xDjAMBgNVBAoTBU15U1NMMSswKQYDVQQLEyJNeVNTTCBU
ZXN0IFJTQSAtIEZvciB0ZXN0IHVzZSBvbmx5MRIwEAYDVQQDEwlNeVNTTC5jb20w
HhcNMjQxMTI2MDczNjA4WhcNMjkxMTI1MDczNjA4WjAgMQswCQYDVQQGEwJDTjER
MA8GA1UEAxMIdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDa7HDGbQ1Km0f4ZaFzYbjVN0q8KkvZ+oQUd4naGOZnlH5k0XFwmjg+TWf88YX3
5IF8c45/rXrTWucPLg7FeqR96Wq9HZEmzEhs6VG031V9Hqa32saRScCOAyhiW7Hj
OWf6BZveuxbZNbgQCR59QzX4CeAIC68xavIDAy3wcTAH9cIkD71BxEPJGGR7BIVH
9DcWXaMAnJqQfrkth0xHBjflZABHAI0wPYPfaw8fd9DRkMYOIkfjwrrcL5IvhI1u
D3wdHJQWA2vR8hjoU4dHiJLbUtQ+xV1UGVkF67CpQ6LDjSQdX7xlZ7WJMc/7dCJ9
a7tr0ZTwq4/3KSgcRvm62oGvAgMBAAGjgc0wgcowDgYDVR0PAQH/BAQDAgWgMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSMEGDAWgBQogSYF0TQa
P8FzD7uTzxUcPwO/fzBjBggrBgEFBQcBAQRXMFUwIQYIKwYBBQUHMAGGFWh0dHA6
Ly9vY3NwLm15c3NsLmNvbTAwBggrBgEFBQcwAoYkaHR0cDovL2NhLm15c3NsLmNv
bS9teXNzbHRlc3Ryc2EuY3J0MBMGA1UdEQQMMAqCCHRlc3QuY29tMA0GCSqGSIb3
DQEBCwUAA4IBAQAxPOlK5WBA9kITzxYyjqe/YvWzfMlsmj0yvpyHrPeZf7HZTTFz
ebYkzrHL8ZLyOHBhag0nL7Poj6ek98NoXTuCYCi8LspdadapOeYQzLce3beu/frk
sqU0A6WLHG9Ol9yUDMCX7xvLoAY/LDrcOM3Z87C/u/ykB4wKfFN2XfR3EZx3PQqw
sV77LOnyQixB4FMHpHlKuDoUkSN9uvxwEPOeGnLZXm96hPsjPwk1bDM8qerNPpVI
CwJ6kNuZ2eLz2Umqu2Gh3l4aADdIwxRY1OOjjZNut8STosABKWVGIwQbbAdRPQze
qHZ05oVTjFy9L1DAzhQ5Zn3oUjLl5KW4tYBA
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2uxwxm0NSptH+GWhc2G41TdKvCpL2fqEFHeJ2hjmZ5R+ZNFx
cJo4Pk1n/PGF9+SBfHOOf61601rnDy4OxXqkfelqvR2RJsxIbOlRtN9VfR6mt9rG
kUnAjgMoYlux4zln+gWb3rsW2TW4EAkefUM1+AngCAuvMWryAwMt8HEwB/XCJA+9
QcRDyRhkewSFR/Q3Fl2jAJyakH65LYdMRwY35WQARwCNMD2D32sPH3fQ0ZDGDiJH
48K63C+SL4SNbg98HRyUFgNr0fIY6FOHR4iS21LUPsVdVBlZBeuwqUOiw40kHV+8
ZWe1iTHP+3QifWu7a9GU8KuP9ykoHEb5utqBrwIDAQABAoIBAQCErEfIKOymKybZ
pZXLnAxswt563FMtngGPecZEM1TmrvpOVROffwbY0wZTJ3fd/FBwwIM6Y0MNdYiU
DYCMM0AewmeahqGh1qmJv3hx2eswMXQt9driz8RvDADcYt+SagbWYbHNsKovJrwO
k8gzd5jsYeewWIxqsXpLUxDzJ1VJbIqoHgkrirRRPo0onpixPWeA0RbElSwjwIUw
y43cC4WF8N7wot3cTST8yeKM8ujtqpN22ZtKnbkHTd03vnwQTMeUMJeDQmSmY5aJ
yFr7yw/Z66+7Amh6pkWhzZSDHsjI4y/S3CCdpwFlMA7ID590umJB6HFxWsmVacSe
MSs2vIJZAoGBAOiecPH1HVDQqH6PcrN/X9E3pDKSyAj+nHsVDGIZsie9f5g/qA0A
tcJtQLS0CzrpMTLsAnsfdh2T7Lg6pYFz5jnOUyMjOImAEbCtgvqBxqgFea//OhdP
8s/RmxKIAenBsk7Wbwx8/KPhbZLUNe8OnILVHDfS6kLSa49Iu+4UvrpNAoGBAPDt
mky5MMHKdHwbqxPo9jYrz1m3gqqIvv+VihO4t/DE6t2Zg43ctfFm1BVEDSwPjYs/
YV69KfVrVRUnzMZVdtHZ/dBK784YTY0OujemoaIzMKFIL8tbJFldVv2IgB+IelTX
e675hVdHjNUqZhHwccd8X6d/8icohZw62SNHb/HrAoGBAN1HSt1/c6Gau42Y212Q
fw9ARLuvEQYtXaFfxmXTV7uh8axccXndAQmwb+r1kfE6PojYJQwGQ4+jVX1ynFnm
bEz0zfUQ3gk+gJV2mK+/n7/ZZYZb3WCrtqimFUOtiVRZ40pHhV91zcX+/QK9R4je
d1elbbBUvG9QRu0IHW0+4qfJAoGAOmlQvIM1l/ZOsXw/yO71KoMKnXTJYDERJYQK
2ucw6VXEn39FjtJQ5jsI9jLugp0usvDl2YNBNfgUw7FHi1pTGWOhjqtsYmov+x/z
8+QZUerZQnDu7X2mXWgs3AEJFxwOlJ09pllmg5ecRF4oKvdBjpzP0BtMCURgyFTY
Kh56vIsCgYBMbneMvFY6PCESKIAXj16BF4lqYVXFqHVoxyfxIuVlAy3TMNwxvpbS
yDETk05Ux9yNES0WyTb1SWVG1o1wXc0dnDXCwJqLC1tzJUNUSD1AYvktoNIFErcN
gs3ercrzBTX5ezORPj9ErRAPrSq+V3z1Lge5Gl+EqgDvAfnknww75w==
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_ssl_certificates_service_certificate" "update" {
  		certificate_name = "${var.name}-update"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID1jCCAr6gAwIBAgIQGKYS2rt7QuCbV3mpxs2D9DANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yNDExMjYwNzM3NTFaFw0yOTExMjUwNzM3NTFaMCAxCzAJBgNVBAYTAkNOMREw
DwYDVQQDEwh0ZXN0LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
ANqXuMEuRfqQ94tlv44PmYbbuJN6d1qaYu5DozPfOrqKMUD5TRkY1+hZtmC+36Ze
EuwQplYK6O+eaMXAloXL9Ofo4oz5Ny6fo6vjN32dcwD3iCYuQY6YpNQlnpl2jb7K
yh8CQYWbkGQ+U3Yg7K2ewp2HjWLBR0ODzGrcej0csbQ2WJtVzm5ptAbRfdLQADQ0
Q9ZmQ2RU4vmCqHGN7xZdnCEoWUMlvec++DsRB94URyAEsU+Z7hDzRR7723HSszry
Q+3aZfqlu4iq852lRQGUYJ8KoGyUWGlynnREB93KyGchG+x/lgADAYWlJh/19CgM
ElY4s1bqTbCUrltlgSA5qZMCAwEAAaOBzTCByjAOBgNVHQ8BAf8EBAMCBaAwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQYMBaAFCiBJgXRNBo/
wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEFBQcwAYYVaHR0cDov
L29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8vY2EubXlzc2wuY29t
L215c3NsdGVzdHJzYS5jcnQwEwYDVR0RBAwwCoIIdGVzdC5jb20wDQYJKoZIhvcN
AQELBQADggEBABhk5x58ctbUkoM5z18bT8Ny0Ko9p0P6wn5XduK7JWD9QwjM5ZKr
kA39pHQU9D4sGhEhLR9SlWvSmrVQmSRn5tn03eHRXhhGv87IWmkTPHBYkoz8LP4L
ArYjAZpo9odmWpH6C+IkhqUw9nPg31na9wwVdUBCYxuIlL36PoII16FNsWwBnKMi
X81UCm+1UHp4qF3dT6s34ttEVNRoYw/u3rnwqVtnwTDs4svcLMaRyyNZrgV2RG5L
LC5tM9mrqvbKQIvQRxc47V1FV+t4jNun7se4St5nWEAavdLwmS1K/1QLb9UmYJOv
Nw8ocKgnHvrCoI59SQSO+oin+weDMchDK6U=
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2pe4wS5F+pD3i2W/jg+Zhtu4k3p3Wppi7kOjM986uooxQPlN
GRjX6Fm2YL7fpl4S7BCmVgro755oxcCWhcv05+jijPk3Lp+jq+M3fZ1zAPeIJi5B
jpik1CWemXaNvsrKHwJBhZuQZD5TdiDsrZ7CnYeNYsFHQ4PMatx6PRyxtDZYm1XO
bmm0BtF90tAANDRD1mZDZFTi+YKocY3vFl2cIShZQyW95z74OxEH3hRHIASxT5nu
EPNFHvvbcdKzOvJD7dpl+qW7iKrznaVFAZRgnwqgbJRYaXKedEQH3crIZyEb7H+W
AAMBhaUmH/X0KAwSVjizVupNsJSuW2WBIDmpkwIDAQABAoIBAQCK4uuYknYUBhfC
khtrf63kaaaUzbMX9g/1ozQGuUbvTu6MgdnioE5OavHd9mjTo+IR62JEORpXZSbc
vsjkqfopf2aye4X8MaIkjHGtdmSjsKLo32r31zSjNmPWzeSx3NcfbKeE5JqRlqgg
3jqC9eRhgsbqgDNvSkaPfxaLzbd67/+KSdituRNqMGvKCgyZAT63yLiO7ArdtEaY
Ij+BSECjABmhue+sBWtObmovI+MGJ7RetnBRaFh5/3I6rd0bY3dyhwab0A2rWuM7
T1usQSZ/Z8c4s1V2anQ8AgvcAe2bAfSSRoCUNwuPMtyj1LJVk2MaxaiLfpUNpKb7
r5P3fP2BAoGBAN18C/Tp9duDoAf+LYUtV4riXQ1CeCw/wxL6g9X0e+dnP7L4kLe4
m+/YSZUkv7IlRf56p4t9r0if2+w7u95zzXt3k7PLuFRigjSpxnZ1hrYKKctNY2Oj
urEUV+dkoekplFC0kSFtOFYaNwewaTV/fkWa0Apd/ZnccbViLB2zYOIZAoGBAPyo
ThxDE69gAOEiQm8B29bAi1lM98Dx69KSrOXP8Yf/mOUJyRnRphYVmYHKA/T5Gubt
Rn6o849Le4mJNWyyrgllg5QEMfShBBzndc5tL4ltsIrzQTKu/we/GlwxqRVccJ1B
Tn4+76gvMKpFvFEDtHQp9/XWy/FhYY3/VO7qTtaLAoGAUe/TKI7pIoVmTa6tvmgQ
y9OEYyRk+tG35CyDW0KwF+JtgVNNjnogTjGwvxkyRcBeTY+orgUYNIDXRmSu0tP6
f6O0I77I+Ybb7omkXyyJYo0N+yUtEK6AoYQKJRNohq6YLOcwDbKvNcNK+nA768u3
th5Yuo0dBa+07UpdUbuLqvkCgYEApLKN4Gx1S5AYYqmzhqs+hDoVXEwJANRytlx4
qoIn31BleYAsgFEipCjGXU2z0KAFwl0P5Ab8Zf99c0Vm9wlu258562XkrqO7i5/y
MnMIVtyTBbDWYlSi2IjhhRG2N79/hXMJ2M/r58WDQqucu27f1g15nt67KQkiz66O
zgMdC0sCgYAd3QLHQfHBxqlHBokcjdHxWoX2fkKwdQlKlKuk6Q+quyrY0dIF2dxr
/suURAMr4407dP4cjrG9LfWGGYfpcqt79/QDa7rbp9z6zdu6CU+RzqZyfgAtcd6r
1LeiSMDF5dMPJoxkrA9/aKMmp4UbYv/UTexUQ41tK/PFTG6fye44pA==
-----END RSA PRIVATE KEY-----
EOF
	}

	resource "alicloud_alb_acl" "default" {
  		count    = 2
  		acl_name = var.name
	}

	locals {
  		load_balancer_id       = alicloud_alb_load_balancer.default_3.id
  		server_group_id        = alicloud_alb_server_group.default_4.0.id
  		server_group_id_update = alicloud_alb_server_group.default_4.1.id
  		certificate_id         = join("", [alicloud_ssl_certificates_service_certificate.default.id, "-%s"])
  		certificate_id_update  = join("-", [alicloud_ssl_certificates_service_certificate.update.id, "%s"])
	}

	resource "alicloud_alb_listener" "quic" {
	  load_balancer_id     = local.load_balancer_id
	  listener_protocol    = "QUIC"
	  listener_port        = 443
	  listener_description = var.name
	  certificates {
		certificate_id = local.certificate_id_update
	  }
	  default_actions {
		type = "ForwardGroup"
		forward_group_config {
		  server_group_tuples {
			server_group_id = local.server_group_id
		  }
		}
	  }
	}

`, name, defaultRegionToTest, defaultRegionToTest)
}

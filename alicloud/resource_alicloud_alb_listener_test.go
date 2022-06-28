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
		return fmt.Errorf("error getting Alicloud client: %s", err)
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

func TestAccAlicloudALBListener_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBListenerBasicDependence0)
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id":     CHECKSET,
						"listener_protocol":    "HTTPS",
						"listener_port":        port,
						"listener_description": "tf-testAccListener_new",
						"default_actions.#":    "1",
						"certificates.#":       "1",
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
											"server_group_id": "${local.server_group_id}",
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
							"certificate_id": "${local.certificate_id}",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run", "xforwarded_for_config"},
			},
		},
	})
}
func TestAccAlicloudALBListener_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBListenerBasicDependence0)
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
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run", "xforwarded_for_config"},
			},
		},
	})
}

func TestAccAlicloudALBListener_basic2(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testaccalblistener%d", rand)
	port := fmt.Sprintf("%d", acctest.RandIntRange(1, 1000))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBListenerBasicDependence0)
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
					"dry_run":            "false",
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
						"dry_run":                 "false",
						"gzip_enabled":            "true",
						"http2_enabled":           "true",
						"idle_timeout":            "20",
						"request_timeout":         "60",
						"security_policy_id":      "tls_cipher_policy_1_0",
						"xforwarded_for_config.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run", "xforwarded_for_config"},
			},
		},
	})
}

func TestAccAlicloudALBListener_basic3(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_alb_listener.default"
	ra := resourceAttrInit(resourceId, AlicloudALBListenerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbListener")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 1000)
	port := fmt.Sprintf("%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, port, AlicloudALBListenerBasicDependence0)
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
					"dry_run":            "false",
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
						"dry_run":                  "false",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBListenerMap0 = map[string]string{
	"dry_run": NOSET,
}

func AlicloudALBListenerBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {	
	default = "tf-testaccalblistener%s"
}

data "alicloud_alb_zones" "default"{}

data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count             = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id 			=  data.alicloud_alb_zones.default.zones.0.id
  vswitch_name      = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count             = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name              = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {}

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
  vpc_id =                data.alicloud_vpcs.default.ids.0
  address_type =        "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name =    var.name
  load_balancer_edition = "Standard"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = 	"PayAsYouGo"
  }
  tags = {
		Created = "TF"
  }
  zone_mappings{
		vswitch_id =  length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
		zone_id =  data.alicloud_alb_zones.default.zones.0.id
	}
  zone_mappings{
		vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
		zone_id =   data.alicloud_alb_zones.default.zones.1.id
	}
  modification_protection_config{
	status = "NonProtection"
  }
  access_log_config{
  	log_project = alicloud_log_project.default.name
  	log_store =   alicloud_log_store.default.name
  }
}

resource "alicloud_alb_server_group" "default_4" {
	protocol = "HTTP"
	vpc_id = data.alicloud_vpcs.default.vpcs.0.id
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
  cert = <<EOF
-----BEGIN CERTIFICATE-----
MIIDRjCCAq+gAwIBAgIJAJn3ox4K13PoMA0GCSqGSIb3DQEBBQUAMHYxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCQkoxDDAKBgNVBAoTA0FMSTEP
MA0GA1UECxMGQUxJWVVOMQ0wCwYDVQQDEwR0ZXN0MR8wHQYJKoZIhvcNAQkBFhB0
ZXN0QGhvdG1haWwuY29tMB4XDTE0MTEyNDA2MDQyNVoXDTI0MTEyMTA2MDQyNVow
djELMAkGA1UEBhMCQ04xCzAJBgNVBAgTAkJKMQswCQYDVQQHEwJCSjEMMAoGA1UE
ChMDQUxJMQ8wDQYDVQQLEwZBTElZVU4xDTALBgNVBAMTBHRlc3QxHzAdBgkqhkiG
9w0BCQEWEHRlc3RAaG90bWFpbC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJ
AoGBAM7SS3e9+Nj0HKAsRuIDNSsS3UK6b+62YQb2uuhKrp1HMrOx61WSDR2qkAnB
coG00Uz38EE+9DLYNUVQBK7aSgLP5M1Ak4wr4GqGyCgjejzzh3DshUzLCCy2rook
KOyRTlPX+Q5l7rE1fcSNzgepcae5i2sE1XXXzLRIDIvQxcspAgMBAAGjgdswgdgw
HQYDVR0OBBYEFBdy+OuMsvbkV7R14f0OyoLoh2z4MIGoBgNVHSMEgaAwgZ2AFBdy
+OuMsvbkV7R14f0OyoLoh2z4oXqkeDB2MQswCQYDVQQGEwJDTjELMAkGA1UECBMC
QkoxCzAJBgNVBAcTAkJKMQwwCgYDVQQKEwNBTEkxDzANBgNVBAsTBkFMSVlVTjEN
MAsGA1UEAxMEdGVzdDEfMB0GCSqGSIb3DQEJARYQdGVzdEBob3RtYWlsLmNvbYIJ
AJn3ox4K13PoMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAY7KOsnyT
cQzfhiiG7ASjiPakw5wXoycHt5GCvLG5htp2TKVzgv9QTliA3gtfv6oV4zRZx7X1
Ofi6hVgErtHaXJheuPVeW6eAW8mHBoEfvDAfU3y9waYrtUevSl07643bzKL6v+Qd
DUBTxOAvSYfXTtI90EAxEG/bJJyOm5LqoiA=
-----END CERTIFICATE-----
EOF
  key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDO0kt3vfjY9BygLEbiAzUrEt1Cum/utmEG9rroSq6dRzKzsetV
kg0dqpAJwXKBtNFM9/BBPvQy2DVFUASu2koCz+TNQJOMK+BqhsgoI3o884dw7IVM
ywgstq6KJCjskU5T1/kOZe6xNX3Ejc4HqXGnuYtrBNV118y0SAyL0MXLKQIDAQAB
AoGAfe3NxbsGKhN42o4bGsKZPQDfeCHMxayGp5bTd10BtQIE/ST4BcJH+ihAS7Bd
6FwQlKzivNd4GP1MckemklCXfsVckdL94e8ZbJl23GdWul3v8V+KndJHqv5zVJmP
hwWoKimwIBTb2s0ctVryr2f18N4hhyFw1yGp0VxclGHkjgECQQD9CvllsnOwHpP4
MdrDHbdb29QrobKyKW8pPcDd+sth+kP6Y8MnCVuAKXCKj5FeIsgVtfluPOsZjPzz
71QQWS1dAkEA0T0KXO8gaBQwJhIoo/w6hy5JGZnrNSpOPp5xvJuMAafs2eyvmhJm
Ev9SN/Pf2VYa1z6FEnBaLOVD6hf6YQIsPQJAX/CZPoW6dzwgvimo1/GcY6eleiWE
qygqjWhsh71e/3bz7yuEAnj5yE3t7Zshcp+dXR3xxGo0eSuLfLFxHgGxwQJAAxf8
9DzQ5NkPkTCJi0sqbl8/03IUKTgT6hcbpWdDXa7m8J3wRr3o5nUB+TPQ5nzAbthM
zWX931YQeACcwhxvHQJBAN5mTzzJD4w4Ma6YTaNHyXakdYfyAWrOkPIWZxfhMfXe
DrlNdiysTI4Dd1dLeErVpjsckAaOW/JDG5PCSwkaMxk=
-----END RSA PRIVATE KEY-----
EOF
}
resource "alicloud_alb_acl" "default" {
	acl_name = var.name
	count = 2
}

locals {
 load_balancer_id = alicloud_alb_load_balancer.default_3.id
 server_group_id = alicloud_alb_server_group.default_4.id
 certificate_id = join("",[alicloud_ssl_certificates_service_certificate.default.id,"-%s"])
}
`, name, defaultRegionToTest)
}

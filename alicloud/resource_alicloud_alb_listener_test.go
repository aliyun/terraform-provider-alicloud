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
			{
				Config: testAccConfig(map[string]interface{}{
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
			{
				Config: testAccConfig(map[string]interface{}{
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
  		key              = <<EOF
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

	resource "alicloud_ssl_certificates_service_certificate" "update" {
  		certificate_name = "${var.name}-update"
  		cert             = <<EOF
-----BEGIN CERTIFICATE-----
MIID7jCCAtagAwIBAgIQUNnSVa/sQNeb9pBN9NhkwTANBgkqhkiG9w0BAQsFADBe
MQswCQYDVQQGEwJDTjEOMAwGA1UEChMFTXlTU0wxKzApBgNVBAsTIk15U1NMIFRl
c3QgUlNBIC0gRm9yIHRlc3QgdXNlIG9ubHkxEjAQBgNVBAMTCU15U1NMLmNvbTAe
Fw0yMzA4MDkwMzM4MThaFw0yODA4MDcwMzM4MThaMCwxCzAJBgNVBAYTAkNOMR0w
GwYDVQQDExRhbGljbG91ZC1wcm92aWRlci5jbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAOgskr8dEfZYdjr0xaIqlCkmE802vABoj3SQNn3rLWnUj+1v
Wqbpsj6Bu61Scb8mtl/OZOOM7sgq0Q1hpdO8xvMGxTMuZ2bjX0EqCMqh4AvFofHL
a/iVD07hfoM1Jo8CEidh1uvcOuXP1TlaqU020x1TX3a3niJu4JVkmCkCOwAbWYuj
O8IsgBCsFaF9d4+C1JRYOtRbIHCNhd0sxG8AGovUDLvlkePeH5NF7DNvFXgGJ4iv
EQcY9pP08RBFUkaznOw/r64Up7zhLb+Ie4SyAvs1FulhMAmIXOcbsND39hJ+/WIP
8beWvIN1eCS8zcvgAvDgMkV8oqqVbQu1dqx5WuMCAwEAAaOB2TCB1jAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFCiBJgXRNBo/wXMPu5PPFRw/A79/MGMGCCsGAQUFBwEBBFcwVTAhBggrBgEF
BQcwAYYVaHR0cDovL29jc3AubXlzc2wuY29tMDAGCCsGAQUFBzAChiRodHRwOi8v
Y2EubXlzc2wuY29tL215c3NsdGVzdHJzYS5jcnQwHwYDVR0RBBgwFoIUYWxpY2xv
dWQtcHJvdmlkZXIuY24wDQYJKoZIhvcNAQELBQADggEBALd0hFZAd2XHJgETbHQs
h4YUBNKxrIy6JiWfxffhIL1ZK5pI443DC4VRGfxVi3zWqs01WbNtJ2b1KdfSoovH
Zwi3hdMF1IwoAB/Y2sS4zjqS0H1od7MN9KKHes6bl3yCgpmaYs5cHbyg0IJHmeq3
rCgbKsvHfUwtzBNNPHlpANakAYd/5O1pztmUskWMUVaExfpMoQLo/AX9Lqm8pVjw
xs921I703l/E5zEnd3PVSYagy/KQJrwVt+wQZS11HsAryfO9kct/9f+c85VDo6Ht
iRirW/EnNPQRSno4z0V2x1Rn5+ZaoJo8cWzPvKrdfCG9TUozt4AR/LIudNLb6NNW
n7g=
-----END CERTIFICATE-----
EOF
  		key              = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA6CySvx0R9lh2OvTFoiqUKSYTzTa8AGiPdJA2festadSP7W9a
pumyPoG7rVJxvya2X85k44zuyCrRDWGl07zG8wbFMy5nZuNfQSoIyqHgC8Wh8ctr
+JUPTuF+gzUmjwISJ2HW69w65c/VOVqpTTbTHVNfdreeIm7glWSYKQI7ABtZi6M7
wiyAEKwVoX13j4LUlFg61FsgcI2F3SzEbwAai9QMu+WR494fk0XsM28VeAYniK8R
Bxj2k/TxEEVSRrOc7D+vrhSnvOEtv4h7hLIC+zUW6WEwCYhc5xuw0Pf2En79Yg/x
t5a8g3V4JLzNy+AC8OAyRXyiqpVtC7V2rHla4wIDAQABAoIBABKGQ+sluaIrKrvH
feFTfmDOHfRYsqVhslh9jSt80THJePZb1SLOMJ+WIFBS7Kpwv0pjoF8bho3IBMgJ
i36aaFFJsABGao+mApqjbPIl+kdWLHarYWEDG6aSjVKQshPk+WfVAZ3uA3EEpSGf
XzS+9Bc56LsDKYXbzOV+kjlraSO35AMec3CpISdx4K1caEAhKX6it9bvPq4pSYXi
PQspba0Jv46VV7MaabVjLzsinz5/md4vxyYHNIJAukHUfwJIsVC9ZNxukwSw+CzE
MMO64ylq2DGokNerGsLetuViV8UWi7qmUmms2fAmchodW16olgNkYTz27+V/A42S
eex63pkCgYEA+CqKhqp3qPe2E9KVrycrwjoycxmhOn3Iz1xiN7uAEv+DzfKtfZVf
mcOIiqw4Z82RkgjHb9vJuTigKdDkB1zE2gSDnep44sDWJM/5nPjGlMgnkiJWJhci
CnD0P4d6cT5wyDt7Q0/tS6ql2UrCpW4ktw1AP0Rm/z/VBD8jGkVenjcCgYEA74DM
Z2Qmh3bPt1TykpOlw+H+sEuvlkYxqMlbtn3Rv3WgEPIBekOFrgP7n/uLW1Aizn8w
EhNBBAE8w5jvklqZWYbpFMJQc09eqUkI8aTbLooZbzYj1f3CrzBRKn1GoTPmN9V0
j9r+TbH3/5CEoqlsJdmeQPofuv5Qid2oEutZcrUCgYBuZ16hco0xmqJiRzlYZvDM
w99V3X0g7Hy947e+W6gqy4nzwZb1W9LgMWE5cEzXwViVw1oWpY0k3dBDSi9oJxlc
dM2pH3sQRgH+9pdyAis2XaVdGfGBmKEITCAdc0RBxSmfqva3h4NmOlD2TpAx0MJ8
vWRrwR6hR+CYtw4CzgG+GQKBgQDGmi5lugW9JUe/xeBUrbyyv0+cT1auLUz2ouq7
XIA23Mo74wJYqW9Lyp+4nTWFJeGHDK8G/hJWyNPjeomG+jvZombbQPrHc9SSWi7h
eowKfpfywZlb1M7AyTc1HacY+9l3CTlcJQPl16NHuEZUQFue02NIjGENhd+xQy4h
ainFVQKBgAoPs9ebtWbBOaqGpOnquzb7WibvW/ifOfzv5/aOhkbpp+KCjcSON6CB
QF3BEXMcNMGWlpPrd8PaxCAzR4MyU///ekJri2icS9lrQhGSz2TtYhdED4pv1Aag
7eTPl5L7xAwphCSwy8nfCKmvlqcX/MSJ7A+LHB/2hdbuuEOyhpbu
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

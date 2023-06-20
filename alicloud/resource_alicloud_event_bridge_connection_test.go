package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test EventBridge Connection. >>> Resource test cases, automatically generated.
// Case 3084
func TestAccAlicloudEventBridgeConnection_basic3084(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeConnectionMap3084)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEventBridgeConnectionBasicDependence3084)
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
					"connection_name": name,
					"network_parameters": []map[string]interface{}{
						{
							"network_type":      "PublicNetwork",
							"vpc_id":            "eb-cn-huhehaote/vpc-hp3bdy0vbee0vb87fq2i6",
							"vswitche_id":       "vsw-hp3uinuttt9qbl27482v9",
							"security_group_id": "/sg-hp37abv61w1mpsuc0zco",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-connection-basic-pre",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-connection-basic-pre",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token",
									"api_key_value": "Token-value",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"password": "admin",
									"username": "admin",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"client_parameters": []map[string]interface{}{
										{
											"client_secret": "ClientSecret",
											"client_id":     "ClientId",
										},
									},
									"http_method": "POST",
									"oauth_http_parameters": []map[string]interface{}{
										{
											"body_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
											"header_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
										},
									},
								},
							},
							"authorization_type": "BASIC_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_parameters": []map[string]interface{}{
						{
							"network_type":      "PublicNetwork",
							"vpc_id":            "eb-cn-huhehaote/vpc-hp3bdy0vbee0vb87fq2i6",
							"vswitche_id":       "vsw-hp3uinuttt9qbl27482v9",
							"security_group_id": "/sg-hp37abv61w1mpsuc0zco",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-connection-basic-pre-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-connection-basic-pre-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_parameters": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"basic_auth_parameters": []map[string]interface{}{
								{},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-connection-basic-pre-update-api-key",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-connection-basic-pre-update-api-key",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_parameters": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"api_key_auth_parameters": []map[string]interface{}{
								{},
							},
							"authorization_type": "API_KEY_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-connection-basic-pre-update-oauth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-connection-basic-pre-update-oauth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_parameters": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"oauth_parameters": []map[string]interface{}{
								{
									"client_parameters": []map[string]interface{}{
										{
											"client_secret": "clientSecret",
											"client_id":     "clientId",
										},
									},
									"oauth_http_parameters": []map[string]interface{}{
										{
											"body_parameters": []map[string]interface{}{
												{},
											},
											"header_parameters": []map[string]interface{}{
												{
													"key":   "age",
													"value": "18",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{},
											},
										},
									},
								},
							},
							"authorization_type": "OAUTH_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_name": name + "_update",
					"description":     "test-connection-basic-pre",
					"network_parameters": []map[string]interface{}{
						{
							"network_type":      "PublicNetwork",
							"vpc_id":            "eb-cn-huhehaote/vpc-hp3bdy0vbee0vb87fq2i6",
							"vswitche_id":       "vsw-hp3uinuttt9qbl27482v9",
							"security_group_id": "/sg-hp37abv61w1mpsuc0zco",
						},
					},
					"auth_parameters": []map[string]interface{}{
						{
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token",
									"api_key_value": "Token-value",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"password": "admin",
									"username": "admin",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"client_parameters": []map[string]interface{}{
										{
											"client_secret": "ClientSecret",
											"client_id":     "ClientId",
										},
									},
									"http_method": "POST",
									"oauth_http_parameters": []map[string]interface{}{
										{
											"body_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
											"header_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
										},
									},
								},
							},
							"authorization_type": "BASIC_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name": name + "_update",
						"description":     "test-connection-basic-pre",
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

var AlicloudEventBridgeConnectionMap3084 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEventBridgeConnectionBasicDependence3084(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3084  twin
func TestAccAlicloudEventBridgeConnection_basic3084_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeConnectionMap3084)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEventBridgeConnectionBasicDependence3084)
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
					"connection_name": name,
					"description":     "test-connection-basic-pre-update-oauth",
					"network_parameters": []map[string]interface{}{
						{
							"network_type":      "PublicNetwork",
							"vpc_id":            "eb-cn-huhehaote/vpc-hp3bdy0vbee0vb87fq2i6",
							"vswitche_id":       "vsw-hp3uinuttt9qbl27482v9",
							"security_group_id": "/sg-hp37abv61w1mpsuc0zco",
						},
					},
					"auth_parameters": []map[string]interface{}{
						{
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token",
									"api_key_value": "Token-value",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"password": "admin",
									"username": "admin",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"client_parameters": []map[string]interface{}{
										{
											"client_secret": "clientSecret",
											"client_id":     "clientId",
										},
									},
									"http_method": "POST",
									"oauth_http_parameters": []map[string]interface{}{
										{
											"body_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
											"header_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "age",
													"value":           "18",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"is_value_secret": "true",
													"key":             "name",
													"value":           "name",
												},
											},
										},
									},
								},
							},
							"authorization_type": "OAUTH_AUTH",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name": name,
						"description":     "test-connection-basic-pre-update-oauth",
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

// Test EventBridge Connection. <<< Resource test cases, automatically generated.

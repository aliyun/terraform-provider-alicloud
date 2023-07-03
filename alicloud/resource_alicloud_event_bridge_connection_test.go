package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEventBridgeConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	resourceId := "alicloud_event_bridge_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeConnectionBasicDependence0)
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
							"network_type": "PublicNetwork",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name":      name,
						"network_parameters.#": "1",
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
					"network_parameters": []map[string]interface{}{
						{
							"network_type":      "PrivateNetwork",
							"vpc_id":            "${alicloud_vpc.default.id}",
							"vswitche_id":       "${alicloud_vswitch.default.id}",
							"security_group_id": "${alicloud_security_group.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"authorization_type": "API_KEY_AUTH",
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token",
									"api_key_value": "Token-value",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"username": "admin",
									"password": "admin",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"http_method":            "POST",
									"client_parameters": []map[string]interface{}{
										{
											"client_id":     "ClientId",
											"client_secret": "ClientSecret",
										},
									},
									"oauth_http_parameters": []map[string]interface{}{
										{
											"header_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
											"body_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_parameters.#": "1",
					}),
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
						{
							"network_type":      "PublicNetwork",
							"vpc_id":            "${alicloud_vpc.default.id}",
							"vswitche_id":       "${alicloud_vswitch.default.id}",
							"security_group_id": "${alicloud_security_group.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"authorization_type": "BASIC_AUTH",
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token-update",
									"api_key_value": "Token-value-update",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"username": "admin-update",
									"password": "admin-update",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"http_method":            "POST",
									"client_parameters": []map[string]interface{}{
										{
											"client_id":     "clientId",
											"client_secret": "clientSecret",
										},
									},
									"oauth_http_parameters": []map[string]interface{}{
										{
											"header_parameters": []map[string]interface{}{
												{
													"key":             "name-update",
													"value":           "name-update",
													"is_value_secret": "false",
												},
											},
											"body_parameters": []map[string]interface{}{
												{
													"key":             "name-update",
													"value":           "name-update",
													"is_value_secret": "false",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"key":             "name-update",
													"value":           "name-update",
													"is_value_secret": "false",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_parameters": []map[string]interface{}{
						{
							"authorization_type": "OAUTH_AUTH",
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token",
									"api_key_value": "Token-value",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"username": "admin",
									"password": "admin",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"http_method":            "POST",
									"client_parameters": []map[string]interface{}{
										{
											"client_id":     "ClientId",
											"client_secret": "ClientSecret",
										},
									},
									"oauth_http_parameters": []map[string]interface{}{
										{
											"header_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
											"body_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_parameters.#": "1",
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

func TestAccAliCloudEventBridgeConnection_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	resourceId := "alicloud_event_bridge_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeConnectionBasicDependence0)
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
							"network_type":      "PrivateNetwork",
							"vpc_id":            "${alicloud_vpc.default.id}",
							"vswitche_id":       "${alicloud_vswitch.default.id}",
							"security_group_id": "${alicloud_security_group.default.id}",
						},
					},
					"auth_parameters": []map[string]interface{}{
						{
							"authorization_type": "API_KEY_AUTH",
							"api_key_auth_parameters": []map[string]interface{}{
								{
									"api_key_name":  "Token",
									"api_key_value": "Token-value",
								},
							},
							"basic_auth_parameters": []map[string]interface{}{
								{
									"username": "admin",
									"password": "admin",
								},
							},
							"oauth_parameters": []map[string]interface{}{
								{
									"authorization_endpoint": "http://127.0.0.1:8080",
									"http_method":            "POST",
									"client_parameters": []map[string]interface{}{
										{
											"client_id":     "ClientId",
											"client_secret": "ClientSecret",
										},
									},
									"oauth_http_parameters": []map[string]interface{}{
										{
											"header_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
											"body_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
											"query_string_parameters": []map[string]interface{}{
												{
													"key":             "name",
													"value":           "name",
													"is_value_secret": "true",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_name":      name,
						"description":          "test-connection-basic-pre-update-oauth",
						"network_parameters.#": "1",
						"auth_parameters.#":    "1",
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

var AliCloudEventBridgeConnectionMap0 = map[string]string{
	"create_time": CHECKSET,
}

func AliCloudEventBridgeConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones[0].id
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vswitch.default.vpc_id
	}
`, name)
}

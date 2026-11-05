package eventbridge_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eventbridge"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The following TypeList attributes are pre-existing stock fields of
// alicloud_event_bridge_connection migrated as-is from the root alicloud
// package; this migration does not add or change their schema. Convergent
// reorder acceptance coverage is tracked separately and is outside this
// migration change, so they are marked as known stock debt per
// scripts/collection-order/README.md.
// collection-order: bypass auth_parameters.oauth_parameters.oauth_http_parameters.body_parameters -- pre-existing attribute outside this migration change
// collection-order: bypass auth_parameters.oauth_parameters.oauth_http_parameters.header_parameters -- pre-existing attribute outside this migration change
// collection-order: bypass auth_parameters.oauth_parameters.oauth_http_parameters.query_string_parameters -- pre-existing attribute outside this migration change

func TestAccAliCloudEventBridgeConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	ResourceId := "alicloud_event_bridge_connection.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeConnectionMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventBridgeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeconnection%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AliCloudEventBridgeConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
				ResourceName:      ResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudEventBridgeConnection_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.EventBridgeConnectionSupportRegions)
	ResourceId := "alicloud_event_bridge_connection.default"
	ra := acctest.ResourceAttrInit(ResourceId, AliCloudEventBridgeConnectionMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(ResourceId, &v, func() interface{} {
		return eventbridge.NewEventBridgeServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEventBridgeConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgeconnection%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(ResourceId, name, AliCloudEventBridgeConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
				ResourceName:      ResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudEventBridgeConnectionMap0 = map[string]string{
	"create_time": acctest.CHECKSET,
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

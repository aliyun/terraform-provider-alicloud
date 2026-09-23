package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketWebsite. >>> Resource test cases, automatically generated.
// Case 测试AliCDN回源 9176
func TestAccAliCloudOssBucketWebsite_basic9176(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_website.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketWebsiteMap9176)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketWebsite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketwebsite%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketWebsiteBasicDependence9176)
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
					"index_document": []map[string]interface{}{
						{
							"suffix":          "index.html",
							"support_sub_dir": "true",
							"type":            "0",
						},
					},
					"error_document": []map[string]interface{}{
						{
							"key":         "error.html",
							"http_status": "404",
						},
					},
					"bucket": "${alicloud_oss_bucket.defaultnVj9x3.bucket}",
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "1",
									"condition": []map[string]interface{}{
										{
											"http_error_code_returned_equals": "404",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type":      "AliCDN",
											"host_name":          "www.alicdn-master.com",
											"protocol":           "https",
											"http_redirect_code": "305",
										},
									},
									"lua_config": []map[string]interface{}{
										{
											"script": "test.lua",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "2",
									"condition": []map[string]interface{}{
										{
											"http_error_code_returned_equals": "405",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type":      "AliCDN",
											"protocol":           "http",
											"host_name":          "www.alicdn-slave.com",
											"http_redirect_code": "303",
										},
									},
								},
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
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "1",
									"condition": []map[string]interface{}{
										{
											"key_prefix_equals": "1",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type": "AliCDN",
											"protocol":      "http",
											"host_name":     "www.alicdn-test.com",
										},
									},
								},
								{
									"rule_number": "2",
									"condition": []map[string]interface{}{
										{
											"key_prefix_equals": "2",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type": "AliCDN",
											"protocol":      "https",
											"host_name":     "www.alicdn-test.com",
										},
									},
								},
								{
									"rule_number": "3",
									"condition": []map[string]interface{}{
										{
											"key_prefix_equals": "3",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"protocol":      "http",
											"host_name":     "www.alicdn-test.com",
											"redirect_type": "AliCDN",
										},
									},
								},
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
					"index_document": []map[string]interface{}{
						{
							"suffix":          "index.html",
							"support_sub_dir": "true",
							"type":            "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudOssBucketWebsiteMap9176 = map[string]string{}

func AlicloudOssBucketWebsiteBasicDependence9176(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaultnVj9x3" {
  storage_class = "Standard"
  lifecycle {
	ignore_changes = [website]
  }
}


`, name)
}

// Case 测试静态主页与镜像回源 9119
func TestAccAliCloudOssBucketWebsite_basic9119(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_website.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketWebsiteMap9119)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketWebsite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketwebsite%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketWebsiteBasicDependence9119)
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
					"index_document": []map[string]interface{}{
						{
							"suffix":          "index.html",
							"support_sub_dir": "true",
							"type":            "0",
						},
					},
					"bucket": "${alicloud_oss_bucket.defaultnVj9x3.bucket}",
					"error_document": []map[string]interface{}{
						{
							"key":         "error.html",
							"http_status": "404",
						},
					},
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "1",
									"condition": []map[string]interface{}{
										{
											"key_prefix_equals":               "test-prefix",
											"key_suffix_equals":               ".target",
											"http_error_code_returned_equals": "404",
											"include_headers": []map[string]interface{}{
												{
													"key":         "test-key-1",
													"equals":      "test-value-1",
													"starts_with": "test-",
													"ends_with":   "-test",
												},
											},
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type":            "Mirror",
											"pass_query_string":        "true",
											"mirror_pass_query_string": "true",
											"mirror_follow_redirect":   "true",
											"mirror_check_md5":         "true",
											"mirror_headers": []map[string]interface{}{
												{
													"pass_all": "true",
													"pass": []string{
														"pass-header"},
													"remove": []string{
														"remove-header"},
													"set": []map[string]interface{}{
														{
															"key":   "add-header",
															"value": "add-value",
														},
													},
												},
											},
											"host_name":                         "www.test-mirror:8900",
											"replace_key_prefix_with":           "abc-",
											"enable_replace_prefix":             "true",
											"replace_key_with":                  "$${key}-def",
											"mirror_pass_original_slashes":      "true",
											"mirror_url_slave":                  "www.mirror-slave.com",
											"mirror_url_probe":                  "www.mirror-probe.com",
											"mirror_save_oss_meta":              "true",
											"mirror_proxy_pass":                 "false",
											"mirror_allow_get_image_info":       "true",
											"mirror_allow_video_snapshot":       "true",
											"mirror_is_express_tunnel":          "false",
											"mirror_user_last_modified":         "true",
											"mirror_switch_all_errors":          "true",
											"mirror_using_role":                 "true",
											"mirror_role":                       "role-abc",
											"mirror_allow_head_object":          "true",
											"transparent_mirror_response_codes": "404,424",
											"mirror_async_status":               "0",
											"mirror_taggings": []map[string]interface{}{
												{
													"taggings": []map[string]interface{}{
														{
															"key":   "add-tag",
															"value": "tag-balue",
														},
													},
												},
											},
											"mirror_return_headers": []map[string]interface{}{
												{
													"return_header": []map[string]interface{}{
														{
															"key":   "test-eahd",
															"value": "test-value",
														},
													},
												},
											},
											"mirror_auth": []map[string]interface{}{
												{
													"auth_type":         "1",
													"region":            "ap-southeat-1",
													"access_key_id":     "testak",
													"access_key_secret": "testsk",
												},
											},
											"mirror_multi_alternates": []map[string]interface{}{
												{
													"mirror_multi_alternate": []map[string]interface{}{
														{
															"mirror_multi_alternate_number":     "2",
															"mirror_multi_alternate_url":        "http://www.alternate-mirror.com",
															"mirror_multi_alternate_vpc_id":     "vpc-xxx",
															"mirror_multi_alternate_dst_region": "oss-cn-beijing",
														},
													},
												},
											},
											"mirror_url": "http://test-mirror-url.com",
											"mirror_sni": "true",
										},
									},
									"lua_config": []map[string]interface{}{
										{
											"script": "test.lua",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"index_document": []map[string]interface{}{
						{
							"suffix":          "other-index.html",
							"support_sub_dir": "false",
							"type":            "2",
						},
					},
					"error_document": []map[string]interface{}{
						{
							"key":         "other-error.html",
							"http_status": "200",
						},
					},
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "2",
									"condition": []map[string]interface{}{
										{
											"include_headers": []map[string]interface{}{
												{
													"key":         "x-oss-img-process",
													"equals":      "abc-test",
													"starts_with": "starts-with-this",
													"ends_with":   "ends-with-this",
												},
											},
											"http_error_code_returned_equals": "404",
											"key_prefix_equals":               "other-prefix",
											"key_suffix_equals":               "other-suffix",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type":     "Mirror",
											"pass_query_string": "false",
											"mirror_headers": []map[string]interface{}{
												{
													"set": []map[string]interface{}{
														{
															"key":   "other-add-key",
															"value": "other-add-value",
														},
													},
													"pass_all": "false",
													"pass": []string{
														"other-header"},
													"remove": []string{
														"othear-remove-header"},
												},
											},
											"replace_key_prefix_with": "$${key}_L.jpg",
											"enable_replace_prefix":   "false",
											"mirror_taggings": []map[string]interface{}{
												{
													"taggings": []map[string]interface{}{
														{
															"key":   "other-tag",
															"value": "other-tag-value",
														},
													},
												},
											},
											"mirror_return_headers": []map[string]interface{}{
												{
													"return_header": []map[string]interface{}{
														{
															"key":   "other-header",
															"value": "other-header-value",
														},
													},
												},
											},
											"mirror_multi_alternates": []map[string]interface{}{
												{
													"mirror_multi_alternate": []map[string]interface{}{
														{
															"mirror_multi_alternate_number":     "1",
															"mirror_multi_alternate_url":        "http://www.another-mirror.com/",
															"mirror_multi_alternate_vpc_id":     "vpc-testxxx",
															"mirror_multi_alternate_dst_region": "oss-cn-hangzhou",
														},
													},
												},
											},
											"mirror_pass_query_string":          "false",
											"mirror_follow_redirect":            "false",
											"mirror_check_md5":                  "false",
											"host_name":                         "www.test-bucket.com",
											"replace_key_with":                  "xyz",
											"mirror_save_oss_meta":              "false",
											"mirror_proxy_pass":                 "false",
											"mirror_allow_get_image_info":       "false",
											"mirror_allow_video_snapshot":       "false",
											"mirror_is_express_tunnel":          "false",
											"mirror_user_last_modified":         "false",
											"mirror_switch_all_errors":          "false",
											"mirror_using_role":                 "true",
											"mirror_allow_head_object":          "false",
											"transparent_mirror_response_codes": "404",
											"mirror_async_status":               "0",
											"mirror_sni":                        "false",
											"mirror_url_slave":                  "www.aliyuncs.com/slave",
											"mirror_url_probe":                  "www.aliyuncs.com/hearbeat",
											"mirror_pass_original_slashes":      "false",
											"mirror_url":                        "https://other-mirror-test.com",
											"mirror_auth": []map[string]interface{}{
												{
													"auth_type":         "s3v4",
													"region":            "ap-southeats-2",
													"access_key_id":     "test-ak",
													"access_key_secret": "test-sk",
												},
											},
											"mirror_role": "oss-abc",
										},
									},
									"lua_config": []map[string]interface{}{
										{
											"script": "new.lua",
										},
									},
								},
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
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "1",
									"condition": []map[string]interface{}{
										{
											"key_prefix_equals":               "abc-",
											"http_error_code_returned_equals": "404",
										},
									},
									"redirect": []map[string]interface{}{
										{
											"redirect_type": "Mirror",
											"mirror_headers": []map[string]interface{}{
												{
													"pass_all": "true",
												},
											},
											"mirror_using_role":        "false",
											"pass_query_string":        "true",
											"mirror_pass_query_string": "true",
											"mirror_follow_redirect":   "false",
											"mirror_check_md5":         "false",
											"mirror_sni":               "false",
											"mirror_url":               "https://origin-mirror.com/",
										},
									},
								},
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
					"routing_rules": []map[string]interface{}{
						{
							"routing_rule": []map[string]interface{}{
								{
									"rule_number": "1",
									"redirect": []map[string]interface{}{
										{
											"redirect_type":      "AliCDN",
											"protocol":           "http",
											"http_redirect_code": "302",
											"host_name":          "www.aliyuncs.com",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudOssBucketWebsiteMap9119 = map[string]string{}

func AlicloudOssBucketWebsiteBasicDependence9119(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaultnVj9x3" {
  storage_class = "Standard"
  lifecycle {
	ignore_changes = [website]
  }
}


`, name)
}

// Test Oss BucketWebsite. <<< Resource test cases, automatically generated.

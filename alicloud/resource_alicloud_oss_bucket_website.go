// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketWebsite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketWebsiteCreate,
		Read:   resourceAliCloudOssBucketWebsiteRead,
		Update: resourceAliCloudOssBucketWebsiteUpdate,
		Delete: resourceAliCloudOssBucketWebsiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"error_document": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_status": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntAtLeast(0),
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"index_document": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"suffix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"support_sub_dir": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"routing_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"routing_rule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"include_headers": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"equals": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"starts_with": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"key": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"ends_with": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"key_prefix_equals": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"http_error_code_returned_equals": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"key_suffix_equals": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"lua_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"script": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"rule_number": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"redirect": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mirror_url": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_headers": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"pass_all": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"set": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"key": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"pass": {
																Type:     schema.TypeList,
																Optional: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"remove": {
																Type:     schema.TypeList,
																Optional: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"http_redirect_code": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"enable_replace_prefix": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_allow_get_image_info": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_dst_vpc_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"pass_query_string": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"transparent_mirror_response_codes": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_allow_head_object": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_pass_query_string": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_proxy_pass": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_multi_alternates": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mirror_multi_alternate": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mirror_multi_alternate_dst_region": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"mirror_multi_alternate_number": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																		"mirror_multi_alternate_url": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"mirror_multi_alternate_vpc_id": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"mirror_async_status": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_user_last_modified": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"host_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_sni": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_is_express_tunnel": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"replace_key_with": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_auth": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"access_key_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"access_key_secret": {
																Type:      schema.TypeString,
																Optional:  true,
																Sensitive: true,
															},
															"region": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"auth_type": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"mirror_allow_video_snapshot": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_tunnel_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_using_role": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_dst_region": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_return_headers": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"return_header": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"key": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"mirror_dst_slave_vpc_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_pass_original_slashes": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_url_probe": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_taggings": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"taggings": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"key": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"mirror_follow_redirect": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"replace_key_prefix_with": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"redirect_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_url_slave": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mirror_save_oss_meta": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_switch_all_errors": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_check_md5": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"mirror_role": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudOssBucketWebsiteCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?website")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("index_document"); !IsNil(v) {
		indexDocument := make(map[string]interface{})
		suffix1, _ := jsonpath.Get("$[0].suffix", d.Get("index_document"))
		if suffix1 != nil && suffix1 != "" {
			indexDocument["Suffix"] = suffix1
		}
		supportSubDir1, _ := jsonpath.Get("$[0].support_sub_dir", d.Get("index_document"))
		if supportSubDir1 != nil && supportSubDir1 != "" {
			indexDocument["SupportSubDir"] = supportSubDir1
		}
		type1, _ := jsonpath.Get("$[0].type", d.Get("index_document"))
		if type1 != nil && type1 != "" {
			indexDocument["Type"] = type1
		}

		objectDataLocalMap["IndexDocument"] = indexDocument
	}

	if v := d.Get("error_document"); !IsNil(v) {
		errorDocument := make(map[string]interface{})
		httpStatus1, _ := jsonpath.Get("$[0].http_status", d.Get("error_document"))
		if httpStatus1 != nil && httpStatus1 != "" {
			errorDocument["HttpStatus"] = httpStatus1
		}
		key1, _ := jsonpath.Get("$[0].key", d.Get("error_document"))
		if key1 != nil && key1 != "" {
			errorDocument["Key"] = key1
		}

		objectDataLocalMap["ErrorDocument"] = errorDocument
	}

	if v := d.Get("routing_rules"); !IsNil(v) {
		routingRules := make(map[string]interface{})
		if v, ok := d.GetOk("routing_rules"); ok {
			localData, err := jsonpath.Get("$[0].routing_rule", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["RuleNumber"] = dataLoopTmp["rule_number"]
				localData1 := make(map[string]interface{})
				keyPrefixEquals1, _ := jsonpath.Get("$[0].key_prefix_equals", dataLoopTmp["condition"])
				if keyPrefixEquals1 != nil && keyPrefixEquals1 != "" {
					localData1["KeyPrefixEquals"] = keyPrefixEquals1
				}
				keySuffixEquals1, _ := jsonpath.Get("$[0].key_suffix_equals", dataLoopTmp["condition"])
				if keySuffixEquals1 != nil && keySuffixEquals1 != "" {
					localData1["KeySuffixEquals"] = keySuffixEquals1
				}
				httpErrorCodeReturnedEquals1, _ := jsonpath.Get("$[0].http_error_code_returned_equals", dataLoopTmp["condition"])
				if httpErrorCodeReturnedEquals1 != nil && httpErrorCodeReturnedEquals1 != "" {
					localData1["HttpErrorCodeReturnedEquals"] = httpErrorCodeReturnedEquals1
				}
				if v, ok := dataLoopTmp["condition"]; ok {
					localData2, err := jsonpath.Get("$[0].include_headers", v)
					if err != nil {
						localData2 = make([]interface{}, 0)
					}
					localMaps2 := make([]interface{}, 0)
					for _, dataLoop2 := range localData2.([]interface{}) {
						dataLoop2Tmp := make(map[string]interface{})
						if dataLoop2 != nil {
							dataLoop2Tmp = dataLoop2.(map[string]interface{})
						}
						dataLoop2Map := make(map[string]interface{})
						dataLoop2Map["Key"] = dataLoop2Tmp["key"]
						dataLoop2Map["Equals"] = dataLoop2Tmp["equals"]
						dataLoop2Map["StartsWith"] = dataLoop2Tmp["starts_with"]
						dataLoop2Map["EndsWith"] = dataLoop2Tmp["ends_with"]
						localMaps2 = append(localMaps2, dataLoop2Map)
					}
					localData1["IncludeHeader"] = localMaps2
				}

				dataLoopMap["Condition"] = localData1
				localData3 := make(map[string]interface{})
				redirectType1, _ := jsonpath.Get("$[0].redirect_type", dataLoopTmp["redirect"])
				if redirectType1 != nil && redirectType1 != "" {
					localData3["RedirectType"] = redirectType1
				}
				passQueryString1, _ := jsonpath.Get("$[0].pass_query_string", dataLoopTmp["redirect"])
				if passQueryString1 != nil && passQueryString1 != "" {
					localData3["PassQueryString"] = passQueryString1
				}
				mirrorUrl, _ := jsonpath.Get("$[0].mirror_url", dataLoopTmp["redirect"])
				if mirrorUrl != nil && mirrorUrl != "" {
					localData3["MirrorURL"] = mirrorUrl
				}
				mirrorSni, _ := jsonpath.Get("$[0].mirror_sni", dataLoopTmp["redirect"])
				if mirrorSni != nil && mirrorSni != "" {
					localData3["MirrorSNI"] = mirrorSni
				}
				mirrorPassQueryString1, _ := jsonpath.Get("$[0].mirror_pass_query_string", dataLoopTmp["redirect"])
				if mirrorPassQueryString1 != nil && mirrorPassQueryString1 != "" {
					localData3["MirrorPassQueryString"] = mirrorPassQueryString1
				}
				mirrorFollowRedirect1, _ := jsonpath.Get("$[0].mirror_follow_redirect", dataLoopTmp["redirect"])
				if mirrorFollowRedirect1 != nil && mirrorFollowRedirect1 != "" {
					localData3["MirrorFollowRedirect"] = mirrorFollowRedirect1
				}
				mirrorCheckMd51, _ := jsonpath.Get("$[0].mirror_check_md5", dataLoopTmp["redirect"])
				if mirrorCheckMd51 != nil && mirrorCheckMd51 != "" {
					localData3["MirrorCheckMd5"] = mirrorCheckMd51
				}
				mirrorHeaders := make(map[string]interface{})
				passAll1, _ := jsonpath.Get("$[0].mirror_headers[0].pass_all", dataLoopTmp["redirect"])
				if passAll1 != nil && passAll1 != "" {
					mirrorHeaders["PassAll"] = passAll1
				}
				pass1, _ := jsonpath.Get("$[0].mirror_headers[0].pass", dataLoopTmp["redirect"])
				if pass1 != nil && pass1 != "" {
					mirrorHeaders["Pass"] = pass1
				}
				remove1, _ := jsonpath.Get("$[0].mirror_headers[0].remove", dataLoopTmp["redirect"])
				if remove1 != nil && remove1 != "" {
					mirrorHeaders["Remove"] = remove1
				}
				if v, ok := dataLoopTmp["redirect"]; ok {
					localData4, err := jsonpath.Get("$[0].mirror_headers[0].set", v)
					if err != nil {
						localData4 = make([]interface{}, 0)
					}
					localMaps4 := make([]interface{}, 0)
					for _, dataLoop4 := range localData4.([]interface{}) {
						dataLoop4Tmp := make(map[string]interface{})
						if dataLoop4 != nil {
							dataLoop4Tmp = dataLoop4.(map[string]interface{})
						}
						dataLoop4Map := make(map[string]interface{})
						dataLoop4Map["Value"] = dataLoop4Tmp["value"]
						dataLoop4Map["Key"] = dataLoop4Tmp["key"]
						localMaps4 = append(localMaps4, dataLoop4Map)
					}
					mirrorHeaders["Set"] = localMaps4
				}

				localData3["MirrorHeaders"] = mirrorHeaders
				protocol1, _ := jsonpath.Get("$[0].protocol", dataLoopTmp["redirect"])
				if protocol1 != nil && protocol1 != "" {
					localData3["Protocol"] = protocol1
				}
				hostName1, _ := jsonpath.Get("$[0].host_name", dataLoopTmp["redirect"])
				if hostName1 != nil && hostName1 != "" {
					localData3["HostName"] = hostName1
				}
				replaceKeyPrefixWith1, _ := jsonpath.Get("$[0].replace_key_prefix_with", dataLoopTmp["redirect"])
				if replaceKeyPrefixWith1 != nil && replaceKeyPrefixWith1 != "" {
					localData3["ReplaceKeyPrefixWith"] = replaceKeyPrefixWith1
				}
				enableReplacePrefix1, _ := jsonpath.Get("$[0].enable_replace_prefix", dataLoopTmp["redirect"])
				if enableReplacePrefix1 != nil && enableReplacePrefix1 != "" {
					localData3["EnableReplacePrefix"] = enableReplacePrefix1
				}
				replaceKeyWith1, _ := jsonpath.Get("$[0].replace_key_with", dataLoopTmp["redirect"])
				if replaceKeyWith1 != nil && replaceKeyWith1 != "" {
					localData3["ReplaceKeyWith"] = replaceKeyWith1
				}
				httpRedirectCode1, _ := jsonpath.Get("$[0].http_redirect_code", dataLoopTmp["redirect"])
				if httpRedirectCode1 != nil && httpRedirectCode1 != "" {
					localData3["HttpRedirectCode"] = httpRedirectCode1
				}
				mirrorPassOriginalSlashes1, _ := jsonpath.Get("$[0].mirror_pass_original_slashes", dataLoopTmp["redirect"])
				if mirrorPassOriginalSlashes1 != nil && mirrorPassOriginalSlashes1 != "" {
					localData3["MirrorPassOriginalSlashes"] = mirrorPassOriginalSlashes1
				}
				mirrorUrlSlave, _ := jsonpath.Get("$[0].mirror_url_slave", dataLoopTmp["redirect"])
				if mirrorUrlSlave != nil && mirrorUrlSlave != "" {
					localData3["MirrorURLSlave"] = mirrorUrlSlave
				}
				mirrorUrlProbe, _ := jsonpath.Get("$[0].mirror_url_probe", dataLoopTmp["redirect"])
				if mirrorUrlProbe != nil && mirrorUrlProbe != "" {
					localData3["MirrorURLProbe"] = mirrorUrlProbe
				}
				mirrorSaveOssMeta1, _ := jsonpath.Get("$[0].mirror_save_oss_meta", dataLoopTmp["redirect"])
				if mirrorSaveOssMeta1 != nil && mirrorSaveOssMeta1 != "" {
					localData3["MirrorSaveOssMeta"] = mirrorSaveOssMeta1
				}
				mirrorProxyPass1, _ := jsonpath.Get("$[0].mirror_proxy_pass", dataLoopTmp["redirect"])
				if mirrorProxyPass1 != nil && mirrorProxyPass1 != "" {
					localData3["MirrorProxyPass"] = mirrorProxyPass1
				}
				mirrorAllowGetImageInfo1, _ := jsonpath.Get("$[0].mirror_allow_get_image_info", dataLoopTmp["redirect"])
				if mirrorAllowGetImageInfo1 != nil && mirrorAllowGetImageInfo1 != "" {
					localData3["MirrorAllowGetImageInfo"] = mirrorAllowGetImageInfo1
				}
				mirrorAllowVideoSnapshot1, _ := jsonpath.Get("$[0].mirror_allow_video_snapshot", dataLoopTmp["redirect"])
				if mirrorAllowVideoSnapshot1 != nil && mirrorAllowVideoSnapshot1 != "" {
					localData3["MirrorAllowVideoSnapshot"] = mirrorAllowVideoSnapshot1
				}
				mirrorIsExpressTunnel1, _ := jsonpath.Get("$[0].mirror_is_express_tunnel", dataLoopTmp["redirect"])
				if mirrorIsExpressTunnel1 != nil && mirrorIsExpressTunnel1 != "" {
					localData3["MirrorIsExpressTunnel"] = mirrorIsExpressTunnel1
				}
				mirrorDstRegion1, _ := jsonpath.Get("$[0].mirror_dst_region", dataLoopTmp["redirect"])
				if mirrorDstRegion1 != nil && mirrorDstRegion1 != "" {
					localData3["MirrorDstRegion"] = mirrorDstRegion1
				}
				mirrorDstVpcId1, _ := jsonpath.Get("$[0].mirror_dst_vpc_id", dataLoopTmp["redirect"])
				if mirrorDstVpcId1 != nil && mirrorDstVpcId1 != "" {
					localData3["MirrorDstVpcId"] = mirrorDstVpcId1
				}
				mirrorDstSlaveVpcId1, _ := jsonpath.Get("$[0].mirror_dst_slave_vpc_id", dataLoopTmp["redirect"])
				if mirrorDstSlaveVpcId1 != nil && mirrorDstSlaveVpcId1 != "" {
					localData3["MirrorDstSlaveVpcId"] = mirrorDstSlaveVpcId1
				}
				mirrorUserLastModified1, _ := jsonpath.Get("$[0].mirror_user_last_modified", dataLoopTmp["redirect"])
				if mirrorUserLastModified1 != nil && mirrorUserLastModified1 != "" {
					localData3["MirrorUserLastModified"] = mirrorUserLastModified1
				}
				mirrorSwitchAllErrors1, _ := jsonpath.Get("$[0].mirror_switch_all_errors", dataLoopTmp["redirect"])
				if mirrorSwitchAllErrors1 != nil && mirrorSwitchAllErrors1 != "" {
					localData3["MirrorSwitchAllErrors"] = mirrorSwitchAllErrors1
				}
				mirrorTunnelId1, _ := jsonpath.Get("$[0].mirror_tunnel_id", dataLoopTmp["redirect"])
				if mirrorTunnelId1 != nil && mirrorTunnelId1 != "" {
					localData3["MirrorTunnelId"] = mirrorTunnelId1
				}
				mirrorUsingRole1, _ := jsonpath.Get("$[0].mirror_using_role", dataLoopTmp["redirect"])
				if mirrorUsingRole1 != nil && mirrorUsingRole1 != "" {
					localData3["MirrorUsingRole"] = mirrorUsingRole1
				}
				mirrorRole1, _ := jsonpath.Get("$[0].mirror_role", dataLoopTmp["redirect"])
				if mirrorRole1 != nil && mirrorRole1 != "" {
					localData3["MirrorRole"] = mirrorRole1
				}
				mirrorAllowHeadObject1, _ := jsonpath.Get("$[0].mirror_allow_head_object", dataLoopTmp["redirect"])
				if mirrorAllowHeadObject1 != nil && mirrorAllowHeadObject1 != "" {
					localData3["MirrorAllowHeadObject"] = mirrorAllowHeadObject1
				}
				transparentMirrorResponseCodes1, _ := jsonpath.Get("$[0].transparent_mirror_response_codes", dataLoopTmp["redirect"])
				if transparentMirrorResponseCodes1 != nil && transparentMirrorResponseCodes1 != "" {
					localData3["TransparentMirrorResponseCodes"] = transparentMirrorResponseCodes1
				}
				mirrorAsyncStatus1, _ := jsonpath.Get("$[0].mirror_async_status", dataLoopTmp["redirect"])
				if mirrorAsyncStatus1 != nil && mirrorAsyncStatus1 != "" {
					localData3["MirrorAsyncStatus"] = mirrorAsyncStatus1
				}
				mirrorTaggings := make(map[string]interface{})
				if v, ok := dataLoopTmp["redirect"]; ok {
					localData5, err := jsonpath.Get("$[0].mirror_taggings[0].taggings", v)
					if err != nil {
						localData5 = make([]interface{}, 0)
					}
					localMaps5 := make([]interface{}, 0)
					for _, dataLoop5 := range localData5.([]interface{}) {
						dataLoop5Tmp := make(map[string]interface{})
						if dataLoop5 != nil {
							dataLoop5Tmp = dataLoop5.(map[string]interface{})
						}
						dataLoop5Map := make(map[string]interface{})
						dataLoop5Map["Key"] = dataLoop5Tmp["key"]
						dataLoop5Map["Value"] = dataLoop5Tmp["value"]
						localMaps5 = append(localMaps5, dataLoop5Map)
					}
					mirrorTaggings["Taggings"] = localMaps5
				}

				localData3["MirrorTaggings"] = mirrorTaggings
				mirrorReturnHeaders := make(map[string]interface{})
				if v, ok := dataLoopTmp["redirect"]; ok {
					localData6, err := jsonpath.Get("$[0].mirror_return_headers[0].return_header", v)
					if err != nil {
						localData6 = make([]interface{}, 0)
					}
					localMaps6 := make([]interface{}, 0)
					for _, dataLoop6 := range localData6.([]interface{}) {
						dataLoop6Tmp := make(map[string]interface{})
						if dataLoop6 != nil {
							dataLoop6Tmp = dataLoop6.(map[string]interface{})
						}
						dataLoop6Map := make(map[string]interface{})
						dataLoop6Map["Value"] = dataLoop6Tmp["value"]
						dataLoop6Map["Key"] = dataLoop6Tmp["key"]
						localMaps6 = append(localMaps6, dataLoop6Map)
					}
					mirrorReturnHeaders["ReturnHeader"] = localMaps6
				}

				localData3["MirrorReturnHeaders"] = mirrorReturnHeaders
				mirrorAuth := make(map[string]interface{})
				authType1, _ := jsonpath.Get("$[0].mirror_auth[0].auth_type", dataLoopTmp["redirect"])
				if authType1 != nil && authType1 != "" {
					mirrorAuth["AuthType"] = authType1
				}
				region1, _ := jsonpath.Get("$[0].mirror_auth[0].region", dataLoopTmp["redirect"])
				if region1 != nil && region1 != "" {
					mirrorAuth["Region"] = region1
				}
				accessKeyId1, _ := jsonpath.Get("$[0].mirror_auth[0].access_key_id", dataLoopTmp["redirect"])
				if accessKeyId1 != nil && accessKeyId1 != "" {
					mirrorAuth["AccessKeyId"] = accessKeyId1
				}
				accessKeySecret1, _ := jsonpath.Get("$[0].mirror_auth[0].access_key_secret", dataLoopTmp["redirect"])
				if accessKeySecret1 != nil && accessKeySecret1 != "" {
					mirrorAuth["AccessKeySecret"] = accessKeySecret1
				}

				localData3["MirrorAuth"] = mirrorAuth
				mirrorMultiAlternates := make(map[string]interface{})
				if v, ok := dataLoopTmp["redirect"]; ok {
					localData7, err := jsonpath.Get("$[0].mirror_multi_alternates[0].mirror_multi_alternate", v)
					if err != nil {
						localData7 = make([]interface{}, 0)
					}
					localMaps7 := make([]interface{}, 0)
					for _, dataLoop7 := range localData7.([]interface{}) {
						dataLoop7Tmp := make(map[string]interface{})
						if dataLoop7 != nil {
							dataLoop7Tmp = dataLoop7.(map[string]interface{})
						}
						dataLoop7Map := make(map[string]interface{})
						dataLoop7Map["MirrorMultiAlternateNumber"] = dataLoop7Tmp["mirror_multi_alternate_number"]
						dataLoop7Map["MirrorMultiAlternateURL"] = dataLoop7Tmp["mirror_multi_alternate_url"]
						dataLoop7Map["MirrorMultiAlternateVpcId"] = dataLoop7Tmp["mirror_multi_alternate_vpc_id"]
						dataLoop7Map["MirrorMultiAlternateDstRegion"] = dataLoop7Tmp["mirror_multi_alternate_dst_region"]
						localMaps7 = append(localMaps7, dataLoop7Map)
					}
					mirrorMultiAlternates["MirrorMultiAlternate"] = localMaps7
				}

				localData3["MirrorMultiAlternates"] = mirrorMultiAlternates
				dataLoopMap["Redirect"] = localData3
				localData8 := make(map[string]interface{})
				script1, _ := jsonpath.Get("$[0].script", dataLoopTmp["lua_config"])
				if script1 != nil && script1 != "" {
					localData8["Script"] = script1
				}
				dataLoopMap["LuaConfig"] = localData8
				localMaps = append(localMaps, dataLoopMap)
			}
			routingRules["RoutingRule"] = localMaps
		}

		objectDataLocalMap["RoutingRules"] = routingRules
	}

	request["WebsiteConfiguration"] = objectDataLocalMap
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketWebsite", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_website", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketWebsiteRead(d, meta)
}

func resourceAliCloudOssBucketWebsiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketWebsite(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_website DescribeOssBucketWebsite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	accessKeySecretMap := make(map[string]interface{})
	if v := d.Get("routing_rules"); !IsNil(v) {
		if v, ok := d.GetOk("routing_rules"); ok {
			localData, err := jsonpath.Get("$[0].routing_rule", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
					ruleNumber := fmt.Sprint(dataLoopTmp["rule_number"])
					accessKeySecret1, _ := jsonpath.Get("$[0].mirror_auth[0].access_key_secret", dataLoopTmp["redirect"])
					if accessKeySecret1 != nil && accessKeySecret1 != "" {
						accessKeySecretMap[ruleNumber] = accessKeySecret1
					}
				}
			}
		}
	}
	errorDocumentMaps := make([]map[string]interface{}, 0)
	errorDocumentMap := make(map[string]interface{})
	errorDocument1Raw := make(map[string]interface{})
	if objectRaw["ErrorDocument"] != nil {
		errorDocument1Raw = objectRaw["ErrorDocument"].(map[string]interface{})
	}
	if len(errorDocument1Raw) > 0 {
		errorDocumentMap["http_status"] = errorDocument1Raw["HttpStatus"]
		errorDocumentMap["key"] = errorDocument1Raw["Key"]

		errorDocumentMaps = append(errorDocumentMaps, errorDocumentMap)
	}
	if objectRaw["ErrorDocument"] != nil {
		if err := d.Set("error_document", errorDocumentMaps); err != nil {
			return err
		}
	}
	indexDocumentMaps := make([]map[string]interface{}, 0)
	indexDocumentMap := make(map[string]interface{})
	indexDocument1Raw := make(map[string]interface{})
	if objectRaw["IndexDocument"] != nil {
		indexDocument1Raw = objectRaw["IndexDocument"].(map[string]interface{})
	}
	if len(indexDocument1Raw) > 0 {
		indexDocumentMap["suffix"] = indexDocument1Raw["Suffix"]
		indexDocumentMap["support_sub_dir"] = indexDocument1Raw["SupportSubDir"]
		indexDocumentMap["type"] = indexDocument1Raw["Type"]

		indexDocumentMaps = append(indexDocumentMaps, indexDocumentMap)
	}
	if objectRaw["IndexDocument"] != nil {
		if err := d.Set("index_document", indexDocumentMaps); err != nil {
			return err
		}
	}
	routingRulesMaps := make([]map[string]interface{}, 0)
	routingRulesMap := make(map[string]interface{})
	routingRule1Raw, _ := jsonpath.Get("$.RoutingRules.RoutingRule", objectRaw)

	routingRuleMaps := make([]map[string]interface{}, 0)
	if routingRule1Raw != nil {
		for _, routingRuleChild1Raw := range routingRule1Raw.([]interface{}) {
			routingRuleMap := make(map[string]interface{})
			routingRuleChild1Raw := routingRuleChild1Raw.(map[string]interface{})
			routingRuleMap["rule_number"] = routingRuleChild1Raw["RuleNumber"]

			conditionMaps := make([]map[string]interface{}, 0)
			conditionMap := make(map[string]interface{})
			condition1Raw := make(map[string]interface{})
			if routingRuleChild1Raw["Condition"] != nil {
				condition1Raw = routingRuleChild1Raw["Condition"].(map[string]interface{})
			}
			if len(condition1Raw) > 0 {
				conditionMap["http_error_code_returned_equals"] = condition1Raw["HttpErrorCodeReturnedEquals"]
				conditionMap["key_prefix_equals"] = condition1Raw["KeyPrefixEquals"]
				conditionMap["key_suffix_equals"] = condition1Raw["KeySuffixEquals"]

				includeHeader1Raw := condition1Raw["IncludeHeader"]
				includeHeadersMaps := make([]map[string]interface{}, 0)
				if includeHeader1Raw != nil {
					for _, includeHeaderChild1Raw := range includeHeader1Raw.([]interface{}) {
						includeHeadersMap := make(map[string]interface{})
						includeHeaderChild1Raw := includeHeaderChild1Raw.(map[string]interface{})
						includeHeadersMap["ends_with"] = includeHeaderChild1Raw["EndsWith"]
						includeHeadersMap["equals"] = includeHeaderChild1Raw["Equals"]
						includeHeadersMap["key"] = includeHeaderChild1Raw["Key"]
						includeHeadersMap["starts_with"] = includeHeaderChild1Raw["StartsWith"]

						includeHeadersMaps = append(includeHeadersMaps, includeHeadersMap)
					}
				}
				conditionMap["include_headers"] = includeHeadersMaps
				conditionMaps = append(conditionMaps, conditionMap)
			}
			routingRuleMap["condition"] = conditionMaps
			luaConfigMaps := make([]map[string]interface{}, 0)
			luaConfigMap := make(map[string]interface{})
			luaConfig1Raw := make(map[string]interface{})
			if routingRuleChild1Raw["LuaConfig"] != nil {
				luaConfig1Raw = routingRuleChild1Raw["LuaConfig"].(map[string]interface{})
			}
			if len(luaConfig1Raw) > 0 {
				luaConfigMap["script"] = luaConfig1Raw["Script"]

				luaConfigMaps = append(luaConfigMaps, luaConfigMap)
			}
			routingRuleMap["lua_config"] = luaConfigMaps
			redirectMaps := make([]map[string]interface{}, 0)
			redirectMap := make(map[string]interface{})
			redirect1Raw := make(map[string]interface{})
			if routingRuleChild1Raw["Redirect"] != nil {
				redirect1Raw = routingRuleChild1Raw["Redirect"].(map[string]interface{})
			}
			if len(redirect1Raw) > 0 {
				redirectMap["enable_replace_prefix"] = redirect1Raw["EnableReplacePrefix"]
				redirectMap["host_name"] = redirect1Raw["HostName"]
				redirectMap["http_redirect_code"] = redirect1Raw["HttpRedirectCode"]
				redirectMap["mirror_allow_get_image_info"] = redirect1Raw["MirrorAllowGetImageInfo"]
				redirectMap["mirror_allow_head_object"] = redirect1Raw["MirrorAllowHeadObject"]
				redirectMap["mirror_allow_video_snapshot"] = redirect1Raw["MirrorAllowVideoSnapshot"]
				redirectMap["mirror_async_status"] = redirect1Raw["MirrorAsyncStatus"]
				redirectMap["mirror_check_md5"] = redirect1Raw["MirrorCheckMd5"]
				redirectMap["mirror_dst_region"] = redirect1Raw["MirrorDstRegion"]
				redirectMap["mirror_dst_slave_vpc_id"] = redirect1Raw["MirrorDstSlaveVpcId"]
				redirectMap["mirror_dst_vpc_id"] = redirect1Raw["MirrorDstVpcId"]
				redirectMap["mirror_follow_redirect"] = redirect1Raw["MirrorFollowRedirect"]
				redirectMap["mirror_is_express_tunnel"] = redirect1Raw["MirrorIsExpressTunnel"]
				redirectMap["mirror_pass_original_slashes"] = redirect1Raw["MirrorPassOriginalSlashes"]
				redirectMap["mirror_pass_query_string"] = redirect1Raw["MirrorPassQueryString"]
				redirectMap["mirror_proxy_pass"] = redirect1Raw["MirrorProxyPass"]
				redirectMap["mirror_role"] = redirect1Raw["MirrorRole"]
				redirectMap["mirror_save_oss_meta"] = redirect1Raw["MirrorSaveOssMeta"]
				redirectMap["mirror_sni"] = redirect1Raw["MirrorSNI"]
				redirectMap["mirror_switch_all_errors"] = redirect1Raw["MirrorSwitchAllErrors"]
				redirectMap["mirror_tunnel_id"] = redirect1Raw["MirrorTunnelId"]
				redirectMap["mirror_url"] = redirect1Raw["MirrorURL"]
				redirectMap["mirror_url_probe"] = redirect1Raw["MirrorURLProbe"]
				redirectMap["mirror_url_slave"] = redirect1Raw["MirrorURLSlave"]
				redirectMap["mirror_user_last_modified"] = redirect1Raw["MirrorUserLastModified"]
				redirectMap["mirror_using_role"] = redirect1Raw["MirrorUsingRole"]
				redirectMap["pass_query_string"] = redirect1Raw["PassQueryString"]
				redirectMap["protocol"] = redirect1Raw["Protocol"]
				redirectMap["redirect_type"] = redirect1Raw["RedirectType"]
				redirectMap["replace_key_prefix_with"] = redirect1Raw["ReplaceKeyPrefixWith"]
				redirectMap["replace_key_with"] = redirect1Raw["ReplaceKeyWith"]
				redirectMap["transparent_mirror_response_codes"] = redirect1Raw["TransparentMirrorResponseCodes"]

				mirrorAuthMaps := make([]map[string]interface{}, 0)
				mirrorAuthMap := make(map[string]interface{})
				mirrorAuth1Raw := make(map[string]interface{})
				if redirect1Raw["MirrorAuth"] != nil {
					mirrorAuth1Raw = redirect1Raw["MirrorAuth"].(map[string]interface{})
				}
				if len(mirrorAuth1Raw) > 0 {
					mirrorAuthMap["access_key_id"] = mirrorAuth1Raw["AccessKeyId"]
					ruleNumber := fmt.Sprint(routingRuleMap["rule_number"])
					if v, ok := accessKeySecretMap[ruleNumber]; ok {
						mirrorAuthMap["access_key_secret"] = v
					}
					mirrorAuthMap["auth_type"] = mirrorAuth1Raw["AuthType"]
					mirrorAuthMap["region"] = mirrorAuth1Raw["Region"]

					mirrorAuthMaps = append(mirrorAuthMaps, mirrorAuthMap)
				}
				redirectMap["mirror_auth"] = mirrorAuthMaps
				mirrorHeadersMaps := make([]map[string]interface{}, 0)
				mirrorHeadersMap := make(map[string]interface{})
				mirrorHeaders1Raw := make(map[string]interface{})
				if redirect1Raw["MirrorHeaders"] != nil {
					mirrorHeaders1Raw = redirect1Raw["MirrorHeaders"].(map[string]interface{})
				}
				if len(mirrorHeaders1Raw) > 0 {
					mirrorHeadersMap["pass_all"] = mirrorHeaders1Raw["PassAll"]

					pass1Raw := make([]interface{}, 0)
					if mirrorHeaders1Raw["Pass"] != nil {
						pass1Raw = mirrorHeaders1Raw["Pass"].([]interface{})
					}

					mirrorHeadersMap["pass"] = pass1Raw
					remove1Raw := make([]interface{}, 0)
					if mirrorHeaders1Raw["Remove"] != nil {
						remove1Raw = mirrorHeaders1Raw["Remove"].([]interface{})
					}

					mirrorHeadersMap["remove"] = remove1Raw
					set1Raw := mirrorHeaders1Raw["Set"]
					setMaps := make([]map[string]interface{}, 0)
					if set1Raw != nil {
						for _, setChild1Raw := range set1Raw.([]interface{}) {
							setMap := make(map[string]interface{})
							setChild1Raw := setChild1Raw.(map[string]interface{})
							setMap["key"] = setChild1Raw["Key"]
							setMap["value"] = setChild1Raw["Value"]

							setMaps = append(setMaps, setMap)
						}
					}
					mirrorHeadersMap["set"] = setMaps
					mirrorHeadersMaps = append(mirrorHeadersMaps, mirrorHeadersMap)
				}
				redirectMap["mirror_headers"] = mirrorHeadersMaps
				mirrorMultiAlternatesMaps := make([]map[string]interface{}, 0)
				mirrorMultiAlternatesMap := make(map[string]interface{})
				mirrorMultiAlternate1Raw, _ := jsonpath.Get("$.Redirect.MirrorMultiAlternates.MirrorMultiAlternate", routingRuleChild1Raw)

				mirrorMultiAlternateMaps := make([]map[string]interface{}, 0)
				if mirrorMultiAlternate1Raw != nil {
					for _, mirrorMultiAlternateChild1Raw := range mirrorMultiAlternate1Raw.([]interface{}) {
						mirrorMultiAlternateMap := make(map[string]interface{})
						mirrorMultiAlternateChild1Raw := mirrorMultiAlternateChild1Raw.(map[string]interface{})
						mirrorMultiAlternateMap["mirror_multi_alternate_dst_region"] = mirrorMultiAlternateChild1Raw["MirrorMultiAlternateDstRegion"]
						mirrorMultiAlternateMap["mirror_multi_alternate_number"] = mirrorMultiAlternateChild1Raw["MirrorMultiAlternateNumber"]
						mirrorMultiAlternateMap["mirror_multi_alternate_url"] = mirrorMultiAlternateChild1Raw["MirrorMultiAlternateURL"]
						mirrorMultiAlternateMap["mirror_multi_alternate_vpc_id"] = mirrorMultiAlternateChild1Raw["MirrorMultiAlternateVpcId"]

						mirrorMultiAlternateMaps = append(mirrorMultiAlternateMaps, mirrorMultiAlternateMap)
					}
				}
				mirrorMultiAlternatesMap["mirror_multi_alternate"] = mirrorMultiAlternateMaps
				mirrorMultiAlternatesMaps = append(mirrorMultiAlternatesMaps, mirrorMultiAlternatesMap)
				redirectMap["mirror_multi_alternates"] = mirrorMultiAlternatesMaps
				mirrorReturnHeadersMaps := make([]map[string]interface{}, 0)
				mirrorReturnHeadersMap := make(map[string]interface{})
				returnHeader1Raw, _ := jsonpath.Get("$.Redirect.MirrorReturnHeaders.ReturnHeader", routingRuleChild1Raw)

				returnHeaderMaps := make([]map[string]interface{}, 0)
				if returnHeader1Raw != nil {
					for _, returnHeaderChild1Raw := range returnHeader1Raw.([]interface{}) {
						returnHeaderMap := make(map[string]interface{})
						returnHeaderChild1Raw := returnHeaderChild1Raw.(map[string]interface{})
						returnHeaderMap["key"] = returnHeaderChild1Raw["Key"]
						returnHeaderMap["value"] = returnHeaderChild1Raw["Value"]

						returnHeaderMaps = append(returnHeaderMaps, returnHeaderMap)
					}
				}
				mirrorReturnHeadersMap["return_header"] = returnHeaderMaps
				mirrorReturnHeadersMaps = append(mirrorReturnHeadersMaps, mirrorReturnHeadersMap)
				redirectMap["mirror_return_headers"] = mirrorReturnHeadersMaps
				mirrorTaggingsMaps := make([]map[string]interface{}, 0)
				mirrorTaggingsMap := make(map[string]interface{})
				taggings1Raw, _ := jsonpath.Get("$.Redirect.MirrorTaggings.Taggings", routingRuleChild1Raw)

				taggingsMaps := make([]map[string]interface{}, 0)
				if taggings1Raw != nil {
					for _, taggingsChild1Raw := range taggings1Raw.([]interface{}) {
						taggingsMap := make(map[string]interface{})
						taggingsChild1Raw := taggingsChild1Raw.(map[string]interface{})
						taggingsMap["key"] = taggingsChild1Raw["Key"]
						taggingsMap["value"] = taggingsChild1Raw["Value"]

						taggingsMaps = append(taggingsMaps, taggingsMap)
					}
				}
				mirrorTaggingsMap["taggings"] = taggingsMaps
				mirrorTaggingsMaps = append(mirrorTaggingsMaps, mirrorTaggingsMap)
				redirectMap["mirror_taggings"] = mirrorTaggingsMaps
				redirectMaps = append(redirectMaps, redirectMap)
			}
			routingRuleMap["redirect"] = redirectMaps
			routingRuleMaps = append(routingRuleMaps, routingRuleMap)
		}
	}
	routingRulesMap["routing_rule"] = routingRuleMaps
	routingRulesMaps = append(routingRulesMaps, routingRulesMap)
	if routingRule1Raw != nil {
		if err := d.Set("routing_rules", routingRulesMaps); err != nil {
			return err
		}
	}

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketWebsiteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	action := fmt.Sprintf("/?website")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("routing_rules") {
		update = true
	}
	if v := d.Get("routing_rules"); !IsNil(v) {
		routingRules := make(map[string]interface{})
		if v, ok := d.GetOk("routing_rules"); ok {
			localData, err := jsonpath.Get("$[0].routing_rule", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				if !IsNil(dataLoopTmp["lua_config"]) {
					localData1 := make(map[string]interface{})
					script1, _ := jsonpath.Get("$[0].script", dataLoopTmp["lua_config"])
					if script1 != nil && script1 != "" {
						localData1["Script"] = script1
					}
					dataLoopMap["LuaConfig"] = localData1
				}
				if !IsNil(dataLoopTmp["redirect"]) {
					localData2 := make(map[string]interface{})
					mirrorMultiAlternates := make(map[string]interface{})
					if v, ok := dataLoopTmp["redirect"]; ok {
						localData3, err := jsonpath.Get("$[0].mirror_multi_alternates[0].mirror_multi_alternate", v)
						if err != nil {
							localData3 = make([]interface{}, 0)
						}
						localMaps3 := make([]interface{}, 0)
						for _, dataLoop3 := range localData3.([]interface{}) {
							dataLoop3Tmp := make(map[string]interface{})
							if dataLoop3 != nil {
								dataLoop3Tmp = dataLoop3.(map[string]interface{})
							}
							dataLoop3Map := make(map[string]interface{})
							dataLoop3Map["MirrorMultiAlternateNumber"] = dataLoop3Tmp["mirror_multi_alternate_number"]
							dataLoop3Map["MirrorMultiAlternateURL"] = dataLoop3Tmp["mirror_multi_alternate_url"]
							dataLoop3Map["MirrorMultiAlternateVpcId"] = dataLoop3Tmp["mirror_multi_alternate_vpc_id"]
							dataLoop3Map["MirrorMultiAlternateDstRegion"] = dataLoop3Tmp["mirror_multi_alternate_dst_region"]
							localMaps3 = append(localMaps3, dataLoop3Map)
						}
						mirrorMultiAlternates["MirrorMultiAlternate"] = localMaps3
					}

					localData2["MirrorMultiAlternates"] = mirrorMultiAlternates
					mirrorAuth := make(map[string]interface{})
					accessKeySecret1, _ := jsonpath.Get("$[0].mirror_auth[0].access_key_secret", dataLoopTmp["redirect"])
					if accessKeySecret1 != nil && accessKeySecret1 != "" {
						mirrorAuth["AccessKeySecret"] = accessKeySecret1
					}
					accessKeyId1, _ := jsonpath.Get("$[0].mirror_auth[0].access_key_id", dataLoopTmp["redirect"])
					if accessKeyId1 != nil && accessKeyId1 != "" {
						mirrorAuth["AccessKeyId"] = accessKeyId1
					}
					authType1, _ := jsonpath.Get("$[0].mirror_auth[0].auth_type", dataLoopTmp["redirect"])
					if authType1 != nil && authType1 != "" {
						mirrorAuth["AuthType"] = authType1
					}
					region1, _ := jsonpath.Get("$[0].mirror_auth[0].region", dataLoopTmp["redirect"])
					if region1 != nil && region1 != "" {
						mirrorAuth["Region"] = region1
					}

					localData2["MirrorAuth"] = mirrorAuth
					mirrorReturnHeaders := make(map[string]interface{})
					if v, ok := dataLoopTmp["redirect"]; ok {
						localData4, err := jsonpath.Get("$[0].mirror_return_headers[0].return_header", v)
						if err != nil {
							localData4 = make([]interface{}, 0)
						}
						localMaps4 := make([]interface{}, 0)
						for _, dataLoop4 := range localData4.([]interface{}) {
							dataLoop4Tmp := make(map[string]interface{})
							if dataLoop4 != nil {
								dataLoop4Tmp = dataLoop4.(map[string]interface{})
							}
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["Value"] = dataLoop4Tmp["value"]
							dataLoop4Map["Key"] = dataLoop4Tmp["key"]
							localMaps4 = append(localMaps4, dataLoop4Map)
						}
						mirrorReturnHeaders["ReturnHeader"] = localMaps4
					}

					localData2["MirrorReturnHeaders"] = mirrorReturnHeaders
					mirrorTaggings := make(map[string]interface{})
					if v, ok := dataLoopTmp["redirect"]; ok {
						localData5, err := jsonpath.Get("$[0].mirror_taggings[0].taggings", v)
						if err != nil {
							localData5 = make([]interface{}, 0)
						}
						localMaps5 := make([]interface{}, 0)
						for _, dataLoop5 := range localData5.([]interface{}) {
							dataLoop5Tmp := make(map[string]interface{})
							if dataLoop5 != nil {
								dataLoop5Tmp = dataLoop5.(map[string]interface{})
							}
							dataLoop5Map := make(map[string]interface{})
							dataLoop5Map["Key"] = dataLoop5Tmp["key"]
							dataLoop5Map["Value"] = dataLoop5Tmp["value"]
							localMaps5 = append(localMaps5, dataLoop5Map)
						}
						mirrorTaggings["Taggings"] = localMaps5
					}

					localData2["MirrorTaggings"] = mirrorTaggings
					mirrorAsyncStatus1, _ := jsonpath.Get("$[0].mirror_async_status", dataLoopTmp["redirect"])
					if mirrorAsyncStatus1 != nil && mirrorAsyncStatus1 != "" {
						localData2["MirrorAsyncStatus"] = mirrorAsyncStatus1
					}
					transparentMirrorResponseCodes1, _ := jsonpath.Get("$[0].transparent_mirror_response_codes", dataLoopTmp["redirect"])
					if transparentMirrorResponseCodes1 != nil && transparentMirrorResponseCodes1 != "" {
						localData2["TransparentMirrorResponseCodes"] = transparentMirrorResponseCodes1
					}
					mirrorAllowHeadObject1, _ := jsonpath.Get("$[0].mirror_allow_head_object", dataLoopTmp["redirect"])
					if mirrorAllowHeadObject1 != nil && mirrorAllowHeadObject1 != "" {
						localData2["MirrorAllowHeadObject"] = mirrorAllowHeadObject1
					}
					mirrorRole1, _ := jsonpath.Get("$[0].mirror_role", dataLoopTmp["redirect"])
					if mirrorRole1 != nil && mirrorRole1 != "" {
						localData2["MirrorRole"] = mirrorRole1
					}
					mirrorUsingRole1, _ := jsonpath.Get("$[0].mirror_using_role", dataLoopTmp["redirect"])
					if mirrorUsingRole1 != nil && mirrorUsingRole1 != "" {
						localData2["MirrorUsingRole"] = mirrorUsingRole1
					}
					mirrorTunnelId1, _ := jsonpath.Get("$[0].mirror_tunnel_id", dataLoopTmp["redirect"])
					if mirrorTunnelId1 != nil && mirrorTunnelId1 != "" {
						localData2["MirrorTunnelId"] = mirrorTunnelId1
					}
					mirrorSwitchAllErrors1, _ := jsonpath.Get("$[0].mirror_switch_all_errors", dataLoopTmp["redirect"])
					if mirrorSwitchAllErrors1 != nil && mirrorSwitchAllErrors1 != "" {
						localData2["MirrorSwitchAllErrors"] = mirrorSwitchAllErrors1
					}
					mirrorUserLastModified1, _ := jsonpath.Get("$[0].mirror_user_last_modified", dataLoopTmp["redirect"])
					if mirrorUserLastModified1 != nil && mirrorUserLastModified1 != "" {
						localData2["MirrorUserLastModified"] = mirrorUserLastModified1
					}
					mirrorDstSlaveVpcId1, _ := jsonpath.Get("$[0].mirror_dst_slave_vpc_id", dataLoopTmp["redirect"])
					if mirrorDstSlaveVpcId1 != nil && mirrorDstSlaveVpcId1 != "" {
						localData2["MirrorDstSlaveVpcId"] = mirrorDstSlaveVpcId1
					}
					mirrorDstVpcId1, _ := jsonpath.Get("$[0].mirror_dst_vpc_id", dataLoopTmp["redirect"])
					if mirrorDstVpcId1 != nil && mirrorDstVpcId1 != "" {
						localData2["MirrorDstVpcId"] = mirrorDstVpcId1
					}
					mirrorDstRegion1, _ := jsonpath.Get("$[0].mirror_dst_region", dataLoopTmp["redirect"])
					if mirrorDstRegion1 != nil && mirrorDstRegion1 != "" {
						localData2["MirrorDstRegion"] = mirrorDstRegion1
					}
					mirrorIsExpressTunnel1, _ := jsonpath.Get("$[0].mirror_is_express_tunnel", dataLoopTmp["redirect"])
					if mirrorIsExpressTunnel1 != nil && mirrorIsExpressTunnel1 != "" {
						localData2["MirrorIsExpressTunnel"] = mirrorIsExpressTunnel1
					}
					mirrorAllowVideoSnapshot1, _ := jsonpath.Get("$[0].mirror_allow_video_snapshot", dataLoopTmp["redirect"])
					if mirrorAllowVideoSnapshot1 != nil && mirrorAllowVideoSnapshot1 != "" {
						localData2["MirrorAllowVideoSnapshot"] = mirrorAllowVideoSnapshot1
					}
					mirrorAllowGetImageInfo1, _ := jsonpath.Get("$[0].mirror_allow_get_image_info", dataLoopTmp["redirect"])
					if mirrorAllowGetImageInfo1 != nil && mirrorAllowGetImageInfo1 != "" {
						localData2["MirrorAllowGetImageInfo"] = mirrorAllowGetImageInfo1
					}
					mirrorProxyPass1, _ := jsonpath.Get("$[0].mirror_proxy_pass", dataLoopTmp["redirect"])
					if mirrorProxyPass1 != nil && mirrorProxyPass1 != "" {
						localData2["MirrorProxyPass"] = mirrorProxyPass1
					}
					mirrorSaveOssMeta1, _ := jsonpath.Get("$[0].mirror_save_oss_meta", dataLoopTmp["redirect"])
					if mirrorSaveOssMeta1 != nil && mirrorSaveOssMeta1 != "" {
						localData2["MirrorSaveOssMeta"] = mirrorSaveOssMeta1
					}
					mirrorUrlProbe, _ := jsonpath.Get("$[0].mirror_url_probe", dataLoopTmp["redirect"])
					if mirrorUrlProbe != nil && mirrorUrlProbe != "" {
						localData2["MirrorURLProbe"] = mirrorUrlProbe
					}
					mirrorUrlSlave, _ := jsonpath.Get("$[0].mirror_url_slave", dataLoopTmp["redirect"])
					if mirrorUrlSlave != nil && mirrorUrlSlave != "" {
						localData2["MirrorURLSlave"] = mirrorUrlSlave
					}
					mirrorUrl, _ := jsonpath.Get("$[0].mirror_url", dataLoopTmp["redirect"])
					if mirrorUrl != nil && mirrorUrl != "" {
						localData2["MirrorURL"] = mirrorUrl
					}
					mirrorSni, _ := jsonpath.Get("$[0].mirror_sni", dataLoopTmp["redirect"])
					if mirrorSni != nil && mirrorSni != "" {
						localData2["MirrorSNI"] = mirrorSni
					}
					mirrorPassOriginalSlashes1, _ := jsonpath.Get("$[0].mirror_pass_original_slashes", dataLoopTmp["redirect"])
					if mirrorPassOriginalSlashes1 != nil && mirrorPassOriginalSlashes1 != "" {
						localData2["MirrorPassOriginalSlashes"] = mirrorPassOriginalSlashes1
					}
					httpRedirectCode1, _ := jsonpath.Get("$[0].http_redirect_code", dataLoopTmp["redirect"])
					if httpRedirectCode1 != nil && httpRedirectCode1 != "" {
						localData2["HttpRedirectCode"] = httpRedirectCode1
					}
					replaceKeyWith1, _ := jsonpath.Get("$[0].replace_key_with", dataLoopTmp["redirect"])
					if replaceKeyWith1 != nil && replaceKeyWith1 != "" {
						localData2["ReplaceKeyWith"] = replaceKeyWith1
					}
					enableReplacePrefix1, _ := jsonpath.Get("$[0].enable_replace_prefix", dataLoopTmp["redirect"])
					if enableReplacePrefix1 != nil && enableReplacePrefix1 != "" {
						localData2["EnableReplacePrefix"] = enableReplacePrefix1
					}
					replaceKeyPrefixWith1, _ := jsonpath.Get("$[0].replace_key_prefix_with", dataLoopTmp["redirect"])
					if replaceKeyPrefixWith1 != nil && replaceKeyPrefixWith1 != "" {
						localData2["ReplaceKeyPrefixWith"] = replaceKeyPrefixWith1
					}
					hostName1, _ := jsonpath.Get("$[0].host_name", dataLoopTmp["redirect"])
					if hostName1 != nil && hostName1 != "" {
						localData2["HostName"] = hostName1
					}
					protocol1, _ := jsonpath.Get("$[0].protocol", dataLoopTmp["redirect"])
					if protocol1 != nil && protocol1 != "" {
						localData2["Protocol"] = protocol1
					}
					mirrorHeaders := make(map[string]interface{})
					if v, ok := dataLoopTmp["redirect"]; ok {
						localData6, err := jsonpath.Get("$[0].mirror_headers[0].set", v)
						if err != nil {
							localData6 = make([]interface{}, 0)
						}
						localMaps6 := make([]interface{}, 0)
						for _, dataLoop6 := range localData6.([]interface{}) {
							dataLoop6Tmp := make(map[string]interface{})
							if dataLoop6 != nil {
								dataLoop6Tmp = dataLoop6.(map[string]interface{})
							}
							dataLoop6Map := make(map[string]interface{})
							dataLoop6Map["Key"] = dataLoop6Tmp["key"]
							dataLoop6Map["Value"] = dataLoop6Tmp["value"]
							localMaps6 = append(localMaps6, dataLoop6Map)
						}
						mirrorHeaders["Set"] = localMaps6
					}

					remove1, _ := jsonpath.Get("$[0].mirror_headers[0].remove", dataLoopTmp["redirect"])
					if remove1 != nil && remove1 != "" {
						mirrorHeaders["Remove"] = remove1
					}
					pass1, _ := jsonpath.Get("$[0].mirror_headers[0].pass", dataLoopTmp["redirect"])
					if pass1 != nil && pass1 != "" {
						mirrorHeaders["Pass"] = pass1
					}
					passAll1, _ := jsonpath.Get("$[0].mirror_headers[0].pass_all", dataLoopTmp["redirect"])
					if passAll1 != nil && passAll1 != "" {
						mirrorHeaders["PassAll"] = passAll1
					}

					localData2["MirrorHeaders"] = mirrorHeaders
					mirrorCheckMd51, _ := jsonpath.Get("$[0].mirror_check_md5", dataLoopTmp["redirect"])
					if mirrorCheckMd51 != nil && mirrorCheckMd51 != "" {
						localData2["MirrorCheckMd5"] = mirrorCheckMd51
					}
					mirrorFollowRedirect1, _ := jsonpath.Get("$[0].mirror_follow_redirect", dataLoopTmp["redirect"])
					if mirrorFollowRedirect1 != nil && mirrorFollowRedirect1 != "" {
						localData2["MirrorFollowRedirect"] = mirrorFollowRedirect1
					}
					mirrorPassQueryString1, _ := jsonpath.Get("$[0].mirror_pass_query_string", dataLoopTmp["redirect"])
					if mirrorPassQueryString1 != nil && mirrorPassQueryString1 != "" {
						localData2["MirrorPassQueryString"] = mirrorPassQueryString1
					}
					passQueryString1, _ := jsonpath.Get("$[0].pass_query_string", dataLoopTmp["redirect"])
					if passQueryString1 != nil && passQueryString1 != "" {
						localData2["PassQueryString"] = passQueryString1
					}
					redirectType1, _ := jsonpath.Get("$[0].redirect_type", dataLoopTmp["redirect"])
					if redirectType1 != nil && redirectType1 != "" {
						localData2["RedirectType"] = redirectType1
					}
					dataLoopMap["Redirect"] = localData2
				}
				if !IsNil(dataLoopTmp["condition"]) {
					localData7 := make(map[string]interface{})
					if v, ok := dataLoopTmp["condition"]; ok {
						localData8, err := jsonpath.Get("$[0].include_headers", v)
						if err != nil {
							localData8 = make([]interface{}, 0)
						}
						localMaps8 := make([]interface{}, 0)
						for _, dataLoop8 := range localData8.([]interface{}) {
							dataLoop8Tmp := make(map[string]interface{})
							if dataLoop8 != nil {
								dataLoop8Tmp = dataLoop8.(map[string]interface{})
							}
							dataLoop8Map := make(map[string]interface{})
							dataLoop8Map["EndsWith"] = dataLoop8Tmp["ends_with"]
							dataLoop8Map["StartsWith"] = dataLoop8Tmp["starts_with"]
							dataLoop8Map["Equals"] = dataLoop8Tmp["equals"]
							dataLoop8Map["Key"] = dataLoop8Tmp["key"]
							localMaps8 = append(localMaps8, dataLoop8Map)
						}
						localData7["IncludeHeader"] = localMaps8
					}

					httpErrorCodeReturnedEquals1, _ := jsonpath.Get("$[0].http_error_code_returned_equals", dataLoopTmp["condition"])
					if httpErrorCodeReturnedEquals1 != nil && httpErrorCodeReturnedEquals1 != "" {
						localData7["HttpErrorCodeReturnedEquals"] = httpErrorCodeReturnedEquals1
					}
					keySuffixEquals1, _ := jsonpath.Get("$[0].key_suffix_equals", dataLoopTmp["condition"])
					if keySuffixEquals1 != nil && keySuffixEquals1 != "" {
						localData7["KeySuffixEquals"] = keySuffixEquals1
					}
					keyPrefixEquals1, _ := jsonpath.Get("$[0].key_prefix_equals", dataLoopTmp["condition"])
					if keyPrefixEquals1 != nil && keyPrefixEquals1 != "" {
						localData7["KeyPrefixEquals"] = keyPrefixEquals1
					}
					dataLoopMap["Condition"] = localData7
				}
				dataLoopMap["RuleNumber"] = dataLoopTmp["rule_number"]
				localMaps = append(localMaps, dataLoopMap)
			}
			routingRules["RoutingRule"] = localMaps
		}

		objectDataLocalMap["RoutingRules"] = routingRules
	}

	if d.HasChange("error_document") {
		update = true
	}
	if v := d.Get("error_document"); !IsNil(v) {
		errorDocument := make(map[string]interface{})
		httpStatus1, _ := jsonpath.Get("$[0].http_status", v)
		if httpStatus1 != nil && (d.HasChange("error_document.0.http_status") || httpStatus1 != "") {
			errorDocument["HttpStatus"] = httpStatus1
		}
		key9, _ := jsonpath.Get("$[0].key", v)
		if key9 != nil && (d.HasChange("error_document.0.key") || key9 != "") {
			errorDocument["Key"] = key9
		}

		objectDataLocalMap["ErrorDocument"] = errorDocument
	}

	if d.HasChange("index_document") {
		update = true
	}
	if v := d.Get("index_document"); !IsNil(v) {
		indexDocument := make(map[string]interface{})
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && (d.HasChange("index_document.0.type") || type1 != "") {
			indexDocument["Type"] = type1
		}
		suffix1, _ := jsonpath.Get("$[0].suffix", v)
		if suffix1 != nil && (d.HasChange("index_document.0.suffix") || suffix1 != "") {
			indexDocument["Suffix"] = suffix1
		}
		supportSubDir1, _ := jsonpath.Get("$[0].support_sub_dir", v)
		if supportSubDir1 != nil && (d.HasChange("index_document.0.support_sub_dir") || supportSubDir1 != "") {
			indexDocument["SupportSubDir"] = supportSubDir1
		}

		objectDataLocalMap["IndexDocument"] = indexDocument
	}

	request["WebsiteConfiguration"] = objectDataLocalMap
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketWebsite", action), query, body, nil, hostMap, false)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudOssBucketWebsiteRead(d, meta)
}

func resourceAliCloudOssBucketWebsiteDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?website")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketWebsite", action), query, nil, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchWebsiteConfiguration"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

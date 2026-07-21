// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudApigDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigDomainRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_cert_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ca_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"force_https": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"http2_option": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"m_tls_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls_cipher_suites_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tls_cipher_suite": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"support_versions": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"config_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tls_max": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls_min": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudApigDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListDomains
	action := fmt.Sprintf("/v1/domains")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["pageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["domainId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["domainId"]

		mapping["cert_identifier"] = objectRaw["certIdentifier"]
		mapping["client_ca_cert"] = objectRaw["clientCACert"]
		mapping["domain_name"] = objectRaw["name"]
		mapping["domain_scope"] = objectRaw["domainScope"]
		mapping["force_https"] = objectRaw["forceHttps"]
		mapping["m_tls_enabled"] = objectRaw["mTLSEnabled"]
		mapping["protocol"] = objectRaw["protocol"]
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["domain_id"] = objectRaw["domainId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["domainId"])
		mapping, err = dataSourceAliCloudApigDomainReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("domains", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudApigDomainReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	apigServiceV2 := ApigServiceV2{client}
	getResp, err := apigServiceV2.DescribeApigDomain(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["ca_cert_identifier"] = objectRaw["caCertIdentifier"]
	mapping["cert_identifier"] = objectRaw["certIdentifier"]
	mapping["client_ca_cert"] = objectRaw["clientCACert"]
	mapping["domain_name"] = objectRaw["name"]
	mapping["domain_scope"] = objectRaw["domainScope"]
	mapping["force_https"] = objectRaw["forceHttps"]
	mapping["http2_option"] = objectRaw["http2Option"]
	mapping["m_tls_enabled"] = objectRaw["mTLSEnabled"]
	mapping["protocol"] = objectRaw["protocol"]
	mapping["resource_group_id"] = objectRaw["resourceGroupId"]
	mapping["tls_max"] = objectRaw["tlsMax"]
	mapping["tls_min"] = objectRaw["tlsMin"]
	mapping["domain_id"] = objectRaw["domainId"]

	tlsCipherSuitesConfigMaps := make([]map[string]interface{}, 0)
	tlsCipherSuitesConfigMap := make(map[string]interface{})
	tlsCipherSuitesConfigRaw := make(map[string]interface{})
	if objectRaw["tlsCipherSuitesConfig"] != nil {
		tlsCipherSuitesConfigRaw = objectRaw["tlsCipherSuitesConfig"].(map[string]interface{})
	}
	if len(tlsCipherSuitesConfigRaw) > 0 {
		tlsCipherSuitesConfigMap["config_type"] = tlsCipherSuitesConfigRaw["configType"]

		tlsCipherSuiteRaw := tlsCipherSuitesConfigRaw["tlsCipherSuite"]
		tlsCipherSuiteMaps := make([]map[string]interface{}, 0)
		if tlsCipherSuiteRaw != nil {
			for _, tlsCipherSuiteChildRaw := range convertToInterfaceArray(tlsCipherSuiteRaw) {
				tlsCipherSuiteMap := make(map[string]interface{})
				tlsCipherSuiteChildRaw := tlsCipherSuiteChildRaw.(map[string]interface{})
				tlsCipherSuiteMap["name"] = tlsCipherSuiteChildRaw["name"]

				supportVersionsRaw := make([]interface{}, 0)
				if tlsCipherSuiteChildRaw["supportVersions"] != nil {
					supportVersionsRaw = convertToInterfaceArray(tlsCipherSuiteChildRaw["supportVersions"])
				}

				tlsCipherSuiteMap["support_versions"] = supportVersionsRaw
				tlsCipherSuiteMaps = append(tlsCipherSuiteMaps, tlsCipherSuiteMap)
			}
		}
		tlsCipherSuitesConfigMap["tls_cipher_suite"] = tlsCipherSuiteMaps
		tlsCipherSuitesConfigMaps = append(tlsCipherSuitesConfigMaps, tlsCipherSuitesConfigMap)
	}
	mapping["tls_cipher_suites_config"] = tlsCipherSuitesConfigMaps

	return mapping, nil
}

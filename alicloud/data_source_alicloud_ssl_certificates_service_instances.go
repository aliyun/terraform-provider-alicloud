// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudSslCertificatesServiceInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSslCertificatesServiceInstanceRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"brand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_reissue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"average_waiting_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"brand": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"csr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_domain_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"generate_csr_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order_end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"order_start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upgrade_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wildcard_domain_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAliCloudSslCertificatesServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

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
	var query map[string]interface{}
	action := "ListInstances"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("brand"); ok {
		request["Brand"] = v
	}
	if v, ok := d.GetOk("certificate_status"); ok {
		request["CertificateStatus"] = v
	}
	if v, ok := d.GetOk("certificate_type"); ok {
		request["CertificateType"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["ShowSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.InstanceList[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["InstanceId"]

		mapping["instance_id"] = objectRaw["InstanceId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["InstanceId"])
		mapping, err = dataSourceAliCloudSslCertificatesServiceInstanceReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudSslCertificatesServiceInstanceReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	sslCertificatesServiceServiceV2 := SslCertificatesServiceServiceV2{client}
	getResp, err := sslCertificatesServiceServiceV2.DescribeSslCertificatesServiceInstance(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["auto_reissue"] = objectRaw["AutoReissue"]
	mapping["average_waiting_time"] = objectRaw["AverageWaitingTime"]
	mapping["brand"] = objectRaw["Brand"]
	mapping["certificate_type"] = objectRaw["CertificateType"]
	mapping["city"] = objectRaw["City"]
	mapping["company_id"] = objectRaw["CompanyId"]
	mapping["country_code"] = objectRaw["CountryCode"]
	mapping["csr"] = objectRaw["Csr"]
	mapping["domain"] = objectRaw["Domain"]
	mapping["full_domain_count"] = objectRaw["FullDomainCount"]
	mapping["generate_csr_method"] = objectRaw["GenerateCsrMethod"]
	mapping["instance_end_time"] = objectRaw["InstanceEndTime"]
	mapping["instance_name"] = objectRaw["CertificateName"]
	mapping["instance_start_time"] = objectRaw["InstanceStartTime"]
	mapping["instance_type"] = objectRaw["InstanceType"]
	mapping["key_algorithm"] = objectRaw["KeyAlgorithm"]
	mapping["order_end_time"] = objectRaw["OrderEndTime"]
	mapping["order_start_time"] = objectRaw["OrderStartTime"]
	mapping["province"] = objectRaw["Province"]
	mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
	mapping["spec"] = objectRaw["Spec"]
	mapping["status"] = objectRaw["Status"]
	mapping["upgrade_status"] = objectRaw["UpgradeStatus"]
	mapping["validation_method"] = objectRaw["ValidationMethod"]
	mapping["wildcard_domain_count"] = objectRaw["WildcardDomainCount"]
	mapping["instance_id"] = objectRaw["InstanceId"]

	contactIdListRaw := make([]interface{}, 0)
	if objectRaw["ContactIdList"] != nil {
		contactIdListRaw = convertToInterfaceArray(objectRaw["ContactIdList"])
	}

	mapping["contact_id_list"] = contactIdListRaw
	tagsMaps := objectRaw["Tags"]
	mapping["tags"] = tagsToMap(tagsMaps)

	return mapping, nil
}

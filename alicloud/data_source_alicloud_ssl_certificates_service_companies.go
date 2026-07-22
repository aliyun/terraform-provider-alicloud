// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudSslCertificatesServiceCompanies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSslCertificatesServiceCompanyRead,
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
			"company_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"companies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"company_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"company_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"department": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"post_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudSslCertificatesServiceCompanyRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]interface{}
	action := "ListCompanies"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v := d.Get("company_id").(int); v != 0 {
		request["CompanyId"] = v
	}
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ShowSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.CompanyList[*]", response)

		result, ok := resp.([]interface{})
		if !ok {
			return WrapError(fmt.Errorf("CompanyList is not a list, got %T", resp))
		}
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["CompanyName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CompanyId"])]; !ok {
					continue
				}
			}
			// API ListCompanies does not honor CompanyId server-side (returns full
			// list even when CompanyId is set on the request), so filter client-side
			// to make the company_id argument functional.
			if cid := d.Get("company_id").(int); cid != 0 {
				if fmt.Sprint(item["CompanyId"]) != fmt.Sprint(cid) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["CompanyId"]

		mapping["city"] = objectRaw["City"]
		mapping["company_address"] = objectRaw["CompanyAddress"]
		mapping["company_code"] = objectRaw["CompanyCode"]
		mapping["company_email"] = objectRaw["CompanyEmail"]
		mapping["company_name"] = objectRaw["CompanyName"]
		mapping["company_phone"] = objectRaw["CompanyPhone"]
		mapping["company_type"] = objectRaw["CompanyType"]
		mapping["country_code"] = objectRaw["CountryCode"]
		mapping["department"] = objectRaw["Department"]
		mapping["lang"] = objectRaw["Lang"]
		mapping["post_code"] = objectRaw["PostCode"]
		mapping["province"] = objectRaw["Province"]
		mapping["company_id"] = objectRaw["CompanyId"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["CompanyName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("companies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

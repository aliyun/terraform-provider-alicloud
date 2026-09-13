package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudEsaCustomHostnames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaCustomHostnamesRead,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_match_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"record_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostnames": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cas_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cas_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_apply_code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cert_apply_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_http_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_http_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_not_after": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_txt_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_txt_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conflict_with": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"offline_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"record_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"site_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_flag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"verify_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"verify_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAliCloudEsaCustomHostnamesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListCustomHostnames"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["SiteId"] = d.Get("site_id")

	if v, ok := d.GetOk("record_id"); ok {
		request["RecordId"] = v
	}
	if v, ok := d.GetOk("hostname"); ok {
		request["Hostname"] = v
	}
	if v, ok := d.GetOk("name_match_type"); ok {
		request["MatchType"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
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

	var hostnameRegex *regexp.Regexp
	if v, ok := d.GetOk("hostname_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		hostnameRegex = r
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("ESA", "2024-09-10", action, request, nil)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_custom_hostnames", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Data.Content.Data", response)
		if err != nil {
			if NotFoundError(err) {
				break
			}
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Content.Data", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			hostnameId := fmt.Sprint(item["HostnameId"])
			if len(idsMap) > 0 {
				if _, ok := idsMap[hostnameId]; !ok {
					continue
				}
			}
			if hostnameRegex != nil {
				if !hostnameRegex.MatchString(fmt.Sprint(item["Hostname"])) {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}

		pageNumber, ok := request["PageNumber"].(int)
		if !ok {
			break
		}
		request["PageNumber"] = pageNumber + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cas_id":             object["CasId"],
			"cas_region":         object["CasRegion"],
			"cert_apply_code":    object["CertApplyCode"],
			"cert_apply_message": object["CertApplyMessage"],
			"cert_http_key":      object["CertHttpKey"],
			"cert_http_value":    object["CertHttpValue"],
			"cert_id":            object["CertId"],
			"cert_not_after":     object["CertNotAfter"],
			"cert_status":        object["CertStatus"],
			"cert_txt_key":       object["CertTxtKey"],
			"cert_txt_value":     object["CertTxtValue"],
			"cert_type":          object["CertType"],
			"conflict_with":      object["ConflictWith"],
			"create_time":        object["GmtCreate"],
			"hostname":           object["Hostname"],
			"hostname_id":        fmt.Sprint(object["HostnameId"]),
			"offline_reason":     object["OfflineReason"],
			"record_id":          object["RecordId"],
			"record_name":        object["RecordName"],
			"site_id":            object["SiteId"],
			"site_name":          object["SiteName"],
			"ssl_flag":           object["SslFlag"],
			"status":             object["Status"],
			"update_time":        object["GmtModified"],
			"verify_code":        object["VerifyCode"],
			"verify_host":        object["VerifyHost"],
		}
		ids = append(ids, fmt.Sprint(mapping["hostname_id"]))
		names = append(names, object["Hostname"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("hostnames", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

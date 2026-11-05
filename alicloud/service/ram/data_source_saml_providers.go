package ram

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceSamlProviders() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceSamlProvidersRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encodedsaml_metadata_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"saml_provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func DataSourceSamlProvidersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSAMLProviders"
	request := make(map[string]interface{})
	request["MaxItems"] = helper.PageSizeLarge
	var objects []map[string]interface{}
	var sAMLProviderNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		sAMLProviderNameRegex = r
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
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})

		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_saml_providers", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.SAMLProviders.SAMLProvider", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.SAMLProviders.SAMLProvider", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if sAMLProviderNameRegex != nil {
				if !sAMLProviderNameRegex.MatchString(fmt.Sprint(item["SAMLProviderName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SAMLProviderName"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if marker, ok := response["Marker"].(string); ok && marker != "" {
			request["Marker"] = marker
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"arn":                object["Arn"],
			"description":        object["Description"],
			"id":                 fmt.Sprint(object["SAMLProviderName"]),
			"saml_provider_name": fmt.Sprint(object["SAMLProviderName"]),
			"update_date":        object["UpdateDate"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["SAMLProviderName"]))
			names = append(names, object["SAMLProviderName"])
			s = append(s, mapping)
			continue
		}

		imsService := RamServiceV2{client}
		id := fmt.Sprint(object["SAMLProviderName"])
		getResp, err := imsService.DescribeRamSamlProvider(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["encodedsaml_metadata_document"] = getResp["EncodedSAMLMetadataDocument"]
		ids = append(ids, fmt.Sprint(object["SAMLProviderName"]))
		names = append(names, object["SAMLProviderName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("providers", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

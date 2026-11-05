package alidns

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudAlidnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsRecordsRead,
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rr_key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"type_key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value_key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rr_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"value_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlidnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateDescribeDomainRecordsRequest()
	if v, ok := d.GetOk("direction"); ok {
		request.Direction = v.(string)
	}
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("key_word"); ok {
		request.KeyWord = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("line"); ok {
		request.Line = v.(string)
	}
	if v, ok := d.GetOk("order_by"); ok {
		request.OrderBy = v.(string)
	}
	if v, ok := d.GetOk("rr_key_word"); ok {
		request.RRKeyWord = v.(string)
	}
	if v, ok := d.GetOk("search_mode"); ok {
		request.SearchMode = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		request.Type = v.(string)
	}
	if v, ok := d.GetOk("type_key_word"); ok {
		request.TypeKeyWord = v.(string)
	}
	if v, ok := d.GetOk("value_key_word"); ok {
		request.ValueKeyWord = v.(string)
	}
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []alidnssdk.Record
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var rrRegex *regexp.Regexp
	if v, ok := d.GetOk("rr_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		rrRegex = r
	}
	var valueRegex *regexp.Regexp
	if v, ok := d.GetOk("value_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		valueRegex = r
	}
	for {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainRecords(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_alidns_records", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidnssdk.DescribeDomainRecordsResponse)

		for _, item := range response.DomainRecords.Record {
			if rrRegex != nil {
				if !rrRegex.MatchString(item.RR) {
					continue
				}
			}
			if valueRegex != nil {
				if !valueRegex.MatchString(item.Value) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.RecordId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.DomainRecords.Record) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"domain_name": object.DomainName,
			"line":        object.Line,
			"locked":      object.Locked,
			"priority":    object.Priority,
			"rr":          object.RR,
			"id":          object.RecordId,
			"record_id":   object.RecordId,
			"status":      object.Status,
			"ttl":         object.TTL,
			"type":        object.Type,
			"value":       object.Value,
			"remark":      object.Remark,
		}
		ids[i] = object.RecordId
		s[i] = mapping
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("records", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

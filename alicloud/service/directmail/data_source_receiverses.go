package directmail

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceReceiverses() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceReceiversesRead,
		Schema: map[string]*schema.Schema{
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
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
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"receiverses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"receivers_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"receivers_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"receivers_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceReceiversesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QueryReceiverByParam"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNo"] = 1
	if v, ok := d.GetOk("key_word"); ok {
		request["KeyWord"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	var objects []map[string]interface{}
	var receiversNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		receiversNameRegex = r
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
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_direct_mail_receiverses", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.data.receiver", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.data.receiver", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if receiversNameRegex != nil {
				if !receiversNameRegex.MatchString(fmt.Sprint(item["ReceiversName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ReceiverId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		statusVal, err := convert.FormatInt(object["ReceiversStatus"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"create_time":     object["CreateTime"],
			"description":     object["Desc"],
			"receivers_alias": object["ReceiversAlias"],
			"id":              fmt.Sprint(object["ReceiverId"]),
			"receivers_id":    fmt.Sprint(object["ReceiverId"]),
			"receivers_name":  object["ReceiversName"],
			"status":          statusVal,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ReceiversName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("receiverses", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

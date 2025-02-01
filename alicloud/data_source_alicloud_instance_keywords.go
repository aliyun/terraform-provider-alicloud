package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudInstanceKeywords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstanceKeywordsRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"account", "database"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// response value
			"keywords": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudInstanceKeywordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeInstanceKeywords"
	request := map[string]interface{}{
		"SourceIp": client.SourceIp,
		"Key":      d.Get("key"),
	}
	var ids []string
	var keywords []string
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instance_keywords", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Words.word", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Words.word", response)
	}
	for _, r := range resp.([]interface{}) {
		ids = append(ids, fmt.Sprint(r))
		keywords = append(keywords, fmt.Sprint(r))
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("keywords", keywords); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), keywords)
	}

	return nil
}

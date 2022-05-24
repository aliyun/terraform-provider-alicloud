package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudTagMetaTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudTagMetaTagsRead,
		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudTagMetaTagsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTagKeys"
	request := make(map[string]interface{})
	request["QueryType"] = "MetaTag"
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}

	keyName, keyNameOk := d.GetOk("key_name")
	var response map[string]interface{}
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_tag_meta_tags", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Keys.Key", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Keys.Key", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if keyNameOk && keyName.(string) != "" && keyName.(string) != item["Key"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		tagService := TagService{client}
		getResp, err := tagService.ListTagValues(fmt.Sprint(object["Key"]))
		if err != nil {
			return WrapError(err)
		}
		for _, item := range getResp {
			mapping := map[string]interface{}{
				"key_name":   fmt.Sprint(object["Key"]),
				"value_name": fmt.Sprint(item),
				"category":   fmt.Sprint(object["Category"]),
			}

			s = append(s, mapping)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("tags", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

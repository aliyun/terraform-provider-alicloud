package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCdnBlockedRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCdnBlockedRegionsRead,
		Schema: map[string]*schema.Schema{
			"language": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en", "jp"}, false),
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"countries_and_regions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"continent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"countries_and_regions_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCdnBlockedRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBlockedRegions"
	request := make(map[string]interface{})

	request["Language"] = d.Get("language")

	var response map[string]interface{}
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), request, nil, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_blocked_regions", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.InfoList.InfoItem", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InfoList.InfoItem", response)
	}
	objects := resp.([]interface{})
	s := make([]map[string]interface{}, 0)
	for _, item := range objects {
		object := item.(map[string]interface{})
		mapping := map[string]interface{}{
			"continent":                  fmt.Sprint(object["Continent"]),
			"countries_and_regions_name": fmt.Sprint(object["CountriesAndRegionsName"]),
			"countries_and_regions":      fmt.Sprint(object["CountriesAndRegions"]),
		}

		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("regions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

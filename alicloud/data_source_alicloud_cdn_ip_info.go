package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCdnIpInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCdnIpInfoRead,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cdn_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isp_ename": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_ename": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCdnIpInfoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeIpInfo"
	request := make(map[string]interface{})

	id := ""
	if v, ok := d.GetOk("ip"); ok {
		request["IP"] = v.(string)
		id = v.(string)
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("Cdn", "2018-05-10", action, request, nil)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_ip_info", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "", response)
	}
	object := resp.(map[string]interface{})
	d.SetId(id)
	d.Set("cdn_ip", object["CdnIp"])
	d.Set("isp", object["ISP"])
	d.Set("isp_ename", object["IspEname"])
	d.Set("region", object["Region"])
	d.Set("region_ename", object["RegionEname"])
	d.Set("ip", id)

	return nil
}

package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudHbrReplicationVaultRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrReplicationVaultRegionsRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replication_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrReplicationVaultRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeVaultReplicationRegions"
	request := make(map[string]interface{})
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_replication_vault_regions", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Regions.RegionId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Regions", response)
	}
	result, _ := resp.([]interface{})
	s := make([]map[string]interface{}, 0)
	for _, v := range result {
		mapping := map[string]interface{}{
			"replication_region_id": fmt.Sprint(v),
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

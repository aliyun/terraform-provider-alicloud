package cloudstoragegateway

import (
	"fmt"
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

func DataSourceAlicloudCloudStorageGatewayGatewayCacheDisks() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAlicloudCloudStorageGatewayGatewayCacheDisksRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_disk_size_in_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cache_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_used": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renew_url": {
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

func DataSourceAlicloudCloudStorageGatewayGatewayCacheDisksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGatewayCaches"
	request := make(map[string]interface{})
	request["GatewayId"] = d.Get("gateway_id")
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_cache_disks", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.Caches.Cache", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Caches.Cache", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["GatewayId"], ":", item["CacheId"], ":", item["LocalFilePath"])]; !ok {
				continue
			}
		}
		if statusOk && fmt.Sprint(status) != fmt.Sprint(item["ExpireStatus"]) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		cache_disk_size_in_gbVal, err := convert.FormatInt(object["SizeInGB"])
		if err != nil {
			return err
		}

		expired_timeVal, err := convert.FormatInt(object["ExpiredTime"])
		if err != nil {
			return err
		}

		iopsVal, err := convert.FormatInt(object["Iops"])
		if err != nil {
			return err
		}

		statusVal, err := convert.FormatInt(object["ExpireStatus"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"cache_disk_category":   object["CacheType"],
			"cache_disk_size_in_gb": cache_disk_size_in_gbVal,
			"cache_id":              object["CacheId"],
			"expired_time":          expired_timeVal,
			"gateway_id":            request["GatewayId"],
			"iops":                  iopsVal,
			"is_used":               object["IsUsed"],
			"id":                    fmt.Sprint(request["GatewayId"], ":", object["CacheId"], ":", object["LocalFilePath"]),
			"local_file_path":       fmt.Sprint(object["LocalFilePath"]),
			"renew_url":             object["RenewURL"],
			"status":                statusVal,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("disks", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

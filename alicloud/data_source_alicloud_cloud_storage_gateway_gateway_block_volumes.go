package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumesRead,
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 8),
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chap_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"chap_in_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chunk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gateway_block_volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lun_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_bucket_ssl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_download": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_upload": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_state": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGatewayBlockVolumes"
	request := make(map[string]interface{})
	request["GatewayId"] = d.Get("gateway_id")
	var objects []map[string]interface{}
	var gatewayBlockVolumeNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		gatewayBlockVolumeNameRegex = r
	}

	status, statusOk := d.GetOk("status")
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_block_volumes", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.BlockVolumes.BlockVolume", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BlockVolumes.BlockVolume", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if gatewayBlockVolumeNameRegex != nil && !gatewayBlockVolumeNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["GatewayId"], ":", item["IndexId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(int) != item["VolumeState"] {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		chunk_sizeVal, err240 := formatInt(object["ChunkSize"])
		if err240 != nil {
			return err240
		}

		lun_idVal, err249 := formatInt(object["LunId"])
		if err249 != nil {
			return err249
		}

		portVal, err253 := formatInt(object["Port"])
		if err253 != nil {
			return err253
		}

		sizeVal, err255 := formatInt(object["Size"])
		if err255 != nil {
			return err255
		}

		statusVal, err257 := formatInt(object["Status"])
		if err257 != nil {
			return err257
		}

		total_downloadVal, err259 := formatInt(object["TotalDownload"])
		if err259 != nil {
			return err259
		}
		total_uploadVal, err260 := formatInt(object["TotalUpload"])
		if err260 != nil {
			return err260
		}
		volume_stateVal, err261 := formatInt(object["VolumeState"])
		if err261 != nil {
			return err261
		}
		mapping := map[string]interface{}{
			"address":                   object["Address"],
			"cache_mode":                object["CacheMode"],
			"chap_enabled":              object["ChapEnabled"],
			"chap_in_user":              fmt.Sprint(object["ChapInUser"]),
			"chunk_size":                chunk_sizeVal,
			"disk_id":                   object["DiskId"],
			"disk_type":                 object["DiskType"],
			"enabled":                   object["Enabled"],
			"gateway_block_volume_name": object["Name"],
			"gateway_id":                request["GatewayId"],
			"id":                        fmt.Sprint(request["GatewayId"], ":", object["IndexId"]),
			"index_id":                  fmt.Sprint(object["IndexId"]),
			"local_path":                object["LocalPath"],
			"lun_id":                    lun_idVal,
			"oss_bucket_name":           object["OssBucketName"],
			"oss_bucket_ssl":            object["OssBucketSsl"],
			"oss_endpoint":              object["OssEndpoint"],
			"port":                      portVal,
			"protocol":                  object["Protocol"],
			"size":                      sizeVal,
			"state":                     object["State"],
			"status":                    statusVal,
			"target":                    object["Target"],
			"total_download":            total_downloadVal,
			"total_upload":              total_uploadVal,
			"volume_state":              volume_stateVal,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("volumes", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

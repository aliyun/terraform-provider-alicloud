package cloudstoragegateway

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

func DataSourceAlicloudCloudStorageGatewayGatewayFileShares() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAlicloudCloudStorageGatewayGatewayFileSharesRead,
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
			"shares": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_based_enumeration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backend_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"browsable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bucket_infos": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buckets_stub": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cache_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_side_cmk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_side_encryption": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"direct_io": {
							Type:     schema.TypeBool,
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
						"download_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"express_sync_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fast_reclaim": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"fe_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_num_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fs_size_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_file_share_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ignore_delete": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"in_place": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"in_rate": {
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
						"kms_rotate_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lag_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mns_health": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nfs_v4_optimization": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"obsolete_buckets": {
							Type:     schema.TypeString,
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
						"oss_health": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"out_rate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"partial_sync_paths": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"polling_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remaining_meta_space": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_sync": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remote_sync_download": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ro_client_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ro_user_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rw_client_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rw_user_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_side_cmk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_side_encryption": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"squash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_archive": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sync_progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_download": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_upload": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transfer_acceleration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"windows_acl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"bypass_cache_read": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceAlicloudCloudStorageGatewayGatewayFileSharesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGatewayFileShares"
	request := make(map[string]interface{})
	request["GatewayId"] = d.Get("gateway_id")
	request["Refresh"] = true
	var objects []map[string]interface{}
	var gatewayFileShareNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		gatewayFileShareNameRegex = r
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_file_shares", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.FileShares.FileShare", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.FileShares.FileShare", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if gatewayFileShareNameRegex != nil && !gatewayFileShareNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["GatewayId"], ":", item["IndexId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		backend_limitVal, err := convert.FormatInt(object["BeLimit"])
		if err != nil {
			return err
		}

		download_limitVal, err := convert.FormatInt(object["DownloadLimit"])
		if err != nil {
			return err
		}

		fe_limitVal, err := convert.FormatInt(object["FeLimit"])
		if err != nil {
			return err
		}

		kms_rotate_periodVal, err := convert.FormatInt(object["KmsRotatePeriod"])
		if err != nil {
			return err
		}

		polling_intervalVal, err := convert.FormatInt(object["PollingInterval"])
		if err != nil {
			return err
		}

		sync_progressVal, err := convert.FormatInt(object["SyncProgress"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"access_based_enumeration": object["AccessBasedEnumeration"],
			"address":                  object["Address"],
			"backend_limit":            backend_limitVal,
			"browsable":                object["Browsable"],
			"bucket_infos":             object["BucketInfos"],
			"buckets_stub":             object["BucketsStub"],
			"cache_mode":               object["CacheMode"],
			"client_side_cmk":          object["ClientSideCmk"],
			"client_side_encryption":   object["ClientSideEncryption"],
			"direct_io":                object["DirectIO"],
			"disk_id":                  object["DiskId"],
			"disk_type":                object["DiskType"],
			"download_limit":           download_limitVal,
			"enabled":                  object["Enabled"],
			"fast_reclaim":             object["FastReclaim"],
			"fe_limit":                 fe_limitVal,
			"file_num_limit":           fmt.Sprint(object["FileNumLimit"]),
			"fs_size_limit":            fmt.Sprint(object["FsSizeLimit"]),
			"gateway_file_share_name":  object["Name"],
			"gateway_id":               request["GatewayId"],
			"ignore_delete":            object["IgnoreDelete"],
			"in_place":                 object["InPlace"],
			"in_rate":                  fmt.Sprint(object["InRate"]),
			"id":                       fmt.Sprint(request["GatewayId"], ":", object["IndexId"]),
			"index_id":                 fmt.Sprint(object["IndexId"]),
			"kms_rotate_period":        kms_rotate_periodVal,
			"lag_period":               fmt.Sprint(object["LagPeriod"]),
			"local_path":               object["LocalPath"],
			"mns_health":               object["MnsHealth"],
			"nfs_v4_optimization":      object["NfsV4Optimization"],
			"obsolete_buckets":         object["ObsoleteBuckets"],
			"oss_bucket_name":          object["OssBucketName"],
			"oss_bucket_ssl":           object["OssBucketSsl"],
			"oss_endpoint":             object["OssEndpoint"],
			"oss_health":               object["OssHealth"],
			"oss_used":                 fmt.Sprint(object["OssUsed"]),
			"out_rate":                 fmt.Sprint(object["OutRate"]),
			"partial_sync_paths":       object["PartialSyncPaths"],
			"path_prefix":              object["PathPrefix"],
			"polling_interval":         polling_intervalVal,
			"protocol":                 object["Protocol"],
			"remaining_meta_space":     fmt.Sprint(object["RemainingMetaSpace"]),
			"remote_sync":              object["RemoteSync"],
			"remote_sync_download":     object["RemoteSyncDownload"],
			"server_side_encryption":   object["ServerSideEncryption"],
			"size":                     fmt.Sprint(object["Size"]),
			"squash":                   object["Squash"],
			"state":                    object["State"],
			"support_archive":          object["SupportArchive"],
			"sync_progress":            sync_progressVal,
			"total_download":           fmt.Sprint(object["TotalDownload"]),
			"total_upload":             fmt.Sprint(object["TotalUpload"]),
			"transfer_acceleration":    object["TransferAcceleration"],
			"used":                     fmt.Sprint(object["Used"]),
			"windows_acl":              object["WindowsAcl"],
			"bypass_cache_read":        object["BypassCacheRead"],
		}
		if v, ok := object["ExpressSyncId"]; ok {
			mapping["express_sync_id"] = v
		}
		if v, ok := object["RoUserList"]; ok {
			mapping["ro_user_list"] = v
		}
		if v, ok := object["RoClientList"]; ok {
			mapping["ro_client_list"] = v
		}
		if v, ok := object["ServerSideCmk"]; ok {
			mapping["server_side_cmk"] = v
		}
		if v, ok := object["RwUserList"]; ok {
			mapping["rw_user_list"] = v
		}
		if v, ok := object["RwClientList"]; ok {
			mapping["rw_client_list"] = v
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("shares", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

package cloudstoragegateway

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAlicloudCloudStorageGatewayGatewayBlockVolume() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeCreate,
		Read:   ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead,
		Update: ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeUpdate,
		Delete: ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cache_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Cache", "WriteThrough"}, false),
			},
			"chap_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"chap_in_password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(12, 16),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("chap_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"chap_in_user": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9]{0,32}$`), "The chap_in_user must be 1 to 32 characters in length, and can contain letters and digits."),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("chap_enabled"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"chunk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{131072, 16384, 32768, 65536, 8192}),
			},
			"gateway_block_volume_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9]{0,31}$`), "The name must be 1 to 32 characters in length, and can contain lower case letters and digits."),
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"index_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_source_deletion": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"local_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("cache_mode"); ok && v.(string) == "Cache" {
						return false
					}
					return true
				},
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_bucket_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"oss_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"recovery": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 262144),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGatewayBlockVolume"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("cache_mode"); ok {
		request["CacheMode"] = v
	}
	if v, ok := d.GetOkExists("chap_enabled"); ok {
		request["ChapEnabled"] = v
	}
	if v, ok := d.GetOk("chap_in_user"); ok {
		request["ChapInUser"] = v
	}
	if v, ok := d.GetOk("chap_in_password"); ok {
		request["ChapInPassword"] = v
	}
	if v, ok := d.GetOk("chunk_size"); ok {
		request["ChunkSize"] = v
	}
	request["Name"] = d.Get("gateway_block_volume_name")
	request["GatewayId"] = d.Get("gateway_id")
	if v, ok := d.GetOk("local_path"); ok {
		request["LocalFilePath"] = v
	}
	request["OssBucketName"] = d.Get("oss_bucket_name")
	if v, ok := d.GetOkExists("oss_bucket_ssl"); ok {
		request["OssBucketSsl"] = v
	}
	request["OssEndpoint"] = d.Get("oss_endpoint")
	request["VolumeProtocol"] = d.Get("protocol")
	if v, ok := d.GetOkExists("recovery"); ok {
		request["Recovery"] = v
	}
	if v, ok := d.GetOk("size"); ok {
		request["Size"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_block_volume", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	sgwService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(request["GatewayId"].(string), response["TaskId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	response, err = sgwService.DescribeTasks(request["GatewayId"].(string), response["TaskId"].(string))
	if err != nil {
		d.SetId("")
		return errs.WrapError(err)
	}
	d.SetId(fmt.Sprint(request["GatewayId"], ":", response["RelatedResourceId"]))
	return ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	object, err := sgwService.DescribeCloudStorageGatewayGatewayBlockVolume(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_block_volume sgwService.DescribeCloudStorageGatewayGatewayBlockVolume Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("gateway_id", parts[0])
	d.Set("index_id", parts[1])
	d.Set("cache_mode", object["CacheMode"])
	d.Set("chap_enabled", object["ChapEnabled"])
	d.Set("chap_in_user", object["ChapInUser"])
	if v, ok := object["ChunkSize"]; ok && fmt.Sprint(v) != "0" {
		chunk_sizeVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("chunk_size", chunk_sizeVal)
	}
	d.Set("gateway_block_volume_name", object["Name"])
	d.Set("local_path", object["LocalPath"])
	d.Set("oss_bucket_name", object["OssBucketName"])
	d.Set("oss_bucket_ssl", object["OssBucketSsl"])
	d.Set("oss_endpoint", object["OssEndpoint"])
	d.Set("protocol", object["Protocol"])
	volumeState, err := convert.FormatInt(object["VolumeState"])
	if err != nil {
		return err
	}
	d.Set("status", fmt.Sprint(volumeState))
	return nil
}
func ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"GatewayId": parts[0],
		"IndexId":   parts[1],
	}
	if v, ok := d.GetOkExists("chap_enabled"); ok {
		request["ChapEnabled"] = v
	}
	if d.HasChange("chap_enabled") {
		update = true
	}
	if v, ok := d.GetOk("chap_in_password"); ok {
		request["ChapInPassword"] = v
	}
	if d.HasChange("chap_in_password") {
		update = true
	}
	if v, ok := d.GetOk("chap_in_user"); ok {
		request["ChapInUser"] = v
	}
	if d.HasChange("chap_in_user") {
		update = true
	}
	if d.HasChange("size") {
		update = true
		if v, ok := d.GetOk("size"); ok {
			request["Size"] = v
		}
	}
	if update {
		action := "UpdateGatewayBlockVolume"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		sgwService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(request["GatewayId"].(string), response["TaskId"].(string), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	return ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeRead(d, meta)
}
func ResourceAlicloudCloudStorageGatewayGatewayBlockVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteGatewayBlockVolumes"
	var response map[string]interface{}
	request := map[string]interface{}{
		"GatewayId": parts[0],
		"IndexId":   parts[1],
	}

	if v, ok := d.GetOkExists("is_source_deletion"); ok {
		request["IsSourceDeletion"] = v
	} else {
		request["IsSourceDeletion"] = true
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	sgwService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(request["GatewayId"].(string), response["TaskId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

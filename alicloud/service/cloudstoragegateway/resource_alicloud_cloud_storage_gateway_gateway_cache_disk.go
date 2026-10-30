package cloudstoragegateway

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudCloudStorageGatewayGatewayCacheDisk() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAliCloudCloudStorageGatewayGatewayCacheDiskCreate,
		Read:   ResourceAliCloudCloudStorageGatewayGatewayCacheDiskRead,
		Update: ResourceAliCloudCloudStorageGatewayGatewayCacheDiskUpdate,
		Delete: ResourceAliCloudCloudStorageGatewayGatewayCacheDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cache_disk_size_in_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validators.IntBetween(20, 32768),
			},
			"cache_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd"}, false),
			},
			"performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
			},
			"cache_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func ResourceAliCloudCloudStorageGatewayGatewayCacheDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	action := "CreateGatewayCacheDisk"
	request := make(map[string]interface{})
	var err error

	request["GatewayId"] = d.Get("gateway_id")
	request["CacheDiskSizeInGB"] = d.Get("cache_disk_size_in_gb")

	if v, ok := d.GetOk("cache_disk_category"); ok {
		request["CacheDiskCategory"] = v
	}

	if v, ok := d.GetOk("performance_level"); ok {
		request["PerformanceLevel"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_cache_disk", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	task, err := sgwService.DescribeTasks(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		return errs.WrapError(err)
	}

	object, err := sgwService.DescribeCloudStorageGatewayGatewayCacheDisk(fmt.Sprintf("%v:%v:", request["GatewayId"], task["RelatedResourceId"]))
	if err != nil {
		return errs.WrapError(err)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["GatewayId"], task["RelatedResourceId"], object["LocalFilePath"]))

	return ResourceAliCloudCloudStorageGatewayGatewayCacheDiskRead(d, meta)
}

func ResourceAliCloudCloudStorageGatewayGatewayCacheDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)

	object, err := sgwService.DescribeCloudStorageGatewayGatewayCacheDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_cache_disk sgwService.DescribeCloudStorageGatewayGatewayCacheDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("gateway_id", parts[0])
	cache_disk_size_in_gbVal, err := convert.FormatInt(object["SizeInGB"])
	if err != nil {
		return err
	}
	d.Set("cache_disk_size_in_gb", cache_disk_size_in_gbVal)
	d.Set("cache_disk_category", object["CacheType"])
	d.Set("performance_level", object["PerformanceLevel"])
	d.Set("cache_id", object["CacheId"])
	d.Set("local_file_path", object["LocalFilePath"])
	statusVal, err := convert.FormatInt(object["ExpireStatus"])
	if err != nil {
		return err
	}
	d.Set("status", statusVal)

	return nil
}

func ResourceAliCloudCloudStorageGatewayGatewayCacheDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"GatewayId":     parts[0],
		"LocalFilePath": parts[2],
	}

	if d.HasChange("cache_disk_size_in_gb") {
		update = true
	}
	request["NewSizeInGB"] = d.Get("cache_disk_size_in_gb")

	if update {
		action := "ExpandCacheDisk"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return ResourceAliCloudCloudStorageGatewayGatewayCacheDiskRead(d, meta)
}

func ResourceAliCloudCloudStorageGatewayGatewayCacheDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := NewService(client)
	action := "DeleteGatewayCacheDisk"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"GatewayId":     parts[0],
		"CacheId":       parts[1],
		"LocalFilePath": parts[2],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

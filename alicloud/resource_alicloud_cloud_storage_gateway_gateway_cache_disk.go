package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudStorageGatewayGatewayCacheDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudStorageGatewayGatewayCacheDiskCreate,
		Read:   resourceAliCloudCloudStorageGatewayGatewayCacheDiskRead,
		Update: resourceAliCloudCloudStorageGatewayGatewayCacheDiskUpdate,
		Delete: resourceAliCloudCloudStorageGatewayGatewayCacheDiskDelete,
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
				ValidateFunc: IntBetween(20, 32768),
			},
			"cache_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "cloud_essd"}, false),
			},
			"performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
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

func resourceAliCloudCloudStorageGatewayGatewayCacheDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_cache_disk", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	task, err := sgwService.DescribeTasks(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		return WrapError(err)
	}

	object, err := sgwService.DescribeCloudStorageGatewayGatewayCacheDisk(fmt.Sprintf("%v:%v:", request["GatewayId"], task["RelatedResourceId"]))
	if err != nil {
		return WrapError(err)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["GatewayId"], task["RelatedResourceId"], object["LocalFilePath"]))

	return resourceAliCloudCloudStorageGatewayGatewayCacheDiskRead(d, meta)
}

func resourceAliCloudCloudStorageGatewayGatewayCacheDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}

	object, err := sgwService.DescribeCloudStorageGatewayGatewayCacheDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_cache_disk sgwService.DescribeCloudStorageGatewayGatewayCacheDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	d.Set("gateway_id", parts[0])
	d.Set("cache_disk_size_in_gb", formatInt(object["SizeInGB"]))
	d.Set("cache_disk_category", object["CacheType"])
	d.Set("performance_level", object["PerformanceLevel"])
	d.Set("cache_id", object["CacheId"])
	d.Set("local_file_path", object["LocalFilePath"])
	d.Set("status", formatInt(object["ExpireStatus"]))

	return nil
}

func resourceAliCloudCloudStorageGatewayGatewayCacheDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCloudStorageGatewayGatewayCacheDiskRead(d, meta)
}

func resourceAliCloudCloudStorageGatewayGatewayCacheDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	action := "DeleteGatewayCacheDisk"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"GatewayId":     parts[0],
		"CacheId":       parts[1],
		"LocalFilePath": parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("sgw", "2018-05-11", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

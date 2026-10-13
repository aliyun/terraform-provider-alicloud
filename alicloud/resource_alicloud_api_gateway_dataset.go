// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApiGatewayDataset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayDatasetCreate,
		Read:   resourceAliCloudApiGatewayDatasetRead,
		Update: resourceAliCloudApiGatewayDatasetUpdate,
		Delete: resourceAliCloudApiGatewayDatasetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataset_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"JWT_BLOCKING", "IP_WHITELIST_CIDR", "PARAMETER_ACCESS"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudApiGatewayDatasetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateDataset"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["DatasetName"] = d.Get("dataset_name")
	request["DatasetType"] = d.Get("dataset_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_dataset", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DatasetId"]))

	return resourceAliCloudApiGatewayDatasetUpdate(d, meta)
}

func resourceAliCloudApiGatewayDatasetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayDataset(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_dataset DescribeApiGatewayDataset Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dataset_id", objectRaw["DatasetId"])
	d.Set("dataset_name", objectRaw["DatasetName"])
	d.Set("dataset_type", objectRaw["DatasetType"])
	d.Set("description", objectRaw["Description"])
	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("modified_time", objectRaw["ModifiedTime"])

	tagsMap, err := apiGatewayServiceV2.ListApiGatewayDatasetTags(d.Id())
	if err != nil {
		log.Printf("[WARN] Failed to read tags for ApiGateway Dataset %s: %s", d.Id(), err)
	}
	d.Set("tags", tagsMap)

	return nil
}

func resourceAliCloudApiGatewayDatasetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	if d.HasChange("dataset_name") || d.HasChange("description") {
		action := "ModifyDataset"
		var request map[string]interface{}
		var response map[string]interface{}
		var query map[string]interface{}
		var err error
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["DatasetId"] = d.Id()
		request["DatasetName"] = d.Get("dataset_name")
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		if err := apiGatewayServiceV2.SetResourceTags(d, "dataset"); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudApiGatewayDatasetRead(d, meta)
}

func resourceAliCloudApiGatewayDatasetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDataset"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DatasetId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundDataset"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

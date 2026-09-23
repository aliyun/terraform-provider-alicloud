// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksDwResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksDwResourceGroupCreate,
		Read:   resourceAliCloudDataWorksDwResourceGroupRead,
		Update: resourceAliCloudDataWorksDwResourceGroupUpdate,
		Delete: resourceAliCloudDataWorksDwResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"payment_duration_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"specification": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudDataWorksDwResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["Remark"] = d.Get("remark")
	request["PaymentType"] = d.Get("payment_type")
	if v, ok := d.GetOkExists("specification"); ok {
		request["Spec"] = v
	}
	request["Name"] = d.Get("resource_group_name")
	request["VpcId"] = d.Get("default_vpc_id")
	request["VswitchId"] = d.Get("default_vswitch_id")
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request["PaymentDurationUnit"] = v
	}
	if v, ok := d.GetOkExists("payment_duration"); ok {
		request["PaymentDuration"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["AliyunResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenewEnabled"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_dw_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.ResourceGroupOrder.Id", response)
	d.SetId(fmt.Sprint(id))

	dataWorksServiceV2 := DataWorksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, dataWorksServiceV2.DataWorksDwResourceGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDataWorksDwResourceGroupUpdate(d, meta)
}

func resourceAliCloudDataWorksDwResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDwResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_dw_resource_group DescribeDataWorksDwResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["DefaultVpcId"] != nil {
		d.Set("default_vpc_id", objectRaw["DefaultVpcId"])
	}
	if objectRaw["DefaultVswitchId"] != nil {
		d.Set("default_vswitch_id", objectRaw["DefaultVswitchId"])
	}
	if objectRaw["PaymentType"] != nil {
		d.Set("payment_type", objectRaw["PaymentType"])
	}
	if objectRaw["Remark"] != nil {
		d.Set("remark", objectRaw["Remark"])
	}
	if objectRaw["AliyunResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["AliyunResourceGroupId"])
	}
	if objectRaw["Name"] != nil {
		d.Set("resource_group_name", objectRaw["Name"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	tagsMaps := objectRaw["AliyunResourceTags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudDataWorksDwResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "UpdateResourceGroup"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
	}
	request["Remark"] = d.Get("remark")
	if !d.IsNewResource() && d.HasChange("resource_group_name") {
		update = true
		request["Name"] = d.Get("resource_group_name")
	}

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["AliyunResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
	}

	if d.HasChange("tags") {
		dataWorksServiceV2 := DataWorksServiceV2{client}
		if err := dataWorksServiceV2.SetResourceTags(d, "dwresourcegroup"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudDataWorksDwResourceGroupRead(d, meta)
}

func resourceAliCloudDataWorksDwResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

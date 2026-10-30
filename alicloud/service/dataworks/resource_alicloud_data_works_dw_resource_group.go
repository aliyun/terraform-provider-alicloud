// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package dataworks

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDwResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDwResourceGroupCreate,
		Read:   resourceDwResourceGroupRead,
		Update: resourceDwResourceGroupUpdate,
		Delete: resourceDwResourceGroupDelete,
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
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceDwResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

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
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenewEnabled"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_data_works_dw_resource_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.ResourceGroupOrder.Id", response)
	d.SetId(fmt.Sprint(id))

	dataWorksServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, dataWorksServiceV2.DataWorksDwResourceGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDwResourceGroupUpdate(d, meta)
}

func resourceDwResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := NewService(client)

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDwResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_dw_resource_group DescribeDataWorksDwResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceDwResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
	}

	if d.HasChange("tags") {
		dataWorksServiceV2 := NewService(client)
		if err := dataWorksServiceV2.SetResourceTags(d, "dwresourcegroup"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceDwResourceGroupRead(d, meta)
}

func resourceDwResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

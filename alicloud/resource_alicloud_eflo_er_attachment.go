// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEfloErAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloErAttachmentCreate,
		Read:   resourceAliCloudEfloErAttachmentRead,
		Update: resourceAliCloudEfloErAttachmentUpdate,
		Delete: resourceAliCloudEfloErAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"across": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_receive_all_route": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"er_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"er_attachment_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"er_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VPD", "VCC"}, false),
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEfloErAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateErAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ErId"] = d.Get("er_id")
	request["RegionId"] = client.RegionId
	request["InstanceType"] = d.Get("instance_type")
	request["InstanceId"] = d.Get("instance_id")
	request["ErAttachmentName"] = d.Get("er_attachment_name")
	request["AutoReceiveAllRoute"] = d.Get("auto_receive_all_route")
	if v, ok := d.GetOk("resource_tenant_id"); ok {
		request["ResourceTenantId"] = v
	}

	// The Lingjun HUB must finish provisioning before an attachment can be created.
	efloServiceV2 := EfloServiceV2{client}
	erStateConf := BuildStateConf([]string{"Executing"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloErStateRefreshFunc(fmt.Sprint(request["ErId"]), "Status", []string{"Not Available"}))
	if _, err := erStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, fmt.Sprint(request["ErId"]))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_er_attachment", action, AlibabaCloudSdkGoERROR)
	}

	contentErAttachmentId, err := jsonpath.Get("$.Content.ErAttachmentId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Content.ErAttachmentId", response)
	}
	d.SetId(fmt.Sprintf("%v:%v", request["ErId"], contentErAttachmentId))

	stateConf := BuildStateConf([]string{"Executing", "Deleting"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloErAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEfloErAttachmentRead(d, meta)
}

func resourceAliCloudEfloErAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloErAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_er_attachment DescribeEfloErAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("across", objectRaw["Across"])
	d.Set("auto_receive_all_route", objectRaw["AutoReceiveAllRoute"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("er_attachment_name", objectRaw["ErAttachmentName"])
	d.Set("gmt_modified", objectRaw["GmtModified"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("message", objectRaw["Message"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("resource_tenant_id", objectRaw["ResourceTenantId"])
	d.Set("status", objectRaw["Status"])
	d.Set("tenant_id", objectRaw["TenantId"])
	d.Set("er_attachment_id", objectRaw["ErAttachmentId"])
	d.Set("er_id", objectRaw["ErId"])

	return nil
}

func resourceAliCloudEfloErAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateErAttachment"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if len(parts) < 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected format <ErId>:<ErAttachmentId>", d.Id()))
	}
	request["ErAttachmentId"] = parts[1]
	request["ErId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("er_attachment_name") {
		update = true
		request["ErAttachmentName"] = d.Get("er_attachment_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"1042", "1010"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		efloServiceV2 := EfloServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, efloServiceV2.EfloErAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEfloErAttachmentRead(d, meta)
}

func resourceAliCloudEfloErAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) < 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected format <ErId>:<ErAttachmentId>", d.Id()))
	}
	action := "DeleteErAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ErAttachmentId"] = parts[1]
	request["ErId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"1010"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"1003"}) || NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	efloServiceV2 := EfloServiceV2{client}
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"Exists"},
		Target:     []string{"Deleted"},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			object, err := efloServiceV2.DescribeEfloErAttachment(d.Id())
			if err != nil {
				if NotFoundError(err) {
					return "deleted", "Deleted", nil
				}
				return nil, "", WrapError(err)
			}
			return object, "Exists", nil
		},
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), action, ProviderERROR)
	}

	d.SetId("")
	return nil
}

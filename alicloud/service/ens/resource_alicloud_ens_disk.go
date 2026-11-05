// Package ens. This file is generated automatically. Please do not modify it manually, thank you!
package ens

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnsDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnsDiskCreate,
		Read:   resourceEnsDiskRead,
		Update: resourceEnsDiskUpdate,
		Delete: resourceEnsDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(7 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"cloud_efficiency", "cloud_ssd"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceEnsDiskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["InstanceChargeType"] = convertEnsInstanceInstanceChargeTypeRequest(d.Get("payment_type").(string))
	request["EnsRegionId"] = d.Get("ens_region_id")
	request["Category"] = d.Get("category")
	if v, ok := d.GetOkExists("size"); ok {
		request["Size"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["Encrypted"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}
	if v, ok := d.GetOk("disk_name"); ok {
		request["DiskName"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ens_disk", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.InstanceIds[0]", response)
	d.SetId(fmt.Sprint(id))

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 4*time.Minute, ensServiceV2.EnsDiskStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceEnsDiskRead(d, meta)
}

func resourceEnsDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	objectRaw, err := ensServiceV2.DescribeEnsDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_disk DescribeEnsDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("category", objectRaw["Category"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("disk_name", objectRaw["DiskName"])
	d.Set("encrypted", objectRaw["Encrypted"])
	d.Set("ens_region_id", objectRaw["EnsRegionId"])
	d.Set("kms_key_id", objectRaw["EncryptedKeyId"])
	d.Set("payment_type", convertEnsDisksDisksDiskChargeTypeResponse(objectRaw["DiskChargeType"]))
	d.Set("snapshot_id", objectRaw["SnapshotId"])
	d.Set("status", objectRaw["Status"])

	e := jsonata.MustCompile("$number($.Size & '') / 1024")
	evaluation, _ := e.Eval(objectRaw)
	sizeVal, formatErr := convert.FormatInt(evaluation)
	if formatErr != nil {
		return formatErr
	}
	d.Set("size", sizeVal)

	objectRaw, err = ensServiceV2.DescribeDiskListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceEnsDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ResizeDisk"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	if d.HasChange("size") {
		update = true
	}
	if v, ok := d.GetOkExists("size"); ok {
		request["NewSize"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		ensServiceV2 := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("size"))}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, ensServiceV2.EnsDiskStateRefreshFunc(d.Id(), "$.Size", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDiskAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	if d.HasChange("disk_name") {
		update = true
	}
	if v, ok := d.GetOk("disk_name"); ok {
		request["DiskName"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		ensServiceV2 := NewService(client)
		if err := ensServiceV2.SetResourceTags(d, "disk"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceEnsDiskRead(d, meta)
}

func resourceEnsDiskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DiskId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsDiskStateRefreshFunc(d.Id(), "$.DiskId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertEnsDisksDisksDiskChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "postpay":
		return "PayAsYouGo"
	}
	return source
}

func convertEnsDiskChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}

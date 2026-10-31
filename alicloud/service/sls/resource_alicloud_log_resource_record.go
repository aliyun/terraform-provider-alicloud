package sls

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudLogResourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogResourceRecordCreate,
		Read:   resourceAlicloudLogResourceRecordRead,
		Update: resourceAlicloudLogResourceRecordUpdate,
		Delete: resourceAlicloudLogResourceRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"record_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudLogResourceRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	recordId := d.Get("record_id").(string)
	tag := d.Get("tag").(string)
	value := d.Get("value").(string)
	recordResourceName := d.Get("resource_name").(string)

	record := &sls.ResourceRecord{
		Id:    recordId,
		Tag:   tag,
		Value: value,
	}
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateResourceRecord(recordResourceName, record)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_resource_record", "CreateResourceRecord", errs.AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", recordResourceName, COLON_SEPARATED, recordId))
	return resourceAlicloudLogResourceRecordRead(d, meta)
}

func resourceAlicloudLogResourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	resourceName := parts[0]
	object, err := logService.DescribeLogResourceRecord(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("record_id", object.Id)
	d.Set("tag", object.Tag)
	d.Set("value", object.Value)
	d.Set("resource_name", resourceName)
	return nil
}

func resourceAlicloudLogResourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	resourceName, recordId := parts[0], parts[1]
	client := meta.(*connectivity.AliyunClient)
	params := &sls.ResourceRecord{
		Id:    recordId,
		Tag:   d.Get("tag").(string),
		Value: d.Get("value").(string),
	}

	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateResourceRecord(resourceName, params)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UpdateResourceRecord", errs.AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogResourceRecordRead(d, meta)
}

func resourceAlicloudLogResourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	var requestInfo *sls.Client
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteResourceRecord(parts[0], parts[1])
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("DeleteResourceRecord", raw, requestInfo, map[string]interface{}{
				"resource_name": parts[0],
				"id":            parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_resource_record", "DeleteResourceRecord", errs.AliyunLogGoSdkERROR)
	}
	return errs.WrapError(logService.WaitForLogResourceRecord(d.Id(), Deleted, DefaultTimeout))
}

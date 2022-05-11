package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
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
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateResourceRecord(recordResourceName, record)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_resource_record", "CreateResourceRecord", AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", recordResourceName, COLON_SEPARATED, recordId))
	return resourceAlicloudLogResourceRecordRead(d, meta)
}

func resourceAlicloudLogResourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	resourceName := parts[0]
	object, err := logService.DescribeLogResourceRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("record_id", object.Id)
	d.Set("tag", object.Tag)
	d.Set("value", object.Value)
	d.Set("resource_name", resourceName)
	return nil
}

func resourceAlicloudLogResourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	resourceName, recordId := parts[0], parts[1]
	client := meta.(*connectivity.AliyunClient)
	params := &sls.ResourceRecord{
		Id:    recordId,
		Tag:   d.Get("tag").(string),
		Value: d.Get("value").(string),
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateResourceRecord(resourceName, params)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateResourceRecord", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogResourceRecordRead(d, meta)
}

func resourceAlicloudLogResourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteResourceRecord(parts[0], parts[1])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteResourceRecord", raw, requestInfo, map[string]interface{}{
				"resource_name": parts[0],
				"id":            parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_resource_record", "DeleteResourceRecord", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogResourceRecord(d.Id(), Deleted, DefaultTimeout))
}

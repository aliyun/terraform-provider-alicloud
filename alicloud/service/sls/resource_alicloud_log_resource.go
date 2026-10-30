package sls

import (
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

func resourceAlicloudLogResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogResourceCreate,
		Read:   resourceAlicloudLogResourceRead,
		Update: resourceAlicloudLogResourceUpdate,
		Delete: resourceAlicloudLogResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"ext_info": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudLogResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceName := d.Get("name").(string)
	resourceType := d.Get("type").(string)
	resourceSchema := d.Get("schema").(string)
	description := d.Get("description").(string)
	extInfo := d.Get("ext_info").(string)

	record := &sls.Resource{
		Name:        resourceName,
		Type:        resourceType,
		Schema:      resourceSchema,
		Description: description,
		ExtInfo:     extInfo,
	}
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateResource(record)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_resource", "CreateResource", errs.AliyunLogGoSdkERROR)
	}

	d.SetId(resourceName)
	return resourceAlicloudLogResourceRead(d, meta)
}

func resourceAlicloudLogResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}

	resourceName := d.Id()
	object, err := logService.DescribeLogResource(resourceName)
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("type", object.Type)
	d.Set("schema", object.Schema)
	d.Set("ext_info", object.ExtInfo)
	d.Set("description", object.Description)
	return nil
}

func resourceAlicloudLogResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	params := &sls.Resource{
		Name:        d.Id(),
		Type:        d.Get("type").(string),
		Schema:      d.Get("schema").(string),
		ExtInfo:     d.Get("ext_info").(string),
		Description: d.Get("description").(string),
	}

	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateResource(params)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UpdateResource", errs.AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogResourceRead(d, meta)
}

func resourceAlicloudLogResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	var requestInfo *sls.Client
	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteResource(d.Id())
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("DeleteResource", raw, requestInfo, map[string]interface{}{
				"resource_name": d.Id(),
			})
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_resource", "DeleteResource", errs.AliyunLogGoSdkERROR)
	}
	return errs.WrapError(logService.WaitForLogResource(d.Id(), Deleted, DefaultTimeout))
}

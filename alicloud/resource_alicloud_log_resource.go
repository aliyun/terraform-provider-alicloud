package alicloud

import (
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"ext_info": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
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
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateResource(record)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_resource", "CreateResource", AliyunLogGoSdkERROR)
	}

	d.SetId(resourceName)
	return resourceAlicloudLogResourceRead(d, meta)
}

func resourceAlicloudLogResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	resourceName := d.Id()
	object, err := logService.DescribeLogResource(resourceName)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.UpdateResource(params)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateResource", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogResourceRead(d, meta)
}

func resourceAlicloudLogResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	var requestInfo *sls.Client
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteResource(d.Id())
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteResource", raw, requestInfo, map[string]interface{}{
				"resource_name": d.Id(),
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_resource", "DeleteResource", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogResource(d.Id(), Deleted, DefaultTimeout))
}

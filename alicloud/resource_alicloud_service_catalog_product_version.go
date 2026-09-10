// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudServiceCatalogProductVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudServiceCatalogProductVersionCreate,
		Read:   resourceAliCloudServiceCatalogProductVersionRead,
		Update: resourceAliCloudServiceCatalogProductVersionUpdate,
		Delete: resourceAliCloudServiceCatalogProductVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guidance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "Administrator guidance"),
			},
			"product_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "Product ID"),
			},
			"product_version_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "Template Type"),
			},
			"template_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudServiceCatalogProductVersionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateProductVersion"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("active"); ok {
		request["Active"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("guidance"); ok {
		request["Guidance"] = v
	}
	request["ProductId"] = d.Get("product_id")
	request["ProductVersionName"] = d.Get("product_version_name")
	request["TemplateType"] = d.Get("template_type")
	request["TemplateUrl"] = d.Get("template_url")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("servicecatalog", "2021-09-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_catalog_product_version", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ProductVersionId"]))

	return resourceAliCloudServiceCatalogProductVersionRead(d, meta)
}

func resourceAliCloudServiceCatalogProductVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serviceCatalogServiceV2 := ServiceCatalogServiceV2{client}

	objectRaw, err := serviceCatalogServiceV2.DescribeServiceCatalogProductVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_catalog_product_version DescribeServiceCatalogProductVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Active"] != nil {
		d.Set("active", objectRaw["Active"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["Guidance"] != nil {
		d.Set("guidance", objectRaw["Guidance"])
	}
	if objectRaw["ProductId"] != nil {
		d.Set("product_id", objectRaw["ProductId"])
	}
	if objectRaw["ProductVersionName"] != nil {
		d.Set("product_version_name", objectRaw["ProductVersionName"])
	}
	if objectRaw["TemplateType"] != nil {
		d.Set("template_type", objectRaw["TemplateType"])
	}
	if objectRaw["TemplateUrl"] != nil {
		d.Set("template_url", objectRaw["TemplateUrl"])
	}

	return nil
}

func resourceAliCloudServiceCatalogProductVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateProductVersion"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ProductVersionId"] = d.Id()
	query["RegionId"] = client.RegionId
	if d.HasChange("active") {
		update = true
		request["Active"] = d.Get("active")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("guidance") {
		update = true
		request["Guidance"] = d.Get("guidance")
	}

	if d.HasChange("product_version_name") {
		update = true
	}
	request["ProductVersionName"] = d.Get("product_version_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicecatalog", "2021-09-01", action, query, request, true)
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

	return resourceAliCloudServiceCatalogProductVersionRead(d, meta)
}

func resourceAliCloudServiceCatalogProductVersionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteProductVersion"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ProductVersionId"] = d.Id()
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("servicecatalog", "2021-09-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

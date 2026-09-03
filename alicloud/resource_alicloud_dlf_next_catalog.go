// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudDlfNextCatalog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDlfNextCatalogCreate,
		Read:   resourceAliCloudDlfNextCatalogRead,
		Update: resourceAliCloudDlfNextCatalogUpdate,
		Delete: resourceAliCloudDlfNextCatalogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PAIMON", "ICEBERG"}, false),
			},
			"options": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDlfNextCatalogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/dlf/v1/catalogs"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	request = make(map[string]interface{})
	request["name"] = d.Get("name")
	request["type"] = d.Get("type")
	if v, ok := d.GetOk("options"); ok {
		request["options"] = v
	}
	if v, ok := d.GetOkExists("is_shared"); ok {
		request["isShared"] = v
	}
	if v, ok := d.GetOk("share_id"); ok {
		request["shareId"] = v
	}
	body = request

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("DlfNext", "2025-03-10", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dlf_next_catalog", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["name"]))

	return resourceAliCloudDlfNextCatalogRead(d, meta)
}

func resourceAliCloudDlfNextCatalogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dlfNextServiceV2 := DlfNextServiceV2{client}

	objectRaw, err := dlfNextServiceV2.DescribeDlfNextCatalog(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dlf_next_catalog DescribeDlfNextCatalog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", objectRaw["name"])
	d.Set("type", objectRaw["type"])
	d.Set("options", objectRaw["options"])
	d.Set("is_shared", objectRaw["isShared"])
	d.Set("share_id", objectRaw["shareId"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("status", objectRaw["status"])
	d.Set("owner", objectRaw["owner"])
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("created_by", objectRaw["createdBy"])
	d.Set("updated_at", objectRaw["updatedAt"])
	d.Set("updated_by", objectRaw["updatedBy"])

	return nil
}

func resourceAliCloudDlfNextCatalogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	catalog := d.Id()
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error

	if d.HasChange("options") {
		oldRaw, newRaw := d.GetChange("options")
		oldOptions, _ := oldRaw.(map[string]interface{})
		newOptions, _ := newRaw.(map[string]interface{})

		updates := make(map[string]interface{})
		for k, v := range newOptions {
			updates[k] = v
		}

		removals := make([]interface{}, 0)
		for k := range oldOptions {
			if _, exists := newOptions[k]; !exists {
				removals = append(removals, k)
			}
		}

		body["updates"] = updates
		if len(removals) > 0 {
			body["removals"] = removals
		}

		request = body
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("DlfNext", "2025-03-10", action, query, nil, body, true)
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

	return resourceAliCloudDlfNextCatalogRead(d, meta)
}

func resourceAliCloudDlfNextCatalogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	catalog := d.Id()
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, true)
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

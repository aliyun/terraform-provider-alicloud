package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDataWorksMetaCategory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDataWorksMetaCategoryCreate,
		Read:   resourceAlicloudDataWorksMetaCategoryRead,
		Update: resourceAlicloudDataWorksMetaCategoryUpdate,
		Delete: resourceAlicloudDataWorksMetaCategoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"category_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_category_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDataWorksMetaCategoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMetaCategory"
	request := make(map[string]interface{})
	var err error

	request["Name"] = d.Get("name").(string)
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	// The CreateMetaCategory API accepts "ParentId" (not "ParentCategoryId") for the parent reference.
	// This naming mismatch with GetMetaCategory's "ParentCategoryId" parameter is an upstream API
	// inconsistency and is preserved here intentionally.
	parentCategoryId := d.Get("parent_category_id").(int)
	request["ParentId"] = parentCategoryId

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_meta_category", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_data_works_meta_category", "$.Data", response)
	}
	dataMap, ok := v.(map[string]interface{})
	if !ok || len(dataMap) < 1 {
		return WrapErrorf(NotFoundErr("dataworks meta category", "create response"), NotFoundWithResponse, response)
	}
	categoryIdRaw, ok := dataMap["CategoryId"]
	if !ok || categoryIdRaw == nil {
		return WrapErrorf(NotFoundErr("dataworks meta category", "CategoryId"), NotFoundWithResponse, response)
	}
	// RpcPost responses are decoded with json.UseNumber, so CategoryId arrives as
	// json.Number (not float64); a direct .(float64) assertion panics. toInt handles
	// json.Number / float64 / string uniformly.
	categoryId, err := toInt(categoryIdRaw)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_data_works_meta_category", "$.Data.CategoryId", response)
	}

	d.SetId(fmt.Sprintf("%d:%d", categoryId, parentCategoryId))
	d.Set("category_id", categoryId)

	return resourceAlicloudDataWorksMetaCategoryRead(d, meta)
}

func resourceAlicloudDataWorksMetaCategoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksMetaCategory(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_meta_category DescribeDataWorksMetaCategory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	categoryId, _ := strconv.Atoi(parts[0])
	parentCategoryId, _ := strconv.Atoi(parts[1])

	d.Set("category_id", categoryId)
	d.Set("parent_category_id", parentCategoryId)
	d.Set("name", object["Name"])
	d.Set("comment", object["Comment"])
	if v, ok := object["CreateTime"]; ok && v != nil {
		// CreateTime from GetMetaCategory is a json.Number (UseNumber decoding); the
		// previous float64 assertion silently failed and left create_time unset,
		// which caused refresh plan-not-empty drift after create. toInt handles the
		// json.Number case so create_time is reliably back-filled.
		if cv, err := toInt(v); err == nil {
			d.Set("create_time", cv)
		}
	}
	return nil
}

func resourceAlicloudDataWorksMetaCategoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	categoryId, _ := strconv.Atoi(parts[0])

	// UpdateMetaCategory only accepts CategoryId, Name and Comment. The API does not support
	// changing the parent of an existing category, so parent_category_id is ForceNew.
	if d.HasChange("name") || d.HasChange("comment") {
		request := map[string]interface{}{
			"CategoryId": categoryId,
		}
		// UpdateMetaCategory requires Name on every call: the API returns
		// InvalidParameter.Meta.CommonError (400) when Name is omitted, even
		// though the OpenAPI help marks it optional. Always send Name so that
		// updating only the comment does not fail with a missing-Name error.
		request["Name"] = d.Get("name")
		if d.HasChange("comment") {
			// Use HasChange instead of GetOk so clearing the comment (empty
			// string) actually sends an empty Comment to the API. GetOk returns
			// (false, "") for an empty string, which previously dropped the
			// field from the request and left the old comment in place.
			request["Comment"] = d.Get("comment")
		}
		action := "UpdateMetaCategory"
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
	return resourceAlicloudDataWorksMetaCategoryRead(d, meta)
}

func resourceAlicloudDataWorksMetaCategoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	categoryId, _ := strconv.Atoi(parts[0])
	action := "DeleteMetaCategory"
	var response map[string]interface{}
	request := map[string]interface{}{
		"CategoryId": categoryId,
	}
	// DeleteMetaCategory only accepts GET (its CloudSpec @http annotation declares
	// methods: ["get"]). The previous RpcPost sent a POST body and the gateway
	// rejected it with 403 UnsupportedHTTPMethod, leaving every deleted category as
	// an orphan. RpcGet puts CategoryId in the query string and uses GET, which the
	// gateway accepts. NotFoundError is still treated as a successful delete so the
	// resource is removed from state when the category is already gone.
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcGet("dataworks-public", "2020-05-18", action, request, nil)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
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

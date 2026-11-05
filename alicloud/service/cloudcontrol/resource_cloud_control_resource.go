package cloudcontrol

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCloudControlResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudControlResourceCreate,
		Read:   resourceCloudControlResourceRead,
		Update: resourceCloudControlResourceUpdate,
		Delete: resourceCloudControlResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"desire_attributes": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"product": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceCloudControlResourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := genResourceAction("aliyun", d.Get("product").(string), d.Get("resource_code").(string), d.Get("resource_id").(string))

	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["regionId"] = helper.StringPointer(client.RegionId)
	request["body"] = convert.ConvertJsonStringToObject(d.Get("desire_attributes"))
	body = request["body"].(map[string]interface{})
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("cloudcontrol", "2022-08-30", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_control_resource", action, errs.AlibabaCloudSdkGoERROR)
	}

	resourceIdVar, _ := jsonpath.Get("$.resourceId", response)
	d.SetId(convertActionToId(fmt.Sprintf("%v/%v", action, resourceIdVar)))

	return resourceCloudControlResourceRead(d, meta)
}

func resourceCloudControlResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeCloudControlResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_control_resource DescribeCloudControlResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["resourceAttributes"] != nil {
		d.Set("resource_attributes", convert.ConvertObjectToJsonString(objectRaw["resourceAttributes"]))
	}

	action := convertIdToAction(d.Id())
	d.Set("product", parseProductCodeFromAction(action))
	d.Set("resource_code", parseResourceCodeFromAction(action))
	d.Set("resource_id", parseParentIdFromAction(action))

	return nil
}

func resourceCloudControlResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	action := convertIdToAction(d.Id())
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["regionId"] = helper.StringPointer(client.RegionId)
	if d.HasChange("desire_attributes") {
		update = true
	}
	if v, ok := d.GetOk("desire_attributes"); ok || d.HasChange("desire_attributes") {
		request["body"] = v
	}
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("cloudcontrol", "2022-08-30", action, query, nil, body, true)
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

	return resourceCloudControlResourceRead(d, meta)
}

func resourceCloudControlResourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := convertIdToAction(d.Id())
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["regionId"] = helper.StringPointer(client.RegionId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("cloudcontrol", "2022-08-30", action, query, nil, nil, true)

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

func genResourceAction(provider, product, resourceCodes, parentResourceIds string) string {
	codeParts := strings.Split(resourceCodes, "::")
	idParts := strings.Split(parentResourceIds, ":")
	if len(codeParts) == 0 || len(idParts) == 0 {
		return ""
	}
	path := fmt.Sprintf("/api/v1/providers/%s/products/%s/resources", provider, product)
	if parentResourceIds != "" {
		for i := 0; i < len(idParts); i++ {
			path += "/" + codeParts[i] + "/" + idParts[i]
		}
	}
	path += "/" + codeParts[len(codeParts)-1]
	return path
}

func parseProductCodeFromAction(action string) string {
	pathParts := strings.Split(action, "/")
	if len(pathParts) < 7 {
		return ""
	}
	return pathParts[6]
}

func parseResourceCodeFromAction(action string) string {
	resourceCode := ""
	pathParts := strings.Split(action, "/")
	if len(pathParts) < 9 {
		return resourceCode
	}
	resourceCode += pathParts[8]
	for i := 10; i < len(pathParts); i++ {
		if i%2 == 0 {
			resourceCode += "::" + pathParts[i]
		}
	}
	return resourceCode
}

func parseParentIdFromAction(action string) string {
	id := convertActionToId(action)
	pathParts := strings.Split(id, ":")

	if len(pathParts) < 5 {
		return pathParts[3]
	}

	parentId := pathParts[3]
	for i := 5; i < (len(pathParts) - 2); i++ {
		if i%2 == 1 {
			id += ":" + pathParts[i]
		}
	}
	return parentId
}

func convertIdToAction(id string) string {
	idParts := strings.Split(id, ":")
	providerName := idParts[0]
	product := idParts[1]
	return "/api/v1/providers/" + providerName + "/products/" + product + "/resources/" + strings.Join(idParts[2:], "/")
}

func convertActionToId(action string) string {
	pathParts := strings.Split(action, "/")
	id := pathParts[4]
	for i := 6; i < len(pathParts); i++ {
		if i != 7 {
			id += ":" + pathParts[i]
		}
	}
	return id
}

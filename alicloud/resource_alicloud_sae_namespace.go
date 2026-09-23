package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeNamespaceCreate,
		Read:   resourceAlicloudSaeNamespaceRead,
		Update: resourceAlicloudSaeNamespaceUpdate,
		Delete: resourceAlicloudSaeNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"namespace_id", "namespace_short_id"},
			},
			"namespace_short_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"namespace_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_micro_registration": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSaeNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/paas/namespace"
	request := make(map[string]*string)
	var err error

	request["NamespaceName"] = StringPointer(d.Get("namespace_name").(string))

	if v, ok := d.GetOk("namespace_id"); ok {
		request["NamespaceId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("namespace_short_id"); ok {
		request["NameSpaceShortId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("namespace_description"); ok {
		request["NamespaceDescription"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOkExists("enable_micro_registration"); ok {
		request["EnableMicroRegistration"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_namespace", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if responseData, ok := response["Data"].(map[string]interface{}); ok {
		d.SetId(fmt.Sprint(responseData["NamespaceId"]))
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}

	return resourceAlicloudSaeNamespaceRead(d, meta)
}

func resourceAlicloudSaeNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_namespace saeService.DescribeSaeNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("namespace_name", object["NamespaceName"])
	d.Set("namespace_id", object["NamespaceId"])
	d.Set("namespace_short_id", object["NameSpaceShortId"])
	d.Set("namespace_description", object["NamespaceDescription"])
	d.Set("enable_micro_registration", object["EnableMicroRegistration"])

	return nil
}

func resourceAlicloudSaeNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]*string{
		"NamespaceId": StringPointer(d.Id()),
	}

	if d.HasChange("namespace_name") {
		update = true
	}
	request["NamespaceName"] = StringPointer(d.Get("namespace_name").(string))

	if d.HasChange("namespace_description") {
		update = true
	}
	if v, ok := d.GetOk("namespace_description"); ok {
		request["NamespaceDescription"] = StringPointer(v.(string))
	}

	if d.HasChange("enable_micro_registration") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_micro_registration"); ok {
		request["EnableMicroRegistration"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if update {
		action := "/pop/v1/paas/namespace"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, true)
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
			if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) {
				return WrapErrorf(NotFoundErr("SAE:Namespace", d.Id()), NotFoundMsg, ProviderERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudSaeNamespaceRead(d, meta)
}

func resourceAlicloudSaeNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/paas/namespace"
	var response map[string]interface{}
	var err error
	request := map[string]*string{
		"NamespaceId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

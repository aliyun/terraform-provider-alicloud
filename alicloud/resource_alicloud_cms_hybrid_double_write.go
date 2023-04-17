package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsHybridDoubleWrite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsHybridDoubleWriteCreate,
		Read:   resourceAlicloudCmsHybridDoubleWriteRead,
		Delete: resourceAlicloudCmsHybridDoubleWriteDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"namespace": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"source_namespace": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"source_user_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"user_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCmsHybridDoubleWriteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}
	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
	}
	if v, ok := d.GetOk("source_namespace"); ok {
		request["SourceNamespace"] = v
	}
	accountId, err := meta.(*connectivity.AliyunClient).AccountId()

	if err != nil {
		return WrapError(err)
	}
	request["SourceUserId"] = accountId

	var response map[string]interface{}
	action := "CreateHybridDoubleWrite"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_hybrid_double_write", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["SourceNamespace"]))

	return resourceAlicloudCmsHybridDoubleWriteRead(d, meta)
}

func resourceAlicloudCmsHybridDoubleWriteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	object, err := cmsService.DescribeCmsHybridDoubleWrite(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_hybrid_double_write cmsService.DescribeCmsHybridDoubleWrite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("source_namespace", object["SourceNamespace"])
	d.Set("namespace", object["Namespace"])
	d.Set("source_user_id", object["SourceUserId"])
	d.Set("user_id", object["UserId"])

	return nil
}

func resourceAlicloudCmsHybridDoubleWriteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"SourceNamespace": d.Id(),
	}

	action := "DeleteHybridDoubleWrite"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

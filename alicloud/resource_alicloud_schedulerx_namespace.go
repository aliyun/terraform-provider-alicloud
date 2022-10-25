package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSchedulerxNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSchedulerxNamespaceCreate,
		Read:   resourceAlicloudSchedulerxNamespaceRead,
		Update: resourceAlicloudSchedulerxNamespaceUpdate,
		Delete: resourceAlicloudSchedulerxNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"namespace_name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudSchedulerxNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	conn, err := client.NewEdasschedulerxClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("namespace_name"); ok {
		request["Name"] = v
	}

	var response map[string]interface{}
	action := "CreateNamespace"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_schedulerx_namespace", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Data.NamespaceUid", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_schedulerx_namespace")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudSchedulerxNamespaceRead(d, meta)
}

func resourceAlicloudSchedulerxNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	schedulerx2Service := Schedulerx2Service{client}

	object, err := schedulerx2Service.DescribeSchedulerxNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_schedulerx_namespace schedulerx2Service.DescribeSchedulerxNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("namespace_name", object["Name"])

	return nil
}

func resourceAlicloudSchedulerxNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEdasschedulerxClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"Uid":      d.Id(),
		"RegionId": client.RegionId,
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("namespace_name") {
		update = true
	}
	request["Name"] = d.Get("namespace_name")

	if update {
		action := "CreateNamespace"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudSchedulerxNamespaceRead(d, meta)
}

func resourceAlicloudSchedulerxNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudSchedulerxNamespace. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

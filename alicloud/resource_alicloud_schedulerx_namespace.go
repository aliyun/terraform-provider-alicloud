// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSchedulerxNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSchedulerxNamespaceCreate,
		Read:   resourceAliCloudSchedulerxNamespaceRead,
		Delete: resourceAliCloudSchedulerxNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_uid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudSchedulerxNamespaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNamespace"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("namespace_uid"); ok {
		request["Uid"] = v
	}
	request["RegionId"] = client.RegionId

	request["Name"] = d.Get("namespace_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("schedulerx2", "2019-04-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_schedulerx_namespace", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.NamespaceUid", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudSchedulerxNamespaceRead(d, meta)
}

func resourceAliCloudSchedulerxNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	schedulerxServiceV2 := SchedulerxServiceV2{client}

	objectRaw, err := schedulerxServiceV2.DescribeSchedulerxNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_schedulerx_namespace DescribeSchedulerxNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["Name"] != nil {
		d.Set("namespace_name", objectRaw["Name"])
	}
	if objectRaw["UId"] != nil {
		d.Set("namespace_uid", objectRaw["UId"])
	}

	d.Set("namespace_uid", d.Id())

	return nil
}

func resourceAliCloudSchedulerxNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Namespace. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

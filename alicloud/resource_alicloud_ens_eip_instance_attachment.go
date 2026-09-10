// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEnsEipInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsEipInstanceAttachmentCreate,
		Read:   resourceAliCloudEnsEipInstanceAttachmentRead,
		Delete: resourceAliCloudEnsEipInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Nat", "SlbInstance", "NetworkInterface", "EnsInstance"}, false),
			},
			"standby": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEnsEipInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateEnsEipAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AllocationId"] = d.Get("allocation_id")
	query["InstanceId"] = d.Get("instance_id")
	query["InstanceType"] = d.Get("instance_type")

	if v, ok := d.GetOkExists("standby"); ok {
		request["Standby"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_eip_instance_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", query["AllocationId"], query["InstanceId"], query["InstanceType"]))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsEipInstanceAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsEipInstanceAttachmentRead(d, meta)
}

func resourceAliCloudEnsEipInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsEipInstanceAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_eip_instance_attachment DescribeEnsEipInstanceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Standby"] != nil {
		d.Set("standby", objectRaw["Standby"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["AllocationId"] != nil {
		d.Set("allocation_id", objectRaw["AllocationId"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}
	if objectRaw["InstanceType"] != nil {
		d.Set("instance_type", objectRaw["InstanceType"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("allocation_id", parts[0])
	d.Set("instance_id", parts[1])
	d.Set("instance_type", parts[2])

	return nil
}

func resourceAliCloudEnsEipInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UnAssociateEnsEipAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AllocationId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ensServiceV2.EnsEipInstanceAttachmentStateRefreshFunc(d.Id(), "InstanceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

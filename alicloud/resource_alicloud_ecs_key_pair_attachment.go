package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsKeyPairAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsKeyPairAttachmentCreate,
		Read:   resourceAliCloudEcsKeyPairAttachmentRead,
		Delete: resourceAliCloudEcsKeyPairAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_pair_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  StringLenBetween(2, 128),
				ConflictsWith: []string{"key_name"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  StringLenBetween(2, 128),
				ConflictsWith: []string{"key_pair_name"},
				Deprecated:    "Field `key_name` has been deprecated from provider version 1.121.0. New field `key_pair_name` instead.",
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEcsKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AttachKeyPair"
	request := make(map[string]interface{})
	var err error
	instanceIds := convertToInterfaceArray(d.Get("instance_ids"))

	request["RegionId"] = client.RegionId
	request["InstanceIds"] = convertListToJsonString(instanceIds)

	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	} else if v, ok := d.GetOk("key_name"); ok {
		request["KeyPairName"] = v
	} else {
		return WrapError(Error(`[ERROR] Field "key_pair_name" or "key_name" must be set one!`))
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_key_pair_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["KeyPairName"], request["InstanceIds"]))

	if d.Get("force").(bool) {
		for _, instanceId := range instanceIds {
			err := ecsService.RebootEcsInstance(fmt.Sprint(instanceId))
			if err != nil {
				return WrapError(err)
			}
		}

		for _, instanceId := range instanceIds {
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.InstanceStateRefreshFunc(fmt.Sprint(instanceId), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}

	return resourceAliCloudEcsKeyPairAttachmentRead(d, meta)
}

func resourceAliCloudEcsKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsKeyPairAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_key_pair_attachment DescribeEcsKeyPairAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	instanceIds, err := convertJsonStringToList(parts[1])
	if err != nil {
		return WrapError(err)
	}

	d.Set("instance_ids", instanceIds)
	d.Set("key_pair_name", object["KeyPairName"])
	d.Set("key_name", object["KeyPairName"])

	return nil
}

func resourceAliCloudEcsKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachKeyPair"
	var response map[string]interface{}
	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"KeyPairName": parts[0],
		"InstanceIds": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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

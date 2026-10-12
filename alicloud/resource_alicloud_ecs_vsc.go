// Package alicloud. Handwritten resource for ECS Vsc (Virtual Storage Channel).
// ACube generation was unavailable for this resource type (see governance task
// 86548773); this file is authored manually against the pre-PASS CloudSpec model
// and the Ecs 2014-05-26 OpenAPI (CreateVsc/DescribeVscs/DeleteVsc/ModifyVsc/
// TagResources/UntagResources/JoinResourceGroup).
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

func resourceAliCloudEcsVsc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsVscCreate,
		Read:   resourceAliCloudEcsVscRead,
		Update: resourceAliCloudEcsVscUpdate,
		Delete: resourceAliCloudEcsVscDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vsc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vsc_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vsc_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Primary", "Standard"}, false),
			},
		},
	}
}

func resourceAliCloudEcsVscCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateVsc"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("vsc_name"); ok {
		request["VscName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("vsc_type"); ok {
		request["VscType"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_vsc", action, AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["VscId"]; !ok || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_vsc")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAliCloudEcsVscRead(d, meta)
}

func resourceAliCloudEcsVscRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsVsc(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_vsc DescribeEcsVsc Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vsc_id", objectRaw["VscId"])
	d.Set("vsc_name", objectRaw["VscName"])
	d.Set("vsc_type", objectRaw["VscType"])
	d.Set("region_id", client.RegionId)

	if v, ok := objectRaw["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	return nil
}

func resourceAliCloudEcsVscUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	// Modify updatable attributes (VscName / Description) via ModifyVsc.
	if d.HasChange("vsc_name") || d.HasChange("description") {
		action := "ModifyVsc"
		request := make(map[string]interface{})
		query := make(map[string]interface{})
		request["VscId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)
		if v, ok := d.GetOk("vsc_name"); ok {
			request["VscName"] = v
		}
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, e := client.RpcPost("Ecs", "2014-05-26", action, query, request, false)
			if e != nil {
				if NeedRetry(e) {
					wait()
					return resource.RetryableError(e)
				}
				return resource.NonRetryableError(e)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("vsc_name")
		d.SetPartial("description")
	}

	// Move the resource to another resource group via JoinResourceGroup.
	if d.HasChange("resource_group_id") {
		action := "JoinResourceGroup"
		request := make(map[string]interface{})
		query := make(map[string]interface{})
		request["ResourceId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["ResourceGroupId"] = d.Get("resource_group_id")
		request["ResourceType"] = "vsc"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, e := client.RpcPost("Ecs", "2014-05-26", action, query, request, false)
			if e != nil {
				if NeedRetry(e) {
					wait()
					return resource.RetryableError(e)
				}
				return resource.NonRetryableError(e)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	ecsServiceV2 := EcsServiceV2{client}
	if d.HasChange("tags") {
		if err := ecsServiceV2.SetResourceTags(d, "vsc"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	d.Partial(false)
	return resourceAliCloudEcsVscRead(d, meta)
}

func resourceAliCloudEcsVscDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVsc"
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["VscId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, e := client.RpcPost("Ecs", "2014-05-26", action, query, request, false)
		if e != nil {
			if NeedRetry(e) {
				wait()
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidVscId.NotFound", "NoSuchVsc"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ecsServiceV2 := EcsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsServiceV2.EcsVscStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

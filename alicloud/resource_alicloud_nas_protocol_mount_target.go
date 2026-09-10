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

func resourceAliCloudNasProtocolMountTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNasProtocolMountTargetCreate,
		Read:   resourceAliCloudNasProtocolMountTargetRead,
		Update: resourceAliCloudNasProtocolMountTargetUpdate,
		Delete: resourceAliCloudNasProtocolMountTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"export_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fset_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protocol_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudNasProtocolMountTargetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateProtocolMountTarget"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("file_system_id"); ok {
		request["FileSystemId"] = v
	}
	if v, ok := d.GetOk("protocol_service_id"); ok {
		request["ProtocolServiceId"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("fset_id"); ok {
		request["FsetId"] = v
	}
	if v, ok := d.GetOk("vswitch_ids"); ok {
		vSwitchIdsMapsArray := convertToInterfaceArray(v)

		request["VSwitchIds"] = vSwitchIdsMapsArray
	}

	if v, ok := d.GetOk("path"); ok {
		request["Path"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("access_group_name"); ok {
		request["AccessGroupName"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_protocol_mount_target", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["FileSystemId"], request["ProtocolServiceId"], response["ExportId"]))

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"AVAILABLE"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, nasServiceV2.NasProtocolMountTargetStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNasProtocolMountTargetRead(d, meta)
}

func resourceAliCloudNasProtocolMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasServiceV2 := NasServiceV2{client}

	objectRaw, err := nasServiceV2.DescribeNasProtocolMountTarget(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_protocol_mount_target DescribeNasProtocolMountTarget Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_group_name", objectRaw["AccessGroupName"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("fset_id", objectRaw["FsetId"])
	d.Set("path", objectRaw["Path"])
	d.Set("status", objectRaw["Status"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("export_id", objectRaw["ExportId"])

	vSwitchIdsRaw := make([]interface{}, 0)
	if objectRaw["VSwitchIds"] != nil {
		vSwitchIdsRaw = convertToInterfaceArray(objectRaw["VSwitchIds"])
	}

	d.Set("vswitch_ids", vSwitchIdsRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("file_system_id", parts[0])
	d.Set("protocol_service_id", parts[1])

	return nil
}

func resourceAliCloudNasProtocolMountTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyProtocolMountTarget"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ExportId"] = parts[2]
	request["FileSystemId"] = parts[0]
	request["ProtocolServiceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
		nasServiceV2 := NasServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"AVAILABLE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, nasServiceV2.NasProtocolMountTargetStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudNasProtocolMountTargetRead(d, meta)
}

func resourceAliCloudNasProtocolMountTargetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteProtocolMountTarget"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ExportId"] = parts[2]
	request["FileSystemId"] = parts[0]
	request["ProtocolServiceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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

	nasServiceV2 := NasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 10*time.Second, nasServiceV2.NasProtocolMountTargetStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

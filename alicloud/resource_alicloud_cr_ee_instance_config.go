package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCrEEInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrEEInstanceConfigCreate,
		Read:   resourceAlicloudCrEEInstanceConfigRead,
		Update: resourceAlicloudCrEEInstanceConfigUpdate,
		Delete: resourceAlicloudCrEEInstanceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replication_acceleration": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCrEEInstanceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return updateInstanceConfig(d, meta, d.Get("instance_id").(string), d.Get("replication_acceleration").(bool), true)
}

func resourceAlicloudCrEEInstanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	request := &GetInstanceConfigRequest{
		InstanceID: d.Id(),
	}
	response, err := crService.GetInstanceConfig(request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetInstanceConfig", AlibabaCloudSdkGoERROR)
	}
	addDebug("GetInstanceConfig", response, request)
	if !response.Data.IsSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), "GetInstanceConfig", AlibabaCloudSdkGoERROR)
	}
	d.Set("instance_id", response.Data.InstanceId)
	d.Set("replication_acceleration", response.Data.SyncTransAccelerate)

	return nil
}

func resourceAlicloudCrEEInstanceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChanges("replication_acceleration") {
		return nil
	}

	return updateInstanceConfig(d, meta, d.Id(), d.Get("replication_acceleration").(bool), false)
}

func resourceAlicloudCrEEInstanceConfigDelete(_ *schema.ResourceData, _ interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudCrEEInstanceConfig. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func updateInstanceConfig(d *schema.ResourceData, meta interface{}, instanceID string, syncTransAccelerate bool, create bool) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	request := &UpdateInstanceConfigRequest{
		InstanceID:          instanceID,
		SyncTransAccelerate: syncTransAccelerate,
	}
	response, err := crService.UpdateInstanceConfig(request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceID, "UpdateInstanceConfig", AlibabaCloudSdkGoERROR)
	}
	addDebug("UpdateInstanceConfig", response, request)
	if create {
		d.SetId(instanceID)
	}

	return resourceAlicloudCrEEInstanceConfigRead(d, meta)
}

package config

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceConfigurationRecorder() *schema.Resource {
	return &schema.Resource{
		Create: ResourceConfigurationRecorderCreate,
		Read:   ResourceConfigurationRecorderRead,
		Update: ResourceConfigurationRecorderUpdate,
		Delete: ResourceConfigurationRecorderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enterprise_edition": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"organization_enable_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_master_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceConfigurationRecorderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	action := "StartConfigurationRecorder"
	request := make(map[string]interface{})

	if v, ok := d.GetOkExists("enterprise_edition"); ok {
		request["EnterpriseEdition"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2019-01-08", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_config_configuration_recorder", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	response = response["ConfigurationRecorder"].(map[string]interface{})
	accountId, err := convert.FormatInt(response["AccountId"])
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(accountId))

	return ResourceConfigurationRecorderUpdate(d, meta)
}
func ResourceConfigurationRecorderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := NewConfigService(client)
	object, err := configService.DescribeConfigConfigurationRecorder(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_configuration_recorder configService.DescribeConfigConfigurationRecorder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("organization_enable_status", object["OrganizationEnableStatus"])
	organization_master_idVal, err := convert.FormatInt(object["OrganizationMasterId"])
	if err != nil {
		return err
	}
	d.Set("organization_master_id", organization_master_idVal)
	d.Set("resource_types", object["ResourceTypes"])
	d.Set("status", object["ConfigurationRecorderStatus"])
	return nil
}
func ResourceConfigurationRecorderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := NewConfigService(client)
	var response map[string]interface{}
	var err error
	update := false
	request := make(map[string]interface{})
	if !d.IsNewResource() && d.HasChange("resource_types") {
		update = true
	}
	request["ResourceTypes"] = convert.ConvertListToCommaSeparate(d.Get("resource_types").(*schema.Set).List())
	if update {
		action := "PutConfigurationRecorder"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Config", "2019-01-08", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"REGISTERED"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, configService.ConfigConfigurationRecorderStateRefreshFunc(d.Id(), []string{"REGISTRABLE"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	return ResourceConfigurationRecorderRead(d, meta)
}
func ResourceConfigurationRecorderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy ResourceConfigurationRecorder. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

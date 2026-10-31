package sls

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudLogMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogMachineGroupCreate,
		Read:   resourceAlicloudLogMachineGroupRead,
		Update: resourceAlicloudLogMachineGroupUpdate,
		Delete: resourceAlicloudLogMachineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      sls.MachineIDTypeIP,
				ValidateFunc: validation.StringInSlice([]string{sls.MachineIDTypeIP, sls.MachineIDTypeUserDefined}, false),
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identify_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MinItems: 1,
			},
		},
	}
}

func resourceAlicloudLogMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	params := &sls.MachineGroup{
		Name:          d.Get("name").(string),
		MachineIDType: d.Get("identify_type").(string),
		MachineIDList: convert.ExpandStringList(d.Get("identify_list").(*schema.Set).List()),
		Attribute: sls.MachinGroupAttribute{
			TopicName: d.Get("topic").(string),
		},
	}
	var requestInfo *sls.Client
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.CreateMachineGroup(d.Get("project").(string), params)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("CreateMachineGroup", raw, requestInfo, map[string]interface{}{
				"project":      d.Get("project").(string),
				"MachineGroup": params,
			})
		}
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_machine_group", "CreateMachineGroup", errs.AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := logService.DescribeLogMachineGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("identify_type", object.MachineIDType)
	d.Set("identify_list", object.MachineIDList)
	d.Set("topic", object.Attribute.TopicName)

	return nil
}

func resourceAlicloudLogMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("identify_type") || d.HasChange("identify_list") || d.HasChange("topic") {
		parts, err := helper.ParseResourceId(d.Id(), 2)
		if err != nil {
			return errs.WrapError(err)
		}

		client := meta.(*connectivity.AliyunClient)
		var requestInfo *sls.Client
		params := &sls.MachineGroup{
			Name:          parts[1],
			MachineIDType: d.Get("identify_type").(string),
			MachineIDList: convert.ExpandStringList(d.Get("identify_list").(*schema.Set).List()),
			Attribute: sls.MachinGroupAttribute{
				TopicName: d.Get("topic").(string),
			},
		}
		if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.UpdateMachineGroup(parts[0], params)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
					time.Sleep(5 * time.Second)
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			if logs.DebugOn() {
				logs.AddDebug("UpdateMachineGroup", raw, requestInfo, map[string]interface{}{
					"project":      parts[0],
					"MachineGroup": params,
				})
			}
			return nil
		}); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UpdateMachineGroup", errs.AliyunLogGoSdkERROR)
		}
	}

	return resourceAlicloudLogMachineGroupRead(d, meta)
}

func resourceAlicloudLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	var requestInfo *sls.Client
	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteMachineGroup(parts[0], parts[1])
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("DeleteMachineGroup", raw, requestInfo, map[string]interface{}{
				"project":      parts[0],
				"machineGroup": parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_store", "ListShards", errs.AliyunLogGoSdkERROR)
	}
	return errs.WrapError(logService.WaitForLogMachineGroup(d.Id(), Deleted, DefaultTimeout))
}

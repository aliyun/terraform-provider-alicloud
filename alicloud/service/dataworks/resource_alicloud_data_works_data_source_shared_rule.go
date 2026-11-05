// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package dataworks

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDataSourceSharedRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataSourceSharedRuleCreate,
		Read:   resourceDataSourceSharedRuleRead,
		Delete: resourceDataSourceSharedRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_source_shared_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"env_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"shared_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDataSourceSharedRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDataSourceSharedRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DataSourceId"] = d.Get("data_source_id")
	request["RegionId"] = client.RegionId

	request["TargetProjectId"] = d.Get("target_project_id")
	if v, ok := d.GetOk("shared_user"); ok {
		request["SharedUser"] = v
	}
	request["EnvType"] = d.Get("env_type")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_data_works_data_source_shared_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DataSourceId"], response["Id"]))

	return resourceDataSourceSharedRuleRead(d, meta)
}

func resourceDataSourceSharedRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := NewService(client)

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDataSourceSharedRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_data_source_shared_rule DescribeDataWorksDataSourceSharedRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["EnvType"] != nil {
		d.Set("env_type", objectRaw["EnvType"])
	}
	if objectRaw["SharedUser"] != nil {
		d.Set("shared_user", objectRaw["SharedUser"])
	}
	if objectRaw["TargetProjectId"] != nil {
		d.Set("target_project_id", objectRaw["TargetProjectId"])
	}
	if objectRaw["DataSourceId"] != nil {
		d.Set("data_source_id", objectRaw["DataSourceId"])
	}
	if objectRaw["Id"] != nil {
		d.Set("data_source_shared_rule_id", objectRaw["Id"])
	}

	return nil
}

func resourceDataSourceSharedRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDataSourceSharedRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"1000011038"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

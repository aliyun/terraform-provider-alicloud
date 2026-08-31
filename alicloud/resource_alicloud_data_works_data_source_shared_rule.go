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

func resourceAliCloudDataWorksDataSourceSharedRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksDataSourceSharedRuleCreate,
		Read:   resourceAliCloudDataWorksDataSourceSharedRuleRead,
		Delete: resourceAliCloudDataWorksDataSourceSharedRuleDelete,
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

func resourceAliCloudDataWorksDataSourceSharedRuleCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_data_source_shared_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DataSourceId"], response["Id"]))

	return resourceAliCloudDataWorksDataSourceSharedRuleRead(d, meta)
}

func resourceAliCloudDataWorksDataSourceSharedRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDataSourceSharedRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_data_source_shared_rule DescribeDataWorksDataSourceSharedRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudDataWorksDataSourceSharedRuleDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"1000011038"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

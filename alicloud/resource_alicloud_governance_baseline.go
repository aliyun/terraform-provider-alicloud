// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGovernanceBaseline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGovernanceBaselineCreate,
		Read:   resourceAliCloudGovernanceBaselineRead,
		Update: resourceAliCloudGovernanceBaselineUpdate,
		Delete: resourceAliCloudGovernanceBaselineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"baseline_items": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"1.0"}, false),
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"baseline_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudGovernanceBaselineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccountFactoryBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("baseline_items"); ok {
		baselineItemsMaps := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Config"] = dataLoopTmp["config"]
			dataLoopMap["Name"] = dataLoopTmp["name"]
			dataLoopMap["Version"] = dataLoopTmp["version"]
			baselineItemsMaps = append(baselineItemsMaps, dataLoopMap)
		}
		request["BaselineItems"] = baselineItemsMaps
	}

	if v, ok := d.GetOk("baseline_name"); ok {
		request["BaselineName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_governance_baseline", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BaselineId"]))

	return resourceAliCloudGovernanceBaselineRead(d, meta)
}

func resourceAliCloudGovernanceBaselineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	governanceServiceV2 := GovernanceServiceV2{client}

	objectRaw, err := governanceServiceV2.DescribeGovernanceBaseline(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_governance_baseline DescribeGovernanceBaseline Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["BaselineName"] != nil {
		d.Set("baseline_name", objectRaw["BaselineName"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}

	baselineItems1Raw := objectRaw["BaselineItems"]
	baselineItemsMaps := make([]map[string]interface{}, 0)
	if baselineItems1Raw != nil {
		for _, baselineItemsChild1Raw := range baselineItems1Raw.([]interface{}) {
			baselineItemsMap := make(map[string]interface{})
			baselineItemsChild1Raw := baselineItemsChild1Raw.(map[string]interface{})
			baselineItemsMap["config"] = baselineItemsChild1Raw["Config"]
			baselineItemsMap["name"] = baselineItemsChild1Raw["Name"]
			baselineItemsMap["version"] = baselineItemsChild1Raw["Version"]

			baselineItemsMaps = append(baselineItemsMaps, baselineItemsMap)
		}
	}
	if objectRaw["BaselineItems"] != nil {
		if err := d.Set("baseline_items", baselineItemsMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudGovernanceBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateAccountFactoryBaseline"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["BaselineId"] = d.Id()
	query["RegionId"] = client.RegionId
	if d.HasChange("baseline_items") {
		update = true
		if v, ok := d.GetOk("baseline_items"); ok {
			baselineItemsMaps := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Config"] = dataLoopTmp["config"]
				dataLoopMap["Name"] = dataLoopTmp["name"]
				dataLoopMap["Version"] = dataLoopTmp["version"]
				baselineItemsMaps = append(baselineItemsMaps, dataLoopMap)
			}
			request["BaselineItems"] = baselineItemsMaps
		}
	}

	if d.HasChange("baseline_name") {
		update = true
		request["BaselineName"] = d.Get("baseline_name")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)
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
	}

	return resourceAliCloudGovernanceBaselineRead(d, meta)
}

func resourceAliCloudGovernanceBaselineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccountFactoryBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["BaselineId"] = d.Id()
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)

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
		if IsExpectedErrors(err, []string{"InvalidBaseline.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

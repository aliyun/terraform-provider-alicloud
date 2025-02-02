// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcFlowLogCreate,
		Read:   resourceAlicloudVpcFlowLogRead,
		Update: resourceAlicloudVpcFlowLogUpdate,
		Delete: resourceAlicloudVpcFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregation_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_log_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"NetworkInterface", "VPC", "VSwitch"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Inactive"}, false),
			},
			"tags": tagsSchema(),
			"traffic_path": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"traffic_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"All", "Allow", "Drop"}, false),
			},
		},
	}
}

func resourceAlicloudVpcFlowLogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateFlowLog"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["ProjectName"] = d.Get("project_name")
	request["ResourceId"] = d.Get("resource_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["LogStoreName"] = d.Get("log_store_name")
	request["TrafficType"] = d.Get("traffic_type")
	if v, ok := d.GetOk("flow_log_name"); ok {
		request["FlowLogName"] = v
	}
	if v, ok := d.GetOk("aggregation_interval"); ok {
		request["AggregationInterval"] = v
	}
	if v, ok := d.GetOk("traffic_path"); ok {
		trafficPathMaps := v.([]interface{})
		request["TrafficPath"] = trafficPathMaps
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ResourceType"] = d.Get("resource_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationFailed.LastTokenProcessing", "LastTokenProcessing", "IncorrectStatus", "InvalidHdMonitorStatus", "OperationConflict", "ServiceUnavailable", "SystemBusy", "UnknownError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_flow_log", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FlowLogId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcFlowLogUpdate(d, meta)
}

func resourceAlicloudVpcFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcFlowLog(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_flow_log DescribeVpcFlowLog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("aggregation_interval", objectRaw["AggregationInterval"])
	d.Set("business_status", objectRaw["BusinessStatus"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("flow_log_name", objectRaw["FlowLogName"])
	d.Set("log_store_name", objectRaw["LogStoreName"])
	d.Set("project_name", objectRaw["ProjectName"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("resource_id", objectRaw["ResourceId"])
	d.Set("resource_type", objectRaw["ResourceType"])
	d.Set("status", objectRaw["Status"])
	d.Set("traffic_type", objectRaw["TrafficType"])
	d.Set("flow_log_id", objectRaw["FlowLogId"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	trafficPathList1Raw, _ := jsonpath.Get("$.TrafficPath.TrafficPathList", objectRaw)
	d.Set("traffic_path", trafficPathList1Raw)

	return nil
}

func resourceAlicloudVpcFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyFlowLogAttribute"
	var err error
	request = make(map[string]interface{})

	request["FlowLogId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("flow_log_name") {
		update = true
		request["FlowLogName"] = d.Get("flow_log_name")
	}

	if !d.IsNewResource() && d.HasChange("aggregation_interval") {
		update = true
		request["AggregationInterval"] = d.Get("aggregation_interval")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active", " Inactive"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("flow_log_name")
		d.SetPartial("aggregation_interval")
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "FLOWLOG"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

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
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		vpcServiceV2 := VpcServiceV2{client}
		object, err := vpcServiceV2.DescribeVpcFlowLog(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Active" {
				action = "ActiveFlowLog"
				request = make(map[string]interface{})
				request["FlowLogId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

					if err != nil {
						if IsExpectedErrors(err, []string{"TaskConflict", "LastTokenProcessing", "IncorrectStatus", "IncorrectStatus.flowlog", "InvalidStatus", "OperationConflict", "ServiceUnavailable", "SystemBusy", "UnknownError"}) || NeedRetry(err) {
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Inactive" {
				action = "DeactiveFlowLog"
				request = make(map[string]interface{})

				request["FlowLogId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

					if err != nil {
						if IsExpectedErrors(err, []string{"TaskConflict", "LastTokenProcessing", "IncorrectStatus", "IncorrectStatus.flowlog", "InvalidStatus", "OperationConflict", "ServiceUnavailable", "SystemBusy", "UnknownError"}) || NeedRetry(err) {
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Inactive"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "FLOWLOG"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudVpcFlowLogRead(d, meta)
}

func resourceAlicloudVpcFlowLogDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteFlowLog"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["FlowLogId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.LastTokenProcessing", "LastTokenProcessing", "InvalidHdMonitorStatus", "IncorrectStatus", "IncorrectStatus.flowlog", "InvalidStatus", "OperationConflict", "SystemBusy", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.IsNotAvailable", "Instance.IsNotPostPay"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcFlowLogStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

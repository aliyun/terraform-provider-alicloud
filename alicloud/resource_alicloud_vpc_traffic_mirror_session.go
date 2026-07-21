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

func resourceAliCloudVpcTrafficMirrorSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcTrafficMirrorSessionCreate,
		Read:   resourceAliCloudVpcTrafficMirrorSessionRead,
		Update: resourceAliCloudVpcTrafficMirrorSessionUpdate,
		Delete: resourceAliCloudVpcTrafficMirrorSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"packet_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(1, 32766),
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
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_session_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"traffic_mirror_session_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 128),
			},
			"traffic_mirror_source_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"traffic_mirror_target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"NetworkInterface", "SLB"}, false),
			},
			"virtual_network_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 16777215),
			},
		},
	}
}

func resourceAliCloudVpcTrafficMirrorSessionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrafficMirrorSession"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("traffic_mirror_session_description"); ok {
		request["TrafficMirrorSessionDescription"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_session_name"); ok {
		request["TrafficMirrorSessionName"] = v
	}
	request["TrafficMirrorTargetId"] = d.Get("traffic_mirror_target_id")
	request["TrafficMirrorTargetType"] = d.Get("traffic_mirror_target_type")
	request["TrafficMirrorFilterId"] = d.Get("traffic_mirror_filter_id")
	if v, ok := d.GetOk("virtual_network_id"); ok {
		request["VirtualNetworkId"] = v
	}
	request["Priority"] = d.Get("priority")
	if v, ok := d.GetOkExists("enabled"); ok {
		request["Enabled"] = v
	}
	if v, ok := d.GetOk("packet_length"); ok {
		request["PacketLength"] = v
	}
	if v, ok := d.GetOk("traffic_mirror_source_ids"); ok {
		trafficMirrorSourceIdsMaps := v.([]interface{})
		request["TrafficMirrorSourceIds"] = trafficMirrorSourceIdsMaps
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_traffic_mirror_session", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TrafficMirrorSessionId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), "TrafficMirrorSessionStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcTrafficMirrorSessionUpdate(d, meta)
}

func resourceAliCloudVpcTrafficMirrorSessionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcTrafficMirrorSession(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_traffic_mirror_session DescribeVpcTrafficMirrorSession Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enabled", objectRaw["Enabled"])
	d.Set("packet_length", objectRaw["PacketLength"])
	d.Set("priority", objectRaw["Priority"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["TrafficMirrorSessionStatus"])
	d.Set("traffic_mirror_filter_id", objectRaw["TrafficMirrorFilterId"])
	d.Set("traffic_mirror_session_description", objectRaw["TrafficMirrorSessionDescription"])
	d.Set("traffic_mirror_session_name", objectRaw["TrafficMirrorSessionName"])
	d.Set("traffic_mirror_target_id", objectRaw["TrafficMirrorTargetId"])
	d.Set("traffic_mirror_target_type", objectRaw["TrafficMirrorTargetType"])
	d.Set("virtual_network_id", objectRaw["VirtualNetworkId"])
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	trafficMirrorSourceIds1Raw := make([]interface{}, 0)
	if objectRaw["TrafficMirrorSourceIds"] != nil {
		trafficMirrorSourceIds1Raw = objectRaw["TrafficMirrorSourceIds"].([]interface{})
	}

	d.Set("traffic_mirror_source_ids", trafficMirrorSourceIds1Raw)

	return nil
}

func resourceAliCloudVpcTrafficMirrorSessionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateTrafficMirrorSessionAttribute"
	var err error
	request = make(map[string]interface{})

	request["TrafficMirrorSessionId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("traffic_mirror_session_description") {
		update = true
		request["TrafficMirrorSessionDescription"] = d.Get("traffic_mirror_session_description")
	}

	if !d.IsNewResource() && d.HasChange("traffic_mirror_session_name") {
		update = true
		request["TrafficMirrorSessionName"] = d.Get("traffic_mirror_session_name")
	}

	if !d.IsNewResource() && d.HasChange("traffic_mirror_target_id") {
		update = true
	}
	request["TrafficMirrorTargetId"] = d.Get("traffic_mirror_target_id")
	if !d.IsNewResource() && d.HasChange("traffic_mirror_target_type") {
		update = true
	}
	request["TrafficMirrorTargetType"] = d.Get("traffic_mirror_target_type")
	if !d.IsNewResource() && d.HasChange("traffic_mirror_filter_id") {
		update = true
	}
	request["TrafficMirrorFilterId"] = d.Get("traffic_mirror_filter_id")
	if !d.IsNewResource() && d.HasChange("virtual_network_id") {
		update = true
		request["VirtualNetworkId"] = d.Get("virtual_network_id")
	}

	if !d.IsNewResource() && d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")
	if !d.IsNewResource() && d.HasChange("enabled") {
		update = true
		request["Enabled"] = d.Get("enabled")
	}

	if d.HasChange("traffic_mirror_target_id") || d.HasChange("traffic_mirror_target_type") {
		update = true
		request["TrafficMirrorTargetId"] = d.Get("traffic_mirror_target_id")
		request["TrafficMirrorTargetType"] = d.Get("traffic_mirror_target_type")
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), "TrafficMirrorSessionStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("traffic_mirror_session_description")
		d.SetPartial("traffic_mirror_session_name")
		d.SetPartial("traffic_mirror_target_id")
		d.SetPartial("traffic_mirror_target_type")
		d.SetPartial("traffic_mirror_filter_id")
		d.SetPartial("virtual_network_id")
		d.SetPartial("priority")
		d.SetPartial("enabled")
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

	request["ResourceType"] = "TrafficMirrorSession"

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

	update = false
	if !d.IsNewResource() && d.HasChange("traffic_mirror_source_ids") {
		update = true
		oldEntry, newEntry := d.GetChange("traffic_mirror_source_ids")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			action = "RemoveSourcesFromTrafficMirrorSession"
			request = make(map[string]interface{})

			request["TrafficMirrorSessionId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed
			trafficMirrorSourceIdsMaps := localData.([]interface{})
			request["TrafficMirrorSourceIds"] = trafficMirrorSourceIdsMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable"}) || NeedRetry(err) {
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
			stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), "TrafficMirrorSessionStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if len(added.([]interface{})) > 0 {
			action = "AddSourcesToTrafficMirrorSession"
			request = make(map[string]interface{})

			request["TrafficMirrorSessionId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added
			trafficMirrorSourceIdsMaps := localData.([]interface{})
			request["TrafficMirrorSourceIds"] = trafficMirrorSourceIdsMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable"}) || NeedRetry(err) {
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
			stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), "TrafficMirrorSessionStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
	}
	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "TrafficMirrorSession"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudVpcTrafficMirrorSessionRead(d, meta)
}

func resourceAliCloudVpcTrafficMirrorSessionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteTrafficMirrorSession"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["TrafficMirrorSessionId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TrafficMirrorSession", "OperationConflict", "IncorrectStatus", "SystemBusy", "LastTokenProcessing", "OperationFailed.LastTokenProcessing", "ServiceUnavailable", "IncorrectStatus.TrafficMirrorFilter"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.TrafficMirrorSession"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcTrafficMirrorSessionStateRefreshFunc(d.Id(), "TrafficMirrorSessionStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

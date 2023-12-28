// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcPeerConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcPeerConnectionCreate,
		Read:   resourceAliCloudVpcPeerConnectionRead,
		Update: resourceAliCloudVpcPeerConnectionUpdate,
		Delete: resourceAliCloudVpcPeerConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accepting_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accepting_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accepting_ali_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntAtLeast(0),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peer_connection_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudVpcPeerConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}
	var response map[string]interface{}
	action := "CreateVpcPeerConnection"
	request := make(map[string]interface{})
	conn, err := client.NewVpcpeerClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateVpcPeerConnection")
	request["VpcId"] = d.Get("vpc_id")
	request["AcceptingVpcId"] = d.Get("accepting_vpc_id")
	request["AcceptingRegionId"] = d.Get("accepting_region_id")

	if v, ok := d.GetOkExists("accepting_ali_uid"); ok {
		request["AcceptingAliUid"] = v
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("peer_connection_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_peer_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	stateConf := BuildStateConf([]string{}, []string{"Activated", "Accepting"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcPeerConnectionUpdate(d, meta)
}

func resourceAliCloudVpcPeerConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	object, err := vpcServiceV2.DescribeVpcPeerConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_peer_connection DescribeVpcPeerConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accepting_region_id", object["AcceptingRegionId"])
	d.Set("accepting_ali_uid", object["AcceptingOwnerUid"])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("peer_connection_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])
	d.Set("create_time", object["GmtCreate"])

	if vpc, ok := object["Vpc"]; ok {
		vpcArg := vpc.(map[string]interface{})

		if vpcId, ok := vpcArg["VpcId"]; ok {
			d.Set("vpc_id", vpcId)
		}
	}

	if acceptingVpc, ok := object["AcceptingVpc"]; ok {
		acceptingVpcArg := acceptingVpc.(map[string]interface{})

		if vpcId, ok := acceptingVpcArg["VpcId"]; ok {
			d.Set("accepting_vpc_id", vpcId)
		}
	}

	if tags, ok := object["Tags"]; ok {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliCloudVpcPeerConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := vpcServiceV2.SetVpcPeerResourceTags(d, "PeerConnection"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	update := false
	modifyVpcPeerConnectionReq := map[string]interface{}{
		"ClientToken": buildClientToken("ModifyVpcPeerConnection"),
		"InstanceId":  d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true

		if v, ok := d.GetOkExists("bandwidth"); ok {
			modifyVpcPeerConnectionReq["Bandwidth"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("peer_connection_name") {
		update = true
	}
	if v, ok := d.GetOk("peer_connection_name"); ok {
		modifyVpcPeerConnectionReq["Name"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		modifyVpcPeerConnectionReq["Description"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		modifyVpcPeerConnectionReq["DryRun"] = v
	}

	if update {
		action := "ModifyVpcPeerConnection"
		conn, err := client.NewVpcpeerClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, modifyVpcPeerConnectionReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyVpcPeerConnectionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Accepting", "Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("bandwidth")
		d.SetPartial("peer_connection_name")
		d.SetPartial("description")
	}

	update = false
	moveResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": "PeerConnection",
		"ResourceId":   d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		moveResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewVpcpeerClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, moveResourceGroupReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, moveResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	if d.HasChange("status") {
		conn, err := client.NewVpcpeerClient()
		if err != nil {
			return WrapError(err)
		}

		object, err := vpcServiceV2.DescribeVpcPeerConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Rejected" {
				action := "RejectVpcPeerConnection"

				request := map[string]interface{}{
					"ClientToken": buildClientToken("RejectVpcPeerConnection"),
					"InstanceId":  d.Id(),
				}

				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
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

				stateConf := BuildStateConf([]string{}, []string{"Rejected"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if target == "Activated" {
				action := "AcceptVpcPeerConnection"

				request := map[string]interface{}{
					"ClientToken": buildClientToken("AcceptVpcPeerConnection"),
					"InstanceId":  d.Id(),
				}

				if v, ok := d.GetOk("resource_group_id"); ok {
					request["ResourceGroupId"] = v
				}

				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
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

				stateConf := BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}

	d.Partial(false)

	return resourceAliCloudVpcPeerConnectionRead(d, meta)
}

func resourceAliCloudVpcPeerConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}
	action := "DeleteVpcPeerConnection"
	var response map[string]interface{}

	conn, err := client.NewVpcpeerClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken": buildClientToken("DeleteVpcPeerConnection"),
		"InstanceId":  d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.InstanceId"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

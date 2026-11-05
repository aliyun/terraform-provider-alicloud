// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package vpcpeer

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePeerConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourcePeerConnectionAccepterCreate,
		Read:   resourcePeerConnectionAccepterRead,
		Update: resourcePeerConnectionAccepterUpdate,
		Delete: resourcePeerConnectionAccepterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accepting_owner_uid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"accepting_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accepting_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntAtLeast(1),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"link_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Platinum", "Gold", "Silver"}, false),
			},
			"peer_connection_accepter_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePeerConnectionAccepterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AcceptVpcPeerConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, true)
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
		if !errs.IsExpectedErrors(err, []string{"IncorrectStatus.VpcPeer"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_peer_connection_accepter", action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	vpcPeerServiceV2 := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionAccepterStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourcePeerConnectionAccepterUpdate(d, meta)
}

func resourcePeerConnectionAccepterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcPeerServiceV2 := Service{client}

	objectRaw, err := vpcPeerServiceV2.DescribeVpcPeerPeerConnectionAccepter(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_peer_connection_accepter DescribeVpcPeerPeerConnectionAccepter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["AcceptingOwnerUid"] != nil {
		d.Set("accepting_owner_uid", objectRaw["AcceptingOwnerUid"])
	}
	if objectRaw["AcceptingRegionId"] != nil {
		d.Set("accepting_region_id", objectRaw["AcceptingRegionId"])
	}
	if objectRaw["Bandwidth"] != nil {
		d.Set("bandwidth", objectRaw["Bandwidth"])
	}
	if objectRaw["GmtCreate"] != nil {
		d.Set("create_time", objectRaw["GmtCreate"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["LinkType"] != nil {
		d.Set("link_type", objectRaw["LinkType"])
	}
	if objectRaw["Name"] != nil {
		d.Set("peer_connection_accepter_name", objectRaw["Name"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}

	acceptingVpc1RawObj, _ := jsonpath.Get("$.AcceptingVpc", objectRaw)
	acceptingVpc1Raw := make(map[string]interface{})
	if acceptingVpc1RawObj != nil {
		acceptingVpc1Raw = acceptingVpc1RawObj.(map[string]interface{})
	}
	if acceptingVpc1Raw["VpcId"] != nil {
		d.Set("accepting_vpc_id", acceptingVpc1Raw["VpcId"])
	}

	vpc1RawObj, _ := jsonpath.Get("$.Vpc", objectRaw)
	vpc1Raw := make(map[string]interface{})
	if vpc1RawObj != nil {
		vpc1Raw = vpc1RawObj.(map[string]interface{})
	}
	if vpc1Raw["VpcId"] != nil {
		d.Set("vpc_id", vpc1Raw["VpcId"])
	}

	d.Set("instance_id", d.Id())

	return nil
}

func resourcePeerConnectionAccepterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ModifyVpcPeerConnection"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("peer_connection_accepter_name") {
		update = true
		request["Name"] = d.Get("peer_connection_accepter_name")
	}

	if d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if d.HasChange("link_type") {
		update = true
		request["LinkType"] = d.Get("link_type")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectStatus.VpcPeer"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		vpcPeerServiceV2 := Service{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionAccepterStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "PeerConnection"
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourcePeerConnectionAccepterRead(d, meta)
}

func resourcePeerConnectionAccepterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcPeerConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("force_delete"); ok {
		request["Force"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectStatus.VpcPeer", "IncorrectStatus", "OperationDenied.RouteEntryExist"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	vpcPeerServiceV2 := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionAccepterStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

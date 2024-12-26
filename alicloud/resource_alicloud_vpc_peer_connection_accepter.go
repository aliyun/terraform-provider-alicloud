// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcPeerPeerConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcPeerPeerConnectionAccepterCreate,
		Read:   resourceAliCloudVpcPeerPeerConnectionAccepterRead,
		Update: resourceAliCloudVpcPeerPeerConnectionAccepterUpdate,
		Delete: resourceAliCloudVpcPeerPeerConnectionAccepterDelete,
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
				ValidateFunc: IntAtLeast(1),
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
				ValidateFunc: StringInSlice([]string{"Platinum", "Gold", "Silver"}, false),
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

func resourceAliCloudVpcPeerPeerConnectionAccepterCreate(d *schema.ResourceData, meta interface{}) error {

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

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, true)
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
		if !IsExpectedErrors(err, []string{"IncorrectStatus.VpcPeer"}) {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_peer_connection_accepter", action, AlibabaCloudSdkGoERROR)
		}
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	vpcPeerServiceV2 := VpcPeerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionAccepterStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcPeerPeerConnectionAccepterUpdate(d, meta)
}

func resourceAliCloudVpcPeerPeerConnectionAccepterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcPeerServiceV2 := VpcPeerServiceV2{client}

	objectRaw, err := vpcPeerServiceV2.DescribeVpcPeerPeerConnectionAccepter(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_peer_connection_accepter DescribeVpcPeerPeerConnectionAccepter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudVpcPeerPeerConnectionAccepterUpdate(d *schema.ResourceData, meta interface{}) error {
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

	request["ClientToken"] = buildClientToken(action)
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.VpcPeer"}) || NeedRetry(err) {
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
		vpcPeerServiceV2 := VpcPeerServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionAccepterStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, false)
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
	}

	d.Partial(false)
	return resourceAliCloudVpcPeerPeerConnectionAccepterRead(d, meta)
}

func resourceAliCloudVpcPeerPeerConnectionAccepterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcPeerConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("force_delete"); ok {
		request["Force"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("VpcPeer", "2022-01-01", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.VpcPeer", "IncorrectStatus", "OperationDenied.RouteEntryExist"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcPeerServiceV2 := VpcPeerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionAccepterStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

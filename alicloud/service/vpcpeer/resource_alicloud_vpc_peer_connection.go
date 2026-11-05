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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePeerConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourcePeerConnectionCreate,
		Read:   resourcePeerConnectionRead,
		Update: resourcePeerConnectionUpdate,
		Delete: resourcePeerConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accepting_ali_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"accepting_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accepting_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"link_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Platinum", "Gold", "Silver"}, false),
			},
			"peer_connection_name": {
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
				Optional: true,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePeerConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcPeerConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["VpcId"] = d.Get("vpc_id")
	request["AcceptingAliUid"] = d.Get("accepting_ali_uid")
	request["AcceptingRegionId"] = d.Get("accepting_region_id")
	request["AcceptingVpcId"] = d.Get("accepting_vpc_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("peer_connection_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOkExists("bandwidth"); ok && v.(int) > 0 {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("link_type"); ok {
		request["LinkType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_peer_connection", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	vpcPeerServiceV2 := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Activated", "Accepting"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourcePeerConnectionUpdate(d, meta)
}

func resourcePeerConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcPeerServiceV2 := Service{client}

	objectRaw, err := vpcPeerServiceV2.DescribeVpcPeerPeerConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_peer_connection DescribeVpcPeerPeerConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["AcceptingOwnerUid"] != nil {
		d.Set("accepting_ali_uid", objectRaw["AcceptingOwnerUid"])
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
		d.Set("peer_connection_name", objectRaw["Name"])
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

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourcePeerConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		vpcPeerServiceV2 := Service{client}
		object, err := vpcPeerServiceV2.DescribeVpcPeerPeerConnection(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Activated" {
				action := "AcceptVpcPeerConnection"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()

				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := d.GetOk("resource_group_id"); ok {
					request["ResourceGroupId"] = v
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
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
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				vpcPeerServiceV2 := Service{client}
				stateConf := helper.BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "Rejected" {
				action := "RejectVpcPeerConnection"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["InstanceId"] = d.Id()

				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
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
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				vpcPeerServiceV2 := Service{client}
				stateConf := helper.BuildStateConf([]string{}, []string{"Rejected"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	action := "ModifyVpcPeerConnection"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if !d.IsNewResource() && d.HasChange("peer_connection_name") {
		update = true
		request["Name"] = d.Get("peer_connection_name")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if !d.IsNewResource() && d.HasChange("link_type") {
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
		vpcPeerServiceV2 := Service{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
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
	request["ResourceType"] = "PeerConnection"
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
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

	if d.HasChange("tags") {
		vpcPeerServiceV2 := Service{client}
		if err := vpcPeerServiceV2.SetResourceTags(d, "PeerConnection"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourcePeerConnectionRead(d, meta)
}

func resourcePeerConnectionDelete(d *schema.ResourceData, meta interface{}) error {

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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.InstanceId"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	vpcPeerServiceV2 := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcPeerServiceV2.VpcPeerPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

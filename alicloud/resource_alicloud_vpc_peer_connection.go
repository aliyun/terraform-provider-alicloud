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
				ValidateFunc: IntAtLeast(0),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"peer_connection_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(2, 128),
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
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudVpcPeerConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcPeerConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcpeerClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("accepting_ali_uid"); ok {
		request["AcceptingAliUid"] = v
	}
	request["AcceptingRegionId"] = d.Get("accepting_region_id")
	request["AcceptingVpcId"] = d.Get("accepting_vpc_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("peer_connection_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_peer_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Activated", "Accepting"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcPeerConnectionUpdate(d, meta)
}

func resourceAliCloudVpcPeerConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcPeerConnection(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_peer_connection DescribeVpcPeerConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accepting_ali_uid", objectRaw["AcceptingOwnerUid"])
	d.Set("accepting_region_id", objectRaw["AcceptingRegionId"])
	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("description", objectRaw["Description"])
	d.Set("peer_connection_name", objectRaw["Name"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	acceptingVpc1RawObj, _ := jsonpath.Get("$.AcceptingVpc", objectRaw)
	acceptingVpc1Raw := make(map[string]interface{})
	if acceptingVpc1RawObj != nil {
		acceptingVpc1Raw = acceptingVpc1RawObj.(map[string]interface{})
	}
	d.Set("accepting_vpc_id", acceptingVpc1Raw["VpcId"])
	vpc1RawObj, _ := jsonpath.Get("$.Vpc", objectRaw)
	vpc1Raw := make(map[string]interface{})
	if vpc1RawObj != nil {
		vpc1Raw = vpc1RawObj.(map[string]interface{})
	}
	d.Set("vpc_id", vpc1Raw["VpcId"])
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcPeerConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyVpcPeerConnection"
	conn, err := client.NewVpcpeerClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["InstanceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("bandwidth") {
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
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
		stateConf := BuildStateConf([]string{}, []string{"Activated", "Accepting"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("bandwidth")
		d.SetPartial("peer_connection_name")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewVpcpeerClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "PeerConnection"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		object, err := vpcServiceV2.DescribeVpcPeerConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Activated" {
				action = "AcceptVpcPeerConnection"
				conn, err = client.NewVpcpeerClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})

				request["InstanceId"] = d.Id()
				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOk("resource_group_id"); ok {
					request["ResourceGroupId"] = v
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

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
				stateConf := BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Rejected" {
				action = "RejectVpcPeerConnection"
				conn, err = client.NewVpcpeerClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})

				request["InstanceId"] = d.Id()
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

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
				stateConf := BuildStateConf([]string{}, []string{"Rejected"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
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
		if err := vpcServiceV2.SetVpcPeerResourceTags(d, "PeerConnection"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudVpcPeerConnectionRead(d, meta)
}

func resourceAliCloudVpcPeerConnectionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteVpcPeerConnection"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcpeerClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.InstanceId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcPeerConnectionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

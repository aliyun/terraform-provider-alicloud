// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package expressconnectrouter

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGrantAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceGrantAssociationCreate,
		Read:   resourceGrantAssociationRead,
		Delete: resourceGrantAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ecr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ecr_owner_ali_uid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGrantAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "GrantInstanceToExpressConnectRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrOwnerAliUid"] = d.Get("ecr_owner_ali_uid")
	request["InstanceId"] = d.Get("instance_id")
	request["InstanceType"] = d.Get("instance_type")
	request["InstanceRegionId"] = d.Get("instance_region_id")
	request["EcrId"] = d.Get("ecr_id")

	request["ClientToken"] = helper.BuildClientToken(action)

	request["CallerType"] = "OTHER"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Conflict.Lock"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_express_connect_router_grant_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v", request["EcrId"], request["InstanceId"], request["InstanceRegionId"], request["EcrOwnerAliUid"], request["InstanceType"]))

	return resourceGrantAssociationRead(d, meta)
}

func resourceGrantAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := NewService(client)

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterGrantAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_grant_association DescribeExpressConnectRouterGrantAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["EcrId"] != nil {
		d.Set("ecr_id", objectRaw["EcrId"])
	}
	if objectRaw["EcrOwnerAliUid"] != nil {
		d.Set("ecr_owner_ali_uid", objectRaw["EcrOwnerAliUid"])
	}
	if objectRaw["NodeId"] != nil {
		d.Set("instance_id", objectRaw["NodeId"])
	}
	if objectRaw["NodeRegionId"] != nil {
		d.Set("instance_region_id", objectRaw["NodeRegionId"])
	}
	if objectRaw["NodeType"] != nil {
		d.Set("instance_type", objectRaw["NodeType"])
	}

	return nil
}

func resourceGrantAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "RevokeInstanceFromExpressConnectRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrOwnerAliUid"] = parts[3]
	request["InstanceId"] = parts[1]
	request["InstanceType"] = parts[4]
	request["InstanceRegionId"] = parts[2]
	request["EcrId"] = parts[0]

	request["ClientToken"] = helper.BuildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Conflict.Lock"}) || helper.NeedRetry(err) {
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

	return nil
}

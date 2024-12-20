// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudExpressConnectRouterGrantAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectRouterGrantAssociationCreate,
		Read:   resourceAliCloudExpressConnectRouterGrantAssociationRead,
		Delete: resourceAliCloudExpressConnectRouterGrantAssociationDelete,
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

func resourceAliCloudExpressConnectRouterGrantAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "GrantInstanceToExpressConnectRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["EcrOwnerAliUid"] = d.Get("ecr_owner_ali_uid")
	request["InstanceId"] = d.Get("instance_id")
	request["InstanceType"] = d.Get("instance_type")
	request["InstanceRegionId"] = d.Get("instance_region_id")
	request["EcrId"] = d.Get("ecr_id")

	request["ClientToken"] = buildClientToken(action)

	request["CallerType"] = "OTHER"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_router_grant_association", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v", request["EcrId"], request["InstanceId"], request["InstanceRegionId"], request["EcrOwnerAliUid"], request["InstanceType"]))

	return resourceAliCloudExpressConnectRouterGrantAssociationRead(d, meta)
}

func resourceAliCloudExpressConnectRouterGrantAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterGrantAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_grant_association DescribeExpressConnectRouterGrantAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudExpressConnectRouterGrantAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "RevokeInstanceFromExpressConnectRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["EcrOwnerAliUid"] = parts[3]
	request["InstanceId"] = parts[1]
	request["InstanceType"] = parts[4]
	request["InstanceRegionId"] = parts[2]
	request["EcrId"] = parts[0]

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
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

	return nil
}

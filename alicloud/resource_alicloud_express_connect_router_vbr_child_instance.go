// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceCreate,
		Read:   resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceRead,
		Update: resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceUpdate,
		Delete: resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"child_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "The ID of the leased line gateway subinstance."),
			},
			"child_instance_owner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "The ID of the Alibaba Cloud account (primary account) to which the VBR instance belongs.> This parameter is required if you want to load a cross-account network instance."),
			},
			"child_instance_region_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Region of the leased line gateway sub-instance"),
			},
			"child_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "The type of the network instance. Value: **VBR**: VBR instance."),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecr_id": {
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

func resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AttachExpressConnectRouterChildInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ChildInstanceId"] = d.Get("child_instance_id")
	request["ChildInstanceType"] = d.Get("child_instance_type")
	request["EcrId"] = d.Get("ecr_id")

	request["ClientToken"] = buildClientToken(action)

	request["ChildInstanceRegionId"] = d.Get("child_instance_region_id")
	if v, ok := d.GetOkExists("child_instance_owner_id"); ok {
		request["ChildInstanceOwnerId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_router_vbr_child_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["EcrId"], request["ChildInstanceId"], request["ChildInstanceType"]))

	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceRead(d, meta)
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_vbr_child_instance DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["ChildInstanceOwnerId"] != nil {
		d.Set("child_instance_owner_id", objectRaw["ChildInstanceOwnerId"])
	}
	if objectRaw["ChildInstanceRegionId"] != nil {
		d.Set("child_instance_region_id", objectRaw["ChildInstanceRegionId"])
	}
	if objectRaw["GmtCreate"] != nil {
		d.Set("create_time", objectRaw["GmtCreate"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["ChildInstanceId"] != nil {
		d.Set("child_instance_id", objectRaw["ChildInstanceId"])
	}
	if objectRaw["ChildInstanceType"] != nil {
		d.Set("child_instance_type", objectRaw["ChildInstanceType"])
	}
	if objectRaw["EcrId"] != nil {
		d.Set("ecr_id", objectRaw["EcrId"])
	}

	return nil
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := "ModifyExpressConnectRouterChildInstance"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["ChildInstanceId"] = parts[1]
	request["ChildInstanceType"] = parts[2]

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
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

	return resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceRead(d, meta)
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DetachExpressConnectRouterChildInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["ChildInstanceId"] = parts[1]
	request["ChildInstanceType"] = parts[2]

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId", "ResourceNotFound.ChildInstanceId", "ResourceNotFound.EcrId"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"0"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectRouterServiceV2.DescribeAsyncExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(d, response, "$.TotalCount", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}

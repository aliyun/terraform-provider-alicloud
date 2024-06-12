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
		Delete: resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"child_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "The ID of the leased line gateway subinstance."),
			},
			"child_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "The ID of the leased line gateway subinstance."),
			},
			"child_instance_owner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9]*$"), "The ID of the subinstance of the leased line gateway."),
			},
			"child_instance_region_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Region of the leased line gateway sub-instance"),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
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
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ChildInstanceId"] = d.Get("child_instance_id")
	request["ChildInstanceType"] = d.Get("child_instance_type")
	request["EcrId"] = d.Get("ecr_id")

	request["ClientToken"] = buildClientToken(action)

	request["ChildInstanceRegionId"] = d.Get("child_instance_region_id")
	if v, ok := d.GetOk("child_instance_owner_id"); ok {
		request["ChildInstanceOwnerId"] = v
	}
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
		addDebug(action, response, request)
		return nil
	})

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

	d.Set("child_instance_owner_id", objectRaw["ChildInstanceOwnerId"])
	d.Set("child_instance_region_id", objectRaw["ChildInstanceRegionId"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("status", objectRaw["Status"])
	d.Set("child_instance_id", objectRaw["ChildInstanceId"])
	d.Set("ecr_id", objectRaw["EcrId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("ecr_id", parts[0])
	d.Set("child_instance_id", parts[1])
	d.Set("child_instance_type", parts[2])

	return nil
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVbrChildInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DetachExpressConnectRouterChildInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["ChildInstanceId"] = parts[1]
	request["ChildInstanceType"] = parts[2]

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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId", "ResourceNotFound.ChildInstanceId", "ResourceNotFound.EcrId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.DescribeAsyncExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(d, response, "$.TotalCount", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}
	return nil
}

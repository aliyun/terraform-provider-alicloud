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

func resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationCreate,
		Read:   resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationRead,
		Update: resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationUpdate,
		Delete: resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allowed_prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExpressConnectRouterAssociation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = d.Get("ecr_id")
	request["VpcId"] = d.Get("vpc_id")

	request["ClientToken"] = buildClientToken(action)

	request["AssociationRegionId"] = d.Get("association_region_id")
	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request["VpcOwnerId"] = v
	}
	if v, ok := d.GetOk("allowed_prefixes"); ok {
		allowedPrefixesMaps := v.([]interface{})
		request["AllowedPrefixes"] = allowedPrefixesMaps
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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_router_vpc_association", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["EcrId"], response["AssociationId"], request["VpcId"]))

	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationRead(d, meta)
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_vpc_association DescribeExpressConnectRouterExpressConnectRouterVpcAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("association_region_id", objectRaw["RegionId"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_owner_id", objectRaw["VpcOwnerId"])
	d.Set("association_id", objectRaw["AssociationId"])
	d.Set("ecr_id", objectRaw["EcrId"])
	d.Set("vpc_id", objectRaw["VpcId"])

	allowedPrefixes1Raw := make([]interface{}, 0)
	if objectRaw["AllowedPrefixes"] != nil {
		allowedPrefixes1Raw = objectRaw["AllowedPrefixes"].([]interface{})
	}

	d.Set("allowed_prefixes", allowedPrefixes1Raw)

	parts := strings.Split(d.Id(), ":")
	d.Set("ecr_id", parts[0])
	d.Set("association_id", parts[1])
	d.Set("vpc_id", parts[2])

	return nil
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyExpressConnectRouterAssociationAllowedPrefix"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["AssociationId"] = parts[1]
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("allowed_prefixes") {
		update = true
		if v, ok := d.GetOk("allowed_prefixes"); ok {
			allowedPrefixesMaps := v.([]interface{})
			request["AllowedPrefixes"] = allowedPrefixesMaps
		}
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationRead(d, meta)
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterVpcAssociationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExpressConnectRouterAssociation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["AssociationId"] = parts[1]

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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

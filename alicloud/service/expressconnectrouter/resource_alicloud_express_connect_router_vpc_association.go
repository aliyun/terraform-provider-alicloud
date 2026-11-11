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

func ResourceVpcAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcAssociationCreate,
		Read:   resourceVpcAssociationRead,
		Update: resourceVpcAssociationUpdate,
		Delete: resourceVpcAssociationDelete,
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
			"allowed_prefixes_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func resourceVpcAssociationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExpressConnectRouterAssociation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = d.Get("ecr_id")
	request["VpcId"] = d.Get("vpc_id")

	request["ClientToken"] = helper.BuildClientToken(action)

	request["AssociationRegionId"] = d.Get("association_region_id")
	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request["VpcOwnerId"] = v
	}
	if v, ok := d.GetOk("allowed_prefixes"); ok {
		allowedPrefixesMaps := v.([]interface{})
		request["AllowedPrefixes"] = allowedPrefixesMaps
	}
	if v, ok := d.GetOk("allowed_prefixes_mode"); ok {
		request["AllowedPrefixesMode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

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
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_express_connect_router_vpc_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["EcrId"], response["AssociationId"], request["VpcId"]))

	expressConnectRouterServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceVpcAssociationRead(d, meta)
}

func resourceVpcAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := NewService(client)

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_vpc_association DescribeExpressConnectRouterExpressConnectRouterVpcAssociation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("association_region_id", objectRaw["RegionId"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("description", objectRaw["Description"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_owner_id", objectRaw["VpcOwnerId"])
	d.Set("association_id", objectRaw["AssociationId"])
	d.Set("ecr_id", objectRaw["EcrId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("allowed_prefixes_mode", objectRaw["AllowedPrefixesMode"])

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

func resourceVpcAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	expressConnectRouterServiceV2 := NewService(client)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)

	// Update description via ModifyExpressConnectRouterAssociation.
	if d.HasChange("description") {
		action := "ModifyExpressConnectRouterAssociation"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["EcrId"] = parts[0]
		request["AssociationId"] = parts[1]
		request["Description"] = d.Get("description")
		request["ClientToken"] = helper.BuildClientToken(action)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				// IncorrectStatus.Ecr: another ECR operation is in progress; retry within the update timeout.
				if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Ecr", "Conflict.Lock"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	// Update allowed_prefixes / allowed_prefixes_mode via ModifyExpressConnectRouterAssociationAllowedPrefix.
	if d.HasChange("allowed_prefixes") || d.HasChange("allowed_prefixes_mode") {
		action := "ModifyExpressConnectRouterAssociationAllowedPrefix"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["EcrId"] = parts[0]
		request["AssociationId"] = parts[1]
		request["ClientToken"] = helper.BuildClientToken(action)
		if v, ok := d.GetOk("allowed_prefixes"); ok {
			request["AllowedPrefixes"] = v.([]interface{})
		}
		// The service falls back to MatchMode when AllowedPrefixesMode is omitted,
		// so always send the effective mode together with the prefixes to keep the
		// mode already applied on the association. An explicit empty string is
		// treated as unset and is never sent.
		if v, ok := d.GetOk("allowed_prefixes_mode"); ok {
			request["AllowedPrefixesMode"] = v
		}
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				// IncorrectStatus.Ecr / IncorrectStatus.EcrAssociation: a concurrent ECR or
				// association operation is in progress; retry within the update timeout.
				if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Ecr", "IncorrectStatus.EcrAssociation", "Conflict.Lock"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceVpcAssociationRead(d, meta)
}

func resourceVpcAssociationDelete(d *schema.ResourceData, meta interface{}) error {

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
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.AssociationId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	expressConnectRouterServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

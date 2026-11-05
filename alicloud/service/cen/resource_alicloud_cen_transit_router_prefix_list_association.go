package cen

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudCenTransitRouterPrefixListAssociation() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudCenTransitRouterPrefixListAssociationCreate,
		Read:   ResourceAlicloudCenTransitRouterPrefixListAssociationRead,
		Delete: ResourceAlicloudCenTransitRouterPrefixListAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"prefix_list_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BlackHole", "VPC", "VBR", "TR"}, false),
			},
			"owner_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceAlicloudCenTransitRouterPrefixListAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	var response map[string]interface{}
	action := "CreateTransitRouterPrefixListAssociation"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateTransitRouterPrefixListAssociation")
	request["PrefixListId"] = d.Get("prefix_list_id")
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["TransitRouterTableId"] = d.Get("transit_router_table_id")
	request["NextHop"] = d.Get("next_hop")

	if v, ok := d.GetOk("next_hop_type"); ok {
		request["NextHopType"] = v
	}

	if v, ok := d.GetOk("owner_uid"); ok {
		request["OwnerUid"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidStatus.Prefixlist", "IncorrectStatus.RouteTable", "ResourceNotFound.PrefixList"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_router_prefix_list_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["PrefixListId"], request["TransitRouterId"], request["TransitRouterTableId"], request["NextHop"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterPrefixListAssociationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAlicloudCenTransitRouterPrefixListAssociationRead(d, meta)
}

func ResourceAlicloudCenTransitRouterPrefixListAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	object, err := cbnService.DescribeCenTransitRouterPrefixListAssociation(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("prefix_list_id", object["PrefixListId"])
	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("transit_router_table_id", object["TransitRouterTableId"])
	d.Set("next_hop", object["NextHop"])
	d.Set("next_hop_type", object["NextHopType"])
	d.Set("owner_uid", object["OwnerUid"])
	d.Set("status", object["Status"])

	return nil
}

func ResourceAlicloudCenTransitRouterPrefixListAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterPrefixListAssociation"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"PrefixListId":         parts[0],
		"TransitRouterId":      parts[1],
		"TransitRouterTableId": parts[2],
		"NextHop":              parts[3],
	}
	request["ClientToken"] = helper.BuildClientToken("DeleteTransitRouterPrefixListAssociation")

	if v, ok := d.GetOk("next_hop_type"); ok {
		request["NextHopType"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.RouteTable", "IncorrectStatus.TransitRouter", "InvalidStatus.Prefixlist", "InvalidStatus.PrefixlistAssociation"}) || helper.NeedRetry(err) {
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

	return nil
}

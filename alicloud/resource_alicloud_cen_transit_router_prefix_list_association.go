package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenTransitRouterPrefixListAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterPrefixListAssociationCreate,
		Read:   resourceAlicloudCenTransitRouterPrefixListAssociationRead,
		Delete: resourceAlicloudCenTransitRouterPrefixListAssociationDelete,
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

func resourceAlicloudCenTransitRouterPrefixListAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterPrefixListAssociation"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateTransitRouterPrefixListAssociation")
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidStatus.Prefixlist", "IncorrectStatus.RouteTable", "ResourceNotFound.PrefixList"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_prefix_list_association", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["PrefixListId"], request["TransitRouterId"], request["TransitRouterTableId"], request["NextHop"]))

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterPrefixListAssociationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouterPrefixListAssociationRead(d, meta)
}

func resourceAlicloudCenTransitRouterPrefixListAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterPrefixListAssociation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAlicloudCenTransitRouterPrefixListAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterPrefixListAssociation"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"PrefixListId":         parts[0],
		"TransitRouterId":      parts[1],
		"TransitRouterTableId": parts[2],
		"NextHop":              parts[3],
	}
	request["ClientToken"] = buildClientToken("DeleteTransitRouterPrefixListAssociation")

	if v, ok := d.GetOk("next_hop_type"); ok {
		request["NextHopType"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.RouteTable", "IncorrectStatus.TransitRouter", "InvalidStatus.Prefixlist", "InvalidStatus.PrefixlistAssociation"}) || NeedRetry(err) {
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

	return nil
}

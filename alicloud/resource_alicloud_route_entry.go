package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteEntryCreate,
		Read:   resourceAliyunRouteEntryRead,
		Delete: resourceAliyunRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidrblock": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(2, 128),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"router_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute router_id has been deprecated and suggest removing it from your template.",
			},
		},
	}
}

func resourceAliyunRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	var cidr string
	rtId := d.Get("route_table_id").(string)
	nt := d.Get("nexthop_type").(string)
	ni := d.Get("nexthop_id").(string)

	table, err := vpcService.QueryRouteTableById(rtId)
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateCreateRouteEntryRequest()
	request.RegionId = client.RegionId
	request.RouteTableId = rtId

	if v, ok := d.GetOk("destination_cidrblock"); ok && v.(string) != "" {
		cidr = v.(string)
		if strings.Contains(v.(string), ":") {
			cidr = strings.Replace(v.(string), ":", "_", -1)
		}
		request.DestinationCidrBlock = v.(string)
	}

	request.NextHopType = nt
	request.NextHopId = ni
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("name"); ok {
		request.RouteEntryName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	// retry 10 min to create lots of entries concurrently
	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
			return resource.NonRetryableError(err)
		}

		args := *request

		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateRouteEntry(&args)
		})
		if err != nil {
			// Route Entry does not support concurrence when creating or deleting it;
			// Route Entry does not support creating or deleting within 5 seconds frequently
			// It must ensure all the route entries, vpc, vswitches' status must be available before creating or deleting route entry.
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectRouteEntryStatus", "IncorrectVpcStatus", "IncorrectVSwitchStatus", "IncorrectHaVipStatus", "IncorrectInstanceStatus", "InvalidVBRStatus", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.Ipv6Address", "LastTokenProcessing", "IncorrectStatus.VpcPeer", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectStatus.RouteTableStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "SystemBusy", "UnknownError", "IncorrectStatus.RouterInterface", "IncorrectStatus.PrefixList"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"RouterEntryConflict.Duplicated"}) {
			en, err := vpcService.DescribeRouteEntry(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)
			if err != nil {
				return WrapError(err)
			}
			return WrapError(Error("The route entry %s has already existed. "+
				"Please import it using ID '%s:%s:%s:%s:%s' or specify a new 'destination_cidrblock' and try again.",
				en.DestinationCidrBlock, en.RouteTableId, table.VRouterId, en.DestinationCidrBlock, en.NextHopType, ni))
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_route_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	// route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id

	d.SetId(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)

	if err := vpcService.WaitForRouteEntry(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunRouteEntryRead(d, meta)
}

func resourceAliyunRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}

	object, err := vpcService.DescribeRouteEntry(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("route_table_id", object.RouteTableId)
	d.Set("destination_cidrblock", object.DestinationCidrBlock)
	d.Set("nexthop_type", object.NextHopType)
	d.Set("nexthop_id", object.InstanceId)
	d.Set("name", object.RouteEntryName)
	d.Set("description", object.Description)
	d.Set("router_id", parts[1])

	return nil
}

func resourceAliyunRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}

	rtId := parts[0]
	if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	request, err := buildAliyunRouteEntryDeleteArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	request.RegionId = client.RegionId

	var raw interface{}
	retryTimes := 7
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouteEntry(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectVpcStatus", "TaskConflict", "OperationConflict", "SystemBusy", "IncorrectRouteEntryStatus", "IncorrectInstanceStatus", "Forbbiden", "UnknownError", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectHaVipStatus", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer", "IncorrectStatus.PrefixList"}) || NeedRetry(err) {
				// Route Entry does not support creating or deleting within 5 seconds frequently
				time.Sleep(time.Duration(retryTimes) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRouteEntry.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(vpcService.WaitForRouteEntry(d.Id(), Deleted, DefaultTimeout))
}

func buildAliyunRouteEntryDeleteArgs(d *schema.ResourceData, meta interface{}) (*vpc.DeleteRouteEntryRequest, error) {

	request := vpc.CreateDeleteRouteEntryRequest()
	request.RouteTableId = d.Get("route_table_id").(string)
	request.DestinationCidrBlock = d.Get("destination_cidrblock").(string)

	if v := d.Get("destination_cidrblock").(string); v != "" {
		request.DestinationCidrBlock = v
	}

	if v := d.Get("nexthop_id").(string); v != "" {
		request.NextHopId = v
	}

	return request, nil
}

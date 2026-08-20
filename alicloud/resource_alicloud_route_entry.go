package alicloud

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// routeEntryTaskLock serializes CreateRouteEntry/DeleteRouteEntry calls made
// by this provider process. The VPC service handles route entry operations as
// serialized asynchronous tasks and rejects concurrent submissions with
// "TaskConflict" ("The operation is too frequent"). Without client-side
// serialization, a plan that creates several route entries in parallel makes
// every resource burn its own retry budget on rejections caused by its
// siblings; when the retry window is short (e.g. the provider-level
// "max_retry_timeout" is set), the budget runs out while the server is still
// busy and the apply fails with a spurious TaskConflict. Queueing on this
// mutex happens before the retry loop starts, so waiting for a sibling
// operation does not consume the resource's own retry budget. TaskConflict
// remains in the retryable lists below because other actors (a second
// Terraform run, the console, other tools) can still trigger it.
var routeEntryTaskLock sync.Mutex

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

	// Serialize with every other route entry create/delete in this process
	// before the retry clock starts, so sibling operations do not eat into
	// this resource's retry budget by triggering TaskConflict rejections.
	routeEntryTaskLock.Lock()
	defer routeEntryTaskLock.Unlock()

	// retry 10 min to create lots of entries concurrently
	var raw interface{}
	attempt := 0
	var abandoned atomic.Bool
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		attempt++
		if err := vpcService.WaitForAllRouteEntriesAvailable(rtId, DefaultTimeout); err != nil {
			return retry.NonRetryableError(err)
		}

		// retry.Retry does not cancel an in-flight callback when its
		// deadline expires: it returns the last error to Terraform and leaves
		// the callback goroutine running. If that goroutine was blocked in
		// WaitForAllRouteEntriesAvailable above when the deadline fired, it
		// would - without this guard - still send CreateRouteEntry afterwards
		// and could succeed on the cloud long after this apply already
		// reported the resource as failed, leaving a route that exists
		// server-side with no record in state (the next apply then fails
		// with Duplicated.VpcNextHop). Skip the call once the loop has been
		// abandoned so no request is sent whose outcome nobody records.
		if abandoned.Load() {
			return retry.NonRetryableError(Error("the retry deadline of creating route entry has expired; skip sending CreateRouteEntry"))
		}

		args := *request

		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateRouteEntry(&args)
		})
		// Log every attempt (success or failure) instead of only the outcome of
		// the final one. Without this, an intermediate attempt that actually
		// succeeded on the server (e.g. after a retried "TaskConflict") leaves
		// no trace in TF_LOG, making it impossible to tell from the client side
		// whether the route was really created before a later step failed.
		addDebug(fmt.Sprintf("%s (attempt %d)", request.GetActionName(), attempt), raw, args.RpcRequest, args)
		if err != nil {
			// Route Entry does not support concurrence when creating or deleting it;
			// Route Entry does not support creating or deleting within 5 seconds frequently
			// It must ensure all the route entries, vpc, vswitches' status must be available before creating or deleting route entry.
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectRouteEntryStatus", "IncorrectVpcStatus", "IncorrectVSwitchStatus", "IncorrectHaVipStatus", "IncorrectInstanceStatus", "InvalidVBRStatus", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.Ipv6Address", "LastTokenProcessing", "IncorrectStatus.VpcPeer", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectStatus.RouteTableStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "SystemBusy", "UnknownError", "IncorrectStatus.RouterInterface", "IncorrectStatus.PrefixList"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	// Mark the loop abandoned before acting on its result, so a callback that
	// outlived the deadline (see the guard above) cannot create the route
	// after the outcome has been decided here.
	abandoned.Store(true)

	// Compute the resource id from the request parameters (all known up-front).
	// The id is recorded into state only after CreateRouteEntry has been
	// confirmed successful (err == nil below) or a client-side wait timeout
	// has been recovered by reading the route back. Setting it earlier -
	// before knowing whether the retry loop succeeded - left partial state
	// behind on every error path that forgot to reset it.
	// route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id
	id := buildRouteEntryResourceId(rtId, table.VRouterId, cidr, nt, ni)

	if err != nil {
		// A client-side wait timeout means the create result is unknown: an
		// earlier retry attempt may already have been committed on the server
		// side even though the loop ultimately gave up. Read the route back
		// before deciding to fail; if it does exist, keep it in state.
		if isRouteEntryWaitTimeout(err) {
			if recoverErr := waitAndReadRouteEntry(d, meta, id); recoverErr == nil {
				return nil
			}
		}
		// CreateRouteEntry returns RouterEntryConflict.Duplicated /
		// Duplicated.VpcNextHop / InvalidCIDRBlock.Duplicate only when a
		// matching route already exists on the cloud. buildClientToken is
		// called exactly once before the retry loop, so every attempt inside
		// the loop reuses the same ClientToken; with ClientToken idempotency,
		// an idempotent re-attempt of *this* apply would either succeed
		// (returning the same RouteEntryId) or keep hitting TaskConflict -
		// it would never return a duplicate. A duplicate therefore proves the
		// route was created by a different request (another apply, the console,
		// or another tool), not by this apply. Surface that as an error with
		// an import hint instead of silently adopting state we did not create.
		if IsExpectedErrors(err, []string{"RouterEntryConflict.Duplicated", "Duplicated.VpcNextHop", "InvalidCIDRBlock.Duplicate"}) {
			return WrapError(Error("The route entry %s already exists on the cloud, but it was not created by this apply (a duplicate error means another request - such as a prior apply, the console, or another tool - created it). "+
				"Please import it using ID '%s' or specify a new 'destination_cidrblock' and try again.",
				id, id))
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_route_entry", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(id)

	if err := vpcService.WaitForRouteEntry(d.Id(), Available, DefaultTimeout); err != nil {
		// CreateRouteEntry has already succeeded at this point, so the route is
		// known to exist; a timeout merely waiting for it to become Available
		// must not discard it from state. Fall back to reading it directly.
		if isRouteEntryWaitTimeout(err) {
			if recoverErr := resourceAliyunRouteEntryRead(d, meta); recoverErr == nil {
				return nil
			}
		}
		return WrapError(err)
	}

	return resourceAliyunRouteEntryRead(d, meta)
}

func buildRouteEntryResourceId(routeTableId, routerId, cidrBlock, nextHopType, nextHopId string) string {
	return routeTableId + ":" + routerId + ":" + cidrBlock + ":" + nextHopType + ":" + nextHopId
}

// isRouteEntryWaitTimeout reports whether err is a client-side timeout from
// waiting for a state change. Such a timeout says nothing about whether the
// route entry exists on the server side, so callers may read the route back
// before treating the operation as failed. Server-side errors - including
// duplicate errors, which only prove a matching route pre-exists, not that
// this Create call created it - are deliberately excluded.
func isRouteEntryWaitTimeout(err error) bool {
	return err != nil && strings.Contains(err.Error(), "timeout while waiting for state to become")
}

func waitAndReadRouteEntry(d *schema.ResourceData, meta interface{}, id string) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	if err := vpcService.WaitForRouteEntry(id, Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	// The route is confirmed to exist and be Available on the server side, so
	// it is safe to record the id now. Doing it only here - rather than before
	// the wait - guarantees a failed wait never leaves a partial id behind.
	d.SetId(id)
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

	// Serialize with every other route entry create/delete in this process
	// before the retry clock starts; see routeEntryTaskLock.
	routeEntryTaskLock.Lock()
	defer routeEntryTaskLock.Unlock()

	var raw interface{}
	retryTimes := 7
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouteEntry(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectVpcStatus", "TaskConflict", "OperationConflict", "SystemBusy", "IncorrectRouteEntryStatus", "IncorrectInstanceStatus", "Forbbiden", "UnknownError", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectHaVipStatus", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer", "IncorrectStatus.PrefixList"}) || NeedRetry(err) {
				// Route Entry does not support creating or deleting within 5 seconds frequently
				time.Sleep(time.Duration(retryTimes) * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

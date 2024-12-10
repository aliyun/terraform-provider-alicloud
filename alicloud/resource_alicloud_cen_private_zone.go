package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenPrivateZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenPrivateZoneCreate,
		Read:   resourceAliCloudCenPrivateZoneRead,
		Delete: resourceAliCloudCenPrivateZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_region_id": {
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

func resourceAliCloudCenPrivateZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	request := cbn.CreateRoutePrivateZoneInCenToVpcRequest()
	request.CenId = d.Get("cen_id").(string)
	request.AccessRegionId = d.Get("access_region_id").(string)
	request.HostVpcId = d.Get("host_vpc_id").(string)
	request.HostRegionId = d.Get("host_region_id").(string)

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.RoutePrivateZoneInCenToVpc(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.CenInstanceStatus", "InvalidOperation.NoChildInstanceEitherRegion"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_private_zone", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request.CenId, request.AccessRegionId))

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenPrivateZoneStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenPrivateZoneRead(d, meta)
}

func resourceAliCloudCenPrivateZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenPrivateZone(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("cen_id", parts[0])
	d.Set("access_region_id", object.AccessRegionId)
	d.Set("host_vpc_id", object.HostVpcId)
	d.Set("host_region_id", object.HostRegionId)
	d.Set("status", object.Status)

	return nil
}

func resourceAliCloudCenPrivateZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := cbn.CreateUnroutePrivateZoneInCenToVpcRequest()
	request.CenId = parts[0]
	request.AccessRegionId = parts[1]

	var raw interface{}
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err = client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.UnroutePrivateZoneInCenToVpc(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.CenInstanceStatus"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}

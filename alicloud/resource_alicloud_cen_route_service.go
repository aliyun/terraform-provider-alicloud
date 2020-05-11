package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCenRouteService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenRouteServiceCreate,
		Read:   resourceAlicloudCenRouteServiceRead,
		Delete: resourceAlicloudCenRouteServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"access_region_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_vpc_id": {
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

func resourceAlicloudCenRouteServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	request := cbn.CreateResolveAndRouteServiceInCenRequest()
	accessRegionIds := expandStringList(d.Get("access_region_ids").(*schema.Set).List())
	request.AccessRegionIds = &accessRegionIds
	request.CenId = d.Get("cen_id").(string)
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	request.Host = d.Get("host").(string)
	request.HostRegionId = d.Get("host_region_id").(string)
	request.HostVpcId = d.Get("host_vpc_id").(string)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ResolveAndRouteServiceInCen(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		d.SetId(d.Get("cen_id").(string) + ":" + d.Get("host_vpc_id").(string))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_route_service", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenRouteServiceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenRouteServiceRead(d, meta)
}
func resourceAlicloudCenRouteServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenRouteService(d.Id())
	if err != nil {
		if NotFoundError(err) {
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
	d.Set("host_vpc_id", parts[1])
	d.Set("access_region_id", object.AccessRegionId)
	//d.Set("access_region_ids", object.AccessRegionId)
	d.Set("description", object.Description)
	d.Set("host", object.Host)
	d.Set("host_region_id", object.HostRegionId)
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudCenRouteServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := cbn.CreateDeleteRouteServiceInCenRequest()
	request.CenId = parts[0]
	request.HostVpcId = parts[1]
	if v, ok := d.GetOk("access_region_id"); ok {
		request.AccessRegionId = v.(string)
	}
	request.Host = d.Get("host").(string)
	request.HostRegionId = d.Get("host_region_id").(string)
	err = resource.Retry(300*time.Second, func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteRouteServiceInCen(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.CloudRouteStatusNotAllow", "Operation.Blocking", "InvalidOperation.CenInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

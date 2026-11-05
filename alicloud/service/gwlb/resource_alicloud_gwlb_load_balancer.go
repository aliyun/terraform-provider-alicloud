package gwlb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerCreate,
		Read:   resourceLoadBalancerRead,
		Update: resourceLoadBalancerUpdate,
		Delete: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_mappings": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"load_balancer_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eni_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ipv4_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ZoneId"] = dataLoop1Tmp["zone_id"]
			dataLoop1Map["VSwitchId"] = dataLoop1Tmp["vswitch_id"]
			zoneMappingsMapsArray = append(zoneMappingsMapsArray, dataLoop1Map)
		}
		request["ZoneMappings"] = zoneMappingsMapsArray
	}

	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gwlb_load_balancer", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))

	gwlbService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gwlbService.LoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gwlbService := NewService(client)

	objectRaw, err := gwlbService.DescribeLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gwlb_load_balancer DescribeLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["AddressIpVersion"] != nil {
		d.Set("address_ip_version", objectRaw["AddressIpVersion"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["LoadBalancerName"] != nil {
		d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["LoadBalancerStatus"] != nil {
		d.Set("status", objectRaw["LoadBalancerStatus"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))
	zoneMappings1Raw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappings1Raw != nil {
		for _, zoneMappingsChild1Raw := range zoneMappings1Raw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChild1Raw := zoneMappingsChild1Raw.(map[string]interface{})
			zoneMappingsMap["vswitch_id"] = zoneMappingsChild1Raw["VSwitchId"]
			zoneMappingsMap["zone_id"] = zoneMappingsChild1Raw["ZoneId"]

			loadBalancerAddresses1Raw := zoneMappingsChild1Raw["LoadBalancerAddresses"]
			loadBalancerAddressesMaps := make([]map[string]interface{}, 0)
			if loadBalancerAddresses1Raw != nil {
				for _, loadBalancerAddressesChild1Raw := range loadBalancerAddresses1Raw.([]interface{}) {
					loadBalancerAddressesMap := make(map[string]interface{})
					loadBalancerAddressesChild1Raw := loadBalancerAddressesChild1Raw.(map[string]interface{})
					loadBalancerAddressesMap["eni_id"] = loadBalancerAddressesChild1Raw["EniId"]
					loadBalancerAddressesMap["private_ipv4_address"] = loadBalancerAddressesChild1Raw["PrivateIpv4Address"]

					loadBalancerAddressesMaps = append(loadBalancerAddressesMaps, loadBalancerAddressesMap)
				}
			}
			zoneMappingsMap["load_balancer_addresses"] = loadBalancerAddressesMaps
			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}
	}
	if objectRaw["ZoneMappings"] != nil {
		if err := d.Set("zone_mappings", zoneMappingsMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateLoadBalancerAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
		gwlbService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gwlbService.LoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerZones"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("zone_mappings") {
		update = true
	}
	if v, ok := d.GetOk("zone_mappings"); ok || d.HasChange("zone_mappings") {
		zoneMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			zoneMappingsMapsArray = append(zoneMappingsMapsArray, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMapsArray
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
		gwlbService := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gwlbService.LoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ResourceType"] = "loadbalancer"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
	}

	if d.HasChange("tags") {
		gwlbService := NewService(client)
		if err := gwlbService.SetResourceTags(d, "loadbalancer"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Gwlb", "2024-04-15", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.LoadBalancer"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

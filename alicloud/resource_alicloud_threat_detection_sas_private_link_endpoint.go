package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionSasPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionSasPrivateLinkEndpointCreate,
		Read:   resourceAlicloudThreatDetectionSasPrivateLinkEndpointRead,
		Update: resourceAlicloudThreatDetectionSasPrivateLinkEndpointUpdate,
		Delete: resourceAlicloudThreatDetectionSasPrivateLinkEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"node_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"region_id": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"vpc_id": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"security_group_id": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"zones": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"v_switch_id": {
							Required: true,
							Type:     schema.TypeString,
						},
						"zone_id": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"update_domain": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"jsrv_domain": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudThreatDetectionSasPrivateLinkEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	request := make(map[string]interface{})

	if v, ok := d.GetOk("node_name"); ok {
		request["NodeName"] = v
	}
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v
	} else {
		request["RegionId"] = client.RegionId
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("zones"); ok {
		zonesMaps := make([]map[string]interface{}, 0)
		for _, value := range v.(*schema.Set).List() {
			zone := value.(map[string]interface{})
			zoneMap := map[string]interface{}{
				"VSwitchId": zone["v_switch_id"],
				"ZoneId":    zone["zone_id"],
			}
			zonesMaps = append(zonesMaps, zoneMap)
		}
		request["Zones"] = zonesMaps
	}

	var response map[string]interface{}
	action := "CreateSasPrivateLinkEndpoint"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_sas_private_link_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.NodeId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_sas_private_link_endpoint")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	sasService := SasService{client}
	stateConf := BuildStateConf([]string{}, []string{"enable"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, sasService.ThreatDetectionSasPrivateLinkEndpointStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudThreatDetectionSasPrivateLinkEndpointRead(d, meta)
}

func resourceAlicloudThreatDetectionSasPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionSasPrivateLinkEndpoint(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_sas_private_link_endpoint sasService.DescribeThreatDetectionSasPrivateLinkEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("node_name", object["NodeName"])
	d.Set("region_id", object["RegionId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("update_domain", object["UpdateDomain"])
	d.Set("jsrv_domain", object["JsrvDomain"])
	if zonesRaw, ok := object["Zones"]; ok && zonesRaw != nil {
		zonesMaps := make([]map[string]interface{}, 0)
		if zonesList, ok := zonesRaw.([]interface{}); ok {
			for _, value := range zonesList {
				zone := value.(map[string]interface{})
				zoneMap := map[string]interface{}{
					"v_switch_id": zone["VSwitchId"],
					"zone_id":     zone["ZoneId"],
				}
				zonesMaps = append(zonesMaps, zoneMap)
			}
		}
		d.Set("zones", zonesMaps)
	}
	return nil
}

func resourceAlicloudThreatDetectionSasPrivateLinkEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("node_name") {
		update = true
		if v, ok := d.GetOk("node_name"); ok {
			request["NodeName"] = v
		}
	}

	if update {
		action := "UpdateSasPrivateLinkEndpoint"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudThreatDetectionSasPrivateLinkEndpointRead(d, meta)
}

func resourceAlicloudThreatDetectionSasPrivateLinkEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	var err error

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	action := "DeleteSasPrivateLinkEndpoint"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, sasService.ThreatDetectionSasPrivateLinkEndpointStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

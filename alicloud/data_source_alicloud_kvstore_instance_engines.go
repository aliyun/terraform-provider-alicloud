package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"time"
)

func dataSourceAlicloudKVStoreInstanceEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStoreInstanceEnginesRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
				ValidateFunc: validateAllowedStringValue([]string{string(PostPaid), string(PrePaid)}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"instance_engines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKVStoreInstanceEnginesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := r_kvstore.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	request.InstanceChargeType = instanceChargeType
	var response = &r_kvstore.DescribeAvailableResourceResponse{}
	err := resource.Retry(time.Minute*5, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAvailableResource(request)
		})
		if err != nil {
			if IsExceptedError(err, Throttling) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response = raw.(*r_kvstore.DescribeAvailableResourceResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instance_engines", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var infos []map[string]interface{}
	var ids []string

	engine, engineGot := d.GetOk("engine")
	engine = strings.ToLower(engine)
	engineVersion, engineVersionGot := d.GetOk("engine_version")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		info := make(map[string]interface{})
		zondId := AvailableZone.ZoneId
		info["zone_id"] = AvailableZone.ZoneId
		ids = append(ids, zondId)
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine.(string) != SupportedEngine.Engine {
				continue
			}
			info["engine"] = SupportedEngine.Engine
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersions.SupportedEngineVersion {
				if engineVersionGot && engineVersion != SupportedEngineVersion.Version {
					continue
				}
				info["engine_version"] = SupportedEngineVersion.Version
				ids = append(ids, SupportedEngineVersion.Version)
				temp := make(map[string]interface{}, len(info))
				for key, value := range info {
					temp[key] = value
				}
				infos = append(infos, temp)
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_engines", infos)
	if err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

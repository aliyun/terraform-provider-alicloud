package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDBInstanceEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstanceEngiesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateAllowedStringValue([]string{string(PostPaid), string(PrePaid)}),
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
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
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstanceEngiesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := rds.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.ZoneId = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	if instanceChargeType == string(PostPaid) {
		instanceChargeType = string(Postpaid)
	} else {
		instanceChargeType = string(Prepaid)
	}
	request.InstanceChargeType = instanceChargeType

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instance_engines", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	response := raw.(*rds.DescribeAvailableResourceResponse)

	infos := []map[string]interface{}{}
	ids := []string{}

	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		info := make(map[string]interface{})
		info["zone_id"] = AvailableZone.ZoneId
		ids = append(ids, AvailableZone.ZoneId)
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine != SupportedEngine.Engine {
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

				for _, SupportedCategory := range SupportedEngineVersion.SupportedCategorys.SupportedCategory {
					info["category"] = SupportedCategory.Category
					temp := make(map[string]interface{}, len(info))
					for key, value := range info {
						temp[key] = value
					}
					infos = append(infos, temp)
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_engines", infos)
	if err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		err = writeToFile(output.(string), infos)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

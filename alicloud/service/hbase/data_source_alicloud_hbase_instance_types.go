package hbase

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceHBaseInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHBaseInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, true),
				Optional:     true,
				Default:      PostPaid,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"singlehbase", "hbase", "hbaseue", "bds"}, false),
				Optional:     true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"cloud_ssd", "cloud_essd_pl1", "cloud_efficiency", "local_hdd_pro", "local_ssd_pro"}, false),
				Optional:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"master_instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"core_instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHBaseInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	response, err := describeAvaiableResource(d, meta)
	if err != nil {
		return errs.WrapErrorf(err, "describeAvaiableResource failed")
	}

	var ids []string
	var masterInstanceTypes []map[string]interface{}
	typeSet := mapset.NewSet()
	var types []map[string]interface{}
	masterTypeSet := mapset.NewSet()
	var coreInstanceTypes []map[string]interface{}
	for _, zone := range response.AvailableZones.AvailableZone {
		for _, master := range zone.MasterResources.MasterResource {
			if masterTypeSet.Contains(master.InstanceType) {
				continue
			}
			e := map[string]interface{}{
				"value":    master.InstanceType,
				"cpu_size": master.InstanceTypeDetail.Cpu,
				"mem_size": master.InstanceTypeDetail.Mem,
			}
			types = append(types, e)

			e = map[string]interface{}{
				"instance_type": master.InstanceType,
				"cpu_size":      master.InstanceTypeDetail.Cpu,
				"mem_size":      master.InstanceTypeDetail.Mem,
			}
			masterInstanceTypes = append(masterInstanceTypes, e)
			masterTypeSet.Add(master.InstanceType)
			if !typeSet.Contains(master.InstanceType) {
				typeSet.Add(master.InstanceType)
			}

		}
		for _, engine := range zone.SupportedEngines.SupportedEngine {
			for _, verion := range engine.SupportedEngineVersions.SupportedEngineVersion {
				for _, category := range verion.SupportedCategories.SupportedCategoriesItem {
					for _, storageType := range category.SupportedStorageTypes.SupportedStorageType {
						for _, coreResource := range storageType.CoreResources.CoreResource {
							if !typeSet.Contains(coreResource.InstanceType) {
								typeSet.Add(coreResource.InstanceType)
							}
							e := map[string]interface{}{
								"zone":           zone.ZoneId,
								"engine":         engine.Engine,
								"version":        verion.Version,
								"category":       category.Category,
								"storage_type":   storageType.StorageType,
								"instance_type":  coreResource.InstanceType,
								"cpu_size":       coreResource.InstanceTypeDetail.Cpu,
								"mem_size":       coreResource.InstanceTypeDetail.Mem,
								"max_core_count": coreResource.MaxCoreCount,
							}
							coreInstanceTypes = append(coreInstanceTypes, e)
						}
					}
				}
			}
		}
	}

	for _, v := range typeSet.ToSlice() {
		ids = append(ids, v.(string))
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("types", types); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("master_instance_types", masterInstanceTypes); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("core_instance_types", coreInstanceTypes); err != nil {
		return errs.WrapError(err)
	}

	result := map[string]interface{}{
		"master_instance_types": masterInstanceTypes,
		"core_instance_types":   coreInstanceTypes,
		"types":                 types,
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = helper.WriteToFile(output.(string), result)
	}
	return nil
}

func describeAvaiableResource(d *schema.ResourceData, meta interface{}) (result hbase.DescribeAvailableResourceResponse, err error) {
	client := meta.(*connectivity.AliyunClient)
	request := hbase.CreateDescribeAvailableResourceRequest()
	if chargeType, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = strings.TrimSpace(chargeType.(string))
	}
	if regionId, ok := d.GetOk("region_id"); ok {
		request.RegionId = strings.TrimSpace(regionId.(string))
	} else {
		request.RegionId = string(client.Region)
	}
	if zoneId, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = strings.TrimSpace(zoneId.(string))
	}
	if engine, ok := d.GetOk("engine"); ok {
		request.Engine = strings.TrimSpace(engine.(string))
	}
	if version, ok := d.GetOk("version"); ok {
		request.EngineVersion = strings.TrimSpace(version.(string))
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.CoreInstanceType = strings.TrimSpace(instanceType.(string))
	}
	if diskType, ok := d.GetOk("disk_type"); ok {
		request.DiskType = strings.TrimSpace(diskType.(string))
	}

	raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
		return hbaseClient.DescribeAvailableResource(request)
	})

	response, _ := raw.(*hbase.DescribeAvailableResourceResponse)
	if err != nil {
		return result, errs.WrapErrorf(err, errs.DefaultErrorMsg, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response == nil {
		return result, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_hbase_instance_types", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return *response, nil
}

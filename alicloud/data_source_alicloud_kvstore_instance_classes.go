package alicloud

import (
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudKVStoreInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStoreAvailableResourceRead,
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
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"standard", "cluster", "rwsplit"}),
			},
			"performance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"standard_performance_type", "enhance_performance_type"}),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"inmemory", "hybrid"}),
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"double", "single", "readone", "readthree", "readfive"}),
			},
			"package_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"standard", "customized"}),
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
			"instance_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAlicloudKVStoreAvailableResourceRead(d *schema.ResourceData, meta interface{}) error {
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var instanceClasses []string
	var ids []string

	engine, engineGot := d.GetOk("engine")
	engineVersion, engineVersionGot := d.GetOk("engine_version")
	architecture, architectureGot := d.GetOk("architecture")
	performanceType, performanceTypeGot := d.GetOk("performance_type")
	storageType, storageTypeGot := d.GetOk("storage_type")
	nodeType, nodeTypeGot := d.GetOk("node_type")
	packageType, packageTypeGot := d.GetOk("package_type")

	for _, AvailableZone := range response.AvailableZones.AvailableZone {
		zondId := AvailableZone.ZoneId
		ids = append(ids, zondId)
		for _, SupportedEngine := range AvailableZone.SupportedEngines.SupportedEngine {
			if engineGot && engine.(string) != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersions.SupportedEngineVersion {
				if engineVersionGot && engineVersion != SupportedEngineVersion.Version {
					continue
				}
				ids = append(ids, SupportedEngineVersion.Version)
				for _, SupportedArchitectureType := range SupportedEngineVersion.SupportedArchitectureTypes.SupportedArchitectureType {
					if architectureGot && architecture.(string) != SupportedArchitectureType.Architecture {
						continue
					}
					for _, SupportedPerformanceType := range SupportedArchitectureType.SupportedPerformanceTypes.SupportedPerformanceType {
						if performanceTypeGot && performanceType.(string) != SupportedPerformanceType.PerformanceType {
							continue
						}
						for _, SupportedStorageType := range SupportedPerformanceType.SupportedStorageTypes.SupportedStorageType {
							if storageTypeGot && storageType.(string) != SupportedStorageType.StorageType {
								continue
							}
							for _, SupportedNodeType := range SupportedStorageType.SupportedNodeTypes.SupportedNodeType {
								if nodeTypeGot && nodeType.(string) != SupportedNodeType.NodeType {
									continue
								}
								for _, SupportedPackageType := range SupportedNodeType.SupportedPackageTypes.SupportedPackageType {
									if packageTypeGot && packageType.(string) != SupportedPackageType.PackageType {
										continue
									}
									for _, AvailableResource := range SupportedPackageType.AvailableResources.AvailableResource {
										instanceClasses = append(instanceClasses, AvailableResource.InstanceClass)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	d.SetId(dataResourceIdHash(ids))
	err = d.Set("instance_classes", instanceClasses)
	if err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), instanceClasses)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

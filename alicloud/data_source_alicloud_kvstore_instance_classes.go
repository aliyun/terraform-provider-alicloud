package alicloud

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.StringInSlice([]string{"standard", "cluster", "rwsplit"}, false),
			},
			"performance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard_performance_type", "enhance_performance_type"}, false),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"inmemory", "hybrid"}, false),
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"double", "single", "readone", "readthree", "readfive"}, false),
			},
			"package_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "customized"}, false),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PrePaid,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Price"}, false),
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
			"classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
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
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*r_kvstore.DescribeAvailableResourceResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_instance_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	var instanceClasses []string
	var ids []string

	engine, engineGot := d.GetOk("engine")
	engine = strings.ToLower(engine.(string))
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
			if engineGot && engine != SupportedEngine.Engine {
				continue
			}
			ids = append(ids, SupportedEngine.Engine)
			for _, SupportedEngineVersion := range SupportedEngine.SupportedEngineVersions.SupportedEngineVersion {
				if engineVersionGot && engineVersion.(string) != SupportedEngineVersion.Version {
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

	var instanceClassPrices []map[string]interface{}
	sortedBy := d.Get("sorted_by").(string)
	if sortedBy == "Price" && len(instanceClasses) > 0 {
		bssopenapiService := BssopenapiService{client}
		priceList, err := getKVStoreInstanceClassPrice(bssopenapiService, instanceChargeType, instanceClasses)
		if err != nil {
			return WrapError(err)
		}
		for i, instanceClass := range instanceClasses {
			classPrice := map[string]interface{}{
				"instance_class": instanceClass,
				"price":          fmt.Sprintf("%.4f", priceList[i]),
			}
			instanceClassPrices = append(instanceClassPrices, classPrice)
		}
		sort.SliceStable(instanceClassPrices, func(i, j int) bool {
			iPrice, _ := strconv.ParseFloat(instanceClassPrices[i]["price"].(string), 64)
			jPrice, _ := strconv.ParseFloat(instanceClassPrices[j]["price"].(string), 64)
			return iPrice < jPrice
		})

		err = d.Set("classes", instanceClassPrices)
		if err != nil {
			return WrapError(err)
		}

		instanceClasses = instanceClasses[:0]
		for _, instanceClass := range instanceClassPrices {
			instanceClasses = append(instanceClasses, instanceClass["instance_class"].(string))
		}
	}

	err = d.Set("instance_classes", instanceClasses)
	if err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), instanceClassPrices)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func getKVStoreInstanceClassPrice(bssopenapiService BssopenapiService, instanceChargeType string, instanceClasses []string) ([]float64, error) {
	client := bssopenapiService.client
	var modules interface{}
	moduleCode := "InstanceClass"
	var payAsYouGo []bssopenapi.GetPayAsYouGoPriceModuleList
	var subsciption []bssopenapi.GetSubscriptionPriceModuleList
	for _, instanceClass := range instanceClasses {
		config := fmt.Sprintf("InstanceClass:%s,Region:%s", instanceClass, client.Region)
		if instanceChargeType == string(PostPaid) {
			payAsYouGo = append(payAsYouGo, bssopenapi.GetPayAsYouGoPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
				PriceType:  "Hour",
			})
		} else {
			subsciption = append(subsciption, bssopenapi.GetSubscriptionPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
			})

		}
	}

	if len(payAsYouGo) != 0 {
		modules = payAsYouGo
	} else {
		modules = subsciption
	}

	return bssopenapiService.GetInstanceTypePrice("redisa", "", modules)
}

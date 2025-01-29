package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"sort"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.StringInSlice([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}, false),
				Default: string(KVStoreRedis),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Tair_rdb", "Tair_scm", "Tair_essd", "OnECS"}, false),
			},
			"performance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard_performance_type", "enhance_performance_type"}, false),
				Deprecated:   "The parameter 'performance_type' has been deprecated from 1.68.0.",
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"inmemory", "hybrid"}, false),
				Deprecated:   "The parameter 'storage_type' has been deprecated from 1.68.0.",
			},
			"package_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "customized"}, false),
				Deprecated:   "The parameter 'package_type' has been deprecated from 1.68.0.",
			},
			"architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "cluster", "rwsplit"}, false),
			},
			"edition_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Community", "Enterprise"}, false),
			},
			"series_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enhanced_performance_type", "hybrid_storage"}, false),
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"double", "single", "readone", "readthree", "readfive"}, false),
			},
			"shard_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 4, 8, 16, 32, 64, 128, 256}),
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
	var response map[string]interface{}
	var err error
	action := "DescribeAvailableResource"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ZoneId"] = d.Get("zone_id").(string)
	instanceChargeType := d.Get("instance_charge_type").(string)
	request["InstanceChargeType"] = instanceChargeType
	request["Engine"] = d.Get("engine").(string)
	request["ProductType"] = d.Get("product_type").(string)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var instanceClasses []string
	var ids []string

	engineVersion, engineVersionGot := d.GetOk("engine_version")
	architecture, architectureGot := d.GetOk("architecture")
	editionType, editionTypeGot := d.GetOk("edition_type")
	seriesType, seriesTypeGot := d.GetOk("series_type")
	shardNumber, shardNumberGot := d.GetOk("shard_number")
	nodeType, nodeTypeGot := d.GetOk("node_type")

	resp, err := jsonpath.Get("$.AvailableZones.AvailableZone", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AvailableZones.AvailableZone", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		AvailableZone := v.(map[string]interface{})
		zondId := AvailableZone["ZoneId"].(string)
		ids = append(ids, zondId)
		for _, v := range AvailableZone["SupportedEngines"].(map[string]interface{})["SupportedEngine"].([]interface{}) {
			SupportedEngine := v.(map[string]interface{})
			ids = append(ids, SupportedEngine["Engine"].(string))
			for _, v := range SupportedEngine["SupportedEditionTypes"].(map[string]interface{})["SupportedEditionType"].([]interface{}) {
				SupportedEditionType := v.(map[string]interface{})
				if editionTypeGot && editionType.(string) != fmt.Sprint(SupportedEditionType["EditionType"]) {
					continue
				}
				ids = append(ids, SupportedEditionType["EditionType"].(string))
				for _, v := range SupportedEditionType["SupportedSeriesTypes"].(map[string]interface{})["SupportedSeriesType"].([]interface{}) {
					SupportedSeriesType := v.(map[string]interface{})
					if seriesTypeGot && seriesType.(string) != fmt.Sprint(SupportedSeriesType["SeriesType"]) {
						continue
					}
					for _, v := range SupportedSeriesType["SupportedEngineVersions"].(map[string]interface{})["SupportedEngineVersion"].([]interface{}) {
						SupportedEngineVersion := v.(map[string]interface{})
						if engineVersionGot && engineVersion.(string) != fmt.Sprint(SupportedEngineVersion["Version"]) {
							continue
						}
						for _, v := range SupportedEngineVersion["SupportedArchitectureTypes"].(map[string]interface{})["SupportedArchitectureType"].([]interface{}) {
							SupportedArchitectureType := v.(map[string]interface{})
							if architectureGot && architecture.(string) != fmt.Sprint(SupportedArchitectureType["Architecture"]) {
								continue
							}
							for _, v := range SupportedArchitectureType["SupportedShardNumbers"].(map[string]interface{})["SupportedShardNumber"].([]interface{}) {
								SupportedShardNumber := v.(map[string]interface{})
								number, _ := strconv.Atoi(fmt.Sprint(SupportedShardNumber["ShardNumber"]))
								if shardNumberGot && shardNumber.(int) != number {
									continue
								}
								for _, v := range SupportedShardNumber["SupportedNodeTypes"].(map[string]interface{})["SupportedNodeType"].([]interface{}) {
									SupportedNodeType := v.(map[string]interface{})
									if nodeTypeGot && nodeType.(string) != fmt.Sprint(SupportedNodeType["SupportedNodeType"]) {
										continue
									}
									for _, v := range SupportedNodeType["AvailableResources"].(map[string]interface{})["AvailableResource"].([]interface{}) {
										AvailableResource := v.(map[string]interface{})
										instanceClasses = append(instanceClasses, AvailableResource["InstanceClass"].(string))
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
		bssopenapiService := BssOpenApiService{client}
		moduleCode := "InstanceClass"
		modules := make([]map[string]interface{}, 0)
		for _, instanceClass := range instanceClasses {
			config := fmt.Sprintf("InstanceClass:%s,Region:%s", instanceClass, client.Region)
			if instanceChargeType == string(PostPaid) {
				modules = append(modules, map[string]interface{}{
					"ModuleCode": moduleCode,
					"Config":     config,
					"PriceType":  "Hour",
				})
			} else {
				modules = append(modules, map[string]interface{}{
					"ModuleCode": moduleCode,
					"Config":     config,
				})

			}
		}
		paymentType := "PayAsYouGo"
		if instanceChargeType == string(PrePaid) {
			paymentType = "Subscription"
		}

		priceList, err := bssopenapiService.GetInstanceTypePrice("redisa", "", paymentType, modules)
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

package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type instanceTypeWithOriginalPrice struct {
	InstanceType  map[string]interface{}
	OriginalPrice float64
}

func dataSourceAlicloudInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type_family": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"memory_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"gpu_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gpu_spec": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  PostPaid,
				// %q must contain a valid InstanceChargeType, expected common.PrePaid, common.PostPaid
				ValidateFunc: StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Vpc", "Classic"}, false),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      NoSpot,
				ValidateFunc: StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			},
			"eni_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"kubernetes_node_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: StringInSlice([]string{
					string(KubernetesNodeMaster),
					string(KubernetesNodeWorker),
				}, false),
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sorted_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: StringInSlice([]string{
					"CPU",
					"Memory",
					"Price",
				}, false),
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "cloud_essd_entry", "cloud_auto"}, false),
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
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"minimum_eni_ipv6_address_quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"minimum_eni_private_ip_address_quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			// Computed values.
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"gpu": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
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
						"burstable_instance": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_credit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"baseline_credit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"eni_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eni_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"primary_eni_queue_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"secondary_eni_queue_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eni_ipv6_address_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"maximum_queue_number_per_eni": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_eni_queue_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eni_private_ip_address_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nvme_support": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_storage": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"amount": {
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
				},
			},
		},
	}
}

func dataSourceAlicloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	zoneId, validZones, _, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}

	mapInstanceTypes := make(map[string][]string)
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			for _, t := range r.SupportedResources.SupportedResource {
				if t.Status == string(SoldOut) {
					continue
				}
				if v, ok := mapInstanceTypes[t.Value]; ok {
					v = append(v, zone.ZoneId)
					mapInstanceTypes[t.Value] = v
				} else {
					mapInstanceTypes[t.Value] = []string{zone.ZoneId}
				}
			}
		}
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := d.Get("memory_size").(float64)
	gpuAmount := d.Get("gpu_amount").(int)

	request := make(map[string]interface{})

	if v, ok := d.GetOkExists("minimum_eni_ipv6_address_quantity"); ok {
		request["MinimumEniIpv6AddressQuantity"] = v
	}

	if v, ok := d.GetOkExists("minimum_eni_private_ip_address_quantity"); ok {
		request["MinimumEniPrivateIpAddressQuantity"] = v
	}

	if v, ok := d.GetOk("instance_type_family"); ok {
		request["InstanceTypeFamily"] = v
	}

	if v, ok := d.GetOk("gpu_spec"); ok {
		request["GPUSpec"] = v
	}
	if v, ok := d.GetOk("system_disk_category"); ok {
		request["SystemDiskCategory"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeInstanceTypes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instance_types", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceTypes.InstanceType", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceTypes.InstanceType", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	var instanceTypes []instanceTypeWithOriginalPrice
	imageSupportInstanceTypesMap := make(map[string]struct{}, 0)
	imageId := strings.TrimSpace(d.Get("image_id").(string))
	if imageId != "" {
		request = map[string]interface{}{
			"ImageId":  imageId,
			"RegionId": client.RegionId,
		}
		action := "DescribeImageSupportInstanceTypes"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instance_types", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.InstanceTypes.InstanceType", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceTypes.InstanceType", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			imageSupportInstanceTypesMap[fmt.Sprint(v.(map[string]interface{})["InstanceTypeId"])] = struct{}{}
		}
	}

	eniAmount := d.Get("eni_amount").(int)
	k8sNode := strings.TrimSpace(d.Get("kubernetes_node_role").(string))

	for _, v := range objects {
		object := v.(map[string]interface{})
		if _, ok := mapInstanceTypes[object["InstanceTypeId"].(string)]; !ok {
			continue
		}

		if imageId != "" {
			if len(imageSupportInstanceTypesMap) > 0 {
				if _, ok := imageSupportInstanceTypesMap[object["InstanceTypeId"].(string)]; !ok {
					continue
				}
			} else {
				continue
			}
		}

		if cpu > 0 && formatInt(object["CpuCoreCount"]) != cpu {
			continue
		}

		if mem > 0 && formatFloat64(object["MemorySize"]) != mem {
			continue
		}
		if eniAmount > formatInt(object["EniQuantity"]) {
			continue
		}
		if gpuAmount > 0 && formatInt(object["GPUAmount"]) != gpuAmount {
			continue
		}
		// Kubernetes node does not support instance types which family is "ecs.t5" and spec less that c2g4
		// Kubernetes master node does not support gpu instance types which family prefixes with "ecs.gn"
		if k8sNode != "" {
			if object["InstanceTypeFamily"].(string) == "ecs.t5" {
				continue
			}
			if formatInt(object["CpuCoreCount"]) < 2 || formatFloat64(object["MemorySize"]) < 4 {
				continue
			}
			if k8sNode == string(KubernetesNodeMaster) && strings.HasPrefix(object["InstanceTypeFamily"].(string), "ecs.gn") {
				continue
			}
		}

		instanceTypes = append(instanceTypes, instanceTypeWithOriginalPrice{
			InstanceType: object,
		})
	}
	sortedBy := d.Get("sorted_by").(string)

	if sortedBy == "Price" && len(instanceTypes) > 0 {
		bssopenapiService := BssOpenApiService{client}
		instanceChargeType := d.Get("instance_charge_type").(string)
		moduleCode := "InstanceType"
		modules := make([]map[string]interface{}, 0)
		for _, types := range instanceTypes {
			config := fmt.Sprintf("InstanceType:%s,IoOptimized:IoOptimized,ImageOs:linux,Region:%s",
				fmt.Sprint(types.InstanceType["InstanceTypeId"]), client.RegionId)
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

		priceList, err := bssopenapiService.GetInstanceTypePrice("ecs", "", paymentType, modules)
		if err != nil {
			return WrapError(err)
		}
		for i := 0; i < len(instanceTypes); i++ {
			instanceTypes[i].OriginalPrice = priceList[i]
		}
	}

	return instanceTypesDescriptionAttributes(d, instanceTypes, mapInstanceTypes)
}

func instanceTypesDescriptionAttributes(d *schema.ResourceData, types []instanceTypeWithOriginalPrice, mapTypes map[string][]string) error {
	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "Price":
				return types[i].OriginalPrice < types[j].OriginalPrice
			case "CPU":
				return formatInt(types[i].InstanceType["CpuCoreCount"]) < formatInt(types[j].InstanceType["CpuCoreCount"])
			case "Memory":
				return formatFloat64(types[i].InstanceType["MemorySize"]) < formatFloat64(types[j].InstanceType["MemorySize"])
			}
			return false
		})
	}

	var ids []string
	var s []map[string]interface{}
	for _, t := range types {
		mapping := map[string]interface{}{
			"id":                              t.InstanceType["InstanceTypeId"],
			"cpu_core_count":                  formatInt(t.InstanceType["CpuCoreCount"]),
			"memory_size":                     formatFloat64(t.InstanceType["MemorySize"]),
			"family":                          t.InstanceType["InstanceTypeFamily"],
			"eni_amount":                      t.InstanceType["EniQuantity"],
			"nvme_support":                    t.InstanceType["NvmeSupport"],
			"eni_quantity":                    t.InstanceType["EniQuantity"],
			"primary_eni_queue_number":        t.InstanceType["PrimaryEniQueueNumber"],
			"secondary_eni_queue_number":      t.InstanceType["SecondaryEniQueueNumber"],
			"eni_ipv6_address_quantity":       t.InstanceType["EniIpv6AddressQuantity"],
			"maximum_queue_number_per_eni":    t.InstanceType["MaximumQueueNumberPerEni"],
			"total_eni_queue_quantity":        t.InstanceType["TotalEniQueueQuantity"],
			"eni_private_ip_address_quantity": t.InstanceType["EniPrivateIpAddressQuantity"],
		}
		if sortedBy == "Price" {
			mapping["price"] = fmt.Sprintf("%.4f", t.OriginalPrice)
		}
		zoneIds := mapTypes[t.InstanceType["InstanceTypeId"].(string)]
		sort.Strings(zoneIds)
		mapping["availability_zones"] = zoneIds
		gpu := map[string]interface{}{
			"amount":   strconv.Itoa(formatInt(t.InstanceType["GPUAmount"])),
			"category": t.InstanceType["GPUSpec"],
		}
		mapping["gpu"] = gpu
		brust := map[string]interface{}{
			"initial_credit":  strconv.Itoa(formatInt(t.InstanceType["InitialCredit"])),
			"baseline_credit": strconv.Itoa(formatInt(t.InstanceType["BaselineCredit"])),
		}
		mapping["burstable_instance"] = brust
		local := map[string]interface{}{
			"amount":   strconv.Itoa(formatInt(t.InstanceType["LocalStorageAmount"])),
			"category": t.InstanceType["LocalStorageCategory"],
		}
		if v, ok := t.InstanceType["LocalStorageCapacity"]; ok {
			local["capacity"] = v.(json.Number).String()
		} else {
			local["capacity"] = "0"
		}
		mapping["local_storage"] = local

		ids = append(ids, fmt.Sprint(t.InstanceType["InstanceTypeId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instance_types", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

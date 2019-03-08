package alicloud

import (
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

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
				ValidateFunc: validateInstanceType,
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
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateInstanceChargeType,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Vpc), string(Classic)}),
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      NoSpot,
				ValidateFunc: validateInstanceSpotStrategy,
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
				ValidateFunc: validateAllowedStringValue([]string{
					string(KubernetesNodeMaster),
					string(KubernetesNodeWorker),
				}),
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
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

	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}

	mapInstanceTypes := make(map[string][]string)
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}

					zones, _ := mapInstanceTypes[t.Value]
					zones = append(zones, zone.ZoneId)
					mapInstanceTypes[t.Value] = zones
				}
			}
		}
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := d.Get("memory_size").(float64)
	family := strings.TrimSpace(d.Get("instance_type_family").(string))

	req := ecs.CreateDescribeInstanceTypesRequest()
	req.InstanceTypeFamily = family

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceTypes(req)
	})
	if err != nil {
		return err
	}
	var instanceTypes []ecs.InstanceType
	resp, _ := raw.(*ecs.DescribeInstanceTypesResponse)
	if resp != nil {

		eniAmount := d.Get("eni_amount").(int)
		k8sNode := strings.TrimSpace(d.Get("kubernetes_node_role").(string))
		for _, types := range resp.InstanceTypes.InstanceType {
			if _, ok := mapInstanceTypes[types.InstanceTypeId]; !ok {
				continue
			}

			if cpu > 0 && types.CpuCoreCount != cpu {
				continue
			}

			if mem > 0 && types.MemorySize != mem {
				continue
			}
			if eniAmount > types.EniQuantity {
				continue
			}
			// Kubernetes node does not support instance types which family is "ecs.t5" and spec less that 2c4g
			// Kubernetes master node does not support gpu instance types which family prefixes with "ecs.gn"
			if k8sNode != "" {
				if types.InstanceTypeFamily == "ecs.t5" {
					continue
				}
				if types.CpuCoreCount < 2 || types.MemorySize < 4 {
					continue
				}
				if k8sNode == string(KubernetesNodeMaster) && strings.HasPrefix(types.InstanceTypeFamily, "ecs.gn") {
					continue
				}
			}
			instanceTypes = append(instanceTypes, types)
		}
	}

	return instanceTypesDescriptionAttributes(d, instanceTypes, mapInstanceTypes)
}

func instanceTypesDescriptionAttributes(d *schema.ResourceData, types []ecs.InstanceType, mapTypes map[string][]string) error {
	var ids []string
	var s []map[string]interface{}
	for _, t := range types {
		mapping := map[string]interface{}{
			"id":             t.InstanceTypeId,
			"cpu_core_count": t.CpuCoreCount,
			"memory_size":    t.MemorySize,
			"family":         t.InstanceTypeFamily,
			"eni_amount":     t.EniQuantity,
		}
		zoneIds := mapTypes[t.InstanceTypeId]
		sort.Strings(zoneIds)
		mapping["availability_zones"] = zoneIds
		gpu := map[string]interface{}{
			"amount":   strconv.Itoa(t.GPUAmount),
			"category": t.GPUSpec,
		}
		mapping["gpu"] = gpu
		brust := map[string]interface{}{
			"initial_credit":  strconv.Itoa(t.InitialCredit),
			"baseline_credit": strconv.Itoa(t.BaselineCredit),
		}
		mapping["burstable_instance"] = brust
		local := map[string]interface{}{
			"capacity": strconv.Itoa(t.LocalStorageCapacity),
			"amount":   strconv.Itoa(t.LocalStorageAmount),
			"category": t.LocalStorageCategory,
		}
		mapping["local_storage"] = local

		ids = append(ids, t.InstanceTypeId)
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

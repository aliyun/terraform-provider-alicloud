package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
)

func dataSourceAlicloudKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},

			"finger_print": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//Computed value
			"key_pairs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Resource{Schema: outputInstancesSchema()},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	var regex *regexp.Regexp
	if name, ok := d.GetOk("name_regex"); ok {
		regex = regexp.MustCompile(name.(string))
	}

	args := &ecs.DescribeKeyPairsArgs{
		RegionId: getRegion(d, meta),
	}
	if fingerPrint, ok := d.GetOk("finger_print"); ok {
		args.KeyPairFingerPrint = fingerPrint.(string)
	}
	var keyPairs []ecs.KeyPairItemType
	pagination := getPagination(1, 50)
	for true {
		args.Pagination = pagination
		results, _, err := conn.DescribeKeyPairs(args)
		if err != nil {
			return fmt.Errorf("Error DescribekeyPairs: %#v", err)
		}
		for _, key := range results {
			if regex == nil || (regex != nil && regex.MatchString(key.KeyPairName)) {
				keyPairs = append(keyPairs, key)
			}
		}
		if len(results) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	if len(keyPairs) < 1 {
		return fmt.Errorf("Your query key pairs returned no results. Please change your search criteria and try again.")
	}

	keyPairsAttach := make(map[string][]map[string]interface{})
	pagination.PageNumber = 1
	for true {
		instances, _, err := conn.DescribeInstances(&ecs.DescribeInstancesArgs{
			RegionId:   getRegion(d, meta),
			Pagination: pagination,
		})
		if err != nil {
			return fmt.Errorf("Error DescribeInstances: %#v", err)
		}
		for _, inst := range instances {
			if inst.KeyPairName != "" {
				public_ip := inst.EipAddress.IpAddress
				if public_ip == "" && len(inst.PublicIpAddress.IpAddress) > 0 {
					public_ip = inst.PublicIpAddress.IpAddress[0]
				}
				var private_ip string
				if len(inst.InnerIpAddress.IpAddress) > 0 {
					private_ip = inst.InnerIpAddress.IpAddress[0]
				} else if len(inst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
					private_ip = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
				}
				mapping := map[string]interface{}{
					"availability_zone": inst.ZoneId,
					"instance_id":       inst.InstanceId,
					"instance_name":     inst.InstanceName,
					"vswitch_id":        inst.VpcAttributes.VSwitchId,
					"public_ip":         public_ip,
					"private_ip":        private_ip,
				}
				if val, ok := keyPairsAttach[inst.KeyPairName]; ok {
					val = append(val, mapping)
					keyPairsAttach[inst.KeyPairName] = val
				} else {
					keyPairsAttach[inst.KeyPairName] = append(make([]map[string]interface{}, 0, 1), mapping)
				}
			}
		}
		if len(instances) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	return keyPairsDescriptionAttributes(d, keyPairs, keyPairsAttach)
}

func keyPairsDescriptionAttributes(d *schema.ResourceData, keyPairs []ecs.KeyPairItemType, keyPairsAttach map[string][]map[string]interface{}) error {
	var names []string
	var s []map[string]interface{}
	for _, key := range keyPairs {
		mapping := map[string]interface{}{
			"id":           key.KeyPairName,
			"key_name":     key.KeyPairName,
			"finger_print": key.KeyPairFingerPrint,
			"instances":    keyPairsAttach[key.KeyPairName],
		}

		log.Printf("[DEBUG] alicloud_key_pairs - adding keypair mapping: %v", mapping)
		names = append(names, string(key.KeyPairName))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("key_pairs", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

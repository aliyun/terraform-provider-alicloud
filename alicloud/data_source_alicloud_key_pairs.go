package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	client := meta.(*connectivity.AliyunClient)

	var regex *regexp.Regexp
	if name, ok := d.GetOk("name_regex"); ok {
		regex = regexp.MustCompile(name.(string))
	}

	args := ecs.CreateDescribeKeyPairsRequest()
	if fingerPrint, ok := d.GetOk("finger_print"); ok {
		args.KeyPairFingerPrint = fingerPrint.(string)
	}
	args.PageNumber = requests.NewInteger(1)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	var keyPairs []ecs.KeyPair
	keyPairsAttach := make(map[string][]map[string]interface{})

	for true {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeKeyPairs(args)
		})
		if err != nil {
			return fmt.Errorf("Error DescribekeyPairs: %#v", err)
		}
		results, _ := raw.(*ecs.DescribeKeyPairsResponse)
		if results == nil || len(results.KeyPairs.KeyPair) < 1 {
			break
		}
		for _, key := range results.KeyPairs.KeyPair {
			if regex == nil || (regex != nil && regex.MatchString(key.KeyPairName)) {
				keyPairs = append(keyPairs, key)
				keyPairsAttach[key.KeyPairName] = make([]map[string]interface{}, 0)
			}
		}
		if len(results.KeyPairs.KeyPair) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	req := ecs.CreateDescribeInstancesRequest()
	req.PageNumber = requests.NewInteger(1)
	req.PageSize = requests.NewInteger(PageSizeLarge)

	for true {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error DescribeInstances: %#v", err)
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		for _, inst := range resp.Instances.Instance {
			if _, ok := keyPairsAttach[inst.KeyPairName]; ok {
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
		if len(resp.Instances.Instance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	return keyPairsDescriptionAttributes(d, keyPairs, keyPairsAttach)
}

func keyPairsDescriptionAttributes(d *schema.ResourceData, keyPairs []ecs.KeyPair, keyPairsAttach map[string][]map[string]interface{}) error {
	var names []string
	var s []map[string]interface{}
	for _, key := range keyPairs {
		mapping := map[string]interface{}{
			"id":           key.KeyPairName,
			"key_name":     key.KeyPairName,
			"finger_print": key.KeyPairFingerPrint,
			"instances":    keyPairsAttach[key.KeyPairName],
		}

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

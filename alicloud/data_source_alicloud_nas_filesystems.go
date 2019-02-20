package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFileSystemsRead,

		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},


			// Computed values
			"filesystems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destription": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metered_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mounttarget_domain": {
							Type:     schema.TypeString,
							Computed: true,
							//Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := nas.CreateDescribeFileSystemsRequest()
	args.RegionId = string(client.Region)
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allfss []nas.FileSystem
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeFileSystems(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*nas.DescribeFileSystemsResponse)
		if resp == nil || len(resp.FileSystems.FileSystem) < 1 {
			break
		}

		for _, e := range resp.FileSystems.FileSystem {
			if len(idsMap) > 0 {
				if _, ok := idsMap[e.FileSystemId]; !ok {
					continue
				}
			}
			if storage_type, ok := d.GetOk("storage_type"); ok && e.StorageType != storage_type.(string) {
				continue
			}

			if protocol_type, ok := d.GetOk("protocol_type"); ok && string(e.ProtocolType) != protocol_type.(string) {
				continue
			}

			allfss = append(allfss, e)
		}

		if len(resp.FileSystems.FileSystem) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return err
		} else {
			args.PageNumber = page
		}
	}

	var filteredFss []nas.FileSystem
	var mount_target_domain []string

	for _, v := range allfss {
		if protocol_type, ok := d.GetOk("protocol_type"); ok && v.ProtocolType != protocol_type.(string) {
			continue
		}

		if storage_type, ok := d.GetOk("storage_type"); ok && string(v.StorageType) != storage_type.(string) {
			continue
		}

		request := nas.CreateDescribeMountTargetsRequest()
		request.FileSystemId = v.FileSystemId
		request.RegionId = string(client.Region)

		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeMountTargets(request)
		})
		if err != nil {
			return fmt.Errorf("Error DescribeMountTargets by filesystem %s: %#v", v.MountTargets.MountTarget[0].MountTargetDomain, err)
		}
		vrs, _ := raw.(*nas.DescribeMountTargetsResponse)
		if vrs != nil && len(vrs.MountTargets.MountTarget) > 0 {
			mount_target_domain = append(mount_target_domain, vrs.MountTargets.MountTarget[0].MountTargetDomain)
		} else {
			mount_target_domain = append(mount_target_domain, "")
		}
		filteredFss = append(filteredFss, v)
	}

	return fileSystemsDecriptionAttributes(d, filteredFss,mount_target_domain, meta)
}


func fileSystemsDecriptionAttributes(d *schema.ResourceData, fssSetTypes []nas.FileSystem,mount_target_domain []string, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for index, fs := range fssSetTypes {
		mapping := map[string]interface{}{
			"id":             		fs.FileSystemId,
			"region_id":      		fs.RegionId,
			"create_time":    		fs.CreateTime,
			"destription":    		fs.Destription,
			"protocol_type":  		fs.ProtocolType,
			"storage_type":   		fs.StorageType,
			"metered_size":   		fs.MeteredSize,
			"mounttarget_domain":   mount_target_domain[index],
		}
		ids = append(ids, fs.FileSystemId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("filesystems", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

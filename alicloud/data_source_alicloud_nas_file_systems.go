package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"systems": {
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
						"description": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := nas.CreateDescribeFileSystemsRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allfss []nas.FileSystem
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
				return nasClient.DescribeFileSystems(request)
			})
			raw = rsp
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_file_systems", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		destription, ok := d.GetOk("description_regex")
		var r *regexp.Regexp
		if ok && destription.(string) != "" {
			r = regexp.MustCompile(destription.(string))
		}
		response, _ := raw.(*nas.DescribeFileSystemsResponse)
		if len(response.FileSystems.FileSystem) < 1 {
			break
		}
		for _, file_system := range response.FileSystems.FileSystem {
			if v, ok := d.GetOk("storage_type"); ok && file_system.StorageType != Trim(v.(string)) {
				continue
			}
			if v, ok := d.GetOk("protocol_type"); ok && string(file_system.ProtocolType) != Trim(v.(string)) {
				continue
			}
			if r != nil && !r.MatchString(file_system.Description) {
				continue
			}
			if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
				id_found := false
				for _, id := range v.([]interface{}) {
					if id == nil {
						continue
					}
					if string(file_system.FileSystemId) == id.(string) {
						id_found = true
						break
					}
				}
				if !id_found {
					continue
				}
			}
			allfss = append(allfss, file_system)
		}

		if len(response.FileSystems.FileSystem) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	return fileSystemsDecriptionAttributes(d, allfss, meta)
}

func fileSystemsDecriptionAttributes(d *schema.ResourceData, fssSetTypes []nas.FileSystem, meta interface{}) error {
	var ids []string
	var descriptions []string
	var s []map[string]interface{}
	for _, fs := range fssSetTypes {
		mapping := map[string]interface{}{
			"id":            fs.FileSystemId,
			"region_id":     fs.RegionId,
			"create_time":   fs.CreateTime,
			"description":   fs.Description,
			"protocol_type": fs.ProtocolType,
			"storage_type":  fs.StorageType,
			"metered_size":  fs.MeteredSize,
		}
		ids = append(ids, fs.FileSystemId)
		descriptions = append(descriptions, fs.Description)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("systems", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("descriptions", descriptions); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

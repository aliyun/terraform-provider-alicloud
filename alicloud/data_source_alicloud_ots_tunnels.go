package alicloud

import (
	"fmt"
	"regexp"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOtsTunnels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsTunnelsRead,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSTableName,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tunnels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_rpo": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tunnel_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_stage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"channels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"channel_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"channel_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"channel_rpo": {
										Type:     schema.TypeInt,
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

type OtsTunnelInfo struct {
	instanceName string
	tableName    string
	tunnelName   string
	tunnelId     string
	tunnelType   string
	tunnelRpo    int64
	tunnelStage  string
	expired      bool
	createTime   int64
	channels     []*otsTunnel.ChannelInfo
}

func dataSourceAlicloudOtsTunnelsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	instanceName := d.Get("instance_name").(string)
	tableName := d.Get("table_name").(string)

	object, err := otsService.ListOtsTunnels(instanceName, tableName)
	if err != nil {
		return WrapError(err)
	}

	idsMap := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, x := range v.([]interface{}) {
			if x == nil {
				continue
			}
			idsMap[x.(string)] = true
		}
	}

	var nameReg *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		nameReg = regexp.MustCompile(v.(string))
	}

	var filteredTunnelNames []string
	for _, tunnel := range object.Tunnels {
		// name_regex mismatch
		if nameReg != nil && !nameReg.MatchString(tunnel.TunnelName) {
			continue
		}
		// ids mismatch
		if len(idsMap) != 0 {
			id := fmt.Sprintf("%s%s%s%s%s", instanceName, COLON_SEPARATED, tableName, COLON_SEPARATED, tunnel.TunnelName)
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}
		filteredTunnelNames = append(filteredTunnelNames, tunnel.TunnelName)
	}

	// get full tunnelInfo via DescribeTunnel
	allTunnelInfos := make([]*OtsTunnelInfo, 0, len(filteredTunnelNames))
	for _, tunnelName := range filteredTunnelNames {
		id := fmt.Sprintf("%s%s%s%s%s", instanceName, COLON_SEPARATED, tableName, COLON_SEPARATED, tunnelName)
		object, err := otsService.DescribeOtsTunnel(id)
		if err != nil {
			return WrapError(err)
		}
		allTunnelInfos = append(allTunnelInfos, &OtsTunnelInfo{
			instanceName: instanceName,
			tableName:    tableName,
			tunnelName:   object.Tunnel.TunnelName,
			tunnelId:     object.Tunnel.TunnelId,
			tunnelType:   object.Tunnel.TunnelType,
			tunnelRpo:    object.TunnelRPO,
			tunnelStage:  object.Tunnel.Stage,
			expired:      object.Tunnel.Expired,
			createTime:   object.Tunnel.CreateTime.UnixNano(),
			channels:     object.Channels,
		})
	}

	return otsTunnelDescriptionAttributes(d, allTunnelInfos, meta)
}

func otsTunnelDescriptionAttributes(d *schema.ResourceData, tunnelInfos []*OtsTunnelInfo, meta interface{}) error {
	ids := make([]string, 0, len(tunnelInfos))
	names := make([]string, 0, len(tunnelInfos))
	s := make([]map[string]interface{}, 0, len(tunnelInfos))
	for _, tunnel := range tunnelInfos {
		id := fmt.Sprintf("%s%s%s%s%s", tunnel.instanceName, COLON_SEPARATED, tunnel.tableName, COLON_SEPARATED, tunnel.tunnelName)
		mapping := map[string]interface{}{
			"id":            id,
			"instance_name": tunnel.instanceName,
			"table_name":    tunnel.tableName,
			"tunnel_name":   tunnel.tunnelName,
			"tunnel_id":     tunnel.tunnelId,
			"tunnel_rpo":    tunnel.tunnelRpo,
			"tunnel_type":   tunnel.tunnelType,
			"tunnel_stage":  tunnel.tunnelStage,
			"expired":       tunnel.expired,
			"create_time":   tunnel.createTime,
		}

		channels := make([]map[string]interface{}, 0, len(tunnel.channels))
		for _, c := range tunnel.channels {
			channel := make(map[string]interface{})
			channel["channel_id"] = c.ChannelId
			channel["channel_type"] = c.ChannelType
			channel["channel_status"] = c.ChannelStatus
			channel["client_id"] = c.ClientId
			channel["channel_rpo"] = c.ChannelRPO
			channels = append(channels, channel)
		}
		mapping["channels"] = channels

		names = append(names, tunnel.tunnelName)
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("tunnels", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

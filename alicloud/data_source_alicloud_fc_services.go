package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudFcServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcServicesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"logstore": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
						"vpc_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"vswitch_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFcServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var serviceMappings []map[string]interface{}
	nextToken := ""
	for {
		args := fc.NewListServicesInput()
		if nextToken != "" {
			args.NextToken = &nextToken
		}

		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.ListServices(args)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*fc.ListServicesOutput)

		if resp.Services == nil || len(resp.Services) < 1 {
			break
		}

		for _, service := range resp.Services {
			mapping := map[string]interface{}{
				"id":                     *service.ServiceID,
				"name":                   *service.ServiceName,
				"description":            *service.Description,
				"role":                   *service.Role,
				"internet_access":        *service.InternetAccess,
				"creation_time":          *service.CreatedTime,
				"last_modification_time": *service.LastModifiedTime,
			}

			var logConfigMappings []map[string]interface{}
			if service.LogConfig != nil {
				logConfigMappings = append(logConfigMappings, map[string]interface{}{
					"project":  *service.LogConfig.Project,
					"logstore": *service.LogConfig.Logstore,
				})
			}
			mapping["log_config"] = logConfigMappings

			var vpcConfigMappings []map[string]interface{}
			if service.VPCConfig != nil &&
				(service.VPCConfig.VPCID != nil || service.VPCConfig.SecurityGroupID != nil) {
				vpcConfigMappings = append(vpcConfigMappings, map[string]interface{}{
					"vpc_id":            *service.VPCConfig.VPCID,
					"vswitch_ids":       service.VPCConfig.VSwitchIDs,
					"security_group_id": *service.VPCConfig.SecurityGroupID,
				})
			}
			mapping["vpc_config"] = vpcConfigMappings

			serviceMappings = append(serviceMappings, mapping)

			log.Printf("[DEBUG] alicloud_fc_services - adding service mapping: %v", mapping)
			ids = append(ids, *service.ServiceID)
		}

		nextToken = ""
		if resp.NextToken != nil {
			nextToken = *resp.NextToken
		}
		if nextToken == "" {
			break
		}
	}

	var filteredServiceMappings []map[string]interface{}
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, mapping := range serviceMappings {
			if r != nil && !r.MatchString(mapping["name"].(string)) {
				continue
			}
			filteredServiceMappings = append(filteredServiceMappings, mapping)
		}
	} else {
		filteredServiceMappings = serviceMappings
	}

	if len(filteredServiceMappings) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	log.Printf("[DEBUG] alicloud_fc_services - Services found: %#v", filteredServiceMappings)

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("services", filteredServiceMappings); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), filteredServiceMappings)
	}
	return nil
}

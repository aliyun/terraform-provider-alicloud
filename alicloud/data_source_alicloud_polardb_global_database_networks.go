package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPolarDBGlobalDatabaseNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBGlobalDatabaseNetworksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gdn_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"creating", "active", "deleting", "locked", "removing_member"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gdn_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
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

func dataSourceAlicloudPolarDBGlobalDatabaseNetworksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGlobalDatabaseNetworks"
	request := make(map[string]interface{})
	setPagingRequest(d, request, PageSizeLarge)
	if v, ok := d.GetOk("gdn_id"); ok {
		request["GDNId"] = v
	}
	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["GDNDescription"] = v
	}
	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")

	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_global_database_networks", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GDNId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["GDNStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["GDNId"]),
			"gdn_id":      fmt.Sprint(object["GDNId"]),
			"description": object["GDNDescription"],
			"db_type":     object["DBType"],
			"db_version":  object["DBVersion"],
			"create_time": object["CreateTime"],
			"status":      object["GDNStatus"],
		}
		if v, ok := object["DBClusters"]; ok {
			dbClustersMaps := make([]map[string]interface{}, 0)
			dbClustersMap := make(map[string]interface{})
			for _, dbClusters := range v.([]interface{}) {
				dbClustersArg := dbClusters.(map[string]interface{})
				dbClustersMap["db_cluster_id"] = dbClustersArg["DBClusterId"]
				dbClustersMap["role"] = dbClustersArg["Role"]
				dbClustersMap["region_id"] = dbClustersArg["RegionId"]
				dbClustersMaps = append(dbClustersMaps, dbClustersMap)
			}
			mapping["db_clusters"] = dbClustersMaps
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("networks", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

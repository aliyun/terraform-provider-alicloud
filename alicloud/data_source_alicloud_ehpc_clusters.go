package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEhpcClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEhpcClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"creating", "exception", "initing", "releasing", "running", "stopped", "uninit"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compute_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"compute_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_owner_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"login_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manager_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manager_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"os_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scc_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_mountpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"post_install_script": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"args": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudEhpcClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListClusters"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var clusterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		clusterNameRegex = r
	}

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
	conn, err := client.NewEhsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ehpc_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Clusters.ClusterInfoSimple", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Clusters.ClusterInfoSimple", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if clusterNameRegex != nil && !clusterNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_type":      object["AccountType"],
			"client_version":    object["ClientVersion"],
			"id":                fmt.Sprint(object["Id"]),
			"cluster_id":        fmt.Sprint(object["Id"]),
			"cluster_name":      object["Name"],
			"create_time":       object["CreateTime"],
			"deploy_mode":       object["DeployMode"],
			"description":       object["Description"],
			"image_id":          object["ImageId"],
			"image_owner_alias": object["ImageOwnerAlias"],
			"os_tag":            object["OsTag"],
			"scheduler_type":    object["SchedulerType"],
			"status":            object["Status"],
			"vswitch_id":        object["VSwitchId"],
			"vpc_id":            object["VpcId"],
			"zone_id":           object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["Id"])
		ehpcService := EhpcService{client}
		getResp, err := ehpcService.DescribeEhpcCluster(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["ha_enable"] = getResp["HaEnable"]

		if ecsInfo, ok := getResp["EcsInfo"]; ok {
			if compute, ok := ecsInfo.(map[string]interface{})["Compute"]; ok {
				mapping["compute_count"] = compute.(map[string]interface{})["Count"]
				mapping["compute_instance_type"] = compute.(map[string]interface{})["InstanceType"]
			}
			if login, ok := ecsInfo.(map[string]interface{})["Login"]; ok {
				mapping["login_count"] = login.(map[string]interface{})["Count"]
				mapping["login_instance_type"] = login.(map[string]interface{})["InstanceType"]
			}
			if manager, ok := ecsInfo.(map[string]interface{})["Manager"]; ok {
				mapping["manager_instance_type"] = manager.(map[string]interface{})["InstanceType"]
				mapping["manager_count"] = manager.(map[string]interface{})["Count"]
			}
		}

		mapping["remote_directory"] = getResp["RemoteDirectory"]
		mapping["scc_cluster_id"] = getResp["SccClusterId"]
		mapping["security_group_id"] = getResp["SecurityGroupId"]
		mapping["volume_id"] = getResp["VolumeId"]
		mapping["volume_mountpoint"] = getResp["VolumeMountpoint"]
		mapping["volume_protocol"] = getResp["VolumeProtocol"]
		mapping["volume_type"] = getResp["VolumeType"]

		postInstallScriptsList := make([]map[string]interface{}, 0)
		if postInstallScripts, ok := object["PostInstallScripts"]; ok {
			if postInstallScriptInfo, ok := postInstallScripts.(map[string]interface{})["PostInstallScriptInfo"]; ok {
				for _, v := range postInstallScriptInfo.([]interface{}) {
					postInstallScriptInfoArg := v.(map[string]interface{})
					postInstallScriptsList = append(postInstallScriptsList, map[string]interface{}{
						"url":  postInstallScriptInfoArg["Url"],
						"args": postInstallScriptInfoArg["Args"],
					})
				}
				mapping["post_install_script"] = postInstallScriptsList
			}
		}
		applicationList := make([]map[string]interface{}, 0)
		if applications, ok := object["Applications"]; ok {
			if applicationInfo, ok := applications.(map[string]interface{})["ApplicationInfo"]; ok {
				for _, v := range applicationInfo.([]interface{}) {
					applicationArg := v.(map[string]interface{})
					applicationList = append(applicationList, map[string]interface{}{
						"tag": applicationArg["Tag"],
					})
				}
				mapping["application"] = applicationList
			}
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

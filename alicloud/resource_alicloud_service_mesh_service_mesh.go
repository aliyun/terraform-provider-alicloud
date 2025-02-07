// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudServiceMeshServiceMesh() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudServiceMeshServiceMeshCreate,
		Read:   resourceAliCloudServiceMeshServiceMeshRead,
		Update: resourceAliCloudServiceMeshServiceMeshUpdate,
		Delete: resourceAliCloudServiceMeshServiceMeshDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"standard", "enterprise", "ultimate"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customized_prometheus": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"edition": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Default", "Pro"}, false),
			},
			"extra_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cr_aggregation_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_loadbalancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pilot_public_eip": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"pilot_public_loadbalancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server_public_eip": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"pilot_public_eip_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"mesh_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_log": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gateway_lifecycle": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(0, 365),
									},
									"gateway_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"sidecar_lifecycle": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(0, 365),
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"sidecar_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"pilot": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http10_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"trace_sampling": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
						"opa": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"request_memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"limit_cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"log_level": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"request_cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"prometheus": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use_external": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"external_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"telemetry": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"outbound_traffic_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"ALLOW_ANY", "REGISTRY_ONLY"}, false),
						},
						"sidecar_injector": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"auto_injection_policy_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"request_memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable_namespaces_by_default": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"limit_cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"init_cni_configuration": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"exclude_namespaces": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"request_cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sidecar_injector_webhook_as_yaml": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"audit": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"kiali": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aggregated_kiali_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"integrate_clb": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"auth_strategy": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"token", "openid", "ramoauth"}, false),
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"kiali_arms_auth_tokens": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"open_id_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"issuer_uri": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"client_secret": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"scopes": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"client_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"distributed_kiali_addresses": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"use_populated_arms_prometheus": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"server_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"web_schema": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringInSlice([]string{"http", "https"}, false),
												},
												"web_root": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"web_fqdn": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"web_port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"kiali_service_annotations": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"distributed_kiali_access_tokens": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ram_oauth_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"redirect_uris": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"custom_prometheus_url": {
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
						"proxy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"request_memory": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cluster_domain": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"limit_cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"request_cpu": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"include_ip_ranges": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_locality_lb": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"control_plane_log": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringMatch(regexp.MustCompile("^[\\w.-]+$"), "The name of the SLS Project to which the control plane logs are collected."),
									},
									"log_ttl_in_day": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: IntBetween(0, 365),
									},
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"tracing": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"customized_zipkin": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitche_list": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"prometheus_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_mesh_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudServiceMeshServiceMeshCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateServiceMesh"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("service_mesh_name"); ok {
		request["Name"] = v
	}
	jsonPathResult1, err := jsonpath.Get("$[0].vpc_id", d.Get("network"))
	if err == nil {
		request["VpcId"] = jsonPathResult1
	}

	if v, ok := d.GetOk("load_balancer"); ok {
		jsonPathResult2, err := jsonpath.Get("$[0].api_server_public_eip", v)
		if err == nil && jsonPathResult2 != "" {
			request["ApiServerPublicEip"] = jsonPathResult2
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult3, err := jsonpath.Get("$[0].tracing", v)
		if err == nil && jsonPathResult3 != "" {
			request["Tracing"] = jsonPathResult3
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult4, err := jsonpath.Get("$[0].pilot[0].trace_sampling", v)
		if err == nil && jsonPathResult4 != "" {
			request["TraceSampling"] = jsonPathResult4
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult5, err := jsonpath.Get("$[0].customized_zipkin", v)
		if err == nil && jsonPathResult5 != "" {
			request["CustomizedZipkin"] = jsonPathResult5
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult6, err := jsonpath.Get("$[0].telemetry", v)
		if err == nil && jsonPathResult6 != "" {
			request["Telemetry"] = jsonPathResult6
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult7, err := jsonpath.Get("$[0].include_ip_ranges", v)
		if err == nil && jsonPathResult7 != "" {
			request["IncludeIPRanges"] = jsonPathResult7
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult8, err := jsonpath.Get("$[0].opa[0].log_level", v)
		if err == nil && jsonPathResult8 != "" {
			request["OPALogLevel"] = jsonPathResult8
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult9, err := jsonpath.Get("$[0].opa[0].request_cpu", v)
		if err == nil && jsonPathResult9 != "" {
			request["OPARequestCPU"] = jsonPathResult9
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult10, err := jsonpath.Get("$[0].opa[0].limit_cpu", v)
		if err == nil && jsonPathResult10 != "" {
			request["OPALimitCPU"] = jsonPathResult10
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult11, err := jsonpath.Get("$[0].opa[0].limit_memory", v)
		if err == nil && jsonPathResult11 != "" {
			request["OPALimitMemory"] = jsonPathResult11
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult12, err := jsonpath.Get("$[0].opa[0].request_memory", v)
		if err == nil && jsonPathResult12 != "" {
			request["OPARequestMemory"] = jsonPathResult12
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult13, err := jsonpath.Get("$[0].proxy[0].request_cpu", v)
		if err == nil && jsonPathResult13 != "" {
			request["ProxyRequestCPU"] = jsonPathResult13
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult14, err := jsonpath.Get("$[0].proxy[0].limit_cpu", v)
		if err == nil && jsonPathResult14 != "" {
			request["ProxyLimitCPU"] = jsonPathResult14
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult15, err := jsonpath.Get("$[0].proxy[0].limit_memory", v)
		if err == nil && jsonPathResult15 != "" {
			request["ProxyLimitMemory"] = jsonPathResult15
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult16, err := jsonpath.Get("$[0].proxy[0].request_memory", v)
		if err == nil && jsonPathResult16 != "" {
			request["ProxyRequestMemory"] = jsonPathResult16
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult17, err := jsonpath.Get("$[0].kiali[0].enabled", v)
		if err == nil && jsonPathResult17 != "" {
			request["KialiEnabled"] = jsonPathResult17
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult18, err := jsonpath.Get("$[0].access_log[0].enabled", v)
		if err == nil && jsonPathResult18 != "" {
			request["AccessLogEnabled"] = jsonPathResult18
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult19, err := jsonpath.Get("$[0].enable_locality_lb", v)
		if err == nil && jsonPathResult19 != "" {
			request["LocalityLoadBalancing"] = jsonPathResult19
		}
	}
	if v, ok := d.GetOk("version"); ok {
		request["IstioVersion"] = v
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult21, err := jsonpath.Get("$[0].opa[0].enabled", v)
		if err == nil && jsonPathResult21 != "" {
			request["OpaEnabled"] = jsonPathResult21
		}
	}
	if v, ok := d.GetOk("edition"); ok {
		request["Edition"] = v
	}
	if v, ok := d.GetOk("cluster_spec"); ok {
		request["ClusterSpec"] = v
	}
	if v, ok := d.GetOkExists("customized_prometheus"); ok {
		request["CustomizedPrometheus"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult25, err := jsonpath.Get("$[0].control_plane_log[0].enabled", v)
		if err == nil && jsonPathResult25 != "" {
			request["ControlPlaneLogEnabled"] = jsonPathResult25
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult26, err := jsonpath.Get("$[0].control_plane_log[0].project", v)
		if err == nil && jsonPathResult26 != "" {
			request["ControlPlaneLogProject"] = jsonPathResult26
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult27, err := jsonpath.Get("$[0].audit[0].enabled", v)
		if err == nil && jsonPathResult27 != "" {
			request["EnableAudit"] = jsonPathResult27
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult28, err := jsonpath.Get("$[0].audit[0].project", v)
		if err == nil && jsonPathResult28 != "" {
			request["AuditProject"] = jsonPathResult28
		}
	}
	jsonPathResult29, err := jsonpath.Get("$[0].vswitche_list", d.Get("network"))
	if err == nil {
		request["VSwitches"] = convertListToJsonString(jsonPathResult29.([]interface{}))
	}

	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult30, err := jsonpath.Get("$[0].access_log[0].project", v)
		if err == nil && jsonPathResult30 != "" {
			request["AccessLogProject"] = jsonPathResult30
		}
	}
	if v, ok := d.GetOk("mesh_config"); ok {
		jsonPathResult31, err := jsonpath.Get("$[0].proxy[0].cluster_domain", v)
		if err == nil && jsonPathResult31 != "" {
			request["ClusterDomain"] = jsonPathResult31
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NameAlreadyExist", "InvalidActiveState.ACK", "ERR404"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_mesh_service_mesh", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceMeshId"]))

	serviceMeshServiceV2 := ServiceMeshServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudServiceMeshServiceMeshUpdate(d, meta)
}

func resourceAliCloudServiceMeshServiceMeshRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serviceMeshServiceV2 := ServiceMeshServiceV2{client}

	objectRaw, err := serviceMeshServiceV2.DescribeServiceMeshServiceMesh(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_service_mesh_service_mesh DescribeServiceMeshServiceMesh Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["ClusterSpec"] != nil {
		d.Set("cluster_spec", objectRaw["ClusterSpec"])
	}

	serviceMeshInfo1RawObj, _ := jsonpath.Get("$.ServiceMeshInfo", objectRaw)
	serviceMeshInfo1Raw := make(map[string]interface{})
	if serviceMeshInfo1RawObj != nil {
		serviceMeshInfo1Raw = serviceMeshInfo1RawObj.(map[string]interface{})
	}
	if serviceMeshInfo1Raw["CreationTime"] != nil {
		d.Set("create_time", serviceMeshInfo1Raw["CreationTime"])
	}
	if serviceMeshInfo1Raw["Profile"] != nil {
		d.Set("edition", serviceMeshInfo1Raw["Profile"])
	}
	if serviceMeshInfo1Raw["Name"] != nil {
		d.Set("service_mesh_name", serviceMeshInfo1Raw["Name"])
	}
	if serviceMeshInfo1Raw["State"] != nil {
		d.Set("status", serviceMeshInfo1Raw["State"])
	}
	if serviceMeshInfo1Raw["Version"] != nil {
		d.Set("version", serviceMeshInfo1Raw["Version"])
	}

	clusters1Raw := make([]interface{}, 0)
	if objectRaw["Clusters"] != nil {
		clusters1Raw = objectRaw["Clusters"].([]interface{})
	}

	d.Set("cluster_ids", clusters1Raw)

	extraConfigurationMaps := make([]map[string]interface{}, 0)
	extraConfigurationMap := make(map[string]interface{})
	cRAggregationConfiguration1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.ExtraConfiguration.CRAggregationConfiguration", objectRaw)
	cRAggregationConfiguration1Raw := make(map[string]interface{})
	if cRAggregationConfiguration1RawObj != nil {
		cRAggregationConfiguration1Raw = cRAggregationConfiguration1RawObj.(map[string]interface{})
	}
	if len(cRAggregationConfiguration1Raw) > 0 {
		extraConfigurationMap["cr_aggregation_enabled"] = cRAggregationConfiguration1Raw["Enabled"]

		extraConfigurationMaps = append(extraConfigurationMaps, extraConfigurationMap)
	}
	if cRAggregationConfiguration1RawObj != nil {
		if err := d.Set("extra_configuration", extraConfigurationMaps); err != nil {
			return err
		}
	}
	loadBalancerMaps := make([]map[string]interface{}, 0)
	loadBalancerMap := make(map[string]interface{})
	loadBalancer1RawObj, _ := jsonpath.Get("$.Spec.LoadBalancer", objectRaw)
	loadBalancer1Raw := make(map[string]interface{})
	if loadBalancer1RawObj != nil {
		loadBalancer1Raw = loadBalancer1RawObj.(map[string]interface{})
	}
	if len(loadBalancer1Raw) > 0 {
		loadBalancerMap["api_server_loadbalancer_id"] = loadBalancer1Raw["ApiServerLoadbalancerId"]
		loadBalancerMap["api_server_public_eip"] = loadBalancer1Raw["ApiServerPublicEip"]
		loadBalancerMap["pilot_public_eip"] = loadBalancer1Raw["PilotPublicEip"]
		loadBalancerMap["pilot_public_eip_id"] = loadBalancer1Raw["PilotPublicEipId"]
		loadBalancerMap["pilot_public_loadbalancer_id"] = loadBalancer1Raw["PilotPublicLoadbalancerId"]

		loadBalancerMaps = append(loadBalancerMaps, loadBalancerMap)
	}
	if loadBalancer1RawObj != nil {
		if err := d.Set("load_balancer", loadBalancerMaps); err != nil {
			return err
		}
	}
	meshConfigMaps := make([]map[string]interface{}, 0)
	meshConfigMap := make(map[string]interface{})
	meshConfig1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig", objectRaw)
	meshConfig1Raw := make(map[string]interface{})
	if meshConfig1RawObj != nil {
		meshConfig1Raw = meshConfig1RawObj.(map[string]interface{})
	}
	if len(meshConfig1Raw) > 0 {
		meshConfigMap["customized_zipkin"] = meshConfig1Raw["CustomizedZipkin"]
		meshConfigMap["enable_locality_lb"] = meshConfig1Raw["EnableLocalityLB"]
		meshConfigMap["include_ip_ranges"] = meshConfig1Raw["IncludeIPRanges"]
		meshConfigMap["outbound_traffic_policy"] = meshConfig1Raw["OutboundTrafficPolicy"]
		meshConfigMap["telemetry"] = meshConfig1Raw["Telemetry"]
		meshConfigMap["tracing"] = meshConfig1Raw["Tracing"]

		accessLogMaps := make([]map[string]interface{}, 0)
		accessLogMap := make(map[string]interface{})
		accessLog1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.AccessLog", objectRaw)
		accessLog1Raw := make(map[string]interface{})
		if accessLog1RawObj != nil {
			accessLog1Raw = accessLog1RawObj.(map[string]interface{})
		}
		if len(accessLog1Raw) > 0 {
			accessLogMap["enabled"] = accessLog1Raw["Enabled"]
			accessLogMap["project"] = accessLog1Raw["Project"]

			accessLogExtraConf1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.ExtraConfiguration.AccessLogExtraConf", objectRaw)
			accessLogExtraConf1Raw := make(map[string]interface{})
			if accessLogExtraConf1RawObj != nil {
				accessLogExtraConf1Raw = accessLogExtraConf1RawObj.(map[string]interface{})
			}
			if len(accessLogExtraConf1Raw) > 0 {
				accessLogMap["gateway_enabled"] = accessLogExtraConf1Raw["GatewayEnabled"]
				accessLogMap["gateway_lifecycle"] = accessLogExtraConf1Raw["GatewayLifecycle"]
				accessLogMap["sidecar_enabled"] = accessLogExtraConf1Raw["SidecarEnabled"]
				accessLogMap["sidecar_lifecycle"] = accessLogExtraConf1Raw["SidecarLifecycle"]
			}
			accessLogMaps = append(accessLogMaps, accessLogMap)
		}
		meshConfigMap["access_log"] = accessLogMaps
		auditMaps := make([]map[string]interface{}, 0)
		auditMap := make(map[string]interface{})
		audit1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.Audit", objectRaw)
		audit1Raw := make(map[string]interface{})
		if audit1RawObj != nil {
			audit1Raw = audit1RawObj.(map[string]interface{})
		}
		if len(audit1Raw) > 0 {
			auditMap["enabled"] = audit1Raw["Enabled"]
			auditMap["project"] = audit1Raw["Project"]

			auditMaps = append(auditMaps, auditMap)
		}
		meshConfigMap["audit"] = auditMaps
		controlPlaneLogMaps := make([]map[string]interface{}, 0)
		controlPlaneLogMap := make(map[string]interface{})
		controlPlaneLogInfo1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.ControlPlaneLogInfo", objectRaw)
		controlPlaneLogInfo1Raw := make(map[string]interface{})
		if controlPlaneLogInfo1RawObj != nil {
			controlPlaneLogInfo1Raw = controlPlaneLogInfo1RawObj.(map[string]interface{})
		}
		if len(controlPlaneLogInfo1Raw) > 0 {
			controlPlaneLogMap["enabled"] = controlPlaneLogInfo1Raw["Enabled"]
			controlPlaneLogMap["log_ttl_in_day"] = controlPlaneLogInfo1Raw["LogTTL"]
			controlPlaneLogMap["project"] = controlPlaneLogInfo1Raw["Project"]

			controlPlaneLogMaps = append(controlPlaneLogMaps, controlPlaneLogMap)
		}
		meshConfigMap["control_plane_log"] = controlPlaneLogMaps
		oPAMaps := make([]map[string]interface{}, 0)
		oPAMap := make(map[string]interface{})
		oPA1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.OPA", objectRaw)
		oPA1Raw := make(map[string]interface{})
		if oPA1RawObj != nil {
			oPA1Raw = oPA1RawObj.(map[string]interface{})
		}
		if len(oPA1Raw) > 0 {
			oPAMap["enabled"] = oPA1Raw["Enabled"]
			oPAMap["limit_cpu"] = oPA1Raw["LimitCPU"]
			oPAMap["limit_memory"] = oPA1Raw["LimitMemory"]
			oPAMap["log_level"] = oPA1Raw["LogLevel"]
			oPAMap["request_cpu"] = oPA1Raw["RequestCPU"]
			oPAMap["request_memory"] = oPA1Raw["RequestMemory"]

			oPAMaps = append(oPAMaps, oPAMap)
		}
		meshConfigMap["opa"] = oPAMaps
		pilotMaps := make([]map[string]interface{}, 0)
		pilotMap := make(map[string]interface{})
		pilot1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.Pilot", objectRaw)
		pilot1Raw := make(map[string]interface{})
		if pilot1RawObj != nil {
			pilot1Raw = pilot1RawObj.(map[string]interface{})
		}
		if len(pilot1Raw) > 0 {
			pilotMap["http10_enabled"] = pilot1Raw["Http10Enabled"]
			pilotMap["trace_sampling"] = pilot1Raw["TraceSampling"]

			pilotMaps = append(pilotMaps, pilotMap)
		}
		meshConfigMap["pilot"] = pilotMaps
		prometheusMaps := make([]map[string]interface{}, 0)
		prometheusMap := make(map[string]interface{})
		prometheus1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.Prometheus", objectRaw)
		prometheus1Raw := make(map[string]interface{})
		if prometheus1RawObj != nil {
			prometheus1Raw = prometheus1RawObj.(map[string]interface{})
		}
		if len(prometheus1Raw) > 0 {
			prometheusMap["external_url"] = prometheus1Raw["ExternalUrl"]
			prometheusMap["use_external"] = prometheus1Raw["UseExternal"]

			prometheusMaps = append(prometheusMaps, prometheusMap)
		}
		meshConfigMap["prometheus"] = prometheusMaps
		proxyMaps := make([]map[string]interface{}, 0)
		proxyMap := make(map[string]interface{})
		proxy1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.Proxy", objectRaw)
		proxy1Raw := make(map[string]interface{})
		if proxy1RawObj != nil {
			proxy1Raw = proxy1RawObj.(map[string]interface{})
		}
		if len(proxy1Raw) > 0 {
			proxyMap["cluster_domain"] = proxy1Raw["ClusterDomain"]
			proxyMap["limit_cpu"] = proxy1Raw["LimitCPU"]
			proxyMap["limit_memory"] = proxy1Raw["LimitMemory"]
			proxyMap["request_cpu"] = proxy1Raw["RequestCPU"]
			proxyMap["request_memory"] = proxy1Raw["RequestMemory"]

			proxyMaps = append(proxyMaps, proxyMap)
		}
		meshConfigMap["proxy"] = proxyMaps
		sidecarInjectorMaps := make([]map[string]interface{}, 0)
		sidecarInjectorMap := make(map[string]interface{})
		sidecarInjector1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.SidecarInjector", objectRaw)
		sidecarInjector1Raw := make(map[string]interface{})
		if sidecarInjector1RawObj != nil {
			sidecarInjector1Raw = sidecarInjector1RawObj.(map[string]interface{})
		}
		if len(sidecarInjector1Raw) > 0 {
			sidecarInjectorMap["auto_injection_policy_enabled"] = sidecarInjector1Raw["AutoInjectionPolicyEnabled"]
			sidecarInjectorMap["enable_namespaces_by_default"] = sidecarInjector1Raw["EnableNamespacesByDefault"]
			sidecarInjectorMap["limit_cpu"] = sidecarInjector1Raw["LimitCPU"]
			sidecarInjectorMap["limit_memory"] = sidecarInjector1Raw["LimitMemory"]
			sidecarInjectorMap["request_cpu"] = sidecarInjector1Raw["RequestCPU"]
			sidecarInjectorMap["request_memory"] = sidecarInjector1Raw["RequestMemory"]
			sidecarInjectorMap["sidecar_injector_webhook_as_yaml"] = sidecarInjector1Raw["SidecarInjectorWebhookAsYaml"]

			initCNIConfigurationMaps := make([]map[string]interface{}, 0)
			initCNIConfigurationMap := make(map[string]interface{})
			initCNIConfiguration1RawObj, _ := jsonpath.Get("$.Spec.MeshConfig.SidecarInjector.InitCNIConfiguration", objectRaw)
			initCNIConfiguration1Raw := make(map[string]interface{})
			if initCNIConfiguration1RawObj != nil {
				initCNIConfiguration1Raw = initCNIConfiguration1RawObj.(map[string]interface{})
			}
			if len(initCNIConfiguration1Raw) > 0 {
				initCNIConfigurationMap["enabled"] = initCNIConfiguration1Raw["Enabled"]
				initCNIConfigurationMap["exclude_namespaces"] = initCNIConfiguration1Raw["ExcludeNamespaces"]

				initCNIConfigurationMaps = append(initCNIConfigurationMaps, initCNIConfigurationMap)
			}
			sidecarInjectorMap["init_cni_configuration"] = initCNIConfigurationMaps
			sidecarInjectorMaps = append(sidecarInjectorMaps, sidecarInjectorMap)
		}
		meshConfigMap["sidecar_injector"] = sidecarInjectorMaps
	}

	networkMaps := make([]map[string]interface{}, 0)
	networkMap := make(map[string]interface{})
	network1RawObj, _ := jsonpath.Get("$.Spec.Network", objectRaw)
	network1Raw := make(map[string]interface{})
	if network1RawObj != nil {
		network1Raw = network1RawObj.(map[string]interface{})
	}
	if len(network1Raw) > 0 {
		networkMap["security_group_id"] = network1Raw["SecurityGroupId"]
		networkMap["vpc_id"] = network1Raw["VpcId"]

		vSwitches1Raw, _ := jsonpath.Get("$.Spec.Network.VSwitches", objectRaw)
		networkMap["vswitche_list"] = vSwitches1Raw
		networkMaps = append(networkMaps, networkMap)
	}
	if network1RawObj != nil {
		if err := d.Set("network", networkMaps); err != nil {
			return err
		}
	}

	objectRaw, err = serviceMeshServiceV2.DescribeDescribeKialiConfiguration(d.Id())
	if err != nil {
		return WrapError(err)
	}

	kialiMaps := make([]map[string]interface{}, 0)
	kialiMap := make(map[string]interface{})

	kialiMap["aggregated_kiali_address"] = objectRaw["AggregatedKialiAddress"]
	kialiMap["auth_strategy"] = objectRaw["AuthStrategy"]
	kialiMap["custom_prometheus_url"] = objectRaw["CustomPrometheusUrl"]
	kialiMap["distributed_kiali_access_tokens"] = objectRaw["DistributedKialiAccessTokens"]
	kialiMap["distributed_kiali_addresses"] = objectRaw["DistributedKialiAddresses"]
	kialiMap["enabled"] = objectRaw["Enabled"]
	kialiMap["integrate_clb"] = objectRaw["IntegrateCLB"]
	kialiMap["kiali_arms_auth_tokens"] = objectRaw["KialiArmsAuthTokens"]
	kialiMap["kiali_service_annotations"] = objectRaw["KialiServiceAnnotations"]
	kialiMap["url"] = objectRaw["KialiIstioIngressAddress"]
	kialiMap["use_populated_arms_prometheus"] = objectRaw["UsePopulatedArmsPrometheus"]

	openIdConfigMaps := make([]map[string]interface{}, 0)
	openIdConfigMap := make(map[string]interface{})
	openIdConfig1RawObj, _ := jsonpath.Get("$.OpenIdConfig", objectRaw)
	openIdConfig1Raw := make(map[string]interface{})
	if openIdConfig1RawObj != nil {
		openIdConfig1Raw = openIdConfig1RawObj.(map[string]interface{})
	}
	if len(openIdConfig1Raw) > 0 {
		openIdConfigMap["client_id"] = openIdConfig1Raw["ClientId"]
		openIdConfigMap["client_secret"] = openIdConfig1Raw["ClientSecret"]
		openIdConfigMap["issuer_uri"] = openIdConfig1Raw["IssuerUri"]

		scopes1Raw, _ := jsonpath.Get("$.OpenIdConfig.Scopes", objectRaw)
		openIdConfigMap["scopes"] = scopes1Raw
		openIdConfigMaps = append(openIdConfigMaps, openIdConfigMap)
	}
	kialiMap["open_id_config"] = openIdConfigMaps
	ramOAuthConfigMaps := make([]map[string]interface{}, 0)
	ramOAuthConfigMap := make(map[string]interface{})
	rAMOAuthConfig1RawObj, _ := jsonpath.Get("$.RAMOAuthConfig", objectRaw)
	rAMOAuthConfig1Raw := make(map[string]interface{})
	if rAMOAuthConfig1RawObj != nil {
		rAMOAuthConfig1Raw = rAMOAuthConfig1RawObj.(map[string]interface{})
	}
	if len(rAMOAuthConfig1Raw) > 0 {
		ramOAuthConfigMap["redirect_uris"] = rAMOAuthConfig1Raw["RedirectUris"]

		ramOAuthConfigMaps = append(ramOAuthConfigMaps, ramOAuthConfigMap)
	}
	kialiMap["ram_oauth_config"] = ramOAuthConfigMaps
	serverConfigMaps := make([]map[string]interface{}, 0)
	serverConfigMap := make(map[string]interface{})
	serverConfig1RawObj, _ := jsonpath.Get("$.ServerConfig", objectRaw)
	serverConfig1Raw := make(map[string]interface{})
	if serverConfig1RawObj != nil {
		serverConfig1Raw = serverConfig1RawObj.(map[string]interface{})
	}
	if len(serverConfig1Raw) > 0 {
		serverConfigMap["web_fqdn"] = serverConfig1Raw["WebFQDN"]
		serverConfigMap["web_port"] = serverConfig1Raw["WebPort"]
		serverConfigMap["web_root"] = serverConfig1Raw["WebRoot"]
		serverConfigMap["web_schema"] = serverConfig1Raw["WebSchema"]

		serverConfigMaps = append(serverConfigMaps, serverConfigMap)
	}
	kialiMap["server_config"] = serverConfigMaps
	kialiMaps = append(kialiMaps, kialiMap)
	meshConfigMap["kiali"] = kialiMaps
	meshConfigMaps = append(meshConfigMaps, meshConfigMap)
	if err := d.Set("mesh_config", meshConfigMaps); err != nil {
		return err
	}

	objectRaw, err = serviceMeshServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = serviceMeshServiceV2.DescribeDescribeServiceMeshKubeconfig(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["Kubeconfig"] != nil {
		d.Set("kubeconfig", objectRaw["Kubeconfig"])
	}

	return nil
}

func resourceAliCloudServiceMeshServiceMeshUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateMeshFeature"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("mesh_config.0.tracing") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].tracing", d.Get("mesh_config"))
		if err == nil {
			request["Tracing"] = jsonPathResult
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.pilot.0.trace_sampling") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].pilot[0].trace_sampling", d.Get("mesh_config"))
		if err == nil {
			request["TraceSampling"] = jsonPathResult1
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.telemetry") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].telemetry", d.Get("mesh_config"))
		if err == nil {
			request["Telemetry"] = jsonPathResult2
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.customized_zipkin") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].customized_zipkin", d.Get("mesh_config"))
		if err == nil {
			request["CustomizedZipkin"] = jsonPathResult3
		}
	}

	if d.HasChange("mesh_config.0.outbound_traffic_policy") {
		update = true
		jsonPathResult4, err := jsonpath.Get("$[0].outbound_traffic_policy", d.Get("mesh_config"))
		if err == nil {
			request["OutboundTrafficPolicy"] = jsonPathResult4
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.include_ip_ranges") {
		update = true
		jsonPathResult5, err := jsonpath.Get("$[0].include_ip_ranges", d.Get("mesh_config"))
		if err == nil {
			request["IncludeIPRanges"] = jsonPathResult5
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.enable_namespaces_by_default") {
		update = true
		jsonPathResult6, err := jsonpath.Get("$[0].sidecar_injector[0].enable_namespaces_by_default", d.Get("mesh_config"))
		if err == nil {
			request["EnableNamespacesByDefault"] = jsonPathResult6
		}
	}

	if d.HasChange("mesh_config.0.pilot.0.http10_enabled") {
		update = true
		jsonPathResult7, err := jsonpath.Get("$[0].pilot[0].http10_enabled", d.Get("mesh_config"))
		if err == nil {
			request["Http10Enabled"] = jsonPathResult7
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.opa.0.log_level") {
		update = true
		jsonPathResult8, err := jsonpath.Get("$[0].opa[0].log_level", d.Get("mesh_config"))
		if err == nil {
			request["OPALogLevel"] = jsonPathResult8
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.opa.0.request_cpu") {
		update = true
		jsonPathResult9, err := jsonpath.Get("$[0].opa[0].request_cpu", d.Get("mesh_config"))
		if err == nil {
			request["OPARequestCPU"] = jsonPathResult9
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.opa.0.request_memory") {
		update = true
		jsonPathResult10, err := jsonpath.Get("$[0].opa[0].request_memory", d.Get("mesh_config"))
		if err == nil {
			request["OPARequestMemory"] = jsonPathResult10
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.opa.0.limit_cpu") {
		update = true
		jsonPathResult11, err := jsonpath.Get("$[0].opa[0].limit_cpu", d.Get("mesh_config"))
		if err == nil {
			request["OPALimitCPU"] = jsonPathResult11
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.opa.0.limit_memory") {
		update = true
		jsonPathResult12, err := jsonpath.Get("$[0].opa[0].limit_memory", d.Get("mesh_config"))
		if err == nil {
			request["OPALimitMemory"] = jsonPathResult12
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.proxy.0.limit_cpu") {
		update = true
		jsonPathResult13, err := jsonpath.Get("$[0].proxy[0].limit_cpu", d.Get("mesh_config"))
		if err == nil {
			request["ProxyLimitCPU"] = jsonPathResult13
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.proxy.0.request_cpu") {
		update = true
		jsonPathResult14, err := jsonpath.Get("$[0].proxy[0].request_cpu", d.Get("mesh_config"))
		if err == nil {
			request["ProxyRequestCPU"] = jsonPathResult14
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.proxy.0.limit_memory") {
		update = true
		jsonPathResult15, err := jsonpath.Get("$[0].proxy[0].limit_memory", d.Get("mesh_config"))
		if err == nil {
			request["ProxyLimitMemory"] = jsonPathResult15
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.access_log.0.enabled") {
		update = true
		jsonPathResult16, err := jsonpath.Get("$[0].access_log[0].enabled", d.Get("mesh_config"))
		if err == nil {
			request["AccessLogEnabled"] = jsonPathResult16
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.init_cni_configuration.0.exclude_namespaces") && (d.Get("cluster_spec") == "enterprise" || d.Get("cluster_spec") == "ultimate") {
		update = true
		jsonPathResult18, err := jsonpath.Get("$[0].sidecar_injector[0].init_cni_configuration[0].exclude_namespaces", d.Get("mesh_config"))
		if err == nil {
			request["CniExcludeNamespaces"] = jsonPathResult18
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.init_cni_configuration.0.enabled") && (d.Get("cluster_spec") == "enterprise" || d.Get("cluster_spec") == "ultimate") {
		update = true
	}
	jsonPathResult19, err := jsonpath.Get("$[0].sidecar_injector[0].init_cni_configuration[0].enabled", d.Get("mesh_config"))
	if err == nil {
		request["CniEnabled"] = jsonPathResult19
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.proxy.0.request_memory") {
		update = true
		jsonPathResult19, err := jsonpath.Get("$[0].proxy[0].request_memory", d.Get("mesh_config"))
		if err == nil {
			request["ProxyRequestMemory"] = jsonPathResult19
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.request_memory") {
		update = true
		jsonPathResult20, err := jsonpath.Get("$[0].sidecar_injector[0].request_memory", d.Get("mesh_config"))
		if err == nil {
			request["SidecarInjectorRequestMemory"] = jsonPathResult20
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.limit_memory") {
		update = true
		jsonPathResult21, err := jsonpath.Get("$[0].sidecar_injector[0].limit_memory", d.Get("mesh_config"))
		if err == nil {
			request["SidecarInjectorLimitMemory"] = jsonPathResult21
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.limit_cpu") {
		update = true
		jsonPathResult22, err := jsonpath.Get("$[0].sidecar_injector[0].limit_cpu", d.Get("mesh_config"))
		if err == nil {
			request["SidecarInjectorLimitCPU"] = jsonPathResult22
		}
	}

	if d.HasChange("mesh_config.0.sidecar_injector.0.request_cpu") {
		update = true
		jsonPathResult23, err := jsonpath.Get("$[0].sidecar_injector[0].request_cpu", d.Get("mesh_config"))
		if err == nil {
			request["SidecarInjectorRequestCPU"] = jsonPathResult23
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.opa.0.enabled") {
		update = true
	}
	jsonPathResult25, err := jsonpath.Get("$[0].opa[0].enabled", d.Get("mesh_config"))
	if err == nil {
		request["OpaEnabled"] = jsonPathResult25
	}

	if v, ok := d.GetOkExists("customized_prometheus"); ok {
		request["CustomizedPrometheus"] = v
	}
	if v, ok := d.GetOk("prometheus_url"); ok {
		request["PrometheusUrl"] = v
	}
	if d.HasChange("mesh_config.0.sidecar_injector.0.auto_injection_policy_enabled") {
		update = true
	}
	jsonPathResult28, err := jsonpath.Get("$[0].sidecar_injector[0].auto_injection_policy_enabled", d.Get("mesh_config"))
	if err == nil {
		request["AutoInjectionPolicyEnabled"] = jsonPathResult28
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.access_log.0.project") {
		update = true
		jsonPathResult28, err := jsonpath.Get("$[0].access_log[0].project", d.Get("mesh_config"))
		if err == nil {
			request["AccessLogProject"] = jsonPathResult28
		}
	}

	if d.HasChange("mesh_config.0.access_log.0.gateway_enabled") {
		update = true
		jsonPathResult29, err := jsonpath.Get("$[0].access_log[0].gateway_enabled", d.Get("mesh_config"))
		if err == nil {
			request["AccessLogGatewayEnabled"] = jsonPathResult29
		}
	}

	if d.HasChange("mesh_config.0.access_log.0.sidecar_enabled") {
		update = true
		jsonPathResult30, err := jsonpath.Get("$[0].access_log[0].sidecar_enabled", d.Get("mesh_config"))
		if err == nil {
			request["AccessLogSidecarEnabled"] = jsonPathResult30
		}
	}

	if d.HasChange("mesh_config.0.access_log.0.gateway_lifecycle") {
		update = true
		jsonPathResult31, err := jsonpath.Get("$[0].access_log[0].gateway_lifecycle", d.Get("mesh_config"))
		if err == nil && jsonPathResult31.(int) > 0 {
			request["AccessLogGatewayLifecycle"] = jsonPathResult31
		}
	}

	if d.HasChange("mesh_config.0.access_log.0.sidecar_lifecycle") {
		update = true
		jsonPathResult32, err := jsonpath.Get("$[0].access_log[0].sidecar_lifecycle", d.Get("mesh_config"))
		if err == nil && jsonPathResult32.(int) > 0 {
			request["AccessLogSidecarLifecycle"] = jsonPathResult32
		}
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.audit.0.project") {
		update = true
	}
	jsonPathResult34, err := jsonpath.Get("$[0].audit[0].project", d.Get("mesh_config"))
	if err == nil {
		request["AuditProject"] = jsonPathResult34
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.audit.0.enabled") {
		update = true
	}
	jsonPathResult35, err := jsonpath.Get("$[0].audit[0].enabled", d.Get("mesh_config"))
	if err == nil {
		request["EnableAudit"] = jsonPathResult35
	}

	if !d.IsNewResource() && d.HasChange("cluster_spec") {
		update = true
		request["ClusterSpec"] = d.Get("cluster_spec")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		serviceMeshServiceV2 := ServiceMeshServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateMeshCRAggregation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if d.HasChange("extra_configuration.0.cr_aggregation_enabled") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].cr_aggregation_enabled", d.Get("extra_configuration"))
		if err == nil {
			request["Enabled"] = jsonPathResult
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		serviceMeshServiceV2 := ServiceMeshServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpgradeMeshEditionPartially"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("version") {
		update = true
		request["ExpectedVersion"] = d.Get("version")
	}

	request["ASMGatewayContinue"] = "false"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		serviceMeshServiceV2 := ServiceMeshServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateControlPlaneLogConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("mesh_config.0.control_plane_log.0.enabled") {
		update = true
	}
	jsonPathResult, err := jsonpath.Get("$[0].control_plane_log[0].enabled", d.Get("mesh_config"))
	if err == nil {
		request["Enabled"] = jsonPathResult
	}

	if !d.IsNewResource() && d.HasChange("mesh_config.0.control_plane_log.0.project") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].control_plane_log[0].project", d.Get("mesh_config"))
		if err == nil {
			request["Project"] = jsonPathResult1
		}
	}

	if d.HasChange("mesh_config.0.control_plane_log.0.log_ttl_in_day") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].control_plane_log[0].log_ttl_in_day", d.Get("mesh_config"))
		if err == nil && jsonPathResult2.(int) > 0 {
			request["LogTTLInDay"] = jsonPathResult2
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		serviceMeshServiceV2 := ServiceMeshServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyServiceMeshName"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("service_mesh_name") {
		update = true
	}
	request["Name"] = d.Get("service_mesh_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "UpdateKialiConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("mesh_config.0.kiali.0.enabled") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].kiali[0].enabled", d.Get("mesh_config"))
		if err == nil {
			request["Enabled"] = jsonPathResult
		}
	}

	if d.HasChange("mesh_config.0.kiali.0.custom_prometheus_url") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].kiali[0].custom_prometheus_url", d.Get("mesh_config"))
		if err == nil {
			request["CustomPrometheusUrl"] = jsonPathResult1
		}
	}

	if d.HasChange("mesh_config.0.kiali.0.integrate_clb") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].kiali[0].integrate_clb", d.Get("mesh_config"))
		if err == nil {
			request["IntegrateCLB"] = jsonPathResult2
		}
	}

	if d.HasChange("mesh_config.0.kiali.0.kiali_service_annotations") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].kiali[0].kiali_service_annotations", d.Get("mesh_config"))
		if err == nil {
			request["KialiServiceAnnotations"] = jsonPathResult3
		}
	}

	if d.HasChange("mesh_config.0.kiali.0.kiali_arms_auth_tokens") {
		update = true
		jsonPathResult4, err := jsonpath.Get("$[0].kiali[0].kiali_arms_auth_tokens", d.Get("mesh_config"))
		if err == nil {
			request["KialiArmsAuthTokens"] = jsonPathResult4
		}
	}

	if d.HasChange("mesh_config.0.kiali.0.auth_strategy") {
		update = true
		jsonPathResult5, err := jsonpath.Get("$[0].kiali[0].auth_strategy", d.Get("mesh_config"))
		if err == nil {
			request["AuthStrategy"] = jsonPathResult5
		}
	}

	if d.HasChange("mesh_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("mesh_config"); !IsNil(v) {
			redirectUris1, _ := jsonpath.Get("$[0].kiali[0].ram_oauth_config[0].redirect_uris", v)
			if redirectUris1 != nil && (d.HasChange("mesh_config.0.kiali.0.ram_oauth_config.0.redirect_uris") || redirectUris1 != "") {
				objectDataLocalMap["RedirectUris"] = redirectUris1
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			request["RAMOAuthConfig"] = string(objectDataLocalMapJson)
		}
	}

	if d.HasChange("mesh_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("mesh_config"); !IsNil(v) {
			clientId1, _ := jsonpath.Get("$[0].kiali[0].open_id_config[0].client_id", v)
			if clientId1 != nil && (d.HasChange("mesh_config.0.kiali.0.open_id_config.0.client_id") || clientId1 != "") {
				objectDataLocalMap1["ClientId"] = clientId1
			}
			clientSecret1, _ := jsonpath.Get("$[0].kiali[0].open_id_config[0].client_secret", v)
			if clientSecret1 != nil && (d.HasChange("mesh_config.0.kiali.0.open_id_config.0.client_secret") || clientSecret1 != "") {
				objectDataLocalMap1["ClientSecret"] = clientSecret1
			}
			issuerUri1, _ := jsonpath.Get("$[0].kiali[0].open_id_config[0].issuer_uri", v)
			if issuerUri1 != nil && (d.HasChange("mesh_config.0.kiali.0.open_id_config.0.issuer_uri") || issuerUri1 != "") {
				objectDataLocalMap1["IssuerUri"] = issuerUri1
			}
			scopes1, _ := jsonpath.Get("$[0].kiali[0].open_id_config[0].scopes", d.Get("mesh_config"))
			if scopes1 != nil && (d.HasChange("mesh_config.0.kiali.0.open_id_config.0.scopes") || scopes1 != "") {
				objectDataLocalMap1["Scopes"] = scopes1
			}

			objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
			if err != nil {
				return WrapError(err)
			}
			request["OpenIdConfig"] = string(objectDataLocalMap1Json)
		}
	}

	if d.HasChange("mesh_config") {
		update = true
		objectDataLocalMap2 := make(map[string]interface{})

		if v := d.Get("mesh_config"); !IsNil(v) {
			webFqdn, _ := jsonpath.Get("$[0].kiali[0].server_config[0].web_fqdn", v)
			if webFqdn != nil && (d.HasChange("mesh_config.0.kiali.0.server_config.0.web_fqdn") || webFqdn != "") {
				objectDataLocalMap2["WebFQDN"] = webFqdn
			}
			webPort1, _ := jsonpath.Get("$[0].kiali[0].server_config[0].web_port", v)
			if webPort1 != nil && (d.HasChange("mesh_config.0.kiali.0.server_config.0.web_port") || webPort1 != "") {
				objectDataLocalMap2["WebPort"] = webPort1
			}
			webRoot1, _ := jsonpath.Get("$[0].kiali[0].server_config[0].web_root", v)
			if webRoot1 != nil && (d.HasChange("mesh_config.0.kiali.0.server_config.0.web_root") || webRoot1 != "") {
				objectDataLocalMap2["WebRoot"] = webRoot1
			}
			webSchema1, _ := jsonpath.Get("$[0].kiali[0].server_config[0].web_schema", v)
			if webSchema1 != nil && (d.HasChange("mesh_config.0.kiali.0.server_config.0.web_schema") || webSchema1 != "") {
				objectDataLocalMap2["WebSchema"] = webSchema1
			}

			objectDataLocalMap2Json, err := json.Marshal(objectDataLocalMap2)
			if err != nil {
				return WrapError(err)
			}
			request["ServerConfig"] = string(objectDataLocalMap2Json)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		serviceMeshServiceV2 := ServiceMeshServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("cluster_ids") {
		oldEntry, newEntry := d.GetChange("cluster_ids")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			clusterIds := removed.([]interface{})

			for _, item := range clusterIds {
				action := "RemoveClusterFromServiceMesh"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ServiceMeshId"] = d.Id()

				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["ClusterId"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				serviceMeshServiceV2 := ServiceMeshServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}

		if len(added.([]interface{})) > 0 {
			clusterIds := added.([]interface{})

			for _, item := range clusterIds {
				action := "AddClusterIntoServiceMesh"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ServiceMeshId"] = d.Id()

				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["ClusterId"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				serviceMeshServiceV2 := ServiceMeshServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}
	if d.HasChange("tags") {
		serviceMeshServiceV2 := ServiceMeshServiceV2{client}
		if err := serviceMeshServiceV2.SetResourceTags(d, "servicemesh"); err != nil {
			return WrapError(err)
		}
	}
	if !d.IsNewResource() && d.HasChange("load_balancer.0.pilot_public_eip_id") {
		oldEntry, newEntry := d.GetChange("load_balancer.0.pilot_public_eip_id")
		removed, removedOk := oldEntry.(string)
		added, addedOk := newEntry.(string)

		log.Printf("[INFO] Resource load_balancer.0.pilot_public_eip_id  removed %s removedOk %v added %s addedOk %v", removed, removedOk, added, addedOk)
		if removedOk && len(removed) > 0 && added != removed {
			action := "ModifyPilotEipResource"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServiceMeshId"] = d.Id()

			request["Operation"] = "UnBindEip"
			if v, ok := d.GetOk("load_balancer"); ok {
				jsonPathResult, err := jsonpath.Get("$[0].pilot_public_eip_id", v)
				if err == nil && jsonPathResult != "" {
					request["EipId"] = jsonPathResult
				}
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			serviceMeshServiceV2 := ServiceMeshServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if addedOk && len(added) > 0 && added != removed {
			action := "ModifyPilotEipResource"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServiceMeshId"] = d.Id()

			if v, ok := d.GetOk("load_balancer"); ok {
				jsonPathResult, err := jsonpath.Get("$[0].pilot_public_eip_id", v)
				if err == nil && jsonPathResult != "" {
					request["EipId"] = jsonPathResult
				}
			}
			request["Operation"] = "BindEip"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			serviceMeshServiceV2 := ServiceMeshServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	d.Partial(false)
	return resourceAliCloudServiceMeshServiceMeshRead(d, meta)
}

func resourceAliCloudServiceMeshServiceMeshDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceMesh"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ServiceMeshId"] = d.Id()

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"RelatedResourceReused", "StillInitializing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"StatusForbidden", "403"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	serviceMeshServiceV2 := ServiceMeshServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 60*time.Second, serviceMeshServiceV2.ServiceMeshServiceMeshStateRefreshFunc(d.Id(), "$.ServiceMeshInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

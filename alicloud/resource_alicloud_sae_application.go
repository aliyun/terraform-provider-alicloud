package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSaeApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSaeApplicationCreate,
		Read:   resourceAliCloudSaeApplicationRead,
		Update: resourceAliCloudSaeApplicationUpdate,
		Delete: resourceAliCloudSaeApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"FatJar", "War", "Image", "PhpZip", "IMAGE_PHP_5_4", "IMAGE_PHP_5_4_ALPINE", "IMAGE_PHP_5_5", "IMAGE_PHP_5_5_ALPINE", "IMAGE_PHP_5_6", "IMAGE_PHP_5_6_ALPINE", "IMAGE_PHP_7_0", "IMAGE_PHP_7_0_ALPINE", "IMAGE_PHP_7_1", "IMAGE_PHP_7_1_ALPINE", "IMAGE_PHP_7_2", "IMAGE_PHP_7_2_ALPINE", "IMAGE_PHP_7_3", "IMAGE_PHP_7_3_ALPINE", "PythonZip"}, false),
			},
			"replicas": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"package_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{500, 1000, 2000, 4000, 8000, 16000, 32000}),
			},
			"memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1024, 2048, 4096, 8192, 12288, 16384, 24576, 32768, 65536, 131072}),
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdk": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jar_start_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jar_start_args": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_config": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_enable_application_scaling_rule": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"batch_wait_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"change_order_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"edas_container_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ahas": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_grey_tag_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("enable_grey_tag_route").(bool)
				},
			},
			"min_ready_instances": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_ready_instance_ratio": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"oss_ak_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"oss_ak_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"php_arms_config_location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"php_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"php_config_location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 60),
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"war_start_options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acr_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acr_assume_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"micro_registration": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1", "2"}, false),
			},
			"envs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sls_configs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"php": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_pull_secrets": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"programming_language": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"command_args_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"command_args"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"custom_host_alias_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"custom_host_alias"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"oss_mount_descs_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"oss_mount_descs"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bucket_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"config_map_mount_desc_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"config_map_mount_desc"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_map_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"liveness_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"liveness"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"initial_delay_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"exec": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"tcp_socket": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"http_get": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"scheme": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key_word": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"is_contain_key_word": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"readiness_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"readiness"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"initial_delay_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"timeout_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"exec": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"tcp_socket": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"http_get": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"scheme": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key_word": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"is_contain_key_word": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"post_start_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"post_start"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exec": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"pre_stop_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pre_stop"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exec": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"tomcat_config_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"tomcat_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_threads": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"context_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uri_encoding": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"use_body_encoding_for_uri": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"update_strategy_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"update_strategy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"batch_update": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"batch": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"batch_wait_time": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"nas_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nas_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nas_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mount_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"kafka_configs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kafka_instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kafka_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kafka_configs": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"log_dir": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"kafka_topic": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"pvtz_discovery_svc": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"namespace_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"port_protocols": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": tagsSchema(),
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"RUNNING", "STOPPED", "UNKNOWN"}, false),
			},
			"command_args": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"command_args_v2"},
				Deprecated:    "Field `command_args` has been deprecated from provider version 1.211.0. New field `command_args_v2` instead.",
			},
			"custom_host_alias": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"custom_host_alias_v2"},
				Deprecated:    "Field `custom_host_alias` has been deprecated from provider version 1.211.0. New field `custom_host_alias_v2` instead.",
			},
			"oss_mount_descs": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"oss_mount_descs_v2"},
				Deprecated:    "Field `oss_mount_descs` has been deprecated from provider version 1.211.0. New field `oss_mount_descs_v2` instead.",
			},
			"config_map_mount_desc": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"config_map_mount_desc_v2"},
				Deprecated:    "Field `config_map_mount_desc` has been deprecated from provider version 1.211.0. New field `config_map_mount_desc_v2` instead.",
			},
			"liveness": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"liveness_v2"},
				Deprecated:    "Field `liveness` has been deprecated from provider version 1.211.0. New field `liveness_v2` instead.",
			},
			"readiness": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"readiness_v2"},
				Deprecated:    "Field `readiness` has been deprecated from provider version 1.211.0. New field `readiness_v2` instead.",
			},
			"post_start": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"post_start_v2"},
				Deprecated:    "Field `post_start` has been deprecated from provider version 1.211.0. New field `post_start_v2` instead.",
			},
			"pre_stop": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"pre_stop_v2"},
				Deprecated:    "Field `pre_stop` has been deprecated from provider version 1.211.0. New field `pre_stop_v2` instead.",
			},
			"tomcat_config": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"tomcat_config_v2"},
				Deprecated:    "Field `tomcat_config` has been deprecated from provider version 1.211.0. New field `tomcat_config_v2` instead.",
			},
			"update_strategy": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"update_strategy_v2"},
				Deprecated:    "Field `update_strategy` has been deprecated from provider version 1.211.0. New field `update_strategy_v2` instead.",
			},
			"nas_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `nas_id` has been removed from provider version 1.211.0.",
			},
			"mount_host": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `mount_host` has been removed from provider version 1.211.0.",
			},
			"mount_desc": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `mount_desc` has been removed from provider version 1.211.0.",
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Field `version_id` has been removed from provider version 1.211.0.",
			},
		},
	}
}

func resourceAliCloudSaeApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	var response map[string]interface{}
	action := "/pop/v1/sam/app/createApplication"
	request := make(map[string]*string)
	var err error

	request["AppName"] = StringPointer(d.Get("app_name").(string))
	request["PackageType"] = StringPointer(d.Get("package_type").(string))
	request["Replicas"] = StringPointer(strconv.Itoa(d.Get("replicas").(int)))
	request["PackageVersion"] = StringPointer(strconv.FormatInt(time.Now().Unix(), 10))

	if v, ok := d.GetOk("namespace_id"); ok {
		request["NamespaceId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("package_url"); ok {
		request["PackageUrl"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("image_url"); ok {
		request["ImageUrl"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request["Cpu"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("memory"); ok {
		request["Memory"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("command"); ok {
		request["Command"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("web_container"); ok {
		request["WebContainer"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("jdk"); ok {
		request["Jdk"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("jar_start_options"); ok {
		request["JarStartOptions"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("jar_start_args"); ok {
		request["JarStartArgs"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("app_description"); ok {
		request["AppDescription"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOkExists("auto_config"); ok {
		request["AutoConfig"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOkExists("deploy"); ok {
		request["Deploy"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("edas_container_version"); ok {
		request["EdasContainerVersion"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("oss_ak_id"); ok {
		request["OssAkId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("oss_ak_secret"); ok {
		request["OssAkSecret"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("php_arms_config_location"); ok {
		request["PhpArmsConfigLocation"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("php_config"); ok {
		request["PhpConfig"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("php_config_location"); ok {
		request["PhpConfigLocation"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("termination_grace_period_seconds"); ok {
		request["TerminationGracePeriodSeconds"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("timezone"); ok {
		request["Timezone"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("war_start_options"); ok {
		request["WarStartOptions"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("acr_instance_id"); ok {
		request["AcrInstanceId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("acr_assume_role_arn"); ok {
		request["AcrAssumeRoleArn"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("micro_registration"); ok {
		request["MicroRegistration"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("envs"); ok {
		request["Envs"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("sls_configs"); ok {
		request["SlsConfigs"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("php"); ok {
		request["Php"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("image_pull_secrets"); ok {
		request["ImagePullSecrets"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("programming_language"); ok {
		request["ProgrammingLanguage"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("command_args_v2"); ok {
		request["CommandArgs"] = StringPointer(convertListToJsonString(v.([]interface{})))
	} else if v, ok := d.GetOk("command_args"); ok {
		request["CommandArgs"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("custom_host_alias_v2"); ok {
		customHostAliasMaps := make([]map[string]interface{}, 0)
		for _, customHostAlias := range v.([]interface{}) {
			customHostAliasMap := map[string]interface{}{}
			customHostAliasArg := customHostAlias.(map[string]interface{})

			if hostName, ok := customHostAliasArg["host_name"]; ok && hostName.(string) != "" {
				customHostAliasMap["hostName"] = hostName
			}

			if ip, ok := customHostAliasArg["ip"]; ok && ip.(string) != "" {
				customHostAliasMap["ip"] = ip
			}

			customHostAliasMaps = append(customHostAliasMaps, customHostAliasMap)
		}

		customHostAliasJson, err := convertListMapToJsonString(customHostAliasMaps)
		if err != nil {
			return WrapError(err)
		}

		request["CustomHostAlias"] = StringPointer(customHostAliasJson)
	} else if v, ok := d.GetOk("custom_host_alias"); ok {
		request["CustomHostAlias"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("oss_mount_descs_v2"); ok {
		ossMountDescsMaps := make([]map[string]interface{}, 0)
		for _, ossMountDescs := range v.([]interface{}) {
			ossMountDescsMap := map[string]interface{}{}
			ossMountDescsArg := ossMountDescs.(map[string]interface{})

			if bucketName, ok := ossMountDescsArg["bucket_name"]; ok && bucketName.(string) != "" {
				ossMountDescsMap["bucketName"] = bucketName
			}

			if bucketPath, ok := ossMountDescsArg["bucket_path"]; ok && bucketPath.(string) != "" {
				ossMountDescsMap["bucketPath"] = bucketPath
			}

			if mountPath, ok := ossMountDescsArg["mount_path"]; ok && mountPath.(string) != "" {
				ossMountDescsMap["mountPath"] = mountPath
			}

			if readOnly, ok := ossMountDescsArg["read_only"]; ok {
				ossMountDescsMap["readOnly"] = readOnly
			}

			ossMountDescsMaps = append(ossMountDescsMaps, ossMountDescsMap)
		}

		ossMountDescsJson, err := convertListMapToJsonString(ossMountDescsMaps)
		if err != nil {
			return WrapError(err)
		}

		request["OssMountDescs"] = StringPointer(ossMountDescsJson)
	} else if v, ok := d.GetOk("oss_mount_descs"); ok {
		request["OssMountDescs"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("config_map_mount_desc_v2"); ok {
		configMapMountDescMaps := make([]map[string]interface{}, 0)
		for _, configMapMountDesc := range v.([]interface{}) {
			configMapMountDescMap := map[string]interface{}{}
			configMapMountDescArg := configMapMountDesc.(map[string]interface{})

			if configMapId, ok := configMapMountDescArg["config_map_id"]; ok && configMapId.(string) != "" {
				configMapMountDescMap["configMapId"] = configMapId
			}

			if mountPath, ok := configMapMountDescArg["mount_path"]; ok && mountPath.(string) != "" {
				configMapMountDescMap["mountPath"] = mountPath
			}

			if key, ok := configMapMountDescArg["key"]; ok && key.(string) != "" {
				configMapMountDescMap["key"] = key
			}

			configMapMountDescMaps = append(configMapMountDescMaps, configMapMountDescMap)
		}

		configMapMountDescJson, err := convertListMapToJsonString(configMapMountDescMaps)
		if err != nil {
			return WrapError(err)
		}

		request["ConfigMapMountDesc"] = StringPointer(configMapMountDescJson)
	} else if v, ok := d.GetOk("config_map_mount_desc"); ok {
		request["ConfigMapMountDesc"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("liveness_v2"); ok {
		livenessMap := map[string]interface{}{}
		for _, livenessList := range v.([]interface{}) {
			livenessArg := livenessList.(map[string]interface{})

			if initialDelaySeconds, ok := livenessArg["initial_delay_seconds"]; ok {
				livenessMap["initialDelaySeconds"] = initialDelaySeconds
			}

			if periodSeconds, ok := livenessArg["period_seconds"]; ok {
				livenessMap["periodSeconds"] = periodSeconds
			}

			if timeoutSeconds, ok := livenessArg["timeout_seconds"]; ok {
				livenessMap["timeoutSeconds"] = timeoutSeconds
			}

			if exec, ok := livenessArg["exec"]; ok && len(exec.([]interface{})) > 0 {
				execMap := map[string]interface{}{}
				for _, execList := range exec.([]interface{}) {
					execArg := execList.(map[string]interface{})

					if command, ok := execArg["command"]; ok {
						execMap["command"] = command
					}
				}

				livenessMap["exec"] = execMap
			}

			if tcpSocket, ok := livenessArg["tcp_socket"]; ok && len(tcpSocket.([]interface{})) > 0 {
				tcpSocketMap := map[string]interface{}{}
				for _, tcpSocketList := range tcpSocket.([]interface{}) {
					tcpSocketArg := tcpSocketList.(map[string]interface{})

					if port, ok := tcpSocketArg["port"]; ok {
						tcpSocketMap["port"] = port
					}
				}

				livenessMap["tcpSocket"] = tcpSocketMap
			}

			if httpGet, ok := livenessArg["http_get"]; ok && len(httpGet.([]interface{})) > 0 {
				httpGetMap := map[string]interface{}{}
				for _, httpGetList := range httpGet.([]interface{}) {
					httpGetArg := httpGetList.(map[string]interface{})

					if path, ok := httpGetArg["path"]; ok && path.(string) != "" {
						httpGetMap["path"] = path
					}

					if port, ok := httpGetArg["port"]; ok {
						httpGetMap["port"] = port
					}

					if scheme, ok := httpGetArg["scheme"]; ok && scheme.(string) != "" {
						httpGetMap["scheme"] = scheme
					}

					if keyWord, ok := httpGetArg["key_word"]; ok && keyWord.(string) != "" {
						httpGetMap["keyWord"] = keyWord

						if isContainKeyWord, ok := d.GetOkExists("liveness_v2.0.http_get.0.is_contain_key_word"); ok {
							httpGetMap["isContainKeyWord"] = isContainKeyWord
						}
					}
				}

				livenessMap["httpGet"] = httpGetMap
			}
		}

		livenessJson, err := convertMaptoJsonString(livenessMap)
		if err != nil {
			return WrapError(err)
		}

		request["Liveness"] = StringPointer(livenessJson)
	} else if v, ok := d.GetOk("liveness"); ok {
		request["Liveness"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("readiness_v2"); ok {
		readinessMap := map[string]interface{}{}
		for _, readinessList := range v.([]interface{}) {
			readinessArg := readinessList.(map[string]interface{})

			if initialDelaySeconds, ok := readinessArg["initial_delay_seconds"]; ok {
				readinessMap["initialDelaySeconds"] = initialDelaySeconds
			}

			if periodSeconds, ok := readinessArg["period_seconds"]; ok {
				readinessMap["periodSeconds"] = periodSeconds
			}

			if timeoutSeconds, ok := readinessArg["timeout_seconds"]; ok {
				readinessMap["timeoutSeconds"] = timeoutSeconds
			}

			if exec, ok := readinessArg["exec"]; ok && len(exec.([]interface{})) > 0 {
				execMap := map[string]interface{}{}
				for _, execList := range exec.([]interface{}) {
					execArg := execList.(map[string]interface{})

					if command, ok := execArg["command"]; ok {
						execMap["command"] = command
					}
				}

				readinessMap["exec"] = execMap
			}

			if tcpSocket, ok := readinessArg["tcp_socket"]; ok && len(tcpSocket.([]interface{})) > 0 {
				tcpSocketMap := map[string]interface{}{}
				for _, tcpSocketList := range tcpSocket.([]interface{}) {
					tcpSocketArg := tcpSocketList.(map[string]interface{})

					if port, ok := tcpSocketArg["port"]; ok {
						tcpSocketMap["port"] = port
					}
				}

				readinessMap["tcpSocket"] = tcpSocketMap
			}

			if httpGet, ok := readinessArg["http_get"]; ok && len(httpGet.([]interface{})) > 0 {
				httpGetMap := map[string]interface{}{}
				for _, httpGetList := range httpGet.([]interface{}) {
					httpGetArg := httpGetList.(map[string]interface{})

					if path, ok := httpGetArg["path"]; ok && path.(string) != "" {
						httpGetMap["path"] = path
					}

					if port, ok := httpGetArg["port"]; ok {
						httpGetMap["port"] = port
					}

					if scheme, ok := httpGetArg["scheme"]; ok && scheme.(string) != "" {
						httpGetMap["scheme"] = scheme
					}

					if keyWord, ok := httpGetArg["key_word"]; ok && keyWord.(string) != "" {
						httpGetMap["keyWord"] = keyWord

						if isContainKeyWord, ok := d.GetOkExists("readiness_v2.0.http_get.0.is_contain_key_word"); ok {
							httpGetMap["isContainKeyWord"] = isContainKeyWord
						}
					}
				}

				readinessMap["httpGet"] = httpGetMap
			}
		}

		readinessJson, err := convertMaptoJsonString(readinessMap)
		if err != nil {
			return WrapError(err)
		}

		request["Readiness"] = StringPointer(readinessJson)

	} else if v, ok := d.GetOk("readiness"); ok {
		request["Readiness"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("post_start_v2"); ok {
		postStartMap := map[string]interface{}{}
		for _, postStartList := range v.([]interface{}) {
			postStartArg := postStartList.(map[string]interface{})

			if exec, ok := postStartArg["exec"]; ok && len(exec.([]interface{})) > 0 {
				execMap := map[string]interface{}{}
				for _, execList := range exec.([]interface{}) {
					execArg := execList.(map[string]interface{})

					if command, ok := execArg["command"]; ok {
						execMap["command"] = command
					}
				}

				postStartMap["exec"] = execMap
			}
		}

		postStartJson, err := convertMaptoJsonString(postStartMap)
		if err != nil {
			return WrapError(err)
		}

		request["PostStart"] = StringPointer(postStartJson)
	} else if v, ok := d.GetOk("post_start"); ok {
		request["PostStart"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("pre_stop_v2"); ok {
		preStopMap := map[string]interface{}{}
		for _, preStopList := range v.([]interface{}) {
			preStopArg := preStopList.(map[string]interface{})

			if exec, ok := preStopArg["exec"]; ok && len(exec.([]interface{})) > 0 {
				execMap := map[string]interface{}{}
				for _, execList := range exec.([]interface{}) {
					execArg := execList.(map[string]interface{})

					if command, ok := execArg["command"]; ok {
						execMap["command"] = command
					}
				}

				preStopMap["exec"] = execMap
			}
		}

		preStopJson, err := convertMaptoJsonString(preStopMap)
		if err != nil {
			return WrapError(err)
		}

		request["PreStop"] = StringPointer(preStopJson)
	} else if v, ok := d.GetOk("pre_stop"); ok {
		request["PreStop"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("tomcat_config_v2"); ok {
		tomcatConfigMap := map[string]interface{}{}
		for _, tomcatConfigList := range v.([]interface{}) {
			tomcatConfigArg := tomcatConfigList.(map[string]interface{})

			if port, ok := tomcatConfigArg["port"]; ok {
				tomcatConfigMap["port"] = port
			}

			if maxThreads, ok := tomcatConfigArg["max_threads"]; ok {
				tomcatConfigMap["maxThreads"] = maxThreads
			}

			if contextPath, ok := tomcatConfigArg["context_path"]; ok && contextPath.(string) != "" {
				tomcatConfigMap["contextPath"] = contextPath
			}

			if uriEncoding, ok := tomcatConfigArg["uri_encoding"]; ok && uriEncoding.(string) != "" {
				tomcatConfigMap["uriEncoding"] = uriEncoding
			}

			if useBodyEncodingForUri, ok := tomcatConfigArg["use_body_encoding_for_uri"]; ok && useBodyEncodingForUri.(string) != "" {
				tomcatConfigMap["useBodyEncodingForUri"] = useBodyEncodingForUri
			}
		}

		tomcatConfigJson, err := convertMaptoJsonString(tomcatConfigMap)
		if err != nil {
			return WrapError(err)
		}

		request["TomcatConfig"] = StringPointer(tomcatConfigJson)
	} else if v, ok := d.GetOk("tomcat_config"); ok {
		request["TomcatConfig"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("nas_configs"); ok {
		nasConfigsMaps := make([]map[string]interface{}, 0)
		for _, nasConfigs := range v.([]interface{}) {
			nasConfigsMap := map[string]interface{}{}
			nasConfigsArg := nasConfigs.(map[string]interface{})

			if nasId, ok := nasConfigsArg["nas_id"]; ok && nasId.(string) != "" {
				nasConfigsMap["nasId"] = nasId
			}

			if nasPath, ok := nasConfigsArg["nas_path"]; ok && nasPath.(string) != "" {
				nasConfigsMap["nasPath"] = nasPath
			}

			if mountPath, ok := nasConfigsArg["mount_path"]; ok && mountPath.(string) != "" {
				nasConfigsMap["mountPath"] = mountPath
			}

			if mountDomain, ok := nasConfigsArg["mount_domain"]; ok && mountDomain.(string) != "" {
				nasConfigsMap["mountDomain"] = mountDomain
			}

			if readOnly, ok := nasConfigsArg["read_only"]; ok {
				nasConfigsMap["readOnly"] = readOnly
			}

			nasConfigsMaps = append(nasConfigsMaps, nasConfigsMap)
		}

		nasConfigsJson, err := convertListMapToJsonString(nasConfigsMaps)
		if err != nil {
			return WrapError(err)
		}

		request["NasConfigs"] = StringPointer(nasConfigsJson)
	}

	if v, ok := d.GetOk("kafka_configs"); ok {
		kafkaConfigsMap := map[string]interface{}{}
		for _, kafkaConfigsList := range v.([]interface{}) {
			kafkaConfigsArg := kafkaConfigsList.(map[string]interface{})

			if kafkaInstanceId, ok := kafkaConfigsArg["kafka_instance_id"]; ok && kafkaInstanceId.(string) != "" {
				kafkaConfigsMap["kafkaInstanceId"] = kafkaInstanceId
			}

			if kafkaEndpoint, ok := kafkaConfigsArg["kafka_endpoint"]; ok && kafkaEndpoint.(string) != "" {
				kafkaConfigsMap["kafkaEndpoint"] = kafkaEndpoint
			}

			if kafkaConfigsValue, ok := kafkaConfigsArg["kafka_configs"]; ok && len(kafkaConfigsValue.([]interface{})) > 0 {
				kafkaConfigsValueMaps := make([]map[string]interface{}, 0)
				for _, kafkaConfigsValueList := range kafkaConfigsValue.([]interface{}) {
					kafkaConfigsValueMap := map[string]interface{}{}
					kafkaConfigsValueArg := kafkaConfigsValueList.(map[string]interface{})

					if logType, ok := kafkaConfigsValueArg["log_type"]; ok && logType.(string) != "" {
						kafkaConfigsValueMap["logType"] = logType
					}

					if logDir, ok := kafkaConfigsValueArg["log_dir"]; ok && logDir.(string) != "" {
						kafkaConfigsValueMap["logDir"] = logDir
					}

					if kafkaTopic, ok := kafkaConfigsValueArg["kafka_topic"]; ok && kafkaTopic.(string) != "" {
						kafkaConfigsValueMap["kafkaTopic"] = kafkaTopic
					}

					kafkaConfigsValueMaps = append(kafkaConfigsValueMaps, kafkaConfigsValueMap)
				}

				kafkaConfigsMap["kafkaConfigs"] = kafkaConfigsValueMaps
			}
		}

		kafkaConfigsJson, err := convertMaptoJsonString(kafkaConfigsMap)
		if err != nil {
			return WrapError(err)
		}

		request["KafkaConfigs"] = StringPointer(kafkaConfigsJson)
	}

	if v, ok := d.GetOk("pvtz_discovery_svc"); ok {
		pvtzDiscoverySvcMap := map[string]interface{}{}
		for _, pvtzDiscoverySvcList := range v.([]interface{}) {
			pvtzDiscoverySvcArg := pvtzDiscoverySvcList.(map[string]interface{})

			if serviceName, ok := pvtzDiscoverySvcArg["service_name"]; ok && serviceName.(string) != "" {
				pvtzDiscoverySvcMap["serviceName"] = serviceName
			}

			if namespaceId, ok := pvtzDiscoverySvcArg["namespace_id"]; ok && namespaceId.(string) != "" {
				pvtzDiscoverySvcMap["namespaceId"] = namespaceId
			}

			if enable, ok := pvtzDiscoverySvcArg["enable"]; ok {
				pvtzDiscoverySvcMap["enable"] = enable
			}

			if portProtocols, ok := pvtzDiscoverySvcArg["port_protocols"]; ok && len(portProtocols.([]interface{})) > 0 {
				portProtocolsMaps := make([]map[string]interface{}, 0)
				for _, portProtocolsList := range portProtocols.([]interface{}) {
					portProtocolsMap := map[string]interface{}{}
					portProtocolsArg := portProtocolsList.(map[string]interface{})

					if port, ok := portProtocolsArg["port"]; ok {
						portProtocolsMap["port"] = port
					}

					if protocol, ok := portProtocolsArg["protocol"]; ok && protocol.(string) != "" {
						portProtocolsMap["protocol"] = protocol
					}

					portProtocolsMaps = append(portProtocolsMaps, portProtocolsMap)
				}

				pvtzDiscoverySvcMap["portProtocols"] = portProtocolsMaps
			}
		}

		pvtzDiscoverySvcJson, err := convertMaptoJsonString(pvtzDiscoverySvcMap)
		if err != nil {
			return WrapError(err)
		}

		request["PvtzDiscoverySvc"] = StringPointer(pvtzDiscoverySvcJson)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_application", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["AppId"]))

	stateConf := BuildStateConf([]string{}, []string{"2", "8", "11", "12"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, saeService.SaeApplicationChangeOrderStateRefreshFunc(fmt.Sprint(responseData["ChangeOrderId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSaeApplicationUpdate(d, meta)
}

func resourceAliCloudSaeApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}

	object, err := saeService.DescribeSaeApplication(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_application saeService.DescribeSaeApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("app_name", object["AppName"])
	d.Set("package_type", object["PackageType"])
	d.Set("namespace_id", object["NamespaceId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("package_version", object["PackageVersion"])
	d.Set("package_url", object["PackageUrl"])
	d.Set("image_url", object["ImageUrl"])
	d.Set("command", object["Command"])
	d.Set("web_container", object["WebContainer"])
	d.Set("jdk", object["Jdk"])
	d.Set("jar_start_options", object["JarStartOptions"])
	d.Set("jar_start_args", object["JarStartArgs"])
	d.Set("app_description", object["AppDescription"])
	d.Set("batch_wait_time", object["BatchWaitTime"])
	d.Set("edas_container_version", object["EdasContainerVersion"])
	d.Set("enable_ahas", fmt.Sprint(object["EnableAhas"]))
	d.Set("enable_grey_tag_route", object["EnableGreyTagRoute"])
	d.Set("oss_ak_id", object["OssAkId"])
	d.Set("oss_ak_secret", object["OssAkSecret"])
	d.Set("php_arms_config_location", object["PhpArmsConfigLocation"])
	d.Set("php_config", object["PhpConfig"])
	d.Set("php_config_location", object["PhpConfigLocation"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("timezone", object["Timezone"])
	d.Set("war_start_options", object["WarStartOptions"])
	d.Set("acr_instance_id", object["AcrInstanceId"])
	d.Set("acr_assume_role_arn", object["AcrAssumeRoleArn"])
	d.Set("micro_registration", object["MicroRegistration"])
	d.Set("php", object["Php"])
	d.Set("image_pull_secrets", object["ImagePullSecrets"])
	d.Set("programming_language", object["ProgrammingLanguage"])
	d.Set("command_args", object["CommandArgs"])
	d.Set("envs", object["Envs"])
	d.Set("custom_host_alias", object["CustomHostAlias"])
	d.Set("liveness", object["Liveness"])
	d.Set("readiness", object["Readiness"])
	d.Set("post_start", object["PostStart"])
	d.Set("pre_stop", object["PreStop"])
	d.Set("sls_configs", object["SlsConfigs"])
	d.Set("tomcat_config", object["TomcatConfig"])
	d.Set("update_strategy", object["UpdateStrategy"])

	if v, ok := object["Replicas"]; ok && fmt.Sprint(v) != "0" {
		d.Set("replicas", formatInt(v))
	}

	if v, ok := object["Cpu"]; ok && fmt.Sprint(v) != "0" {
		d.Set("cpu", formatInt(v))
	}

	if v, ok := object["Memory"]; ok && fmt.Sprint(v) != "0" {
		d.Set("memory", formatInt(v))
	}

	if v, ok := object["MinReadyInstances"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_ready_instances", formatInt(v))
	}

	if v, ok := object["MinReadyInstanceRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_ready_instance_ratio", formatInt(v))
	}

	if v, ok := object["TerminationGracePeriodSeconds"]; ok && fmt.Sprint(v) != "0" {
		d.Set("termination_grace_period_seconds", formatInt(v))
	}

	if v, ok := object["OssMountDescs"].([]interface{}); ok {
		ossMountDescs, err := convertListObjectToCommaSeparate(v)
		if err != nil {
			return WrapError(err)
		}

		d.Set("oss_mount_descs", ossMountDescs)
	}

	if configMapMountDescList, ok := object["ConfigMapMountDesc"]; ok {
		configMapMountDescMaps := make([]map[string]interface{}, 0)
		for _, configMapMountDesc := range configMapMountDescList.([]interface{}) {
			configMapMountDescArg := configMapMountDesc.(map[string]interface{})
			configMapMountDescMap := map[string]interface{}{}

			if configMapId, ok := configMapMountDescArg["ConfigMapId"]; ok {
				configMapMountDescMap["configMapId"] = configMapId
			}

			if key, ok := configMapMountDescArg["Key"]; ok {
				configMapMountDescMap["key"] = key
			}

			if mountPath, ok := configMapMountDescArg["MountPath"]; ok {
				configMapMountDescMap["mountPath"] = mountPath
			}

			configMapMountDescMaps = append(configMapMountDescMaps, configMapMountDescMap)
		}

		configMapMountDescJson, err := convertListMapToJsonString(configMapMountDescMaps)
		if err != nil {
			return WrapError(err)
		}

		d.Set("config_map_mount_desc", configMapMountDescJson)
	}

	if v, ok := object["CommandArgs"].(string); ok && v != "" {
		commandArgsList, err := convertJsonStringToList(v)
		if err != nil {
			return WrapError(err)
		}
		d.Set("command_args_v2", commandArgsList)
	} else {
		d.Set("command_args_v2", nil)
	}

	if v, ok := object["CustomHostAlias"].(string); ok {
		customHostAliasList, err := convertJsonStringToList(v)
		if err != nil {
			return WrapError(err)
		}

		customHostAliasMaps := make([]map[string]interface{}, 0)
		for _, customHostAlias := range customHostAliasList {
			customHostAliasArg := customHostAlias.(map[string]interface{})
			customHostAliasMap := map[string]interface{}{}

			if hostName, ok := customHostAliasArg["hostName"]; ok {
				customHostAliasMap["host_name"] = hostName
			}

			if ip, ok := customHostAliasArg["ip"]; ok {
				customHostAliasMap["ip"] = ip
			}

			customHostAliasMaps = append(customHostAliasMaps, customHostAliasMap)
		}

		d.Set("custom_host_alias_v2", customHostAliasMaps)
	}

	if ossMountDescsList, ok := object["OssMountDescs"]; ok {
		ossMountDescsMaps := make([]map[string]interface{}, 0)
		for _, ossMountDescs := range ossMountDescsList.([]interface{}) {
			ossMountDescsArg := ossMountDescs.(map[string]interface{})
			ossMountDescsMap := map[string]interface{}{}

			if bucketName, ok := ossMountDescsArg["bucketName"]; ok {
				ossMountDescsMap["bucket_name"] = bucketName
			}

			if bucketPath, ok := ossMountDescsArg["bucketPath"]; ok {
				ossMountDescsMap["bucket_path"] = bucketPath
			}

			if mountPath, ok := ossMountDescsArg["mountPath"]; ok {
				ossMountDescsMap["mount_path"] = mountPath
			}

			if readOnly, ok := ossMountDescsArg["readOnly"]; ok {
				ossMountDescsMap["read_only"] = readOnly
			}

			ossMountDescsMaps = append(ossMountDescsMaps, ossMountDescsMap)
		}

		d.Set("oss_mount_descs_v2", ossMountDescsMaps)
	}

	if configMapMountDescList, ok := object["ConfigMapMountDesc"]; ok {
		configMapMountDescMaps := make([]map[string]interface{}, 0)
		for _, configMapMountDesc := range configMapMountDescList.([]interface{}) {
			configMapMountDescArg := configMapMountDesc.(map[string]interface{})
			configMapMountDescMap := map[string]interface{}{}

			if configMapId, ok := configMapMountDescArg["ConfigMapId"]; ok {
				configMapMountDescMap["config_map_id"] = configMapId
			}

			if mountPath, ok := configMapMountDescArg["MountPath"]; ok {
				configMapMountDescMap["mount_path"] = mountPath
			}

			if key, ok := configMapMountDescArg["Key"]; ok {
				configMapMountDescMap["key"] = key
			}

			configMapMountDescMaps = append(configMapMountDescMaps, configMapMountDescMap)
		}

		d.Set("config_map_mount_desc_v2", configMapMountDescMaps)
	}

	if v, ok := object["Liveness"].(string); ok && v != "" {
		livenessArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		livenessMaps := make([]map[string]interface{}, 0)
		livenessMap := map[string]interface{}{}

		if initialDelaySeconds, ok := livenessArg["initialDelaySeconds"]; ok {
			livenessMap["initial_delay_seconds"] = initialDelaySeconds
		}

		if periodSeconds, ok := livenessArg["periodSeconds"]; ok {
			livenessMap["period_seconds"] = periodSeconds
		}

		if timeoutSeconds, ok := livenessArg["timeoutSeconds"]; ok {
			livenessMap["timeout_seconds"] = timeoutSeconds
		}

		if exec, ok := livenessArg["exec"]; ok {
			execMaps := make([]map[string]interface{}, 0)
			execArg := exec.(map[string]interface{})
			execMap := map[string]interface{}{}

			if command, ok := execArg["command"]; ok {
				execMap["command"] = command
			}

			execMaps = append(execMaps, execMap)

			livenessMap["exec"] = execMaps
		}

		if tcpSocket, ok := livenessArg["tcpSocket"]; ok {
			tcpSocketMaps := make([]map[string]interface{}, 0)
			tcpSocketArg := tcpSocket.(map[string]interface{})
			tcpSocketMap := map[string]interface{}{}

			if port, ok := tcpSocketArg["port"]; ok {
				tcpSocketMap["port"] = port
			}

			tcpSocketMaps = append(tcpSocketMaps, tcpSocketMap)

			livenessMap["tcp_socket"] = tcpSocketMaps
		}

		if httpGet, ok := livenessArg["httpGet"]; ok {
			httpGetMaps := make([]map[string]interface{}, 0)
			httpGetArg := httpGet.(map[string]interface{})
			httpGetMap := map[string]interface{}{}

			if path, ok := httpGetArg["path"]; ok {
				httpGetMap["path"] = path
			}

			if port, ok := httpGetArg["port"]; ok {
				httpGetMap["port"] = port
			}

			if scheme, ok := httpGetArg["scheme"]; ok {
				httpGetMap["scheme"] = scheme
			}

			if keyWord, ok := httpGetArg["keyWord"]; ok {
				httpGetMap["key_word"] = keyWord
			}

			if isContainKeyWord, ok := httpGetArg["isContainKeyWord"]; ok {
				httpGetMap["is_contain_key_word"] = isContainKeyWord
			} else {
				httpGetMap["is_contain_key_word"] = nil
			}

			httpGetMaps = append(httpGetMaps, httpGetMap)

			livenessMap["http_get"] = httpGetMaps
		}

		livenessMaps = append(livenessMaps, livenessMap)

		d.Set("liveness_v2", livenessMaps)
	}

	if v, ok := object["Readiness"].(string); ok && v != "" {
		readinessArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		readinessMaps := make([]map[string]interface{}, 0)
		readinessMap := map[string]interface{}{}

		if initialDelaySeconds, ok := readinessArg["initialDelaySeconds"]; ok {
			readinessMap["initial_delay_seconds"] = initialDelaySeconds
		}

		if periodSeconds, ok := readinessArg["periodSeconds"]; ok {
			readinessMap["period_seconds"] = periodSeconds
		}

		if timeoutSeconds, ok := readinessArg["timeoutSeconds"]; ok {
			readinessMap["timeout_seconds"] = timeoutSeconds
		}

		if exec, ok := readinessArg["exec"]; ok {
			execMaps := make([]map[string]interface{}, 0)
			execArg := exec.(map[string]interface{})
			execMap := map[string]interface{}{}

			if command, ok := execArg["command"]; ok {
				execMap["command"] = command
			}

			execMaps = append(execMaps, execMap)

			readinessMap["exec"] = execMaps
		}

		if tcpSocket, ok := readinessArg["tcpSocket"]; ok {
			tcpSocketMaps := make([]map[string]interface{}, 0)
			tcpSocketArg := tcpSocket.(map[string]interface{})
			tcpSocketMap := map[string]interface{}{}

			if port, ok := tcpSocketArg["port"]; ok {
				tcpSocketMap["port"] = port
			}

			tcpSocketMaps = append(tcpSocketMaps, tcpSocketMap)

			readinessMap["tcp_socket"] = tcpSocketMaps
		}

		if httpGet, ok := readinessArg["httpGet"]; ok {
			httpGetMaps := make([]map[string]interface{}, 0)
			httpGetArg := httpGet.(map[string]interface{})
			httpGetMap := map[string]interface{}{}

			if path, ok := httpGetArg["path"]; ok {
				httpGetMap["path"] = path
			}

			if port, ok := httpGetArg["port"]; ok {
				httpGetMap["port"] = port
			}

			if scheme, ok := httpGetArg["scheme"]; ok {
				httpGetMap["scheme"] = scheme
			}

			if keyWord, ok := httpGetArg["keyWord"]; ok {
				httpGetMap["key_word"] = keyWord
			}

			if isContainKeyWord, ok := httpGetArg["isContainKeyWord"]; ok {
				httpGetMap["is_contain_key_word"] = isContainKeyWord
			} else {
				httpGetMap["is_contain_key_word"] = nil
			}

			httpGetMaps = append(httpGetMaps, httpGetMap)

			readinessMap["http_get"] = httpGetMaps
		}

		readinessMaps = append(readinessMaps, readinessMap)

		d.Set("readiness_v2", readinessMaps)
	}

	if v, ok := object["PostStart"].(string); ok && v != "" {
		postStartArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		postStartMaps := make([]map[string]interface{}, 0)
		postStartMap := map[string]interface{}{}

		if exec, ok := postStartArg["exec"]; ok {
			execMaps := make([]map[string]interface{}, 0)
			execArg := exec.(map[string]interface{})
			execMap := map[string]interface{}{}

			if command, ok := execArg["command"]; ok {
				execMap["command"] = command
			}

			execMaps = append(execMaps, execMap)

			postStartMap["exec"] = execMaps
		}

		postStartMaps = append(postStartMaps, postStartMap)

		d.Set("post_start_v2", postStartMaps)
	}

	if v, ok := object["PreStop"].(string); ok && v != "" {
		preStopArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		preStopMaps := make([]map[string]interface{}, 0)
		preStopMap := map[string]interface{}{}

		if exec, ok := preStopArg["exec"]; ok {
			execMaps := make([]map[string]interface{}, 0)
			execArg := exec.(map[string]interface{})
			execMap := map[string]interface{}{}

			if command, ok := execArg["command"]; ok {
				execMap["command"] = command
			}

			execMaps = append(execMaps, execMap)

			preStopMap["exec"] = execMaps
		}

		preStopMaps = append(preStopMaps, preStopMap)

		d.Set("pre_stop_v2", preStopMaps)
	}

	if v, ok := object["TomcatConfig"].(string); ok && v != "" {
		tomcatConfigArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		tomcatConfigMaps := make([]map[string]interface{}, 0)
		tomcatConfigMap := map[string]interface{}{}

		if port, ok := tomcatConfigArg["port"]; ok {
			tomcatConfigMap["port"] = port
		}

		if maxThreads, ok := tomcatConfigArg["maxThreads"]; ok {
			tomcatConfigMap["max_threads"] = maxThreads
		}

		if contextPath, ok := tomcatConfigArg["contextPath"]; ok {
			tomcatConfigMap["context_path"] = contextPath
		}

		if uriEncoding, ok := tomcatConfigArg["uriEncoding"]; ok {
			tomcatConfigMap["uri_encoding"] = uriEncoding
		}

		if useBodyEncodingForUri, ok := tomcatConfigArg["useBodyEncodingForUri"]; ok {
			tomcatConfigMap["use_body_encoding_for_uri"] = useBodyEncodingForUri
		}

		tomcatConfigMaps = append(tomcatConfigMaps, tomcatConfigMap)

		d.Set("tomcat_config_v2", tomcatConfigMaps)
	}

	if v, ok := object["UpdateStrategy"].(string); ok && v != "" {
		updateStrategyArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		updateStrategyMaps := make([]map[string]interface{}, 0)
		updateStrategyMap := map[string]interface{}{}

		if updateStrategyType, ok := updateStrategyArg["type"]; ok {
			updateStrategyMap["type"] = updateStrategyType
		}

		if batchUpdate, ok := updateStrategyArg["batchUpdate"]; ok {
			batchUpdateMaps := make([]map[string]interface{}, 0)
			batchUpdateArg := batchUpdate.(map[string]interface{})
			batchUpdateMap := map[string]interface{}{}

			if releaseType, ok := batchUpdateArg["releaseType"]; ok {
				batchUpdateMap["release_type"] = releaseType
			}

			if batch, ok := batchUpdateArg["batch"]; ok {
				batchUpdateMap["batch"] = batch
			}

			if batchWaitTime, ok := batchUpdateArg["batchWaitTime"]; ok {
				batchUpdateMap["batch_wait_time"] = batchWaitTime
			}

			batchUpdateMaps = append(batchUpdateMaps, batchUpdateMap)

			updateStrategyMap["batch_update"] = batchUpdateMaps
		}

		updateStrategyMaps = append(updateStrategyMaps, updateStrategyMap)

		d.Set("update_strategy_v2", updateStrategyMaps)
	}

	if v, ok := object["NasConfigs"].(string); ok {
		nasConfigsList, err := convertJsonStringToList(v)
		if err != nil {
			return WrapError(err)
		}

		nasConfigsMaps := make([]map[string]interface{}, 0)
		for _, nasConfigs := range nasConfigsList {
			nasConfigsArg := nasConfigs.(map[string]interface{})
			nasConfigsMap := map[string]interface{}{}

			if nasId, ok := nasConfigsArg["nasId"]; ok {
				nasConfigsMap["nas_id"] = nasId
			}

			if nasPath, ok := nasConfigsArg["nasPath"]; ok {
				nasConfigsMap["nas_path"] = nasPath
			}

			if mountPath, ok := nasConfigsArg["mountPath"]; ok {
				nasConfigsMap["mount_path"] = mountPath
			}

			if mountDomain, ok := nasConfigsArg["mountDomain"]; ok {
				nasConfigsMap["mount_domain"] = mountDomain
			}

			if readOnly, ok := nasConfigsArg["readOnly"]; ok {
				nasConfigsMap["read_only"] = readOnly
			}

			nasConfigsMaps = append(nasConfigsMaps, nasConfigsMap)
		}

		d.Set("nas_configs", nasConfigsMaps)
	}

	if v, ok := object["KafkaConfigs"].(string); ok && v != "" {
		kafkaConfigsArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		kafkaConfigsMaps := make([]map[string]interface{}, 0)
		kafkaConfigsMap := map[string]interface{}{}

		if kafkaInstanceId, ok := kafkaConfigsArg["kafkaInstanceId"]; ok {
			kafkaConfigsMap["kafka_instance_id"] = kafkaInstanceId
		}

		if kafkaEndpoint, ok := kafkaConfigsArg["kafkaEndpoint"]; ok {
			kafkaConfigsMap["kafka_endpoint"] = kafkaEndpoint
		}

		if kafkaConfigsValue, ok := kafkaConfigsArg["kafkaConfigs"]; ok {
			kafkaConfigsValueMaps := make([]map[string]interface{}, 0)
			for _, kafkaConfigsValueList := range kafkaConfigsValue.([]interface{}) {
				kafkaConfigsValueArg := kafkaConfigsValueList.(map[string]interface{})
				kafkaConfigsValueMap := map[string]interface{}{}

				if logType, ok := kafkaConfigsValueArg["logType"]; ok {
					kafkaConfigsValueMap["log_type"] = logType
				}

				if logDir, ok := kafkaConfigsValueArg["logDir"]; ok {
					kafkaConfigsValueMap["log_dir"] = logDir
				}

				if kafkaTopic, ok := kafkaConfigsValueArg["kafkaTopic"]; ok {
					kafkaConfigsValueMap["kafka_topic"] = kafkaTopic
				}

				kafkaConfigsValueMaps = append(kafkaConfigsValueMaps, kafkaConfigsValueMap)
			}

			kafkaConfigsMap["kafka_configs"] = kafkaConfigsValueMaps
		}

		kafkaConfigsMaps = append(kafkaConfigsMaps, kafkaConfigsMap)

		d.Set("kafka_configs", kafkaConfigsMaps)
	}

	if v, ok := object["PvtzDiscovery"].(string); ok && v != "" {
		pvtzDiscoveryArg, err := convertJsonStringToMap(v)
		if err != nil {
			return WrapError(err)
		}

		pvtzDiscoveryMaps := make([]map[string]interface{}, 0)
		pvtzDiscoveryMap := map[string]interface{}{}

		if serviceName, ok := pvtzDiscoveryArg["serviceName"]; ok {
			pvtzDiscoveryMap["service_name"] = serviceName
		}

		if namespaceId, ok := pvtzDiscoveryArg["namespaceId"]; ok {
			pvtzDiscoveryMap["namespace_id"] = namespaceId
		}

		if enable, ok := pvtzDiscoveryArg["enable"]; ok {
			enableBool, err := strconv.ParseBool(fmt.Sprint(enable))
			if err != nil {
				return WrapError(err)
			}

			pvtzDiscoveryMap["enable"] = enableBool
		}

		if portProtocols, ok := pvtzDiscoveryArg["portProtocols"]; ok {
			portProtocolsMaps := make([]map[string]interface{}, 0)
			for _, portProtocolsList := range portProtocols.([]interface{}) {
				portProtocolsArg := portProtocolsList.(map[string]interface{})
				portProtocolsMap := map[string]interface{}{}

				if port, ok := portProtocolsArg["port"]; ok {
					portProtocolsMap["port"] = port
				}

				if protocol, ok := portProtocolsArg["protocol"]; ok {
					portProtocolsMap["protocol"] = protocol
				}

				portProtocolsMaps = append(portProtocolsMaps, portProtocolsMap)
			}

			pvtzDiscoveryMap["port_protocols"] = portProtocolsMaps
		}

		pvtzDiscoveryMaps = append(pvtzDiscoveryMaps, pvtzDiscoveryMap)

		d.Set("pvtz_discovery_svc", pvtzDiscoveryMaps)
	}

	listTagResourcesObject, err := saeService.ListTagResources(d.Id(), "application")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	describeApplicationStatusObject, err := saeService.DescribeApplicationStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("status", describeApplicationStatusObject["CurrentStatus"])

	return nil
}

func resourceAliCloudSaeApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := saeService.SetResourceTags(d, "application"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	deployApplicationReq := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	if !d.IsNewResource() && d.HasChange("replicas") {
		update = true
	}
	deployApplicationReq["Replicas"] = StringPointer(strconv.Itoa(d.Get("replicas").(int)))

	if !d.IsNewResource() && d.HasChange("vswitch_id") {
		update = true
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		deployApplicationReq["VSwitchId"] = StringPointer(v.(string))
	}

	if d.HasChange("package_version") {
		update = true
	}
	if v, ok := d.GetOk("package_version"); ok {
		deployApplicationReq["PackageVersion"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("package_url") {
		update = true
	}
	if v, ok := d.GetOk("package_url"); ok {
		deployApplicationReq["PackageUrl"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("image_url") {
		update = true
	}
	if v, ok := d.GetOk("image_url"); ok {
		deployApplicationReq["ImageUrl"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("cpu") {
		update = true
	}
	if v, ok := d.GetOk("cpu"); ok {
		deployApplicationReq["Cpu"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if !d.IsNewResource() && d.HasChange("memory") {
		update = true
	}
	if v, ok := d.GetOk("memory"); ok {
		deployApplicationReq["Memory"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if !d.IsNewResource() && d.HasChange("command") {
		update = true
	}
	if v, ok := d.GetOk("command"); ok {
		deployApplicationReq["Command"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("web_container") {
		update = true
	}
	if v, ok := d.GetOk("web_container"); ok {
		deployApplicationReq["WebContainer"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("jdk") {
		update = true
	}
	if v, ok := d.GetOk("jdk"); ok {
		deployApplicationReq["Jdk"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("jar_start_options") {
		update = true
	}
	if v, ok := d.GetOk("jar_start_options"); ok {
		deployApplicationReq["JarStartOptions"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("jar_start_args") {
		update = true
	}
	if v, ok := d.GetOk("jar_start_args"); ok {
		deployApplicationReq["JarStartArgs"] = StringPointer(v.(string))
	}

	if d.HasChange("auto_enable_application_scaling_rule") {
		update = true
	}
	if v, ok := d.GetOk("auto_enable_application_scaling_rule"); ok {
		deployApplicationReq["AutoEnableApplicationScalingRule"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if d.HasChange("batch_wait_time") {
		update = true
	}
	if v, ok := d.GetOk("batch_wait_time"); ok {
		deployApplicationReq["BatchWaitTime"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if d.HasChange("change_order_desc") {
		update = true
	}
	if v, ok := d.GetOk("change_order_desc"); ok {
		deployApplicationReq["ChangeOrderDesc"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("edas_container_version") {
		update = true
	}
	if v, ok := d.GetOk("edas_container_version"); ok {
		deployApplicationReq["EdasContainerVersion"] = StringPointer(v.(string))
	}

	if d.HasChange("enable_ahas") {
		update = true

		if v, ok := d.GetOk("enable_ahas"); ok {
			deployApplicationReq["EnableAhas"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("enable_grey_tag_route") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_grey_tag_route"); ok {
		deployApplicationReq["EnableGreyTagRoute"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if d.HasChange("min_ready_instances") {
		update = true
	}
	if v, ok := d.GetOk("min_ready_instances"); ok {
		deployApplicationReq["MinReadyInstances"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if d.HasChange("min_ready_instance_ratio") {
		update = true
	}
	if v, ok := d.GetOk("min_ready_instance_ratio"); ok {
		deployApplicationReq["MinReadyInstanceRatio"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if !d.IsNewResource() && d.HasChange("oss_ak_id") {
		update = true
	}
	if v, ok := d.GetOk("oss_ak_id"); ok {
		deployApplicationReq["OssAkId"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("oss_ak_secret") {
		update = true
	}
	if v, ok := d.GetOk("oss_ak_secret"); ok {
		deployApplicationReq["OssAkSecret"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("php_arms_config_location") {
		update = true
	}
	if v, ok := d.GetOk("php_arms_config_location"); ok {
		deployApplicationReq["PhpArmsConfigLocation"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("php_config") {
		update = true
	}
	if v, ok := d.GetOk("php_config"); ok {
		deployApplicationReq["PhpConfig"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("php_config_location") {
		update = true
	}
	if v, ok := d.GetOk("php_config_location"); ok {
		deployApplicationReq["PhpConfigLocation"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("security_group_id") {
		update = true
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		deployApplicationReq["SecurityGroupId"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("termination_grace_period_seconds") {
		update = true
	}
	if v, ok := d.GetOk("termination_grace_period_seconds"); ok {
		deployApplicationReq["TerminationGracePeriodSeconds"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if !d.IsNewResource() && d.HasChange("timezone") {
		update = true
	}
	if v, ok := d.GetOk("timezone"); ok {
		deployApplicationReq["Timezone"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("war_start_options") {
		update = true
	}
	if v, ok := d.GetOk("war_start_options"); ok {
		deployApplicationReq["WarStartOptions"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("acr_instance_id") {
		update = true
	}
	if v, ok := d.GetOk("acr_instance_id"); ok {
		deployApplicationReq["AcrInstanceId"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("acr_assume_role_arn") {
		update = true
	}
	if v, ok := d.GetOk("acr_assume_role_arn"); ok {
		deployApplicationReq["AcrAssumeRoleArn"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("micro_registration") {
		update = true
	}
	if v, ok := d.GetOk("micro_registration"); ok {
		deployApplicationReq["MicroRegistration"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("envs") {
		update = true
	}
	if v, ok := d.GetOk("envs"); ok {
		deployApplicationReq["Envs"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("sls_configs") {
		update = true
	}
	if v, ok := d.GetOk("sls_configs"); ok {
		deployApplicationReq["SlsConfigs"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("php") {
		update = true
	}
	if v, ok := d.GetOk("php"); ok {
		deployApplicationReq["Php"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("image_pull_secrets") {
		update = true
	}
	if v, ok := d.GetOk("image_pull_secrets"); ok {
		deployApplicationReq["ImagePullSecrets"] = StringPointer(v.(string))
	}

	if !d.IsNewResource() && d.HasChange("command_args_v2") {
		update = true

		if v, ok := d.GetOk("command_args_v2"); ok {
			deployApplicationReq["CommandArgs"] = StringPointer(convertListToJsonString(v.([]interface{})))
		}
	}

	if !d.IsNewResource() && d.HasChange("custom_host_alias_v2") {
		update = true

		if v, ok := d.GetOk("custom_host_alias_v2"); ok {
			customHostAliasMaps := make([]map[string]interface{}, 0)
			for _, customHostAlias := range v.([]interface{}) {
				customHostAliasMap := map[string]interface{}{}
				customHostAliasArg := customHostAlias.(map[string]interface{})

				if hostName, ok := customHostAliasArg["host_name"]; ok && hostName.(string) != "" {
					customHostAliasMap["hostName"] = hostName
				}

				if ip, ok := customHostAliasArg["ip"]; ok && ip.(string) != "" {
					customHostAliasMap["ip"] = ip
				}

				customHostAliasMaps = append(customHostAliasMaps, customHostAliasMap)
			}

			customHostAliasJson, err := convertListMapToJsonString(customHostAliasMaps)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["CustomHostAlias"] = StringPointer(customHostAliasJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("oss_mount_descs_v2") {
		update = true

		if v, ok := d.GetOk("oss_mount_descs_v2"); ok {
			ossMountDescsMaps := make([]map[string]interface{}, 0)
			for _, ossMountDescs := range v.([]interface{}) {
				ossMountDescsMap := map[string]interface{}{}
				ossMountDescsArg := ossMountDescs.(map[string]interface{})

				if bucketName, ok := ossMountDescsArg["bucket_name"]; ok && bucketName.(string) != "" {
					ossMountDescsMap["bucketName"] = bucketName
				}

				if bucketPath, ok := ossMountDescsArg["bucket_path"]; ok && bucketPath.(string) != "" {
					ossMountDescsMap["bucketPath"] = bucketPath
				}

				if mountPath, ok := ossMountDescsArg["mount_path"]; ok && mountPath.(string) != "" {
					ossMountDescsMap["mountPath"] = mountPath
				}

				if readOnly, ok := ossMountDescsArg["read_only"]; ok {
					ossMountDescsMap["readOnly"] = readOnly
				}

				ossMountDescsMaps = append(ossMountDescsMaps, ossMountDescsMap)
			}

			ossMountDescsJson, err := convertListMapToJsonString(ossMountDescsMaps)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["OssMountDescs"] = StringPointer(ossMountDescsJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("config_map_mount_desc_v2") {
		update = true

		if v, ok := d.GetOk("config_map_mount_desc_v2"); ok {
			configMapMountDescMaps := make([]map[string]interface{}, 0)
			for _, configMapMountDesc := range v.([]interface{}) {
				configMapMountDescMap := map[string]interface{}{}
				configMapMountDescArg := configMapMountDesc.(map[string]interface{})

				if configMapId, ok := configMapMountDescArg["config_map_id"]; ok && configMapId.(string) != "" {
					configMapMountDescMap["configMapId"] = configMapId
				}

				if mountPath, ok := configMapMountDescArg["mount_path"]; ok && mountPath.(string) != "" {
					configMapMountDescMap["mountPath"] = mountPath
				}

				if key, ok := configMapMountDescArg["key"]; ok && key.(string) != "" {
					configMapMountDescMap["key"] = key
				}

				configMapMountDescMaps = append(configMapMountDescMaps, configMapMountDescMap)
			}

			configMapMountDescJson, err := convertListMapToJsonString(configMapMountDescMaps)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["ConfigMapMountDesc"] = StringPointer(configMapMountDescJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("liveness_v2") {
		update = true

		if v, ok := d.GetOk("liveness_v2"); ok {
			livenessMap := map[string]interface{}{}
			for _, livenessList := range v.([]interface{}) {
				livenessArg := livenessList.(map[string]interface{})

				if initialDelaySeconds, ok := livenessArg["initial_delay_seconds"]; ok {
					livenessMap["initialDelaySeconds"] = initialDelaySeconds
				}

				if periodSeconds, ok := livenessArg["period_seconds"]; ok {
					livenessMap["periodSeconds"] = periodSeconds
				}

				if timeoutSeconds, ok := livenessArg["timeout_seconds"]; ok {
					livenessMap["timeoutSeconds"] = timeoutSeconds
				}

				if exec, ok := livenessArg["exec"]; ok && len(exec.([]interface{})) > 0 {
					execMap := map[string]interface{}{}
					for _, execList := range exec.([]interface{}) {
						execArg := execList.(map[string]interface{})

						if command, ok := execArg["command"]; ok {
							execMap["command"] = command
						}
					}

					livenessMap["exec"] = execMap
				}

				if tcpSocket, ok := livenessArg["tcp_socket"]; ok && len(tcpSocket.([]interface{})) > 0 {
					tcpSocketMap := map[string]interface{}{}
					for _, tcpSocketList := range tcpSocket.([]interface{}) {
						tcpSocketArg := tcpSocketList.(map[string]interface{})

						if port, ok := tcpSocketArg["port"]; ok {
							tcpSocketMap["port"] = port
						}
					}

					livenessMap["tcpSocket"] = tcpSocketMap
				}

				if httpGet, ok := livenessArg["http_get"]; ok && len(httpGet.([]interface{})) > 0 {
					httpGetMap := map[string]interface{}{}
					for _, httpGetList := range httpGet.([]interface{}) {
						httpGetArg := httpGetList.(map[string]interface{})

						if path, ok := httpGetArg["path"]; ok && path.(string) != "" {
							httpGetMap["path"] = path
						}

						if port, ok := httpGetArg["port"]; ok {
							httpGetMap["port"] = port
						}

						if scheme, ok := httpGetArg["scheme"]; ok && scheme.(string) != "" {
							httpGetMap["scheme"] = scheme
						}

						if keyWord, ok := httpGetArg["key_word"]; ok && keyWord.(string) != "" {
							httpGetMap["keyWord"] = keyWord

							if isContainKeyWord, ok := d.GetOkExists("liveness_v2.0.http_get.0.is_contain_key_word"); ok {
								httpGetMap["isContainKeyWord"] = isContainKeyWord
							}
						}
					}

					livenessMap["httpGet"] = httpGetMap
				}
			}

			livenessJson, err := convertMaptoJsonString(livenessMap)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["Liveness"] = StringPointer(livenessJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("readiness_v2") {
		update = true

		if v, ok := d.GetOk("readiness_v2"); ok {
			readinessMap := map[string]interface{}{}
			for _, readinessList := range v.([]interface{}) {
				readinessArg := readinessList.(map[string]interface{})

				if initialDelaySeconds, ok := readinessArg["initial_delay_seconds"]; ok {
					readinessMap["initialDelaySeconds"] = initialDelaySeconds
				}

				if periodSeconds, ok := readinessArg["period_seconds"]; ok {
					readinessMap["periodSeconds"] = periodSeconds
				}

				if timeoutSeconds, ok := readinessArg["timeout_seconds"]; ok {
					readinessMap["timeoutSeconds"] = timeoutSeconds
				}

				if exec, ok := readinessArg["exec"]; ok && len(exec.([]interface{})) > 0 {
					execMap := map[string]interface{}{}
					for _, execList := range exec.([]interface{}) {
						execArg := execList.(map[string]interface{})

						if command, ok := execArg["command"]; ok {
							execMap["command"] = command
						}
					}

					readinessMap["exec"] = execMap
				}

				if tcpSocket, ok := readinessArg["tcp_socket"]; ok && len(tcpSocket.([]interface{})) > 0 {
					tcpSocketMap := map[string]interface{}{}
					for _, tcpSocketList := range tcpSocket.([]interface{}) {
						tcpSocketArg := tcpSocketList.(map[string]interface{})

						if port, ok := tcpSocketArg["port"]; ok {
							tcpSocketMap["port"] = port
						}
					}

					readinessMap["tcpSocket"] = tcpSocketMap
				}

				if httpGet, ok := readinessArg["http_get"]; ok && len(httpGet.([]interface{})) > 0 {
					httpGetMap := map[string]interface{}{}
					for _, httpGetList := range httpGet.([]interface{}) {
						httpGetArg := httpGetList.(map[string]interface{})

						if path, ok := httpGetArg["path"]; ok && path.(string) != "" {
							httpGetMap["path"] = path
						}

						if port, ok := httpGetArg["port"]; ok {
							httpGetMap["port"] = port
						}

						if scheme, ok := httpGetArg["scheme"]; ok && scheme.(string) != "" {
							httpGetMap["scheme"] = scheme
						}

						if keyWord, ok := httpGetArg["key_word"]; ok && keyWord.(string) != "" {
							httpGetMap["keyWord"] = keyWord

							if isContainKeyWord, ok := d.GetOkExists("readiness_v2.0.http_get.0.is_contain_key_word"); ok {
								httpGetMap["isContainKeyWord"] = isContainKeyWord
							}
						}
					}

					readinessMap["httpGet"] = httpGetMap
				}
			}

			readinessJson, err := convertMaptoJsonString(readinessMap)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["Readiness"] = StringPointer(readinessJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("post_start_v2") {
		update = true

		if v, ok := d.GetOk("post_start_v2"); ok {
			postStartMap := map[string]interface{}{}
			for _, postStartList := range v.([]interface{}) {
				postStartArg := postStartList.(map[string]interface{})

				if exec, ok := postStartArg["exec"]; ok && len(exec.([]interface{})) > 0 {
					execMap := map[string]interface{}{}
					for _, execList := range exec.([]interface{}) {
						execArg := execList.(map[string]interface{})

						if command, ok := execArg["command"]; ok {
							execMap["command"] = command
						}
					}

					postStartMap["exec"] = execMap
				}
			}

			postStartJson, err := convertMaptoJsonString(postStartMap)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["PostStart"] = StringPointer(postStartJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("pre_stop_v2") {
		update = true

		if v, ok := d.GetOk("pre_stop_v2"); ok {
			preStopMap := map[string]interface{}{}
			for _, preStopList := range v.([]interface{}) {
				preStopArg := preStopList.(map[string]interface{})

				if exec, ok := preStopArg["exec"]; ok && len(exec.([]interface{})) > 0 {
					execMap := map[string]interface{}{}
					for _, execList := range exec.([]interface{}) {
						execArg := execList.(map[string]interface{})

						if command, ok := execArg["command"]; ok {
							execMap["command"] = command
						}
					}

					preStopMap["exec"] = execMap
				}
			}

			preStopJson, err := convertMaptoJsonString(preStopMap)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["PreStop"] = StringPointer(preStopJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("tomcat_config_v2") {
		update = true

		if v, ok := d.GetOk("tomcat_config_v2"); ok {
			tomcatConfigMap := map[string]interface{}{}
			for _, tomcatConfigList := range v.([]interface{}) {
				tomcatConfigArg := tomcatConfigList.(map[string]interface{})

				if port, ok := tomcatConfigArg["port"]; ok {
					tomcatConfigMap["port"] = port
				}

				if maxThreads, ok := tomcatConfigArg["max_threads"]; ok {
					tomcatConfigMap["maxThreads"] = maxThreads
				}

				if contextPath, ok := tomcatConfigArg["context_path"]; ok && contextPath.(string) != "" {
					tomcatConfigMap["contextPath"] = contextPath
				}

				if uriEncoding, ok := tomcatConfigArg["uri_encoding"]; ok && uriEncoding.(string) != "" {
					tomcatConfigMap["uriEncoding"] = uriEncoding
				}

				if useBodyEncodingForUri, ok := tomcatConfigArg["use_body_encoding_for_uri"]; ok && useBodyEncodingForUri.(string) != "" {
					tomcatConfigMap["useBodyEncodingForUri"] = useBodyEncodingForUri
				}
			}

			tomcatConfigJson, err := convertMaptoJsonString(tomcatConfigMap)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["TomcatConfig"] = StringPointer(tomcatConfigJson)
		}
	}

	if d.HasChange("update_strategy_v2") {
		update = true

		if v, ok := d.GetOk("update_strategy_v2"); ok {
			updateStrategyMap := map[string]interface{}{}
			for _, updateStrategyList := range v.([]interface{}) {
				updateStrategyArg := updateStrategyList.(map[string]interface{})

				if updateStrategyType, ok := updateStrategyArg["type"]; ok && updateStrategyType.(string) != "" {
					updateStrategyMap["type"] = updateStrategyType
				}

				if batchUpdate, ok := updateStrategyArg["batch_update"]; ok && len(batchUpdate.([]interface{})) > 0 {
					batchUpdateMap := map[string]interface{}{}
					for _, batchUpdateList := range batchUpdate.([]interface{}) {
						batchUpdateArg := batchUpdateList.(map[string]interface{})

						if releaseType, ok := batchUpdateArg["release_type"]; ok && releaseType.(string) != "" {
							batchUpdateMap["releaseType"] = releaseType
						}

						if batch, ok := batchUpdateArg["batch"]; ok {
							batchUpdateMap["batch"] = batch
						}

						if batchWaitTime, ok := batchUpdateArg["batch_wait_time"]; ok {
							batchUpdateMap["batchWaitTime"] = batchWaitTime
						}
					}

					updateStrategyMap["batchUpdate"] = batchUpdateMap
				}
			}

			updateStrategyJson, err := convertMaptoJsonString(updateStrategyMap)
			if err != nil {
				return WrapError(err)
			}

			deployApplicationReq["UpdateStrategy"] = StringPointer(updateStrategyJson)
		}
	}

	if !d.IsNewResource() && d.HasChange("nas_configs") {
		update = true
	}
	if v, ok := d.GetOk("nas_configs"); ok {
		nasConfigsMaps := make([]map[string]interface{}, 0)
		for _, nasConfigs := range v.([]interface{}) {
			nasConfigsMap := map[string]interface{}{}
			nasConfigsArg := nasConfigs.(map[string]interface{})

			if nasId, ok := nasConfigsArg["nas_id"]; ok && nasId.(string) != "" {
				nasConfigsMap["nasId"] = nasId
			}

			if nasPath, ok := nasConfigsArg["nas_path"]; ok && nasPath.(string) != "" {
				nasConfigsMap["nasPath"] = nasPath
			}

			if mountPath, ok := nasConfigsArg["mount_path"]; ok && mountPath.(string) != "" {
				nasConfigsMap["mountPath"] = mountPath
			}

			if mountDomain, ok := nasConfigsArg["mount_domain"]; ok && mountDomain.(string) != "" {
				nasConfigsMap["mountDomain"] = mountDomain
			}

			if readOnly, ok := nasConfigsArg["read_only"]; ok {
				nasConfigsMap["readOnly"] = readOnly
			}

			nasConfigsMaps = append(nasConfigsMaps, nasConfigsMap)
		}

		nasConfigsJson, err := convertListMapToJsonString(nasConfigsMaps)
		if err != nil {
			return WrapError(err)
		}

		deployApplicationReq["NasConfigs"] = StringPointer(nasConfigsJson)
	}

	if !d.IsNewResource() && d.HasChange("kafka_configs") {
		update = true
	}
	if v, ok := d.GetOk("kafka_configs"); ok {
		kafkaConfigsMap := map[string]interface{}{}
		for _, kafkaConfigsList := range v.([]interface{}) {
			kafkaConfigsArg := kafkaConfigsList.(map[string]interface{})

			if kafkaInstanceId, ok := kafkaConfigsArg["kafka_instance_id"]; ok && kafkaInstanceId.(string) != "" {
				kafkaConfigsMap["kafkaInstanceId"] = kafkaInstanceId
			}

			if kafkaEndpoint, ok := kafkaConfigsArg["kafka_endpoint"]; ok && kafkaEndpoint.(string) != "" {
				kafkaConfigsMap["kafkaEndpoint"] = kafkaEndpoint
			}

			if kafkaConfigsValue, ok := kafkaConfigsArg["kafka_configs"]; ok && len(kafkaConfigsValue.([]interface{})) > 0 {
				kafkaConfigsValueMaps := make([]map[string]interface{}, 0)
				for _, kafkaConfigsValueList := range kafkaConfigsValue.([]interface{}) {
					kafkaConfigsValueMap := map[string]interface{}{}
					kafkaConfigsValueArg := kafkaConfigsValueList.(map[string]interface{})

					if logType, ok := kafkaConfigsValueArg["log_type"]; ok && logType.(string) != "" {
						kafkaConfigsValueMap["logType"] = logType
					}

					if logDir, ok := kafkaConfigsValueArg["log_dir"]; ok && logDir.(string) != "" {
						kafkaConfigsValueMap["logDir"] = logDir
					}

					if kafkaTopic, ok := kafkaConfigsValueArg["kafka_topic"]; ok && kafkaTopic.(string) != "" {
						kafkaConfigsValueMap["kafkaTopic"] = kafkaTopic
					}

					kafkaConfigsValueMaps = append(kafkaConfigsValueMaps, kafkaConfigsValueMap)
				}

				kafkaConfigsMap["kafkaConfigs"] = kafkaConfigsValueMaps
			}
		}

		kafkaConfigsJson, err := convertMaptoJsonString(kafkaConfigsMap)
		if err != nil {
			return WrapError(err)
		}

		deployApplicationReq["KafkaConfigs"] = StringPointer(kafkaConfigsJson)
	}

	if !d.IsNewResource() && d.HasChange("pvtz_discovery_svc") {
		update = true
	}
	if v, ok := d.GetOk("pvtz_discovery_svc"); ok {
		pvtzDiscoverySvcMap := map[string]interface{}{}
		for _, pvtzDiscoverySvcList := range v.([]interface{}) {
			pvtzDiscoverySvcArg := pvtzDiscoverySvcList.(map[string]interface{})

			if serviceName, ok := pvtzDiscoverySvcArg["service_name"]; ok && serviceName.(string) != "" {
				pvtzDiscoverySvcMap["serviceName"] = serviceName
			}

			if namespaceId, ok := pvtzDiscoverySvcArg["namespace_id"]; ok && namespaceId.(string) != "" {
				pvtzDiscoverySvcMap["namespaceId"] = namespaceId
			}

			if enable, ok := pvtzDiscoverySvcArg["enable"]; ok {
				pvtzDiscoverySvcMap["enable"] = enable
			}

			if portProtocols, ok := pvtzDiscoverySvcArg["port_protocols"]; ok && len(portProtocols.([]interface{})) > 0 {
				portProtocolsMaps := make([]map[string]interface{}, 0)
				for _, portProtocolsList := range portProtocols.([]interface{}) {
					portProtocolsMap := map[string]interface{}{}
					portProtocolsArg := portProtocolsList.(map[string]interface{})

					if port, ok := portProtocolsArg["port"]; ok {
						portProtocolsMap["port"] = port
					}

					if protocol, ok := portProtocolsArg["protocol"]; ok && protocol.(string) != "" {
						portProtocolsMap["protocol"] = protocol
					}

					portProtocolsMaps = append(portProtocolsMaps, portProtocolsMap)
				}

				pvtzDiscoverySvcMap["portProtocols"] = portProtocolsMaps
			}
		}

		pvtzDiscoverySvcJson, err := convertMaptoJsonString(pvtzDiscoverySvcMap)
		if err != nil {
			return WrapError(err)
		}

		deployApplicationReq["PvtzDiscoverySvc"] = StringPointer(pvtzDiscoverySvcJson)
	}

	if !d.IsNewResource() && d.HasChange("command_args") {
		update = true

		if v, ok := d.GetOk("command_args"); ok {
			deployApplicationReq["CommandArgs"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("custom_host_alias") {
		update = true

		if v, ok := d.GetOk("custom_host_alias"); ok {
			deployApplicationReq["CustomHostAlias"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("oss_mount_descs") {
		update = true

		if v, ok := d.GetOk("oss_mount_descs"); ok {
			deployApplicationReq["OssMountDescs"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("config_map_mount_desc") {
		update = true

		if v, ok := d.GetOk("config_map_mount_desc"); ok {
			deployApplicationReq["ConfigMapMountDesc"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("liveness") {
		update = true

		if v, ok := d.GetOk("liveness"); ok {
			deployApplicationReq["Liveness"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("readiness") {
		update = true

		if v, ok := d.GetOk("readiness"); ok {
			deployApplicationReq["Readiness"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("post_start") {
		update = true

		if v, ok := d.GetOk("post_start"); ok {
			deployApplicationReq["PostStart"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("pre_stop") {
		update = true

		if v, ok := d.GetOk("pre_stop"); ok {
			deployApplicationReq["PreStop"] = StringPointer(v.(string))
		}
	}

	if !d.IsNewResource() && d.HasChange("tomcat_config") {
		update = true

		if v, ok := d.GetOk("tomcat_config"); ok {
			deployApplicationReq["TomcatConfig"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("update_strategy") {
		update = true

		if v, ok := d.GetOk("update_strategy"); ok {
			deployApplicationReq["UpdateStrategy"] = StringPointer(v.(string))
		}
	}

	if update {
		action := "/pop/v1/sam/app/deployApplication"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RoaPost("sae", "2019-05-06", action, deployApplicationReq, nil, nil, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, deployApplicationReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
		responseData := response["Data"].(map[string]interface{})

		stateConf := BuildStateConf([]string{}, []string{"2", "8", "11", "12"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, saeService.SaeApplicationChangeOrderStateRefreshFunc(fmt.Sprint(responseData["ChangeOrderId"]), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("replicas")
		d.SetPartial("vswitch_id")
		d.SetPartial("package_version")
		d.SetPartial("package_url")
		d.SetPartial("image_url")
		d.SetPartial("cpu")
		d.SetPartial("memory")
		d.SetPartial("command")
		d.SetPartial("web_container")
		d.SetPartial("jdk")
		d.SetPartial("jar_start_options")
		d.SetPartial("jar_start_args")
		d.SetPartial("auto_enable_application_scaling_rule")
		d.SetPartial("batch_wait_time")
		d.SetPartial("change_order_desc")
		d.SetPartial("edas_container_version")
		d.SetPartial("enable_ahas")
		d.SetPartial("enable_grey_tag_route")
		d.SetPartial("min_ready_instances")
		d.SetPartial("min_ready_instance_ratio")
		d.SetPartial("oss_ak_id")
		d.SetPartial("oss_ak_secret")
		d.SetPartial("php_arms_config_location")
		d.SetPartial("php_config")
		d.SetPartial("php_config_location")
		d.SetPartial("security_group_id")
		d.SetPartial("termination_grace_period_seconds")
		d.SetPartial("timezone")
		d.SetPartial("war_start_options")
		d.SetPartial("acr_instance_id")
		d.SetPartial("acr_assume_role_arn")
		d.SetPartial("micro_registration")
		d.SetPartial("envs")
		d.SetPartial("sls_configs")
		d.SetPartial("php")
		d.SetPartial("image_pull_secrets")
		d.SetPartial("command_args_v2")
		d.SetPartial("custom_host_alias_v2")
		d.SetPartial("oss_mount_descs_v2")
		d.SetPartial("config_map_mount_desc_v2")
		d.SetPartial("liveness_v2")
		d.SetPartial("readiness_v2")
		d.SetPartial("post_start_v2")
		d.SetPartial("pre_stop_v2")
		d.SetPartial("tomcat_config_v2")
		d.SetPartial("update_strategy_v2")
		d.SetPartial("nas_configs")
		d.SetPartial("kafka_configs")
		d.SetPartial("pvtz_discovery_svc")
		d.SetPartial("command_args")
		d.SetPartial("custom_host_alias")
		d.SetPartial("oss_mount_descs")
		d.SetPartial("config_map_mount_desc")
		d.SetPartial("liveness")
		d.SetPartial("readiness")
		d.SetPartial("post_start")
		d.SetPartial("pre_stop")
		d.SetPartial("tomcat_config")
		d.SetPartial("update_strategy")
	}

	update = false
	updateApplicationDescriptionReq := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	if !d.IsNewResource() && d.HasChange("app_description") {
		update = true
	}
	if v, ok := d.GetOk("app_description"); ok {
		updateApplicationDescriptionReq["AppDescription"] = StringPointer(v.(string))
	}

	if update {
		action := "/pop/v1/sam/app/updateAppDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RoaPut("sae", "2019-05-06", action, updateApplicationDescriptionReq, nil, nil, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateApplicationDescriptionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("app_description")
	}

	//	【Exists】update status
	if d.HasChange("status") {
		object, err := saeService.DescribeApplicationStatus(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["CurrentStatus"].(string) != target {
			if target == "RUNNING" {
				request := map[string]*string{
					"AppId": StringPointer(d.Id()),
				}

				action := "/pop/v1/sam/app/startApplication"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}

			if target == "STOPPED" {
				request := map[string]*string{
					"AppId": StringPointer(d.Id()),
				}

				action := "/pop/v1/sam/app/stopApplication"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"Application.InvalidStatus"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})

				addDebug(action, response, request)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
				}
			}

			d.SetPartial("status")
		}
	}

	d.Partial(false)

	stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Id(), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSaeApplicationRead(d, meta)
}

func resourceAliCloudSaeApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/app/deleteApplication"
	var response map[string]interface{}

	var err error

	request := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}

	action = "/pop/v1/sam/app/describeApplicationConfig"
	request = map[string]*string{
		"AppId": StringPointer(d.Id()),
	}

	wait = incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(5*time.Minute), func() *resource.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if response != nil {
			err = fmt.Errorf("application have not been destroyed yet")
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapError(err)
	}

	return nil
}

func convertListObjectToCommaSeparate(configured []interface{}) (string, error) {
	if len(configured) < 1 {
		return "", nil
	}
	result := "["
	for i, v := range configured {
		rail := ","
		if i == len(configured)-1 {
			rail = ""
		}
		vv, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		result += string(vv) + rail
	}
	return result + "]", nil
}

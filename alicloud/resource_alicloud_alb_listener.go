// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbListenerCreate,
		Read:   resourceAliCloudAlbListenerRead,
		Update: resourceAliCloudAlbListenerUpdate,
		Delete: resourceAliCloudAlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_log_record_customized_headers_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_log_tracing_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tracing_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"Zipkin"}, false),
						},
						"tracing_sample": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 10000),
						},
						"tracing_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"ca_certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ca_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"default_actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"forward_group_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_group_tuples": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gzip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"http2_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 60),
			},
			"listener_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 65535),
			},
			"listener_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS", "QUIC"}, false),
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quic_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quic_listener_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"quic_upgrade_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 180),
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"x_forwarded_for_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"x_forwarded_for_client_source_ips_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"x_forwarded_for_host_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"x_forwarded_for_client_source_ips_trusted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"x_forwarded_for_client_cert_subject_dn_alias": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-z-_]{1,40}"), "The Custom Header Field Name,"),
						},
						"x_forwarded_for_slb_id_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_client_cert_issuer_dn_alias": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-z-_]{1,40}"), "The Custom Header Field Names Only When xforwardedforclientcertsubjectdnenabled, Which Evaluates to True When the Entry into Force of."),
						},
						"x_forwarded_for_client_cert_client_verify_alias": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-z-_]{1,40}"), "The Custom Header Field Names Only When xforwardedforclientcertclientverifyenabled Has a Value of True, this Value Will Not Take Effect until."),
						},
						"x_forwarded_for_client_cert_finger_print_alias": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-z-_]{1,40}"), "The Custom Header Field Names Only When xforwardedforclientcertfingerprintenabled, Which Evaluates to True When the Entry into Force of."),
						},
						"x_forwarded_for_proto_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_slb_port_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_client_src_port_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_processing_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_client_cert_client_verify_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_client_cert_subject_dn_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_client_cert_finger_print_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"x_forwarded_for_client_cert_issuer_dn_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"xforwarded_for_config": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				MaxItems:      1,
				ConflictsWith: []string{"x_forwarded_for_config"},
				Deprecated:    "Field 'xforwarded_for_config' has been deprecated from provider version 1.161.0. Use 'x_forwarded_for_config' instead.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("listener_protocol"); ok && v.(string) == "HTTPS" {
						return false
					}
					return true
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"xforwardedforclientcert_issuerdnalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcert_issuerdnenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientcertclientverifyalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcertclientverifyenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientcertfingerprintalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcertfingerprintenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientcertsubjectdnalias": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"xforwardedforclientcertsubjectdnenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforclientsrcportenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforprotoenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforslbidenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"xforwardedforslbportenabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"acl_config": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				Deprecated: "Field 'acl_config' has been deprecated from provider version 1.163.0 and it will be removed in the future version. Please use the new resource 'alicloud_alb_listener_acl_attachment'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_relations": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"acl_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"acl_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"White", "Black"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudAlbListenerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("certificates"); ok {
		certificatesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["CertificateId"] = dataLoopTmp["certificate_id"]
			certificatesMapsArray = append(certificatesMapsArray, dataLoopMap)
		}
		request["Certificates"] = certificatesMapsArray
	}

	request["LoadBalancerId"] = d.Get("load_balancer_id")
	request["ListenerProtocol"] = d.Get("listener_protocol")
	request["ListenerPort"] = d.Get("listener_port")
	if v, ok := d.GetOkExists("request_timeout"); ok && v.(int) > 0 {
		request["RequestTimeout"] = v
	}
	if v, ok := d.GetOkExists("gzip_enabled"); ok {
		request["GzipEnabled"] = v
	}
	if v, ok := d.GetOkExists("http2_enabled"); ok {
		request["Http2Enabled"] = v
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request["SecurityPolicyId"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("quic_config"); !IsNil(v) {
		quicListenerId1, _ := jsonpath.Get("$[0].quic_listener_id", v)
		if quicListenerId1 != nil && quicListenerId1 != "" {
			objectDataLocalMap["QuicListenerId"] = quicListenerId1
		}
		quicUpgradeEnabled1, _ := jsonpath.Get("$[0].quic_upgrade_enabled", v)
		if quicUpgradeEnabled1 != nil && quicUpgradeEnabled1 != "" {
			objectDataLocalMap["QuicUpgradeEnabled"] = quicUpgradeEnabled1
		}

		request["QuicConfig"] = objectDataLocalMap
	}

	if v, ok := d.GetOkExists("idle_timeout"); ok && v.(int) > 0 {
		request["IdleTimeout"] = v
	}
	if v, ok := d.GetOk("default_actions"); ok {
		defaultActionsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Type"] = dataLoop1Tmp["type"]
			localData2 := make(map[string]interface{})
			if v, ok := dataLoop1Tmp["forward_group_config"]; ok {
				localData3, err := jsonpath.Get("$[0].server_group_tuples", v)
				if err != nil {
					localData3 = make([]interface{}, 0)
				}
				localMaps1 := make([]interface{}, 0)
				for _, dataLoop3 := range localData3.([]interface{}) {
					dataLoop3Tmp := make(map[string]interface{})
					if dataLoop3 != nil {
						dataLoop3Tmp = dataLoop3.(map[string]interface{})
					}
					dataLoop3Map := make(map[string]interface{})
					dataLoop3Map["ServerGroupId"] = dataLoop3Tmp["server_group_id"]
					localMaps1 = append(localMaps1, dataLoop3Map)
				}
				localData2["ServerGroupTuples"] = localMaps1
			}

			dataLoop1Map["ForwardGroupConfig"] = localData2
			defaultActionsMapsArray = append(defaultActionsMapsArray, dataLoop1Map)
		}
		request["DefaultActions"] = defaultActionsMapsArray
	}

	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request["CaEnabled"] = v
	}
	if v, ok := d.GetOk("ca_certificates"); ok {
		caCertificatesMapsArray := make([]interface{}, 0)
		for _, dataLoop4 := range v.([]interface{}) {
			dataLoop4Tmp := dataLoop4.(map[string]interface{})
			dataLoop4Map := make(map[string]interface{})
			dataLoop4Map["CertificateId"] = dataLoop4Tmp["certificate_id"]
			caCertificatesMapsArray = append(caCertificatesMapsArray, dataLoop4Map)
		}
		request["CaCertificates"] = caCertificatesMapsArray
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("listener_description"); ok {
		request["ListenerDescription"] = v
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("x_forwarded_for_config"); !IsNil(v) {
		xForwardedForClientCertClientVerifyAlias1, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_client_verify_alias", v)
		if xForwardedForClientCertClientVerifyAlias1 != nil && xForwardedForClientCertClientVerifyAlias1 != "" {
			objectDataLocalMap1["XForwardedForClientCertClientVerifyAlias"] = xForwardedForClientCertClientVerifyAlias1
		}
		xForwardedForClientCertClientVerifyEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_client_verify_enabled", v)
		if xForwardedForClientCertClientVerifyEnabled1 != nil && xForwardedForClientCertClientVerifyEnabled1 != "" {
			objectDataLocalMap1["XForwardedForClientCertClientVerifyEnabled"] = xForwardedForClientCertClientVerifyEnabled1
		}
		xForwardedForClientCertFingerPrintAlias, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_finger_print_alias", v)
		if xForwardedForClientCertFingerPrintAlias != nil && xForwardedForClientCertFingerPrintAlias != "" {
			objectDataLocalMap1["XForwardedForClientCertFingerprintAlias"] = xForwardedForClientCertFingerPrintAlias
		}
		xForwardedForClientCertFingerPrintEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_finger_print_enabled", v)
		if xForwardedForClientCertFingerPrintEnabled != nil && xForwardedForClientCertFingerPrintEnabled != "" {
			objectDataLocalMap1["XForwardedForClientCertFingerprintEnabled"] = xForwardedForClientCertFingerPrintEnabled
		}
		xForwardedForClientCertIssuerDnAlias, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_issuer_dn_alias", v)
		if xForwardedForClientCertIssuerDnAlias != nil && xForwardedForClientCertIssuerDnAlias != "" {
			objectDataLocalMap1["XForwardedForClientCertIssuerDNAlias"] = xForwardedForClientCertIssuerDnAlias
		}
		xForwardedForClientCertIssuerDnEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_issuer_dn_enabled", v)
		if xForwardedForClientCertIssuerDnEnabled != nil && xForwardedForClientCertIssuerDnEnabled != "" {
			objectDataLocalMap1["XForwardedForClientCertIssuerDNEnabled"] = xForwardedForClientCertIssuerDnEnabled
		}
		xForwardedForClientCertSubjectDnAlias, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_subject_dn_alias", v)
		if xForwardedForClientCertSubjectDnAlias != nil && xForwardedForClientCertSubjectDnAlias != "" {
			objectDataLocalMap1["XForwardedForClientCertSubjectDNAlias"] = xForwardedForClientCertSubjectDnAlias
		}
		xForwardedForClientCertSubjectDnEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_subject_dn_enabled", v)
		if xForwardedForClientCertSubjectDnEnabled != nil && xForwardedForClientCertSubjectDnEnabled != "" {
			objectDataLocalMap1["XForwardedForClientCertSubjectDNEnabled"] = xForwardedForClientCertSubjectDnEnabled
		}
		xForwardedForClientSrcPortEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_client_src_port_enabled", v)
		if xForwardedForClientSrcPortEnabled1 != nil && xForwardedForClientSrcPortEnabled1 != "" {
			objectDataLocalMap1["XForwardedForClientSrcPortEnabled"] = xForwardedForClientSrcPortEnabled1
		}
		xForwardedForEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_enabled", v)
		if xForwardedForEnabled1 != nil && xForwardedForEnabled1 != "" {
			objectDataLocalMap1["XForwardedForEnabled"] = xForwardedForEnabled1
		}
		xForwardedForProtoEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_proto_enabled", v)
		if xForwardedForProtoEnabled1 != nil && xForwardedForProtoEnabled1 != "" {
			objectDataLocalMap1["XForwardedForProtoEnabled"] = xForwardedForProtoEnabled1
		}
		xForwardedForSlbIdEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_slb_id_enabled", v)
		if xForwardedForSlbIdEnabled != nil && xForwardedForSlbIdEnabled != "" {
			objectDataLocalMap1["XForwardedForSLBIdEnabled"] = xForwardedForSlbIdEnabled
		}
		xForwardedForSlbPortEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_slb_port_enabled", v)
		if xForwardedForSlbPortEnabled != nil && xForwardedForSlbPortEnabled != "" {
			objectDataLocalMap1["XForwardedForSLBPortEnabled"] = xForwardedForSlbPortEnabled
		}
		xForwardedForClientSourceIpsEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_client_source_ips_enabled", v)
		if xForwardedForClientSourceIpsEnabled1 != nil && xForwardedForClientSourceIpsEnabled1 != "" {
			objectDataLocalMap1["XForwardedForClientSourceIpsEnabled"] = xForwardedForClientSourceIpsEnabled1
		}
		xForwardedForProcessingMode1, _ := jsonpath.Get("$[0].x_forwarded_for_processing_mode", v)
		if xForwardedForProcessingMode1 != nil && xForwardedForProcessingMode1 != "" {
			objectDataLocalMap1["XForwardedForProcessingMode"] = xForwardedForProcessingMode1
		}
		xForwardedForClientSourceIpsTrusted1, _ := jsonpath.Get("$[0].x_forwarded_for_client_source_ips_trusted", v)
		if xForwardedForClientSourceIpsTrusted1 != nil && xForwardedForClientSourceIpsTrusted1 != "" {
			objectDataLocalMap1["XForwardedForClientSourceIpsTrusted"] = xForwardedForClientSourceIpsTrusted1
		}
		xForwardedForHostEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_host_enabled", v)
		if xForwardedForHostEnabled1 != nil && xForwardedForHostEnabled1 != "" {
			objectDataLocalMap1["XForwardedForHostEnabled"] = xForwardedForHostEnabled1
		}

		request["XForwardedForConfig"] = objectDataLocalMap1
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("xforwarded_for_config"); ok {
		xforwardedForConfigMap := map[string]interface{}{}
		for _, xforwardedForConfig := range v.(*schema.Set).List() {
			xforwardedForConfigArg := xforwardedForConfig.(map[string]interface{})
			xforwardedForConfigMap["XForwardedForClientCertIssuerDNAlias"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnalias"]
			xforwardedForConfigMap["XForwardedForClientCertIssuerDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnenabled"]
			xforwardedForConfigMap["XForwardedForClientCertClientVerifyAlias"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyalias"]
			xforwardedForConfigMap["XForwardedForClientCertClientVerifyEnabled"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyenabled"]
			xforwardedForConfigMap["XForwardedForClientCertFingerprintAlias"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintalias"]
			xforwardedForConfigMap["XForwardedForClientCertFingerprintEnabled"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintenabled"]
			xforwardedForConfigMap["XForwardedForClientCertSubjectDNAlias"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnalias"]
			xforwardedForConfigMap["XForwardedForClientCertSubjectDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnenabled"]
			xforwardedForConfigMap["XForwardedForClientSrcPortEnabled"] = xforwardedForConfigArg["xforwardedforclientsrcportenabled"]
			xforwardedForConfigMap["XForwardedForEnabled"] = xforwardedForConfigArg["xforwardedforenabled"]
			xforwardedForConfigMap["XForwardedForProtoEnabled"] = xforwardedForConfigArg["xforwardedforprotoenabled"]
			xforwardedForConfigMap["XForwardedForSLBIdEnabled"] = xforwardedForConfigArg["xforwardedforslbidenabled"]
			xforwardedForConfigMap["XForwardedForSLBPortEnabled"] = xforwardedForConfigArg["xforwardedforslbportenabled"]
		}

		request["XForwardedForConfig"] = xforwardedForConfigMap
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing", "IncorrectBusinessStatus.LoadBalancer", "-21020"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ListenerId"]))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albServiceV2.DescribeAsyncAlbListenerStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudAlbListenerUpdate(d, meta)
}

func resourceAliCloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_listener DescribeAlbListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CaEnabled"] != nil {
		d.Set("ca_enabled", objectRaw["CaEnabled"])
	}
	if objectRaw["GzipEnabled"] != nil {
		d.Set("gzip_enabled", objectRaw["GzipEnabled"])
	}
	if objectRaw["Http2Enabled"] != nil {
		d.Set("http2_enabled", objectRaw["Http2Enabled"])
	}
	if objectRaw["IdleTimeout"] != nil {
		d.Set("idle_timeout", objectRaw["IdleTimeout"])
	}
	if objectRaw["ListenerDescription"] != nil {
		d.Set("listener_description", objectRaw["ListenerDescription"])
	}
	if objectRaw["ListenerPort"] != nil {
		d.Set("listener_port", objectRaw["ListenerPort"])
	}
	if objectRaw["ListenerProtocol"] != nil {
		d.Set("listener_protocol", objectRaw["ListenerProtocol"])
	}
	if objectRaw["LoadBalancerId"] != nil {
		d.Set("load_balancer_id", objectRaw["LoadBalancerId"])
	}
	if objectRaw["RequestTimeout"] != nil {
		d.Set("request_timeout", objectRaw["RequestTimeout"])
	}
	if objectRaw["SecurityPolicyId"] != nil {
		d.Set("security_policy_id", objectRaw["SecurityPolicyId"])
	}
	if objectRaw["ListenerStatus"] != nil {
		d.Set("status", objectRaw["ListenerStatus"])
	}

	logConfig1RawObj, _ := jsonpath.Get("$.LogConfig", objectRaw)
	logConfig1Raw := make(map[string]interface{})
	if logConfig1RawObj != nil {
		logConfig1Raw = logConfig1RawObj.(map[string]interface{})
	}
	if logConfig1Raw["AccessLogRecordCustomizedHeadersEnabled"] != nil {
		d.Set("access_log_record_customized_headers_enabled", logConfig1Raw["AccessLogRecordCustomizedHeadersEnabled"])
	}

	accessLogTracingConfigMaps := make([]map[string]interface{}, 0)
	accessLogTracingConfigMap := make(map[string]interface{})
	accessLogTracingConfig1RawObj, _ := jsonpath.Get("$.LogConfig.AccessLogTracingConfig", objectRaw)
	accessLogTracingConfig1Raw := make(map[string]interface{})
	if accessLogTracingConfig1RawObj != nil {
		accessLogTracingConfig1Raw = accessLogTracingConfig1RawObj.(map[string]interface{})
	}
	if len(accessLogTracingConfig1Raw) > 0 {
		accessLogTracingConfigMap["tracing_enabled"] = accessLogTracingConfig1Raw["TracingEnabled"]
		accessLogTracingConfigMap["tracing_sample"] = accessLogTracingConfig1Raw["TracingSample"]
		accessLogTracingConfigMap["tracing_type"] = accessLogTracingConfig1Raw["TracingType"]

		accessLogTracingConfigMaps = append(accessLogTracingConfigMaps, accessLogTracingConfigMap)
	}
	if accessLogTracingConfig1RawObj != nil {
		if err := d.Set("access_log_tracing_config", accessLogTracingConfigMaps); err != nil {
			return err
		}
	}
	caCertificates1Raw := objectRaw["CaCertificates"]
	caCertificatesMaps := make([]map[string]interface{}, 0)
	if caCertificates1Raw != nil {
		for _, caCertificatesChild1Raw := range caCertificates1Raw.([]interface{}) {
			caCertificatesMap := make(map[string]interface{})
			caCertificatesChild1Raw := caCertificatesChild1Raw.(map[string]interface{})
			caCertificatesMap["certificate_id"] = caCertificatesChild1Raw["CertificateId"]

			caCertificatesMaps = append(caCertificatesMaps, caCertificatesMap)
		}
	}
	if objectRaw["CaCertificates"] != nil {
		if err := d.Set("ca_certificates", caCertificatesMaps); err != nil {
			return err
		}
	}
	certificates1Raw := objectRaw["Certificates"]
	certificatesMaps := make([]map[string]interface{}, 0)
	if certificates1Raw != nil {
		for _, certificatesChild1Raw := range certificates1Raw.([]interface{}) {
			certificatesMap := make(map[string]interface{})
			certificatesChild1Raw := certificatesChild1Raw.(map[string]interface{})
			certificatesMap["certificate_id"] = certificatesChild1Raw["CertificateId"]

			certificatesMaps = append(certificatesMaps, certificatesMap)
		}
	}
	if objectRaw["Certificates"] != nil {
		if err := d.Set("certificates", certificatesMaps); err != nil {
			return err
		}
	}
	defaultActions1Raw := objectRaw["DefaultActions"]
	defaultActionsMaps := make([]map[string]interface{}, 0)
	if defaultActions1Raw != nil {
		for _, defaultActionsChild1Raw := range defaultActions1Raw.([]interface{}) {
			defaultActionsMap := make(map[string]interface{})
			defaultActionsChild1Raw := defaultActionsChild1Raw.(map[string]interface{})
			defaultActionsMap["type"] = defaultActionsChild1Raw["Type"]

			forwardGroupConfigMaps := make([]map[string]interface{}, 0)
			forwardGroupConfigMap := make(map[string]interface{})
			serverGroupTuples1Raw, _ := jsonpath.Get("$.ForwardGroupConfig.ServerGroupTuples", defaultActionsChild1Raw)

			serverGroupTuplesMaps := make([]map[string]interface{}, 0)
			if serverGroupTuples1Raw != nil {
				for _, serverGroupTuplesChild1Raw := range serverGroupTuples1Raw.([]interface{}) {
					serverGroupTuplesMap := make(map[string]interface{})
					serverGroupTuplesChild1Raw := serverGroupTuplesChild1Raw.(map[string]interface{})
					serverGroupTuplesMap["server_group_id"] = serverGroupTuplesChild1Raw["ServerGroupId"]

					serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
				}
			}
			forwardGroupConfigMap["server_group_tuples"] = serverGroupTuplesMaps
			forwardGroupConfigMaps = append(forwardGroupConfigMaps, forwardGroupConfigMap)
			defaultActionsMap["forward_group_config"] = forwardGroupConfigMaps
			defaultActionsMaps = append(defaultActionsMaps, defaultActionsMap)
		}
	}
	if objectRaw["DefaultActions"] != nil {
		if err := d.Set("default_actions", defaultActionsMaps); err != nil {
			return err
		}
	}
	quicConfigMaps := make([]map[string]interface{}, 0)
	quicConfigMap := make(map[string]interface{})
	quicConfig1Raw := make(map[string]interface{})
	if objectRaw["QuicConfig"] != nil {
		quicConfig1Raw = objectRaw["QuicConfig"].(map[string]interface{})
	}
	if len(quicConfig1Raw) > 0 {
		quicConfigMap["quic_listener_id"] = quicConfig1Raw["QuicListenerId"]
		quicConfigMap["quic_upgrade_enabled"] = quicConfig1Raw["QuicUpgradeEnabled"]

		quicConfigMaps = append(quicConfigMaps, quicConfigMap)
	}
	if objectRaw["QuicConfig"] != nil {
		if err := d.Set("quic_config", quicConfigMaps); err != nil {
			return err
		}
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	xForwardedForConfigMaps := make([]map[string]interface{}, 0)
	xForwardedForConfigMap := make(map[string]interface{})
	xForwardedForConfig1Raw := make(map[string]interface{})
	if objectRaw["XForwardedForConfig"] != nil {
		xForwardedForConfig1Raw = objectRaw["XForwardedForConfig"].(map[string]interface{})
	}
	if len(xForwardedForConfig1Raw) > 0 {
		xForwardedForConfigMap["x_forwarded_for_client_cert_client_verify_alias"] = xForwardedForConfig1Raw["XForwardedForClientCertClientVerifyAlias"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_client_verify_enabled"] = xForwardedForConfig1Raw["XForwardedForClientCertClientVerifyEnabled"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_finger_print_alias"] = xForwardedForConfig1Raw["XForwardedForClientCertFingerprintAlias"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_finger_print_enabled"] = xForwardedForConfig1Raw["XForwardedForClientCertFingerprintEnabled"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_issuer_dn_alias"] = xForwardedForConfig1Raw["XForwardedForClientCertIssuerDNAlias"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_issuer_dn_enabled"] = xForwardedForConfig1Raw["XForwardedForClientCertIssuerDNEnabled"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_subject_dn_alias"] = xForwardedForConfig1Raw["XForwardedForClientCertSubjectDNAlias"]
		xForwardedForConfigMap["x_forwarded_for_client_cert_subject_dn_enabled"] = xForwardedForConfig1Raw["XForwardedForClientCertSubjectDNEnabled"]
		xForwardedForConfigMap["x_forwarded_for_client_source_ips_enabled"] = xForwardedForConfig1Raw["XForwardedForClientSourceIpsEnabled"]
		xForwardedForConfigMap["x_forwarded_for_client_source_ips_trusted"] = xForwardedForConfig1Raw["XForwardedForClientSourceIpsTrusted"]
		xForwardedForConfigMap["x_forwarded_for_client_src_port_enabled"] = xForwardedForConfig1Raw["XForwardedForClientSrcPortEnabled"]
		xForwardedForConfigMap["x_forwarded_for_enabled"] = xForwardedForConfig1Raw["XForwardedForEnabled"]
		xForwardedForConfigMap["x_forwarded_for_host_enabled"] = xForwardedForConfig1Raw["XForwardedForHostEnabled"]
		xForwardedForConfigMap["x_forwarded_for_processing_mode"] = xForwardedForConfig1Raw["XForwardedForProcessingMode"]
		xForwardedForConfigMap["x_forwarded_for_proto_enabled"] = xForwardedForConfig1Raw["XForwardedForProtoEnabled"]
		xForwardedForConfigMap["x_forwarded_for_slb_id_enabled"] = xForwardedForConfig1Raw["XForwardedForSLBIdEnabled"]
		xForwardedForConfigMap["x_forwarded_for_slb_port_enabled"] = xForwardedForConfig1Raw["XForwardedForSLBPortEnabled"]

		xForwardedForConfigMaps = append(xForwardedForConfigMaps, xForwardedForConfigMap)
	}
	if objectRaw["XForwardedForConfig"] != nil {
		if err := d.Set("x_forwarded_for_config", xForwardedForConfigMaps); err != nil {
			return err
		}
	}

	if xforwardedForConfig, ok := objectRaw["XForwardedForConfig"]; ok && len(xforwardedForConfig.(map[string]interface{})) > 0 {
		xforwardedForConfigSli := make([]map[string]interface{}, 0)
		xforwardedForConfigMap := make(map[string]interface{})
		xforwardedForConfigMap["xforwardedforclientcert_issuerdnalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertIssuerDNAlias"]
		xforwardedForConfigMap["xforwardedforclientcert_issuerdnenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertIssuerDNEnabled"]
		xforwardedForConfigMap["xforwardedforclientcertclientverifyalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertClientVerifyAlias"]
		xforwardedForConfigMap["xforwardedforclientcertclientverifyenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertClientVerifyEnabled"]
		xforwardedForConfigMap["xforwardedforclientcertfingerprintalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertFingerprintAlias"]
		xforwardedForConfigMap["xforwardedforclientcertfingerprintenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertFingerprintEnabled"]
		xforwardedForConfigMap["xforwardedforclientcertsubjectdnalias"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertSubjectDNAlias"]
		xforwardedForConfigMap["xforwardedforclientcertsubjectdnenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientCertSubjectDNEnabled"]
		xforwardedForConfigMap["xforwardedforclientsrcportenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForClientSrcPortEnabled"]
		xforwardedForConfigMap["xforwardedforenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForEnabled"]
		xforwardedForConfigMap["xforwardedforprotoenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForProtoEnabled"]
		xforwardedForConfigMap["xforwardedforslbidenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForSLBIdEnabled"]
		xforwardedForConfigMap["xforwardedforslbportenabled"] = xforwardedForConfig.(map[string]interface{})["XForwardedForSLBPortEnabled"]
		xforwardedForConfigSli = append(xforwardedForConfigSli, xforwardedForConfigMap)
		d.Set("xforwarded_for_config", xforwardedForConfigSli)
	}

	aclConfigSli := make([]map[string]interface{}, 0)
	if aclConfig, ok := objectRaw["AclConfig"]; ok && len(aclConfig.(map[string]interface{})) > 0 {
		aclConfigMap := make(map[string]interface{})

		aclRelationsSli := make([]map[string]interface{}, 0)
		if v, ok := aclConfig.(map[string]interface{})["AclRelations"]; ok && len(v.([]interface{})) > 0 {
			for _, aclRelations := range v.([]interface{}) {
				aclRelationsMap := make(map[string]interface{})
				aclRelationsMap["acl_id"] = aclRelations.(map[string]interface{})["AclId"]
				aclRelationsMap["status"] = aclRelations.(map[string]interface{})["Status"]
				aclRelationsSli = append(aclRelationsSli, aclRelationsMap)
			}
		}
		aclConfigMap["acl_relations"] = aclRelationsSli
		aclConfigMap["acl_type"] = aclConfig.(map[string]interface{})["AclType"]
		aclConfigSli = append(aclConfigSli, aclConfigMap)
	}
	d.Set("acl_config", aclConfigSli)
	return nil
}

func resourceAliCloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		albServiceV2 := AlbServiceV2{client}
		object, err := albServiceV2.DescribeAlbListener(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["ListenerStatus"].(string) != target {
			if target == "Running" {
				action := "StartListener"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ListenerId"] = d.Id()

				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "IdempotenceProcessing", "IncorrectStatus.Listener", "VipStatusNotSupport", "-22001"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				action := "StopListener"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ListenerId"] = d.Id()

				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "IdempotenceProcessing", "IncorrectStatus.Listener", "-22001", "VipStatusNotSupport"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbListenerStateRefreshFunc(d.Id(), "ListenerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	action := "UpdateListenerAttribute"
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ListenerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("certificates") {
		update = true
		if v, ok := d.GetOk("certificates"); ok || d.HasChange("certificates") {
			certificatesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["CertificateId"] = dataLoopTmp["certificate_id"]
				certificatesMapsArray = append(certificatesMapsArray, dataLoopMap)
			}
			request["Certificates"] = certificatesMapsArray
		}
	}

	if !d.IsNewResource() && d.HasChange("request_timeout") {
		update = true
		request["RequestTimeout"] = d.Get("request_timeout")
	}

	if !d.IsNewResource() && d.HasChange("gzip_enabled") {
		update = true
		request["GzipEnabled"] = d.Get("gzip_enabled")
	}

	if !d.IsNewResource() && d.HasChange("http2_enabled") {
		update = true
		request["Http2Enabled"] = d.Get("http2_enabled")
	}

	if !d.IsNewResource() && d.HasChange("security_policy_id") {
		update = true
		request["SecurityPolicyId"] = d.Get("security_policy_id")
	}

	if !d.IsNewResource() && d.HasChange("quic_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("quic_config"); v != nil {
			quicUpgradeEnabled1, _ := jsonpath.Get("$[0].quic_upgrade_enabled", v)
			if quicUpgradeEnabled1 != nil && (d.HasChange("quic_config.0.quic_upgrade_enabled") || quicUpgradeEnabled1 != "") {
				objectDataLocalMap["QuicUpgradeEnabled"] = quicUpgradeEnabled1
			}
			quicListenerId1, _ := jsonpath.Get("$[0].quic_listener_id", v)
			if quicListenerId1 != nil && (d.HasChange("quic_config.0.quic_listener_id") || quicListenerId1 != "") {
				objectDataLocalMap["QuicListenerId"] = quicListenerId1
			}

			request["QuicConfig"] = objectDataLocalMap
		}
	}

	if !d.IsNewResource() && d.HasChange("default_actions") {
		update = true
	}
	if v, ok := d.GetOk("default_actions"); ok || d.HasChange("default_actions") {
		defaultActionsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			if !IsNil(dataLoop1Tmp["forward_group_config"]) {
				localData2 := make(map[string]interface{})
				if v, ok := dataLoop1Tmp["forward_group_config"]; ok {
					localData3, err := jsonpath.Get("$[0].server_group_tuples", v)
					if err != nil {
						localData3 = make([]interface{}, 0)
					}
					localMaps1 := make([]interface{}, 0)
					for _, dataLoop3 := range localData3.([]interface{}) {
						dataLoop3Tmp := make(map[string]interface{})
						if dataLoop3 != nil {
							dataLoop3Tmp = dataLoop3.(map[string]interface{})
						}
						dataLoop3Map := make(map[string]interface{})
						dataLoop3Map["ServerGroupId"] = dataLoop3Tmp["server_group_id"]
						localMaps1 = append(localMaps1, dataLoop3Map)
					}
					localData2["ServerGroupTuples"] = localMaps1
				}

				dataLoop1Map["ForwardGroupConfig"] = localData2
			}
			dataLoop1Map["Type"] = dataLoop1Tmp["type"]
			defaultActionsMapsArray = append(defaultActionsMapsArray, dataLoop1Map)
		}
		request["DefaultActions"] = defaultActionsMapsArray
	}

	if !d.IsNewResource() && d.HasChange("idle_timeout") {
		update = true
		request["IdleTimeout"] = d.Get("idle_timeout")
	}

	if !d.IsNewResource() && d.HasChange("ca_enabled") {
		update = true
		request["CaEnabled"] = d.Get("ca_enabled")
	}

	if !d.IsNewResource() && d.HasChange("ca_certificates") {
		update = true
		if v, ok := d.GetOk("ca_certificates"); ok || d.HasChange("ca_certificates") {
			caCertificatesMapsArray := make([]interface{}, 0)
			for _, dataLoop4 := range v.([]interface{}) {
				dataLoop4Tmp := dataLoop4.(map[string]interface{})
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["CertificateId"] = dataLoop4Tmp["certificate_id"]
				caCertificatesMapsArray = append(caCertificatesMapsArray, dataLoop4Map)
			}
			request["CaCertificates"] = caCertificatesMapsArray
		}
	}

	if !d.IsNewResource() && d.HasChange("listener_description") {
		update = true
		request["ListenerDescription"] = d.Get("listener_description")
	}

	if !d.IsNewResource() && d.HasChange("x_forwarded_for_config") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("x_forwarded_for_config"); v != nil {
			xForwardedForClientCertClientVerifyAlias1, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_client_verify_alias", v)
			if xForwardedForClientCertClientVerifyAlias1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_client_verify_alias") || xForwardedForClientCertClientVerifyAlias1 != "") {
				objectDataLocalMap1["XForwardedForClientCertClientVerifyAlias"] = xForwardedForClientCertClientVerifyAlias1
			}
			xForwardedForClientCertClientVerifyEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_client_verify_enabled", v)
			if xForwardedForClientCertClientVerifyEnabled1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_client_verify_enabled") || xForwardedForClientCertClientVerifyEnabled1 != "") {
				objectDataLocalMap1["XForwardedForClientCertClientVerifyEnabled"] = xForwardedForClientCertClientVerifyEnabled1
			}
			xForwardedForClientCertFingerPrintAlias, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_finger_print_alias", v)
			if xForwardedForClientCertFingerPrintAlias != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_finger_print_alias") || xForwardedForClientCertFingerPrintAlias != "") {
				objectDataLocalMap1["XForwardedForClientCertFingerprintAlias"] = xForwardedForClientCertFingerPrintAlias
			}
			xForwardedForClientCertFingerPrintEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_finger_print_enabled", v)
			if xForwardedForClientCertFingerPrintEnabled != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_finger_print_enabled") || xForwardedForClientCertFingerPrintEnabled != "") {
				objectDataLocalMap1["XForwardedForClientCertFingerprintEnabled"] = xForwardedForClientCertFingerPrintEnabled
			}
			xForwardedForClientCertIssuerDnAlias, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_issuer_dn_alias", v)
			if xForwardedForClientCertIssuerDnAlias != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_issuer_dn_alias") || xForwardedForClientCertIssuerDnAlias != "") {
				objectDataLocalMap1["XForwardedForClientCertIssuerDNAlias"] = xForwardedForClientCertIssuerDnAlias
			}
			xForwardedForClientCertIssuerDnEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_issuer_dn_enabled", v)
			if xForwardedForClientCertIssuerDnEnabled != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_issuer_dn_enabled") || xForwardedForClientCertIssuerDnEnabled != "") {
				objectDataLocalMap1["XForwardedForClientCertIssuerDNEnabled"] = xForwardedForClientCertIssuerDnEnabled
			}
			xForwardedForClientCertSubjectDnAlias, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_subject_dn_alias", v)
			if xForwardedForClientCertSubjectDnAlias != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_subject_dn_alias") || xForwardedForClientCertSubjectDnAlias != "") {
				objectDataLocalMap1["XForwardedForClientCertSubjectDNAlias"] = xForwardedForClientCertSubjectDnAlias
			}
			xForwardedForClientCertSubjectDnEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_client_cert_subject_dn_enabled", v)
			if xForwardedForClientCertSubjectDnEnabled != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_cert_subject_dn_enabled") || xForwardedForClientCertSubjectDnEnabled != "") {
				objectDataLocalMap1["XForwardedForClientCertSubjectDNEnabled"] = xForwardedForClientCertSubjectDnEnabled
			}
			xForwardedForClientSrcPortEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_client_src_port_enabled", v)
			if xForwardedForClientSrcPortEnabled1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_src_port_enabled") || xForwardedForClientSrcPortEnabled1 != "") {
				objectDataLocalMap1["XForwardedForClientSrcPortEnabled"] = xForwardedForClientSrcPortEnabled1
			}
			xForwardedForEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_enabled", v)
			if xForwardedForEnabled1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_enabled") || xForwardedForEnabled1 != "") {
				objectDataLocalMap1["XForwardedForEnabled"] = xForwardedForEnabled1
			}
			xForwardedForProtoEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_proto_enabled", v)
			if xForwardedForProtoEnabled1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_proto_enabled") || xForwardedForProtoEnabled1 != "") {
				objectDataLocalMap1["XForwardedForProtoEnabled"] = xForwardedForProtoEnabled1
			}
			xForwardedForSlbIdEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_slb_id_enabled", v)
			if xForwardedForSlbIdEnabled != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_slb_id_enabled") || xForwardedForSlbIdEnabled != "") {
				objectDataLocalMap1["XForwardedForSLBIdEnabled"] = xForwardedForSlbIdEnabled
			}
			xForwardedForSlbPortEnabled, _ := jsonpath.Get("$[0].x_forwarded_for_slb_port_enabled", v)
			if xForwardedForSlbPortEnabled != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_slb_port_enabled") || xForwardedForSlbPortEnabled != "") {
				objectDataLocalMap1["XForwardedForSLBPortEnabled"] = xForwardedForSlbPortEnabled
			}
			xForwardedForClientSourceIpsEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_client_source_ips_enabled", v)
			if xForwardedForClientSourceIpsEnabled1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_source_ips_enabled") || xForwardedForClientSourceIpsEnabled1 != "") {
				objectDataLocalMap1["XForwardedForClientSourceIpsEnabled"] = xForwardedForClientSourceIpsEnabled1
			}
			xForwardedForProcessingMode1, _ := jsonpath.Get("$[0].x_forwarded_for_processing_mode", v)
			if xForwardedForProcessingMode1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_processing_mode") || xForwardedForProcessingMode1 != "") {
				objectDataLocalMap1["XForwardedForProcessingMode"] = xForwardedForProcessingMode1
			}
			xForwardedForClientSourceIpsTrusted1, _ := jsonpath.Get("$[0].x_forwarded_for_client_source_ips_trusted", v)
			if xForwardedForClientSourceIpsTrusted1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_client_source_ips_trusted") || xForwardedForClientSourceIpsTrusted1 != "") {
				objectDataLocalMap1["XForwardedForClientSourceIpsTrusted"] = xForwardedForClientSourceIpsTrusted1
			}
			xForwardedForHostEnabled1, _ := jsonpath.Get("$[0].x_forwarded_for_host_enabled", v)
			if xForwardedForHostEnabled1 != nil && (d.HasChange("x_forwarded_for_config.0.x_forwarded_for_host_enabled") || xForwardedForHostEnabled1 != "") {
				objectDataLocalMap1["XForwardedForHostEnabled"] = xForwardedForHostEnabled1
			}

			request["XForwardedForConfig"] = objectDataLocalMap1
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if !d.IsNewResource() && d.HasChange("xforwarded_for_config") {
		update = true

		if v, ok := d.GetOk("xforwarded_for_config"); ok {
			xforwardedForConfigMap := map[string]interface{}{}
			for _, xforwardedForConfig := range v.(*schema.Set).List() {
				xforwardedForConfigArg := xforwardedForConfig.(map[string]interface{})
				xforwardedForConfigMap["XForwardedForClientCertIssuerDNAlias"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnalias"]
				xforwardedForConfigMap["XForwardedForClientCertIssuerDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcert_issuerdnenabled"]
				xforwardedForConfigMap["XForwardedForClientCertClientVerifyAlias"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyalias"]
				xforwardedForConfigMap["XForwardedForClientCertClientVerifyEnabled"] = xforwardedForConfigArg["xforwardedforclientcertclientverifyenabled"]
				xforwardedForConfigMap["XForwardedForClientCertFingerprintAlias"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintalias"]
				xforwardedForConfigMap["XForwardedForClientCertFingerprintEnabled"] = xforwardedForConfigArg["xforwardedforclientcertfingerprintenabled"]
				xforwardedForConfigMap["XForwardedForClientCertSubjectDNAlias"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnalias"]
				xforwardedForConfigMap["XForwardedForClientCertSubjectDNEnabled"] = xforwardedForConfigArg["xforwardedforclientcertsubjectdnenabled"]
				xforwardedForConfigMap["XForwardedForClientSrcPortEnabled"] = xforwardedForConfigArg["xforwardedforclientsrcportenabled"]
				xforwardedForConfigMap["XForwardedForEnabled"] = xforwardedForConfigArg["xforwardedforenabled"]
				xforwardedForConfigMap["XForwardedForProtoEnabled"] = xforwardedForConfigArg["xforwardedforprotoenabled"]
				xforwardedForConfigMap["XForwardedForSLBIdEnabled"] = xforwardedForConfigArg["xforwardedforslbidenabled"]
				xforwardedForConfigMap["XForwardedForSLBPortEnabled"] = xforwardedForConfigArg["xforwardedforslbportenabled"]
			}

			request["XForwardedForConfig"] = xforwardedForConfigMap
		}
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "IdempotenceProcessing", "IncorrectStatus.Listener", "VipStatusNotSupport", "-22001"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.DescribeAsyncAlbListenerStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "UpdateListenerLogConfig"
	conn, err = client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ListenerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("access_log_record_customized_headers_enabled") {
		update = true
		request["AccessLogRecordCustomizedHeadersEnabled"] = d.Get("access_log_record_customized_headers_enabled")
	}

	if d.HasChange("access_log_tracing_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("access_log_tracing_config"); v != nil {
			tracingSample1, _ := jsonpath.Get("$[0].tracing_sample", v)
			if tracingSample1 != nil && (d.HasChange("access_log_tracing_config.0.tracing_sample") || tracingSample1 != "") && tracingSample1.(int) > 0 {
				objectDataLocalMap["TracingSample"] = tracingSample1
			}
			tracingType1, _ := jsonpath.Get("$[0].tracing_type", v)
			if tracingType1 != nil && (d.HasChange("access_log_tracing_config.0.tracing_type") || tracingType1 != "") {
				objectDataLocalMap["TracingType"] = tracingType1
			}
			tracingEnabled1, _ := jsonpath.Get("$[0].tracing_enabled", v)
			if tracingEnabled1 != nil && (d.HasChange("access_log_tracing_config.0.tracing_enabled") || tracingEnabled1 != "") {
				objectDataLocalMap["TracingEnabled"] = tracingEnabled1
			}

			request["AccessLogTracingConfig"] = objectDataLocalMap
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "IdempotenceProcessing", "IncorrectStatus.Listener", "-22001", "VipStatusNotSupport"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutUpdate), 0, albServiceV2.DescribeAsyncAlbListenerStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	if d.HasChange("tags") {
		albServiceV2 := AlbServiceV2{client}
		if err := albServiceV2.SetResourceTags(d, "listener"); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("acl_config") {
		albService := AlbService{client}
		albServiceV2 := AlbServiceV2{client}
		oldAssociateAcls, newAssociateVpcs := d.GetChange("acl_config")
		oldAssociateAclsSet := oldAssociateAcls.(*schema.Set)
		newAssociateAclsSet := newAssociateVpcs.(*schema.Set)
		removed := oldAssociateAclsSet.Difference(newAssociateAclsSet)
		added := newAssociateAclsSet.Difference(oldAssociateAclsSet)

		if removed.Len() > 0 {
			action := "DissociateAclsFromListener"
			dissociateAclsFromListenerReq := map[string]interface{}{
				"ListenerId": d.Id(),
			}
			dissociateAclsFromListenerReq["ClientToken"] = buildClientToken("DissociateAclsFromListener")
			associateAclIds := make([]string, 0)
			for _, aclConfig := range removed.List() {
				if aclRelationsMaps, ok := aclConfig.(map[string]interface{})["acl_relations"]; ok {
					for _, aclRelationsMap := range aclRelationsMaps.(*schema.Set).List() {
						associateAclIds = append(associateAclIds, aclRelationsMap.(map[string]interface{})["acl_id"].(string))
					}
				}
			}
			dissociateAclsFromListenerReq["AclIds"] = associateAclIds
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, dissociateAclsFromListenerReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"VipStatusNotSupport"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, dissociateAclsFromListenerReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

			stateConf = BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbJobStateRefreshFunc(d.Id(), "listener", response["JobId"].(string), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if added.Len() > 0 {
			action := "AssociateAclsWithListener"
			associateAclsWithListenerReq := map[string]interface{}{
				"ListenerId": d.Id(),
			}
			associateAclsWithListenerReq["ClientToken"] = buildClientToken("AssociateAclsWithListener")
			associateAclIds := make([]string, 0)
			for _, aclConfig := range added.List() {
				if aclRelationsMaps, ok := aclConfig.(map[string]interface{})["acl_relations"]; ok {
					for _, aclRelationsMap := range aclRelationsMaps.(*schema.Set).List() {
						associateAclIds = append(associateAclIds, aclRelationsMap.(map[string]interface{})["acl_id"].(string))
					}
				}
				associateAclsWithListenerReq["AclType"] = aclConfig.(map[string]interface{})["acl_type"]
			}
			associateAclsWithListenerReq["AclIds"] = associateAclIds
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, associateAclsWithListenerReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"VipStatusNotSupport"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, associateAclsWithListenerReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbListenerStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

			stateConf = BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbJobStateRefreshFunc(d.Id(), "listener", response["JobId"].(string), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	d.Partial(false)
	return resourceAliCloudAlbListenerRead(d, meta)
}

func resourceAliCloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteListener"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ListenerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectBusinessStatus.LoadBalancer", "SystemBusy", "IdempotenceProcessing", "ResourceInConfiguring.Listener", "IncorrectStatus.LoadBalancer", "-22031"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.DescribeAsyncAlbListenerStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}

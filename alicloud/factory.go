package alicloud

import (
	"context"
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProtoV5ProviderServerFactory(ctx context.Context) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	primary, err := New(ctx)

	if err != nil {
		return nil, nil, err
	}

	servers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider(primary)),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		return nil, nil, err
	}

	return muxServer.ProviderServer, primary, nil
}

func New(ctx context.Context) (*schema.Provider, error) {
	provider := &schema.Provider{
		// This schema must match exactly the Terraform Protocol v6 (Terraform Plugin Framework) provider's schema.
		// Notably the attributes can have no Default values.
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCESS_KEY", os.Getenv("ALIBABACLOUD_ACCESS_KEY_ID")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECRET_KEY", os.Getenv("ALIBABACLOUD_ACCESS_KEY_SECRET")),
				Description: descriptions["secret_key"],
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECURITY_TOKEN", os.Getenv("ALIBABACLOUD_SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ECS_ROLE_NAME", os.Getenv("ALICLOUD_ECS_ROLE_NAME")),
				Description: descriptions["ecs_role_name"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_REGION", os.Getenv("ALICLOUD_REGION")),
				Description: descriptions["region"],
			},
			"ots_instance_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'ots_instance_name' has been deprecated from provider version 1.10.0. New field 'instance_name' of resource 'alicloud_ots_table' instead.",
			},
			"log_endpoint": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'log_endpoint' has been deprecated from provider version 1.28.0. New field 'log' which in nested endpoints instead.",
			},
			"mns_endpoint": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'mns_endpoint' has been deprecated from provider version 1.28.0. New field 'mns' which in nested endpoints instead.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCOUNT_ID", os.Getenv("ALICLOUD_ACCOUNT_ID")),
				Description: descriptions["account_id"],
			},
			"assume_role":           assumeRoleSchema(),
			"sign_version":          signVersionSchema(),
			"assume_role_with_oidc": assumeRoleWithOidcSchema(),
			"fc": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'fc' has been deprecated from provider version 1.28.0. New field 'fc' which in nested endpoints instead.",
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SHARED_CREDENTIALS_FILE", ""),
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_PROFILE", ""),
			},
			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},
			"configuration_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 128),
				DefaultFunc:  schema.EnvDefaultFunc("TF_APPEND_USER_AGENT", ""),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTPS",
				Description:  descriptions["protocol"],
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"client_read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_READ_TIMEOUT", 60000),
				Description: descriptions["client_read_timeout"],
			},
			"client_connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_CONNECT_TIMEOUT", 60000),
				Description: descriptions["client_connect_timeout"],
			},
			"source_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SOURCE_IP", os.Getenv("ALICLOUD_SOURCE_IP")),
				Description: descriptions["source_ip"],
			},
			"security_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECURITY_TRANSPORT", os.Getenv("ALICLOUD_SECURITY_TRANSPORT")),
				//Deprecated:  "It has been deprecated from version 1.136.0 and using new field secure_transport instead.",
			},
			"secure_transport": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECURE_TRANSPORT", os.Getenv("ALICLOUD_SECURE_TRANSPORT")),
				Description: descriptions["secure_transport"],
			},
			"credentials_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_CREDENTIALS_URI", os.Getenv("ALICLOUD_CREDENTIALS_URI")),
				Description: descriptions["credentials_uri"],
			},
			"max_retry_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MAX_RETRY_TIMEOUT", 0),
				Description: descriptions["max_retry_timeout"],
			},
		},

		// Data sources and resources implemented using Terraform Plugin SDK
		// should use the @SDKDataSource and @SDKResource function-level annotations
		// rather than adding directly to these maps.
		DataSourcesMap: make(map[string]*schema.Resource),
		ResourcesMap:   make(map[string]*schema.Resource),
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return configure(ctx, d, provider)
	}

	// Set the provider Meta (instance data) here.
	// It will be overwritten by the result of the call to ConfigureContextFunc,
	// but can be used pre-configuration by other (non-primary) provider servers.
	// NOTICE: 这里需要修改
	// provider.SetMeta(meta)

	return provider, nil
}

func configure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	log.Println("using terraform version:", p.TerraformVersion)
	var diags diag.Diagnostics

	var getProviderConfig = func(str string, key string) string {
		if str == "" {
			value, err := getConfigFromProfile(d, key)
			if err == nil && value != nil {
				str = value.(string)
			}
		}
		return str
	}

	accessKey := getProviderConfig(d.Get("access_key").(string), "access_key_id")
	secretKey := getProviderConfig(d.Get("secret_key").(string), "access_key_secret")
	region := getProviderConfig(d.Get("region").(string), "region_id")
	if region == "" {
		region = DEFAULT_REGION
	}
	securityToken := getProviderConfig(d.Get("security_token").(string), "sts_token")

	ecsRoleName := getProviderConfig(d.Get("ecs_role_name").(string), "ram_role_name")

	if accessKey == "" || secretKey == "" {
		if v, ok := d.GetOk("credentials_uri"); ok && v.(string) != "" {
			credentialsURIResp, err := getClientByCredentialsURI(v.(string))
			if err != nil {
				return nil, appendFromErr(diags, err)
			}
			accessKey = credentialsURIResp.AccessKeyId
			secretKey = credentialsURIResp.AccessKeySecret
			securityToken = credentialsURIResp.SecurityToken
		}
	}

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
		ConfigurationSource:  d.Get("configuration_source").(string),
		Protocol:             d.Get("protocol").(string),
		ClientReadTimeout:    d.Get("client_read_timeout").(int),
		ClientConnectTimeout: d.Get("client_connect_timeout").(int),
		SourceIp:             strings.TrimSpace(d.Get("source_ip").(string)),
		SecureTransport:      strings.TrimSpace(d.Get("secure_transport").(string)),
		MaxRetryTimeout:      d.Get("max_retry_timeout").(int),
		TerraformTraceId:     strings.Trim(uuid.New().String(), "-"),
		TerraformVersion:     p.TerraformVersion,
	}
	log.Println("alicloud provider trace id:", config.TerraformTraceId)
	if v, ok := d.GetOk("security_transport"); config.SecureTransport == "" && ok && v.(string) != "" {
		config.SecureTransport = v.(string)
	}
	config.SecurityToken = strings.TrimSpace(securityToken)

	config.RamRoleArn = getProviderConfig("", "ram_role_arn")
	config.RamRoleSessionName = getProviderConfig("", "ram_session_name")
	expiredSeconds, err := getConfigFromProfile(d, "expired_seconds")
	if err == nil && expiredSeconds != nil {
		config.RamRoleSessionExpiration = (int)(expiredSeconds.(float64))
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		if assumeRole["role_arn"].(string) != "" {
			config.RamRoleArn = assumeRole["role_arn"].(string)
		}
		if assumeRole["session_name"].(string) != "" {
			config.RamRoleSessionName = assumeRole["session_name"].(string)
		}
		if config.RamRoleSessionName == "" {
			config.RamRoleSessionName = "terraform"
		}
		config.RamRolePolicy = assumeRole["policy"].(string)
		if assumeRole["session_expiration"].(int) == 0 {
			if v := os.Getenv("ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION"); v != "" {
				if expiredSeconds, err := strconv.Atoi(v); err == nil {
					config.RamRoleSessionExpiration = expiredSeconds
				}
			}
			if config.RamRoleSessionExpiration == 0 {
				config.RamRoleSessionExpiration = 3600
			}
		} else {
			config.RamRoleSessionExpiration = assumeRole["session_expiration"].(int)
		}
		if v := assumeRole["external_id"].(string); v != "" {
			config.RamRoleExternalId = v
		}

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d, RamRoleExternalId: %s)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration, config.RamRoleExternalId)
	}

	if v, ok := d.GetOk("assume_role_with_oidc"); ok && len(v.([]interface{})) == 1 {
		config.AssumeRoleWithOidc, err = getAssumeRoleWithOIDCConfig(v.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return nil, appendFromErr(diags, err)
		}
		log.Printf("[INFO] assume_role_with_oidc configuration set: (RoleArn: %q, SessionName: %q, SessionExpiration: %d, OIDCProviderArn: %s)",
			config.AssumeRoleWithOidc.RoleARN, config.AssumeRoleWithOidc.RoleSessionName, config.AssumeRoleWithOidc.DurationSeconds, config.AssumeRoleWithOidc.OIDCProviderArn)
	}

	if err := config.RefreshAuthCredential(); err != nil {
		return nil, appendFromErr(diags, err)
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)
	var endpointInit sync.Map
	config.Endpoints = &endpointInit

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		for key, val := range endpoints {
			endpointInit.Store(key, val)
		}
		config.EcsEndpoint = strings.TrimSpace(endpoints["ecs"].(string))
		config.RdsEndpoint = strings.TrimSpace(endpoints["rds"].(string))
		config.SlbEndpoint = strings.TrimSpace(endpoints["slb"].(string))
		config.VpcEndpoint = strings.TrimSpace(endpoints["vpc"].(string))
		config.EssEndpoint = strings.TrimSpace(endpoints["ess"].(string))
		config.OssEndpoint = strings.TrimSpace(endpoints["oss"].(string))
		config.OnsEndpoint = strings.TrimSpace(endpoints["ons"].(string))
		config.AlikafkaEndpoint = strings.TrimSpace(endpoints["alikafka"].(string))
		config.DnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		config.RamEndpoint = strings.TrimSpace(endpoints["ram"].(string))
		config.CsEndpoint = strings.TrimSpace(endpoints["cs"].(string))
		config.CrEndpoint = strings.TrimSpace(endpoints["cr"].(string))
		config.CdnEndpoint = strings.TrimSpace(endpoints["cdn"].(string))
		config.KmsEndpoint = strings.TrimSpace(endpoints["kms"].(string))
		config.OtsEndpoint = strings.TrimSpace(endpoints["ots"].(string))
		config.CmsEndpoint = strings.TrimSpace(endpoints["cms"].(string))
		config.PvtzEndpoint = strings.TrimSpace(endpoints["pvtz"].(string))
		config.StsEndpoint = strings.TrimSpace(endpoints["sts"].(string))
		config.LogEndpoint = strings.TrimSpace(endpoints["log"].(string))
		config.DrdsEndpoint = strings.TrimSpace(endpoints["drds"].(string))
		config.DdsEndpoint = strings.TrimSpace(endpoints["dds"].(string))
		config.GpdbEnpoint = strings.TrimSpace(endpoints["gpdb"].(string))
		config.KVStoreEndpoint = strings.TrimSpace(endpoints["kvstore"].(string))
		config.PolarDBEndpoint = strings.TrimSpace(endpoints["polardb"].(string))
		config.FcEndpoint = strings.TrimSpace(endpoints["fc"].(string))
		config.ApigatewayEndpoint = strings.TrimSpace(endpoints["apigateway"].(string))
		config.DatahubEndpoint = strings.TrimSpace(endpoints["datahub"].(string))
		config.MnsEndpoint = strings.TrimSpace(endpoints["mns"].(string))
		config.LocationEndpoint = strings.TrimSpace(endpoints["location"].(string))
		config.ElasticsearchEndpoint = strings.TrimSpace(endpoints["elasticsearch"].(string))
		config.NasEndpoint = strings.TrimSpace(endpoints["nas"].(string))
		config.ActiontrailEndpoint = strings.TrimSpace(endpoints["actiontrail"].(string))
		config.BssOpenApiEndpoint = strings.TrimSpace(endpoints["bssopenapi"].(string))
		config.DdoscooEndpoint = strings.TrimSpace(endpoints["ddoscoo"].(string))
		config.DdosbgpEndpoint = strings.TrimSpace(endpoints["ddosbgp"].(string))
		config.EmrEndpoint = strings.TrimSpace(endpoints["emr"].(string))
		config.CasEndpoint = strings.TrimSpace(endpoints["cas"].(string))
		config.MarketEndpoint = strings.TrimSpace(endpoints["market"].(string))
		config.AdbEndpoint = strings.TrimSpace(endpoints["adb"].(string))
		config.CbnEndpoint = strings.TrimSpace(endpoints["cbn"].(string))
		config.MaxComputeEndpoint = strings.TrimSpace(endpoints["maxcompute"].(string))
		config.DmsEnterpriseEndpoint = strings.TrimSpace(endpoints["dms_enterprise"].(string))
		config.WafOpenapiEndpoint = strings.TrimSpace(endpoints["waf_openapi"].(string))
		config.ResourcemanagerEndpoint = strings.TrimSpace(endpoints["resourcemanager"].(string))
		config.EciEndpoint = strings.TrimSpace(endpoints["eci"].(string))
		config.OosEndpoint = strings.TrimSpace(endpoints["oos"].(string))
		config.DcdnEndpoint = strings.TrimSpace(endpoints["dcdn"].(string))
		config.MseEndpoint = strings.TrimSpace(endpoints["mse"].(string))
		config.ConfigEndpoint = strings.TrimSpace(endpoints["config"].(string))
		config.RKvstoreEndpoint = strings.TrimSpace(endpoints["r_kvstore"].(string))
		config.FnfEndpoint = strings.TrimSpace(endpoints["fnf"].(string))
		config.RosEndpoint = strings.TrimSpace(endpoints["ros"].(string))
		config.PrivatelinkEndpoint = strings.TrimSpace(endpoints["privatelink"].(string))
		config.ResourcesharingEndpoint = strings.TrimSpace(endpoints["resourcesharing"].(string))
		config.GaEndpoint = strings.TrimSpace(endpoints["ga"].(string))
		config.HitsdbEndpoint = strings.TrimSpace(endpoints["hitsdb"].(string))
		config.BrainIndustrialEndpoint = strings.TrimSpace(endpoints["brain_industrial"].(string))
		config.EipanycastEndpoint = strings.TrimSpace(endpoints["eipanycast"].(string))
		config.ImsEndpoint = strings.TrimSpace(endpoints["ims"].(string))
		config.QuotasEndpoint = strings.TrimSpace(endpoints["quotas"].(string))
		config.SgwEndpoint = strings.TrimSpace(endpoints["sgw"].(string))
		config.DmEndpoint = strings.TrimSpace(endpoints["dm"].(string))
		config.EventbridgeEndpoint = strings.TrimSpace(endpoints["eventbridge"].(string))
		config.OnsproxyEndpoint = strings.TrimSpace(endpoints["onsproxy"].(string))
		config.CdsEndpoint = strings.TrimSpace(endpoints["cds"].(string))
		config.HbrEndpoint = strings.TrimSpace(endpoints["hbr"].(string))
		config.ArmsEndpoint = strings.TrimSpace(endpoints["arms"].(string))
		config.ServerlessEndpoint = strings.TrimSpace(endpoints["serverless"].(string))
		config.AlbEndpoint = strings.TrimSpace(endpoints["alb"].(string))
		config.RedisaEndpoint = strings.TrimSpace(endpoints["redisa"].(string))
		config.GwsecdEndpoint = strings.TrimSpace(endpoints["gwsecd"].(string))
		config.CloudphoneEndpoint = strings.TrimSpace(endpoints["cloudphone"].(string))
		config.ScdnEndpoint = strings.TrimSpace(endpoints["scdn"].(string))
		config.DataworkspublicEndpoint = strings.TrimSpace(endpoints["dataworkspublic"].(string))
		config.HcsSgwEndpoint = strings.TrimSpace(endpoints["hcs_sgw"].(string))
		config.CddcEndpoint = strings.TrimSpace(endpoints["cddc"].(string))
		config.MscopensubscriptionEndpoint = strings.TrimSpace(endpoints["mscopensubscription"].(string))
		config.SddpEndpoint = strings.TrimSpace(endpoints["sddp"].(string))
		config.BastionhostEndpoint = strings.TrimSpace(endpoints["bastionhost"].(string))
		config.SasEndpoint = strings.TrimSpace(endpoints["sas"].(string))
		config.AlidfsEndpoint = strings.TrimSpace(endpoints["alidfs"].(string))
		config.EhpcEndpoint = strings.TrimSpace(endpoints["ehpc"].(string))
		config.EnsEndpoint = strings.TrimSpace(endpoints["ens"].(string))
		config.IotEndpoint = strings.TrimSpace(endpoints["iot"].(string))
		config.ImmEndpoint = strings.TrimSpace(endpoints["imm"].(string))
		config.ClickhouseEndpoint = strings.TrimSpace(endpoints["clickhouse"].(string))
		config.DtsEndpoint = strings.TrimSpace(endpoints["dts"].(string))
		config.DgEndpoint = strings.TrimSpace(endpoints["dg"].(string))
		config.CloudssoEndpoint = strings.TrimSpace(endpoints["cloudsso"].(string))
		config.WafEndpoint = strings.TrimSpace(endpoints["waf"].(string))
		config.SwasEndpoint = strings.TrimSpace(endpoints["swas"].(string))
		config.VsEndpoint = strings.TrimSpace(endpoints["vs"].(string))
		config.QuickbiEndpoint = strings.TrimSpace(endpoints["quickbi"].(string))
		config.VodEndpoint = strings.TrimSpace(endpoints["vod"].(string))
		config.OpensearchEndpoint = strings.TrimSpace(endpoints["opensearch"].(string))
		config.GdsEndpoint = strings.TrimSpace(endpoints["gds"].(string))
		config.DbfsEndpoint = strings.TrimSpace(endpoints["dbfs"].(string))
		config.DevopsrdcEndpoint = strings.TrimSpace(endpoints["devopsrdc"].(string))
		config.EaisEndpoint = strings.TrimSpace(endpoints["eais"].(string))
		config.CloudauthEndpoint = strings.TrimSpace(endpoints["cloudauth"].(string))
		config.ImpEndpoint = strings.TrimSpace(endpoints["imp"].(string))
		config.MhubEndpoint = strings.TrimSpace(endpoints["mhub"].(string))
		config.ServicemeshEndpoint = strings.TrimSpace(endpoints["servicemesh"].(string))
		config.AcrEndpoint = strings.TrimSpace(endpoints["acr"].(string))
		config.EdsuserEndpoint = strings.TrimSpace(endpoints["edsuser"].(string))
		config.GaplusEndpoint = strings.TrimSpace(endpoints["gaplus"].(string))
		config.DdosbasicEndpoint = strings.TrimSpace(endpoints["ddosbasic"].(string))
		config.SmartagEndpoint = strings.TrimSpace(endpoints["smartag"].(string))
		config.TagEndpoint = strings.TrimSpace(endpoints["tag"].(string))
		config.EdasEndpoint = strings.TrimSpace(endpoints["edas"].(string))
		config.EdasschedulerxEndpoint = strings.TrimSpace(endpoints["edasschedulerx"].(string))
		config.EhsEndpoint = strings.TrimSpace(endpoints["ehs"].(string))
		config.CloudfwEndpoint = strings.TrimSpace(endpoints["cloudfw"].(string))
		config.DysmsEndpoint = strings.TrimSpace(endpoints["dysms"].(string))
		config.CbsEndpoint = strings.TrimSpace(endpoints["cbs"].(string))
		config.NlbEndpoint = strings.TrimSpace(endpoints["nlb"].(string))
		config.VpcpeerEndpoint = strings.TrimSpace(endpoints["vpcpeer"].(string))
		config.EbsEndpoint = strings.TrimSpace(endpoints["ebs"].(string))
		config.DmsenterpriseEndpoint = strings.TrimSpace(endpoints["dmsenterprise"].(string))
		config.BpStudioEndpoint = strings.TrimSpace(endpoints["bpstudio"].(string))
		config.DasEndpoint = strings.TrimSpace(endpoints["das"].(string))
		config.CloudfirewallEndpoint = strings.TrimSpace(endpoints["cloudfirewall"].(string))
		config.SrvcatalogEndpoint = strings.TrimSpace(endpoints["srvcatalog"].(string))
		config.VpcPeerEndpoint = strings.TrimSpace(endpoints["vpcpeer"].(string))
		config.EfloEndpoint = strings.TrimSpace(endpoints["eflo"].(string))
		config.OceanbaseEndpoint = strings.TrimSpace(endpoints["oceanbase"].(string))
		config.BeebotEndpoint = strings.TrimSpace(endpoints["beebot"].(string))
		config.ComputeNestEndpoint = strings.TrimSpace(endpoints["computenest"].(string))
		if endpoint, ok := endpoints["alidns"]; ok {
			config.AlidnsEndpoint = strings.TrimSpace(endpoint.(string))
		} else {
			config.AlidnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		}
		config.CassandraEndpoint = strings.TrimSpace(endpoints["cassandra"].(string))
	}

	var signVersion sync.Map
	config.SignVersion = &signVersion
	for _, version := range d.Get("sign_version").(*schema.Set).List() {
		for key, val := range version.(map[string]interface{}) {
			signVersion.Store(key, val)
		}
	}

	if config.RamRoleArn != "" {
		config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config)
		if err != nil {
			return nil, appendFromErr(diags, err)
		}
	}
	if (config.AccessKey == "" || config.SecretKey == "") && config.EcsRoleName == "" {
		return nil, appendFromErr(diags, fmt.Errorf("configuring Terraform Alibaba Cloud Provider: no valid credential sources for Terraform Alibaba Cloud Provider found.\n\n%s",
			"Please see https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#authentication\n"+
				"for more information about providing credentials."))
	}

	if ots_instance_name, ok := d.GetOk("ots_instance_name"); ok && ots_instance_name.(string) != "" {
		config.OtsInstanceName = strings.TrimSpace(ots_instance_name.(string))
	}

	if logEndpoint, ok := d.GetOk("log_endpoint"); ok && logEndpoint.(string) != "" {
		config.LogEndpoint = strings.TrimSpace(logEndpoint.(string))
	}
	if mnsEndpoint, ok := d.GetOk("mns_endpoint"); ok && mnsEndpoint.(string) != "" {
		config.MnsEndpoint = strings.TrimSpace(mnsEndpoint.(string))
	}

	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		config.AccountId = strings.TrimSpace(account.(string))
	}

	if fcEndpoint, ok := d.GetOk("fc"); ok && fcEndpoint.(string) != "" {
		config.FcEndpoint = strings.TrimSpace(fcEndpoint.(string))
	}

	configurationSources := []string{
		fmt.Sprintf("Default/%s", config.TerraformTraceId),
	}

	// configuration source final value should also contain TF_APPEND_USER_AGENT value
	// there is need to deduplication
	config.ConfigurationSource += " " + strings.TrimSpace(os.Getenv("TF_APPEND_USER_AGENT"))
	if config.ConfigurationSource != "" {
		for _, s := range strings.Split(config.ConfigurationSource, " ") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			exist := false
			for _, con := range configurationSources {
				if s == con {
					exist = true
					break
				}
			}
			if !exist {
				configurationSources = append(configurationSources, s)
			}
		}
	}
	config.ConfigurationSource = strings.Join(configurationSources, " ") + getModuleAddr()

	client, err := config.Client()
	if err != nil {
		return nil, appendFromErr(diags, err)
	}

	return client, nil
}

func appendFromErr(diags diag.Diagnostics, err error) diag.Diagnostics {
	if err == nil {
		return diags
	}
	return append(diags, diag.FromErr(err)...) // nosemgrep:ci.semgrep.pluginsdk.avoid-diag_FromErr
}

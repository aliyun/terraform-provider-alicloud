package config

// productCodeToLocationCode records all products' code mapping to location
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: location code
var productCodeToLocationCode = map[string]string{
	"ecs":                     "ecs",     // ECS
	"adb":                     "ads",     //ADB
	"ess":                     "ess",     //AutoScaling
	"cs":                      "cs",      // ACK
	"polardb":                 "polardb", // PolarDB
	"cr":                      "acr",     // CR
	"dds":                     "dds",     //MongoDB
	"gpdb":                    "gpdb",    //GPDB
	"fc_open":                 "fc",      // FC, FCV2
	"fc":                      "fc",      // FCV3
	"cloudapi":                "apigateway",
	"apig":                    "nativeapigw",       // APIG
	"datahub":                 "datahub",           // DataHub
	"mns_open":                "mns",               // MessageService
	"elasticsearch":           "elasticsearch",     // Elasticsearch
	"ddoscoo":                 "ddoscoo",           // DdosCoo
	"ddosbgp":                 "ddosbgp",           // DdosBgp
	"antiddos_public":         "ddosbasic",         // DdosBasic
	"bssopenapi":              "bssopenapi",        //BssOpenApi
	"alikafka":                "alikafka",          //AliKafka
	"emr":                     "emr",               //EMR
	"smartag":                 "smartag",           // Smartag
	"yundun_dbaudit":          "dbaudit",           //DBAudit
	"yundun_bastionhost":      "bastionhost",       //Bastionhost
	"hbase":                   "hbase",             //HBase
	"edas":                    "edas",              // EDAS
	"alidns":                  "alidns",            //Alidns
	"cassandra":               "cds",               //Cassandra
	"eci":                     "eci",               // ECI
	"dcdn":                    "dcdn",              // DCDN
	"r_kvstore":               "redisa",            // Redis
	"ons":                     "ons",               //Ons
	"config":                  "config",            //Config
	"fnf":                     "fnf",               // FnF
	"ros":                     "ros",               // ROS
	"mse":                     "mse",               // MSE
	"pvtz":                    "pvtz",              //PrivateZone
	"privatelink":             "privatelink",       // PrivateLink
	"maxcompute":              "odps",              //MaxCompute
	"resourcesharing":         "ressharing",        // ResourceManager
	"ga":                      "gaplus",            // Ga
	"actiontrail":             "actiontrail",       //ActionTrail
	"hitsdb":                  "hitsdb",            //Lindorm
	"brain_industrial":        "aistudio",          //BrainIndustrial
	"eipanycast":              "eipanycast",        // Eipanycast
	"oos":                     "oos",               // OOS
	"ims":                     "ims",               //IMS
	"resourcemanager":         "resourcemanager",   // ResourceManager
	"nas":                     "nas",               //NAS
	"dms_enterprise":          "dmsenterprise",     //DMSEnterprise
	"sgw":                     "hcs_sgw",           // CloudStorageGateway
	"slb":                     "slb",               // SLB
	"kms":                     "kms",               //KMS
	"dm":                      "dm",                //DirectMail
	"eventbridge":             "eventbridge",       // EventBridge
	"hbr":                     "hbr",               //HBR
	"cas":                     "cas",               //SSLCertificatesService
	"arms":                    "arms",              // ARMS
	"cloudfw":                 "cloudfirewall",     //CloudFirewall
	"sae":                     "serverless",        //SAE
	"alb":                     "alb",               // ALB
	"ecd":                     "gwsecd",            // ECD
	"cloudphone":              "cloudphone",        // ECP
	"scdn":                    "scdn",              //SCDN
	"dataworks_public":        "dide",              //DataWorks
	"cdn":                     "cdn",               // CDN
	"cddc":                    "cddc",              // CDDC
	"mscopensubscription":     "mscsub",            //MscSub
	"sddp":                    "sddp",              // SDDP
	"sas":                     "sas",               // ThreatDetection
	"ehpc":                    "ehs",               // Ehpc
	"ens":                     "ens",               // ENS
	"iot":                     "iot",               // Iot
	"imm":                     "imm",               // IMM
	"clickhouse":              "clickhouse",        // ClickHouse
	"selectdb":                "selectdb",          //SelectDB
	"dts":                     "dts",               // DTS
	"dg":                      "dg",                // DatabaseGateway
	"cloudsso":                "cloudsso",          // CloudSSO
	"swas_open":               "swas",              // SimpleApplicationServer
	"vs":                      "vs",                // VideoSurveillanceSystem
	"quickbi_public":          "quickbi",           // QuickBI
	"devops_rdc":              "rdcdevops",         // RDC
	"vod":                     "vod",               // VOD
	"opensearch":              "opensearch",        // OpenSearch
	"gdb":                     "gds",               // GraphDatabase
	"dbfs":                    "dbfs",              // DBFS
	"eais":                    "eais",              // EAIS
	"cloudauth":               "cloudauth",         // Cloudauth
	"imp":                     "imp",               // IMP
	"mhub":                    "emas",              // MHUB
	"servicemesh":             "servicemesh",       // ServiceMesh
	"starrocks":               "emrstarrocks",      // StarRocks
	"eds_user":                "edsuser",           // ECD
	"tag":                     "tag",               // Tag
	"schedulerx2":             "edasschedulerx",    // Schedulerx
	"dysmsapi":                "dysms",             // SMS
	"vpcpeer":                 "vpcpeer",           // VpcPeer
	"dbs":                     "cbs",               // DBS
	"nlb":                     "nlb",               // NLB
	"ebs":                     "ebs",               // EBS
	"bpstudio":                "bpstudio",          // BPStudio
	"das":                     "hdm",               // DAS
	"servicecatalog":          "srvcatalog",        // ServiceCatalog
	"eflo":                    "eflo",              //Eflo
	"oceanbasepro":            "oceanbase",         // OceanBase
	"chatbot":                 "beebot",            // Chatbot
	"computenest":             "computenest",       // ComputeNest
	"drds":                    "drds",              // DRDS
	"polardbx":                "polardbx",          // DRDS
	"adcp":                    "adcp",              // AckOne
	"sls":                     "sls",               // SLS
	"rocketmq":                "rmq",               // RocketMQ
	"resourcecenter":          "",                  // ResourceManager
	"hologram":                "hologram",          // Hologram
	"foasconsole":             "foasconsole",       // RealtimeCompute
	"vpc":                     "vpc",               // VPC, VPNGateway,ExpressConnect, CBWP, EIP
	"oss":                     "oss",               // OSS
	"cms":                     "cms",               // CloudMonitorService
	"waf_openapi":             "waf",               //WAFV3,WAF
	"dfs":                     "alidfs",            //DFS
	"amqp":                    "onsproxy",          // Amqp
	"amqp_open":               "onsproxy",          // Amqp
	"cbn":                     "cbn",               // CEN
	"expressconnectrouter":    "ecr",               // ExpressConnectRouter
	"green":                   "green",             // Aligreen
	"governance":              "governance",        // Governance
	"ots":                     "ots",               // OTS
	"tablestore":              "ots",               // OTS
	"ram":                     "ram",               //RAM
	"quotas":                  "quotas",            //Quotas
	"market":                  "market",            //Market
	"aiworkspace":             "paiworkspace",      //PAIWorkspace
	"vpcipam":                 "vpcipam",           //VpcIpam
	"gwlb":                    "gwlb",              // GWLB
	"esa":                     "dcdnservices",      // ESA
	"live":                    "live",              // Live
	"eds_aic":                 "wycloudphone",      // CloudPhone
	"cloudcontrol":            "cloudcontrol",      // CloudControl
	"eflo_controller":         "efloctrl",          // Eflo
	"eflo_cnp":                "eflocnp",           // Eflo-CNP
	"milvus":                  "cloudmilvus",       // Milvus
	"resourcedirectorymaster": "resourcedirectory", // ResourceManager
	"openapiexplorer":         "apiexplorer",       // OpenAPIExplorer
	"agentloop":               "agentloop",         // AgentLoop
}

// productCodeToConfigEndpoints records all products' code mapping to endpoints' subfield name
// Key: product code
// Value: Field endpoints' subfield name, such as vpc, ecs, log
var productCodeToConfigEndpoints = map[string]string{
	"sls": "log",
}

// irregularProductEndpoint specially records those product codes that
// cannot be parsed out by the location service.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpoint = map[string]string{
	"tablestore": "tablestore.%s.aliyuncs.com",
	//"ots":              "tablestore.%s.aliyuncs.com", // the old endpoint ots.%s.aliyuncs.com has been deprecated
	"ram":                     "ram.aliyuncs.com",
	"brain_industrial":        "brain-industrial.cn-hangzhou.aliyuncs.com",
	"cassandra":               "cassandra.aliyuncs.com",
	"cloudfw":                 "cloudfw.aliyuncs.com",
	"scdn":                    "scdn.aliyuncs.com",
	"vpcpeer":                 "vpcpeer.aliyuncs.com",
	"resourcecenter":          "resourcecenter.aliyuncs.com",
	"resourcedirectorymaster": "resourcedirectory.aliyuncs.com",
	"market":                  "market.aliyuncs.com",
	"bssopenapi":              BssOpenAPIEndpointDomestic,
	"esa":                     "esa.cn-hangzhou.aliyuncs.com",
	"cas":                     "cas.aliyuncs.com",
	"sas":                     "tds.aliyuncs.com",
	"ros":                     "ros.aliyuncs.com",
	"eds_aic":                 "eds-aic.cn-shanghai.aliyuncs.com",
	"alidns":                  "alidns.aliyuncs.com",
	"openapiexplorer":         "openapi-mcp.cn-hangzhou.aliyuncs.com",
	"computenestsupplier":     "computenestsupplier.cn-hangzhou.aliyuncs.com",
	"agentloop":               "agentloop.%s.aliyuncs.com",
}

// irregularProductEndpointForIntlRegion specially records those product codes that
// cannot be parsed out by the location service and sensitive to region.
// These products adapt to international region, and conflict with irregularProductEndpointForIntlAccount
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpointForIntlRegion = map[string]string{
	"sas": SaSOpenAPIEndpointInternational,
}

// irregularProductEndpointForIntlAccount specially records those product codes that
// cannot be parsed out by the location service and sensitive to account type.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpointForIntlAccount = map[string]string{
	"cloudfw":                 "cloudfw.ap-southeast-1.aliyuncs.com",
	"resourcecenter":          "resourcecenter-intl.aliyuncs.com",
	"resourcedirectorymaster": "resourcedirectory-intl.aliyuncs.com",
	"bssopenapi":              BssOpenAPIEndpointInternational,
	"esa":                     "esa.ap-southeast-1.aliyuncs.com",
	"eds_aic":                 "eds-aic.ap-southeast-1.aliyuncs.com",
	"ros":                     "ros-intl.aliyuncs.com",
	"openapiexplorer":         "openapi-mcp.ap-southeast-1.aliyuncs.com",
	"computenestsupplier":     "computenestsupplier.ap-southeast-1.aliyuncs.com",
}

// irregularProductEndpointForIntlAccountIntlRegion specially records those product codes that
// cannot be parsed out by the location service and sensitive to account type and region.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is higher than location service, lower than user environment variable configuration
var irregularProductEndpointForIntlAccountIntlRegion = map[string]string{
	"cas": "cas.ap-southeast-1.aliyuncs.com",
}

// regularProductEndpoint specially records those product codes that have been confirmed to be
// regional or central endpoints.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpoint = map[string]string{
	"ecs":                  "ecs.%s.aliyuncs.com",
	"mse":                  "mse.%s.aliyuncs.com",
	"vpc":                  "vpc.%s.aliyuncs.com",
	"oss":                  "oss-%s.aliyuncs.com",
	"cr":                   "cr.%s.aliyuncs.com",
	"sts":                  "sts.%s.aliyuncs.com",
	"cms":                  "metrics.%s.aliyuncs.com",
	"sls":                  "%s.log.aliyuncs.com",
	"drds":                 "drds.%s.aliyuncs.com",
	"polardbx":             "polardbx.%s.aliyuncs.com",
	"fc_open":              "%s.fc.aliyuncs.com",
	"fc":                   "%s.fc.aliyuncs.com",
	"cloudapi":             "apigateway.%s.aliyuncs.com",
	"mns_open":             "mns-open.%s.aliyuncs.com",
	"elasticsearch":        "elasticsearch.%s.aliyuncs.com",
	"alikafka":             "alikafka.%s.aliyuncs.com",
	"emr":                  "emr.%s.aliyuncs.com",
	"smartag":              "smartag.%s.aliyuncs.com",
	"alidns":               "alidns.%s.aliyuncs.com",
	"eci":                  "eci.%s.aliyuncs.com",
	"ons":                  "ons.%s.aliyuncs.com",
	"pvtz":                 "pvtz.aliyuncs.com",
	"privatelink":          "privatelink.%s.aliyuncs.com",
	"maxcompute":           "maxcompute.%s.aliyuncs.com",
	"resourcesharing":      "resourcesharing.%s.aliyuncs.com",
	"actiontrail":          "actiontrail.%s.aliyuncs.com",
	"hitsdb":               "hitsdb.%s.aliyuncs.com",
	"oos":                  "oos.%s.aliyuncs.com",
	"nas":                  "nas.%s.aliyuncs.com",
	"dms_enterprise":       "dms-enterprise.%s.aliyuncs.com",
	"kms":                  "kms.%s.aliyuncs.com",
	"eventbridge":          "eventbridge-console.%s.aliyuncs.com",
	"hbr":                  "hbr.%s.aliyuncs.com",
	"arms":                 "arms.%s.aliyuncs.com",
	"sae":                  "sae.%s.aliyuncs.com",
	"alb":                  "alb.%s.aliyuncs.com",
	"ecd":                  "ecd.%s.aliyuncs.com",
	"cloudphone":           "cloudphone.%s.aliyuncs.com",
	"dataworks_public":     "dataworks.%s.aliyuncs.com",
	"ehpc":                 "ehpc.%s.aliyuncs.com",
	"ens":                  "ens.aliyuncs.com",
	"iot":                  "iot.%s.aliyuncs.com",
	"imm":                  "imm.%s.aliyuncs.com",
	"swas_open":            "swas.%s.aliyuncs.com",
	"vs":                   "vs.%s.aliyuncs.com",
	"vod":                  "vod.%s.aliyuncs.com",
	"opensearch":           "opensearch.%s.aliyuncs.com",
	"dbfs":                 "dbfs.%s.aliyuncs.com",
	"eais":                 "eais.%s.aliyuncs.com",
	"servicemesh":          "servicemesh.aliyuncs.com",
	"tag":                  "tag.%s.aliyuncs.com",
	"schedulerx2":          "schedulerx.%s.aliyuncs.com",
	"dbs":                  "dbs-api.%s.aliyuncs.com",
	"nlb":                  "nlb.%s.aliyuncs.com",
	"ebs":                  "ebs.%s.aliyuncs.com",
	"eflo":                 "eflo.%s.aliyuncs.com",
	"eflo_controller":      "eflo-controller.%s.aliyuncs.com",
	"oceanbasepro":         "oceanbasepro.%s.aliyuncs.com",
	"adcp":                 "adcp.%s.aliyuncs.com",
	"rocketmq":             "rocketmq.%s.aliyuncs.com",
	"hologram":             "hologram.%s.aliyuncs.com",
	"foasconsole":          "foasconsole.aliyuncs.com",
	"cs":                   "cs.%s.aliyuncs.com",
	"waf_openapi":          "wafopenapi.cn-hangzhou.aliyuncs.com",
	"dfs":                  "dfs.%s.aliyuncs.com",
	"amqp":                 "amqp-open.%s.aliyuncs.com",
	"amqp_open":            "amqp-open.%s.aliyuncs.com",
	"cbn":                  "cbn.aliyuncs.com",
	"expressconnectrouter": "expressconnectrouter.cn-shanghai.aliyuncs.com",
	"green":                "green.%s.aliyuncs.com",
	"governance":           "governance.cn-hangzhou.aliyuncs.com",
	"dysmsapi":             "dysmsapi.aliyuncs.com",
	"sddp":                 "sddp.cn-zhangjiakou.aliyuncs.com",
	"ddoscoo":              "ddoscoo.cn-hangzhou.aliyuncs.com",
	"config":               "config.cn-shanghai.aliyuncs.com",
	"ga":                   "ga.cn-hangzhou.aliyuncs.com",
	"dcdn":                 "dcdn.aliyuncs.com",
	"cdn":                  "cdn.aliyuncs.com",
	"cloudauth":            "cloudauth.aliyuncs.com",
	"ims":                  "ims.aliyuncs.com",
	"mhub":                 "mhub.cn-shanghai.aliyuncs.com",
	"eds_user":             "eds-user.cn-shanghai.aliyuncs.com",
	"eipanycast":           "eipanycast.cn-hangzhou.aliyuncs.com",
	"mscopensubscription":  "mscopensubscription.aliyuncs.com",
	"resourcemanager":      "resourcemanager.aliyuncs.com",
	"quotas":               "quotas.aliyuncs.com",
	"imp":                  "imp.aliyuncs.com",
	"das":                  "das.cn-shanghai.aliyuncs.com",
	"servicecatalog":       "servicecatalog.cn-hangzhou.aliyuncs.com",
	"chatbot":              "chatbot.cn-shanghai.aliyuncs.com",
	"computenest":          "computenest.cn-hangzhou.aliyuncs.com",
	"aiworkspace":          "aiworkspace.%s.aliyuncs.com",
	"vpcipam":              "vpcipam.%s.aliyuncs.com",
	"gwlb":                 "gwlb.%s.aliyuncs.com",
	"live":                 "live.aliyuncs.com",
	"dts":                  "dts.%s.aliyuncs.com",
	"dg":                   "dg.%s.aliyuncs.com",
	"cloudsso":             "cloudsso.%s.aliyuncs.com",
	"quickbi_public":       "quickbi.%s.aliyuncs.com",
	"ddosbgp":              "ddosbgp.%s.aliyuncs.com",
	"apig":                 "apig.%s.aliyuncs.com",
	"dds":                  "mongodb.%s.aliyuncs.com",
	"cloudcontrol":         "cloudcontrol.aliyuncs.com",
	"bpstudio":             "bpstudio.cn-hangzhou.aliyuncs.com",
}

// regularProductEndpointForIntlRegion specially records those product codes that have been confirmed to be
// regional or central endpoints. But the endpoints are sensitive to region.
// These products adapt to international region, and conflict with regularProductEndpointForIntlAccount
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpointForIntlRegion = map[string]string{
	"ddoscoo":      "ddoscoo.ap-southeast-1.aliyuncs.com",
	"eds_user":     "eds-user.ap-southeast-1.aliyuncs.com",
	"dysmsapi":     "dysmsapi.ap-southeast-1.aliyuncs.com",
	"sddp":         "sddp.ap-southeast-1.aliyuncs.com",
	"governance":   "governance.ap-southeast-1.aliyuncs.com",
	"waf_openapi":  "wafopenapi.ap-southeast-1.aliyuncs.com",
	"cloudcontrol": "cloudcontrol.ap-southeast-1.aliyuncs.com",
	"bpstudio":     "bpstudio.ap-southeast-1.aliyuncs.com",
}

// regularProductEndpointForIntlAccount specially records those product codes that have been confirmed to be
// regional or central endpoints. But the endpoints are sensitive to account type.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpointForIntlAccount = map[string]string{
	"config":              "config.ap-southeast-1.aliyuncs.com",
	"mscopensubscription": "mscopensubscription.ap-southeast-1.aliyuncs.com",
}

// regularProductEndpointForIntlAccountIntlRegion specially records those product codes that have been confirmed to be
// regional or central endpoints. But the endpoints are sensitive to account type and region.
// These products adapt to international account.
// Key: product code, its value equals to the gateway code of the API after converting it to lowercase and using underscores
// Value: product endpoint
// The priority of this configuration is lower than location service, and as a backup endpoint
var regularProductEndpointForIntlAccountIntlRegion = map[string]string{}

// regularProductEndpointReplace specially records some endpoints need to be replaced results from different reasons.
// Key: source endpoint
// Value: replaced endpoint
// The priority of this configuration is lower than other endpoints mapping rules, and as a backup endpoint
var regularProductEndpointReplace = map[string]string{
	// ecs, ecs.aliyuncs.com has more higher speed
	"ecs-cn-hangzhou.aliyuncs.com": "ecs.aliyuncs.com",
}

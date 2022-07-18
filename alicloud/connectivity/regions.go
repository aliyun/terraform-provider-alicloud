package connectivity

// Region represents ECS region
type Region string

// Constants of region definition
const (
	Hangzhou    = Region("cn-hangzhou")
	Qingdao     = Region("cn-qingdao")
	Beijing     = Region("cn-beijing")
	Hongkong    = Region("cn-hongkong")
	Shenzhen    = Region("cn-shenzhen")
	Shanghai    = Region("cn-shanghai")
	Zhangjiakou = Region("cn-zhangjiakou")
	Huhehaote   = Region("cn-huhehaote")
	ChengDu     = Region("cn-chengdu")
	HeYuan      = Region("cn-heyuan")
	WuLanChaBu  = Region("cn-wulanchabu")
	GuangZhou   = Region("cn-guangzhou")
	NanJing     = Region("cn-nanjing")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APNorthEast2 = Region("ap-northeast-2")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")
	APSouthEast6 = Region("ap-southeast-6")
	APSouthEast7 = Region("ap-southeast-7")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1    = Region("me-east-1")
	MECentral1 = Region("me-central-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	RusWest1 = Region("rus-west-1")

	HangzhouFinance     = Region("cn-hangzhou-finance")
	HangzhouFinanceOSS  = Region("cn-hzfinance")
	HangzhouFinanceOSS1 = Region("cn-hzjbp")

	BeijingFinance1   = Region("cn-beijing-finance-1")
	BeijingFinancePub = Region("cn-beijing-finance-1-pub")

	ShanghaiFinance     = Region("cn-shanghai-finance-1")
	ShanghaiFinance1Pub = Region("cn-shanghai-finance-1-pub")
	ShenZhenFinance1    = Region("cn-shenzhen-finance-1")
	ShenzhenFinance2    = Region("cn-szfinance")
	ShenzhenFinance     = Region("cn-shenzhen-finance")

	CnNorth2Gov1 = Region("cn-north-2-gov-1")
)

var ValidRegions = []Region{
	Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, WuLanChaBu, GuangZhou, NanJing,
	USWest1, USEast1,
	APNorthEast1, APNorthEast2, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APSouthEast6, APSouthEast7,
	APSouth1,
	MEEast1, MECentral1,
	EUCentral1, EUWest1,
	CnNorth2Gov1,
	ShenZhenFinance1, ShenzhenFinance2, ShenzhenFinance,
	ShanghaiFinance1Pub, ShanghaiFinance,
	HangzhouFinance, HangzhouFinanceOSS, HangzhouFinanceOSS1,
	BeijingFinance1, BeijingFinancePub,
}

var EcsClassicSupportedRegions = []Region{Shenzhen, Shanghai, Beijing, Qingdao, Hangzhou, Hongkong, USWest1, APSouthEast1}
var EcsSpotNoSupportedRegions = []Region{APSouth1}
var EcsSccSupportedRegions = []Region{Shanghai, Beijing, Zhangjiakou}
var SlbGuaranteedSupportedRegions = []Region{Qingdao, Beijing, Hangzhou, Shanghai, Shenzhen, Zhangjiakou, Huhehaote, APSouthEast1, USEast1}
var DrdsSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Qingdao, Hongkong, Shanghai, Huhehaote, Zhangjiakou, APSouthEast1}
var DrdsClassicNoSupportedRegions = []Region{Hongkong}
var GpdbSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Shanghai, Hongkong}

// Some Ram resources only one can be owned by one account at the same time,
// skipped here to avoid multi regions concurrency conflict.
var RamNoSkipRegions = []Region{Hangzhou, EUCentral1, APSouth1}
var CenNoSkipRegions = []Region{Shanghai, EUCentral1, APSouth1}
var KmsSkippedRegions = []Region{Beijing, Shanghai}

// Actiontrail only one can be owned by one account at the same time,
// skipped here to avoid multi regions concurrency conflict.
var ActiontrailNoSkipRegions = []Region{Hangzhou, EUCentral1, APSouth1}
var FcNoSupportedRegions = []Region{MEEast1}
var DatahubSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, APSouthEast1}
var RdsClassicNoSupportedRegions = []Region{APSouth1, APSouthEast2, APSouthEast3, APNorthEast1, EUCentral1, EUWest1, MEEast1}
var RdsMultiAzNoSupportedRegions = []Region{Qingdao, APNorthEast1, APSouthEast5, MEEast1}
var RdsPPASNoSupportedRegions = []Region{Qingdao, USEast1, APNorthEast1, EUCentral1, MEEast1, APSouthEast2, APSouthEast3, APSouth1, APSouthEast5, ChengDu, EUWest1}
var RouteTableNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var ApiGatewayNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, USEast1, USWest1, EUWest1, MEEast1}
var OtsHighPerformanceNoSupportedRegions = []Region{Qingdao, Zhangjiakou, Huhehaote, APSouthEast2, APSouthEast5, APNorthEast1, EUCentral1, MEEast1}
var OtsCapacityNoSupportedRegions = []Region{APSouthEast1, USWest1, USEast1}
var PrivateIpNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var SwarmSupportedRegions = []Region{Qingdao, Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1, APSouthEast2,
	APSouthEast3, USWest1, USEast1, EUCentral1}
var ManagedKubernetesSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, ChengDu, Hongkong, APSouthEast1, APSouthEast2, EUCentral1, USWest1}
var ServerlessKubernetesSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, APSouthEast1, APSouthEast3, APSouthEast5, APSouth1, Huhehaote}
var KubernetesSupportedRegions = []Region{Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1,
	APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, EUWest1, MEEast1, EUCentral1}
var NasClassicSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ShenZhenFinance1, ShanghaiFinance}
var CasClassicSupportedRegions = []Region{Hangzhou, APSouth1, MEEast1, EUCentral1, APNorthEast1, APSouthEast2}
var CRNoSupportedRegions = []Region{Beijing, Hangzhou, Qingdao, Huhehaote, Zhangjiakou}
var MongoDBClassicNoSupportedRegions = []Region{Huhehaote, Zhangjiakou, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1}
var MongoDBMultiAzSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, EUCentral1}
var DdoscooSupportedRegions = []Region{Hangzhou, APSouthEast1}
var DdosbgpSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, Qingdao, Shanghai, Zhangjiakou, Huhehaote}

//var NetworkAclSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Hongkong, APSouthEast5, APSouth1}
var EssScalingConfigurationMultiSgSupportedRegions = []Region{APSouthEast1, APSouth1}
var SlbClassicNoSupportedRegions = []Region{APNorthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, MEEast1, EUCentral1, EUWest1, Huhehaote, Zhangjiakou}
var NasNoSupportedRegions = []Region{Qingdao, APSouth1, APSouthEast3, APSouthEast5}
var OssVersioningSupportedRegions = []Region{APSouth1}
var OssSseSupportedRegions = []Region{Qingdao, Hangzhou, Beijing, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouth1, USEast1}
var GpdbClassicNoSupportedRegions = []Region{APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1, EUCentral1}
var OnsNoSupportRegions = []Region{APSouthEast5}
var AlikafkaSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, APNorthEast1, APSouthEast1, APSouthEast3, EUCentral1, EUWest1, USEast1, USWest1}
var SmartagSupportedRegions = []Region{Shanghai, ShanghaiFinance, Hongkong, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, EUCentral1, APNorthEast1}
var YundunDbauditSupportedRegions = []Region{Hangzhou, Beijing, Shanghai}
var HttpHttpsHealthCheckMehtodSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, EUWest1, ChengDu, Qingdao, Hongkong, Shenzhen, APSouthEast5, Zhangjiakou, Huhehaote, MEEast1, APSouth1, EUCentral1, USWest1, APSouthEast3, APSouthEast2, APSouthEast1, APNorthEast1}
var HBaseClassicSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen}
var EdasSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen, Zhangjiakou, Qingdao, Hongkong}
var CloudConfigSupportedRegions = []Region{Shanghai}
var DBReadwriteSplittingConnectionSupportedRegions = []Region{APSouthEast1}
var KVstoreClassicNetworkInstanceSupportRegions = []Region{Hangzhou, Beijing, Shanghai, APSouthEast1, USEast1, USWest1}
var MaxComputeSupportRegions = []Region{}
var FnfSupportRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen, USWest1}
var PrivateLinkRegions = []Region{EUCentral1}
var BrainIndustrialRegions = []Region{Hangzhou}
var EciContainerGroupRegions = []Region{Hangzhou}
var TsdbInstanceSupportRegions = []Region{Beijing, Hangzhou, Shenzhen, Shanghai, ShenZhenFinance1, Qingdao, Zhangjiakou, ShanghaiFinance, Hongkong, USWest1, APNorthEast1, EUWest1, APSouthEast1, APSouthEast2, APSouthEast3, EUCentral1, APSouthEast5, Zhangjiakou, CnNorth2Gov1}
var VpcIpv6SupportRegions = []Region{Hangzhou, Shanghai, Shenzhen, Beijing, Huhehaote, Hongkong, APSouthEast1}
var EssdSupportRegions = []Region{Zhangjiakou, Huhehaote}
var AdbReserverUnSupportRegions = []Region{EUCentral1}
var KmsKeyHSMSupportRegions = []Region{Beijing, Zhangjiakou, Hangzhou, Shanghai, Shenzhen, Hongkong, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, USEast1}
var DmSupportRegions = []Region{Hangzhou}
var BssOpenApiSupportRegions = []Region{Hangzhou, Shanghai, APSouthEast1}
var EipAddressBGPProSupportRegions = []Region{Hongkong}
var CenTransitRouterVpcAttachmentSupportRegions = []Region{EUCentral1} // Not all of APSouthEast1 and HangZhou zones support vpc attachment
var ARMSSupportRegions = []Region{Hangzhou, Shanghai, Beijing, APSouthEast1}
var SaeSupportRegions = []Region{Hangzhou, Shanghai, Beijing, Zhangjiakou, Shenzhen, USWest1}
var HbrSupportRegions = []Region{Hangzhou}
var EcdSupportRegions = []Region{Hangzhou, Shanghai, Beijing, Shenzhen, Hongkong, APSouthEast1, APSouthEast2}
var EcpSupportRegions = []Region{Hangzhou, Shanghai, Beijing, Shenzhen}
var SddpSupportRegions = []Region{Hangzhou, Zhangjiakou, APSouthEast1}
var DfsSupportRegions = []Region{Hangzhou, Zhangjiakou, Shanghai, Beijing, HeYuan, ChengDu, APSouthEast5, USEast1, RusWest1}
var EventBridgeSupportRegions = []Region{Hangzhou, Zhangjiakou, Shanghai, Shenzhen, Beijing, HeYuan, ChengDu, Huhehaote, Hongkong, EUCentral1, USWest1, USEast1}
var AlbSupportRegions = []Region{Hangzhou, Shanghai, Qingdao, Zhangjiakou, Beijing, WuLanChaBu, Shenzhen, ChengDu, Hongkong, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APNorthEast1, EUCentral1, USEast1, APSouth1}
var IMMSupportRegions = []Region{Hangzhou, Zhangjiakou, APSouthEast1, Shenzhen, Beijing, Shanghai}
var CenTRSupportRegions = []Region{EUCentral1, APSouthEast1, Hangzhou, Shanghai, Beijing, Shenzhen, Hongkong, APSouthEast1, USEast1, APSouth1}
var VbrSupportRegions = []Region{Hangzhou}
var ClickHouseSupportRegions = []Region{Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, USWest1, USEast1, APSouthEast1, EUCentral1, EUWest1, APNorthEast1, APSouthEast1, APSouthEast5}
var ClickHouseBackupPolicySupportRegions = []Region{Shanghai}
var DatabaseGatewaySupportRegions = []Region{Hangzhou, Zhangjiakou, Shanghai, Beijing, Qingdao, Huhehaote, Shenzhen, ChengDu, Hongkong, APNorthEast1, APSouth1, APSouthEast1, APSouthEast2, APSouthEast3, EUWest1, EUCentral1, APSouthEast5, USWest1, USEast1}
var CloudSsoSupportRegions = []Region{Shanghai, USWest1}
var SWASSupportRegions = []Region{Qingdao, Hangzhou, Beijing, Shenzhen, Shanghai, GuangZhou, Huhehaote, ChengDu, Zhangjiakou, Hongkong, APSouthEast1}
var SurveillanceSystemSupportRegions = []Region{Beijing, Shenzhen, Qingdao}
var VodSupportRegions = []Region{Shanghai}
var OpenSearchSupportRegions = []Region{Beijing, Shenzhen, Hangzhou, Zhangjiakou, Qingdao, Shanghai, APSouthEast1}
var GraphDatabaseSupportRegions = []Region{Shenzhen, Beijing, Qingdao, Shanghai, Hongkong, Zhangjiakou, Hangzhou, APSouthEast1, APSouthEast5, USWest1, USEast1, APSouth1}
var DBFSSystemSupportRegions = []Region{Hangzhou}
var EAISSystemSupportRegions = []Region{Hangzhou}
var CloudAuthSupportRegions = []Region{Hangzhou}
var MHUBSupportRegions = []Region{Shanghai}
var ActiontrailSupportRegions = []Region{Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, WuLanChaBu, GuangZhou, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APNorthEast1, USWest1, USEast1, EUCentral1, EUWest1, APSouth1, MEEast1}
var VpcTrafficMirrorSupportRegions = []Region{Hangzhou, Beijing, Zhangjiakou, Qingdao, Huhehaote, Shenzhen, Hongkong, APSouthEast2, ChengDu, USEast1, USWest1, EUWest1}
var EcdUserSupportRegions = []Region{Shanghai}
var VpcIpv6GatewaySupportRegions = []Region{Qingdao, Beijing, Zhangjiakou, Huhehaote, WuLanChaBu, Hangzhou, Shanghai, Shenzhen, GuangZhou, Hongkong, ChengDu, HeYuan, APSouthEast1, APSouthEast6, USEast1, EUCentral1}
var CmsDynamicTagGroupSupportRegions = []Region{Shanghai}
var OOSApplicationSupportRegions = []Region{Hangzhou}
var DTSSupportRegions = []Region{Hangzhou, APSouth1, ShenZhenFinance1, CnNorth2Gov1, Qingdao, ShanghaiFinance, USWest1, APNorthEast1, Beijing, Hongkong, APSouthEast1, APSouthEast3, EUCentral1, APSouthEast5, Shenzhen, APSouthEast2, Huhehaote, USEast1, Zhangjiakou, EUWest1, MEEast1, Shanghai}
var OOSSupportRegions = []Region{APSouthEast5, USWest1, EUWest1, Qingdao, ChengDu, Shanghai, Huhehaote, Shenzhen, APNorthEast1, APSouthEast1, EUCentral1, Hangzhou, Beijing, APSouth1, APSouthEast3, USEast1, Zhangjiakou, Hongkong, APSouthEast2}
var MongoDBSupportRegions = []Region{APSouth1, Shanghai, APSouthEast2, WuLanChaBu, CnNorth2Gov1, Hangzhou, Beijing, Qingdao, Zhangjiakou, USWest1, GuangZhou, APSouthEast6, EUWest1, ChengDu, APSouthEast1, APSouthEast3, APSouthEast5, ShanghaiFinance, Hongkong, HeYuan, Huhehaote, USEast1, EUCentral1, APNorthEast1, Shenzhen, ShenZhenFinance1, MEEast1}
var MongoDBServerlessSupportRegions = []Region{APSouthEast5, Shanghai, USEast1, Hongkong, HeYuan, Zhangjiakou, APSouthEast6, GuangZhou, Huhehaote, Beijing, Shenzhen, WuLanChaBu, ChengDu, Hangzhou, Qingdao, USWest1, APSouthEast1}
var FnFSupportRegions = []Region{Shenzhen, Beijing, Shanghai, APSouthEast1, USWest1, Hangzhou}
var GaSupportRegions = []Region{Hangzhou}
var AlidnsSupportRegions = []Region{Hangzhou, APSouthEast1}
var VPCVbrHaSupportRegions = []Region{Hangzhou}
var ROSSupportRegions = []Region{USWest1, HeYuan, Zhangjiakou, Hongkong, APSouthEast3, EUCentral1, Huhehaote, APSouthEast6, Shenzhen, APSouth1, Qingdao, GuangZhou, APSouthEast2, WuLanChaBu, EUWest1, MEEast1, ChengDu, Shanghai, APSouthEast1, APSouthEast5, USEast1, Beijing, APNorthEast1, Hangzhou}
var VPCBgpGroupSupportRegions = []Region{Hangzhou}
var NASSupportRegions = []Region{HeYuan, Huhehaote, APSouthEast5, WuLanChaBu, CnNorth2Gov1, Qingdao, ChengDu, Hangzhou, APSouth1, ShenZhenFinance1, EUCentral1, Shenzhen, APSouthEast2, Beijing, Shanghai, ShanghaiFinance, APSouthEast1, APSouthEast6, APNorthEast1, APSouthEast3, GuangZhou, USEast1, EUWest1, Hongkong, Zhangjiakou, USWest1}
var HBRSupportRegions = []Region{Beijing, ChengDu, Huhehaote, Qingdao, Shanghai, Shenzhen, Zhangjiakou, Hangzhou}
var NASCPFSSupportRegions = []Region{Hangzhou, Shenzhen, Beijing, Shanghai, HeYuan, Huhehaote, WuLanChaBu, Qingdao, ChengDu}
var WAFSupportRegions = []Region{Hangzhou, APSouth1}
var MSCSupportRegions = []Region{Hangzhou}

// Other regions requires the custom should have icp
var FCCustomDomainSupportRegions = []Region{EUCentral1, APSouthEast1}
var RDCupportRegions = []Region{Shanghai}
var MSEGatewaySupportRegions = []Region{Shenzhen, Hangzhou, Shanghai, Beijing}
var BrainIndustrialSupportRegions = []Region{Hangzhou}
var TestSalveRegions = []Region{Hangzhou}
var TestPvtzRegions = []Region{Hangzhou}
var ECPSupportRegions = []Region{Beijing, Hangzhou}
var DCDNSupportRegions = []Region{Hangzhou, APSouthEast1, APNorthEast1}
var GpdbElasticInstanceSupportRegions = []Region{EUCentral1, Beijing, Hangzhou, Shanghai, Shenzhen, APSouthEast1, APSouthEast5, Hongkong}
var PolarDBSupportRegions = []Region{Hangzhou}
var ESSSupportRegions = []Region{Beijing}
var SimpleApplicationServerNotSupportRegions = []Region{EUCentral1}
var CRSupportRegions = []Region{WuLanChaBu, APSouthEast2, Hangzhou, ShenZhenFinance1, MEEast1, APSouth1, ShanghaiFinance, APNorthEast1, APSouthEast5, CnNorth2Gov1, Hongkong, Huhehaote, Beijing, ChengDu, APSouthEast3, Shenzhen, USEast1, GuangZhou, Qingdao, Zhangjiakou, EUWest1, Shanghai, APSouthEast1, HeYuan, EUCentral1, USWest1}
var MSESupportRegions = []Region{Zhangjiakou, USWest1, Shenzhen, ChengDu, Qingdao, APSouthEast3, USEast1, Hangzhou, APNorthEast1, ShenZhenFinance1, APSouthEast1, APSouthEast2, APSouthEast5, Beijing, EUWest1, Shanghai, ShanghaiFinance, Huhehaote, APSouth1, CnNorth2Gov1, Hongkong, HeYuan, EUCentral1}
var LogResourceSupportRegions = []Region{HeYuan}
var AliKafkaSupportRegions = []Region{Beijing, CnNorth2Gov1, Qingdao, APSouthEast3, Huhehaote, APSouth1, EUWest1, ShenZhenFinance1, ChengDu, USEast1, USWest1, Hangzhou, Zhangjiakou, Shenzhen, Shanghai, Hongkong, HeYuan, APSouthEast5, APNorthEast1, ShanghaiFinance, APSouthEast1, EUCentral1}
var BastionhostSupportRegions = []Region{CnNorth2Gov1, Qingdao, ShanghaiFinance, EUCentral1, EUWest1, ChengDu, Shanghai, HeYuan, APNorthEast1, MEEast1, APSouth1, Hongkong, Zhangjiakou, USWest1, APSouthEast1, APSouthEast2, Huhehaote, APSouthEast5, Beijing, Hangzhou, ShenZhenFinance1, APSouthEast3, USEast1, Shenzhen}
var ACKSystemDiskEncryptionSupportRegions = []Region{Hongkong}
var DdosBasicSupportRegions = []Region{WuLanChaBu, APSouth1, HeYuan, Shenzhen, MEEast1, APSouthEast1, Huhehaote, CnNorth2Gov1, ChengDu, USEast1, Hangzhou, ShanghaiFinance, ShenZhenFinance1, GuangZhou, APSouthEast2, Beijing, EUCentral1, USWest1, APNorthEast1, Qingdao, APSouthEast3, APSouthEast5, APSouthEast6, Shanghai, Hongkong, Zhangjiakou, EUWest1}
var TagSupportRegions = []Region{Huhehaote, APSouthEast5, CnNorth2Gov1, HeYuan, APSouthEast2, Beijing, APSouthEast3, USWest1, WuLanChaBu, GuangZhou, MEEast1, ShenZhenFinance1, Shanghai, ShanghaiFinance, EUCentral1, APSouthEast1, USEast1, Hangzhou, Hongkong, Qingdao, Zhangjiakou, Shenzhen, EUWest1, APNorthEast1, APSouth1, ChengDu}
var GraphDatabaseDbInstanceSupportRegions = []Region{Hangzhou}
var SchedulerxSupportRegions = []Region{Hangzhou, Shenzhen, Beijing, Shanghai, USEast1, Zhangjiakou}
var EhpcSupportRegions = []Region{Qingdao, HeYuan, Hongkong, APSouthEast1, Huhehaote, Shenzhen, WuLanChaBu, APNorthEast1, EUCentral1, Zhangjiakou, Hangzhou, Beijing, Shanghai, APSouthEast2}
var EcsActivationsSupportRegions = []Region{Qingdao, Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, HeYuan, Hongkong}
var CloudFirewallSupportRegions = []Region{APSouthEast1, Hangzhou}
var SMSSupportRegions = []Region{Hangzhou, APSouthEast1, APSouthEast5}

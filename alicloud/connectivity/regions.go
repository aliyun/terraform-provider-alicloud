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

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1 = Region("me-east-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	ShenZhenFinance = Region("cn-shenzhen-finance-1")
	ShanghaiFinance = Region("cn-shanghai-finance-1")
)

var ValidRegions = []Region{
	Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote,
	USWest1, USEast1,
	APNorthEast1, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5,
	APSouth1,
	MEEast1,
	EUCentral1, EUWest1,
	ShenZhenFinance, ShanghaiFinance,
}

var EcsClassicSupportedRegions = []Region{Shenzhen, Shanghai, Beijing, Qingdao, Hangzhou, Hongkong, USWest1, APSouthEast1}
var EcsSpotNoSupportedRegions = []Region{APSouth1}
var SlbGuaranteedSupportedRegions = []Region{Qingdao, Beijing, Hangzhou, Shanghai, Shenzhen, Zhangjiakou, Huhehaote, APSouthEast1, USEast1}
var DrdsSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Qingdao, Hongkong}
var DrdsClassicNoSupportedRegions = []Region{Hongkong}
var FcNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, APSouthEast3, APSouthEast5, EUWest1, USEast1, MEEast1}
var DatahubSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, APSouthEast1}
var RdsClassicNoSupportedRegions = []Region{APSouth1, APSouthEast2, APSouthEast3, APNorthEast1, EUCentral1, EUWest1, MEEast1}
var RdsMultiAzNoSupportedRegions = []Region{Qingdao, APNorthEast1, APSouthEast5, MEEast1}
var RdsPPASNoSupportedRegions = []Region{Qingdao, USEast1, EUCentral1, MEEast1, APSouthEast2, APSouthEast3, APSouth1, EUWest1}
var RouteTableNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var ApiGatewayNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, USEast1, USWest1, EUWest1, MEEast1}
var OtsHighPerformanceNoSupportedRegions = []Region{Qingdao, Zhangjiakou, Huhehaote, Hongkong, APSouthEast2, APSouthEast5, APNorthEast1, EUCentral1, MEEast1, APSouth1}
var OtsCapacityNoSupportedRegions = []Region{APSouthEast1, USWest1, USEast1}
var PrivateIpNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var SwarmSupportedRegions = []Region{Qingdao, Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1, APSouthEast2,
	APSouthEast3, USWest1, USEast1, EUCentral1}
var ManagedKubernetesSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, APSouthEast1, APSouthEast3, APSouthEast5, APSouth1}
var KubernetesSupportedRegions = []Region{Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1,
	APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, EUWest1, MEEast1, EUCentral1}
var NasClassicSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ShenZhenFinance, ShanghaiFinance}
var CasClassicSupportedRegions = []Region{Hangzhou, APSouth1, MEEast1, EUCentral1, APNorthEast1, APSouthEast2}
var CRNoSupportedRegions = []Region{Beijing, Hangzhou, Qingdao, Huhehaote, Zhangjiakou}
var MongoDBClassicNoSupportedRegions = []Region{Huhehaote, Zhangjiakou, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1}
var MongoDBMultiAzSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, EUCentral1}
var DdoscooSupportedRegions = []Region{Hangzhou}

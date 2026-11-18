package client

import "github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"

// Region represents ECS region
type Region = config.Region

// Constants of region definition
const (
	Hangzhou    = config.Hangzhou
	Qingdao     = config.Qingdao
	Beijing     = config.Beijing
	Hongkong    = config.Hongkong
	Shenzhen    = config.Shenzhen
	Shanghai    = config.Shanghai
	Zhangjiakou = config.Zhangjiakou
	Huhehaote   = config.Huhehaote
	ChengDu     = config.ChengDu
	HeYuan      = config.HeYuan
	WuLanChaBu  = config.WuLanChaBu
	GuangZhou   = config.GuangZhou
	NanJing     = config.NanJing
	FuZhou      = config.FuZhou
	ZhongWei    = config.ZhongWei

	APSouthEast1 = config.APSouthEast1
	APNorthEast1 = config.APNorthEast1
	APNorthEast2 = config.APNorthEast2
	APSouthEast2 = config.APSouthEast2
	APSouthEast3 = config.APSouthEast3
	APSouthEast5 = config.APSouthEast5
	APSouthEast6 = config.APSouthEast6
	APSouthEast7 = config.APSouthEast7
	APSouthEast8 = config.APSouthEast8

	APSouth1 = config.APSouth1

	USWest1      = config.USWest1
	USEast1      = config.USEast1
	USSouthEast1 = config.USSouthEast1

	NASouth1 = config.NASouth1

	MEEast1    = config.MEEast1
	MECentral1 = config.MECentral1

	EUCentral1 = config.EUCentral1
	EUWest1    = config.EUWest1

	RusWest1 = config.RusWest1

	HangzhouFinance     = config.HangzhouFinance
	HangzhouFinanceOSS  = config.HangzhouFinanceOSS
	HangzhouFinanceOSS1 = config.HangzhouFinanceOSS1

	BeijingFinance1   = config.BeijingFinance1
	BeijingFinancePub = config.BeijingFinancePub

	ShanghaiFinance     = config.ShanghaiFinance
	ShanghaiFinance1Pub = config.ShanghaiFinance1Pub
	ShenZhenFinance1    = config.ShenZhenFinance1
	ShenzhenFinance2    = config.ShenzhenFinance2
	ShenzhenFinance     = config.ShenzhenFinance

	CnNorth2Gov1 = config.CnNorth2Gov1
)

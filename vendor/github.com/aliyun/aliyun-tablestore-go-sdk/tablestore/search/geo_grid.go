package search

// GeoGrid 表示地球上的一个区域，它包含一个 TopLeft 和一个 BottomRight 。
// TopLeft 和 BottomRight 组合成一个网格，所以 TopLeft 的 lat 应该大于 BottomRight 的 lat,
// TopLeft 的 lon 应该小于 BottomRight 的 lat。
type GeoGrid struct {
	TopLeft     GeoPoint
	BottomRight GeoPoint
}

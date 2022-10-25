package tablestore

type TimeseriesConfiguration struct {
	metaCacheMaxDataSize int
}

func NewTimeseriesConfiguration()*TimeseriesConfiguration {
	return &TimeseriesConfiguration{
		metaCacheMaxDataSize: 64 * 1024 * 1024,
	}
}

func (this *TimeseriesConfiguration) SetMetaCacheMaxDataSize(metaCacheMaxDataSize int) {
	this.metaCacheMaxDataSize = metaCacheMaxDataSize
}

func (this *TimeseriesConfiguration) GetMetaCacheMaxDataSize() int {
	return this.metaCacheMaxDataSize
}

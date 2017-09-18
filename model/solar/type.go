package solar

type TypeIndicator struct {
	IndicatorBase
}

func (indicator *TypeIndicator) GetEntity() interface{} {
	return indicator
}

package template

//||------------------------------------------------------------------------------------------------||
//|| Template Markers
//||------------------------------------------------------------------------------------------------||

func (t *TemplateInstance) Add(marker string, value string) {
	t.Markers = append(t.Markers, TemplateMarker{
		Marker: marker,
		Value:  value,
	})
}

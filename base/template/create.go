package template

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

func Create(name, language string) TemplateInstance {
	t := TemplateInstance{
		Data:     Templates[name].Data,
		Markers:  []TemplateMarker{},
		Language: language,
	}
	return t
}

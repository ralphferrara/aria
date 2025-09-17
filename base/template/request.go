package template

import "net/http"

func (t *TemplateInstance) Request(req *http.Request) {
	lang := GetLanguageFromRequest(req)
	t.Language = lang
}

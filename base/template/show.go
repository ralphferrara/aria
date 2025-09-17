package template

import "net/http"

func (t *TemplateInstance) Show(w http.ResponseWriter) string {
	content := t.Compile()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
	return content
}

package webhandler

import (
	"html/template"
	"strings"
)

// Convert simple text to html
func FormatHTML(text string) template.HTML {
	ps := strings.Split(text, "\n")
	for i := range ps {
		if strings.HasPrefix(ps[i], ">") {
			ps[i] = "<i>" + ps[i] + "</i>"
		}
	}
	joinedParagraphs := "<p>" + strings.Join(ps, "</p><p>") + "</p>"

	return template.HTML(joinedParagraphs)
}

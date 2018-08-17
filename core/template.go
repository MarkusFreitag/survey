package core

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/mgutz/ansi"
)

var DisableColor = false

type Icon struct {
  Symbol  string
  Color   string
}

var (
	// HelpInputRune is the rune which the user should enter to trigger
	// more detailed question help
	HelpInputIcon = Icon{Symbol: "?", Color: "cyan"}

	// ErrorIcon will be be shown before an error
	ErrorIcon = Icon{Symbol: "X", Color: "red"}

	// HelpIcon will be shown before more detailed question help
	HelpIcon = Icon{Symbol: "????", Color: "cyan"}
	// QuestionIcon will be shown before a question Message
	QuestionIcon = Icon{Symbol: "?", Color: "green+hb"}

	// MarkedOptionIcon will be prepended before a selected multiselect option
	MarkedOptionIcon = Icon{Symbol: "[x]", Color: "green"}
	// UnmarkedOptionIcon will be prepended before an unselected multiselect option
	UnmarkedOptionIcon = Icon{Symbol: "[ ]", Color: "default+hb"}

	// SelectFocusIcon is prepended to an option to signify the user is
	// currently focusing that option
	SelectFocusIcon = Icon{Symbol: ">", Color: "cyan"}
)

/*
  SetFancyIcons changes the err, help, marked, and focus input icons to their
  fancier forms. These forms may not be compatible with most terminals.
  This function will not touch the QuestionIcon as its fancy and non fancy form
  are the same.
*/
func SetFancyIcons() {
	ErrorIcon.Symbol = "✘"
	HelpIcon.Symbol = "ⓘ"
	// QuestionIcon fancy and non-fancy form are the same

	MarkedOptionIcon.Symbol = "◉"
	UnmarkedOptionIcon.Symbol = "◯"

	SelectFocusIcon.Symbol = "❯"
}

var TemplateFuncs = map[string]interface{}{
	// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
	"color": func(color string) string {
		if DisableColor {
			return ""
		}
		return ansi.ColorCode(color)
	},
	"HelpInputIcon": func() string {
		return HelpInputIcon.Symbol
	},
	"ErrorIcon": func() string {
		return ErrorIcon.Symbol
	},
	"HelpIcon": func() string {
		return HelpIcon.Symbol
	},
	"QuestionIcon": func() string {
		return QuestionIcon.Symbol
	},
	"MarkedOptionIcon": func() string {
		return MarkedOptionIcon.Symbol
	},
	"UnmarkedOptionIcon": func() string {
		return UnmarkedOptionIcon.Symbol
	},
	"SelectFocusIcon": func() string {
		return SelectFocusIcon.Symbol
	},
}

var (
	memoizedGetTemplate = map[string]*template.Template{}

	memoMutex = &sync.RWMutex{}
)

func getTemplate(tmpl string) (*template.Template, error) {
	memoMutex.RLock()
	if t, ok := memoizedGetTemplate[tmpl]; ok {
		memoMutex.RUnlock()
		return t, nil
	}
	memoMutex.RUnlock()

	t, err := template.New("prompt").Funcs(TemplateFuncs).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	memoMutex.Lock()
	memoizedGetTemplate[tmpl] = t
	memoMutex.Unlock()
	return t, nil
}

func RunTemplate(tmpl string, data interface{}) (string, error) {
	t, err := getTemplate(tmpl)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

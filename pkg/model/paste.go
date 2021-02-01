package model

import (
	"strings"

	"github.com/go-enry/go-enry/v2"
	"github.com/tzvetkoff-go/errors"
	"github.com/tzvetkoff-go/pasteur/pkg/codemirror"
)

// Paste ...
type Paste struct {
	ID          int    `json:"-"            yaml:"id"`
	IndentStyle string `json:"indent-style" yaml:"indent-style"`
	IndentSize  string `json:"indent-size"  yaml:"indent-size"`
	MimeType    string `json:"mime-type"    yaml:"mime-type"`
	Filename    string `json:"filename"     yaml:"filename"`
	Content     string `json:"content"      yaml:"content"`
}

// NewPaste ...
func NewPaste() *Paste {
	return &Paste{
		ID:          0,
		IndentStyle: "spaces",
		IndentSize:  "4",
		MimeType:    "",
		Filename:    "",
		Content:     "",
	}
}

// Validate ...
func (p *Paste) Validate() error {
	if p.IndentStyle == "" {
		p.IndentStyle = "spaces"
	}
	if p.IndentSize == "" {
		p.IndentSize = "4"
	}
	if p.MimeType == "" {
		lang := enry.GetLanguage(p.Filename, []byte(p.Content))

		println("foobar")
		println(lang)
		println("foobar")

		if lang != "" {
			lang = strings.ToLower(lang)
			for _, mode := range codemirror.Modes {
				if lang == strings.ToLower(mode.Name) {
					if len(mode.MimeTypes) > 0 {
						p.MimeType = mode.MimeTypes[0]
						goto SearchEnd
					}
				}

				for _, alias := range mode.Aliases {
					if lang == strings.ToLower(alias) {
						if len(mode.MimeTypes) > 0 {
							p.MimeType = mode.MimeTypes[0]
							goto SearchEnd
						}
					}
				}
			}
		}
	SearchEnd:

		if p.MimeType == "" {
			p.MimeType = "text/plain"
		}
	}

	if p.IndentStyle != "tabs" && p.IndentStyle != "spaces" {
		return errors.New("indent-style: invalid value")
	}

	if p.IndentSize != "1" &&
		p.IndentSize != "2" &&
		p.IndentSize != "3" &&
		p.IndentSize != "4" &&
		p.IndentSize != "5" &&
		p.IndentSize != "6" &&
		p.IndentSize != "7" &&
		p.IndentSize != "8" {
		return errors.New("paste-indent-size: invalid value")
	}

	mimeTypeOK := false
	for i := 0; i < len(codemirror.Modes); i++ {
		mode := codemirror.Modes[i]
		for j := 0; j < len(mode.MimeTypes); j++ {
			if p.MimeType == mode.MimeTypes[j] {
				mimeTypeOK = true
				break
			}
		}

		if mimeTypeOK {
			break
		}
	}
	if !mimeTypeOK {
		return errors.New("mime-type: unknown mime type")
	}

	if strings.TrimSpace(p.Content) == "" {
		return errors.New("content: cannot be empty")
	}

	return nil
}

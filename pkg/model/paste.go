package model

import (
	"path"
	"strings"
	"time"

	"github.com/tzvetkoff-go/errors"

	"github.com/tzvetkoff-go/pasteur/pkg/indentdb"
	"github.com/tzvetkoff-go/pasteur/pkg/monaco"
	"github.com/tzvetkoff-go/pasteur/pkg/stringutil"
)

// Paste ...
type Paste struct {
	ID          int       `json:"-"`
	Secret      string    `json:"-"`
	Private     int       `json:"private"`
	Filename    string    `json:"filename"`
	Filetype    string    `json:"filetype"`
	IndentStyle string    `json:"indent-style"`
	IndentSize  int       `json:"indent-size"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created-at"`
	UpdatedAt   time.Time `json:"updated-at"`
}

// PasteList ...
type PasteList []*Paste

// PaginatedPasteList ...
type PaginatedPasteList struct {
	Pastes     []*Paste
	Pagination *Pagination
}

// NewPaste ...
func NewPaste() *Paste {
	return &Paste{
		ID:          0,
		Secret:      "",
		Private:     0,
		Filename:    "",
		Filetype:    "",
		IndentStyle: "",
		IndentSize:  0,
		Content:     "",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
}

// Validate ...
func (p *Paste) Validate() error {
	//
	// Generate secret ...
	//

	if p.Secret == "" {
		p.Secret = stringutil.RandString(32, "")
	}

	//
	// Guess filetype ...
	//

	if p.Filetype == "" && p.Filename != "" {
		for _, monacoLanguage := range monaco.Languages {
			pasteExtension := path.Ext(p.Filename)

			for _, filename := range monacoLanguage.Filenames {
				if p.Filename == filename {
					p.Filetype = monacoLanguage.ID
					goto SearchEnd
				}
			}

			for _, extension := range monacoLanguage.Extensions {
				if pasteExtension == extension {
					p.Filetype = monacoLanguage.ID
					goto SearchEnd
				}
			}
		SearchEnd:
		}
	}

	if p.Filetype == "" {
		p.Filetype = "plain"
	}

	//
	// Guess filename ...
	//

	if p.Filename == "" {
		if monacoLanguage, ok := monaco.Languages[p.Filetype]; ok && len(monacoLanguage.Extensions) > 0 {
			p.Filename = "paste" + monacoLanguage.Extensions[0]
		} else {
			p.Filename = "paste.txt"
		}
	}

	//
	// Guess indent ...
	//

	if p.IndentStyle != "" {
		if strings.Index("spaces", strings.ToLower(p.IndentStyle)) == 0 {
			p.IndentStyle = "spaces"
		} else if strings.Index("tabs", strings.ToLower(p.IndentStyle)) == 0 {
			p.IndentStyle = "tabs"
		}
	} else {
		if indent, ok := indentdb.Known[p.Filetype]; ok {
			p.IndentStyle = indent.Style
			p.IndentSize = indent.Size
		}

		if p.IndentStyle == "" {
			p.IndentStyle = "spaces"
		}
	}

	if p.IndentSize == 0 {
		p.IndentSize = 4
	}

	//
	// Perform validations ...
	//

	if p.IndentStyle != "tabs" && p.IndentStyle != "spaces" {
		return errors.New("indent-style: invalid value")
	}

	if p.IndentSize < 0 || p.IndentSize > 8 {
		return errors.New("indent-size: invalid value")
	}

	if p.Filetype != "" {
		filetypeOK := false
		for _, monacoLanguage := range monaco.Languages {
			if p.Filetype == monacoLanguage.ID {
				filetypeOK = true
				break
			}
		}
		if !filetypeOK {
			return errors.New("filetype: unknown filetype")
		}
	}

	if strings.TrimSpace(p.Content) == "" {
		return errors.New("content: cannot be empty")
	}

	return nil
}

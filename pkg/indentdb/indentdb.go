package indentdb

// Indent ...
type Indent struct {
	Style string
	Size  int
}

// Known ...
var Known = map[string]Indent{
	"coffeescript": {"spaces", 2},
	"c":            {"spaces", 4},
	"cpp":          {"spaces", 4},
	"css":          {"spaces", 2},
	"dart":         {"spaces", 2},
	"dockerfile":   {"spaces", 4},
	"elixir":       {"spaces", 2},
	"go":           {"tabs", 4},
	"hcl":          {"spaces", 2},
	"html":         {"spaces", 2},
	"java":         {"spaces", 4},
	"javascript":   {"spaces", 2},
	"less":         {"spaces", 2},
	"lua":          {"spaces", 2},
	"markdown":     {"spaces", 4},
	"objective-c":  {"spaces", 4},
	"perl":         {"spaces", 4},
	"php":          {"spaces", 4},
	"python":       {"spaces", 4},
	"ruby":         {"spaces", 2},
	"rust":         {"spaces", 4},
	"scss":         {"spaces", 2},
	"shell":        {"spaces", 2},
	"swift":        {"spaces", 2},
	"tcl":          {"spaces", 2},
	"yaml":         {"spaces", 2},
}

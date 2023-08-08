package glamour

//go:generate go run ./internal/generate-style-json

import (
	"github.com/rprtr258/tea/glamour/ansi"
)

var (
	// ASCIIStyleConfig uses only ASCII characters.
	ASCIIStyleConfig = ansi.StyleConfig{ //nolint:dupl
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         ptr[uint](1),
			IndentToken:    ptr("| "),
		},
		Paragraph: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
		},
		List: ansi.StyleList{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			LevelIndent: 4,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "# ",
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
			},
		},
		Strikethrough: ansi.StylePrimitive{
			BlockPrefix: "~~",
			BlockSuffix: "~~",
		},
		Emph: ansi.StylePrimitive{
			BlockPrefix: "*",
			BlockSuffix: "*",
		},
		Strong: ansi.StylePrimitive{
			BlockPrefix: "**",
			BlockSuffix: "**",
		},
		HorizontalRule: ansi.StylePrimitive{
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			Ticked:   "[x] ",
			Unticked: "[ ] ",
		},
		ImageText: ansi.StylePrimitive{
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "`",
				BlockSuffix: "`",
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				Margin: ptr[uint](2),
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: ptr("+"),
			ColumnSeparator: ptr("|"),
			RowSeparator:    ptr("-"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n* ",
		},
	}

	// DarkStyleConfig is the default dark style.
	DarkStyleConfig = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       ptr("252"),
			},
			Margin: ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         ptr[uint](1),
			IndentToken:    ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       ptr("39"),
				Bold:        ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           ptr("228"),
				BackgroundColor: ptr("63"),
				Bold:            ptr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  ptr("35"),
				Bold:   ptr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  ptr("240"),
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{},
			Ticked:         "[✓] ",
			Unticked:       "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     ptr("30"),
			Underline: ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: ptr("35"),
			Bold:  ptr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     ptr("212"),
			Underline: ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  ptr("243"),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           ptr("203"),
				BackgroundColor: ptr("236"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: ptr("244"),
				},
				Margin: ptr[uint](2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: ptr("#C4C4C4"),
				},
				Error: ansi.StylePrimitive{
					Color:           ptr("#F1F1F1"),
					BackgroundColor: ptr("#F05B5B"),
				},
				Comment: ansi.StylePrimitive{
					Color: ptr("#676767"),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					Color: ptr("#00AAFF"),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: ptr("#FF5F87"),
				},
				KeywordType: ansi.StylePrimitive{
					Color: ptr("#6E6ED8"),
				},
				Operator: ansi.StylePrimitive{
					Color: ptr("#EF8080"),
				},
				Punctuation: ansi.StylePrimitive{
					Color: ptr("#E8E8A8"),
				},
				Name: ansi.StylePrimitive{
					Color: ptr("#C4C4C4"),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: ptr("#FF8EC7"),
				},
				NameTag: ansi.StylePrimitive{
					Color: ptr("#B083EA"),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: ptr("#7A7AE6"),
				},
				NameClass: ansi.StylePrimitive{
					Color:     ptr("#F1F1F1"),
					Underline: ptr(true),
					Bold:      ptr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: ptr("#FFFF87"),
				},
				NameFunction: ansi.StylePrimitive{
					Color: ptr("#00D787"),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: ptr("#6EEFC0"),
				},
				LiteralString: ansi.StylePrimitive{
					Color: ptr("#C69669"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: ptr("#AFFFD7"),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: ptr("#777777"),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: ptr("#373737"),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			CenterSeparator: ptr("┼"),
			ColumnSeparator: ptr("│"),
			RowSeparator:    ptr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}

	// LightStyleConfig is the default light style.
	LightStyleConfig = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       ptr("234"),
			},
			Margin: ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         ptr[uint](1),
			IndentToken:    ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       ptr("27"),
				Bold:        ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           ptr("228"),
				BackgroundColor: ptr("63"),
				Bold:            ptr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Bold:   ptr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  ptr("249"),
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{},
			Ticked:         "[✓] ",
			Unticked:       "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     ptr("36"),
			Underline: ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: ptr("29"),
			Bold:  ptr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     ptr("205"),
			Underline: ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  ptr("243"),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           ptr("203"),
				BackgroundColor: ptr("254"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: ptr("242"),
				},
				Margin: ptr[uint](2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: ptr("#2A2A2A"),
				},
				Error: ansi.StylePrimitive{
					Color:           ptr("#F1F1F1"),
					BackgroundColor: ptr("#FF5555"),
				},
				Comment: ansi.StylePrimitive{
					Color: ptr("#8D8D8D"),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					Color: ptr("#279EFC"),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: ptr("#FB406F"),
				},
				KeywordType: ansi.StylePrimitive{
					Color: ptr("#7049C2"),
				},
				Operator: ansi.StylePrimitive{
					Color: ptr("#FF2626"),
				},
				Punctuation: ansi.StylePrimitive{
					Color: ptr("#FA7878"),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: ptr("#0A1BB1"),
				},
				NameTag: ansi.StylePrimitive{
					Color: ptr("#581290"),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: ptr("#8362CB"),
				},
				NameClass: ansi.StylePrimitive{
					Color:     ptr("#212121"),
					Underline: ptr(true),
					Bold:      ptr(true),
				},
				NameConstant: ansi.StylePrimitive{
					Color: ptr("#581290"),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: ptr("#A3A322"),
				},
				NameFunction: ansi.StylePrimitive{
					Color: ptr("#019F57"),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: ptr("#22CCAE"),
				},
				LiteralString: ansi.StylePrimitive{
					Color: ptr("#7E5B38"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: ptr("#00AEAE"),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: ptr("#777777"),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: ptr("#373737"),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			CenterSeparator: ptr("┼"),
			ColumnSeparator: ptr("│"),
			RowSeparator:    ptr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}

	// PinkStyleConfig is the default pink style.
	PinkStyleConfig = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			Margin: ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			Indent:      ptr[uint](1),
			IndentToken: ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       ptr("212"),
				Bold:        ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				BlockPrefix: "\n",
				Prefix:      "",
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "▌ ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "┃ ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "│ ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "┆ ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "┊ ",
				Bold:   ptr(false),
			},
		},
		Text: ansi.StylePrimitive{},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  ptr("212"),
			Format: "\n──────\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     ptr("99"),
			Underline: ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			Bold: ptr(true),
		},
		Image: ansi.StylePrimitive{
			Underline: ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			Format: "Image: {{.text}}",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:           ptr("212"),
				BackgroundColor: ptr("236"),
				Prefix:          " ",
				Suffix:          " ",
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: ptr("┼"),
			ColumnSeparator: ptr("│"),
			RowSeparator:    ptr("─"),
		},
		DefinitionList: ansi.StyleBlock{},
		DefinitionTerm: ansi.StylePrimitive{},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
		HTMLBlock: ansi.StyleBlock{},
		HTMLSpan:  ansi.StyleBlock{},
	}

	// NoTTYStyleConfig is the default notty style.
	NoTTYStyleConfig = ansi.StyleConfig{ //nolint:dupl
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         ptr[uint](1),
			IndentToken:    ptr("│ "),
		},
		Paragraph: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
		},
		List: ansi.StyleList{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			LevelIndent: 4,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "# ",
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
			},
		},
		Strikethrough: ansi.StylePrimitive{
			BlockPrefix: "~~",
			BlockSuffix: "~~",
		},
		Emph: ansi.StylePrimitive{
			BlockPrefix: "*",
			BlockSuffix: "*",
		},
		Strong: ansi.StylePrimitive{
			BlockPrefix: "**",
			BlockSuffix: "**",
		},
		HorizontalRule: ansi.StylePrimitive{
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		ImageText: ansi.StylePrimitive{
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "`",
				BlockSuffix: "`",
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				Margin: ptr[uint](2),
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: ptr("┼"),
			ColumnSeparator: ptr("│"),
			RowSeparator:    ptr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}

	// DefaultStyles are the default styles.
	DefaultStyles = map[string]*ansi.StyleConfig{
		ASCIIStyle:   &ASCIIStyleConfig,
		DarkStyle:    &DarkStyleConfig,
		DraculaStyle: &DraculaStyleConfig,
		LightStyle:   &LightStyleConfig,
		NoTTYStyle:   &NoTTYStyleConfig,
		PinkStyle:    &PinkStyleConfig,
	}
)

func ptr[T any](t T) *T { return &t }

package markdown

//go:generate go run ./internal/generate-style-json

import (
	"github.com/rprtr258/tea/components/markdown/ansi"
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
				BlockPrefix:     "\n",
				BlockSuffix:     "\n",
				ForegroundColor: ptr("252"),
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
				BlockSuffix:     "\n",
				ForegroundColor: ptr("39"),
				Bold:            ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: ptr("228"),
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
				Prefix:          "###### ",
				ForegroundColor: ptr("35"),
				Bold:            ptr(false),
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
			ForegroundColor: ptr("240"),
			Format:          "\n--------\n",
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
			ForegroundColor: ptr("30"),
			Underline:       ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			ForegroundColor: ptr("35"),
			Bold:            ptr(true),
		},
		Image: ansi.StylePrimitive{
			ForegroundColor: ptr("212"),
			Underline:       ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			ForegroundColor: ptr("243"),
			Format:          "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: ptr("203"),
				BackgroundColor: ptr("236"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					ForegroundColor: ptr("244"),
				},
				Margin: ptr[uint](2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					ForegroundColor: ptr("#C4C4C4"),
				},
				Error: ansi.StylePrimitive{
					ForegroundColor: ptr("#F1F1F1"),
					BackgroundColor: ptr("#F05B5B"),
				},
				Comment: ansi.StylePrimitive{
					ForegroundColor: ptr("#676767"),
				},
				CommentPreproc: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					ForegroundColor: ptr("#00AAFF"),
				},
				KeywordReserved: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF5F87"),
				},
				KeywordType: ansi.StylePrimitive{
					ForegroundColor: ptr("#6E6ED8"),
				},
				Operator: ansi.StylePrimitive{
					ForegroundColor: ptr("#EF8080"),
				},
				Punctuation: ansi.StylePrimitive{
					ForegroundColor: ptr("#E8E8A8"),
				},
				Name: ansi.StylePrimitive{
					ForegroundColor: ptr("#C4C4C4"),
				},
				NameBuiltin: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF8EC7"),
				},
				NameTag: ansi.StylePrimitive{
					ForegroundColor: ptr("#B083EA"),
				},
				NameAttribute: ansi.StylePrimitive{
					ForegroundColor: ptr("#7A7AE6"),
				},
				NameClass: ansi.StylePrimitive{
					ForegroundColor: ptr("#F1F1F1"),
					Underline:       ptr(true),
					Bold:            ptr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					ForegroundColor: ptr("#FFFF87"),
				},
				NameFunction: ansi.StylePrimitive{
					ForegroundColor: ptr("#00D787"),
				},
				LiteralNumber: ansi.StylePrimitive{
					ForegroundColor: ptr("#6EEFC0"),
				},
				LiteralString: ansi.StylePrimitive{
					ForegroundColor: ptr("#C69669"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					ForegroundColor: ptr("#AFFFD7"),
				},
				GenericDeleted: ansi.StylePrimitive{
					ForegroundColor: ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					ForegroundColor: ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					ForegroundColor: ptr("#777777"),
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
				BlockPrefix:     "\n",
				BlockSuffix:     "\n",
				ForegroundColor: ptr("234"),
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
				BlockSuffix:     "\n",
				ForegroundColor: ptr("27"),
				Bold:            ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: ptr("228"),
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
			ForegroundColor: ptr("249"),
			Format:          "\n--------\n",
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
			ForegroundColor: ptr("36"),
			Underline:       ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			ForegroundColor: ptr("29"),
			Bold:            ptr(true),
		},
		Image: ansi.StylePrimitive{
			ForegroundColor: ptr("205"),
			Underline:       ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			ForegroundColor: ptr("243"),
			Format:          "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: ptr("203"),
				BackgroundColor: ptr("254"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					ForegroundColor: ptr("242"),
				},
				Margin: ptr[uint](2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					ForegroundColor: ptr("#2A2A2A"),
				},
				Error: ansi.StylePrimitive{
					ForegroundColor: ptr("#F1F1F1"),
					BackgroundColor: ptr("#FF5555"),
				},
				Comment: ansi.StylePrimitive{
					ForegroundColor: ptr("#8D8D8D"),
				},
				CommentPreproc: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					ForegroundColor: ptr("#279EFC"),
				},
				KeywordReserved: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					ForegroundColor: ptr("#FB406F"),
				},
				KeywordType: ansi.StylePrimitive{
					ForegroundColor: ptr("#7049C2"),
				},
				Operator: ansi.StylePrimitive{
					ForegroundColor: ptr("#FF2626"),
				},
				Punctuation: ansi.StylePrimitive{
					ForegroundColor: ptr("#FA7878"),
				},
				NameBuiltin: ansi.StylePrimitive{
					ForegroundColor: ptr("#0A1BB1"),
				},
				NameTag: ansi.StylePrimitive{
					ForegroundColor: ptr("#581290"),
				},
				NameAttribute: ansi.StylePrimitive{
					ForegroundColor: ptr("#8362CB"),
				},
				NameClass: ansi.StylePrimitive{
					ForegroundColor: ptr("#212121"),
					Underline:       ptr(true),
					Bold:            ptr(true),
				},
				NameConstant: ansi.StylePrimitive{
					ForegroundColor: ptr("#581290"),
				},
				NameDecorator: ansi.StylePrimitive{
					ForegroundColor: ptr("#A3A322"),
				},
				NameFunction: ansi.StylePrimitive{
					ForegroundColor: ptr("#019F57"),
				},
				LiteralNumber: ansi.StylePrimitive{
					ForegroundColor: ptr("#22CCAE"),
				},
				LiteralString: ansi.StylePrimitive{
					ForegroundColor: ptr("#7E5B38"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					ForegroundColor: ptr("#00AEAE"),
				},
				GenericDeleted: ansi.StylePrimitive{
					ForegroundColor: ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					ForegroundColor: ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					ForegroundColor: ptr("#777777"),
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
				BlockSuffix:     "\n",
				ForegroundColor: ptr("212"),
				Bold:            ptr(true),
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
			ForegroundColor: ptr("212"),
			Format:          "\n──────\n",
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
			ForegroundColor: ptr("99"),
			Underline:       ptr(true),
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
				ForegroundColor: ptr("212"),
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

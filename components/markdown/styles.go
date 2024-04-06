package markdown

//go:generate go run ./internal/generate-style-json

import (
	. "github.com/rprtr258/fun"

	"github.com/rprtr258/tea/components/markdown/ansi"
)

var (
	// ASCIIStyle uses only ASCII characters.
	ASCIIStyle = ansi.StyleConfig{ //nolint:dupl
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: Ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         Ptr[uint](1),
			IndentToken:    Ptr("| "),
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
				Margin: Ptr[uint](2),
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: Ptr("+"),
			ColumnSeparator: Ptr("|"),
			RowSeparator:    Ptr("-"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n* ",
		},
	}

	// DarkStyle is the default dark style.
	DarkStyle = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix:     "\n",
				BlockSuffix:     "\n",
				ForegroundColor: Ptr("252"),
			},
			Margin: Ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         Ptr[uint](1),
			IndentToken:    Ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix:     "\n",
				ForegroundColor: Ptr("39"),
				Bold:            Ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: Ptr("228"),
				BackgroundColor: Ptr("63"),
				Bold:            Ptr(true),
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
				ForegroundColor: Ptr("35"),
				Bold:            Ptr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: Ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: Ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: Ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			ForegroundColor: Ptr("240"),
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
			ForegroundColor: Ptr("30"),
			Underline:       Ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			ForegroundColor: Ptr("35"),
			Bold:            Ptr(true),
		},
		Image: ansi.StylePrimitive{
			ForegroundColor: Ptr("212"),
			Underline:       Ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			ForegroundColor: Ptr("243"),
			Format:          "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: Ptr("203"),
				BackgroundColor: Ptr("236"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					ForegroundColor: Ptr("244"),
				},
				Margin: Ptr[uint](2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					ForegroundColor: Ptr("#C4C4C4"),
				},
				Error: ansi.StylePrimitive{
					ForegroundColor: Ptr("#F1F1F1"),
					BackgroundColor: Ptr("#F05B5B"),
				},
				Comment: ansi.StylePrimitive{
					ForegroundColor: Ptr("#676767"),
				},
				CommentPreproc: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					ForegroundColor: Ptr("#00AAFF"),
				},
				KeywordReserved: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF5F87"),
				},
				KeywordType: ansi.StylePrimitive{
					ForegroundColor: Ptr("#6E6ED8"),
				},
				Operator: ansi.StylePrimitive{
					ForegroundColor: Ptr("#EF8080"),
				},
				Punctuation: ansi.StylePrimitive{
					ForegroundColor: Ptr("#E8E8A8"),
				},
				Name: ansi.StylePrimitive{
					ForegroundColor: Ptr("#C4C4C4"),
				},
				NameBuiltin: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF8EC7"),
				},
				NameTag: ansi.StylePrimitive{
					ForegroundColor: Ptr("#B083EA"),
				},
				NameAttribute: ansi.StylePrimitive{
					ForegroundColor: Ptr("#7A7AE6"),
				},
				NameClass: ansi.StylePrimitive{
					ForegroundColor: Ptr("#F1F1F1"),
					Underline:       Ptr(true),
					Bold:            Ptr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FFFF87"),
				},
				NameFunction: ansi.StylePrimitive{
					ForegroundColor: Ptr("#00D787"),
				},
				LiteralNumber: ansi.StylePrimitive{
					ForegroundColor: Ptr("#6EEFC0"),
				},
				LiteralString: ansi.StylePrimitive{
					ForegroundColor: Ptr("#C69669"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					ForegroundColor: Ptr("#AFFFD7"),
				},
				GenericDeleted: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: Ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					ForegroundColor: Ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: Ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					ForegroundColor: Ptr("#777777"),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: Ptr("#373737"),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			CenterSeparator: Ptr("┼"),
			ColumnSeparator: Ptr("│"),
			RowSeparator:    Ptr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}

	// LightStyle is the default light style.
	LightStyle = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix:     "\n",
				BlockSuffix:     "\n",
				ForegroundColor: Ptr("234"),
			},
			Margin: Ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         Ptr[uint](1),
			IndentToken:    Ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix:     "\n",
				ForegroundColor: Ptr("27"),
				Bold:            Ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: Ptr("228"),
				BackgroundColor: Ptr("63"),
				Bold:            Ptr(true),
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
				Bold:   Ptr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: Ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: Ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: Ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			ForegroundColor: Ptr("249"),
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
			ForegroundColor: Ptr("36"),
			Underline:       Ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			ForegroundColor: Ptr("29"),
			Bold:            Ptr(true),
		},
		Image: ansi.StylePrimitive{
			ForegroundColor: Ptr("205"),
			Underline:       Ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			ForegroundColor: Ptr("243"),
			Format:          "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				ForegroundColor: Ptr("203"),
				BackgroundColor: Ptr("254"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					ForegroundColor: Ptr("242"),
				},
				Margin: Ptr[uint](2),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					ForegroundColor: Ptr("#2A2A2A"),
				},
				Error: ansi.StylePrimitive{
					ForegroundColor: Ptr("#F1F1F1"),
					BackgroundColor: Ptr("#FF5555"),
				},
				Comment: ansi.StylePrimitive{
					ForegroundColor: Ptr("#8D8D8D"),
				},
				CommentPreproc: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					ForegroundColor: Ptr("#279EFC"),
				},
				KeywordReserved: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FB406F"),
				},
				KeywordType: ansi.StylePrimitive{
					ForegroundColor: Ptr("#7049C2"),
				},
				Operator: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FF2626"),
				},
				Punctuation: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FA7878"),
				},
				NameBuiltin: ansi.StylePrimitive{
					ForegroundColor: Ptr("#0A1BB1"),
				},
				NameTag: ansi.StylePrimitive{
					ForegroundColor: Ptr("#581290"),
				},
				NameAttribute: ansi.StylePrimitive{
					ForegroundColor: Ptr("#8362CB"),
				},
				NameClass: ansi.StylePrimitive{
					ForegroundColor: Ptr("#212121"),
					Underline:       Ptr(true),
					Bold:            Ptr(true),
				},
				NameConstant: ansi.StylePrimitive{
					ForegroundColor: Ptr("#581290"),
				},
				NameDecorator: ansi.StylePrimitive{
					ForegroundColor: Ptr("#A3A322"),
				},
				NameFunction: ansi.StylePrimitive{
					ForegroundColor: Ptr("#019F57"),
				},
				LiteralNumber: ansi.StylePrimitive{
					ForegroundColor: Ptr("#22CCAE"),
				},
				LiteralString: ansi.StylePrimitive{
					ForegroundColor: Ptr("#7E5B38"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					ForegroundColor: Ptr("#00AEAE"),
				},
				GenericDeleted: ansi.StylePrimitive{
					ForegroundColor: Ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: Ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					ForegroundColor: Ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: Ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					ForegroundColor: Ptr("#777777"),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: Ptr("#373737"),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			CenterSeparator: Ptr("┼"),
			ColumnSeparator: Ptr("│"),
			RowSeparator:    Ptr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}

	// PinkStyle is the default pink style.
	PinkStyle = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			Margin: Ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			Indent:      Ptr[uint](1),
			IndentToken: Ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix:     "\n",
				ForegroundColor: Ptr("212"),
				Bold:            Ptr(true),
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
				Bold:   Ptr(false),
			},
		},
		Text: ansi.StylePrimitive{},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: Ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: Ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: Ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			ForegroundColor: Ptr("212"),
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
			ForegroundColor: Ptr("99"),
			Underline:       Ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			Bold: Ptr(true),
		},
		Image: ansi.StylePrimitive{
			Underline: Ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			Format: "Image: {{.text}}",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				ForegroundColor: Ptr("212"),
				BackgroundColor: Ptr("236"),
				Prefix:          " ",
				Suffix:          " ",
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: Ptr("┼"),
			ColumnSeparator: Ptr("│"),
			RowSeparator:    Ptr("─"),
		},
		DefinitionList: ansi.StyleBlock{},
		DefinitionTerm: ansi.StylePrimitive{},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
		HTMLBlock: ansi.StyleBlock{},
		HTMLSpan:  ansi.StyleBlock{},
	}

	// NoTTYStyle is the default notty style.
	NoTTYStyle = ansi.StyleConfig{ //nolint:dupl
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
			},
			Margin: Ptr[uint](2),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         Ptr[uint](1),
			IndentToken:    Ptr("│ "),
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
				Margin: Ptr[uint](2),
			},
		},
		Table: ansi.StyleTable{
			CenterSeparator: Ptr("┼"),
			ColumnSeparator: Ptr("│"),
			RowSeparator:    Ptr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}
)

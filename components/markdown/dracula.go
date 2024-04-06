package markdown

import (
	. "github.com/rprtr258/fun"

	"github.com/rprtr258/tea/components/markdown/ansi"
)

var DraculaStyle = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix:     "\n",
			BlockSuffix:     "\n",
			ForegroundColor: Ptr("#f8f8f2"),
		},
		Margin: Ptr[uint](2),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			ForegroundColor: Ptr("#f1fa8c"),
			Italic:          Ptr(true),
		},
		Indent: Ptr[uint](2),
	},
	List: ansi.StyleList{
		LevelIndent: 2,
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				ForegroundColor: Ptr("#f8f8f2"),
			},
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix:     "\n",
			ForegroundColor: Ptr("#bd93f9"),
			Bold:            Ptr(true),
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
		CrossedOut: Ptr(true),
	},
	Emph: ansi.StylePrimitive{
		ForegroundColor: Ptr("#f1fa8c"),
		Italic:          Ptr(true),
	},
	Strong: ansi.StylePrimitive{
		Bold:            Ptr(true),
		ForegroundColor: Ptr("#ffb86c"),
	},
	HorizontalRule: ansi.StylePrimitive{
		ForegroundColor: Ptr("#6272A4"),
		Format:          "\n--------\n",
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
	},
	Enumeration: ansi.StylePrimitive{
		BlockPrefix:     ". ",
		ForegroundColor: Ptr("#8be9fd"),
	},
	Task: ansi.StyleTask{
		StylePrimitive: ansi.StylePrimitive{},
		Ticked:         "[✓] ",
		Unticked:       "[ ] ",
	},
	Link: ansi.StylePrimitive{
		ForegroundColor: Ptr("#8be9fd"),
		Underline:       Ptr(true),
	},
	LinkText: ansi.StylePrimitive{
		ForegroundColor: Ptr("#ff79c6"),
	},
	Image: ansi.StylePrimitive{
		ForegroundColor: Ptr("#8be9fd"),
		Underline:       Ptr(true),
	},
	ImageText: ansi.StylePrimitive{
		ForegroundColor: Ptr("#ff79c6"),
		Format:          "Image: {{.text}} →",
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			ForegroundColor: Ptr("#50fa7b"),
		},
	},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ffb86c"),
			},
			Margin: Ptr[uint](2),
		},
		Chroma: &ansi.Chroma{
			Text: ansi.StylePrimitive{
				ForegroundColor: Ptr("#f8f8f2"),
			},
			Error: ansi.StylePrimitive{
				ForegroundColor: Ptr("#f8f8f2"),
				BackgroundColor: Ptr("#ff5555"),
			},
			Comment: ansi.StylePrimitive{
				ForegroundColor: Ptr("#6272A4"),
			},
			CommentPreproc: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			Keyword: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			KeywordReserved: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			KeywordNamespace: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			KeywordType: ansi.StylePrimitive{
				ForegroundColor: Ptr("#8be9fd"),
			},
			Operator: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			Punctuation: ansi.StylePrimitive{
				ForegroundColor: Ptr("#f8f8f2"),
			},
			Name: ansi.StylePrimitive{
				ForegroundColor: Ptr("#8be9fd"),
			},
			NameBuiltin: ansi.StylePrimitive{
				ForegroundColor: Ptr("#8be9fd"),
			},
			NameTag: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			NameAttribute: ansi.StylePrimitive{
				ForegroundColor: Ptr("#50fa7b"),
			},
			NameClass: ansi.StylePrimitive{
				ForegroundColor: Ptr("#8be9fd"),
			},
			NameConstant: ansi.StylePrimitive{
				ForegroundColor: Ptr("#bd93f9"),
			},
			NameDecorator: ansi.StylePrimitive{
				ForegroundColor: Ptr("#50fa7b"),
			},
			NameFunction: ansi.StylePrimitive{
				ForegroundColor: Ptr("#50fa7b"),
			},
			LiteralNumber: ansi.StylePrimitive{
				ForegroundColor: Ptr("#6EEFC0"),
			},
			LiteralString: ansi.StylePrimitive{
				ForegroundColor: Ptr("#f1fa8c"),
			},
			LiteralStringEscape: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff79c6"),
			},
			GenericDeleted: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ff5555"),
			},
			GenericEmph: ansi.StylePrimitive{
				ForegroundColor: Ptr("#f1fa8c"),
				Italic:          Ptr(true),
			},
			GenericInserted: ansi.StylePrimitive{
				ForegroundColor: Ptr("#50fa7b"),
			},
			GenericStrong: ansi.StylePrimitive{
				ForegroundColor: Ptr("#ffb86c"),
				Bold:            Ptr(true),
			},
			GenericSubheading: ansi.StylePrimitive{
				ForegroundColor: Ptr("#bd93f9"),
			},
			Background: ansi.StylePrimitive{
				BackgroundColor: Ptr("#282a36"),
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

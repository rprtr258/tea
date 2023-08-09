package glamour

import "github.com/rprtr258/tea/glamour/ansi"

var DraculaStyleConfig = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       ptr("#f8f8f2"),
		},
		Margin: ptr[uint](2),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:  ptr("#f1fa8c"),
			Italic: ptr(true),
		},
		Indent: ptr[uint](2),
	},
	List: ansi.StyleList{
		LevelIndent: 2,
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: ptr("#f8f8f2"),
			},
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       ptr("#bd93f9"),
			Bold:        ptr(true),
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
		CrossedOut: ptr(true),
	},
	Emph: ansi.StylePrimitive{
		Color:  ptr("#f1fa8c"),
		Italic: ptr(true),
	},
	Strong: ansi.StylePrimitive{
		Bold:  ptr(true),
		Color: ptr("#ffb86c"),
	},
	HorizontalRule: ansi.StylePrimitive{
		Color:  ptr("#6272A4"),
		Format: "\n--------\n",
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
	},
	Enumeration: ansi.StylePrimitive{
		BlockPrefix: ". ",
		Color:       ptr("#8be9fd"),
	},
	Task: ansi.StyleTask{
		StylePrimitive: ansi.StylePrimitive{},
		Ticked:         "[✓] ",
		Unticked:       "[ ] ",
	},
	Link: ansi.StylePrimitive{
		Color:     ptr("#8be9fd"),
		Underline: ptr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: ptr("#ff79c6"),
	},
	Image: ansi.StylePrimitive{
		Color:     ptr("#8be9fd"),
		Underline: ptr(true),
	},
	ImageText: ansi.StylePrimitive{
		Color:  ptr("#ff79c6"),
		Format: "Image: {{.text}} →",
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: ptr("#50fa7b"),
		},
	},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: ptr("#ffb86c"),
			},
			Margin: ptr[uint](2),
		},
		Chroma: &ansi.Chroma{
			Text: ansi.StylePrimitive{
				Color: ptr("#f8f8f2"),
			},
			Error: ansi.StylePrimitive{
				Color:           ptr("#f8f8f2"),
				BackgroundColor: ptr("#ff5555"),
			},
			Comment: ansi.StylePrimitive{
				Color: ptr("#6272A4"),
			},
			CommentPreproc: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			Keyword: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			KeywordReserved: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			KeywordNamespace: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			KeywordType: ansi.StylePrimitive{
				Color: ptr("#8be9fd"),
			},
			Operator: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			Punctuation: ansi.StylePrimitive{
				Color: ptr("#f8f8f2"),
			},
			Name: ansi.StylePrimitive{
				Color: ptr("#8be9fd"),
			},
			NameBuiltin: ansi.StylePrimitive{
				Color: ptr("#8be9fd"),
			},
			NameTag: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			NameAttribute: ansi.StylePrimitive{
				Color: ptr("#50fa7b"),
			},
			NameClass: ansi.StylePrimitive{
				Color: ptr("#8be9fd"),
			},
			NameConstant: ansi.StylePrimitive{
				Color: ptr("#bd93f9"),
			},
			NameDecorator: ansi.StylePrimitive{
				Color: ptr("#50fa7b"),
			},
			NameFunction: ansi.StylePrimitive{
				Color: ptr("#50fa7b"),
			},
			LiteralNumber: ansi.StylePrimitive{
				Color: ptr("#6EEFC0"),
			},
			LiteralString: ansi.StylePrimitive{
				Color: ptr("#f1fa8c"),
			},
			LiteralStringEscape: ansi.StylePrimitive{
				Color: ptr("#ff79c6"),
			},
			GenericDeleted: ansi.StylePrimitive{
				Color: ptr("#ff5555"),
			},
			GenericEmph: ansi.StylePrimitive{
				Color:  ptr("#f1fa8c"),
				Italic: ptr(true),
			},
			GenericInserted: ansi.StylePrimitive{
				Color: ptr("#50fa7b"),
			},
			GenericStrong: ansi.StylePrimitive{
				Color: ptr("#ffb86c"),
				Bold:  ptr(true),
			},
			GenericSubheading: ansi.StylePrimitive{
				Color: ptr("#bd93f9"),
			},
			Background: ansi.StylePrimitive{
				BackgroundColor: ptr("#282a36"),
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

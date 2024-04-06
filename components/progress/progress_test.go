package progress

const (
	AnsiReset = "\x1b[0m"
)

// func TestGradient(t *testing.T) {
// 	colAHex := "#FF0000"
// 	colA, _ := colorful.Hex(colAHex)
// 	colBHex := "#00FF00"
// 	colB, _ := colorful.Hex(colBHex)

// 	for name, opts := range map[string][]Option{
// 		"progress bar with gradient": {
// 			WithoutPercentage(),
// 			WithGradient(colA, colB),
// 		},
// 		"progress bar with scaled gradient": {
// 			WithoutPercentage(),
// 			WithScaledGradient(colA, colB),
// 		},
// 	} {
// 		t.Run(name, func(t *testing.T) {
// 			p := New(opts...)

// 			// build the expected colors by colorizing an empty string and then cutting off the following reset sequence
// 			sb := strings.Builder{}
// 			sb.WriteString(scuf.String("", scuf.FgRGB(scuf.MustParseHexRGB(colAHex))))
// 			expFirst := strings.Split(sb.String(), AnsiReset)[0]
// 			sb.Reset()
// 			sb.WriteString(scuf.String("", scuf.FgRGB(scuf.MustParseHexRGB(colBHex))))
// 			expLast := strings.Split(sb.String(), AnsiReset)[0]

// 			for _, width := range []int{3, 5, 50} {
// 				p.Width = width
// 				res := p.ViewAs(1.0)

// 				// extract colors from the progrss bar by splitting at p.Full+AnsiReset, leaving us with just the color sequences
// 				colors := strings.Split(res, string(p.Full)+AnsiReset)

// 				// discard the last color, because it is empty (no new color comes after the last char of the bar)
// 				colors = colors[0 : len(colors)-1]

// 				assert.Equal(t, expFirst, colors[0])
// 				assert.Equal(t, expLast, colors[len(colors)-1])
// 			}
// 		})
// 	}
// }

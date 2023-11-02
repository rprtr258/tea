package layout

// This example demonstrates various Lip Gloss style and layout features.

import (
	"context"
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/rprtr258/tea/styles"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30
)

// Style definitions.
var (
	// General.
	subtle = styles.FgAdaptiveColor("#D9DCCF", "#383838")
	// highlight   = styles.FgAdaptiveColor("#874BFD", "#7D56F4")
	special     = styles.FgAdaptiveColor("#43BF6D", "#73F59F")
	highlightBg = styles.FgAdaptiveColor("#874BFD", "#7D56F4")

	divider = styles.Style{}.
		SetString("•").
		// Padding(0, 1).
		Foreground(subtle).
		String()

	url = styles.Style{}.Foreground(special).Render

	// Tabs.
	// activeTabBorder = box.Border{
	// 	Top:         '─',
	// 	Bottom:      ' ',
	// 	Left:        '│',
	// 	Right:       '│',
	// 	TopLeft:     '╭',
	// 	TopRight:    '╮',
	// 	BottomLeft:  '┘',
	// 	BottomRight: '└',
	// }

	// tabBorder = box.Border{
	// 	Top:         '─',
	// 	Bottom:      '─',
	// 	Left:        '│',
	// 	Right:       '│',
	// 	TopLeft:     '╭',
	// 	TopRight:    '╮',
	// 	BottomLeft:  '┴',
	// 	BottomRight: '┴',
	// }

	tab = styles.Style{}
	// Border(tabBorder, true).
	// BorderForeground(highlight)
	// Padding(0, 1)

	activeTab = tab.Copy() // .Border(activeTabBorder, true)

	tabGap = tab.Copy()
	// BorderTop(false).
	// BorderLeft(false).
	// BorderRight(false)

	// Title.
	titleStyle = styles.Style{}.
		// MarginLeft(1).
		// MarginRight(5).
		// Padding(0, 1).
		Italic().
		Foreground(styles.FgColor("#FFF7DB")).
		SetString("Lip Gloss")

	descStyle = styles.Style{} // .MarginTop(1)

	infoStyle = styles.Style{}
	// BorderStyle(styles.NormalBorder).
	// BorderTop(true).
	// BorderForeground(subtle)

	// Dialog.
	dialogBoxStyle = styles.Style{}
	// Border(styles.RoundedBorder).
	// BorderForeground(styles.FgColor("#874BFD")).
	// Padding(1, 0).
	// BorderTop(true).
	// BorderLeft(true).
	// BorderRight(true).
	// BorderBottom(true)

	buttonStyle = styles.Style{}.
			Foreground(styles.FgColor("#FFF7DB")).
			Background(styles.BgColor("#888B7E"))
		// Padding(0, 3).
		// MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(styles.FgColor("#FFF7DB")).
				Background(styles.BgColor("#F25D94")).
		// MarginRight(2).
		Underline()

	// List.
	list = styles.Style{}
	// Border(styles.NormalBorder, false, true, false, false).
	// BorderForeground(subtle).
	// MarginRight(2).
	// Height(8).
	// Width(columnWidth + 1)
	listHeader = styles.Style{}.
		// BorderStyle(styles.NormalBorder).
		// BorderBottom(true).
		// BorderForeground(subtle).
		// MarginRight(2).
		Render
	listItem = styles.Style{}.
		// PaddingLeft(2).
		Render

	checkMark = styles.Style{}.SetString("✓").
			Foreground(special).
		// PaddingRight(1).
		String()

	listDone = func(s string) string {
		return checkMark + styles.Style{}.
			Strikethrough(true).
			Foreground(styles.FgAdaptiveColor("#969B86", "#696969")).
			Render(s)
	}

	// Paragraphs/History.
	historyStyle = styles.Style{}.
			Align(styles.Left).
			Foreground(styles.FgColor("#FAFAFA")).
			Background(highlightBg)
		// Margin(1, 3, 0, 0).
		// Padding(1, 2).
		// Height(19).
		// Width(columnWidth)

	// Status Bar.
	statusNugget = styles.Style{}.
			Foreground(styles.FgColor("#FFFDF5"))
		// Padding(0, 1)
	statusBarStyle = styles.Style{}.
			Foreground(styles.FgAdaptiveColor("#343433", "#C1C6B2")).
			Background(styles.BgAdaptiveColor("#D9DCCF", "#353533"))
	statusStyle = styles.Style{}.
			Inherit(statusBarStyle).
			Foreground(styles.FgColor("#FFFDF5")).
			Background(styles.BgColor("#FF5F87"))
		// Padding(0, 1).
		// MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(styles.BgColor("#A550DF")).
			Align(styles.Right)

	statusText = styles.Style{}.Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(styles.BgColor("#6124DF"))

	// Page.
	docStyle = styles.Style{}
	// Padding(1, 2, 1, 2)
)

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

func Main(context.Context) error {
	// physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	// Tabs
	{
		row := styles.JoinHorizontal(
			styles.Top,
			strings.Split(activeTab.Render("Lip Gloss"), "\n"),
			strings.Split(tab.Render("Blush"), "\n"),
			strings.Split(tab.Render("Eye Shadow"), "\n"),
			strings.Split(tab.Render("Mascara"), "\n"),
			strings.Split(tab.Render("Foundation"), "\n"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-styles.Width(row)-2)))
		row = styles.JoinHorizontal(styles.Bottom, strings.Split(row, "\n"), strings.Split(gap, "\n"))
		doc.WriteString(row + "\n\n")
	}

	// Title
	{
		var (
			colors = colorGrid(1, 5)
			title  strings.Builder
		)

		for i, v := range colors {
			// const offset = 2
			fmt.Fprint(&title, titleStyle.Copy(). /*.MarginLeft(i*offset)*/ Background(styles.BgColor(v[0])))
			if i < len(colors)-1 {
				title.WriteRune('\n')
			}
		}

		desc := styles.JoinVertical(styles.Left,
			descStyle.Render("Style Definitions for Nice Terminal Layouts"),
			infoStyle.Render("From Charm"+divider+url("https://github.com/rprtr258/tea/styles")),
		)

		row := styles.JoinHorizontal(styles.Top, strings.Split(title.String(), "\n"), strings.Split(desc, "\n"))
		doc.WriteString(row + "\n\n")
	}

	// Dialog
	{
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")

		question := styles.Style{}. /*.Width(50)*/ Align(styles.Center).Render("Are you sure you want to eat marmalade?")
		buttons := styles.JoinHorizontal(styles.Top, strings.Split(okButton, "\n"), strings.Split(cancelButton, "\n"))
		ui := styles.JoinVertical(styles.Center, question, buttons)

		dialog := styles.Place(width, 9,
			styles.Center, styles.Center,
			dialogBoxStyle.Render(ui),
			styles.WithWhitespaceChars("猫咪"),
			styles.WithWhitespaceForeground(subtle),
		)

		doc.WriteString(dialog + "\n\n")
	}

	// Color grid
	colors := func() string {
		colors := colorGrid(14, 8)

		b := strings.Builder{}
		for _, x := range colors {
			for _, y := range x {
				s := styles.Style{}.SetString("  ").Background(styles.BgColor(y))
				b.WriteString(s.String())
			}
			b.WriteRune('\n')
		}

		return b.String()
	}()

	lists := styles.JoinHorizontal(styles.Top,
		strings.Split(list.Render(
			styles.JoinVertical(styles.Left,
				listHeader("Citrus Fruits to Try"),
				listDone("Grapefruit"),
				listDone("Yuzu"),
				listItem("Citron"),
				listItem("Kumquat"),
				listItem("Pomelo"),
			),
		), "\n"),
		strings.Split(list.Copy(). /*.Width(columnWidth)*/ Render(
			styles.JoinVertical(styles.Left,
				listHeader("Actual Lip Gloss Vendors"),
				listItem("Glossier"),
				listItem("Claire‘s Boutique"),
				listDone("Nyx"),
				listItem("Mac"),
				listDone("Milk"),
			),
		), "\n"),
	)

	doc.WriteString(styles.JoinHorizontal(styles.Top, strings.Split(lists, "\n"), strings.Split(colors, "\n")))

	// Marmalade history
	{
		//nolint:lll
		const (
			historyA = `The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos.`
			historyB = `Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac.`
			historyC = `In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting.`
		)

		doc.WriteString(styles.JoinHorizontal(
			styles.Top,
			strings.Split(historyStyle.Copy().Align(styles.Right).Render(historyA), "\n"),
			strings.Split(historyStyle.Copy().Align(styles.Center).Render(historyB), "\n"),
			strings.Split(historyStyle.Copy(). /*.MarginRight(0)*/ Render(historyC), "\n"),
		))

		doc.WriteString("\n\n")
	}

	// Status bar
	{
		// w := styles.Width

		statusKey := statusStyle.Render("STATUS")
		encoding := encodingStyle.Render("UTF-8")
		fishCake := fishCakeStyle.Render("🍥 Fish Cake")
		statusVal := statusText.Copy().
			// Width(width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("Ravishing")

		bar := styles.JoinHorizontal(styles.Top,
			strings.Split(statusKey, "\n"),
			strings.Split(statusVal, "\n"),
			strings.Split(encoding, "\n"),
			strings.Split(fishCake, "\n"),
		)

		doc.WriteString(statusBarStyle. /*.Width(width)*/ Render(bar))
	}

	// if physicalWidth > 0 {
	// 	docStyle = docStyle.MaxWidth(physicalWidth)
	// }

	// Okay, let's print it
	fmt.Println(docStyle.Render(doc.String()))
	return nil
}

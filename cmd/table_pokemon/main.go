package table_pokemon

import (
	"context"

	"github.com/rprtr258/scuf"
	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/tablebox"
	"github.com/rprtr258/tea/styles"
)

type model struct{}

func (*model) Init(func(...tea.Cmd)) {}

func (*model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.MsgKey:
		switch msg.String() {
		case "q", "ctrl+c":
			f(tea.Quit)
			return
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	for i, vb := range tablebox.Box(
		vb.MaxWidth(80).MaxHeight(32),
		tablebox.NormalBorder,
		styles.Style{}.Foreground(styles.ANSIColor(238)),
		tablebox.Grid(
			[]tea.Layout{tea.Fixed(1), tea.Fixed(28)}, // TODO: fixed(1), auto
			[]tea.Layout{tea.Fixed(6), tea.Fixed(14), tea.Fixed(12), tea.Fixed(10), tea.Fixed(14), tea.Fixed(17)},
		)...,
	) {
		y, x := i%2, i/2
		if y == 0 { // header
			headerStyle := styles.Style{}.Foreground(styles.ANSIColor(252)).Bold(true)
			headers := []string{"#", "NAME", "TYPE 1", "TYPE 2", "JAPANESE", "OFFICIAL ROM."}
			vb.PaddingLeft(1).Styled(headerStyle).WriteLine(headers[x])

			continue
		}

		selectedStyle := styles.Style{}.Foreground(styles.FgRGB("#01BE85")).Background(styles.BgRGB("#00432F"))
		typeColors := map[string][2]scuf.Modifier{ // normal, dimmed
			"Bug":      {styles.FgRGB("#D7FF87"), styles.FgRGB("#97AD64")},
			"Electric": {styles.FgRGB("#FDFF90"), styles.FgRGB("#FCFF5F")},
			"Fire":     {styles.FgRGB("#FF7698"), styles.FgRGB("#BA5F75")},
			"Grass":    {styles.FgRGB("#75FBAB"), styles.FgRGB("#59B980")},
			"Ground":   {styles.FgRGB("#FF875F"), styles.FgRGB("#C77252")},
			"Poison":   {styles.FgRGB("#7D5AFC"), styles.FgRGB("#634BD0")},
			"Water":    {styles.FgRGB("#00E2C7"), styles.FgRGB("#439F8E")},

			"Normal": {styles.FgRGB("#929292"), styles.FgRGB("#727272")},
			"Flying": {styles.FgRGB("#FF87D7"), styles.FgRGB("#C97AB2")},
		}
		textColors := [2]scuf.Modifier{styles.ANSIColor(245), styles.ANSIColor(252)}

		for y, cur := range [][6]string{
			{"1", "Bulbasaur", "Grass", "Poison", "フシギダネ", "Bulbasaur"},
			{"2", "Ivysaur", "Grass", "Poison", "フシギソウ", "Ivysaur"},
			{"3", "Venusaur", "Grass", "Poison", "フシギバナ", "Venusaur"},
			{"4", "Charmander", "Fire", "", "ヒトカゲ", "Hitokage"},
			{"5", "Charmeleon", "Fire", "", "リザード", "Lizardo"},
			{"6", "Charizard", "Fire", "Flying", "リザードン", "Lizardon"},
			{"7", "Squirtle", "Water", "", "ゼニガメ", "Zenigame"},
			{"8", "Wartortle", "Water", "", "カメール", "Kameil"},
			{"9", "Blastoise", "Water", "", "カメックス", "Kamex"},
			{"10", "Caterpie", "Bug", "", "キャタピー", "Caterpie"},
			{"11", "Metapod", "Bug", "", "トランセル", "Trancell"},
			{"12", "Butterfree", "Bug", "Flying", "バタフリー", "Butterfree"},
			{"13", "Weedle", "Bug", "Poison", "ビードル", "Beedle"},
			{"14", "Kakuna", "Bug", "Poison", "コクーン", "Cocoon"},
			{"15", "Beedrill", "Bug", "Poison", "スピアー", "Spear"},
			{"16", "Pidgey", "Normal", "Flying", "ポッポ", "Poppo"},
			{"17", "Pidgeotto", "Normal", "Flying", "ピジョン", "Pigeon"},
			{"18", "Pidgeot", "Normal", "Flying", "ピジョット", "Pigeot"},
			{"19", "Rattata", "Normal", "", "コラッタ", "Koratta"},
			{"20", "Raticate", "Normal", "", "ラッタ", "Ratta"},
			{"21", "Spearow", "Normal", "Flying", "オニスズメ", "Onisuzume"},
			{"22", "Fearow", "Normal", "Flying", "オニドリル", "Onidrill"},
			{"23", "Ekans", "Poison", "", "アーボ", "Arbo"},
			{"24", "Arbok", "Poison", "", "アーボック", "Arbok"},
			{"25", "Pikachu", "Electric", "", "ピカチュウ", "Pikachu"},
			{"26", "Raichu", "Electric", "", "ライチュウ", "Raichu"},
			{"27", "Sandshrew", "Ground", "", "サンド", "Sand"},
			{"28", "Sandslash", "Ground", "", "サンドパン", "Sandpan"},
		} {
			vbRow := vb.Row(y)
			switch x {
			case 2, 3:
				vbRow = vbRow.Styled(styles.Style{}.Foreground(typeColors[cur[x]][1-y%2]))
			default:
				vbRow = vbRow.Styled(styles.Style{}.Foreground(textColors[1-y%2]))
			}

			if cur[1] == "Pikachu" {
				vbRow = vbRow.Styled(selectedStyle)
			}

			vbRow.PaddingLeft(1).WriteLine(cur[x])
		}
	}
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}

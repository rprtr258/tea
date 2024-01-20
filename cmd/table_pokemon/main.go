package table_pokemon

import (
	"context"

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
	tablebox.Box(
		vb.MaxWidth(80).MaxHeight(32),
		[]tea.Layout{tea.Fixed(1), tea.Fixed(28)}, // TODO: fixed(1), auto
		[]tea.Layout{tea.Fixed(6), tea.Fixed(14), tea.Fixed(12), tea.Fixed(10), tea.Fixed(14), tea.Fixed(17)},
		func(vb tea.Viewbox, y, x int) {
			if y == 0 { // header
				headerStyle := styles.Style{}.Foreground(styles.ANSIColor(252)).Bold(true)
				headers := []string{"#", "Name", "Type 1", "Type 2", "Japanese", "Official Rom."}
				vb.Styled(headerStyle).WriteLine(headers[x])

				return
			}

			// data := [][]string{
			// 	{"1", "Bulbasaur", "Grass", "Poison", "フシギダネ", "Bulbasaur"},
			// 	{"2", "Ivysaur", "Grass", "Poison", "フシギソウ", "Ivysaur"},
			// 	{"3", "Venusaur", "Grass", "Poison", "フシギバナ", "Venusaur"},
			// 	{"4", "Charmander", "Fire", "", "ヒトカゲ", "Hitokage"},
			// 	{"5", "Charmeleon", "Fire", "", "リザード", "Lizardo"},
			// 	{"6", "Charizard", "Fire", "Flying", "リザードン", "Lizardon"},
			// 	{"7", "Squirtle", "Water", "", "ゼニガメ", "Zenigame"},
			// 	{"8", "Wartortle", "Water", "", "カメール", "Kameil"},
			// 	{"9", "Blastoise", "Water", "", "カメックス", "Kamex"},
			// 	{"10", "Caterpie", "Bug", "", "キャタピー", "Caterpie"},
			// 	{"11", "Metapod", "Bug", "", "トランセル", "Trancell"},
			// 	{"12", "Butterfree", "Bug", "Flying", "バタフリー", "Butterfree"},
			// 	{"13", "Weedle", "Bug", "Poison", "ビードル", "Beedle"},
			// 	{"14", "Kakuna", "Bug", "Poison", "コクーン", "Cocoon"},
			// 	{"15", "Beedrill", "Bug", "Poison", "スピアー", "Spear"},
			// 	{"16", "Pidgey", "Normal", "Flying", "ポッポ", "Poppo"},
			// 	{"17", "Pidgeotto", "Normal", "Flying", "ピジョン", "Pigeon"},
			// 	{"18", "Pidgeot", "Normal", "Flying", "ピジョット", "Pigeot"},
			// 	{"19", "Rattata", "Normal", "", "コラッタ", "Koratta"},
			// 	{"20", "Raticate", "Normal", "", "ラッタ", "Ratta"},
			// 	{"21", "Spearow", "Normal", "Flying", "オニスズメ", "Onisuzume"},
			// 	{"22", "Fearow", "Normal", "Flying", "オニドリル", "Onidrill"},
			// 	{"23", "Ekans", "Poison", "", "アーボ", "Arbo"},
			// 	{"24", "Arbok", "Poison", "", "アーボック", "Arbok"},
			// 	{"25", "Pikachu", "Electric", "", "ピカチュウ", "Pikachu"},
			// 	{"26", "Raichu", "Electric", "", "ライチュウ", "Raichu"},
			// 	{"27", "Sandshrew", "Ground", "", "サンド", "Sand"},
			// 	{"28", "Sandslash", "Ground", "", "サンドパン", "Sandpan"},
			// }

			// selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
			// typeColors := map[string]lipgloss.Color{
			// 	"Bug":      lipgloss.Color("#D7FF87"),
			// 	"Electric": lipgloss.Color("#FDFF90"),
			// 	"Fire":     lipgloss.Color("#FF7698"),
			// 	"Flying":   lipgloss.Color("#FF87D7"),
			// 	"Grass":    lipgloss.Color("#75FBAB"),
			// 	"Ground":   lipgloss.Color("#FF875F"),
			// 	"Normal":   lipgloss.Color("#929292"),
			// 	"Poison":   lipgloss.Color("#7D5AFC"),
			// 	"Water":    lipgloss.Color("#00E2C7"),
			// }
			// dimTypeColors := map[string]lipgloss.Color{
			// 	"Bug":      lipgloss.Color("#97AD64"),
			// 	"Electric": lipgloss.Color("#FCFF5F"),
			// 	"Fire":     lipgloss.Color("#BA5F75"),
			// 	"Flying":   lipgloss.Color("#C97AB2"),
			// 	"Grass":    lipgloss.Color("#59B980"),
			// 	"Ground":   lipgloss.Color("#C77252"),
			// 	"Normal":   lipgloss.Color("#727272"),
			// 	"Poison":   lipgloss.Color("#634BD0"),
			// 	"Water":    lipgloss.Color("#439F8E"),
			// }
		},
		tablebox.NormalBorder,
		styles.Style{}.Foreground(styles.ANSIColor(238)),
	)
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}

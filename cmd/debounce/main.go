package debounce

// This example illustrates how to debounce commands.
//
// When the user presses a key we increment the "tag" value on the model and,
// after a short delay, we include that tag value in the message produced
// by the Tick command.
//
// In a subsequent Update, if the tag in the Msg matches current tag on the
// model's state we know that the debouncing is complete and we can proceed as
// normal. If not, we simply ignore the inbound message.

import (
	"context"
	"fmt"
	"time"

	"github.com/rprtr258/tea"
)

const debounceDuration = time.Second

type msgExit int

type model struct {
	tag int
}

func (m *model) Init(func(...tea.Cmd)) {}

func (m *model) Update(msg tea.Msg, f func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		// Increment the tag on the model...
		m.tag++
		f(tea.Tick(debounceDuration, func(_ time.Time) tea.Msg {
			// ...and include a copy of that tag value in the message.
			return msgExit(m.tag)
		}))
		return
	case msgExit:
		// If the tag in the message doesn't match the tag on the model then we
		// know that this message was not the last one sent and another is on
		// the way. If that's the case we know, we can ignore this message.
		// Otherwise, the debounce timeout has passed and this message is a
		// valid debounced one.
		if int(msg) == m.tag {
			f(tea.Quit)
			return
		}
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb.
		WriteLineX("Key presses: ").
		WriteLineX(fmt.Sprint(m.tag))
	vb.PaddingTop(1).WriteLine("To exit press any key, then wait for one second without pressing anything.")
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, &model{}).Run()
	return err
}

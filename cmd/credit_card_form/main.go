package credit_card_form //nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/textinput"
	"github.com/rprtr258/tea/styles"
)

const (
	_ccn = iota
	_exp
	_cvv
)

var (
	inputStyle    = styles.Style{}.Foreground(styles.FgColor("#FF06B7")) // hot pink
	continueStyle = styles.Style{}.Foreground(styles.FgColor("#767676")) // dark gray
)

type model struct {
	inputs  []textinput.Model
	focused int
}

// Validator functions to ensure valid input
func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if s == "" || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")

	if _, err := strconv.ParseInt(e, 10, 64); err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func initialModel() *model {
	ccnInput := textinput.New()
	ccnInput.Placeholder = "4505 **** **** 1234"
	ccnInput.Focus()
	ccnInput.CharLimit = 20
	ccnInput.Width = 30
	ccnInput.Prompt = ""
	ccnInput.Validate = ccnValidator

	expInput := textinput.New()
	expInput.Placeholder = "MM/YY "
	expInput.CharLimit = 5
	expInput.Width = 5
	expInput.Prompt = ""
	expInput.Validate = expValidator

	cvvInput := textinput.New()
	cvvInput.Placeholder = "XXX"
	cvvInput.CharLimit = 3
	cvvInput.Width = 5
	cvvInput.Prompt = ""
	cvvInput.Validate = cvvValidator

	return &model{
		inputs: []textinput.Model{
			_ccn: ccnInput,
			_exp: expInput,
			_cvv: cvvInput,
		},
		focused: 0,
	}
}

func (m *model) Init(yield func(...tea.Cmd)) {
	yield(textinput.Blink)
}

func (m *model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	if msg, ok := msg.(tea.MsgKey); ok {
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				yield(tea.Quit)
				return
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			yield(tea.Quit)
			return
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	for i := range m.inputs {
		m.inputs[i].Update(msg, yield)
	}
}

func (m *model) View(vb tea.Viewbox) {
	vb = vb.PaddingLeft(1)

	vb.WriteLine("Total: $21.50:")
	vb = vb.PaddingTop(2)
	vb.Styled(inputStyle).WriteLine("Card Number")
	vb = vb.PaddingTop(1)
	m.inputs[_ccn].View(vb)

	vb = vb.PaddingTop(2)

	vbExp := vb.Sub(tea.Rectangle{
		Height: 2,
		Width:  8,
	})
	vbExp.Styled(inputStyle).WriteLine("EXP")
	m.inputs[_exp].View(vb.PaddingTop(1))

	vbCvv := vb.Sub(tea.Rectangle{
		Left:   8,
		Height: 2,
		Width:  8,
	})
	vbCvv.Styled(inputStyle).WriteLine("CVV")
	m.inputs[_cvv].View(vbCvv.PaddingTop(1))

	vb.PaddingTop(3).WriteLine(continueStyle.Render("Continue ->"))
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused++
	if m.focused >= len(m.inputs)-1 {
		m.focused = len(m.inputs) - 1
	}
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}

package pager

type Pager struct {
	total, perPage int
	selected       int // always in 0..total whenever total > 0
}

func New(
	total, perPage int,
) *Pager {
	return &Pager{
		total:    total,
		perPage:  perPage,
		selected: 0,
	}
}

// Total is the total number of pages
func (p *Pager) Total() int { return p.total }
func (p *Pager) TotalSet(total int) {
	p.total = total
	p.selected = min(p.selected, p.total-1)
}

// PerPage is the number of items per page
func (p *Pager) PerPage() int           { return p.perPage }
func (p *Pager) PerPageSet(perPage int) { p.perPage = perPage }

// Page is the current page number
func (p *Pager) Page() int  { return p.selected / p.perPage }
func (p *Pager) Pages() int { return (p.total + p.perPage - 1) / p.perPage }
func (p *Pager) PageSet(page int) {
	switch {
	case page < 0:
		p.selected = 0
	case page >= p.Pages():
		p.selected = p.total - 1
	default:
		p.selected = page * p.perPage
	}
}
func (p *Pager) PagePrev()  { p.PageSet(p.Page() - 1) }
func (p *Pager) PageNext()  { p.PageSet(p.Page() + 1) }
func (p *Pager) PageFirst() { p.PageSet(0) }
func (p *Pager) PageLast()  { p.PageSet(p.Pages() - 1) }

func (p *Pager) SelectValue() int { return p.selected }
func (p *Pager) SelectNext() {
	if p.selected < p.total {
		p.selected++
	}
}
func (p *Pager) SelectPrev() {
	if p.selected > 0 {
		p.selected--
	}
}

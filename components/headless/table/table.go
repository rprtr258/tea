package table

import "slices"

type Column[T any] struct {
	Title string
	Get   func(T) string
	Cmp   func(T, T) int
}

type Table[T any] struct {
	cols   []Column[T]
	order  []int
	rows   []T
	cursor int
}

func New[T any](cols []Column[T], rows []T) *Table[T] {
	order := make([]int, len(rows))
	for i := range order {
		order[i] = i
	}

	return &Table[T]{
		cols:   cols,
		rows:   rows,
		order:  order,
		cursor: -1,
	}
}

func (t *Table[T]) RowsCount() int { return len(t.rows) }
func (t *Table[T]) Rows() []T {
	rows := make([]T, len(t.rows))
	for i := range rows {
		rows[i] = t.rows[t.order[i]]
	}
	return rows
}
func (t *Table[T]) Columns() []Column[T] { return t.cols }

func (t *Table[T]) Cursor() int     { return t.cursor }
func (t *Table[T]) Selected() T     { return t.rows[t.cursor] }
func (t *Table[T]) MoveTo(y int)    { t.cursor = max(0, min(len(t.rows)-1, y)) }
func (t *Table[T]) MoveTop()        { t.MoveTo(0) }
func (t *Table[T]) MoveBottom()     { t.MoveTo(len(t.rows) - 1) }
func (t *Table[T]) MoveUp(dy int)   { t.MoveTo(t.cursor - dy) }
func (t *Table[T]) MoveDown(dy int) { t.MoveTo(t.cursor + dy) }

// TODO: move current?
func (t *Table[T]) SortColumn(i int) {
	slices.SortFunc(t.order, func(i, j int) int {
		return t.cols[i].Cmp(t.rows[t.order[i]], t.rows[t.order[j]])
	})
}
func (t *Table[T]) SortColumnDesc(i int) {
	slices.SortFunc(t.order, func(i, j int) int {
		return -t.cols[i].Cmp(t.rows[t.order[i]], t.rows[t.order[j]])
	})
}

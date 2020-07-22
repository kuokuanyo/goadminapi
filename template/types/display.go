package types

type FieldDisplay struct {
	Display              FieldFilterFn
	DisplayProcessChains DisplayProcessFnChains
}

type DisplayProcessFnChains []FieldFilterFn

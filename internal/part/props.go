package part

type PropSpec struct {
	Items []PropItem
}

type PropItem struct {
	Label   string
	Kind    int
	Value   any
	Choices []string
	Min     int
	Max     int
}

type PropAction struct {
	Index    int
	NewValue any
}

const (
	PropText = iota
	PropInt
	PropChoice
	PropBool
	PropActionButton
)

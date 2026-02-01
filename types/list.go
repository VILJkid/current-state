package types

type ListItem struct {
	PrimaryText   string
	SecondaryText string
	Shortcut      rune
	Action        func()
	Err           error
}

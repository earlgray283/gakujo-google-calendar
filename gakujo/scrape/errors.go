package scrape

type ErrNotFound struct {
	Name string
}

func (e ErrNotFound) Error() string {
	return "not found: " + e.Name
}

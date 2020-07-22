package parameter

type Parameters struct {
	Page        string
	PageInt     int
	PageSize    string
	PageSizeInt int
	SortField   string
	Columns     []string
	SortType    string
	Animation   bool
	URLPath     string
	Fields      map[string][]string
}

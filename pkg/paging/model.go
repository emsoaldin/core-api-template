package paging

// Model model
type Model struct {
	Results   interface{} `json:"results"`
	Paginator *Paginator  `json:"paginator"`
}

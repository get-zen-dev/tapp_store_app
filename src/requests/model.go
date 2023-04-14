package requests

// Model is a box for storing the name, path and url
type Model struct {
	Name string
	Path string
	URL  string
}

func newModel(name, path, url string) Model {
	return Model{name, path, url}
}

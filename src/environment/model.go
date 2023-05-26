package environment

// Model is a box for storing the name, path and description
type Model struct {
	Name        string
	Version     string
	Description string
}

// Response is a box for storing a slice of Models
type Models struct {
	slice []Model
}

// Return a slice of Model
func (r *Models) Value() []Model {
	return r.slice
}

func (r *Models) append(m Model) {
	r.slice = append(r.slice, m)
}

package requests

// Response is a box for storing a slice of Models
type Response struct {
	slice []Model
}

func (r *Response) append(m Model) {
	r.slice = append(r.slice, m)
}

// Return a slice of Models
func (r *Response) Value() []Model {
	return r.slice
}

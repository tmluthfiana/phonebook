package helper

type Result struct {
	Data    interface{}
	Message string
	Total   int
}

func NewResult() *Result {
	r := new(Result)

	return r
}

func (r *Result) SetData(d interface{}) *Result {
	r.Data = d

	return r
}

func (r *Result) SetMessage(m string) *Result {
	r.Message = m

	return r
}

func (r *Result) SetTotal(i int) *Result {
	r.Total = i

	return r
}

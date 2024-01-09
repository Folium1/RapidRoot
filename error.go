package rapidroot

import "net/http"

const (
	internalServerErr = "internal server error"
)

func (r *Request) abortWithErr(code int, err error) {
	r.isAborted = true
	http.Error(r.Writer, err.Error(), code)
}

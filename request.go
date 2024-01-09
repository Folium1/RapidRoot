package rapidroot

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sync"
)

// Request is a struct for handlers, to interact with request.
type Request struct {
	Writer http.ResponseWriter
	Req    *http.Request

	// used to prevent usage of the shared memory of the data by multiple goroutines
	mu *sync.Mutex

	// used to save data in Request, and then it can be used in other handlers or middlewares
	data map[string]any

	// data from r.Req.URL.Query().
	queryValues url.Values

	// used for log, to print the name of the function
	handlerName string

	// used to interact with cookies
	cookie *cookies

	// used to abort request
	isAborted bool
}

var requestPool sync.Pool

func init() {
	requestPool = sync.Pool{
		New: func() interface{} {
			return &Request{}
		},
	}
}

// GetRequest retrieves a Request from the sync pool.
func getRequest(w http.ResponseWriter, req *http.Request) *Request {
	request := requestPool.Get().(*Request)
	request.Writer = w
	request.Req = req
	request.data = make(map[string]any)
	request.queryValues = req.URL.Query()

	return request
}

func (r *Request) reset() {
	r.Writer = nil
	r.Req = nil
	r.data = nil
	r.queryValues = nil
	r.cookie = nil
	r.handlerName = ""
	r.isAborted = false
}

// ReleaseRequest releases a Request back to the sync pool.
func releaseRequest(request *Request) {
	request.reset()
	requestPool.Put(request)
}
func newRequest(writer http.ResponseWriter, req *http.Request) *Request {
	// Initialize a new Request struct
	newRequest := &Request{
		Writer:      writer,
		Req:         req,
		mu:          new(sync.Mutex),
		data:        make(map[string]interface{}),
		queryValues: req.URL.Query(),
		cookie: &cookies{
			defaults: &http.Cookie{
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			},
		},
		handlerName: "",
		isAborted:   false,
	}

	return newRequest
}

/*
//////////////////////////////////
!!!!!Request manipulations!!!!!!!!
//////////////////////////////////
*/

// SetValue puts key-value to the Request.
func (r *Request) SetValue(key string, val any) {
	r.data[key] = val
}

// Value returns value set to the Request struct.
func (r *Request) Value(key string) any {
	return r.data[key]
}

// Values returns all values set to Request struct.
func (r *Request) Values() map[string]any {
	return r.data
}

// SetStatus should be used,if you don't use response functions of this package.
func (r *Request) SetStatus(code int) {
	r.Writer.WriteHeader(code)
}

func (r *Request) GetStatus() int {
	return r.Writer.(*responseCodeWrapper).statusCode
}

// QueryValue returns value from query.
func (r *Request) QueryValue(key string) string {
	return r.queryValues.Get(key)
}

// QueryValues returns all values from query.
func (r *Request) QueryValues() url.Values {
	return r.queryValues
}

// PostFormValues returns all values from post form.
func (r *Request) PostFormValues() url.Values {
	return r.Req.PostForm
}

// PostFormVal returns value from post form.
func (r *Request) PostFormVal(key string) string {
	return r.Req.FormValue(key)
}

// IsAborted returns true if request is aborted.
func (r *Request) IsAborted() bool {
	return r.isAborted
}

// Abort aborts request.
func (r *Request) Abort() {
	r.isAborted = true
}

/*
//////////////////////////
!!!!!!!!!RESPONSE!!!!!!!!!
//////////////////////////
*/

// Redirect redirects request to another url.
// Only codes from 300 to 308 are valid.
func (r *Request) Redirect(code int, url string) {
	http.Redirect(r.Writer, r.Req, url, code)
	return
}

// ERROR return an error with status code.
func (r *Request) ERROR(code int, err error) {
	http.Error(r.Writer, err.Error(), code)
}

// JSON parses data to json format and sends response with a provided code.
func (r *Request) JSON(code int, data any) {
	r.writeJSON(code, data)
}

// XML parses data to xml format and sends response with a provided code.
func (r *Request) XML(code int, data any) {
	r.writeXML(code, data)
}

// XMLIndent parses data to xml format and sends response with a provided code.
func (r *Request) XMLIndent(code int, data any, prefix, indent string) {
	r.writeXMLIndent(code, data, prefix, indent)
}

// HTML parses data to HTML format and sends a response with the provided code.
// If there is no file with such name, it will abort with a 500 error status code.
func (r *Request) HTML(code int, name string, data any) {
	if !fileExists(name) {
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		log.error(fmt.Sprintf("File: %s doesn't exist", name), r.handlerName)
		return
	}

	tmpl, err := template.ParseFiles(name)
	if err != nil {
		log.error(fmt.Sprintf("Failed to parse HTML file: %s, err: %s", name, err.Error()), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}

	r.writeHTML(code, *tmpl, data)
}

// HTMLTemplate same as HTML, but you can put your html template to execute.
// You need to set template name to template.Template struct. And then pass it to this function.
//
// Example:
//
//	  tmpl, err := template.ParseFiles("main.html", "footer.html")
//		if err != nil {
//			// handle err
//		}
//
//	  r.HTMLTemplate(200,"main.html", tmpl, data)
//
// If there is no file with such name, will abort with 500 error status code.
func (r *Request) HTMLTemplate(code int, templateName string, tmpl *template.Template, data any) {
	if templateName == "" {
		log.error("templateName is empty", r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
	r.writeHTMLTemplate(code, templateName, *tmpl, data)
}

// BINARY response with binary data and provided code.
func (r *Request) BINARY(code int, data []byte) {
	r.writeBINARY(code, data)
}

// FILE response with file and provided code.
func (r *Request) FILE(code int, fileName string) {
	if !fileExists(fileName) {
		log.error(fmt.Sprintf("File: %s doesn't exist", fileName), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
	r.writeFILE(code, fileName)
}

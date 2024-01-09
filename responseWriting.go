package rapidroot

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func (r *Request) writeJSON(code int, data any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.SetStatus(code)
	err := json.NewEncoder(r.Writer).Encode(data)
	if err != nil {
		log.error(fmt.Sprint("failed to convert data to JSON: %w", err), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
}

func (r *Request) writeXML(code int, data any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	xmlData, err := xml.Marshal(data)
	if err != nil {
		log.error(fmt.Sprint("failed to convert data to XML: %w", err), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}

	r.Writer.Header().Set("Content-Type", "application/xml")
	r.SetStatus(code)
	_, err = r.Writer.Write(xmlData)
	if err != nil {
		log.error(err.Error(), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
}

func (r *Request) writeXMLIndent(code int, data any, prefix, indent string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	xmlData, err := xml.MarshalIndent(data, prefix, indent)
	if err != nil {
		log.error(fmt.Sprint("failed to convert data to XML: %w", err), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}

	r.Writer.Header().Set("Content-Type", "application/xml")
	r.SetStatus(code)
	_, err = r.Writer.Write(xmlData)
	if err != nil {
		log.error(err.Error(), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
}

func (r *Request) writeHTML(code int, templ template.Template, data any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Writer.Header().Set("Content-Type", "text/html")

	r.SetStatus(code)
	err := templ.Execute(r.Writer, data)
	if err != nil {
		log.error(fmt.Sprint("failed to execute HTML template: %w", err), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
}

func (r *Request) writeHTMLTemplate(code int, name string, templ template.Template, data any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Writer.Header().Set("Content-Type", "text/html")

	r.SetStatus(code)
	err := templ.ExecuteTemplate(r.Writer, name, data)
	if err != nil {
		log.error(fmt.Sprint("failed to execute HTML template: %w", err), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
}

func (r *Request) writeBINARY(code int, data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.SetStatus(code)
	_, err := r.Writer.Write(data)
	if err != nil {
		log.error(err.Error(), r.handlerName)
		r.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
		return
	}
}

func (req *Request) writeFILE(code int, name string) {
	req.mu.Lock()
	defer req.mu.Unlock()

	req.Writer.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", name),
	)

	req.Writer.Header().Set(
		"Content-Type",
		req.Req.Header.Get("Content-Type"),
	)

	fileBytes, err := os.ReadFile(name)
	if err != nil {
		log.error(
			fmt.Sprintf("failed to read file: %s | error: %s", name, err.Error()),
			req.handlerName,
		)
		req.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
	}

	_, err = io.Copy(req.Writer, bytes.NewReader(fileBytes))
	if err != nil {
		log.error(
			fmt.Sprintf("failed to copy file: %s | error: %s", name, err.Error()),
			req.handlerName,
		)
		req.abortWithErr(http.StatusInternalServerError, fmt.Errorf(internalServerErr))
	}
}

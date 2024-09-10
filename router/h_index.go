package router

import "net/http"

func index(domain string) HttpHandleFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := indexTmpl.ExecuteTemplate(writer, "index.html", struct {
			WebAppURL string
		}{
			WebAppURL: domain,
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

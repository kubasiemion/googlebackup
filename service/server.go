package service

import (
	"net/http"
	"text/template"

	"golang.org/x/oauth2"
)

func init() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		r.FormValue("code")
		if r.FormValue("code") != "" {
			tokenreceived <- r.FormValue("code")
			template2.Execute(w, nil)
		} else {
			template3.Execute(w, nil)
		}

	})
	http.HandleFunc("/prompt", promptHtml)
	http.HandleFunc("/selfclose", selfcloseHtml)
}

var srv = http.Server{
	Addr: ":18080",
}

var template1, _ = template.New("prompt").Parse(promottemplate)
var template2, _ = template.New("selfclose").Parse(youcanclosenow)
var template3, _ = template.New("authmissing").Parse(authmissing)

func promptHtml(w http.ResponseWriter, r *http.Request) {
	template1.Execute(w, config.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
}

func selfcloseHtml(w http.ResponseWriter, r *http.Request) {
	template2.Execute(w, nil)
}

const promottemplate = `<!doctype html><html><head><title>Authorization</title></head><body>
<h1>Authorization</h1><p>Please close this tab after you have authorized the app.</p>
<a href="{{.}}">Authorize</a>
</body></html>`

const youcanclosenow = `<!doctype html><html><head><title>Authorization</title></head><body>
<h1>Authorization received</h1><p>You can close this tab now.</p>
</body></html>`

const authmissing = `<!doctype html><html><head><title>Authorization</title></head><body>
<h1>Authorization missing</h1><p>Authorization code is missing.</p>		
</body></html>`

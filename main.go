package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Mock struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	StatusCode  string `json:"statusCode"`
	ContentType string `json:"contentType"`
	Response    string `json:"response"`
}

type MockHandler struct {
	mutex sync.Mutex
	mock  Mock
	f     string
	p     string
}

func (h MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if r.Method == h.mock.Method {
		w.Header().Set("Content-Type", h.mock.ContentType)
		var statusCode int
		statusCode, _ = strconv.Atoi(h.mock.StatusCode)
		w.WriteHeader(statusCode)

		w.Write([]byte(h.mock.Response))
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("JSON Server")

	portTextEntry := widget.NewEntry()
	portTextEntry.SetText("8080")

	var mocks []Mock

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Port number", Widget: portTextEntry},
		},
		OnSubmit: func() {
			p := portTextEntry.Text
			f := "./mock.json"

			mockJsonFile, err := os.Open(f)
			if err != nil {
				dialog.ShowError(err, w)
			}

			defer mockJsonFile.Close()

			mockByteValue, _ := ioutil.ReadAll(mockJsonFile)

			json.Unmarshal(mockByteValue, &mocks)

			for _, mock := range mocks {
				var handler MockHandler
				handler.mock = mock
				handler.f = f
				handler.p = p

				http.Handle(mock.Path, handler)
			}

			http.ListenAndServe(":"+p, nil)
		},
	}
	form.SubmitText = "server start"

	w.SetContent(container.NewVBox(
		form,
	))

	w.Resize(fyne.NewSize(200, 50))
	w.ShowAndRun()
}

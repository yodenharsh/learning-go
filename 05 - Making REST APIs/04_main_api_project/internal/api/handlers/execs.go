package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Placeholder for execs route"))

	if r.Method == http.MethodPost {
		w.Write([]byte("Hello POST Method on Execs Route"))
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wronge when trying to access form values", http.StatusInternalServerError)
			return
		}
		formMap := r.Form
		fmt.Println("formMap:", formMap)

		fmt.Println("Queries: ", r.URL.Query())
	}
}

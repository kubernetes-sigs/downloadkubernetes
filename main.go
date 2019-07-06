package main

import (
	"log"
	"net/http"
	"os/exec"
)

func main() {
	// assume nginx has authed the request
	http.HandleFunc("/update", updateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// run the script
	cmd := exec.Command("updateindex")
	out, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.Write(out)
		return
	}
	// run npm update

	cmd = exec.Command("npm", "run", "build-prod")
	out, err = cmd.CombinedOutput()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.Write(out)
		return
	}
	w.WriteHeader(http.StatusOK)
}

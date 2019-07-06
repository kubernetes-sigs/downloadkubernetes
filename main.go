package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	// assume nginx has authed the request
	http.HandleFunc("/update", updateHandler)
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating template")
	// run the script
	cmd := exec.Command("./scripts/updateindex")
	out, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.Write(out)
		return
	}
	// run npm update

	fmt.Println("regenerating index")
	cmd = exec.Command("npm", "run", "build-prod")
	out, err = cmd.CombinedOutput()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.Write(out)
		return
	}
	w.WriteHeader(http.StatusOK)
}

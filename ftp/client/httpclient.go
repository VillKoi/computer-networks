package client

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jlaffaye/ftp"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Link     string `json:"link"`
}

func (c *FTPClient) Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origqin", "*")

	var err error
	var credentials Credentials

	if err = json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Println("can't unmarshal body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.FTP, err = ftp.Dial(credentials.Link, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = c.FTP.Login(credentials.Login, credentials.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FTPClient) Ls(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	path := query.Get("path")

	err := c.List(path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

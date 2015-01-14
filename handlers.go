package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"text/template"
)

type Page struct {
	Profile    string
	Passphrase string
	Sites      []Site
}

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	passphrase := r.FormValue("p")
	host := r.FormValue("host")
	minimumLength, _ := strconv.Atoi(r.FormValue("minimumLength"))
	maximumLength, _ := strconv.Atoi(r.FormValue("maximumLength"))
	minimumDigits, _ := strconv.Atoi(r.FormValue("minimumDigits"))
	minimumUppercase, _ := strconv.Atoi(r.FormValue("minimumUppercase"))
	minimumSpecialCharacters, _ := strconv.Atoi(r.FormValue("minimumSpecialCharacters"))
	specialCharacters := r.FormValue("specialCharacters")

	if profile == "" || passphrase == "" || host == "" {
		http.Error(w, "Missing credentials", http.StatusUnauthorized)
		return
	}

	url, err := url.Parse(host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	site := Site{
		Host:                      url.Host,
		MinimumLength:             minimumLength,
		MaximumLength:             maximumLength,
		SpecialCharacters:         specialCharacters,
		NumberOfSpecialCharacters: minimumSpecialCharacters,
		NumberOfDigits:            minimumDigits,
		NumberOfUpperCase:         minimumUppercase,
		Revision:                  0,
	}

	book := getBookname(profile)
	sites, err := read(book, passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sites = append(sites, site)
	err = save(book, passphrase, sites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/book", http.StatusSeeOther)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	passphrase := r.FormValue("p")
	newPassphrase := r.FormValue("newPassphrase")
	confirmPassphrase := r.FormValue("confirmPassphrase")
	cmd := r.FormValue("cmd")

	if profile == "" || passphrase == "" || newPassphrase == "" || confirmPassphrase == "" || cmd == "" {
		http.Error(w, "Missing credentials", http.StatusUnauthorized)
		return
	}

	book := getBookname(profile)

	if cmd == "delete" {
		err := os.Remove(book)
		if err != nil {
			// Return an error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if cmd == "update" {
		if newPassphrase != confirmPassphrase {
		}
		sites, err := read(book, passphrase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = save(book, newPassphrase, sites)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/book", http.StatusSeeOther)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	passphrase := r.FormValue("p")
	host := r.FormValue("host")

	if profile == "" || passphrase == "" || host == "" {
		http.Error(w, "Missing credentials", http.StatusUnauthorized)
		return
	}

	// Update the revision number and generate a new password
	book := getBookname(profile)
	sites, err := read(book, passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, s := range sites {
		if s.Host == host {
			sites[i].Revision++
			break
		}
	}
	err = save(book, passphrase, sites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/book", http.StatusSeeOther)
}

func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	passphrase := r.FormValue("p")
	host := r.FormValue("host")

	if profile == "" || passphrase == "" || host == "host" {
		http.Error(w, "Missing credentials", http.StatusUnauthorized)
		return
	}

	// Remove the site from our book and save it
	book := getBookname(profile)
	sites, err := read(book, passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, site := range sites {
		if site.Host == host {
			sites = append(sites[:i], sites[i+1:]...)
			break
		}
	}
	err = save(book, passphrase, sites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/book", http.StatusSeeOther)
}

func SignOffHandler(w http.ResponseWriter, r *http.Request) {
	cookieProfile := &http.Cookie{
		Name:   "profile",
		MaxAge: -1,
	}
	cookiePassphrase := &http.Cookie{
		Name:   "passphrase",
		MaxAge: -1,
	}

	http.SetCookie(w, cookieProfile)
	http.SetCookie(w, cookiePassphrase)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	passphrase := r.FormValue("p")

	if profile == "" || passphrase == "" {
		c, err := r.Cookie("profile")
		if err == nil {
			profile = c.Value
		}

		c, err = r.Cookie("passphrase")
		if err == nil {
			passphrase = c.Value
		}
	}

	if profile == "" || passphrase == "" {
		http.Redirect(w, r, "/book", http.StatusSeeOther)
		return
	}

	// Set cookies
	//expire := time.Now().AddDate(0, 0, 1)
	cookieProfile := &http.Cookie{
		Name:  "profile",
		Value: profile,
		//Path:       "/",
		//Domain:     "localhost",
		//Expires:    expire,
		//RawExpires: expire.Format(time.UnixDate),
		//MaxAge:     0,
		//Secure:     false,
		//HttpOnly:   true,
	}
	http.SetCookie(w, cookieProfile)
	cookiePassphrase := &http.Cookie{
		Name:  "passphrase",
		Value: passphrase,
		//Path:       "/",
		//Domain:     "localhost",
		//Expires:    expire,
		//RawExpires: expire.Format(time.UnixDate),
		//MaxAge:     0,
		//Secure:     false,
		//HttpOnly:   true,
	}
	http.SetCookie(w, cookiePassphrase)
	book := getBookname(profile)

	sites, err := read(book, passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, s := range sites {
		p := s.generatePassphrase(profile, passphrase)
		sites[i].Password = fmt.Sprintf("%s", string(p))
	}

	page := Page{
		Profile:    profile,
		Passphrase: passphrase,
		Sites:      sites,
	}

	t := template.Must(template.New("book").Parse(templateBook))
	err = t.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, templateIndex)
}

package main

import (
	"html/template"
	"net/http"
	"regexp"
	"sync"
)

var vues int
var mu sync.Mutex 

type Etudiant struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string 
}

type Classe struct {
	Nom         string
	Filiere     string
	Niveau      string
	NbEtudiants int
	Etudiants   []Etudiant
}

type User struct {
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
}

func promoHandler(w http.ResponseWriter, r *http.Request) {
	etudiants := []Etudiant{
		{Nom: "Dupont", Prenom: "Jean", Age: 20, Sexe: "masculin"},
		{Nom: "Durand", Prenom: "Marie", Age: 21, Sexe: "feminin"},
	}
	classe := Classe{
		Nom:         "B1 Informatique",
		Filiere:     "Informatique",
		Niveau:      "Bachelor 1",
		NbEtudiants: len(etudiants),
		Etudiants:   etudiants,
	}

	tmpl := template.Must(template.ParseFiles("templates/promo.html"))
	tmpl.Execute(w, classe)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	vues++
	mu.Unlock()

	tmpl := template.Must(template.ParseFiles("templates/change.html"))
	data := struct {
		Vues   int
		IsPair bool
	}{
		Vues:   vues,
		IsPair: vues%2 == 0,
	}

	tmpl.Execute(w, data)
}

func userFormHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/user_form.html"))
    tmpl.Execute(w, nil)
}

func userTreatmentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/user/form", http.StatusSeeOther)
        return
    }

    nom := r.FormValue("nom")
    prenom := r.FormValue("prenom")
    dateNaissance := r.FormValue("date_naissance")
    sexe := r.FormValue("sexe")

    validName := regexp.MustCompile(`^[a-zA-Z]{1,32}$`)
    if !validName.MatchString(nom) || !validName.MatchString(prenom) || (sexe != "masculin" && sexe != "feminin" && sexe != "autre") {
        http.Redirect(w, r, "/error", http.StatusSeeOther)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:  "nom",
        Value: nom,
        Path:  "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:  "prenom",
        Value: prenom,
        Path:  "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:  "date_naissance",
        Value: dateNaissance,
        Path:  "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:  "sexe",
        Value: sexe,
        Path:  "/",
    })

    http.Redirect(w, r, "/user/display", http.StatusSeeOther)
}

func userDisplayHandler(w http.ResponseWriter, r *http.Request) {
    nomCookie, errNom := r.Cookie("nom")
    prenomCookie, errPrenom := r.Cookie("prenom")
    dateNaissanceCookie, errDateNaissance := r.Cookie("date_naissance")
    sexeCookie, errSexe := r.Cookie("sexe")

    if errNom != nil || errPrenom != nil || errDateNaissance != nil || errSexe != nil {
        http.Redirect(w, r, "/user/form", http.StatusSeeOther)
        return
    }

    user := User{
        Nom:           nomCookie.Value,
        Prenom:        prenomCookie.Value,
        DateNaissance: dateNaissanceCookie.Value,
        Sexe:          sexeCookie.Value,
    }

    tmpl := template.Must(template.ParseFiles("templates/user_display.html"))
    tmpl.Execute(w, user)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/error.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/promo", promoHandler)
	http.HandleFunc("/change", changeHandler)
	http.HandleFunc("/user/form", userFormHandler)
	http.HandleFunc("/user/treatment", userTreatmentHandler)
	http.HandleFunc("/user/display", userDisplayHandler)
	http.HandleFunc("/error", errorHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}

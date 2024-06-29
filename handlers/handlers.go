package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func FindFileContents(fileName string) []byte {
	dat, _ := os.ReadFile("static/" + fileName)
	return dat
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("myTemplate").Parse(string(FindFileContents("index.html"))))
	tmpl.Execute(w, "some data")
}

var clientId = "291507503968-gcjvvjat0om9r44kapj2tvuk8ishj9ar.apps.googleusercontent.com"
var err1 = godotenv.Load()
var clientSecret = os.Getenv("GOOGLE_OAUTH_SECRET")
var provider, err = oidc.NewProvider(context.Background(), "https://accounts.google.com")
var verifier = provider.Verifier(&oidc.Config{ClientID: clientId})
var conf = oauth2.Config{
	ClientID:     clientId,
	ClientSecret: clientSecret,
	RedirectURL:  "http://localhost:3000/callback",
	Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	Endpoint:     provider.Endpoint(),
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL("some-user-state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func OauthRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	t, err := conf.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rawIDToken, ok := t.Extra("id_token").(string)
	if !ok {
		// handle missing token
		fmt.Println("token missing!")
	}
	_, err = verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: rawIDToken})
	client := conf.Client(context.Background(), t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var v any

	// Reading the JSON body using JSON decoder
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// printing the json data into std out.
	fmt.Printf("%v", v)
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)

}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("myAdminTemplate").Parse(string(FindFileContents("admin.html"))))
	tmpl.Execute(w, "")
}

func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		idToken, err := verifier.Verify(context.Background(), token.Value)
		if idToken == nil || err != nil {
			// handle error
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

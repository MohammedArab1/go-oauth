package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"

	"go-oauth/handlers"
)

func main() {
	// var tokenAuth *jwtauth.JWTAuth
	godotenv.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handlers.MainPage)
	r.Get("/login", handlers.LoginPage)
	r.Get("/callback", handlers.OauthRedirect)
	r.Group(func(r chi.Router) {
		// tokenAuth = jwtauth.New("RS256", []byte("secret"), nil)
		// tokenAuth.Encode(map[string]interface{}{"user_id": 123})
		// r.Use(jwtauth.Verifier(tokenAuth))
		// r.Use(jwtauth.Authenticator)
		r.Use(handlers.JWTAuthenticator)
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			fmt.Println("CLAIMS IS: ", claims)
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["name"])))
		})
	})
	// r.Get("/admin", handlers.AdminPage)
	http.ListenAndServe(":3000", r)
}

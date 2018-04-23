package muxalice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func init() {
	fmt.Println("Starting VitrnX 0.2 - MUX + Alice backend")
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	// Define our common wrappers with alice
	ws := alice.New(loggingHandler, recoverHandler, checkCredentials, addCrossDomainHeader)

	// Define known routes
	myRouter.Handle("/hello", ws.ThenFunc(sayHello))
	// myRouter.Handle("/posts", ws.ThenFunc(model.GetArticleByTag))
	// myRouter.Handle("/posts/{id}", ws.ThenFunc(model.GetSingleArticle))
	// myRouter.Handle("/by-tag/{tag}", ws.ThenFunc(model.GetArticleByTag))

	// Launch the server
	log.Fatal(http.ListenAndServe(":7777", myRouter))
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

// loggingHandler simply logs every request to stdout
func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

// recoverHandler adds a wrapper to prevent the App to crash in case of a panic
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("got a panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// checkCredentials calls the authentication API to verify the token
func checkCredentials(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO implement
		// jwt := ""
		// gateway.CheckCredentialAgainstFireBase(r.Context(), jwt)
		h.ServeHTTP(w, r) // call original
	}
	return http.HandlerFunc(fn)
}

// addCrossDomainHeader adds custom headers to workaround cross site limitations
func addCrossDomainHeader(h http.Handler) http.Handler {
	// FIXME what are the consequences in term of security
	// see https://stackoverflow.com/questions/22363268/cross-origin-request-blocked

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "True")

		// Cross domain headers
		if acrh, ok := r.Header["Access-Control-Request-Headers"]; ok {
			w.Header().Set("Access-Control-Allow-Headers", acrh[0])
			fmt.Println("Set ac allow headers for: " + acrh[0])
		}

		if acao, ok := r.Header["Access-Control-Allow-Origin"]; ok {
			w.Header().Set("Access-Control-Allow-Origin", acao[0])
			fmt.Println("Set ac allow origin for: " + acao[0])
		} else {
			if _, oko := r.Header["Origin"]; oko {
				w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
				fmt.Println("Set ac allow origin via origin header for: " + r.Header["Origin"][0])
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				fmt.Println("Fall back,  ac allow origin for: *")
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Connection", "Close")
	}
	return http.HandlerFunc(fn)
}

package core

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Router *mux.Router
	Store  *Store
}

/* Authentication related controllers */
const (
	sessionId = "sessionId"
)

var session_store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)
	session_store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	session_store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func Session(r *http.Request) (*sessions.Session, error) {
	return session_store.Get(r, sessionId)
}

func SessionUser(r *http.Request) string {
	session, _ := Session(r)
	username := session.Values["username"].(string)

	return username
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	jsonResponse(w, code, map[string]string{"result": "ERROR", "message": message})
}

func successResponse(w http.ResponseWriter, code int, body interface{}) {
	jsonResponse(w, code, map[string]interface{}{"result": "SUCCESS", "body": body})
}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	io.WriteString(w, "\n")
}

func NewServer(store *Store) *Server {
	_router := mux.NewRouter()
	server := &Server{
		Router: _router,
		Store:  store,
	}

	_router.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/static/app.js")
	})

	_router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/static/style.css")
	})

	_api_auth_subrouter := _router.PathPrefix("/api").Subrouter()
	_api_auth_subrouter.Use(AuthMiddleware)
	_api_auth_subrouter.HandleFunc("/list", server.list)
	_api_auth_subrouter.HandleFunc("/order", server.order)

	_router.HandleFunc("/login", server.LoginPage)
	_router.HandleFunc("/logout", server.logout).Methods("GET")

	_static_auth_subrouter := _router.PathPrefix("/").Subrouter()
	_static_auth_subrouter.Use(AuthMiddleware)
	_static_auth_subrouter.HandleFunc("/", server.ProductsPage).Methods("GET")
	_static_auth_subrouter.HandleFunc("/products", server.ProductsPage).Methods("GET")
	_static_auth_subrouter.HandleFunc("/checkout", server.CheckoutPage).Methods("GET")
	_static_auth_subrouter.HandleFunc("/thank-you", server.ThankYou).Methods("GET")

	return server
}

// AuthMiddleware is a middleware function to check authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request contains an "Authorization" header
		setupCorsResponse(&w, r)
		_s, _ := Session(r)
		_is_auth := _s.Values["authenticated"]
		_authenticated, ok := _is_auth.(bool)
		if !ok || !_authenticated {
			log.Infof("User is not authenticated.")
			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (server *Server) list(w http.ResponseWriter, r *http.Request) {
	_q := r.URL.Query()
	limit := _q.Get("limit")
	offset := _q.Get("offset")
	filter := _q.Get("filter")

	if limit == "" {
		limit = "100"
	}

	if offset == "" {
		offset = "0"
	}

	products, err := server.Store.ListProductsStore(limit, offset, filter)
	if err != nil {
		log.Error(err)
		errorResponse(w, http.StatusInternalServerError, "Server error")
		return
	}

	total := server.Store.TotalProductsStore()
	_response := ProductsResponse{
		Products: products,
		Total:    total,
	}
	successResponse(w, http.StatusOK, _response)
}

type Order struct {
	Products []Product `json:"products"`
	Csrf     string    `json:"csrf"`
	Address  struct {
		Area         string `json:"area"`
		City         string `json:"city"`
		Code         int    `json:"code"`
		Street       string `json:"street"`
		StreetNumber int    `json:"streetNumber"`
	} `json:"address"`
	// Add more fields as needed
}

func (server *Server) order(w http.ResponseWriter, r *http.Request) {
	session, _ := Session(r)

	var order_data Order
	err := json.NewDecoder(r.Body).Decode(&order_data)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Infof("incoming csrf: %s, Existing csrf: %#v", order_data.Csrf, session.Values)
	if session.Values["csrf"] != order_data.Csrf {
		session.Values["csrf"] = nil
		errorResponse(w, http.StatusBadRequest, "Bad CSRF token")
		return
	}
	session.Values["csrf"] = nil
	session.Save(r, w)
	log.Infof("Order placed")
	http.Redirect(w, r, "/thank-you", http.StatusMovedPermanently)
}

func (server *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Session(r)
	if session.Values["authenticated"] == true {
		session.Values["authenticated"] = false
		session.Values["username"] = nil
		session.Values["csrf"] = nil
		session.Options.MaxAge = -1
		session.Save(r, w)
		successResponse(w, http.StatusOK, "Logged out")
		return
	} else {
		successResponse(w, http.StatusOK, "Already logged out")
		return
	}
	errorResponse(w, http.StatusInternalServerError, "Server error")
}

// Pages -----------------------------------------------------------------------------------------

func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		setupCorsResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Infof("Could not parse form")
		}

		_u := r.PostFormValue("username")
		_p := r.PostFormValue("password")
		_logged_in := server.Store.Login(_u, _p)

		if _logged_in {
			session, _ := Session(r)
			session.Values["authenticated"] = true
			session.Values["username"] = _u
			session.Values["csrf"] = ""
			// saves all sessions used during the current request
			session.Save(r, w)
			log.Infof("User logged in !!!")
			http.Redirect(w, r, "/products", http.StatusSeeOther)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Wrong credentials"))
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
	}

	// If not a POST request, serve the login page template.
	_data := struct {
		Title string
	}{
		Title: "Login",
	}
	tmpl, err := template.ParseFiles("public/login.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

func (server *Server) ProductsPage(w http.ResponseWriter, r *http.Request) {
	session, _ := Session(r)
	session.Values["csrf"] = nil
	_data := struct {
		Title string
	}{
		Title: "Products",
	}
	session.Save(r, w)
	tmpl, err := template.ParseFiles("public/products.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

func (server *Server) CheckoutPage(w http.ResponseWriter, r *http.Request) {
	session, _ := Session(r)
	session.Values["csrf"] = uuid.New().String()
	log.Infof("Session current %#v", session.Values)

	_data := struct {
		Title string
		Csrf  string
	}{
		Title: "Checkout",
		Csrf:  session.Values["csrf"].(string),
	}
	session.Save(r, w)

	tmpl, err := template.ParseFiles("public/checkout.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

func (server *Server) ThankYou(w http.ResponseWriter, r *http.Request) {
	session, _ := Session(r)
	session.Values["csrf"] = nil
	_data := struct {
		Title string
	}{
		Title: "Thank you",
	}

	session.Save(r, w)
	tmpl, err := template.ParseFiles("public/thankyou.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

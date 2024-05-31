package core

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const CHECKOUT_FORM = "checkout_form"
const COOKIE_TTL = 300

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Router *mux.Router
	Store  *Store
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
}

const (
	sessionId = "sessionId"
)

// Key-value structure to keep track of the valid csrf tokens
var csrf_to_forms = make(map[string]map[string]string)

var invalid_cookies = make(map[string]int64)

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
		HttpOnly: true,
	}
}

func Session(r *http.Request) (*sessions.Session, error) {
	return session_store.Get(r, sessionId)
}

func NewSession(r *http.Request) (*sessions.Session, error) {
	return session_store.New(r, sessionId)
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
		http.ServeFile(w, r, "./client/static/app.js")
	})

	_router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/static/style.css")
	})

	_api_auth_subrouter := _router.PathPrefix("/api").Subrouter()
	_api_auth_subrouter.Use(AuthMiddleware)
	_api_auth_subrouter.HandleFunc("/list", server.list)
	_api_auth_subrouter.HandleFunc("/order", server.order)

	_router.HandleFunc("/login", server.LoginPage)
	_router.HandleFunc("/logout", server.logout).Methods("GET")

	// Pages. All pages are protected by the AuthMiddleware
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
		setupCorsResponse(&w, r)
		_s, _ := Session(r)
		now := time.Now()
		if _s.Values["time_generated"] != nil && _s.Values["id"] != nil {
			for k, v := range invalid_cookies {
				if now.Unix()-v > COOKIE_TTL {
					delete(invalid_cookies, k)
				}
			}

			if _t, ok := invalid_cookies[_s.Values["id"].(string)]; ok {
				if now.Unix()-_t > COOKIE_TTL {
					delete(invalid_cookies, _s.ID)
				} else {
					logger.Infof("User is not authenticated.")
					http.Redirect(w, r, "/login", http.StatusMovedPermanently)
				}
			}

			if now.Unix()-_s.Values["time_generated"].(int64) > COOKIE_TTL {
				logger.Infof("User is not authenticated.")
				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			}
			next.ServeHTTP(w, r)
			return
		}
		logger.Infof("User is not authenticated.")
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
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
		logger.Error(err)
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

func (server *Server) order(w http.ResponseWriter, r *http.Request) {
	_s, err := Session(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var order_data Order
	err = json.NewDecoder(r.Body).Decode(&order_data)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, ok := csrf_to_forms[_s.Values["user"].(string)]; ok {
		if _, ok := csrf_to_forms[_s.Values["user"].(string)][CHECKOUT_FORM]; ok && _s.Values["csrf"] != csrf_to_forms[_s.Values["user"].(string)][CHECKOUT_FORM] {
			errorResponse(w, http.StatusBadRequest, "Bad CSRF token")
			return
		}
	}
	logger.Infof("Incoming csrf: %s, Existing csrf: %#v", order_data.Csrf, _s.Values)
	if _s.Values["csrf"] != order_data.Csrf {
		csrf_to_forms[_s.Values["user"].(string)][CHECKOUT_FORM] = ""
		_s.Values["csrf"] = nil
		errorResponse(w, http.StatusBadRequest, "Bad CSRF token")
		return
	}

	csrf_to_forms[_s.Values["user"].(string)][CHECKOUT_FORM] = ""
	_s.Values["csrf"] = nil

	err = _s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SendMail(
		_s.Values["user"].(string),
		"Thank you for your order!",
		order_data)

	logger.Infof("Order placed")
	http.Redirect(w, r, "/thank-you", http.StatusMovedPermanently)
}

func (server *Server) logout(w http.ResponseWriter, r *http.Request) {
	_s, err := Session(r)
	now := time.Now()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	invalid_cookies[_s.Values["id"].(string)] = now.Unix()
	_s.Options.MaxAge = -1

	err = _s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// Pages -----------------------------------------------------------------------------------------

func (server *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	_authentication_attempt := false
	if r.Method == http.MethodPost {
		setupCorsResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		err := r.ParseForm()
		if err != nil {
			logger.Infof("Could not parse form")
		}

		_u := r.PostFormValue("username")
		_p := r.PostFormValue("password")
		_logged_in := server.Store.Login(_u, _p)

		if _logged_in {
			session, _ := NewSession(r)
			now := time.Now()
			uuid := uuid.New()

			session.Values["user"] = _u
			session.Values["csrf"] = ""
			session.Values["time_generated"] = now.Unix()
			session.Values["id"] = uuid.String()
			session.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   60 * 5, // 5 minutes cookie
				HttpOnly: true,
			}

			// saves all sessions used during the current request
			session.Save(r, w)
			logger.Infof("User logged in.")
			http.Redirect(w, r, "/products", http.StatusSeeOther)
			return
		} else {
			_authentication_attempt = true
		}
	}

	// If not a POST request, serve the login page template.
	_data := struct {
		Title                  string
		Authentication_attempt bool
	}{
		Title:                  "Login",
		Authentication_attempt: _authentication_attempt,
	}
	tmpl, err := template.ParseFiles("client/login.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

func (server *Server) ProductsPage(w http.ResponseWriter, r *http.Request) {
	_data := struct {
		Title string
	}{
		Title: "Products",
	}
	tmpl, err := template.ParseFiles("client/products.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

func (server *Server) CheckoutPage(w http.ResponseWriter, r *http.Request) {
	_s, err := Session(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	new_csrf_token := uuid.New().String()
	_s.Values["csrf"] = new_csrf_token
	csrf_to_forms[_s.Values["user"].(string)] = map[string]string{
		CHECKOUT_FORM: new_csrf_token,
	}

	logger.Infof("Session current %#v", _s.Values)

	_data := struct {
		Title string
		Csrf  string
	}{
		Title: "Checkout",
		Csrf:  _s.Values["csrf"].(string),
	}
	err = _s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("client/checkout.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

func (server *Server) ThankYou(w http.ResponseWriter, r *http.Request) {
	_s, err := Session(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_s.Values["csrf"] = nil
	_data := struct {
		Title string
	}{
		Title: "Thank you",
	}

	err = _s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("client/thankyou.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, _data)
}

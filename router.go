package main

import (
	"net/http"

	"gorestboilerplate/types"
	"gorestboilerplate/utils/api"

	"github.com/gorilla/mux"
)

// RouterContext represents the current HTTP context used by the middleware
type RouterContext struct {
	Source string
	User   string
}

type route struct {
	route   string
	handler func(http.ResponseWriter, *http.Request)
	method  string
}

func setupRoutes(port string) {

	router := mux.NewRouter()

	/* Health check routes */
	healthRoutes := router.PathPrefix("/health").Subrouter()
	healthRoutes.Methods("OPTIONS")

	logger.Info("Adding handler for '/health/ready'")
	healthRoutes.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		api.SuccessResponse(w, "ready", http.StatusOK)
	}).Methods("GET")

	logger.Info("Adding handler for '/health/alive'")
	healthRoutes.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {

		err := db.Ping()
		if err != nil {
			logger.Fatalf("DB Error: %s", err.Error())
		}
		api.SuccessResponse(w, "All OK", http.StatusOK)
	}).Methods("GET")

	/* User Routes */
	userRoutes := router.PathPrefix("/users").Subrouter()
	userRoutes.Methods("OPTIONS")

	routes := []route{
		{
			route:  "",
			method: "GET",
			handler: func(w http.ResponseWriter, r *http.Request) {
				api.SuccessResponse(w, types.User{"John", "Smith"}, http.StatusOK)
			},
		},
		{
			route:  "",
			method: "POST",
			handler: func(w http.ResponseWriter, r *http.Request) {
				api.SuccessResponse(w, types.User{"John", "Smith"}, http.StatusOK)
			},
		},
		{
			route:  "/{id}",
			method: "PUT",
			handler: func(w http.ResponseWriter, r *http.Request) {
				api.SuccessResponse(w, types.User{"John", "Smith"}, http.StatusOK)
			},
		},
		{
			route:  "/{id}",
			method: "GET",
			handler: func(w http.ResponseWriter, r *http.Request) {
				api.SuccessResponse(w, types.User{"John", "Smith"}, http.StatusOK)
			},
		},
		{
			route:  "/{id}",
			method: "DELETE",
			handler: func(w http.ResponseWriter, r *http.Request) {
				api.SuccessResponse(w, nil, http.StatusOK)
			},
		},
	}

	for _, p := range routes {

		logger.Infof("Adding handler for '%s /users%s'", p.method, p.route)
		userRoutes.HandleFunc(p.route, p.handler).Methods(p.method)
	}

	/* misc routes */
	router.Methods("OPTIONS")

	routes = []route{
		{
			route:  "/register",
			method: "POST",
			handler: func(w http.ResponseWriter, r *http.Request) {
				api.SuccessResponse(w, "ready", http.StatusOK)
			},
		},
	}

	for _, p := range routes {

		logger.Infof("Adding handler for '%s %s'", p.method, p.route)
		router.HandleFunc(p.route, p.handler).Methods(p.method)
	}

	logger.Infof("Startup Complete")
	logger.Infof("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	logger.Fatal(http.ListenAndServe(":"+port, router))
}

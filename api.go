package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAdrr string
	storage    Storage
}

// func that return an ApiServer instance
func NewApiServer(listenAdrr string, storage Storage) *ApiServer {
	return &ApiServer{
		listenAdrr: listenAdrr,
		storage:    storage,
	}
}

// my api entrypoint
func (s *ApiServer) Run() {
	router := mux.NewRouter()
	//creating routes
	router.HandleFunc("/products", makeHttpHandleFunc(s.HandleGetProducts))
	router.HandleFunc("/product/{id}", makeHttpHandleFunc(s.HandleGetProductById))
	router.HandleFunc("/product", makeHttpHandleFunc(s.HandleCreateProduct))
	log.Println("JSON API Server running on port: ", s.listenAdrr)
	http.ListenAndServe(s.listenAdrr, router)
}

//usecases of my api

func (s *ApiServer) HandleProduct(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.HandleGetProducts(w, r)
	}

	if r.Method == "POST" {
		return s.HandleCreateProduct(w, r)
	}

	if r.Method == "DELETE" {
		return s.HandleDeleteProduct(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiServer) HandleGetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := s.storage.GetProducts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

func (s *ApiServer) HandleGetProductById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]

	product, err := s.storage.GetProductById(vars)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, product)
}

func (s *ApiServer) HandleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	// getting the request and fill the product object
	createProductReq := CreateRequestProduct{}
	if err := json.NewDecoder(r.Body).Decode(&createProductReq); err != nil {
		return err
	}
	product := NewProduct(createProductReq.Name, createProductReq.Description, createProductReq.Price)
	res, err := s.storage.CreateProduct(product)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, res)
}

func (s *ApiServer) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// to serve response in JSON format
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}

}

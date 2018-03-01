package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davidsbond/lux"
	mgo "gopkg.in/mgo.v2"
)

type (
	// The CustomerService provides CRUD operations on customers.
	CustomerService struct {
		session  *mgo.Session
		database *mgo.Database
	}

	// Customer represents a customer in the database.
	Customer struct {
		ID        string    `json:"id"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		DOB       time.Time `json:"dob"`
	}
)

// New creates a new instance of the service.
func New() (*CustomerService, error) {
	session, err := mgo.Dial("")

	if err != nil {
		return nil, err
	}

	db := session.DB("customers")

	return &CustomerService{
		session:  session,
		database: db,
	}, nil
}

// Start registeres the lambda handlers and starts the lambda function.
func (s *CustomerService) Start() {
	router := lux.NewRouter()

	router.Handler("POST", s.Insert).Headers("Content-Type", "application/json")
	router.Handler("DELETE", s.Delete)
	router.Handler("GET", s.Get)
	router.Handler("PUT", s.Update).Headers("Content-Type", "application/json")

	lambda.Start(router.HandleRequest)
}

// Get obtains a customer from the datbase.
func (s *CustomerService) Get(req lux.Request) lux.Response {
	var resp lux.Response
	var customer Customer

	id, ok := req.QueryStringParameters["id"]

	if !ok {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = "no id provided"

		return resp
	}

	customers := s.database.C("customers")

	if err := customers.FindId(id).One(&customer); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	json, err := json.Marshal(customer)

	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	resp.StatusCode = http.StatusOK
	resp.Body = string(json)

	return resp
}

// Delete removes a customer from the database.
func (s *CustomerService) Delete(req lux.Request) lux.Response {
	var resp lux.Response

	id, ok := req.QueryStringParameters["id"]

	if !ok {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = "no id provided"

		return resp
	}

	customers := s.database.C("customers")

	if err := customers.RemoveId(id); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	resp.StatusCode = http.StatusOK

	return resp
}

// Insert adds a customer to the database.
func (s *CustomerService) Insert(req lux.Request) lux.Response {
	var customer Customer
	var resp lux.Response

	if err := json.Unmarshal([]byte(req.Body), &customer); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	customers := s.database.C("customers")

	if err := customers.Insert(customer); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	resp.StatusCode = http.StatusOK

	return resp
}

// Update edits a customer in the database.
func (s *CustomerService) Update(req lux.Request) lux.Response {
	var customer Customer
	var resp lux.Response

	if err := json.Unmarshal([]byte(req.Body), &customer); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	customers := s.database.C("customers")

	if err := customers.UpdateId(customer.ID, customer); err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()

		return resp
	}

	resp.StatusCode = http.StatusOK

	return resp
}

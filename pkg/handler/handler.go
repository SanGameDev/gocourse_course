package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SanGameDev/go_lib_response/response"
	"github.com/SanGameDev/gocourse_course/internal/course"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {
	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse, encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse, encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourses, encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse, encodeResponse,
		opts...,
	)).Methods("PATCH")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse, encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf(("invalid request format: '%s'"), err.Error()))
	}

	return req, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (interface{}, error) {

	p := mux.Vars(r)
	req := course.GetReq{
		ID: p["id"],
	}
	return req, nil
}

func decodeGetAllCourses(_ context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Page:  page,
		Limit: limit,
	}

	return req, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf(("invalid request format: '%s'"), err.Error()))
	}

	p := mux.Vars(r)
	req.ID = p["id"]

	return req, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {

	p := mux.Vars(r)

	req := course.DeleteReq{
		ID: p["id"],
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

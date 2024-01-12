package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bootcamp-go/go-web/internal/products"
	"github.com/bootcamp-go/go-web/internal/repository"
	"github.com/go-chi/chi/v5"
)

func NewRequest(method, url string, body io.Reader, urlParams map[string]string, urlQuery map[string]string) *http.Request {
	// old request
	req := httptest.NewRequest(method, url, body)

	// url params
	// - new request with a new context with key chi.RouteCtxKey and value chiCtx -> "id":"1"
	if urlParams != nil {
		chiCtx := chi.NewRouteContext()
		for k, v := range urlParams {
			chiCtx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	}

	// url query
	// r.URL.RawQuery = "name=task 1"
	if urlQuery != nil {
		query := req.URL.Query()
		for k, v := range urlQuery {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode() // raw string
	}

	return req
}

// r.Get("/ping", h.Ping())
func TestPingPong(t *testing.T) {
	//Given
	var data []products.Product
	rp := repository.NewProductRepository(&data)
	h := NewProductHandler(rp)

	//WHEN
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Crear un ResponseRecorder (implementa http.ResponseWriter) para grabar la respuesta.
	rr := httptest.NewRecorder()

	// Crear un http.HandlerFunc que se pueda pasar al método ServeHTTP.
	// Llamar al manipulador (handler) con la solicitud y el ResponseRecorder.
	hdFunc := h.Ping()
	hdFunc(rr, req)
	//handler := http.HandlerFunc(h.Ping()) //equivalente
	//handler.ServeHTTP(rr, req)

	//THEN
	// Verificar el código de estado esperado.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Código de estado incorrecto: obtuvo %v, esperaba %v",
			status, http.StatusOK)
	}
	t.Logf("Código de estado correcto: obtuvo %v, esperaba %v", rr.Code, http.StatusOK)

	// Verificar el cuerpo de la respuesta.
	expected := "pong"
	if rr.Body.String() != expected {
		t.Errorf("Cuerpo de respuesta incorrecto: obtuvo %v, esperaba %v",
			rr.Body.String(), expected)
	}
	t.Logf("Cuerpo de respuesta correcto: obtuvo %v, esperaba %v", rr.Body.String(), expected)
}

func TestGetById(t *testing.T) {
	//Given
	idObj := "10"
	idObj2 := 10
	var data []products.Product
	data = append(data, products.Product{Id: idObj2})
	rp := repository.NewProductRepository(&data)
	h := NewProductHandler(rp)

	//WHEN
	req := NewRequest("GET", "/products/"+idObj, nil, map[string]string{"productId": idObj}, nil)

	// Crear un ResponseRecorder (implementa http.ResponseWriter) para grabar la respuesta.
	rr := httptest.NewRecorder()

	// Crear un http.HandlerFunc que se pueda pasar al método ServeHTTP.
	// Llamar al manipulador (handler) con la solicitud y el ResponseRecorder.
	hdFunc := h.GetById()
	hdFunc(rr, req)
	//handler := http.HandlerFunc(h.Ping()) //equivalente
	//handler.ServeHTTP(rr, req)

	//THEN
	// Verificar el código de estado esperado.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Código de estado incorrecto: obtuvo %v, esperaba %v",
			status, http.StatusOK)
	} else {
		t.Logf("Código de estado correcto: obtuvo %v, esperaba %v", rr.Code, http.StatusOK)
	}

	// Verificar el cuerpo de la respuesta.
	// expected := `{"message":"OK","data":{"id":10,"name":"Soup Bowl Clear 8oz92008","quantity":424,"code_value":"B180","is_published":false,"expiration":"18/10/2021","price":92.8}}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("Cuerpo de respuesta incorrecto: obtuvo %v, esperaba %v",
	// 		rr.Body.String(), expected)
	// } else {
	// 	t.Logf("Cuerpo de respuesta correcto: obtuvo %v, esperaba %v", rr.Body.String(), expected)
	// }
}
func TestDelete(t *testing.T) {
	//Given

	//Given
	idObj := "1"
	var data []products.Product
	rp := repository.NewProductRepository(&data)
	h := NewProductHandler(rp)
	t.Setenv("GO_API_TOKEN_1", "julio")

	//WHEN
	req := NewRequest("DELETE", "/products/"+idObj, nil, map[string]string{"productId": idObj}, nil)
	req.Header.Add("token", "julio")

	// Crear un ResponseRecorder (implementa http.ResponseWriter) para grabar la respuesta.
	rr := httptest.NewRecorder()

	// Crear un http.HandlerFunc que se pueda pasar al método ServeHTTP.
	// Llamar al manipulador (handler) con la solicitud y el ResponseRecorder.
	hdFunc := h.DeleteById()
	hdFunc(rr, req)
	//handler := http.HandlerFunc(h.Ping()) //equivalente
	//handler.ServeHTTP(rr, req)

	//THEN
	// Verificar el código de estado esperado.
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Código de estado incorrecto: obtuvo %v, esperaba %v", status, http.StatusOK)
		t.Errorf("Response: %v", rr.Body)
	} else {
		t.Logf("Código de estado correcto: obtuvo %v, esperaba %v", rr.Code, http.StatusOK)
	}

}
func TestPost(t *testing.T) {
	//Given
	t.Setenv("GO_API_TOKEN_1", "julio")
	var data []products.Product = make([]products.Product, 1)
	rp := repository.NewProductRepository(&data)
	h := NewProductHandler(rp)
	json := `{
		"name": "Cesar",
		"quantity": 1,
		"code_value": "asdfasdf",
		"is_published": true,
		"expiration": "15/04/2022",
		"price": 1.0
	}`
	body := strings.NewReader(json)

	//WHEN
	req := NewRequest("POST", "/products", body, nil, nil)
	//req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Add("token", "julio")

	// Crear un ResponseRecorder (implementa http.ResponseWriter) para grabar la respuesta.
	rr := httptest.NewRecorder()

	// Crear un http.HandlerFunc que se pueda pasar al método ServeHTTP.
	// Llamar al manipulador (handler) con la solicitud y el ResponseRecorder.
	hdFunc := h.Post()
	hdFunc(rr, req)
	//handler := http.HandlerFunc(h.Ping()) //equivalente
	//handler.ServeHTTP(rr, req)

	//THEN
	// Verificar el código de estado esperado.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Código de estado incorrecto: obtuvo %v, esperaba %v", status, http.StatusOK)
		t.Errorf("Response: %v", rr.Body)
	} else {
		t.Logf("Código de estado correcto: obtuvo %v, esperaba %v", rr.Code, http.StatusOK)
	}

}

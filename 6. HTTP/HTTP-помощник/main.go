package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// начало решения

// Handy предоставляет удобный интерфейс
// для выполнения HTTP-запросов
type Handy struct {
	url     string
	client  *http.Client
	headers map[string]string
	params  *url.Values
	body    []byte
	error   error
}

// NewHandy создает новый экземпляр Handy
func NewHandy() *Handy {
	return &Handy{
		client:  &http.Client{Timeout: 3 * time.Millisecond},
		headers: map[string]string{},
		params:  &url.Values{},
		error:   nil,
		body:    nil,
	}
}

// URL устанавливает URL, на который пойдет запрос
func (h *Handy) URL(uri string) *Handy {
	h.url = uri
	return h
}

// Client устанавливает HTTP-клиента
// вместо умолчательного http.DefaultClient
func (h *Handy) Client(client *http.Client) *Handy {
	h.client = client
	return h
}

// Header устанавливает значение заголовка
func (h *Handy) Header(key, value string) *Handy {
	h.headers[key] = value
	return h
}

// Param устанавливает значение URL-параметра
func (h *Handy) Param(key, value string) *Handy {
	h.params.Set(key, value)
	return h
}

// Form устанавливает данные, которые будут закодированы
// как application/x-www-form-urlencoded и отправлены в теле запроса
// с соответствующим content-type
func (h *Handy) Form(form map[string]string) *Handy {
	formBody := &url.Values{}
	for k, v := range form {
		formBody.Add(k, v)
	}

	h.headers["Content-Type"] = "application/x-www-form-urlencoded"
	h.body = []byte(formBody.Encode())

	return h
}

// JSON устанавливает данные, которые будут закодированы
// как application/json и отправлены в теле запроса
// с соответствующим content-type
func (h *Handy) JSON(v any) *Handy {
	body, err := json.Marshal(v)

	if err != nil {
		h.error = err
		return h
	}

	h.headers["Content-Type"] = "application/json"
	h.body = body

	return h
}

// Get выполняет GET-запрос с настроенными ранее параметрами
func (h *Handy) Get() *HandyResponse {
	if h.error != nil {
		return &HandyResponse{0, nil, h.error}
	}

	request, requestError := http.NewRequest(http.MethodGet, h.url, bytes.NewReader(h.body))
	if requestError != nil {
		return &HandyResponse{0, nil, requestError}
	}

	// get parameters
	request.URL.RawQuery = h.params.Encode()

	// headers
	for k, v := range h.headers {
		request.Header.Add(k, v)
	}

	// make request
	resp, responseErr := h.client.Do(request)
	if responseErr != nil {
		return &HandyResponse{0, nil, responseErr}
	}
	defer resp.Body.Close()

	// read response
	body, readResponseError := io.ReadAll(resp.Body)
	if readResponseError != nil {
		return &HandyResponse{0, nil, readResponseError}
	}

	return &HandyResponse{
		StatusCode:   resp.StatusCode,
		ResponseBody: body,
		error:        nil,
	}
}

// Post выполняет POST-запрос с настроенными ранее параметрами
func (h *Handy) Post() *HandyResponse {
	if h.error != nil {
		return &HandyResponse{0, nil, h.error}
	}

	request, requestError := http.NewRequest(http.MethodPost, h.url, bytes.NewReader(h.body))
	if requestError != nil {
		return &HandyResponse{0, nil, requestError}
	}

	// get parameters
	request.URL.RawQuery = h.params.Encode()

	// headers
	for k, v := range h.headers {
		request.Header.Add(k, v)
	}

	// make request
	resp, responseErr := h.client.Do(request)
	if responseErr != nil {
		return &HandyResponse{0, nil, responseErr}
	}
	defer resp.Body.Close()

	// read response
	body, readResponseError := io.ReadAll(resp.Body)
	if readResponseError != nil {
		return &HandyResponse{0, nil, readResponseError}
	}

	return &HandyResponse{
		StatusCode:   resp.StatusCode,
		ResponseBody: body,
		error:        nil,
	}
}

// HandyResponse представляет ответ на HTTP-запрос
type HandyResponse struct {
	StatusCode   int
	ResponseBody []byte
	error        error
}

// OK возвращает true, если во время выполнения запроса
// не произошло ошибок, а код HTTP-статуса ответа равен 200
func (r *HandyResponse) OK() bool {
	return r.error == nil && r.StatusCode == http.StatusOK
}

// Bytes возвращает тело ответа как срез байт
func (r *HandyResponse) Bytes() []byte {
	return r.ResponseBody
}

// String возвращает тело ответа как строку
func (r *HandyResponse) String() string {
	return string(r.ResponseBody)
}

// JSON декодирует тело ответа из JSON и сохраняет
// результат по адресу, на который указывает v
func (r *HandyResponse) JSON(v any) {
	// работает аналогично json.Unmarshal()
	// если при декодировании произошла ошибка,
	// она должна быть доступна через r.Err()

	err := json.Unmarshal(r.ResponseBody, v)

	if err != nil {
		r.error = err
	}
}

// Err возвращает ошибку, которая возникла при выполнении запроса
// или обработке ответа
func (r *HandyResponse) Err() error {
	return r.error
}

// конец решения

func main() {
	{
		// примеры запросов

		// GET-запрос с параметрами
		NewHandy().URL("https://httpbingo.org/get").Param("id", "42").Get()

		// HTTP-заголовки
		NewHandy().
			URL("https://httpbingo.org/get").
			Header("Accept", "text/html").
			Header("Authorization", "Bearer 1234567890").
			Get()

		// POST формы
		params := map[string]string{
			"brand":    "lg",
			"category": "tv",
		}
		NewHandy().URL("https://httpbingo.org/post").Form(params).Post()

		// POST JSON-документа
		NewHandy().URL("https://httpbingo.org/post").JSON(params).Post()
	}

	{
		// пример обработки ответа

		// отправляем GET-запрос с параметрами
		resp := NewHandy().URL("https://httpbingo.org/get").Param("id", "42").Get()
		if !resp.OK() {
			panic(resp.String())
		}

		// декодируем ответ в JSON
		var data map[string]any
		resp.JSON(&data)

		fmt.Println(data["url"])
		// "https://httpbingo.org/get"
		fmt.Println(data["args"])
		// map[id:[42]]
	}
}

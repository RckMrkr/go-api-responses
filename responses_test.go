package responses

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(f))
}

func callTestServer(ts *httptest.Server) (string, int, error) {
	res, err := http.Get(ts.URL)
	if err != nil {
		return "", -1, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return "", -1, err
	}

	return string(body), res.StatusCode, nil
}

func TestCheckCodes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Function http.HandlerFunc
		Code     int
	}{
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondSuccess(w, nil)
			},
			200,
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondNotFound(w)
			},
			404,
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondUnauthorizedBearerJWT(w)
			},
			401,
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondInternalError(w, "")
			},
			500,
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondBadRequest(w, "")
			},
			400,
		},
	}

	for _, test := range tests {
		func() {
			ts := createTestServer(test.Function)
			defer ts.Close()

			_, code, err := callTestServer(ts)
			if !assert.Nil(err) {
				return
			}

			assert.Equal(test.Code, code)
		}()
	}
}

func TestNoBody(t *testing.T) {
	assert := assert.New(t)
	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		RespondSuccess(w, nil)
	})

	defer ts.Close()

	body, _, err := callTestServer(ts)
	if !assert.Nil(err) {
		return
	}

	assert.Equal("", string(body))
}

// Unsure how to check JSON equivalence
func TestMarshaling(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Function http.HandlerFunc
		Json     string
	}{
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondSuccess(w, nil)
			},
			"",
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondSuccess(w,
					struct {
						S string `json:"string"`
						I int    `json:"int"`
					}{
						"Test",
						34,
					})
			},
			"{\"string\":\"Test\",\"int\":34}\n",
		},
		{
			func(w http.ResponseWriter, r *http.Request) {
				RespondBadRequest(w, "Bad request")
			},
			"{\"error\":\"Bad request\"}\n",
		},
	}

	for _, test := range tests {
		func() {
			ts := createTestServer(test.Function)
			defer ts.Close()

			body, _, err := callTestServer(ts)
			if !assert.Nil(err) {
				return
			}

			assert.Equal(test.Json, body)
		}()
	}
}

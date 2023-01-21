package opslevel_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	ol "github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2022"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(output)
	flag.Parse()
	teardown := autopilot.Setup()
	defer teardown()
	os.Exit(m.Run())
}

func ToJson(query autopilot.GraphqlQuery) string {
	bytes, _ := json.Marshal(query)
	return string(bytes)
}

func Parse(r *http.Request) autopilot.GraphqlQuery {
	output := autopilot.GraphqlQuery{
		Variables: map[string]interface{}{},
	}
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return output
	}
	json.Unmarshal(bytes, &output)
	return output
}

func LogRaw() autopilot.RequestValidation {
	return func(r *http.Request) {
		defer r.Body.Close()
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Err(err)
		}
		log.Info().Msg(string(bytes))
	}
}

func Templated(input string) string {
	response, err := autopilot.Templater.Use(input)
	if err != nil {
		panic(err)
	}
	return response
}

func TemplatedResponse(response string) autopilot.ResponseWriter {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, Templated(response))
	}
}

func GraphQLQueryTemplate(request string) autopilot.GraphqlQuery {
	exp := autopilot.GraphqlQuery{
		Variables: nil,
	}
	json.Unmarshal([]byte(Templated(request)), &exp)
	return exp
}

func GraphQLQueryTemplatedValidation(t *testing.T, request string) autopilot.RequestValidation {
	return func(r *http.Request) {
		autopilot.Equals(t, ToJson(GraphQLQueryTemplate(request)), ToJson(Parse(r)))
	}
}

func ABetterTestClient(t *testing.T, endpoint string, request string, response string) *ol.Client {
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(autopilot.RegisterEndpoint(fmt.Sprintf("/LOCAL_TESTING/%s", endpoint),
		TemplatedResponse(response),
		GraphQLQueryTemplatedValidation(t, request))))
}

func ATestClient(t *testing.T, endpoint string) *ol.Client {
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(autopilot.RegisterEndpoint(fmt.Sprintf("/LOCAL_TESTING/%s", endpoint),
		autopilot.FixtureResponse(fmt.Sprintf("%s_response.json", endpoint)),
		autopilot.GraphQLQueryFixtureValidation(t, fmt.Sprintf("%s_request.json", endpoint)))))
}

type TestRequest struct {
	Request  string
	Response string
}

func RegisterPaginatedEndpoint(t *testing.T, endpoint string, requests ...TestRequest) string {
	url := fmt.Sprintf("/LOCAL_TESTING/%s", endpoint)
	requestCount := 0
	autopilot.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		GraphQLQueryTemplatedValidation(t, requests[requestCount].Request)(r)
		TemplatedResponse(requests[requestCount].Response)(w)
		requestCount += 1
	})
	return autopilot.Server.URL + url
}

func APaginatedTestClient(t *testing.T, endpoint string, requests ...TestRequest) *ol.Client {
	url := RegisterPaginatedEndpoint(t, endpoint, requests...)
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(url))
}

func ATestClientAlt(t *testing.T, response string, request string) *ol.Client {
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(autopilot.RegisterEndpoint(fmt.Sprintf("/LOCAL_TESTING/%s__%s", response, request),
		autopilot.FixtureResponse(fmt.Sprintf("%s_response.json", response)),
		autopilot.GraphQLQueryFixtureValidation(t, fmt.Sprintf("%s_request.json", request)))))
}

func ATestClientSkipRequest(t *testing.T, endpoint string) *ol.Client {
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(autopilot.RegisterEndpoint(fmt.Sprintf("/LOCAL_TESTING/%s", endpoint),
		autopilot.FixtureResponse(fmt.Sprintf("%s_response.json", endpoint)),
		autopilot.SkipRequestValidation())))
}

func ATestClientLogRequest(t *testing.T, endpoint string) *ol.Client {
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(autopilot.RegisterEndpoint(fmt.Sprintf("/LOCAL_TESTING/%s", endpoint),
		autopilot.FixtureResponse(fmt.Sprintf("%s_response.json", endpoint)),
		LogRaw())))
}

func TestClientQuery(t *testing.T) {
	// Arrange
	url := autopilot.RegisterEndpoint("/LOCAL_TESTING/account", autopilot.FixtureResponse("account_response.json"), autopilot.GraphQLQueryValidation(t, "{account{id}}"))
	client := ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(url))
	var q struct {
		Account struct {
			Id ol.ID
		}
	}
	var v map[string]interface{}
	// Act
	err := client.Query(&q, v)
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, "1234", string(q.Account.Id))
}

/*
These tests don't work very well with our autopilot endpoint stuff because they need to make recursive calls
We need to figure out a better way to handle mapping different payloads to different requests

// This test should infinitly recurse on the Service pagination call
func TestClientQueryPagination(t *testing.T) {
	t.Parallel()
	// Arrange
	url := autopilot.RegisterEndpoint("/pagination", autopilot.FixtureResponse("pagination_response.json"), autopilot.SkipRequestValidation())
	client := NewGQLClient("X", SetURL(url))
	timeout := time.After(3 * time.Second)
	done := make(chan bool)

	// Act
	go func() {
		_, err := client.ListServices()
		autopilot.Ok(t, err)
		done <- true
	}()

	// Assert
	select {
	case <-timeout:
		// Test Was running infinitely in a pagination recursion - this is a success
	case <-done:
		t.Fatal("TestClientQueryPagination did not infinitely recurse on pagination of Services")
	}
}

// This test should infinitly recurse on the Service.Tags nested pagination call
func TestClientQueryNestedPagination(t *testing.T) {
	t.Parallel()
	// Arrange
	//url := autopilot.RegisterEndpoint("/query_nested_pagination", "query_nested_pagination.json")
	url := autopilot.RegisterEndpoint("/nested_pagination", autopilot.FixtureResponse("nested_pagination_response.json"), autopilot.GraphQLQueryFixtureValidation(t, "nested_pagination_request.json"))
	client := NewGQLClient("X", SetURL(url))
	timeout := time.After(3 * time.Second)
	done := make(chan bool)

	// Act
	go func() {
		_, err := client.ListServices()
		autopilot.Ok(t, err)
		done <- true
	}()

	// Assert
	select {
	case <-timeout:
		// Test Was running infinitely in a nested pagination recursion - this is a success
	case <-done:
		t.Fatal("TestClientQueryNestedPagination did not infinitely recurse on nested pagination of Service.Tags")
	}
}
*/

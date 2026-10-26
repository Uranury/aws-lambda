package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func newReq(method, path, body string) events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest{
		RawPath: path,
		Body:    body,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{Method: method},
		},
	}
}

func TestHealth(t *testing.T) {
	res, _ := handler(context.Background(), newReq("GET", "/health", ""))
	if res.StatusCode != 200 {
		t.Fatalf("got status %d", res.StatusCode)
	}
}

func TestAnalyze(t *testing.T) {
	res, _ := handler(context.Background(), newReq("POST", "/analyze", `{"text":"hello big\nworld"}`))
	if res.StatusCode != 200 {
		t.Fatalf("got status %d: %s", res.StatusCode, res.Body)
	}
	var out analyzeResponse
	if err := json.Unmarshal([]byte(res.Body), &out); err != nil {
		t.Fatal(err)
	}
	if out.Words != 3 || out.Lines != 2 || out.Chars != 15 {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestBadJSON(t *testing.T) {
	res, _ := handler(context.Background(), newReq("POST", "/analyze", `not json`))
	if res.StatusCode != 400 {
		t.Fatalf("got status %d", res.StatusCode)
	}
}

func TestNotFound(t *testing.T) {
	res, _ := handler(context.Background(), newReq("GET", "/nope", ""))
	if res.StatusCode != 404 {
		t.Fatalf("got status %d", res.StatusCode)
	}
}

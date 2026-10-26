package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type analyzeRequest struct {
	Text string `json:"text"`
}

type analyzeResponse struct {
	Words int `json:"words"`
	Chars int `json:"chars"`
	Lines int `json:"lines"`
}

// handler receives the v2 HTTP payload, which is what both Lambda Function
// URLs and API Gateway HTTP APIs send. Routing is just a switch on method+path.
func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch req.RequestContext.HTTP.Method + " " + req.RawPath {
	case "GET /health":
		return respond(http.StatusOK, map[string]string{"status": "ok"})
	case "POST /analyze":
		return analyze(req)
	default:
		return respond(http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

func analyze(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	body := req.Body
	if req.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return respond(http.StatusBadRequest, map[string]string{"error": "invalid base64 body"})
		}
		body = string(decoded)
	}

	var in analyzeRequest
	if err := json.Unmarshal([]byte(body), &in); err != nil {
		return respond(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	lines := 0
	if in.Text != "" {
		lines = strings.Count(in.Text, "\n") + 1
	}

	return respond(http.StatusOK, analyzeResponse{
		Words: len(strings.Fields(in.Text)),
		Chars: utf8.RuneCountInString(in.Text),
		Lines: lines,
	})
}

func respond(status int, v any) (events.APIGatewayV2HTTPResponse, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(b),
	}, nil
}

func main() {
	lambda.Start(handler)
}

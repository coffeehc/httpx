package httpxcommons

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/proto"
)

func TestPublishedListDataSourceAndWireCompatibility(t *testing.T) {
	for _, test := range []struct {
		value ListData[string]
		want  string
	}{
		{ListData[string]{List: nil, Total: int64(0)}, `{"list":null,"total":0}`},
		{ListData[string]{List: []string{}, Total: int64(0)}, `{"list":[],"total":0}`},
		{ListData[string]{List: []string{"one"}, Total: int64(9)}, `{"list":["one"],"total":9}`},
	} {
		got, err := json.Marshal(test.value)
		if err != nil || string(got) != test.want {
			t.Fatalf("%s %v", got, err)
		}
	}
}

func TestPublishedEnvelopeCompatibility(t *testing.T) {
	value := &AjaxResponse{Success: true, Message: "", Payload: []string{}, RequestID: "req", Code: 0}
	var response BaseResponse = value
	if !response.IsSuccess() || response.GetMessage() != "" || response.GetRequestId() != "req" || response.GetCode() != 0 {
		t.Fatal("accessor changed")
	}
	if _, ok := response.GetPayload().([]string); !ok {
		t.Fatal("payload changed")
	}
	got, err := json.Marshal(value)
	if err != nil || string(got) != `{"success":true,"message":"","payload":[],"request_id":"req","code":0}` {
		t.Fatalf("%s %v", got, err)
	}
}

func TestPublishedSendHelpersCompatibility(t *testing.T) {
	for _, test := range []struct {
		name, accept, want string
		status             int
		handler            fiber.Handler
	}{
		{"success", "application/json", `{"success":true,"message":"","payload":[],"request_id":"","code":0}`, 200, func(c fiber.Ctx) error { return SendSuccess(c, []string{}, 0) }},
		{"error", "application/json", `{"success":false,"message":"denied","payload":null,"request_id":"","code":40301}`, 403, func(c fiber.Ctx) error { return SendError(c, "denied", 40301, 403) }},
		{"redirect", "*/*", `{"success":false,"message":"sign in","payload":null,"request_id":"","code":40101,"redirect":"/login"}`, 401, func(c fiber.Ctx) error { return SendErrorWithRedirect(c, "sign in", "/login", 40101, 401) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", test.handler)
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.Header.Set("Accept", test.accept)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil || resp.StatusCode != test.status || string(body) != test.want {
				t.Fatalf("%d %s %v", resp.StatusCode, body, err)
			}
		})
	}
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error { return SendError(c, "denied", 40301, 403) })
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set("Accept", "application/x-protobuf")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var response PBResponse
	if err := proto.Unmarshal(body, &response); err != nil || resp.StatusCode != 403 || response.Success || response.Code != 40301 || response.Message != "denied" {
		t.Fatalf("%d %+v %v", resp.StatusCode, &response, err)
	}
}

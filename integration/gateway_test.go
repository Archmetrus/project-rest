package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"example.com/project-rest/internal/store"
	"example.com/project-rest/internal/testdb"
	"example.com/project-rest/internal/transport"
)

func gateway(t *testing.T) string {
	t.Helper()
	if address := os.Getenv("TEST_GATEWAY_ADDR"); address != "" {
		return "http://" + address
	}
	user := httptest.NewServer(transport.UserHandler(store.UserStore{DB: testdb.New(t, "user")}))
	t.Cleanup(user.Close)
	hr := httptest.NewServer(transport.LeaveHandler(store.LeaveStore{DB: testdb.New(t, "hr")}))
	t.Cleanup(hr.Close)
	handler, err := transport.Gateway(strings.TrimPrefix(user.URL, "http://"), strings.TrimPrefix(hr.URL, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	gateway := httptest.NewServer(handler)
	t.Cleanup(gateway.Close)
	return gateway.URL
}

func request(t *testing.T, method, url, body string, want int, output any) {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != want {
		t.Fatalf("%s %s: got %d, want %d: %s", method, url, response.StatusCode, want, data)
	}
	if output != nil {
		if err := json.Unmarshal(data, output); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGatewayCRUD(t *testing.T) {
	base := gateway(t)
	var user store.User
	request(t, "POST", base+"/users", `{"name":"Ada","email":"ada@example.com","address":"Ankara"}`, 201, &user)
	if user.ID == 0 || user.Address != "Ankara" {
		t.Fatalf("bad user: %+v", user)
	}
	userURL := fmt.Sprintf("%s/users/%d", base, user.ID)
	var gotUser store.User
	request(t, "GET", userURL, "", 200, &gotUser)
	if gotUser != user {
		t.Fatalf("get user mismatch: %+v", gotUser)
	}
	var users []store.User
	request(t, "GET", base+"/users", "", 200, &users)
	found := false
	for _, v := range users {
		if v == user {
			found = true
		}
	}
	if !found {
		t.Fatal("created user missing from list")
	}
	request(t, "PUT", userURL, `{"name":"Grace","email":"grace@example.com","address":"Istanbul"}`, 200, &gotUser)
	if gotUser.ID != user.ID || gotUser.Name != "Grace" || gotUser.Email != "grace@example.com" || gotUser.Address != "Istanbul" {
		t.Fatalf("bad update: %+v", gotUser)
	}
	request(t, "GET", userURL, "", 200, &user)
	if user != gotUser {
		t.Fatal("user update was not persisted")
	}

	// Unknown owners, reversed dates and multiple records are intentionally allowed.
	leaveBody := `{"user_id":999999,"start_date":"2026-10-20","end_date":"2026-10-01"}`
	var leave, second store.Leave
	request(t, "POST", base+"/leaves", leaveBody, 201, &leave)
	request(t, "POST", base+"/leaves", leaveBody, 201, &second)
	if leave.ID == 0 || second.ID == leave.ID {
		t.Fatal("leave IDs are not unique")
	}
	leaveURL := fmt.Sprintf("%s/leaves/%d", base, leave.ID)
	var gotLeave store.Leave
	request(t, "GET", leaveURL, "", 200, &gotLeave)
	if gotLeave != leave {
		t.Fatal("leave round trip mismatch")
	}
	var leaves []store.Leave
	request(t, "GET", base+"/leaves", "", 200, &leaves)
	found = false
	for _, v := range leaves {
		if v == leave {
			found = true
		}
	}
	if !found {
		t.Fatal("created leave missing from list")
	}
	request(t, "PUT", leaveURL, `{"user_id":123456,"start_date":"2026-12-02","end_date":"2026-12-01"}`, 200, &gotLeave)
	if gotLeave.ID != leave.ID || gotLeave.UserID != 123456 || gotLeave.StartDate != "2026-12-02" || gotLeave.EndDate != "2026-12-01" {
		t.Fatalf("bad leave update: %+v", gotLeave)
	}
	request(t, "GET", leaveURL, "", 200, &leave)
	if leave != gotLeave {
		t.Fatal("leave update was not persisted")
	}
	request(t, "DELETE", fmt.Sprintf("%s/leaves/%d", base, second.ID), "", 204, nil)
	for _, item := range []struct{ url, body string }{{userURL, `{}`}, {leaveURL, leaveBody}} {
		request(t, "DELETE", item.url, "", 204, nil)
		request(t, "GET", item.url, "", 404, nil)
		request(t, "PUT", item.url, item.body, 404, nil)
		request(t, "DELETE", item.url, "", 404, nil)
	}
}

func TestMalformedRequests(t *testing.T) {
	base := gateway(t)
	for _, path := range []string{"/users", "/leaves"} {
		request(t, "GET", base+path+"/bad-id", "", 400, nil)
		for _, body := range []string{`{`, `null`, `[]`, `{} {}`, `{"id":"bad"}`} {
			request(t, "POST", base+path, body, 400, nil)
		}
	}
	request(t, "POST", base+"/leaves", `{"user_id":"wrong-type"}`, 400, nil)
}

func TestGatewayUnavailable(t *testing.T) {
	upstream := httptest.NewServer(http.NotFoundHandler())
	address := strings.TrimPrefix(upstream.URL, "http://")
	upstream.Close()
	handler, err := transport.Gateway(address, address)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	request(t, "GET", server.URL+"/users", "", 502, nil)
}

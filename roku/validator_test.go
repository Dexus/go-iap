package roku

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestHandleError(t *testing.T) {
	var expected, actual error

	server, client := testTools(
		497,
		"{\"errorMessage\":\"Purchase token/app user mismatch\",\"status\":\"Failure\"}",
	)
	defer server.Close()

	// status 400
	expected = errors.New("Purchase token/app user mismatch")
	_, actual = client.Verify("{transactionid}")
  
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNew(t *testing.T) {
	expected := Client{
		URL:      ProductionURL,
		TimeOut:  time.Second * 5,
		DevToken: "devToken",
	}

	actual := New("devToken")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithEnvironment(t *testing.T) {
	expected := Client{
		URL:      ProductionURL,
		TimeOut:  time.Second * 5,
		DevToken: "devToken",
	}

	actual := New("devToken")
	os.Clearenv()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithConfig(t *testing.T) {
	config := Config{
		IsProduction: true,
		DevToken:     "devToken",
		TimeOut:      time.Second * 2,
	}

	expected := Client{
		URL:      ProductionURL,
		TimeOut:  time.Second * 2,
		DevToken: "devToken",
	}

	actual := NewWithConfig(config)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithConfigTimeout(t *testing.T) {
	config := Config{
		IsProduction: true,
		DevToken:     "devToken",
	}

	expected := Client{
		URL:      ProductionURL,
		TimeOut:  time.Second * 5,
		DevToken: "devToken",
	}

	actual := NewWithConfig(config)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerify(t *testing.T) {
	server, client := testTools(
		200,
		"{\"transactionId\":\"{transactionid}\",\"purchaseDate\":\"2012-07-22T14:59:50\",\"channelName\":\"123Video\",\"productName\":\"123Video Monthly Subscription\",\"productId\":\"NETMONTH\",\"amount\":9.99,\"currency\":\"USD\",\"quantity\":1,\"expirationDate\":\"2012-08-22T14:59:50\", \"originalPurchaseDate\":\"2010-08-22T14:59:50\", \"status\":\"Success\", \"errorMessage\":\"error_message\"}",
	)
	defer server.Close()

	expected := IAPValidationResponse{
		TransactionID:        "{transactionid}",
		PurchaseDate:         "2012-07-22T14:59:50",
		ChannelName:          "123Video",
		ProductName:          "123Video Monthly Subscription",
		ProductID:            "NETMONTH",
		Amount:               9.99,
		Currency:             "USD",
		Quantity:             1,
		ExpirationDate:       "2012-08-22T14:59:50",
		OriginalPurchaseDate: "2010-08-22T14:59:50",
		Status:               "Success",
		ErrorMessage:         "error_message",
	}

	actual, _ := client.Verify("{transactionid}")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyTimeout(t *testing.T) {
	// HTTP 100 is "continue" so it will time out
	server, client := testTools(100, "timeout response")
	defer server.Close()

	expected := errors.New("")
	_, actual := client.Verify("timeout")
	if !reflect.DeepEqual(reflect.TypeOf(actual), reflect.TypeOf(expected)) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func testTools(code int, body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	client := &Client{URL: server.URL, TimeOut: time.Second * 2, DevToken: "devToken"}
	return server, client
}

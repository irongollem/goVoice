package telnyx

import (
	"io"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestSendCommand(t *testing.T) {
	// Create a new Telnyx instance
	telnyx := &Telnyx{}

	// Create a mock HTTP server
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.telnyx.com/v2/calls/call_control_id/actions/test_command",
		httpmock.NewStringResponder(200, `{"result":"success"}`))

	// Set the Telnyx API URL to the mock server URL
	telnyx.APIUrl, _ = url.Parse("https://api.telnyx.com/v2")

	foo := &CommandPayload{}
	// Call SendCommand with the mock data
	response, err := telnyx.sendCommand("POST",foo, "calls", "call_control_id", "actions", "test_command")
	if err != nil {
		t.Errorf("Error sending command: %s", err)
	}

	// Check that the response is correct
	expectedResponse := `{"result":"success"}`
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
	}
	if string(body) != expectedResponse {
		t.Errorf("Expected response body %s, got %s", expectedResponse, string(body))
	}
}

package addressbooks

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/NomNes/go-sendpulse"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var restyClient = resty.New()
var collection Collection

func TestMain(m *testing.M) {
	collection = New(go_sendpulse.Client{Client: restyClient})
	os.Exit(m.Run())
}

func TestGetList(t *testing.T) {
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	ab := []map[string]interface{}{
		{"id": 1, "name": "Test 1"},
		{"id": 2, "name": "Test 2"},
	}

	httpmock.RegisterResponder("GET", "/addressbooks", httpmock.NewJsonResponderOrPanic(200, ab))

	list, err := collection.GetList(ctx, nil)
	if assert.NoError(t, err) {
		if assert.IsType(t, []Item{}, list) {
			if assert.Len(t, list, 2) {
				assert.Equal(t, ab[0]["id"], list[0].ID())
				assert.Equal(t, ab[1]["id"], list[1].ID())
				assert.Equal(t, ab[0]["name"], list[0].Name())
				assert.Equal(t, ab[1]["name"], list[1].Name())
			}
		}
	}
}

func TestGetOne(t *testing.T) {
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	ab := map[string]interface{}{
		"id":   123,
		"name": "Test Mailing List",
	}

	httpmock.RegisterResponder("GET", "/addressbooks/123", httpmock.NewJsonResponderOrPanic(200, []map[string]interface{}{ab}))

	item, err := collection.GetOne(ctx, 123)
	if assert.NoError(t, err) {
		if assert.IsType(t, &Item{}, item) {
			assert.Equal(t, ab["id"], item.ID())
			assert.Equal(t, ab["name"], item.Name())
		}
	}
}

func TestAddEmails(t *testing.T) {
	httpmock.ActivateNonDefault(restyClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/addressbooks/123/emails", func(req *http.Request) (*http.Response, error) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		var v map[string]interface{}
		err = json.Unmarshal(body, &v)
		if err != nil {
			return nil, err
		}
		assert.IsType(t, []interface{}{}, v["emails"])
		emails := v["emails"].([]interface{})
		assert.Len(t, emails, 1)
		assert.IsType(t, map[string]interface{}{}, emails[0])
		email := emails[0].(map[string]interface{})
		assert.IsType(t, "", email["email"])
		assert.Equal(t, "test@mail.com", email["email"])
		return httpmock.NewJsonResponse(200, go_sendpulse.ResultResponse{Result: true})
	})

	err := collection.AddEmails(ctx, 123, []string{"test@mail.com"}, nil)
	assert.NoError(t, err)
}

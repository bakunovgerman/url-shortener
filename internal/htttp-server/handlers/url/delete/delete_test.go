package delete

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shorter/base/internal/htttp-server/handlers/url/delete/mocks"
	"url-shorter/base/internal/htttp-server/handlers/url/save"
	"url-shorter/base/internal/lib/logger/handlers/slogdiscard"
	"url-shorter/base/internal/storage"

	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "success",
			alias: "alias",
		},
		{
			name:      "invalid alias",
			alias:     "invalid",
			respError: "alias not found",
			mockError: storage.AliasNotFound,
		},
		{
			name:      "error delete alias",
			alias:     "alias",
			respError: "failed to delete",
			mockError: errors.New("failed to delete"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)
			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return("", tc.mockError).
					Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlDeleterMock)
			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)

			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}

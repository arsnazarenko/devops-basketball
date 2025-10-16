package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/arsnazarenko/devops-basketball/internal/apperrors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/oapi-codegen/nethttp-middleware"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTestServer(mockUC *MockPlayer) *httptest.Server {
	// Create server implementation with mock
	serversImpl := NewPlayersServerImpl(mockUC)

	// Create chi router
	r := chi.NewRouter()
	// Add CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Load swagger spec
	swagger, err := gen.GetSwagger()
	if err != nil {
		panic("Error loading swagger spec: " + err.Error())
	}
	swagger.Servers = nil

	// Add validation middleware
	r.Use(nethttpmiddleware.OapiRequestValidator(swagger))

	// Create strict handler
	server := gen.NewStrictHandler(serversImpl, []gen.StrictMiddlewareFunc{})
	gen.HandlerFromMux(server, r)

	return httptest.NewServer(r)
}

func TestCreatePlayer(t *testing.T) {
	mockUC := &MockPlayer{}
	server := setupTestServer(mockUC)
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		playerCreate := &gen.PlayerCreate{
			Name:        "John",
			Surname:     "Doe",
			Age:         25,
			Height:      1900,
			Weight:      85000,
			Citizenship: "USA",
			Role:        "PG",
			TeamId:      1,
		}

		expectedPlayer := &gen.Player{
			Id:          1,
			Name:        "John",
			Surname:     "Doe",
			Age:         25,
			Height:      1900,
			Weight:      85000,
			Citizenship: "USA",
			Role:        gen.PlayerRole("PG"),
			TeamId:      1,
		}

		mockUC.On("CreatePlayer", mock.Anything, playerCreate).Return(expectedPlayer, nil).Once()

		body, _ := json.Marshal(playerCreate)
		resp, err := http.Post(server.URL+"/players", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var response gen.Player
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, *expectedPlayer, response)

		mockUC.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		playerCreate := &gen.PlayerCreate{
			Name:        "Jane",
			Surname:     "Smith",
			Age:         30,
			Height:      1800,
			Weight:      75000,
			Citizenship: "CAN",
			Role:        "C",
			TeamId:      999,
		}

		mockUC.On("CreatePlayer", mock.Anything, playerCreate).Return((*gen.Player)(nil), apperrors.ErrTeamNotFound).Once()

		body, _ := json.Marshal(playerCreate)
		resp, err := http.Post(server.URL+"/players", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockUC.AssertExpectations(t)
	})

	t.Run("invalid name too long", func(t *testing.T) {
		playerCreate := &gen.PlayerCreate{
			Name:        string(make([]byte, 51)), // 51 char > max 50
			Surname:     "Doe",
			Age:         25,
			Height:      1900,
			Weight:      85000,
			Citizenship: "USA",
			Role:        "PG",
			TeamId:      1,
		}

		body, _ := json.Marshal(playerCreate)
		resp, err := http.Post(server.URL+"/players", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid age too low", func(t *testing.T) {
		playerCreate := &gen.PlayerCreate{
			Name:        "John",
			Surname:     "Doe",
			Age:         10, // < min 15
			Height:      1900,
			Weight:      85000,
			Citizenship: "USA",
			Role:        "PG",
			TeamId:      1,
		}

		body, _ := json.Marshal(playerCreate)
		resp, err := http.Post(server.URL+"/players", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid role", func(t *testing.T) {
		playerCreate := &gen.PlayerCreate{
			Name:        "John",
			Surname:     "Doe",
			Age:         25,
			Height:      1900,
			Weight:      85000,
			Citizenship: "USA",
			Role:        "INVALID", // not in enum
			TeamId:      1,
		}

		body, _ := json.Marshal(playerCreate)
		resp, err := http.Post(server.URL+"/players", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetPlayer(t *testing.T) {
	mockUC := &MockPlayer{}
	server := setupTestServer(mockUC)
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		playerID := int64(1)
		expectedPlayer := &gen.Player{
			Id:          1,
			Name:        "John",
			Surname:     "Doe",
			Age:         25,
			Height:      1900,
			Weight:      85000,
			Citizenship: "USA",
			Role:        gen.PlayerRole("PG"),
			TeamId:      1,
		}

		mockUC.On("GetPlayer", mock.Anything, playerID).Return(expectedPlayer, nil).Once()

		resp, err := http.Get(server.URL + "/players/1")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response gen.Player
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, *expectedPlayer, response)

		mockUC.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		playerID := int64(999)

		mockUC.On("GetPlayer", mock.Anything, playerID).Return((*gen.Player)(nil), apperrors.ErrPlayerNotFound).Once()

		resp, err := http.Get(server.URL + "/players/999")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockUC.AssertExpectations(t)
	})

	t.Run("invalid id format", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/players/abc")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestListPlayers(t *testing.T) {
	mockUC := &MockPlayer{}
	server := setupTestServer(mockUC)
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		expectedPlayers := []gen.Player{
			{
				Id:          1,
				Name:        "John",
				Surname:     "Doe",
				Age:         25,
				Height:      1900,
				Weight:      85000,
				Citizenship: "USA",
				Role:        gen.PlayerRole("PG"),
				TeamId:      1,
			},
			{
				Id:          2,
				Name:        "Jane",
				Surname:     "Smith",
				Age:         30,
				Height:      1800,
				Weight:      75000,
				Citizenship: "CAN",
				Role:        gen.PlayerRole("C"),
				TeamId:      2,
			},
		}

		mockUC.On("GetPlayerList", mock.Anything, uint64(20), uint64(1)).Return(expectedPlayers, nil).Once()

		resp, err := http.Get(server.URL + "/players?page_size=20&page_number=1")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response []gen.Player
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, expectedPlayers, response)

		mockUC.AssertExpectations(t)
	})

	t.Run("invalid page number", func(t *testing.T) {
		// Since OpenAPI validation happens in middleware, no usecase call expected
		resp, err := http.Get(server.URL + "/players?page_size=20&page_number=0")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// No usecase call expected due to middleware validation
	})

	t.Run("invalid page size", func(t *testing.T) {
		// Since OpenAPI validation happens in middleware, no usecase call expected
		resp, err := http.Get(server.URL + "/players?page_size=0&page_number=1")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// No usecase call expected due to middleware validation
	})

	t.Run("page size too large", func(t *testing.T) {
		// page_size max 100
		resp, err := http.Get(server.URL + "/players?page_size=101&page_number=1")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestUpdatePlayer(t *testing.T) {
	mockUC := &MockPlayer{}
	server := setupTestServer(mockUC)
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		playerID := int64(1)
		playerUpdate := &gen.PlayerUpdate{
			Name:        stringPtr("John"),
			Surname:     stringPtr("Doe"),
			Age:         intPtr(26),
			Height:      intPtr(1910),
			Weight:      intPtr(86000),
			Citizenship: stringPtr("USA"),
			Role:        playerUpdateRolePtr(gen.PlayerUpdateRole("PG")),
			TeamId:      int64Ptr(1),
		}

		expectedPlayer := &gen.Player{
			Id:          1,
			Name:        "John",
			Surname:     "Doe",
			Age:         26,
			Height:      1910,
			Weight:      86000,
			Citizenship: "USA",
			Role:        gen.PlayerRole("PG"),
			TeamId:      1,
		}

		mockUC.On("UpdatePlayer", mock.Anything, playerID, playerUpdate).Return(expectedPlayer, nil).Once()

		body, _ := json.Marshal(playerUpdate)
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/players/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response gen.Player
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, *expectedPlayer, response)

		mockUC.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		playerID := int64(999)
		playerUpdate := &gen.PlayerUpdate{
			Name: stringPtr("Jane"),
		}

		mockUC.On("UpdatePlayer", mock.Anything, playerID, playerUpdate).Return((*gen.Player)(nil), apperrors.ErrPlayerNotFound).Once()

		body, _ := json.Marshal(playerUpdate)
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/players/999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockUC.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		playerID := int64(1)
		playerUpdate := &gen.PlayerUpdate{
			TeamId: int64Ptr(999),
		}

		mockUC.On("UpdatePlayer", mock.Anything, playerID, playerUpdate).Return((*gen.Player)(nil), apperrors.ErrTeamNotFound).Once()

		body, _ := json.Marshal(playerUpdate)
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/players/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockUC.AssertExpectations(t)
	})
}

// Helper functions for pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func playerUpdateRolePtr(r gen.PlayerUpdateRole) *gen.PlayerUpdateRole {
	return &r
}

func TestDeletePlayer(t *testing.T) {
	mockUC := &MockPlayer{}
	server := setupTestServer(mockUC)
	defer server.Close()

	t.Run("success", func(t *testing.T) {
		playerID := int64(1)

		mockUC.On("DeletePlayer", mock.Anything, playerID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/players/1", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		mockUC.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		playerID := int64(999)

		mockUC.On("DeletePlayer", mock.Anything, playerID).Return(apperrors.ErrPlayerNotFound).Once()

		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/players/999", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockUC.AssertExpectations(t)
	})

	t.Run("invalid id format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/players/xyz", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

package enrollment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Define the structure of the initial HTTP POST request payload for FastAPI
type FastAPIEnrollInitialRequest struct {
	Username string `json:"username"`
	UniqueID string `json:"unique_id"`
}

// Define the structure of the initial HTTP POST response from FastAPI
type FastAPIEnrollInitialResponse struct {
	Message             string `json:"message"`
	EnrollmentSessionID string `json:"enrollment_session_id"`
	WebsocketURL        string `json:"websocket_url"`
}

// Define the structure of WebSocket messages from FastAPI
type WebSocketMessage struct {
	Type      string          `json:"type"` // e.g., "status", "COMPLETED", "FAILED"
	Data      json.RawMessage `json:"data"` // Raw JSON to be parsed based on Type
	Timestamp float64         `json:"timestamp"`
}

// Data for "status" type messages
type WebSocketStatusData struct {
	Stage   string `json:"stage"`
	Details string `json:"details"`
}

// Data for "COMPLETED" type messages
type WebSocketCompletedData struct {
	Message  string `json:"message"`
	UID      string `json:"uid"`
	Username string `json:"username"`
	UniqueID string `json:"unique_id"`
	Success  bool   `json:"success"`
}

// Data for "FAILED" type messages
type WebSocketFailedData struct {
	Message string            `json:"message"`
	Success bool              `json:"success"`
	Details map[string]string `json:"details"` // Optional error details
}

// EnrollmentService is a service that handles enrollment operations.
type EnrollmentService struct {
	FastAPIBaseURL string // Base URL of your FastAPI server (e.g., "http://localhost:8000")
}

// NewEnrollmentService creates a new instance of EnrollmentService.
func NewEnrollmentService(fastAPIBaseURL string) *EnrollmentService {
	return &EnrollmentService{
		FastAPIBaseURL: fastAPIBaseURL,
	}
}

func (es *EnrollmentService) enroll(ctx context.Context, req FastAPIEnrollInitialRequest) (WebSocketCompletedData, error) {
	slog.Info("Initiating enrollment for:", "user", req.Username, "unique id", req.UniqueID)

	initialPayload := FastAPIEnrollInitialRequest{
		Username: req.Username,
		UniqueID: req.UniqueID,
	}

	jsonPayload, err := json.Marshal(initialPayload)
	if err != nil {
		slog.Error("Failed to marshal initial JSON payload", "error", err)
		return WebSocketCompletedData{}, fmt.Errorf("failed to marshal initial JSON payload: %w", err)
	}

	httpPostURL := fmt.Sprintf("%s/cs/enroll", es.FastAPIBaseURL)
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", httpPostURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		slog.Error("Failed to create initial HTTP request", "error", err)
		return WebSocketCompletedData{}, fmt.Errorf("failed to create initial HTTP request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		slog.Error("Failed to send initial HTTP request to FastAPI", "error", err)
		return WebSocketCompletedData{}, fmt.Errorf("failed to send initial HTTP request to FastAPI: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(httpResponse.Body)

		return WebSocketCompletedData{}, fmt.Errorf("FastAPI initial enrollment endpoint returned non-OK status: %s, Body: %s", httpResponse.Status, buf.String())
	}

	var initialFastAPIResponse FastAPIEnrollInitialResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&initialFastAPIResponse)
	if err != nil {
		slog.Error("Failed to decode initial FastAPI response", "error", err)
		return WebSocketCompletedData{}, fmt.Errorf("failed to decode initial FastAPI response: %w", err)
	}

	slog.Info("FastAPI initiated enrollment:", "initmessage", initialFastAPIResponse.Message, "Session ID", initialFastAPIResponse.EnrollmentSessionID, "WebSocket_URL", initialFastAPIResponse.WebsocketURL)

	wsURL, err := url.Parse(initialFastAPIResponse.WebsocketURL)
	if err != nil {
		slog.Error("Invalid WebSocket URL from FastAPI", "error", err)
		return WebSocketCompletedData{}, fmt.Errorf("invalid WebSocket URL received from FastAPI: %w", err)
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, wsURL.String(), nil)
	if err != nil {
		slog.Error("Failed to connect to WebSocket", "url", wsURL.String(), "error", err)
		return WebSocketCompletedData{}, fmt.Errorf("failed to connect to WebSocket %s: %w", wsURL.String(), err)
	}
	defer conn.Close()

	log.Printf("Connected to WebSocket: %s", wsURL.String())

	enrollmentResultChan := make(chan *WebSocketCompletedData)
	errorChan := make(chan error, 1)

	go func() {
		defer close(enrollmentResultChan)
		defer close(errorChan)

		for {
			select {
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						log.Printf("WebSocket closed normally: %v", err)
						errorChan <- fmt.Errorf("WebSocket closed normally: %w", err)
					} else {
						log.Printf("WebSocket read error: %v", err)
						errorChan <- fmt.Errorf("WebSocket read error: %w", err)
					}
					return
				}

				var wsMsg WebSocketMessage
				if err := json.Unmarshal(message, &wsMsg); err != nil {
					log.Printf("Failed to unmarshal WebSocket message: %v, message: %s", err, string(message))
					continue
				}

				log.Printf("Received WS message type: %s, Data: %s", wsMsg.Type, string(wsMsg.Data))

				switch wsMsg.Type {
				case "STATUS":
				case "COMPLETED":
					var completedData WebSocketCompletedData
					if err := json.Unmarshal(wsMsg.Data, &completedData); err != nil {
						log.Printf("Failed to unmarshal WS completed data: %v", err)
						errorChan <- fmt.Errorf("failed to parse completion data: %w", err)
						return
					}
					log.Printf("Enrollment COMPLETE for UID: %s, Message: %s", completedData.UID, completedData.Message)
					enrollmentResultChan <- &WebSocketCompletedData{
						Success:  completedData.Success,
						Message:  completedData.Message,
						UID:      completedData.UID,
						Username: completedData.Username,
						UniqueID: completedData.UniqueID,
					}
					return
				case "FAILED":
					var failedData WebSocketFailedData
					if err := json.Unmarshal(wsMsg.Data, &failedData); err != nil {
						log.Printf("Failed to unmarshal WS failed data: %v", err)
						errorChan <- fmt.Errorf("failed to parse failure data: %w", err)
						return
					}
					log.Printf("Enrollment FAILED: %s", failedData.Message)
					errorChan <- fmt.Errorf("enrollment failed: %s", failedData.Message)
					return
				default:
					log.Printf("Received unknown WS message type: %s", wsMsg.Type)
				}
			}
		}
	}()

	select {
	case res := <-enrollmentResultChan:
		return *res, nil
	case err := <-errorChan:
		slog.Error("Error during WebSocket communication", "error", err)
		return WebSocketCompletedData{}, err
	case <-ctx.Done():
		slog.Error("enrollment process timed out or was cancelled", "error", err) // Context timeout/cancellation
		return WebSocketCompletedData{}, fmt.Errorf("enrollment process timed out or was cancelled: %w", ctx.Err())
	}
}

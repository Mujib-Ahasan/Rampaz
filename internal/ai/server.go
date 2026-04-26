package ai

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	sessionID   string
	chatService *ChatService
	conn        *grpc.ClientConn
}

func NewServer() *Server {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	k8sClient := pb.NewK8SInfoClient(conn)

	llm := NewLLMClient()
	executor := NewToolExecutor(k8sClient)
	chatService := NewChatService(llm, executor)

	return &Server{
		sessionID:   "default-session",
		chatService: chatService,
		conn:        conn,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/chat", s.chatHandler)
	mux.HandleFunc("/api/session", s.sessionHandler)

	// Serve local chatbot UI
	fileServer := http.FileServer(http.Dir("web/chatbot"))
	mux.Handle("/", fileServer)

	return mux
}

func (s *Server) sessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"sessionId": s.sessionID,
	})
}

func (s *Server) chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		writeError(w, "message is required", http.StatusBadRequest)
		return
	}
	fmt.Printf("%s", req.Message)
	answer, err := s.chatService.Chat(r.Context(), req.Message)
	if err != nil {
		log.Printf("chat failed: %v", err)
		writeError(w, "chat failed", http.StatusInternalServerError)
		return
	}

	// Temporary response.
	// Later this will call LLM + Rampaz gRPC tools.
	resp := ChatResponse{
		Answer:    answer,
		SessionID: s.sessionID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

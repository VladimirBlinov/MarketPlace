package fake

import (
	"context"
	"sync"

	"github.com/VladimirBlinov/AuthService/pkg/authservice"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthServiceClientFake struct {
	mu           sync.RWMutex
	sessionStore map[string]*authservice.Session
}

func NewAuthServiceClientFake() *AuthServiceClientFake {
	return &AuthServiceClientFake{
		sessionStore: make(map[string]*authservice.Session),
	}
}

func (ascf *AuthServiceClientFake) Create(ctx context.Context, s *authservice.Session, opts ...grpc.CallOption) (*authservice.SessionID, error) {
	ascf.mu.Lock()
	defer ascf.mu.Unlock()

	sessionID := &authservice.SessionID{
		ID: uuid.New().String(),
	}

	ascf.sessionStore[sessionID.ID] = s
	return sessionID, nil
}

func (ascf *AuthServiceClientFake) Check(ctx context.Context, sID *authservice.SessionID, opts ...grpc.CallOption) (*authservice.Session, error) {
	s, ok := ascf.sessionStore[sID.ID]
	if !ok {
		return nil, grpc.Errorf(codes.NotFound, "SM check: not found")
	}

	return s, nil
}

func (ascf *AuthServiceClientFake) Delete(ctx context.Context, sID *authservice.SessionID, opts ...grpc.CallOption) (*authservice.Nothing, error) {
	_, ok := ascf.sessionStore[sID.ID]
	if !ok {
		return &authservice.Nothing{
			Dummy: false,
		}, grpc.Errorf(codes.NotFound, "SM check: not found")
	}

	delete(ascf.sessionStore, sID.ID)
	return &authservice.Nothing{
		Dummy: true,
	}, nil
}

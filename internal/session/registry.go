package session

import "sync"

type SessionRegistry struct {
	mu       *sync.RWMutex
	sessions map[SessionID]Session
}

func NewSessionRegistry() *SessionRegistry {
	return &SessionRegistry{
		sessions: make(map[SessionID]Session),
		mu:       &sync.RWMutex{},
	}
}

func (r *SessionRegistry) Add(s Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID()] = s
}

func (r *SessionRegistry) Get(id SessionID) (Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[id]
	return s, ok
}

func (r *SessionRegistry) Remove(id SessionID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
}

func (r *SessionRegistry) List() []Session {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Session, 0, len(r.sessions))
	for _, s := range r.sessions {
		out = append(out, s)
	}
	return out
}

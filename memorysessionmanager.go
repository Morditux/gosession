package sessions

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
)

type MemorySessionManager struct {
	keys  map[string]*Session
	mutex sync.RWMutex
}

// NewMemorySessionManager - Create a new MemorySessionManager object
func NewMemorySessionManager() *MemorySessionManager {
	sm := &MemorySessionManager{
		keys: make(map[string]*Session),
	}
	return sm
}

func (sm *MemorySessionManager) Exists(key string) bool {
	sm.mutex.RLock()
	defer sm.mutex.Unlock()
	_, v := sm.keys[key]
	return v
}

func (sm *MemorySessionManager) GetSessionCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return len(sm.keys)
}

func (sm *MemorySessionManager) Add(session *Session) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.keys[session.key] = session
}

func (sm *MemorySessionManager) CreateSession(userName string, isAdmin bool) (*Session, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	session := &Session{
		key:      uuid.New().String(),
		userName: userName,
		lastSeen: time.Now(),
		isAdmin:  isAdmin,
		isLogged: false,
		data:     make(map[string]interface{}),
	}
	sm.keys[session.key] = session
	return session, true
}

func (sm *MemorySessionManager) Remove(key string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	delete(sm.keys, key)
}

func (sm *MemorySessionManager) Get(key string) (*Session, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	session, v := sm.keys[key]
	if v {
		return session, nil
	}
	return &Session{
		key: "",
	}, ErrSessionNotFound
}

func (sm *MemorySessionManager) Clean() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	for k, v := range sm.keys {
		if time.Since(v.lastSeen) > time.Minute*10 {
			delete(sm.keys, k)
		}
	}
}

func (sm *MemorySessionManager) FromBinary(data []byte) *Session {
	type InternalSession struct {
		Key      string
		UserName string
		LastSeen time.Time
		IsAdmin  bool
		IsLogin  bool
		Data     map[string]interface{}
	}
	var internalSession InternalSession
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&internalSession)
	if err != nil {
		panic(err)
	}
	return &Session{
		key:      internalSession.Key,
		userName: internalSession.UserName,
		lastSeen: internalSession.LastSeen,
		isAdmin:  internalSession.IsAdmin,
		isLogged: internalSession.IsLogin,
		data:     internalSession.Data,
	}
}

func (sm *MemorySessionManager) ToBinary(session *Session) []byte {
	type InternalSession struct {
		Key      string
		UserName string
		LastSeen time.Time
		IsAdmin  bool
		IsLogin  bool
		Data     map[string]interface{}
	}
	internalSession := InternalSession{
		Key:      session.key,
		UserName: session.userName,
		LastSeen: session.lastSeen,
		IsAdmin:  session.isAdmin,
		IsLogin:  session.isLogged,
		Data:     session.data,
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(internalSession)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

package sessions

import (
	"sync"
	"time"
)

type Session struct {
	key      string
	userName string
	lastSeen time.Time
	isLogin  bool
	isAdmin  bool
	data     map[string]interface{}
	mutex    *sync.RWMutex
}

func (s *Session) GetID() string {
	return s.key
}

func (s *Session) UpdateLastSeen() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastSeen = time.Now()
}

func (s *Session) IsAdmin() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isAdmin
}

func (s *Session) GetUserName() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.userName
}

func (s *Session) SetAdmin(isAdmin bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.isAdmin = isAdmin
}

func (s *Session) SetUserName(userName string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.userName = userName
}

func (s *Session) IsLogin() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isLogin
}

func (s *Session) SetLogin(isLogin bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.isLogin = isLogin
}

func (s *Session) GetData(key string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}

func (s *Session) SetData(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = value
}

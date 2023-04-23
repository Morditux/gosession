package sessions

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type RedisSessionManager struct {
	client *redis.Client
}

func NewRedisSessionManager(client *redis.Client) *RedisSessionManager {
	return &RedisSessionManager{
		client: client,
	}
}

func (sm *RedisSessionManager) Exists(key string) bool {
	return sm.client.Exists(key).Val() == 1
}

func (sm *RedisSessionManager) GetSessionCount() int {
	return int(sm.client.DBSize().Val())
}

func (sm *RedisSessionManager) Add(session Session) {
	sm.client.Set(session.key, sm.ToBinary(session), 0)
}

func (sm *RedisSessionManager) CreateSession(userName string, isAdmin bool) (Session, bool) {
	session := Session{
		key:      uuid.New().String(),
		userName: userName,
		lastSeen: time.Now(),
		isAdmin:  false,
	}
	sm.Add(session)
	return session, true
}

func (sm *RedisSessionManager) Remove(key string) {
	sm.client.Del(key)
}

func (sm *RedisSessionManager) Get(key string) Session {
	if sm.Exists(key) {
		data, err := sm.client.Get(key).Bytes()
		if err != nil {
			return Session{}
		}
		return sm.FromBinary(data)
	}
	return Session{}
}

func (sm *RedisSessionManager) Clean() {
	sm.client.FlushDB()
}

func (sm *RedisSessionManager) FromBinary(data []byte) Session {
	var session Session
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(&session)
	return session
}

func (sm *RedisSessionManager) ToBinary(session Session) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(session)
	return buf.Bytes()
}

package sessions

type SessionManager interface {
	CreateSession(userName string, isAdmin bool) (*Session, bool)
	Get(key string) (*Session, error)
	GetSessionCount() int
	Clean()
	Add(session *Session)
	Remove(key string)
	Exists(key string) bool
	FromBinary(data []byte) *Session
	ToBinary(session *Session) []byte
}

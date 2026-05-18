package email

import (
	"sync"
	"time"
)

const verificationTTL = 5 * time.Minute

type verificationEntry struct {
	code      string
	expiredAt time.Time
}

var store = &verificationStore{m: sync.Map{}}

type verificationStore struct {
	m sync.Map
}

func (s *verificationStore) save(email, code string) {
	s.m.Store(email, verificationEntry{
		code:      code,
		expiredAt: time.Now().Add(verificationTTL),
	})
}

func (s *verificationStore) verify(email, code string) bool {
	val, ok := s.m.Load(email)
	if !ok {
		return false
	}
	entry := val.(verificationEntry)
	if time.Now().After(entry.expiredAt) {
		s.m.Delete(email)
		return false
	}
	if entry.code != code {
		return false
	}
	s.m.Delete(email)
	return true
}

func (s *verificationStore) exists(email string) bool {
	_, ok := s.m.Load(email)
	return ok
}

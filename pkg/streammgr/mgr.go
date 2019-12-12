package streammgr

import (
	"errors"
	"math/rand"
	"sync"
)

// Mgr control push address generate and auth,
// play url mapping
type Mgr struct {
	m_push map[string]string
	m_play map[string]string

	mutex *sync.Mutex
}

func NewStreamMgr() *Mgr {
	return &Mgr{
		m_push: make(map[string]string),
		m_play: make(map[string]string),
		mutex:  &sync.Mutex{},
	}
}

// GeneratePush return pushurl ,playurl ,error
func (s *Mgr) GeneratePush() (string, string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	pushurl := generate()
	playurl := generate()

	if _, exist := s.m_push[pushurl]; exist {
		return "", "", errors.New("push url exist")
	}

	if _, exist := s.m_play[playurl]; exist {
		return "", "", errors.New("play url exist")
	}

	s.m_push[pushurl] = playurl
	s.m_play[playurl] = pushurl

	// generate two random string
	// s.m.Store(key interface{}, value interface{})
	return pushurl, playurl, nil
}

// get real streaming url from play url
func (s *Mgr) GetStreamMapping(playurl string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if url, exist := s.m_play[playurl]; !exist {
		return "", errors.New("play url not exist")
	} else {
		return url, nil
	}
}

func (s *Mgr) AuthPushURL(url string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exist := s.m_push[url]; !exist {
		return errors.New("push url not exist")
	}
	return nil
}

var strs = `abcdefghijklmnopqrstuvwxyz`
var target_len = 26
var strslen = len(strs)

func generate() string {
	res := ""
	for i := 0; i < target_len; i++ {
		idx := rand.Int() % strslen
		res += strs[idx : idx+1]
	}
	return res
}

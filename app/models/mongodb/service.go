package mongodb

import (
	"gopkg.in/mgo.v2"
	"time"
)

type Service struct {
	baseSession *mgo.Session
	queue       chan int
	URL         string
	Open        int
}

var service Service

func (s *Service) New() error {
	var err error
	s.queue = make(chan int, MaxPool)
	for i := 0; i < MaxPool; i = i + 1 {
		s.queue <- 1
	}
	s.Open = 0
	s.baseSession, err = mgo.DialWithTimeout(s.URL, time.Duration(int64(DialTimeout))*time.Second)
	return err
}

func (s *Service) Session() *mgo.Session {
	<-s.queue
	s.Open++
	return s.baseSession.Copy()
}

func (s *Service) Close(c *Collection) {
	c.db.s.Close()
	s.queue <- 1
	s.Open--
}

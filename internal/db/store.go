package db

import "sync"

type Store struct {
	Mpp map[string]string
	Mu  sync.Mutex
}

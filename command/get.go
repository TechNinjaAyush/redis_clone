package command

import (
	"redis/internal/db"
)

func GetCommand(key string, s *db.Store) string {
	s.Mu.Lock()
	val := s.Mpp[key]
	defer s.Mu.Unlock()
	return val

}

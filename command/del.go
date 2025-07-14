package command

import "redis/internal/db"

func Deletecommand(keys []string, s *db.Store) int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	deletedCount := 0
	for _, key := range keys {
		if _, exists := s.Mpp[key]; exists {
			delete(s.Mpp, key)
			deletedCount++
		}
	}
	return deletedCount

}

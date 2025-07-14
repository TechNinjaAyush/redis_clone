package command

import (
	"fmt"
	"redis/internal/db"
	"strconv"
)

func INCRcommand(key string, s *db.Store) (int, error) {
	s.Mu.Lock()
	val, exists := s.Mpp[key]
	if !exists {

		s.Mpp[key] = "1"
		return 1, nil
	}
	defer s.Mu.Unlock()
	value, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("ERR value is not an integer or  out of range")
	}
	value++

	s.Mpp[key] = strconv.Itoa(value)

	return value, nil

}

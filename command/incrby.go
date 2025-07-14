package command

import (
	"fmt"
	"redis/internal/db"
	"strconv"
)

func INCRBYcommand(key, incr string, s *db.Store) (int, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Convert incr to int
	Incr, err := strconv.Atoi(incr)
	if err != nil {
		return 0, fmt.Errorf("ERR increment value is not an integer")
	}

	val, exists := s.Mpp[key]
	if !exists {
		s.Mpp[key] = incr
		return Incr, nil
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("ERR value is not an integer or out of range")
	}

	value += Incr
	s.Mpp[key] = strconv.Itoa(value)

	return value, nil
}

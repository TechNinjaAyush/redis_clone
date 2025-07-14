package command

import (
	"fmt"
	"redis/internal/db"
)

func SetCommand(key, value string, s *db.Store) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Mpp[key] = value
	fmt.Println("key is:", key, "and value is:", value)

}

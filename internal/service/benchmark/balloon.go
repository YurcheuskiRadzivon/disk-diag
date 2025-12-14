package benchmark

import (
	"fmt"
	"log"
	"math/rand"
)

func (s *service) Balloon() error {
	var balloon [][]byte
	physmem, err := s.physMemGet()
	if err != nil {
		return fmt.Errorf("physMemGet failed - %w", err)
	}

	for i := 0; i < int(physmem>>30); i++ {
		for j := 0; j < 1024; j++ {
			OneMiBBlock := make([]byte, Blocksize)
			lenRandom, err := rand.Read(OneMiBBlock)
			if err != nil {
				return fmt.Errorf("Balloon failed - %w", err)
			}
			if lenRandom != Blocksize {
				panic(fmt.Sprintf("lenRandom %d is not equal to Blocksize %d", lenRandom, Blocksize))
			}
			balloon = append(balloon, OneMiBBlock)
		}
		log.Printf("GiB #%d", i)
	}
	for i := 0; i < 10; i++ {
		x := rand.Intn(len(balloon))
		y := rand.Intn(Blocksize)
		if balloon[x][y] == 'b' {
			log.Print("Today is your lucky day: balloon[x][y] equals 'b'")
		}
	}

	return nil
}

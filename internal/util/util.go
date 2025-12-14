package util

import (
	"context"
	"fmt"

	"hotel/services"
	"math/big"

	"crypto/rand"
)

func GeneratePickUpCode(s *services.Services, HotelID uint) (string, error) {
	ma := big.NewInt(1000000)
	now, _ := rand.Int(rand.Reader, ma)
	code := fmt.Sprintf("%06d", now)
	ex, err := s.RdbRand.Exists(context.Background(), code).Result()
	if err != nil {
		return "", err
	}
	if ex == 1 {
		return GeneratePickUpCode(s, HotelID)

	} else {
		s.RdbRand.Set(context.Background(), fmt.Sprintf("%d:%s", HotelID, code), "1", 0)
		return code, nil
	}
}

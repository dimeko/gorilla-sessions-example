package core

import (
	"crypto/sha256"
	"fmt"
)

type Store struct {
	Rdbms *RDBMS
}

func NewStore() *Store {
	rdbms := NewRDBMS()
	store := &Store{
		Rdbms: rdbms,
	}
	return store
}

func (s *Store) ListProductsStore(limit string, offset string, filter string) ([]*Product, error) {
	var products []*Product
	products, err1 := s.Rdbms.ListProductsRDBMS(limit, offset, filter)
	if err1 != nil {
		log.Error(err1)
		return nil, err1
	}

	return products, nil
}

func (s *Store) TotalProductsStore() int {
	return s.Rdbms.TotalProductsRDBMS()
}
func (s *Store) Login(username string, password string) bool {
	_hash := sha256.New()
	_hash.Write([]byte(password))
	return s.Rdbms.UserExists(username, fmt.Sprintf("%x", _hash.Sum(nil)))
}

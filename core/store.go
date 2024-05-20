package core

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
		logger.Error(err1)
		return nil, err1
	}

	return products, nil
}

func (s *Store) TotalProductsStore() int {
	return s.Rdbms.TotalProductsRDBMS()
}
func (s *Store) Login(username string, password string) bool {
	return s.Rdbms.Login(username, password)
}

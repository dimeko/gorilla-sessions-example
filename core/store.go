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

// func (s *Store) CreateStore(payload interface{}) (*interface{}, error) {
// 	newUser, err := s.Rdbms.CreateRDBMS(payload)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return newUser, nil
// }

// func (s *Store) Update(id string, payload interface{}) (*interface{}, error) {
// 	updatedUser, err := UpdateRDBMS(id, payload)
// 	return updatedUser, nil
// }

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
	return s.Rdbms.UserExists(username, password)
}

// func (s *Store) GetStore(id string) (*interface{}, error) {
// 	unhashedCacheKey := "user:" + id
// 	var user *interface{}

// 	user, err = GetRDBMS(id)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	}

// 	err2 := json.Unmarshal(item.Value, &user)
// 	if err2 != nil {
// 		log.Error(err2)
// 		return nil, err2
// 	}

// 	return user, nil
// }

// func (s *Store) DeleteStore(id string) error {
// 	unhashedCacheKey := "user:" + id
// 	err := DeleteRDBMS(id)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	err1 := s.Cache.Remove(unhashedCacheKey)
// 	if err1 != nil {
// 		log.Error(err1)
// 	}

// 	return nil
// }

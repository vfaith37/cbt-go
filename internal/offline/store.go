package offline

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type OfflineStore struct {
	db *badger.DB
}

func NewOfflineStore(path string) (*OfflineStore, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable logging for production

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &OfflineStore{db: db}, nil
}

func (s *OfflineStore) StoreExam(exam *domain.Exam) error {
	return s.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(exam)
		if err != nil {
			return err
		}

		key := []byte(fmt.Sprintf("exam:%s", exam.ID))
		entry := badger.NewEntry(key, data).WithTTL(24 * time.Hour)
		return txn.SetEntry(entry)
	})
}

func (s *OfflineStore) GetExam(examID uuid.UUID) (*domain.Exam, error) {
	var exam domain.Exam

	err := s.db.View(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("exam:%s", examID))
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &exam)
		})
	})

	if err != nil {
		return nil, err
	}

	return &exam, nil
}

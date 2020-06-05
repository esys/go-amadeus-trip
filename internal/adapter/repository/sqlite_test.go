package repository

import (
	"amadeus-trip-parser/internal/domain/model"
	"database/sql"
	"github.com/google/uuid"
	"testing"
)

func getMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("cannot create memory sqlite database: %s", err)
		return nil
	}
	return db
}

func TestNewSQLiteTripRepo(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"no param",
			args{},
			true,
		},
		{
			"with param",
			args{getMemoryDB(t)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSQLiteTripRepo(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSQLiteTripRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_sqliteTripRepo_Create(t *testing.T) {
	type args struct {
		trip *model.Trip
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty trip",
			args{&model.Trip{
				ID: uuid.New().String(),
			}},
			false,
		},
	}
	s, _ := NewSQLiteTripRepo(getMemoryDB(t))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Create(tt.args.trip); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

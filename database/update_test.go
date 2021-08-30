package database

import (
	"database/sql"
	"testing"
)

func TestUpdateNoEntry(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}

	in := File{
		ID:           "i-am-not-a-real-id",
		Name:         "faux.png",
		Size:         6400,
		Status:       "remote",
		Replications: 5,
	}

	// Test getting an entry from an empty DB
	_, err = db.Update(in)
	if err == nil || err != sql.ErrNoRows {
		t.Error(err)
		t.Fatal("Test did not return a no rows error as expected")
	}
}

func TestUpdateEntry(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}

	in := File{
		ID:           "i-am-not-a-real-id",
		Name:         "faux.png",
		Size:         6400,
		Status:       "remote",
		Replications: 5,
	}

	// Add data to mock db
	err = db.Add(in)
	if err != nil {
		t.Fatal(err)
	}

	in.Replications = 0

	// Ask for data back from db
	old, err := db.Update(in)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the in and out file id's match
	if old.Replications != 5 {
		t.Fatalf("expected file replications to be %d but got %d instead", in.Replications, 5)
	}

	// Ask for data back from db
	out, err := db.Get(in.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the in and out file id's match
	if out.Replications != in.Replications {
		t.Fatalf("expected file replications to be %d but got %d instead", in.Replications, out.Replications)
	}

}

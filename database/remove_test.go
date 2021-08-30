package database

import (
	"database/sql"
	"testing"
)

func TestRemoveNoEntry(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}

	// Test getting an entry from an empty DB
	_, err = db.Remove("i-am-not-a-real-id")
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}
}

func TestRemoveEntry(t *testing.T) {
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

	// Ask for data back from db
	out, err := db.Remove(in.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the in and out file id's match
	if out.ID != in.ID {
		t.Fatalf("expected file id with %s but got %s instead", in.ID, out.ID)
	}
	// Check in and out file names match
	if out.Name != in.Name {
		t.Fatalf("expected file name with %s but got %s instead", in.Name, out.Name)
	}
	// Check in and out file sizes match
	if out.Size != in.Size {
		t.Fatalf("expected file size with %d but got %d instead", in.Size, out.Size)
	}
	// Check in and out status's match
	if out.Status != in.Status {
		t.Fatalf("expected file status with %s but got %s instead", in.Status, out.Status)
	}

	// Test that the entry was removed as expected.
	_, err = db.Get("i-am-not-a-real-id")
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}
}

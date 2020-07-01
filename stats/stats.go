package stats

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// NodeData is a struct of data that should be send to the arkstat server.
type NodeData struct {
	ID         string
	Username   string
	Email      string
	TotalSpace int
	UsedSpace  int
}

// CheckIn sends a JSON post request to the stats location provided.
func CheckIn(location string, input NodeData) (err error) {

	bytesRepresentation, err := json.Marshal(input)
	if err != nil {
		return err
	}

	resp, err := http.Post(location, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}

	var result NodeData

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result != input {
		return errors.New("different arkstat back from server")
	}

	return
}

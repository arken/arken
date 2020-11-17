package stats

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// NodeData is a struct of data that should be send to the arkstat server.
type NodeData struct {
	ID         string
	Username   string
	Email      string
	TotalSpace float64
	UsedSpace  float64
}

// CheckIn sends a JSON post request to the stats location provided.
func CheckIn(location string, input NodeData) (err error) {

	bytesRepresentation, err := json.Marshal(input)
	if err != nil {
		return err
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(location, "application/json", bytes.NewBuffer(bytesRepresentation))
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

	err = resp.Body.Close()
	client.CloseIdleConnections()
	return err
}

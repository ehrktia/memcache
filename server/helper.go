package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Data struct {
	Key   any `json:"key"`
	Value any `json:"value"`
}

func extractReqData(b *bytes.Buffer) (*Data, error) {
	d := new(Data)
	if err := json.Unmarshal(b.Bytes(), d); err != nil {
		return d, err
	}
	return d, nil

}

func writeResponse(res http.ResponseWriter, b []byte) error {
	if _, err := res.Write(b); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "ERROR:[%v]\n", err)
		return err
	}
	return nil

}

func reqGetMethod(req *http.Request) error {
	if req.Method != http.MethodGet {
		return fmt.Errorf("invalid req method supports only-%v", http.MethodGet)

	}
	return nil
}

func reqPostMethod(req *http.Request) error {
	if req.Method != http.MethodPost {
		return fmt.Errorf("invalid req method supports only-%v", http.MethodPost)

	}
	return nil
}

func readReqBody(
	buf *bytes.Buffer,
	req *http.Request, res http.ResponseWriter) error {
	if _, err := buf.ReadFrom(req.Body); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR-%v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return err

	}
	return nil
}

package git

import (
	"encoding/json"
	"fmt"
	"time"
)

type CommitDate time.Time

type Commit struct {
	Hash    string    `json:"hash"`
	Date    time.Time `json:"date"`
	Subject string    `json:"subject"`

	// TODO "message": "%s"
	// Message string `json:"message"`
}

func (c Commit) String() string {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("fail marshal: %v", err).Error()
	}
	return string(jsonData)
}

func (c Commit) AbbreviationHash() string {
	return string(c.Hash)[:7]
}

func (cd *CommitDate) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02T15:04:05Z"
	str := string(b)
	str = str[1 : len(str)-1]
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		return err
	}
	*cd = CommitDate(parsedTime)
	return nil
}

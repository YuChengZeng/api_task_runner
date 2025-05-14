package models

import "time"

type Record struct {
    Blockchain string `json:"blockchain"`
    Address    string `json:"address"`
}

type ProgressRecord struct {
    Blockchain   string                 `bson:"blockchain"`
    Address      string                 `bson:"address"`
    Status       string                 `bson:"status"`
    LastUpdated  time.Time              `bson:"last_updated"`
    Retries      int                    `bson:"retries"`
    ResponseData map[string]interface{} `bson:"response_data,omitempty"`
}
package api

import (
    "context"
    "net/http"
    "time"
    "encoding/json"
    "errors"
)

func MakeRequest(ctx context.Context, url string, headers map[string]string, params map[string]string) (map[string]interface{}, error) {
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    q := req.URL.Query()
    for k, v := range params {
        q.Add(k, v)
    }
    req.URL.RawQuery = q.Encode()

    for key, val := range headers {
        req.Header.Add(key, val)
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("API call failed with status: " + resp.Status)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    return result, nil
}
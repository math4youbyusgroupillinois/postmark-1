package postmark

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
)

const (
    Endpoint   = "https://api.postmarkapp.com/email"
    AuthHeader = "X-Postmark-Server-Token"
)

var (
    MissingOrIncorrectAPIKey = fmt.Errorf("postmark: Missing or incorrect API key header")
    InvalidRequest           = fmt.Errorf("postmark: Unprocessable Entity")
    ServerError              = fmt.Errorf("postmark: Server error")
    client                   http.Client
)

type Postmark struct {
    key string
}

func New(apikey string) *Postmark {
    return &Postmark{key: apikey}
}

func (p *Postmark) Send(m *Message) (*Response, error) {
    data, err := m.Marshal()
    if err != nil {
        return nil, err
    }
    postData := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", Endpoint, postData)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set(AuthHeader, p.key)

    rsp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    switch {
    case rsp.StatusCode == 401:
        return nil, MissingOrIncorrectAPIKey
    case rsp.StatusCode == 500:
        return nil, ServerError
    }

    var body bytes.Buffer
    _, err = io.Copy(&body, rsp.Body)
    rsp.Body.Close()
    if err != nil {
        return nil, err
    }

    prsp, err := UnmarshalResponse([]byte(body.String()))
    if err != nil {
        return nil, err
    }

    if rsp.StatusCode == 422 {
        return prsp, InvalidRequest
    }

    return prsp, nil
}

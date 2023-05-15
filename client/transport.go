package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "time"
)

type Transport struct {
    config *Config
}

func NewTransport(conf *Config) *Transport {
    return &Transport{config: conf}
}

// createRequest need for build http request from JSON
func (t *Transport) createRequest(body interface{}) (*http.Request, error) {
    var requestData []interface{}
    requestData = append(requestData, body)
    requestDataJSON, err := json.Marshal(requestData)
    out := bytes.Buffer{}
    json.Indent(&out, requestDataJSON, "", "\t")
    log.Printf("Request data:\n%s\n", out.String())

    if err != nil {
        return nil, fmt.Errorf("%v: %v", "can not serialize request", err)
    }
    encryptedRequestBody, err := t.config.encrypt(requestDataJSON)
    if err != nil {
        return nil, fmt.Errorf("%v: %v", "can not ecript request", err)
    }

    date := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
    req, err := http.NewRequest(method, t.config.URL, bytes.NewReader([]byte(encryptedRequestBody)))
    if err != nil {
        return nil, fmt.Errorf("%v: %v", "can not create request", err)
    }

    req.Header.Add(headerContentType, valueContentType)
    req.Header.Add(headerAccept, valueAccept)
    req.Header.Add(headerAuthentication, t.config.generateSign(date))
    req.Header.Add(headerClientIdentifier, valueClientIdentifier)
    req.Header.Add(headerUserIdentifier, t.config.UserIdentifier)
    req.Header.Add(headerDate, date)
    req.Header.Add(headerAccessKey, t.config.AccessKey)
    req.Header.Add(headerSdkVersion, valueSdkVersion)
    req.Header.Add(headerAsync, valueAsync)

    req.Close = true

    return req, nil
}

// Do executes action on remote agent
func (t *Transport) Do(body interface{}) (*M3BatchResult, error) {
    req, err := t.createRequest(body)
    if err != nil {
        return nil, err
    }
    log.Println(req)
    resp, err := t.config.Client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("%v: %v", "failed to process request", err)
    }

    defer func() {
        if rerr := resp.Body.Close(); err == nil {
            err = rerr
        }
    }()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("%v: %v", "failed to process response", err)
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("got status code %v instead of 200\nBody: %s",
            resp.StatusCode, string(respBody))
    }

    decryptedResponse, err := t.config.decrypt(respBody)
    if err != nil {
        return nil, err
    }

    out := bytes.Buffer{}
    json.Indent(&out, []byte(decryptedResponse), "", "\t")
    log.Printf("Response data:\n%s\n", out.String())

    result := new(M3BatchResult)
    err = json.Unmarshal([]byte(decryptedResponse), &result)
    if err != nil {
        return nil, fmt.Errorf("%v: %v", "failed to unmarshal response", err)
    }
    return result, nil
}

// MakePayload collect params to DefaultPayload
func (t *Transport) MakePayload(params interface{}, method string) (*DefaultPayload, error) {
    paramsJSON, err := json.Marshal(params)
    if err != nil {
        return nil, err
    }
    payload := &DefaultPayload{
        ID:   generateUUID(),
        Type: method,
        Params: &Params{
            Body: string(paramsJSON),
        },
    }
    return payload, nil
}

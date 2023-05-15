package client

// DEFAULT

//Params contains body with request in json format
type Params struct {
    Body string `json:"body"`
}

//DefaultPayload contains action id and action type for json rpc and params
type DefaultPayload struct {
    ID     string  `json:"id"`
    Type   string  `json:"type"`
    Params *Params `json:"params"`
}

//M3RawResult contains fields that will be returned by server in a response for single request
type M3RawResult struct {
    ID         string `json:"id"`
    Status     string `json:"status"`
    Error      string `json:"error,omitempty"`
    Data       string `json:"data"`
    StatusCode int    `json:"statusCode"`
}

//M3BatchResult contains list of raw results
type M3BatchResult struct {
    Results []*M3RawResult `json:"results"`
}

package client

//go:generate go install github.com/golang/mock/mockgen@v1.6.0
//go:generate mockgen -source ./client.go -destination ./mock/client_mock.go -package cmock

const (
    method = "POST"

    valueContentType      = "application/json"
    valueAccept           = "application/json"
    valueClientIdentifier = "terraform-provider"
    valueSdkVersion       = "3.2.80"
    valueAsync            = "false"

    headerContentType      = "Content-type"
    headerAccept           = "Accept"
    headerAuthentication   = "maestro-authentication"
    headerClientIdentifier = "maestro-request-identifier"
    headerDate             = "maestro-date"
    headerAccessKey        = "maestro-accesskey"
    headerSdkVersion       = "maestro-sdk-version"
    headerAsync            = "maestro-sdk-async"

    defaultHeaderEventGroup = "priv_cloud_action"
    defaultZippedHeader     = "zipped"
)

// Transporter interface witch provide transports possibility
type Transporter interface {
    Do(body interface{}) (*M3BatchResult, error)
    MakePayload(interface{}, string) (*DefaultPayload, error)
}

// Client collect methods and params for working with m3API
type Client struct {
    Transporter
}

func NewClient(conf *Config) *Client {
    return &Client{
        Transporter: NewTransport(conf),
    }
}

package client

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "crypto/tls"
    "encoding/base64"
    "errors"
    "github.com/google/uuid"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

// Config is configuration struct for agent, also provide encrypting methods
type Config struct {
    Client *http.Client

    URL            string `json:"url"`
    AccessKey      string `json:"access_key"`
    SecretKey      string `json:"secret_key"`
    UserIdentifier string `json:"user_identifier"`
    TenantName     string `json:"tenant_name"`
    RegionName     string `json:"regionName"`
    Cloud          string `json:"cloud"`
}

func NewConfig(url, userIdentifier, accessKey, secretKey, tenantName, regionName, cloud string) *Config {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    if len(url) > 0 {
        return &Config{
            Client: &http.Client{
                Transport: tr,
                Timeout:   time.Minute * 2,
            },
            URL:            url,
            UserIdentifier: userIdentifier,
            AccessKey:      accessKey,
            SecretKey:      secretKey,
            TenantName:     tenantName,
            RegionName:     regionName,
            Cloud:          cloud,
        }
    }
    url, exists := os.LookupEnv("TERRAFORM_M3_URL")
    if !exists {
        log.Fatalf("ENV VARIABLE 'TERRAFORM_M3_URL' IS NOT SET")
    }
    userIdentifier, exists = os.LookupEnv("TERRAFORM_M3_USER_IDENTIFIER")
    if !exists {
        log.Fatalf("ENV VARIABLE 'TERRAFORM_M3_USER_IDENTIFIER' IS NOT SET")
    }
    accessKey, exists = os.LookupEnv("TERRAFORM_M3_ACCESS_KEY")
    if !exists {
        log.Fatalf("ENV VARIABLE 'TERRAFORM_M3_ACCESS_KEY' IS NOT SET")
    }
    secretKey, exists = os.LookupEnv("TERRAFORM_M3_SECRET_KEY")
    if !exists {
        log.Fatalf("ENV VARIABLE 'TERRAFORM_M3_SECRET_KEY' IS NOT SET")
    }
    tenantName, exists = os.LookupEnv("TERRAFORM_M3_TENANT_NAME")
    if !exists {
        tenantName = ""
    }
    regionName, exists = os.LookupEnv("TERRAFORM_M3_REGION_NAME")
    if !exists {
        regionName = ""
    }
    cloud, exists = os.LookupEnv("TERRAFORM_M3_CLOUD")
    if !exists {
        cloud = ""
    }

    return &Config{
        Client:         &http.Client{Transport: tr},
        URL:            url,
        UserIdentifier: userIdentifier,
        AccessKey:      accessKey,
        SecretKey:      secretKey,
        TenantName:     tenantName,
        RegionName:     regionName,
        Cloud:          cloud,
    }
}

func generateUUID() string {
    id, err := uuid.NewUUID()
    if err != nil {
        log.Printf("Failed to generate UUID: %v", err)
        return ""
    }
    return id.String()
}

func (c *Config) encrypt(requestDataJSON []byte) (string, error) {
    block, err := aes.NewCipher([]byte(c.SecretKey))
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    iv := make([]byte, 12)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, iv, requestDataJSON, nil)
    encryptedData := append(iv, ciphertext...)

    return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func (c *Config) decrypt(response []byte) (string, error) {
    data, err := base64.StdEncoding.DecodeString(string(response))
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher([]byte(c.SecretKey))
    if err != nil {
        return "", err
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    if len(data) < aesgcm.NonceSize() {
        return "", errors.New("ciphertext too short")
    }
    nonce := data[:aesgcm.NonceSize()]
    ciphertext := data[aesgcm.NonceSize():]

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}

func (c *Config) generateSign(date string) string {
    mac := hmac.New(sha256.New, []byte(c.SecretKey+date))
    message := "M3-POST:" + c.AccessKey + ":" + date
    mac.Write([]byte(message))

    var stringBuilder strings.Builder

    for _, element := range mac.Sum(nil) {
        stringBuilder.WriteString(strconv.FormatInt(((int64(element) & 0xff) + 0x100), 16))
    }
    return stringBuilder.String()
}

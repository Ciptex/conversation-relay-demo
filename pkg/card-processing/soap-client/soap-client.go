package soapclient

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SOAPClient represents a client for making SOAP requests
type SOAPClient struct {
	endpoint string
	username string
	password string
	wsdl     string
	client   *http.Client
	headers  map[string]string
}
type Namespace struct {
	Prefix string
	URI    string
}

// NewSOAPClient creates a new SOAP client with basic authentication
func NewSOAPClient(endpoint, username, password string) *SOAPClient {
	return &SOAPClient{
		endpoint: endpoint,
		username: username,
		password: password,
		client:   &http.Client{},
		headers:  make(map[string]string),
	}
}

// AddHeader adds a custom HTTP header to the SOAP request
func (c *SOAPClient) AddHeader(key, value string) {
	c.headers[key] = value
}

// SOAPResponse represents a generic SOAP response with dynamic operation response
type SOAPResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	Content []byte `xml:",innerxml"`
}

// Call sends a SOAP request with the provided XML envelope
func (c *SOAPClient) Call(action string, xmlEnvelope string, response interface{}) (string, error) {
	// Create a new request
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBufferString(xmlEnvelope))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Add necessary headers
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	if action != "" {
		req.Header.Add("SOAPAction", action)
	}

	// Add basic authentication
	if c.username != "" && c.password != "" {
		auth := c.username + ":" + c.password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Add("Authorization", "Basic "+encodedAuth)
	}

	// Add any custom headers
	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received non-OK status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var envelope SOAPResponse
	err = xml.Unmarshal([]byte(body), &envelope)
	if err != nil {
		return "", err
	}
	err = xml.Unmarshal(envelope.Body.Content, response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}
	return string(envelope.Body.Content), nil
}

func (c *SOAPClient) BuildEnvelope(bodyContent string, namespaces []Namespace) string {
	// Default SOAP namespace if not provided
	hasSOAPNS := false
	for _, ns := range namespaces {
		if ns.Prefix == "soap" {
			hasSOAPNS = true
			break
		}
	}

	if !hasSOAPNS {
		// Add default SOAP namespace if not specified
		namespaces = append(namespaces, Namespace{
			Prefix: "soap",
			URI:    "http://schemas.xmlsoap.org/soap/envelope/",
		})
	}

	// Build namespace declarations
	nsDeclarations := ""
	for _, ns := range namespaces {
		nsDeclarations += fmt.Sprintf(" xmlns:%s=\"%s\"", ns.Prefix, ns.URI)
	}

	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope%s>
	<soap:Header></soap:Header>
	<soap:Body>
%s
	</soap:Body>
</soap:Envelope>`, nsDeclarations, strings.TrimSpace(bodyContent))

	return envelope
}

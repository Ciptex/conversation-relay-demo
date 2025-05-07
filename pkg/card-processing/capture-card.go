package cardprocessing

import (
	soapclient "conversation-relay/pkg/card-processing/soap-client"
	"conversation-relay/pkg/repo"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
)

type CapturePANResponse struct {
	XMLName     xml.Name `xml:"capturePANResponse"`
	ResultCode  string   `xml:"resultCode"`
	ResultHR    string   `xml:"resultHR"`
	BenignPAN   string   `xml:"benignPAN"`
	Pid         string   `xml:"pid,omitempty"`
	Token       string   `xml:"token,omitempty"`
	IssuerName  string   `xml:"issuerName,omitempty"`
	Brand       string   `xml:"brand,omitempty"`
	ProductType string   `xml:"productType,omitempty"`
	Level       string   `xml:"level,omitempty"`
	IsoCountry  string   `xml:"isoCountry,omitempty"`
	LuhnToken   string   `xml:"luhnToken,omitempty"`
	// Include the namespace in the struct definition
	Xmlns string `xml:"xmlns,attr"`
}

type CaptureCV2Response struct {
	XMLName    xml.Name `xml:"captureCV2Response"`
	ResultCode string   `xml:"resultCode,omitempty"`
	ResultHR   string   `xml:"resultHR,omitempty"`
	Xmlns      string   `xml:"xmlns,attr"`
}

func createSoapClient() (*soapclient.SOAPClient, []soapclient.Namespace) {
	namespaces := []soapclient.Namespace{
		{
			Prefix: "tns",
			URI:    "https://soap.syntec.co.uk/",
		},
	}
	return soapclient.NewSOAPClient(
		"https://eu-w2.ceclients.cardeasy.com/soap/cardeasy?wsdl", // SOAP endpoint URL
		os.Getenv("CARD_EASY_UNAME"),                              // Basic auth username
		os.Getenv("CARD_EASY_PWD"),                                // Basic auth password
	), namespaces
}

func CaptureCard(callSid, epid string) (string, error) {
	slog.Info("CaptureCard", "callSid", callSid, "epid", epid)
	client, namespaces := createSoapClient()
	client.AddHeader("Cache-Control", "no-cache")

	bodyContent := `<tns:capturePAN>
	<Customer>ciptex</Customer>
	<epid>` + epid + `</epid>
	<capMethod>DTMF</capMethod>
	<beepFlag>BEEP</beepFlag>
</tns:capturePAN>`
	soapEnvelope := client.BuildEnvelope(bodyContent, namespaces)
	capturePANResponse := new(CapturePANResponse)
	res, err := client.Call(`""`, soapEnvelope, capturePANResponse)
	if err != nil {
		slog.Error("Error capturing card", "error", err)
		return "", err
	}
	_ = res

	slog.Info("capture pan res", "resultCode", capturePANResponse.ResultCode, "resultHR", capturePANResponse.ResultHR, "benignPAN", capturePANResponse.BenignPAN)
	if capturePANResponse.Token == "" {
		slog.Error("Error capturing card", "error", "Token is empty")
		return "", fmt.Errorf("Token is empty")
	}
	pm := repo.GetGloabalRepo().GetPaymentMeta(callSid)
	pm.Pid = capturePANResponse.Pid
	pm.Token = capturePANResponse.Token
	pm.LuhnToken = capturePANResponse.LuhnToken
	repo.GetGloabalRepo().SetPaymentMeta(callSid, pm)
	return "successfull", nil
}

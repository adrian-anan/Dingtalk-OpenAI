package ops

import (
	"crypto/tls"
	"net"
	"regexp"
	"strings"
	"time"
)

// 域名信息
type DomainMsg struct {
	CreateDate string `json:"create_date"`
	ExpiryDate string `json:"expiry_date"`
	Registrar  string `json:"registrar"`
}

// GetDomainMsg 获取域名信息
func GetDomainMsg(domain string) (dm DomainMsg, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", "whois.verisign-grs.com:43")
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(domain + "\r\n"))
	if err != nil {
		return
	}
	buf := make([]byte, 1024)
	var num int
	num, err = conn.Read(buf)
	if err != nil {
		return
	}
	response := string(buf[:num])
	re := regexp.MustCompile(`Creation Date: (.*)\n.*Expiry Date: (.*)\n.*Registrar: (.*)`)
	match := re.FindStringSubmatch(response)
	if len(match) > 3 {
		dm.CreateDate = strings.TrimSpace(strings.Split(match[1], "Creation Date:")[0])
		dm.ExpiryDate = strings.TrimSpace(strings.Split(match[2], "Expiry Date:")[0])
		dm.Registrar = strings.TrimSpace(strings.Split(match[3], "Registrar:")[0])
	}
	return
}

// GetDomainCertMsg 获取域名证书信息
func GetDomainCertMsg(domain string) (cm tls.ConnectionState, err error) {
	var conn net.Conn
	conn, err = net.DialTimeout("tcp", domain+":443", time.Second*10)
	if err != nil {
		return
	}
	defer conn.Close()
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: domain,
	})
	defer tlsConn.Close()
	err = tlsConn.Handshake()
	if err != nil {
		return
	}
	cm = tlsConn.ConnectionState()
	return
}
func search(question string) (string, error) {
	url := "https://aichat.adriantech.uk/api/openapi/v1/chat/completions"
	headers := map[string]string{
		"Authorization": "Bearer fastgpt-apd5320l00ojv21yp8h49fok-64e6b185f1124b2fc0829976",
		"User-Agent":    "Apifox/1.0.0 (https://www.apifox.cn)",
		"Content-Type":  "application/json",
	}

	data := requestData{
		ChatID: "88888888",
		Stream: false,
		Detail: false,
		Messages: []userInput{
			{
				Content: question,
				Role:    "user",
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Custom TLS configuration (not recommended for production)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)

		content := ""
		if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
			if firstChoice, ok := choices[0].(map[string]interface{}); ok {
				if message, ok := firstChoice["message"].(map[string]interface{}); ok {
					content = message["content"].(string)
				}
			}
		}

		return content, nil
	}

	return "", fmt.Errorf("Received %d HTTP status code", resp.StatusCode)
}

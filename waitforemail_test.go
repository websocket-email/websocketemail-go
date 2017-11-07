package nomockemail

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"testing"
	"time"
)

func GetTestToken(t *testing.T) string {
	tok := os.Getenv("NOMOCKEMAIL_TEST_TOKEN")
	if tok == "" {
		t.Fatal("please set NOMOCKEMAIL_TEST_TOKEN env variable")
	}
	return tok
}

type TestEmailCreds struct {
	User     string
	Password string
	Server   string
}

func GetTestEmailCredentials(t *testing.T) TestEmailCreds {
	credsJson := os.Getenv("NOMOCKEMAIL_TEST_EMAIL_CREDS")
	if credsJson == "" {
		format, _ := json.Marshal(TestEmailCreds{})
		t.Fatalf("please set NOMOCKEMAIL_TEST_EMAIL_CREDS env variable. Format:  %s", string(format))
	}

	creds := TestEmailCreds{}
	err := json.Unmarshal([]byte(credsJson), &creds)
	if err != nil {
		t.Fatal(err)
	}

	return creds
}

func SendEmail(t *testing.T, to, subject, content string) {
	creds := GetTestEmailCredentials(t)

	portSepIdx := strings.LastIndex(creds.Server, ":")
	if portSepIdx == -1 {
		t.Fatal("'Server' must contain a port")
	}

	c, err := smtp.Dial(creds.Server)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	host := creds.Server[0:portSepIdx]
	config := &tls.Config{ServerName: host}
	err = c.StartTLS(config)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Auth(smtp.PlainAuth("", creds.User, creds.Password, host))
	if err != nil {
		t.Fatal(err)
	}
	err = c.Mail(creds.User)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Rcpt(to)
	if err != nil {
		t.Fatal(err)
	}
	w, err := c.Data()
	if err != nil {
		t.Fatal(err)
	}

	mail := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, content)
	_, err = w.Write([]byte(mail))
	if err != nil {
		w.Close()
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Quit()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWaitForEmail(t *testing.T) {
	tok := GetTestToken(t)
	addr := MustGenerateEmailAddress()

	ch, cleanup, err := WaitForEmail(tok, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	SendEmail(t, addr, "hi", "hello world")

	select {
	case email, ok := <-ch:
		if !ok {
			t.Fatal("unable to wait for email!")
		}
		if email.To != addr {
			t.Fatalf("bad email to: %#v", email)
		}
		if email.Subject != "hi" {
			t.Fatalf("bad email subject: %#v", email)
		}
		if strings.Trim(email.Body, " \n") != "hello world" {
			t.Fatalf("bad email body: %#v", email)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("email was not recieved after 20 seconds!")
	}
}

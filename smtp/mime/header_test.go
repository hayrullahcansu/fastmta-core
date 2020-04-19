package mime

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestParseHeaders(t *testing.T) {
	tests := map[string]struct {
		input  string
		output HeaderCollection
		err    error
	}{
		"successful conversion 2": {
			input: "Received: by vmta1.localhost;\r\nTo: hayrullah.cansu@gmail.com\r\nFrom: test@test.com\r\nSubject: Test from SMTP Diag ToolDate: Sat, 26 Jan 2019 12:01:37 +0300\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=\"boundaryHtOrkw==\"This is a message encoded in MIME format.\r\n--boundaryHtOrkw==Content-Type: text/plainContent-Transfer-Encoding: quoted-printableThis is a test from SMTP Diag Tool--boundaryHtOrkw==Content-Type: text/htmlContent-Transfer-Encoding: quoted-printableThis is a test from SMTP Diag Tool--boundaryHtOrkw==--",
			output: HeaderCollection{M: map[string]string{
				"Received":     "by vmta1.localhost;",
				"To":           "hayrullah.cansu@gmail.com",
				"From":         "test@test.com",
				"Subject":      "Test from SMTP Diag ToolDate: Sat, 26 Jan 2019 12:01:37 +0300",
				"MIME-Version": "1.0",
				"Content-Type": "multipart/alternative; boundary=\"boundaryHtOrkw==\"This is a message encoded in MIME format.",
			}},
			err: nil,
		},
		"successful conversion": {
			input: "Message-ID: <8d5>\r\n" +
				"To: hayrullah.cansu@gmail.com\r\n" +
				"From: posta@posta.jetmail.com.tr\r\n" +
				"Subject: Test from SMTP Diag Tool\r\n" +
				"MIME-Version: 1.0\r\n\r\ndatadatadata",
			output: HeaderCollection{M: map[string]string{
				"Message-ID":   "<8d5>",
				"To":           "hayrullah.cansu@gmail.com",
				"From":         "posta@posta.jetmail.com.tr",
				"Subject":      "Test from SMTP Diag Tool",
				"MIME-Version": "1.0"}},
			err: nil,
		},
	}

	/*
		M: {
						"Mesage-ID:":    "<8d5>",
						"To:":           "hayrullah.cansu@gmail.com",
						"From:":         "posta@posta.jetmail.com.tr",
						"Subject:":      "Test from SMTP Diag Tool",
						"MIME-Version:": "1.0"},
	*/

	for _, test := range tests {
		a, b := test.output, ParseHeaders(test.input)
		if !assert.EqualValues(t, a.M, (*b).M) {
			t.Errorf("Parsing header not works expected")
		}

		if !reflect.DeepEqual(a.M, (*b).M) {
			t.Errorf("Parsing header not works expected")
		}
		if diff := deep.Equal(a.M, (*b).M); diff != nil {
			t.Error(diff)
		}
	}
}

// func checkDifferent(expected *message.HeaderCollection, output *message.HeaderCollection) (r bool) {
// 	r = true
// 	for key, val := range *expected {
// 		if a, b := (*output).GetData(key); !b || a != val {
// 			r = false
// 			return
// 		}
// 	}
// 	return
// }

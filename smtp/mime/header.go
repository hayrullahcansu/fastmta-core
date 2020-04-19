package mime

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hayrullahcansu/fastmta-core/cross"
	"github.com/hayrullahcansu/fastmta-core/logger"
	"github.com/hayrullahcansu/fastmta-core/util"
)

const (
	_BoundaryTokenPattern = "boundary=\"(.+)\""
	_BoundaryParsePattern = "--%s"
)

//ParseHeaders provides to get headers data from mime message
func ParseHeaders(data string) *HeaderCollection {
	headers := NewHeaderCollection()
	headerSection := GetHeaderSection(data, true)
	reader := bufio.NewReader(strings.NewReader(headerSection))
	for {
		line, _, err := reader.ReadLine()
		if util.IsNullOrEmptyB(line) {
			if err == io.EOF {
				break
			}
			continue
		}
		if err == io.EOF {
			colonIndex := bytes.Index(line, []byte(":"))
			if colonIndex >= 0 {
				key := string(bytes.TrimSpace(line[0:colonIndex]))
				value := string(bytes.TrimSpace(line[colonIndex+1:]))
				headers.Add(key, value)
			}
			break
		} else {
			colonIndex := bytes.Index(line, []byte(":"))
			if colonIndex >= 0 {
				key := string(bytes.TrimSpace(line[0:colonIndex]))
				value := string(bytes.TrimSpace(line[colonIndex+1:]))
				headers.Add(key, value)
			}
		}

	}
	return headers
}

//GetHeaderSection provides to get header section from mime message
func GetHeaderSection(data string, unfold bool) string {
	logger.Instance().Info(data)

	//try detect header section smtp standarts as ref -> rfc
	endOfHeaderIndex := strings.Index(data, cross.CLRF+cross.CLRF)
	if endOfHeaderIndex < 0 {
		//we cant parse rfc
		//Lets try to get using regexp
		r := getHeaderParserRegex()
		o := r.FindSubmatch([]byte(data))
		if o != nil {
			boundary := string(o[1])
			endOfHeaderIndex = strings.Index(data, fmt.Sprintf(_BoundaryParsePattern, boundary))
		}
	}

	//if we detect header section...
	if endOfHeaderIndex >= 0 {
		// logger.Infof("header--> %s\n", data[0:endOfHeaderIndex])
		// logger.Infof("header--> %s\n", string(data[0:endOfHeaderIndex]))
		// logger.Infof("header--> %s\n", string([]rune(data)[0:endOfHeaderIndex]))
		if unfold {
			return unfoldHeaders(data[0:endOfHeaderIndex])
		}
		return data[0:endOfHeaderIndex]
	}

	//there are no headers so return an empty string
	return ""
}

func unfoldHeaders(data string) string {
	sb := &strings.Builder{}
	reader := bufio.NewReader(strings.NewReader(data))
	lines, isPrefix, err := reader.ReadLine()
	for !util.IsNullOrEmptyB(lines) && err != io.EOF && !isPrefix {
		if !util.IsNullOrEmptyB(lines) {
			sb.Write([]byte(cross.CLRF))
		}
		sb.Write(lines)
		lines, isPrefix, err = reader.ReadLine()
	}
	_r := sb.String()
	return _r
}

func getHeaderParserRegex() *regexp.Regexp {
	return regexp.MustCompile(_BoundaryTokenPattern)
}

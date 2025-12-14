package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/pkg/errors"
	"github.com/xinliangnote/go-gin-api/pkg/timeutil"
)

func (s *signature) Verify(authorization, date string, path string, method string, params url.Values) (ok bool, err error) {
	if date == "" {
		err = errors.New("date required")
		return
	}

	if path == "" {
		err = errors.New("path required")
		return
	}

	if method == "" {
		err = errors.New("method required")
		return
	}

	methodName := strings.ToUpper(method)
	if !methods[methodName] {
		err = errors.New("method param error")
		return
	}

	ts, err := timeutil.ParseCSTInLocation(date)
	if err != nil {
		err = errors.New("date must follow '2006-01-02 15:04:05'")
		return
	}

	if timeutil.SubInLocation(ts) > float64(s.ttl/time.Second) {
		err = errors.Errorf("date exceeds limit %v", s.ttl)
		return
	}

	// Encode() 方法中自带 sorted by key
	// params.Encode()：
	// params := url.Values{
	// 	"z": []string{"3"},
	// 	"a": []string{"1"},
	// 	"b": []string{"2"},
	// }
	// params.Encode() 返回 "a=1&b=2&z=3"（按字母顺序排序）

	// QueryUnescape: 将编码后的查询字符串解码回原始格式,去除 URL 编码的特殊字符（如 %20 → 空格）; 如 "name=John%20Doe&age=30" -> "name=John Doe&age=30"
	sortParamsEncode, err := url.QueryUnescape(params.Encode())
	if err != nil {
		err = errors.Errorf("url QueryUnescape error %v", err)
		return
	}

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(path)
	buffer.WriteString(delimiter)
	buffer.WriteString(methodName)
	buffer.WriteString(delimiter)
	buffer.WriteString(sortParamsEncode)
	buffer.WriteString(delimiter)
	buffer.WriteString(date)

	// 对数据进行 hmac 加密，并进行 base64 encode
	hash := hmac.New(sha256.New, []byte(s.secret))
	hash.Write(buffer.Bytes())
	digest := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	ok = authorization == fmt.Sprintf("%s %s", s.key, digest)
	return
}

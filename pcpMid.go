package gopcp_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	KLog "github.com/lock-free/goklog"
	"github.com/lock-free/gopcp"
	"net/http"
	"net/url"
)

var klog = KLog.GetInstance()

type PcpHttpResponse struct {
	Data   interface{} `json:"text"`
	Errno  int         `json:"errno"`
	ErrMsg string      `json:"errMsg"`
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func ResponseToBytes(pcpHttpRes PcpHttpResponse) []byte {
	bytes, jerr := JSONMarshal(pcpHttpRes)

	if jerr != nil {
		ret, _ := json.Marshal(ErrorToResponse(jerr, 530))
		return ret
	} else {
		return bytes
	}
}

type MidFunType = func(http.ResponseWriter, *http.Request, interface{}) (interface{}, error)

func ErrorToResponse(err error, code int) PcpHttpResponse {
	return PcpHttpResponse{nil, code, err.Error()}
}

func GetPcpMid(sandbox *gopcp.Sandbox) MidFunType {
	pcpServer := gopcp.NewPcpServer(sandbox)

	return func(w http.ResponseWriter, r *http.Request, attachment interface{}) (interface{}, error) {
		var pcpHttpRes PcpHttpResponse
		var arr interface{}
		var err error = nil
		var rawQuery string
		var ret interface{}

		// recover if panic
		defer func() {
			if r := recover(); r != nil {
				panicErr := fmt.Errorf("panic happened at pcp mid, return is %v", r)
				klog.LogError("pcp-mid-exception", panicErr)
				w.Write(ResponseToBytes(ErrorToResponse(err, 530)))
			}
		}()

		if r.Method == "GET" {
			rawQuery, err = url.QueryUnescape(r.URL.RawQuery)
			if err == nil {
				// parse url query
				err = json.Unmarshal([]byte(rawQuery), &arr)
			}
		} else {
			// get post body
			arr, err = GetJsonBody(r)
		}

		if err != nil {
			pcpHttpRes = ErrorToResponse(err, 530)
		} else {
			ret, err = pcpServer.ExecuteJsonObj(arr, attachment)

			if err != nil {
				pcpHttpRes = ErrorToResponse(err, 530)
			} else {
				pcpHttpRes = PcpHttpResponse{ret, 0, ""}
			}
		}

		w.Write(ResponseToBytes(pcpHttpRes))

		if err != nil {
			klog.LogError("pcp-mid-exception", err)
		}

		return arr, err
	}
}

package gopcp_service

import (
	"encoding/json"
	"github.com/idata-shopee/gopcp"
	"net/http"
)

type PcpHttpResponse struct {
	Data   interface{} `json:"text"`
	Errno  int         `json:"errno"`
	ErrMsg string      `json:"errMsg"`
}

func executeRequest(pcpServer *gopcp.PcpServer, arr interface{}, attachment interface{}) PcpHttpResponse {
	ret, err := pcpServer.ExecuteJsonObj(arr, attachment)
	if err != nil {
		return PcpHttpResponse{nil, 530, err.Error()}
	}
	return PcpHttpResponse{ret, 0, ""}
}

func toBytes(pcpHttpRes PcpHttpResponse) []byte {
	bytes, jerr := json.Marshal(pcpHttpRes)

	if jerr != nil {
		ret, _ := json.Marshal(PcpHttpResponse{nil, 530, jerr.Error()})
		return ret
	} else {
		return bytes
	}
}

type MidFunType = func(http.ResponseWriter, *http.Request, interface{})

func GetPcpMid(sandbox *gopcp.Sandbox) MidFunType {
	pcpServer := gopcp.NewPcpServer(sandbox)

	return func(w http.ResponseWriter, r *http.Request, attachment interface{}) {
		var pcpHttpRes PcpHttpResponse
		var arr interface{}
		var err error

		if r.Method == "GET" {
			// parse url query
			err = json.Unmarshal([]byte(r.URL.RawQuery), &arr)
		} else {
			// get post body
			arr, err = GetJsonBody(r)
		}

		if err != nil {
			pcpHttpRes = PcpHttpResponse{nil, 530, err.Error()}
		} else {
			pcpHttpRes = executeRequest(pcpServer, arr, attachment)
		}

		w.Write(toBytes(pcpHttpRes))
	}
}

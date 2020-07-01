package route

import (
	"bs/util"
	"encoding/json"

	"github.com/szxby/tools/log"
)

type rawPack struct {
	Sign string
	Data string
}

func checkSignature(msg []byte) bool {
	pkg := rawPack{}
	if err := json.Unmarshal(msg, &pkg); err != nil {
		log.Error("umarshal msg fail %v", err)
		return false
	}
	sign := pkg.Sign
	data := pkg.Data
	// log.Debug("signData:%v", data)
	// log.Debug("sign:%v", signature(data))
	if util.CalculateHash(data) != sign {
		return false
	}
	return true
}

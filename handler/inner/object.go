package inner

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/arkor/utils"
)

func AllocateFileID(ctx *macaron.Context, log *logrus.Logger) (int, []byte) {
	m := make(map[string]string)
	m["file_id"] = utils.MD5ID()

	result, err := json.Marshal(m)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	return http.StatusOK, result
}

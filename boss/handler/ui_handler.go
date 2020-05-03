package handler

import (
	"net/http"

	"github.com/hayrullahcansu/fastmta-core/boss/ui"
	"github.com/hayrullahcansu/fastmta-core/util"
	"github.com/sirupsen/logrus"
)

//JoinUI handles login requests and authorize broker which is valid
func JoinUI(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("A client joint UI\n%s", util.FormatRequest(r))
	c := ui.NewClient()
	c.ServeWs(w, r)
	ui.Manager().ConnectLobby(c)
}

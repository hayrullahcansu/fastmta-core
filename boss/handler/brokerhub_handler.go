package handler

import (
	"net/http"

	"github.com/hayrullahcansu/fastmta-core/boss/brokerhub"
	"github.com/hayrullahcansu/fastmta-core/util"
	"github.com/sirupsen/logrus"
)

//JoinLobby handles login requests and authorize broker which is valid
func JoinLobby(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("A client joint Lobby\n%s", util.FormatRequest(r))
	c := brokerhub.NewClient()
	c.ServeWs(w, r)
	brokerhub.Manager().ConnectLobby(c)
}

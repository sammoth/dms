package dms

import (
	"net/http"
	"github.com/anacrolix/dms/upnp"
)

type fakeContentManagerService struct {
	*Server
	upnp.Eventing
}

func (me *fakeContentManagerService) Handle(action string, argsXML []byte, r *http.Request) (map[string]string, error) {
	switch action {
	case "GetCurrentConnectionIDs":
		return map[string]string{
			"ConnectionIDs": "0",
		}, nil
	case "GetCurrentConnectionInfo":
		return map[string]string{
		}, nil
	case "GetProtocolInfo":
		return map[string]string{
		}, nil
	default:
		return nil, upnp.InvalidActionError
	}
}

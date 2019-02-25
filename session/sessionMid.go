package session

import (
	"encoding/json"
	"net/http"
	"time"
)

type SessionMid struct {
	SessionUserInfoKey string // "sui-v0"
	SessionKey         []byte
	Path               string
	Expire             time.Duration
}

func (sessionMid *SessionMid) GetSessionUserInfo(r *http.Request, sessionUserInfo interface{}) error {
	txt, serr := GetSession(r, sessionMid.SessionKey, sessionMid.SessionUserInfoKey)
	if serr != nil {
		return serr
	}
	if err := json.Unmarshal([]byte(txt), sessionUserInfo); err != nil {
		return err
	}
	return nil
}

func (sessionMid *SessionMid) SetSessionUserInfo(w http.ResponseWriter, sessionUserInfo interface{}) error {
	value, err := json.Marshal(sessionUserInfo)
	if err != nil {
		return err
	}
	return SetSession(w, sessionMid.SessionKey, sessionMid.SessionUserInfoKey, string(value), sessionMid.Path, sessionMid.Expire)
}

func (sessionMid *SessionMid) RemoveUserFromSession(w http.ResponseWriter) {
	RemoveSession(w, sessionMid.SessionUserInfoKey, sessionMid.Path)
}

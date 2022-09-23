package webtoy_base

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	SessionAuthorityUser = "session_authority_user"
	CookieSessionID      = "session_id"
	CookieSessionToken   = "session_token"
)

type SessionHandler struct {
	RedisClient         *redis.Client
	SessionExpireSecond int
}

type Session struct {
	UserID string `json:"user"`
	Token  string `json:"token"`
}

func NewSessionHandler(redisClient *redis.Client, sessionExpireSecond int) *SessionHandler {
	return &SessionHandler{
		RedisClient:         redisClient,
		SessionExpireSecond: sessionExpireSecond,
	}
}

func (this *SessionHandler) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieSessionID, err := r.Cookie(CookieSessionID)
		if err != nil {
			log.Debugf("failed get %v from cookie", cookieSessionID)
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_AUTH,
				ErrMsg: "request without authority",
			})
			return
		}

		cookieToken, err := r.Cookie(CookieSessionToken)
		if err != nil {
			log.Debugf("failed get %v from cookie", CookieSessionToken)
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_AUTH,
				ErrMsg: "request without authority",
			})
			return
		}

		stringCmd := this.RedisClient.Get(cookieSessionID.Value)
		if stringCmd.Err() != nil {
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_AUTH,
				ErrMsg: "session expired",
			})
			return
		}

		var session Session
		err = json.Unmarshal([]byte(stringCmd.Val()), &session)
		if err != nil {
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_INTERNAL,
				ErrMsg: "failed decode session",
			})
			return
		}

		if session.Token != cookieToken.Value {
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_AUTH,
				ErrMsg: "token not equal",
			})
			return
		}

		r.Header.Set(SessionAuthorityUser, session.UserID)

		//this.RedisClient.Expire(stringCmd.Val(), time.Second*time.Duration(this.SessionExpireSecond))
		this.RedisClient.Expire(cookieSessionID.Value, time.Second*time.Duration(this.SessionExpireSecond))

		h.ServeHTTP(w, r)
	})
}

func (this *SessionHandler) GenSession(userID int64, w *http.ResponseWriter) (string, *Session, error) {
	strUserID := strconv.FormatInt(userID, 10)

	sessionID, err := this.GenSessionID()
	if err != nil {
		return "", nil, err
	}
	token := this.GenToken(strUserID)

	session := &Session{
		UserID: strUserID,
		Token:  token,
	}

	err = this.SaveSession(sessionID, session)
	if err != nil {
		return "", nil, err
	}

	cookieSessionID := http.Cookie{Name: CookieSessionID, Value: sessionID, MaxAge: this.SessionExpireSecond, Path: "/"}
	http.SetCookie(*w, &cookieSessionID)

	cookieToken := http.Cookie{Name: CookieSessionToken, Value: token, MaxAge: this.SessionExpireSecond, Path: "/"}
	http.SetCookie(*w, &cookieToken)

	return sessionID, session, nil
}

func (this *SessionHandler) GenSessionID() (string, error) {
	uuid := uuid.NewV4()
	return uuid.String(), nil
}

func (this *SessionHandler) GenToken(plainText string) string {
	salt := time.Now().Unix()
	m5 := md5.New()
	m5.Write([]byte(plainText))
	m5.Write([]byte(fmt.Sprint(salt)))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}

func (this *SessionHandler) SaveSession(sessionID string, session *Session) error {
	jsonBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	statusCmd := this.RedisClient.Set(sessionID, string(jsonBytes), time.Second*time.Duration(this.SessionExpireSecond))
	return statusCmd.Err()
}

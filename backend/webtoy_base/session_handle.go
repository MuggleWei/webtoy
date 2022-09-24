package webtoy_base

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	SESSION_USER  = "userId"
	SESSION_ID    = "session"
	SESSION_TOKEN = "token"
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

func (this *SessionHandler) MiddlewareSessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("middleware session check: url=%v", r.URL.Path)
		userSessionID := r.Header.Get(SESSION_ID)
		userSessionToken := r.Header.Get(SESSION_TOKEN)

		userID, err := this.UpdateSession(userSessionID, userSessionToken)
		if err != nil {
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_AUTH,
				ErrMsg: err.Error(),
			})
			return
		}

		r.Header.Set(SESSION_USER, userID)

		next.ServeHTTP(w, r)
	})
}

func (this *SessionHandler) MiddlewareSessionUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSessionID := r.Header.Get(SESSION_ID)
		userSessionToken := r.Header.Get(SESSION_TOKEN)

		if len(userSessionID) > 0 && len(userSessionToken) > 0 {
			userID, err := this.UpdateSession(userSessionID, userSessionToken)
			if err == nil {
				r.Header.Set(SESSION_USER, userID)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (this *SessionHandler) GenSession(userID int64, w *http.ResponseWriter) (string, *Session, error) {
	strUserID := strconv.FormatInt(userID, 10)

	sessionID, err := this.GenSessionID()
	if err != nil {
		log.Errorf("failed GenSessionID")
		return "", nil, err
	}
	token, err := this.GenToken(strUserID)
	if err != nil {
		log.Errorf("failed GenToken")
		return "", nil, err
	}

	session := &Session{
		UserID: strUserID,
		Token:  token,
	}

	err = this.SaveSession(sessionID, session)
	if err != nil {
		return "", nil, err
	}

	cookieSessionID := http.Cookie{Name: SESSION_ID, Value: sessionID, MaxAge: this.SessionExpireSecond, Path: "/"}
	http.SetCookie(*w, &cookieSessionID)

	cookieToken := http.Cookie{Name: SESSION_TOKEN, Value: token, MaxAge: this.SessionExpireSecond, Path: "/"}
	http.SetCookie(*w, &cookieToken)

	return sessionID, session, nil
}

func (this *SessionHandler) GenSessionID() (string, error) {
	uuid := uuid.NewV4()
	return uuid.String(), nil
}

func (this *SessionHandler) GenToken(plainText string) (string, error) {
	salt, err := this.GenerateRandomBytes(32)
	if err != nil {
		log.Errorf("failed generate random bytes")
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(plainText))
	h.Write([]byte(salt))
	st := h.Sum(nil)
	return hex.EncodeToString(st), nil
}

func (this *SessionHandler) SaveSession(sessionID string, session *Session) error {
	jsonBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	statusCmd := this.RedisClient.Set(sessionID, string(jsonBytes), time.Second*time.Duration(this.SessionExpireSecond))
	return statusCmd.Err()
}

// refresh session
func (this *SessionHandler) UpdateSession(userSessionID, userSessionToken string) (string, error) {
	stringCmd := this.RedisClient.Get(userSessionID)
	if stringCmd.Err() != nil {
		return "", errors.New("session not exists")
	}

	var session Session
	err := json.Unmarshal([]byte(stringCmd.Val()), &session)
	if err != nil {
		return "", err
	}

	if session.Token != userSessionToken {
		return "", errors.New("incorrect session")
	}

	this.RedisClient.Expire(userSessionID, time.Second*time.Duration(this.SessionExpireSecond))

	return session.UserID, nil
}

func (this *SessionHandler) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

package webtoy_base

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	SESSION_USER  = "uid"
	SESSION_ID    = "session"
	SESSION_TOKEN = "token"
)

type SessionHandler struct {
	RedisClient         *redis.Client
	SessionExpireSecond int
}

type Session struct {
	UserID  string `json:"uid"`
	Session string `json:"session"`
	Token   string `json:"token"`
}

func NewSessionHandler(redisClient *redis.Client, sessionExpireSecond int) *SessionHandler {
	return &SessionHandler{
		RedisClient:         redisClient,
		SessionExpireSecond: sessionExpireSecond,
	}
}

func (this *SessionHandler) MiddlewareSessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("middleware session check begin: url=%v", r.URL.Path)
		session := Session{
			UserID:  r.Header.Get(SESSION_USER),
			Session: r.Header.Get(SESSION_ID),
			Token:   r.Header.Get(SESSION_TOKEN),
		}

		err := this.UpdateSession(&session)
		if err != nil {
			log.Debugf("failed update session: %+v", session)
			HttpResponse(w, &MessageRsp{
				Code:   ERROR_AUTH,
				ErrMsg: err.Error(),
			})
			return
		}

		next.ServeHTTP(w, r)

		log.Debugf("middleware session check end: url=%v", r.URL.Path)
	})
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

// generate and save user session
func (this *SessionHandler) GenSession(userID string) (*Session, error) {
	sessionID, err := this.GenSessionID()
	if err != nil {
		log.Errorf("failed GenSessionID")
		return nil, err
	}
	token, err := this.GenToken(userID)
	if err != nil {
		log.Errorf("failed GenToken")
		return nil, err
	}

	session := &Session{
		UserID:  userID,
		Session: sessionID,
		Token:   token,
	}

	err = this.SaveSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (this *SessionHandler) SaveSession(session *Session) error {
	jsonBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	sessionName := "session_" + session.UserID
	statusCmd := this.RedisClient.Set(sessionName, string(jsonBytes), time.Second*time.Duration(this.SessionExpireSecond))
	return statusCmd.Err()
}

// refresh session
func (this *SessionHandler) UpdateSession(session *Session) error {
	sessionName := "session_" + session.UserID
	stringCmd := this.RedisClient.Get(sessionName)
	if stringCmd.Err() != nil {
		return errors.New("session not exists")
	}

	var retSession Session
	err := json.Unmarshal([]byte(stringCmd.Val()), &retSession)
	if err != nil {
		return err
	}

	if session.UserID != retSession.UserID ||
		session.Session != retSession.Session ||
		session.Token != retSession.Token {
		return errors.New("incorrect session")
	}

	this.RedisClient.Expire(session.UserID, time.Second*time.Duration(this.SessionExpireSecond))

	return nil
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

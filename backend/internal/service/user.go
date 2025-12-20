package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"time"

	"github.com/your-team/taskmanager-chat/backend/internal/domain"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/mod/sumdb/storage"
)

type UserStorage interface {
	InsertUser(user domain.User) (int64, error)
	SelectUser(email string) (domain.User, error)
	SelectUserByID(userID int64) (domain.User, error)
	RefreshStore(userID int64, token string, expiresAt time.Time) error
	RefreshGet(token string) (int64, error)
	RefreshDelete(token string) error
	UserBlocked(email string, windowStart time.Time) ([]map[string]interface{}, error)
	LogAttempt(email string, result bool, attemptTime time.Time) error
	GetFailedLogAttempts(email string, windowStart time.Time) (int, error)
	BlockUser(email, blockedUntil string) error
	RenovationTwoFAStatus(userID int64, enabled bool) error
}

type TwoFaStorage interface{
	InsertTwoFaCode(userID int64, code string, expiresAt time.Time) error
	SelectTwoFaCodeByUserID(userID int64) (domain.TwoFaCode, error)
	RenovationTwoFaCodeAttempts(codeID int64, attempts int) error
	MarkTwoFaCodeUsed(codeID int64) error
	SelectRecentCodeRequests(userID int64, since time.Time) (int, error)
	SelectRecentVerificationAttempts(userID int64, since time.Time) (int, error)
}

type User struct {
	storage      UserStorage
	TwoFaStorage TwoFaStorage
	jwtSecret    string
}

func NewUser(storage UserStorage, twoFa TwoFaStorage, jwt string) *User{
	return &User{storage: storage, TwoFaStorage: twoFa, jwtSecret: jwt}
}
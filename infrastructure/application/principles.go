package application

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

var (
	jwtSigningMethods = make(map[string]*jwt.SigningMethodHMAC)
)

func init() {
	jwtSigningMethods[HS256] = jwt.SigningMethodHS256
	jwtSigningMethods[HS384] = jwt.SigningMethodHS384
	jwtSigningMethods[HS512] = jwt.SigningMethodHS512
}

//A Claims create an instance for JWT claim
type Claims struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Roles     []string  `json:"roles"`
	jwt.StandardClaims
}

//CreateJWT creates json web token
func (a *Application) CreateJWT(
	userid uuid.UUID, firstName string, lastName string, lifetime bool, roles ...string,
) (string, error) {
	preError := "app/crypto/crypto.go CreateJWT(), %s"
	var expirationTime int64
	if lifetime {
		expirationTime = time.Now().
			Add(time.Duration(24*365*100) * time.Hour).Unix()
	} else {
		expirationTime = time.Now().
			Add(time.Duration(a.Config.JWT.MaxAge) * time.Second).Unix()
	}
	claims := Claims{
		ID:    userid,
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(
		jwtSigningMethods[a.Config.JWT.Algorithm], claims,
	)
	tokenString, err := token.SignedString(a.Config.JWT.Secret)
	if err != nil {
		return "", fmt.Errorf(preError, err)
	}

	return tokenString, nil
}

//CreateRefreshToken creates refresh token
func (a *Application) CreateRefreshToken(userid uuid.UUID) (string, error) {
	preError := "app/crypto/crypto.go CreateRefreshToken(), %s"
	expirationTime := time.Now().
		Add(time.Duration(a.Config.JWT.RefreshToken.MaxAge) * time.Second)
	claims := Claims{
		ID:    userid,
		Roles: nil,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(
		jwtSigningMethods[a.Config.JWT.RefreshToken.Algorithm], claims,
	)
	tokenString, err := token.SignedString(a.Config.JWT.RefreshToken.Secret)
	if err != nil {
		return "", fmt.Errorf(preError, err)
	}

	return tokenString, nil
}

func (a *Application) secureCookie() *securecookie.SecureCookie {
	return securecookie.New(a.Config.SecretKey, a.Config.BlockSecretKey)
}

//SetCookie sets cookie to ResponseWriter
func (a *Application) SetCookie(
	name, value string, lifetime bool, w http.ResponseWriter,
) error {
	secureCookie := a.secureCookie()
	encoded, err := secureCookie.Encode(name, value)
	if err == nil {
		cookie := http.Cookie{
			Name:  name,
			Value: encoded,
		}
		if lifetime {
			cookie.Expires = time.Now().Add(24 * 365 * 10 * time.Hour)
		}
		http.SetCookie(w, &cookie)
	}

	return err
}

//ReadCookie reads cookie from request
func (a *Application) ReadCookie(r *http.Request, name string) (string, error) {
	secureCookie := a.secureCookie()
	cookie, err := r.Cookie(name)
	if err != nil {
		return EMPTY, err
	}

	var value string
	err = secureCookie.Decode(name, cookie.Value, &value)
	if err != nil {
		return EMPTY, err
	}

	return value, nil
}

//GenerateSymmetricKey create symmetric and secure key
func GenerateSymmetricKey(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	digits := DIGITS
	specials := SPECIALS
	all := ALL
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf) // E.g. "3i[g0|)z"

	return []byte(str)
}

//GenerateSymmetricKey create symmetric and secure key without spicials
func GenerateSymmetricKeyU(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	digits := DIGITS
	all := ALL_UNSPECIALS
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf) // E.g. "3i[g0|)z"

	return []byte(str)
}

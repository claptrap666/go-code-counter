package middleware

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config
var verifier *oidc.IDTokenVerifier
var authpoint map[string]string = make(map[string]string)

type User struct {
	UserID   string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
}

func OidcMiddleware(provider string, clientID string, clientSecret string, redirectURL string) gin.HandlerFunc {
	if verifier == nil {
		InitOidc(provider, clientID, clientSecret, redirectURL, context.Background())
	}
	return handleRedirect
}

func InitOidc(providerUrl string, clientID string, clientSecret string, redirectURL string, ctx context.Context) {

	provider, err := oidc.NewProvider(ctx, providerUrl)
	if err != nil {
		// handle error
		logrus.WithError(err).Errorln()
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})
}

func handleRedirect(c *gin.Context) {
	user := GetUserBySession(c)
	if user == nil {
		p, _ := rand.Int(rand.Reader, big.NewInt(99999999))
		state := p.String()
		authpoint[state] = c.Request.URL.RequestURI()
		c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
	}
	c.Next()
}

func HandleOAuth2Callback(c *gin.Context) {
	// Verify state and errors.
	ctx := context.Background()
	state := c.Query("state")
	oauth2Token, err := oauth2Config.Exchange(ctx, c.Query("code"))
	if err != nil {
		// handle error
		logrus.WithError(err).Warningln()
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
		logrus.Warn("missing token")
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		// handle error
		logrus.WithError(err).Warningln()
	}

	// logrus.Infoln(rawIDToken)
	// Extract custom claims
	claims := &User{}
	if err := idToken.Claims(claims); err != nil {
		// handle error
		logrus.WithError(err).Warningln()
	}
	logrus.Infoln(claims)
	SetUserToSession(c, claims)
	redirectURL := authpoint[state]
	delete(authpoint, state)
	c.Redirect(http.StatusFound, redirectURL)
}

func GetUserBySession(c *gin.Context) *User {
	session := sessions.Default(c)
	userv := session.Get("user")
	userStr := ""
	switch userv.(type) {
	case string:
		userStr = userv.(string)
	default:
		return nil
	}
	user := &User{}
	json.Unmarshal([]byte(userStr), user)
	return user
}

func SetUserToSession(c *gin.Context, user *User) bool {
	session := sessions.Default(c)
	userStr, err := json.Marshal(user)
	session.Set("user", string(userStr[:]))
	session.Save()
	return err == nil
}

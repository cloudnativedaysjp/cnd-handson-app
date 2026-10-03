package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// idp がまだ起動していないなど、1 回目の取得が失敗しても次のリクエストで取り直せること
func TestVerifierRetriesAfterFailedFetch(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	b64 := base64.RawURLEncoding.EncodeToString
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kty": "RSA", "kid": "k1", "n": b64(key.N.Bytes()), "e": b64(big.NewInt(int64(key.E)).Bytes()),
		}}})
	}))
	defer srv.Close()

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer: "iss", Audience: jwt.ClaimStrings{"aud"}, Subject: testUser,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
	})
	tok.Header["kid"] = "k1"
	signed, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}

	verify := newVerifier(srv.URL, "iss", "aud")
	if _, err := verify(context.Background(), signed); err == nil {
		t.Fatal("first call: want error while JWKS is unavailable")
	}
	if sub, err := verify(context.Background(), signed); err != nil || sub != testUser {
		t.Fatalf("second call: got %q, %v", sub, err)
	}
}

// 未知の kid の取得が idp で詰まっていても、キャッシュ済みの kid は待たずに返ること
func TestKnownKidNotBlockedByFetch(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	defer close(release)

	cached := &rsa.PublicKey{}
	j := &jwks{url: srv.URL, client: srv.Client(), keys: map[string]*rsa.PublicKey{"k1": cached}}
	go func() { _, _ = j.key(context.Background(), "unknown") }()
	time.Sleep(100 * time.Millisecond) // 取得が始まるのを待つ

	done := make(chan *rsa.PublicKey)
	go func() { k, _ := j.key(context.Background(), "k1"); done <- k }()
	select {
	case k := <-done:
		if k != cached {
			t.Fatal("got a different key for k1")
		}
	case <-time.After(time.Second):
		t.Fatal("cached kid waited for the in-flight fetch")
	}
}

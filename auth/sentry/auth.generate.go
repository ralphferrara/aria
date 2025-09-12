package sentry

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ralphferrara/aria/base/crypto"
	"github.com/ralphferrara/aria/responses"
)

//||------------------------------------------------------------------------------------------------||
//|| Simple in-memory limiter (per-IP)
//||------------------------------------------------------------------------------------------------||

var genKeyLimit = struct {
	sync.Mutex
	clients map[string]time.Time
}{
	clients: make(map[string]time.Time),
}

const genKeyCooldown = 5 * time.Second // allow 1 request every 5s per IP

//||------------------------------------------------------------------------------------------------||
//|| GenerateKeyPairHandler (with built-in rate limiting)
//||------------------------------------------------------------------------------------------------||

func GenerateKeyPairHandler(w http.ResponseWriter, r *http.Request) {

	//||------------------------------------------------------------------------------------------------||
	//|| IP Rate Limiting Check
	//||------------------------------------------------------------------------------------------------||

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	genKeyLimit.Lock()
	last, exists := genKeyLimit.clients[ip]
	now := time.Now()
	if exists && now.Sub(last) < genKeyCooldown {
		genKeyLimit.Unlock()
		http.Error(w, "Too many requests, slow down", http.StatusTooManyRequests)
		return
	}
	genKeyLimit.clients[ip] = now
	genKeyLimit.Unlock()

	//||------------------------------------------------------------------------------------------------||
	//|| Generate Key Pair using helpers.GenerateKeyPair()
	//||------------------------------------------------------------------------------------------------||

	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Failed to generate key pair")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Respond with Keys
	//||------------------------------------------------------------------------------------------------||
	responses.Success(w, http.StatusOK, map[string]string{
		"privateKey": privateKey,
		"publicKey":  publicKey,
	})
}

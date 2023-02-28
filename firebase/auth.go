package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware : to verify all authorized operations
func AuthMiddleware(c *gin.Context) {
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)

	authorizationToken := c.GetHeader("Authorization")
	reqToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))

	if reqToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not available"})
		c.Abort()
		return
	}

	// verify and get token
	idToken, err := signInWithCustomToken(reqToken)
	if err != nil {
		idToken = reqToken
	}

	fmt.Printf("Firebase idToken %s", idToken)

	token, err := firebaseAuth.VerifyIDToken(c, idToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	fmt.Printf("Token uid %s", token.UID)

	c.Set("UUID", token.UID)

	emails := token.Firebase.Identities["email"].([]interface{})
	c.Set("userEmail", emails[0].(string))

	// set custom claims to request context
	c.Set("userId", token.Claims["userId"])
	c.Set("userRoles", token.Claims["userRoles"])
	c.Set("partnerId", token.Claims["partnerId"])
	c.Set("memberRoles", token.Claims["memberRoles"])

	c.Next()
}

const (
	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
)

func signInWithCustomToken(token string) (string, error) {
	body, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("FIREBASE_API_KEY")
	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, apiKey), body)
	if err != nil {
		return "", err
	}

	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}

	return respBody.IDToken, err
}

func postRequest(url string, body []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

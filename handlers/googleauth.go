// handlers/googleauth.go
package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
    RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
    Endpoint:     google.Endpoint,
}

func GoogleLoginURLHandler(w http.ResponseWriter, r *http.Request) {
    url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Missing code", http.StatusBadRequest)
        return
    }

    token, err := googleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }

    resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var userInfo map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&userInfo)

    // TODO: issue JWT for your app user
    jwt := "your-signed-jwt-here"

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, `<script>
        window.opener.postMessage({ token: "%s", user: %s }, "*");
        window.close();
    </script>`, jwt, toJSON(userInfo))
}

func toJSON(v interface{}) string {
    b, _ := json.Marshal(v)
    return string(b)
}

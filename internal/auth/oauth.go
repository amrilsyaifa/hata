package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// localAuthPort is the fixed port for the OAuth redirect.
// You MUST register exactly "http://localhost:8085" in:
// Google Cloud Console → APIs & Services → Credentials → your OAuth client → Authorized redirect URIs
const localAuthPort = 8085

// credentialsFile represents the structure of a downloaded OAuth credentials JSON.
type credentialsFile struct {
	Installed *credentialsDetails `json:"installed"`
	Web       *credentialsDetails `json:"web"`
}

type credentialsDetails struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
	RedirectURIs []string `json:"redirect_uris"`
}

func parseCredentials(b []byte) (*oauth2.Config, error) {
	var f credentialsFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("invalid credentials JSON: %w", err)
	}

	var details *credentialsDetails
	switch {
	case f.Installed != nil:
		details = f.Installed
	case f.Web != nil:
		details = f.Web
	default:
		return nil, fmt.Errorf("credentials file must have an 'installed' or 'web' key")
	}

	if details.ClientID == "" || details.ClientSecret == "" {
		return nil, fmt.Errorf("credentials file is missing client_id or client_secret")
	}

	authURI := details.AuthURI
	if authURI == "" {
		authURI = google.Endpoint.AuthURL
	}
	tokenURI := details.TokenURI
	if tokenURI == "" {
		tokenURI = google.Endpoint.TokenURL
	}

	return &oauth2.Config{
		ClientID:     details.ClientID,
		ClientSecret: details.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURI,
			TokenURL: tokenURI,
		},
		Scopes: []string{sheets.SpreadsheetsScope},
	}, nil
}

func OAuthOption(ctx context.Context, credentialsPath, tokenPath string) (option.ClientOption, error) {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w", err)
	}

	oauthCfg, err := parseCredentials(b)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials file: %w", err)
	}

	token, err := loadToken(tokenPath)
	if err != nil || !token.Valid() {
		token, err = getTokenFromWeb(ctx, oauthCfg, tokenPath)
		if err != nil {
			return nil, err
		}
	}

	client := oauthCfg.Client(ctx, token)
	return option.WithHTTPClient(client), nil
}

func loadToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	if err := json.NewDecoder(f).Decode(&tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("unable to create token directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to save token: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(ctx context.Context, oauthCfg *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	redirectURL := fmt.Sprintf("http://localhost:%d", localAuthPort)
	oauthCfg.RedirectURL = redirectURL

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", localAuthPort))
	if err != nil {
		return nil, fmt.Errorf("failed to start local server on port %d (already in use?): %w", localAuthPort, err)
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				errCh <- fmt.Errorf("no authorization code in callback")
				http.Error(w, "Authorization failed.", http.StatusBadRequest)
				return
			}
			fmt.Fprintln(w, "<html><body><h2>Authorization successful!</h2><p>You can close this tab and return to the terminal.</p></body></html>")
			codeCh <- code
		}),
	}

	go func() { _ = server.Serve(listener) }()
	defer server.Close()

	authURL := oauthCfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("Redirect URI (must be registered in Google Cloud Console): %s\n\n", redirectURL)
	fmt.Printf("If the browser does not open, visit:\n  %s\n\n", authURL)
	openBrowser(authURL)

	select {
	case code := <-codeCh:
		token, err := oauthCfg.Exchange(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve token: %w", err)
		}
		if err := saveToken(tokenPath, token); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save token: %v\n", err)
		}
		fmt.Println("Authentication successful! Token saved.")
		return token, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	_ = cmd.Start()
}

package auth

import (
	"net/http"
	"path"
	"strings"
)

// NewServeMux generate http.Handler for auth
func (auth *Auth) NewServeMux() http.Handler {
	return &serveMux{Auth: auth}
}

type serveMux struct {
	*Auth
}

// ServeHTTP dispatches the handler registered in the matched route
func (serveMux *serveMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		claims  *Claims
		reqPath = strings.TrimPrefix(req.URL.Path, serveMux.Prefix)
		paths   = strings.Split(reqPath, "/")
		session = &Session{Auth: serveMux.Auth, Claims: claims}
	)

	if len(paths) >= 2 {
		// eg: /phone/login

		if provider := serveMux.Auth.GetProvider(paths[0]); provider != nil {
			session.Provider = provider

			// serve mux
			switch paths[1] {
			case "login":
				provider.Login(req, w, session)
			case "logout":
				provider.Logout(req, w, session)
			case "register":
				provider.Register(req, w, session)
			case "callback":
				provider.Callback(req, w, session)
			default:
				provider.ServeHTTP(req, w, session)
			}
			return
		}
	} else if len(paths) == 1 {
		// eg: /login, /logout

		switch paths[0] {
		case "login":
			// render login page
			serveMux.Auth.Render.Execute("auth/login", session, req, w)
			return
		case "logout":
			// destroy login session
			serveMux.Auth.LogoutHandler(req, w, session)
			return
		case "register":
			// render register page
			serveMux.Auth.Render.Execute("auth/register", session, req, w)
			return
		}
	}

	http.NotFound(w, req)
}

// Session session
type Session struct {
	*Auth
	*Claims
	Provider
	Params map[string]interface{}
}

// AuthURL generate URL for auth
func (auth *Auth) AuthURL(pth string) string {
	return path.Join(auth.Prefix, pth)
}

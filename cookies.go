package rapidroot

import (
	"net/http"
	"time"
)

type cookies struct {
	defaults *http.Cookie
}

// Cookie returns one value from cookies.
func (r *Request) Cookie(key string) (*http.Cookie, error) {
	cookie, err := r.Req.Cookie(key)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

// Cookies returns slice of all cookies.
func (r *Request) Cookies() []*http.Cookie {
	return r.Req.Cookies()
}

// SetCookie put key-value to cookies.
func (r *Request) SetCookie(key string, val string, exp time.Time) {
	http.SetCookie(r.Writer, &http.Cookie{
		Name:     key,
		Value:    val,
		Expires:  exp,
		HttpOnly: r.cookie.defaults.HttpOnly,
		Secure:   r.cookie.defaults.Secure,
		SameSite: r.cookie.defaults.SameSite,
		Path:     r.cookie.defaults.Path,
	})
}

// SetCookieObj puts cookie object to cookies.
func (r *Request) SetCookieObj(cookie *http.Cookie) {
	http.SetCookie(r.Writer, cookie)
}

// SetCookiesHTTPOnly sets httpOnly to all cookies.
func (r *Request) SetCookiesHTTPOnly(httpOnly bool) {
	r.cookie.defaults.HttpOnly = httpOnly
}

// SetCookiesSecure sets secure to all cookies.
func (r *Request) SetCookiesSecure(secure bool) {
	r.cookie.defaults.Secure = secure
}

// SetCookiesSameSite sets sameSite to all cookies.
func (r *Request) SetCookiesSameSite(same http.SameSite) {
	r.cookie.defaults.SameSite = same
}

// RemoveCookie removes cookie by key.
func (r *Request) RemoveCookie(key string) {
	r.SetCookie(key, "", time.Unix(0, 0))
}

// RemoveAllCookies removes all cookies from request.
func (r *Request) RemoveAllCookies() {
	for _, cookie := range r.Cookies() {
		r.RemoveCookie(cookie.Name)
	}
}

// SetCookieWithOptions puts a key-value pair into the cookies with additional options.
func (r *Request) SetCookieWithOptions(key, val string, exp time.Time) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    val,
		Expires:  exp,
		HttpOnly: r.cookie.defaults.HttpOnly,
		Secure:   r.cookie.defaults.Secure,
		SameSite: r.cookie.defaults.SameSite,
		Path:     r.cookie.defaults.Path,
		Domain:   r.cookie.defaults.Domain,
		MaxAge:   r.cookie.defaults.MaxAge,
		Raw:      r.cookie.defaults.Raw,
		RawExpires: r.cookie.defaults.RawExpires,
		Unparsed: r.cookie.defaults.Unparsed,
	}

	http.SetCookie(r.Writer, cookie)
}

// SetDefaultCookieOptions sets default values for cookie attributes.
func (r *Request) SetDefaultCookieOptions(options *http.Cookie) {
	if options == nil {
		options = &http.Cookie{}
	}
	r.cookie.defaults = options
}

// SetCookiePath sets the path for the cookie.
func (r *Request) SetCookiePath(path string) {
	if path == "" {
		path = "/"
	}
	r.cookie.defaults.Path = path
}

// SetSecureFlagAutomatically sets the Secure flag based on the request's scheme.
func (r *Request) SetSecureFlagAutomatically() {
	if r.Req.TLS != nil {
		r.cookie.defaults.Secure = true
	} else {
		r.cookie.defaults.Secure = false
	}
}

// RemoveAllCookiesExcept removes all cookies from the response except the specified ones.
func (r *Request) RemoveAllCookiesExcept(exceptions ...string) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if !contains(exceptions, cookie.Name) {
			r.RemoveCookie(cookie.Name)
		}
	}
}

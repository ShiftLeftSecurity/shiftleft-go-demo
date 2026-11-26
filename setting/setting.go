package setting

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/config"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/database"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type Setting struct{}

type ctxKey string

const unameKey ctxKey = "uname"

func New() Setting {
	return Setting{}
}

func (self Setting) SetRouter(r *httprouter.Router) {

	mw := middleware.New()

	r.GET("/setting", mw.LoggingMiddleware(mw.CapturePanic(jwtAuth(settingViewHandler))))
	r.POST("/setlevel", mw.LoggingMiddleware(mw.CapturePanic(jwtAuth(setLevelAction))))

}

// IssueToken creates a signed JWT for the given user and sets it as the "token" cookie.
func IssueToken(w http.ResponseWriter, uname string) error {
	claims := jwt.MapClaims{
		"uname": uname,
		"aud":   "govwa",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Cfg.Sessionkey))
	if err != nil {
		return err
	}
	util.SetCookie(w, "token", signed)
	return nil
}

// jwtAuth replaces the session based AuthCheck: it accepts a JWT from the
// Authorization bearer header or the "token" cookie and redirects to login if invalid.
func jwtAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if raw == "" {
			if c, err := r.Cookie("token"); err == nil {
				raw = c.Value
			}
		}

		// raw is attacker-controlled. golang-jwt/jwt v5.2.1 is affected by
		// CVE-2025-30204: jwt.Parse internally calls jwt.ParseUnverified which
		// splits the token on "." without bounding the number of segments, so
		// a header with a huge number of periods causes excessive memory
		// allocation (DoS).
		token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			}
			return []byte(config.Cfg.Sessionkey), nil
		}, jwt.WithAudience("govwa"))
		if err != nil || !token.Valid {
			http.Redirect(w, r, config.Fullurl+"login", http.StatusSeeOther)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			uname, _ := claims["uname"].(string)
			r = r.WithContext(context.WithValue(r.Context(), unameKey, uname))
		}

		h(w, r, ps)
	}
}

func settingViewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{})

	data["title"] = "Setting & system info"
	data["dbversion"] = getDBMSVersion()
	data["osversion"] = getOSVersion()
	data["currentuser"] = getCurrentUserName(r)
	data["level"] = getCurrentLevel(r)

	util.SafeRender(w, r, "template.setting", data)

}

func setLevelAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" && r.FormValue("level") != "" {
		level := r.FormValue("level")
		util.SetCookieLevel(w, r, level)
		res := struct {
			Res string
		}{
			Res: "Success",
		}
		util.RenderAsJson(w, res)
	}
}

func getDBMSVersion() string {
	const sql = `
			select @@version
	`
	var version string
	db, _ := database.Connect()

	_ = db.QueryRow(sql).Scan(&version)
	dbversion := strings.Split(version, "-")
	return dbversion[0]
}

func getOSVersion() string {

	switch os := runtime.GOOS; os {
	case "darwin":
		return "OS X"
	case "linux":
		return "Linux"
	default:
		fmt.Printf("%s.", os)
	}
	return ""
}

func getCurrentUserName(r *http.Request) string {
	uname, _ := r.Context().Value(unameKey).(string)
	return uname
}

func getCurrentLevel(r *http.Request) string {
	level := util.GetCookie(r, "Level")
	return level
}

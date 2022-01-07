package setup

import (
	"fmt"
	"net/http"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/user/session"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/config"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/database"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/middleware"

	"github.com/julienschmidt/httprouter"
)

type Setup struct{}

func New() Setup {
	return Setup{}
}

func (self Setup) SetRouter(r *httprouter.Router) {

	mw := middleware.New()

	r.GET("/setup", mw.LoggingMiddleware(mw.CapturePanic(setupViewHandler)))
	r.POST("/setupaction", mw.LoggingMiddleware(mw.CapturePanic(setupActionHandler)))

}

type JsonResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func setupViewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var info string

	data := make(map[string]interface{})
	ok, err := database.CheckDatabase()

	if !ok || err != nil {
		info = fmt.Sprintf(`<div id="info" class="alert alert-danger">%s
							<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a></div>`, err.Error())
		data["error"] = util.ToHTML(info)
	} else {

		info = fmt.Sprintf(`<div id="info" class="alert alert-success">Connection Success Click Reset to reset Database<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a></div>`)
		data["error"] = util.ToHTML(info)
	}

	data["title"] = "Setup/Reset"
	data["weburl"] = config.Fullurl
	util.SafeRender(w, r, "template.setup", data)
}

func setupActionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if r.Method == "POST" && r.FormValue("act") == "cr" {

		res := []JsonResp{}
		loginUrl := util.ToHTML(fmt.Sprintf(`<a href=%slogin>Login</a>`, config.Fullurl))

		err = createUsersTable() //create users table
		if err != nil {
			res = append(res, JsonResp{
				Status:  "0",
				Message: err.Error(),
			})
		} else {
			res = append(res, JsonResp{
				Status:  "1",
				Message: "Create Users Table Success Please " + string(loginUrl),
			})
		}

		err = createProfileTable() //create profilet table

		if err != nil {
			res = append(res, JsonResp{
				Status:  "0",
				Message: err.Error(),
			})
		} else {
			res = append(res, JsonResp{
				Status:  "1",
				Message: "Create Profile Table Success Please " + string(loginUrl),
			})
		}

		/* clear login session when reseting */
		s := session.New()
		s.DeleteSession(w, r)
		cookies := []string{"Level", "Uid"}
		util.DeleteCookie(w, cookies)

		util.RenderAsJson(w, res)

	}
}

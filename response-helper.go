package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"./models"
	"github.com/geo-stanciu/go-utils/utils"
)

// ResponseHelper - HTTPS response utils
type ResponseHelper struct {
	Title           string `sql:"name"`
	Template        string `sql:"request_template"`
	Controller      string `sql:"controller"`
	Action          string `sql:"action"`
	RedirectURL     string `sql:"redirect_url"`
	RedirectOnError string `sql:"redirect_on_error"`
}

func (res *ResponseHelper) getResponse(w http.ResponseWriter, r *http.Request) (models.ResponseModel, error) {
	switch res.Controller {
	case "Home":
		home := HomeController{}
		return res.getResponseValue(home, w, r)

	default:
		return nil, nil
	}
}

func (res *ResponseHelper) getResponseValue(controller interface{}, w http.ResponseWriter, r *http.Request) (models.ResponseModel, error) {
	if len(res.Action) == 0 || res.Action == "-" {
		return nil, nil
	}

	response := utils.InvokeMethodByName(controller, res.Action, w, r, res)

	if len(response) >= 2 {
		r := response[0].Interface()
		i2 := response[1].Interface()

		if r == nil && i2 != nil {
			return nil, i2.(error)
		}

		if i2 != nil {
			return r.(models.ResponseModel), i2.(error)
		}

		return r.(models.ResponseModel), nil
	}

	return nil, fmt.Errorf("Function does not return the requested number of values")
}

func getResponseHelperByURL(sessionData *SessionData, url string, requestType string) (*ResponseHelper, error) {
	var res ResponseHelper
	var sURL string

	if url == "/" {
		sURL = "index"
	} else {
		sURL = strings.Replace(url[1:], ".html", "", 1)
	}

	var suser string
	var lang string
	if sessionData != nil {
		suser = sessionData.User.Username
		lang = sessionData.Lang
	}
	if len(suser) == 0 {
		suser = "-"
	}
	if len(lang) == 0 {
		lang = "EN"
	}

	pq := dbUtils.PQuery(`
		WITH access AS (
			SELECT rr.request_id
			  FROM "user" u, user_role ur, request_role rr
			 WHERE u.user_id = ur.user_id
			   AND ur.role_id = rr.role_id
			   AND u.loweredusername = lower(?)
			UNION ALL
			SELECT rr.request_id 
              FROM role r, request_role rr
			 WHERE r.role_id = rr.role_id
			   AND r.loweredrole = lower(?)
		),
		name AS (
			SELECT request_id, name
			  FROM request_name
			 WHERE language = ?
		)
		SELECT case when nm.name is null then '-' else nm.name end AS name,
			r.request_template,
			r.controller,
			r.action,
			r.redirect_url,
			r.redirect_on_error
		FROM request r
		LEFT OUTER JOIN name nm ON (r.request_id = nm.request_id)
		WHERE r.request_url = ?
		AND r.request_type = ?
		AND r.request_id IN (
			select request_id from access
		)
	`, suser,
		"All",
		lang,
		sURL,
		requestType)

	err := dbUtils.RunQuery(pq, &res)

	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("request \"%s\" - not found or access denied", url)
		return nil, err
	case err != nil:
		return nil, err
	}

	return &res, nil
}

package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

// RequestHelper - GET and POST request helper
type RequestHelper struct {
	sync.RWMutex
	tx        *sql.Tx
	RequestID int `sql:"request_id"`
	urlRequest
}

type urlRequest struct {
	RequestType     string `sql:"request_type"`
	RequestURL      string `sql:"request_url"`
	RequestTemplate string `sql:"request_template"`
	Controller      string `sql:"controller"`
	Action          string `sql:"action"`
	RedirectURL     string `sql:"redirect_url"`
	RedirectOnError string `sql:"redirect_on_error"`
	IndexLevel      int    `sql:"index_level"`
	OrderNumber     int    `sql:"order_number"`
	FireEvent       int    `sql:"fire_event"`
}

var requestLock sync.RWMutex

// Exists - check to see if request already exists
func (r *RequestHelper) Exists() (bool, error) {
	var err error
	var found bool

	pq := dbUtils.PQuery(`
		select CASE WHEN EXISTS (
		select 1
			from request
			where request_url = ?
			and request_type = ?
		) THEN 1 ELSE 0 END
		FROM dual
	`, r.RequestURL,
		r.RequestType)

	err = r.tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// GetByURL - get user by url
func (r *RequestHelper) GetByURL(requestType string, requestURL string) error {
	r.Lock()
	defer r.Unlock()

	pq := dbUtils.PQuery(`
		select request_id,
			   request_type,
			   request_url,
	           request_template,
	           controller,
	           action,
	           redirect_url,
			   redirect_on_error,
			   case when index_level is null then -1 else index_level end AS index_level,
			   case when order_number is null then -1 else order_number end AS order_number,
			   fire_event
	      from request
		 where request_url = ?
		   and request_type = ?
	`, requestURL,
		requestType)

	err := dbUtils.RunQueryTx(r.tx, pq, r)

	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("request not found")
	case err != nil:
		return err
	}

	return nil
}

// GetByID - get user by id
func (r *RequestHelper) GetByID(requestID int) error {
	r.Lock()
	defer r.Unlock()

	pq := dbUtils.PQuery(`
		select request_id,
			   request_type,
			   request_url,
	           request_template,
	           controller,
	           action,
	           redirect_url,
			   redirect_on_error,
			   case when index_level is null then -1 else index_level end AS index_level,
			   case when order_number is null then -1 else order_number end AS order_number,
			   fire_event
	      from request
	     where request_id = ?
	`, requestID)

	err := dbUtils.RunQueryTx(r.tx, pq, r)

	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("request not found")
	case err != nil:
		return err
	}

	return nil
}

// Load - load details from urlRequest struct
func (r *RequestHelper) Load(req *urlRequest) {
	r.RequestType = req.RequestType
	r.RequestURL = req.RequestURL
	r.RequestTemplate = req.RequestTemplate
	r.Controller = req.Controller
	r.Action = req.Action
	r.RedirectURL = req.RedirectURL
	r.RedirectOnError = req.RedirectOnError
	r.IndexLevel = req.IndexLevel
	r.OrderNumber = req.OrderNumber
	r.FireEvent = req.FireEvent
}

func (r *RequestHelper) testSave() error {
	if len(r.RequestURL) == 0 || len(r.RequestType) == 0 {
		return fmt.Errorf("unknown request \"%s\"", r.RequestURL)
	}

	var found bool

	pq := dbUtils.PQuery(`
	    SELECT CASE WHEN EXISTS (
	        SELECT 1
			  from request
			 where request_url = ?
			   and request_type = ?
	           and request_id <> ?
	    ) THEN 1 ELSE 0 END
	    FROM dual
	`, r.RequestURL,
		r.RequestType,
		r.RequestID)

	err := r.tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return err
	}

	if found {
		return fmt.Errorf("duplicate request \"%s:%s\"", r.RequestType, r.RequestURL)
	}

	return nil
}

// Save - save the request
func (r *RequestHelper) Save() error {
	requestLock.Lock()
	defer requestLock.Unlock()

	err := r.testSave()
	if err != nil {
		return err
	}

	if r.RequestID <= 0 {
		pq := dbUtils.PQuery(`
			insert into request (
				request_template,
				request_url,
				request_type,
				controller,
				action,
				redirect_url,
				redirect_on_error,
				index_level,
				order_number,
				fire_event
			)
			values (
				?, ?, ?, ?, ?, ?, ?,
				case when ? <= 0 then CAST(null AS int) else ? end,
				case when ? <= 0 then CAST(null AS int) else ? end,
				?
			)
		`, r.RequestTemplate,
			r.RequestURL,
			r.RequestType,
			r.Controller,
			r.Action,
			r.RedirectURL,
			r.RedirectOnError,
			r.IndexLevel,
			r.IndexLevel,
			r.OrderNumber,
			r.OrderNumber,
			r.FireEvent,
		)

		_, err = dbUtils.ExecTx(r.tx, pq)
		if err != nil {
			return err
		}

		pq = dbUtils.PQuery(`
			SELECT request_id
			  FROM request
			 WHERE request_url = ?
			   AND request_type = ?
		`, r.RequestURL,
			r.RequestType)

		err = r.tx.QueryRow(pq.Query, pq.Args...).Scan(&r.RequestID)

		switch {
		case err == sql.ErrNoRows:
			r.RequestID = -1
		case err != nil:
			return err
		}

		audit.Log(nil, "add-request", "Add new request.", "new", r)
	} else {
		old := RequestHelper{tx: r.tx}
		err = old.GetByID(r.RequestID)
		if err != nil {
			return err
		}

		if !r.Equals(&old) {
			pq := dbUtils.PQuery(`
			    UPDATE request
			       SET request_template = ?,
				   request_url = ?,
				   request_type = ?,
				   controller = ?,
				   action = ?,
				   redirect_url = ?,
				   redirect_on_error = ?,
				   index_level = case when ? <= 0 then CAST(null AS int) else ? end,
				   order_number = case when ? <= 0 then CAST(null AS int) else ? end,
				   fire_event = ?
			     WHERE request_id = ?
			`, r.RequestTemplate,
				r.RedirectURL,
				r.RequestType,
				r.Controller,
				r.Action,
				r.RedirectURL,
				r.RedirectOnError,
				r.IndexLevel,
				r.IndexLevel,
				r.OrderNumber,
				r.OrderNumber,
				r.FireEvent,
				r.RequestID,
			)

			_, err = dbUtils.ExecTx(r.tx, pq)
			if err != nil {
				return err
			}

			audit.Log(nil, "update-request", "Update request.", "old", &old, "new", r)
		}
	}

	return nil
}

// Equals - test equality between Request structs
func (r *RequestHelper) Equals(r1 *RequestHelper) bool {
	if r == nil && r1 != nil ||
		r != nil && r1 == nil ||
		r.RequestID != r1.RequestID ||
		r.RequestType != r1.RequestType ||
		r.RequestURL != r1.RequestURL ||
		r.RequestTemplate != r1.RequestTemplate ||
		r.Controller != r1.Controller ||
		r.Action != r1.Action ||
		r.RedirectURL != r1.RedirectURL ||
		r.RedirectOnError != r1.RedirectOnError ||
		r.IndexLevel != r1.IndexLevel ||
		r.OrderNumber != r1.OrderNumber ||
		r.FireEvent != r1.FireEvent {

		return false
	}

	return true
}

func getClientIP(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")

	ipList := strings.Split(ips, ", ")

	if len(ipList) > 0 && len(ipList[0]) > 0 {
		return ipList[0]
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}

func isRequestFromLocalhost(r *http.Request) bool {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	if ip == "127.0.0.1" || ip == "::1" {
		return true
	}

	return false
}

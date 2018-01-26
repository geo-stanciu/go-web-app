package main

import "database/sql"

func addRequests(tx *sql.Tx) (bool, error) {

	/*
		If no action (method) from controller is to be called, put -
	*/
	requests := []urlRequest{
		// pages
		{
			RequestType:     "GET",
			RequestURL:      "index",
			RequestTemplate: "home/index.html",
			Controller:      "Home",
			Action:          "Index",
			RedirectURL:     "-",
			RedirectOnError: "-",
			IndexLevel:      1,
			OrderNumber:     1,
			FireEvent:       1,
		},
		{
			RequestType:     "GET",
			RequestURL:      "users",
			RequestTemplate: "home/users.html",
			Controller:      "Home",
			Action:          "-",
			RedirectURL:     "-",
			RedirectOnError: "-",
			IndexLevel:      1,
			OrderNumber:     2,
			FireEvent:       1,
		},
		{
			RequestType:     "GET",
			RequestURL:      "about",
			RequestTemplate: "home/about.html",
			Controller:      "Home",
			Action:          "-",
			RedirectURL:     "-",
			RedirectOnError: "-",
			IndexLevel:      1,
			OrderNumber:     3,
			FireEvent:       1,
		},
		{
			RequestType:     "GET",
			RequestURL:      "login",
			RequestTemplate: "home/login.html",
			Controller:      "Home",
			Action:          "-",
			RedirectURL:     "-",
			RedirectOnError: "-",
			IndexLevel:      1,
			OrderNumber:     4,
			FireEvent:       1,
		},
		{
			RequestType:     "GET",
			RequestURL:      "register",
			RequestTemplate: "home/register.html",
			Controller:      "Home",
			Action:          "-",
			RedirectURL:     "-",
			RedirectOnError: "-",
			IndexLevel:      1,
			OrderNumber:     5,
			FireEvent:       1,
		},
		{
			RequestType:     "GET",
			RequestURL:      "change-password",
			RequestTemplate: "home/change-password.html",
			Controller:      "Home",
			Action:          "-",
			RedirectURL:     "-",
			RedirectOnError: "-",
			IndexLevel:      1,
			OrderNumber:     6,
			FireEvent:       1,
		},
		// gets
		{
			RequestType:     "GET",
			RequestURL:      "logout",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "Logout",
			RedirectURL:     "/",
			RedirectOnError: "-",
			FireEvent:       1,
		},
		{
			RequestType:     "GET",
			RequestURL:      "exchange-rates",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "GetExchangeRates",
			RedirectURL:     "-",
			RedirectOnError: "-",
			FireEvent:       1,
		},
		// posts
		{
			RequestType:     "POST",
			RequestURL:      "login",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "Login",
			RedirectURL:     "index",
			RedirectOnError: "login",
			FireEvent:       1,
		},
		{
			RequestType:     "POST",
			RequestURL:      "logout",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "Logout",
			RedirectURL:     "login",
			RedirectOnError: "login",
			FireEvent:       1,
		},
		{
			RequestType:     "POST",
			RequestURL:      "register",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "Register",
			RedirectURL:     "login",
			RedirectOnError: "register",
			FireEvent:       1,
		},
		{
			RequestType:     "POST",
			RequestURL:      "change-password",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "ChangePassword",
			RedirectURL:     "change-password",
			RedirectOnError: "change-password",
			FireEvent:       1,
		},
		{
			RequestType:     "POST",
			RequestURL:      "exchange-rates",
			RequestTemplate: "-",
			Controller:      "Home",
			Action:          "GetExchangeRates",
			RedirectURL:     "-",
			RedirectOnError: "-",
			FireEvent:       1,
		},
	}

	foundNew := false

	for _, req := range requests {
		r := RequestHelper{tx: tx}
		r.Load(&req)

		found, err := r.Exists()
		if err != nil {
			return false, err
		}

		if found {
			continue
		}

		foundNew = true

		err = r.Save()
		if err != nil {
			return false, err
		}
	}

	return foundNew, nil
}

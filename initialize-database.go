package main

func initializeDatabase() error {
	var err error
	var foundNew bool

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if foundNew, err = addRequests(tx); foundNew || err != nil {
		audit.Log(err, "initialize", "requests")
	}

	if foundNew, err = addRoles(tx); foundNew || err != nil {
		audit.Log(err, "initialize", "roles")
	}

	if foundNew, err = addGetRequestsAccessRules(tx); foundNew || err != nil {
		audit.Log(err, "initialize", "access rules - GET")
	}

	if foundNew, err = addPostRequestsAccessRules(tx); foundNew || err != nil {
		audit.Log(err, "initialize", "access rules - POST")
	}

	if foundNew, err = addGeneralMemberRequestsAccessRules(tx); foundNew || err != nil {
		audit.Log(err, "initialize", "access rules - members")
	}

	if foundNew, err = addChildRequestAccessRules(tx); foundNew || err != nil {
		audit.Log(err, "initialize", "access rules - child requests")
	}

	tx.Commit()

	return err
}

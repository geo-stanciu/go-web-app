package main

func initializeDatabase() error {
	var err error
	var foundNew bool

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if foundNew, err = addRequests(tx); foundNew {
		audit.Log(err, "initialize", "requests")
	}
	if err != nil {
		return err
	}

	if foundNew, err = addRoles(tx); foundNew {
		audit.Log(err, "initialize", "roles")
	}
	if err != nil {
		return err
	}

	if foundNew, err = addGetRequestsAccessRules(tx); foundNew {
		audit.Log(err, "initialize", "access rules - GET")
	}
	if err != nil {
		return err
	}

	if foundNew, err = addPostRequestsAccessRules(tx); foundNew {
		audit.Log(err, "initialize", "access rules - POST")
	}
	if err != nil {
		return err
	}

	if foundNew, err = addGeneralMemberRequestsAccessRules(tx); foundNew {
		audit.Log(err, "initialize", "access rules - members")
	}
	if err != nil {
		return err
	}

	if foundNew, err = addChildRequestAccessRules(tx); foundNew {
		audit.Log(err, "initialize", "access rules - child requests")
	}
	if err != nil {
		return err
	}

	tx.Commit()

	return err
}

// using crud queries

	// assigning user struct from models.go 
	var user crud.User

    // assigning queries struct from db.go
	var query *crud.Queries

    // pass the db connection to the new method
	query = crud.New(db)

    // query will now have all the query methods created in the "...query.sql" files
    // provide a context and paramsstruct variable
    // this case only has one value hence no paramsstruct variable
	user, err := query.GetUser(context.Background(), 1)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(user)


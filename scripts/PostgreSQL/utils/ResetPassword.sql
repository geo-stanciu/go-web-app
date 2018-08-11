DO $$
declare
    _user_id       wmeter.user.user_id%type;
    _user          wmeter.user.loweredusername%type;
    _password      wmeter.user_password.password%type;
    _password_salt wmeter.user_password.password_salt%type;
begin
    _user := lower('admin');

    -- password: Parola1!
    _password      := 'JDJhJDEwJENjSTFlZ2hlWXNXNzRHaVVYcDBpZ08zTTNWQjR6Y3g0WXVLQjFGWHlQZ2UvR0xyc3ZrSzBp';
    _password_salt := '0e016728-0703-452a-aacc-553d5e05c083';

    select user_id
      into _user_id
      from wmeter.user
     where loweredusername = _user;

    update wmeter.user_password p
	   set valid_until = (current_timestamp at time zone 'UTC'),
           valid = 0
     where user_id = _user_id
	   and valid = 1;

    insert into wmeter.user_password (
        user_id,
        password,
        password_salt,
        valid_from
    )
    values (
        _user_id,
        _password,
        _password_salt,
        (current_timestamp at time zone 'UTC')
    );

    UPDATE wmeter.user
       SET last_password_change = (current_timestamp at time zone 'UTC'),
           locked_out = 0,
           failed_password_atmpts = 0
     WHERE user_id = _user_id;

end$$;

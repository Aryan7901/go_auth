package auth_handler

//select queries for auth_handler
var (
	get_otp = /* sql */ `SELECT cast(extract(epoch from otp_expiry) as int) as otp_expiry ,otp FROM user_forgot_password where email=$1`

	get_user = /* sql */ `SELECT is_verified,activation_code,cast(extract(epoch from activation_expiry) as int) as activation_expiry 
	FROM user_auth 
	WHERE username=$1;`
	get_user_details = /* sql */ `SELECT email,password_hash,username,is_verified FROM user_auth where email=$1`

	getEmailPassword = /* sql */ `SELECT email,password_hash FROM user_auth where email=$1`

	check_verification = /* sql */ `
	SELECT is_verified,cast(extract(epoch from activation_expiry) as int) 
	as activation_expiry FROM user_auth where email=$1`
)

//update queries for auth_handler
var (
	update_user = /* sql */ `UPDATE user_auth
	SET is_verified=true,activated_at=current_timestamp
	WHERE username=$1;`

	update_password = /* sql */ `UPDATE user_auth SET password_hash=$1 WHERE email=$2`
)

//insert queries for auth_handler
var (
	user_create = /* sql */ `INSERT INTO user_auth (username, password_hash,email,activation_code) 
	VALUES ($1, $2,$3,$4) ON conflict(email)
	DO UPDATE
	SET activation_expiry=current_timestamp+interval '1 DAY',
	username =excluded.username,activation_code=excluded.activation_code
	where user_auth.is_verified =false;
	`

	create_otp = /* sql */ `INSERT INTO user_forgot_password(email,otp) VALUES ($1,$2) 
	ON conflict(email) where user_forgot_password.otp_expiry < current_timestamp
	DO UPDATE 
	SET otp_expiry=current_timestamp+interval '5 minutes',
	otp=excluded.otp;`
	get_otp_expiry = /* sql */ `SELECT cast(extract(epoch from otp_expiry) as int) as otp_expiry 
	FROM user_forgot_password
	WHERE email=$1;`
)

//delete queries for auth_handler
var (
	delete_password = /* sql */ `DELETE FROM user_forgot_password WHERE email=$1`
)

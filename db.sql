create table user_auth (user_id serial primary key,
	username varchar (50) unique not null,
	password_hash varchar (255) not null,
	phone_number int default null,
	email varchar (255) unique not null,
	is_verified  bool   DEFAULT false,
    activation_code   varchar(255) NOT NULL,
    activation_expiry TIMESTAMP WITH TIME ZONE NOT null default current_timestamp+interval '1 DAY' ,
    activated_at      TIMESTAMP WITH TIME ZONE     DEFAULT NULL,
    created_at        timestamp WITH TIME ZONE    NOT NULL DEFAULT current_timestamp)

create table user_forgot_password(
	email varchar (255) unique not null primary key,
	otp varchar (255) not null,
	otp_expiry TIMESTAMP WITH TIME ZONE NOT null default current_timestamp+interval '5 minutes',
	CONSTRAINT fk_user
      FOREIGN KEY(email) 
	  REFERENCES user_auth(email)
	  ON DELETE CASCADE
);
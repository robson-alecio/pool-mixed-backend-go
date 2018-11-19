BEGIN;

CREATE TABLE poll (
	id uuid NOT NULL PRIMARY KEY,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	name text NOT NULL,
	owner uuid NOT NULL,
	published boolean NOT NULL
);


CREATE TABLE poll_option (
	id uuid NOT NULL PRIMARY KEY,
	poll_id uuid REFERENCES poll(id),
	content text NOT NULL
);


CREATE TABLE poll_vote (
	id uuid NOT NULL PRIMARY KEY,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	poll_id uuid NOT NULL,
	user_id uuid NOT NULL,
	chosen_option text NOT NULL
);


CREATE TABLE poll_session (
	id uuid NOT NULL PRIMARY KEY,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	user_id uuid NOT NULL
);


CREATE TABLE poll_user (
	id uuid NOT NULL PRIMARY KEY,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	login text NOT NULL,
	name text NOT NULL,
	password text NOT NULL
);

COMMIT;
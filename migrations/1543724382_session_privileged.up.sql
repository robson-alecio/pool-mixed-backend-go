--session_privileged up
BEGIN;

alter table poll_session add column registered_user boolean;

COMMIT;
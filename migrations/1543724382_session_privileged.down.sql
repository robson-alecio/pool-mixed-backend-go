--session_privileged down
BEGIN;

alter table poll_session drop column registered_user;

COMMIT;
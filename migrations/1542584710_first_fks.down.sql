BEGIN;

alter table poll_vote 
  drop constraint poll_vote_poll_fk;

alter table poll_vote 
  drop constraint poll_vote_user_fk;

alter table poll
  drop constraint poll_user_fk;

COMMIT;
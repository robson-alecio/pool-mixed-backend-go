BEGIN;

alter table poll
  add constraint poll_user_fk 
  foreign key (owner)
  references poll_user(id);

alter table poll_vote 
  add constraint poll_vote_poll_fk 
  foreign key (poll_id)
  references poll(id);

alter table poll_vote 
  add constraint poll_vote_user_fk 
  foreign key (user_id)
  references poll_user(id);

COMMIT;
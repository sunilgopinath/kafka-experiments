# Setup Postgres

1. Create database
2. Create table user_tracking
3. Check table
4. Check insert

```bash
psql -U postgres
postgres=# CREATE DATABASE user_tracking;
postgres=# \c user_tracking;
user_tracking=# \dt
             List of relations
 Schema |     Name      | Type  |  Owner
--------+---------------+-------+----------
 public | user_activity | table | postgres
(1 row)
user_tracking=# INSERT INTO user_activity (user_id, event_type, timestamp)
VALUES ('u1', 'page_view', 1615123456);
INSERT 0 1
```
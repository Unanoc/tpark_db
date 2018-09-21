DROP INDEX IF EXISTS index_on_users_nickname;
DROP INDEX IF EXISTS index_on_forums_slug;
DROP INDEX IF EXISTS index_on_threads_slug;
DROP INDEX IF EXISTS index_on_threads_id;
DROP INDEX IF EXISTS index_on_posts_id;
DROP INDEX IF EXISTS index_on_votes_nickname_and_thread;

DROP TABLE IF EXISTS "errors";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "forums";
DROP TABLE IF EXISTS "threads";
DROP TABLE IF EXISTS "posts";
DROP TABLE IF EXISTS "votes";
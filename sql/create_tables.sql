DROP TABLE IF EXISTS "errors";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "forums";
DROP TABLE IF EXISTS "threads";
DROP TABLE IF EXISTS "posts";
DROP TABLE IF EXISTS "votes";

-- TABLE "errors" --
CREATE TABLE IF NOT EXISTS errors (
  "message" TEXT
);

-- TABLE "users" --
CREATE TABLE IF NOT EXISTS users (
  "nickname" CITEXT UNIQUE PRIMARY KEY,
  "email"    CITEXT UNIQUE NOT NULL,
  "fullname" CITEXT NOT NULL,
  "about"    TEXT
);

-- TABLE "forums" --
CREATE TABLE IF NOT EXISTS forums (
  "posts"   BIGINT DEFAULT 0,
  "slug"    CITEXT UNIQUE NOT NULL,
  "threads" INTEGER DEFAULT 0,
  "title"   TEXT,
  "user"    CITEXT NOT NULL 
);

-- TABLE "threads" --
CREATE TABLE IF NOT EXISTS threads (
  "id"      SERIAL4 UNIQUE PRIMARY KEY,
  "author"  CITEXT NOT NULL,
  "created" TIMESTAMPTZ(3) DEFAULT now(),
  "forum"   CITEXT,
  "message" TEXT,
  "slug"    CITEXT,
  "title"   TEXT,
  "votes"   INTEGER DEFAULT 0
);

-- TABLE "posts" --
CREATE TABLE IF NOT EXISTS posts (
  "id"       SERIAL8 UNIQUE PRIMARY KEY,
  "author"   CITEXT NOT NULL,
  "created"  TIMESTAMPTZ(3) DEFAULT now(),
  "forum"    CITEXT,
  "isEdited" BOOLEAN DEFAULT FALSE,
  "message"  TEXT NOT NULL,
  "parent"   BIGINT DEFAULT 0,
  "thread"   INTEGER
);

-- TABLE "votes" --
CREATE TABLE IF NOT EXISTS votes (
  "voice"    SMALLINT NOT NULL,
  "nickname" CITEXT NOT NULL,
  "thread"   INTEGER
);


-- INDEX on users "nickname"
CREATE INDEX IF NOT EXISTS index_on_users_nickname
  ON users ("nickname");

-- INDEX on forums "nickname"
CREATE INDEX IF NOT EXISTS index_on_forums_slug
  ON forums ("slug");

-- INDEX on threads "slug"
CREATE INDEX IF NOT EXISTS index_on_threads_slug
  ON threads ("slug");

-- INDEX on threads "id"
CREATE INDEX IF NOT EXISTS index_on_threads_id
  ON threads ("id");

-- INDEX on posts "id"
CREATE INDEX IF NOT EXISTS index_on_posts_id
ON posts ("id");

-- INDEX on votes "thread" and "nickname"
CREATE UNIQUE INDEX index_on_votes_nickname_and_thread 
ON votes ("thread", "nickname");
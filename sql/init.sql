DROP TABLE IF EXISTS "errors" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;

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
  "posts"   BIGINT  DEFAULT 0,
  "slug"    CITEXT  UNIQUE NOT NULL,
  "threads" INTEGER DEFAULT 0,
  "title"   TEXT    NOT NULL,
  "user"    CITEXT  NOT NULL REFERENCES users ("nickname")
);

-- TABLE "threads" --
CREATE TABLE IF NOT EXISTS threads (
  "id"      SERIAL         UNIQUE PRIMARY KEY,
  "author"  CITEXT         NOT NULL REFERENCES users ("nickname"),
  "created" TIMESTAMPTZ(3) DEFAULT now(),
  "forum"   CITEXT         NOT NULL REFERENCES forums ("slug"),
  "message" TEXT           NOT NULL,
  "slug"    CITEXT,
  "title"   TEXT           NOT NULL,
  "votes"   INTEGER        DEFAULT 0
);

-- TABLE "posts" --
CREATE TABLE IF NOT EXISTS posts (
  "id"       SERIAL         UNIQUE PRIMARY KEY,
  "author"   CITEXT         NOT NULL REFERENCES users ("nickname"),
  "created"  TIMESTAMPTZ(3) DEFAULT now(),
  "forum"    CITEXT         NOT NULL REFERENCES forums ("slug"),
  "isEdited" BOOLEAN        DEFAULT FALSE,
  "message"  TEXT           NOT NULL,
  "parent"   INTEGER        DEFAULT 0,
  "thread"   INTEGER        NOT NULL REFERENCES threads ("id"),
  "path"     BIGINT []
);

-- TABLE "votes" --
CREATE TABLE IF NOT EXISTS votes (
  "thread" INT NOT NULL REFERENCES threads("id"),
  "voice"    INTEGER NOT NULL,
  "nickname" CITEXT   NOT NULL
);


-- INDEX on users "nickname"
CREATE INDEX IF NOT EXISTS index_on_users_nickname
  ON users ("nickname");

-- -- INDEX on forums "nickname"
-- CREATE INDEX IF NOT EXISTS index_on_forums_slug
--   ON forums ("slug");

-- -- INDEX on threads "slug"
-- CREATE INDEX IF NOT EXISTS index_on_threads_slug
--   ON threads ("slug");

-- -- INDEX on threads "id"
-- CREATE INDEX IF NOT EXISTS index_on_threads_id
--   ON threads ("id");

-- -- INDEX on posts "id"
-- CREATE INDEX IF NOT EXISTS index_on_posts_id
-- ON posts ("id");
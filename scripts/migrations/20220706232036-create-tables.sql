-- +migrate Up
CREATE TABLE "requests" (
  "id" SERIAL PRIMARY KEY,
  "requester_id" int,
  "created_at" timestamp NOT NULL,
  "request_service" varchar NOT NULL,
  "request_details" text NOT NULL
);

CREATE TABLE "requesters" (
  "id" SERIAL PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar
);

ALTER TABLE "requests" ADD FOREIGN KEY ("requester_id") REFERENCES "requesters" ("id");


-- +migrate Down
set schema 'calidum_rotae';
DROP TABLE "requests";
DROP TABLE "requesters";
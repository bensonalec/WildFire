CREATE TABLE "Tables" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "NumberOfColumns" int
);

CREATE TABLE "Columns" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "Hidden" boolean,
  "tableID" int
);

CREATE TABLE "Business" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "URL" varchar,
  "Email" varchar,
  "City" varchar,
  "Zip" varchar,
  "Address" varchar,
  "Phone" varchar,
  "PointofContactID" int
);

CREATE TABLE "Bus_POC" (
  "ID" SERIAL PRIMARY KEY,
  "First" varchar,
  "Last" varchar,
  "Email" varchar,
  "Phone" varchar,
  "businessID" int
);

CREATE TABLE "Bus_Event" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "Date" varchar,
  "Trainer" varchar,
  "PointofContactID" int,
  "businessID" int
);

CREATE TABLE "Sector" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar
);

CREATE TABLE "busToSec" (
  "ID" SERIAL PRIMARY KEY,
  "businessID" int,
  "sectorID" int
);

CREATE TABLE "Schools" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "Address" varchar,
  "City" varchar,
  "URL" varchar,
  "Email" varchar,
  "Phone" varchar,
  "Attendees" int,
  "pointofContact_id" int
);

CREATE TABLE "School_POC" (
  "ID" SERIAL PRIMARY KEY,
  "First" varchar,
  "Last" varchar,
  "Email" varchar,
  "Phone" varchar,
  "schoolID" int
);

CREATE TABLE "School_Event" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "Date" varchar,
  "Trainer" varchar,
  "PointofContactID" int,
  "schoolID" int
);

CREATE TABLE "Students" (
  "ID" SERIAL PRIMARY KEY,
  "First" varchar,
  "Last" varchar,
  "Email" varchar,
  "schoolID" int
);

CREATE TABLE "stuEvents" (
  "id" SERIAL PRIMARY KEY,
  "studentid" int,
  "eventid" int
);

ALTER TABLE "Columns" ADD FOREIGN KEY ("tableID") REFERENCES "Tables" ("ID");

ALTER TABLE "Business" ADD FOREIGN KEY ("PointofContactID") REFERENCES "Bus_POC" ("ID");

ALTER TABLE "Bus_POC" ADD FOREIGN KEY ("businessID") REFERENCES "Business" ("ID");

ALTER TABLE "Bus_POC" ADD FOREIGN KEY ("ID") REFERENCES "Bus_Event" ("PointofContactID");

ALTER TABLE "Bus_Event" ADD FOREIGN KEY ("businessID") REFERENCES "Business" ("ID");

ALTER TABLE "busToSec" ADD FOREIGN KEY ("businessID") REFERENCES "Business" ("ID");

ALTER TABLE "busToSec" ADD FOREIGN KEY ("sectorID") REFERENCES "Sector" ("ID");

ALTER TABLE "Schools" ADD FOREIGN KEY ("pointofContact_id") REFERENCES "School_POC" ("ID");

ALTER TABLE "School_POC" ADD FOREIGN KEY ("schoolID") REFERENCES "Schools" ("ID");

ALTER TABLE "School_POC" ADD FOREIGN KEY ("ID") REFERENCES "School_Event" ("PointofContactID");

ALTER TABLE "School_Event" ADD FOREIGN KEY ("schoolID") REFERENCES "Schools" ("ID");

ALTER TABLE "Students" ADD FOREIGN KEY ("schoolID") REFERENCES "Schools" ("ID");

ALTER TABLE "stuEvents" ADD FOREIGN KEY ("studentid") REFERENCES "Students" ("ID");

ALTER TABLE "stuEvents" ADD FOREIGN KEY ("eventid") REFERENCES "School_Event" ("ID");

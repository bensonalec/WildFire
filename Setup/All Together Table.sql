CREATE TABLE "Types" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar
);

CREATE TABLE "Tables" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "NumberOfColumns" int,
  "Hidden" boolean,
  "DisplayName" varchar,
  "Type" varchar
);

CREATE TABLE "Columns" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "DisplayName" varchar,
  "Hidden" boolean,
  "tableID" int,
  "Type" varchar,
  "relType" varchar
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
  "PointofContactID" int,
  "Access" varchar
);

CREATE TABLE "Bus_POC" (
  "ID" SERIAL PRIMARY KEY,
  "First" varchar,
  "Last" varchar,
  "Email" varchar,
  "Phone" varchar,
  "businessID" int,
  "Access" varchar
);

CREATE TABLE "Bus_Event" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "Date" varchar,
  "Trainer" varchar,
  "PointofContactID" int,
  "businessID" int,
  "Access" varchar
);

CREATE TABLE "Sector" (
  "ID" SERIAL PRIMARY KEY,
  "Name" varchar,
  "Access" varchar
);

CREATE TABLE "busToSec" (
  "ID" SERIAL PRIMARY KEY,
  "businessID" int,
  "sectorID" int
);

ALTER TABLE "Tables" ADD FOREIGN KEY ("Type") REFERENCES "Types" ("Name");

ALTER TABLE "Columns" ADD FOREIGN KEY ("tableID") REFERENCES "Tables" ("ID");

ALTER TABLE "Business" ADD FOREIGN KEY ("PointofContactID") REFERENCES "Bus_POC" ("ID");

ALTER TABLE "Bus_POC" ADD FOREIGN KEY ("businessID") REFERENCES "Business" ("ID");

ALTER TABLE "Bus_POC" ADD FOREIGN KEY ("ID") REFERENCES "Bus_Event" ("PointofContactID");

ALTER TABLE "Bus_Event" ADD FOREIGN KEY ("businessID") REFERENCES "Business" ("ID");

ALTER TABLE "busToSec" ADD FOREIGN KEY ("businessID") REFERENCES "Business" ("ID");

ALTER TABLE "busToSec" ADD FOREIGN KEY ("sectorID") REFERENCES "Sector" ("ID");

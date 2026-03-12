CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "roles" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(50) UNIQUE NOT NULL,
  "description" TEXT,
  "level" INT NOT NULL DEFAULT 1,
  "is_system_role" BOOLEAN DEFAULT false,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "updated_at" TIMESTAMPTZ DEFAULT (now()),
  "deleted_at" TIMESTAMPTZ
);

CREATE TABLE "permissions" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "resource" VARCHAR(50) NOT NULL,
  "action" VARCHAR(50) NOT NULL,
  "description" TEXT,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "updated_at" TIMESTAMPTZ DEFAULT (now()),
  "deleted_at" TIMESTAMPTZ
);

CREATE TABLE "role_permissions" (
  "role_id" BIGINT NOT NULL,
  "permission_id" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  PRIMARY KEY ("role_id", "permission_id")
);

CREATE TABLE "users" (
  "uid" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "full_name" VARCHAR(255) NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "role_id" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "updated_at" TIMESTAMPTZ DEFAULT (now())
);

CREATE TABLE "airports" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "code" VARCHAR(3) UNIQUE NOT NULL,
  "name" VARCHAR(255),
  "city" VARCHAR(255),
  "country" VARCHAR(255),
  "timezone" VARCHAR(100)
);

CREATE TABLE "aircrafts" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "model" VARCHAR(100),
  "manufacturer" VARCHAR(100),
  "total_seats" INT
);

CREATE TABLE "seat_classes" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "code" VARCHAR(10),
  "name" VARCHAR(50)
);

CREATE TABLE "aircraft_seats" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "aircraft_id" UUID NOT NULL,
  "seat_number" VARCHAR(5) NOT NULL,
  "row_number" INT NOT NULL,
  "seat_letter" CHAR(1),
  "x_position" INT,
  "y_position" INT,
  "seat_class_id" UUID,
  "seat_type" VARCHAR(20),
  "is_exit_row" BOOLEAN DEFAULT false
);

CREATE TABLE "flight_schedules" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "flight_number" VARCHAR(10),
  "departure_airport_id" UUID,
  "arrival_airport_id" UUID,
  "departure_time" TIME,
  "arrival_time" TIME,
  "operating_days" VARCHAR(20)
);

CREATE TABLE "flights" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "schedule_id" UUID,
  "aircraft_id" UUID,
  "departure_time" TIMESTAMPTZ,
  "arrival_time" TIMESTAMPTZ,
  "status" VARCHAR(20)
);

CREATE TABLE "flight_seats" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "flight_id" UUID,
  "aircraft_seat_id" UUID,
  "price" NUMERIC(10,2),
  "status" VARCHAR(20) DEFAULT 'available'
);

CREATE TABLE "pnrs" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "record_locator" VARCHAR(6) UNIQUE NOT NULL,
  "status" VARCHAR(20),
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "ttl" TIMESTAMPTZ
);

CREATE TABLE "pnr_contacts" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "pnr_id" UUID,
  "name" VARCHAR(255),
  "email" VARCHAR(255),
  "phone" VARCHAR(50)
);

CREATE TABLE "pnr_passengers" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "pnr_id" UUID,
  "first_name" VARCHAR(100),
  "last_name" VARCHAR(100),
  "passenger_type" VARCHAR(10),
  "birth_date" DATE,
  "passport_number" VARCHAR(50)
);

CREATE TABLE "pnr_segments" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "pnr_id" UUID,
  "flight_id" UUID,
  "segment_order" INT
);

CREATE TABLE "seat_locks" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "flight_seat_id" UUID,
  "pnr_id" UUID,
  "locked_at" TIMESTAMPTZ,
  "expires_at" TIMESTAMPTZ
);

CREATE TABLE "seat_assignments" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "segment_id" UUID,
  "flight_seat_id" UUID,
  "assigned_at" TIMESTAMPTZ
);

CREATE TABLE "tickets" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "ticket_number" VARCHAR(20) UNIQUE,
  "issued_at" TIMESTAMPTZ
);

CREATE TABLE "ticket_segments" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "ticket_id" UUID,
  "segment_id" UUID
);

CREATE TABLE "ssr_types" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "code" VARCHAR(10),
  "name" VARCHAR(100)
);

CREATE TABLE "passenger_ssr" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "segment_id" UUID,
  "ssr_type_id" UUID
);

CREATE TABLE "meals" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "code" VARCHAR(10),
  "name" VARCHAR(100)
);

CREATE TABLE "passenger_meals" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "segment_id" UUID,
  "meal_id" UUID
);

CREATE TABLE "baggage" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "segment_id" UUID,
  "weight" NUMERIC(5,2),
  "tag_number" VARCHAR(50),
  "status" VARCHAR(20)
);

CREATE TABLE "checkins" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "segment_id" UUID,
  "checkin_time" TIMESTAMPTZ
);

CREATE TABLE "boarding_passes" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "passenger_id" UUID,
  "segment_id" UUID,
  "boarding_group" VARCHAR(10),
  "gate" VARCHAR(10),
  "boarding_time" TIMESTAMPTZ,
  "qr_code" TEXT
);

CREATE TABLE "payments" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "pnr_id" UUID,
  "amount" NUMERIC(10,2),
  "method" VARCHAR(50),
  "status" VARCHAR(20),
  "paid_at" TIMESTAMPTZ
);

CREATE INDEX "idx_roles_level" ON "roles" ("level");

CREATE UNIQUE INDEX "unique_permission_resource_action" ON "permissions" ("resource", "action");

CREATE INDEX "idx_role_permissions_role" ON "role_permissions" ("role_id");

CREATE INDEX "idx_users_role" ON "users" ("role_id");

CREATE INDEX "idx_aircraft_seat_aircraft" ON "aircraft_seats" ("aircraft_id");

CREATE INDEX "idx_schedule_airport" ON "flight_schedules" ("departure_airport_id", "arrival_airport_id");

CREATE INDEX "idx_flight_departure" ON "flights" ("departure_time");

CREATE INDEX "idx_flight_seat_flight" ON "flight_seats" ("flight_id");

CREATE INDEX "idx_flight_seat_status" ON "flight_seats" ("flight_id", "status");

CREATE INDEX "idx_pnr_locator" ON "pnrs" ("record_locator");

CREATE INDEX "idx_passenger_pnr" ON "pnr_passengers" ("pnr_id");

CREATE INDEX "idx_seatlock_expire" ON "seat_locks" ("expires_at");

CREATE UNIQUE INDEX "unique_seat_segment" ON "seat_assignments" ("segment_id", "flight_seat_id");

CREATE INDEX "idx_payment_pnr" ON "payments" ("pnr_id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "aircraft_seats" ADD FOREIGN KEY ("aircraft_id") REFERENCES "aircrafts" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "aircraft_seats" ADD FOREIGN KEY ("seat_class_id") REFERENCES "seat_classes" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "flight_schedules" ADD FOREIGN KEY ("departure_airport_id") REFERENCES "airports" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "flight_schedules" ADD FOREIGN KEY ("arrival_airport_id") REFERENCES "airports" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "flights" ADD FOREIGN KEY ("schedule_id") REFERENCES "flight_schedules" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "flights" ADD FOREIGN KEY ("aircraft_id") REFERENCES "aircrafts" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "flight_seats" ADD FOREIGN KEY ("flight_id") REFERENCES "flights" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "flight_seats" ADD FOREIGN KEY ("aircraft_seat_id") REFERENCES "aircraft_seats" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "pnr_contacts" ADD FOREIGN KEY ("pnr_id") REFERENCES "pnrs" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "pnr_passengers" ADD FOREIGN KEY ("pnr_id") REFERENCES "pnrs" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "pnr_segments" ADD FOREIGN KEY ("pnr_id") REFERENCES "pnrs" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "pnr_segments" ADD FOREIGN KEY ("flight_id") REFERENCES "flights" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "seat_locks" ADD FOREIGN KEY ("flight_seat_id") REFERENCES "flight_seats" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "seat_locks" ADD FOREIGN KEY ("pnr_id") REFERENCES "pnrs" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "seat_assignments" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "seat_assignments" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "seat_assignments" ADD FOREIGN KEY ("flight_seat_id") REFERENCES "flight_seats" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "tickets" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ticket_segments" ADD FOREIGN KEY ("ticket_id") REFERENCES "tickets" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ticket_segments" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "passenger_ssr" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "passenger_ssr" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "passenger_ssr" ADD FOREIGN KEY ("ssr_type_id") REFERENCES "ssr_types" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "passenger_meals" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "passenger_meals" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "passenger_meals" ADD FOREIGN KEY ("meal_id") REFERENCES "meals" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "baggage" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "baggage" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "checkins" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "checkins" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "boarding_passes" ADD FOREIGN KEY ("passenger_id") REFERENCES "pnr_passengers" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "boarding_passes" ADD FOREIGN KEY ("segment_id") REFERENCES "pnr_segments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "payments" ADD FOREIGN KEY ("pnr_id") REFERENCES "pnrs" ("id") DEFERRABLE INITIALLY IMMEDIATE;

CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
)
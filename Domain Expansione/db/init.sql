-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS talent;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS admins;

-- Table: users (customer)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    username VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    password TEXT NOT NULL, -- hashed password (bcrypt)
    profile_picture TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: admins
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL -- hashed password
);

-- Table: talent
CREATE TABLE talent (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    age INT,
    hobbies TEXT,
    city VARCHAR(100),
    profile_picture TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20),
    date DATE NOT NULL,
    package VARCHAR(20) CHECK (package IN ('1_hour', '2_hours', '1_day')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: payment
CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    payment_method VARCHAR(50) DEFAULT 'qr',
    payment_confirmed BOOLEAN DEFAULT FALSE,
    confirmed_at TIMESTAMP
);

-- Insert seed data

-- Admin
INSERT INTO admins (email, password)
VALUES ('admin@renta.app', 'hashed_admin_pass');

-- Users
INSERT INTO users (email, username, phone_number, password, profile_picture)
VALUES 
('john@example.com', 'johnny', '08123456789', 'tes12', 'john.jpg'),
('jane@example.com', 'jane88', '08987654321', 'hashed_pass_2', 'jane.jpg');

-- Talents
INSERT INTO talent (email, name, phone_number, age, hobbies, city, profile_picture)
VALUES
('sarah@renta.app', 'Sarah', '0811223344', 24, 'Cat Lover', 'Jakarta', 'sarah.jpg'),
('nadia@renta.app', 'Nadia', '0811223344', 24, 'Gamer', 'Jakarta', 'sarah.jpg'),
('lina@renta.app', 'Lina', '0822334455', 26, 'Storyteller', 'Bandung', 'lina.jpg');

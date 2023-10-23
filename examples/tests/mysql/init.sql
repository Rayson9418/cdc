CREATE DATABASE IF NOT EXISTS demo;

USE demo;

-- Creating table demo1
CREATE TABLE demo1 (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_time DATETIME,
    event_name VARCHAR(255),
    event_desc TEXT
);

-- Creating table demo2
CREATE TABLE demo2 (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    age INT,
    email VARCHAR(255),
    address TEXT
);

-- Inserting sample data into demo1 table
INSERT INTO demo1 (event_time, event_name, event_desc) VALUES
('2022-01-13 08:00:00', 'Event 1', 'Description of event 1'),
('2023-05-13 09:00:00', 'Event 2', 'Description of event 2'),
('2023-09-13 10:00:00', 'Event 3', 'Description of event 3'),
('2023-10-13 11:00:00', 'Event 4', 'Description of event 4'),
('2023-12-13 12:00:00', 'Event 5', 'Description of event 5');

-- Inserting sample data into demo2 table
INSERT INTO demo2 (name, age, email, address) VALUES
('John Doe', 30, 'john@example.com', '123 Street, City, Country'),
('Jane Smith', 25, 'jane@example.com', '456 Avenue, Town, Country'),
('Michael Johnson', 35, 'michael@example.com', '789 Road, Village, Country'),
('Emily Williams', 28, 'emily@example.com', '1011 Lane, County, Country'),
('Robert Brown', 40, 'robert@example.com', '1213 Drive, State, Country');

create table users (
    id serial PRIMARY KEY,
    username text not null,
    password text not null,
    role text not null
);

INSERT INTO users (username, password, role) VALUES 
('admin', 'super-secret-pass', 'admin'),
('alice', 'qwerty123', 'user'),
('bob', 'password', 'user');

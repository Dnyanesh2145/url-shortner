-- CREATE USER urlshortuser WITH ENCRYPTED PASSWORD 'Qa94nsBVjv4peT4';
-- CREATE DATABASE urlshortener OWNER urlshortuser;
-- \c urlshortener

CREATE TABLE IF NOT EXISTS url_shortner (
    id character varying(10) NOT NULL PRIMARY KEY ,
    url character varying(2048) NOT NULL,
    created_at TIMESTAMP,
    constraint u_constraint unique (url)
);

CREATE TABLE IF NOT EXISTS url_shortner_api_logs (
    id  SERIAL PRIMARY KEY,
    user_id character varying(36),
    request_id character varying(36),
    deviceinfo character varying(256),
    apiendpoint character varying(128),
    ip_address character varying(40),
    geolocation character varying(50),
    httpmethod character varying(10),
    referrer character varying(512),
    responsesize integer,
    responsetime integer,
    statuscode integer,
    statusmessage character varying(256),
    useragent character varying(512),
    created_at TIMESTAMP
);
    
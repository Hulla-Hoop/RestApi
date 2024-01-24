-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
id int NOT NULL,
created_at timestamp with time zone,
updated_at timestamp with time zone,
name text,
surname text,
patronymic text,
age integer,
gender text,
nationality text,
PRIMARY KEY (id)
);

CREATE SEQUENCE public.user_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

ALTER SEQUENCE public.user_id_seq OWNED BY users.id;

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users ;
-- +goose StatementEnd

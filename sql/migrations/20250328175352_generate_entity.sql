-- +goose Up
-- +goose StatementBegin
INSERT INTO investing.users(id, name, email)
VALUES ('e0208302-df4a-42a2-9537-e2636c49a203','Kostya', 'kostya@mail.ru');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM investing.users;
-- +goose StatementEnd

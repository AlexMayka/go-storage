-- +goose Up
-- +goose StatementBegin
CREATE TABLE folders (
    id UUID primary key DEFAULT gen_random_uuid(),
    name varchar(255) not null,
    path varchar(255) unique not null,
    parent_id UUID,
    company_id UUID not null,
    user_create UUID not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
    is_active bool not null default true,
    foreign key (company_id) references companies(id) on delete cascade,
    foreign key (user_create) references users(id) on delete restrict on update cascade,
    foreign key (parent_id) references folders(id) on delete cascade
);

create index idx_folders_company_id on folders(company_id);
create index idx_folders_parent_id on folders(parent_id);

alter table folders add constraint unique_folder_name unique (parent_id, name, company_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_folders_company_id;
DROP INDEX IF EXISTS idx_folders_parent_id;
ALTER table folders DROP CONSTRAINT IF EXISTS unique_folder_name;
DROP TABLE IF EXISTS folders;
-- +goose StatementEnd

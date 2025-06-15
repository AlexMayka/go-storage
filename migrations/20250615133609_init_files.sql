-- +goose Up
-- +goose StatementBegin
create table files (
    id UUID primary key default gen_random_uuid(),
    name varchar(255) not null,
    mime_type varchar(255) not null,
    path varchar(255) unique not null,
    folder_id UUID not null,
    company_id UUID not null,
    hash varchar(255) not null,
    size numeric not null,
    user_created UUID not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
    is_active bool not null default true,
    foreign key (folder_id) references folders(id) on delete cascade,
    foreign key (company_id) references companies(id) on delete cascade,
    foreign key (user_created) references users(id) on delete restrict on update cascade
);

create index idx_files_folder_id on files(folder_id);
create index idx_files_company_id on files(company_id);
create index idx_files_hash on files(hash);

alter table files add constraint unique_file_name unique (folder_id, name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_files_folder_id;
drop index idx_files_company_id;
drop index idx_files_hash;

alter table files drop constraint IF exists unique_file_name;

drop table files;
-- +goose StatementEnd

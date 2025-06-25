-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
    id UUID primary key DEFAULT gen_random_uuid(),
    name varchar(255) not null,
    type varchar(20) not null check (type in ('file', 'folder')),
    full_path varchar(1000) not null,
    parent_id UUID,
    company_id UUID not null,
    user_created UUID not null,

    mime_type varchar(255),
    size BIGINT,
    hash varchar(64),
    storage_path varchar(500),
    
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
    is_active bool not null default true,
    
    foreign key (company_id) references companies(id) on delete cascade,
    foreign key (user_created) references users(id) on delete restrict on update cascade,
    foreign key (parent_id) references files(id) on delete cascade
);

create index idx_files_company_active on files(company_id, is_active);
create index idx_files_parent_type on files(parent_id, type) where is_active = true;
create index idx_files_full_path on files(full_path varchar_pattern_ops) where is_active = true;

create unique index idx_unique_name_in_folder 
    on files(company_id, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::UUID), name) 
    where is_active = true;

alter table files add constraint check_file_fields
check (
    (type = 'folder' and mime_type is null and size is null and hash is null and storage_path is null) or
    (type = 'file' and mime_type is not null and size is not null and hash is not null and storage_path is not null)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_files_company_active;
DROP INDEX IF EXISTS idx_files_parent_type;
DROP INDEX IF EXISTS idx_files_full_path;
DROP INDEX IF EXISTS idx_unique_name_in_folder;
ALTER table files DROP CONSTRAINT IF EXISTS check_file_fields;
DROP TABLE IF EXISTS files;
-- +goose StatementEnd

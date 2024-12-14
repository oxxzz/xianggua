create table yw_books (
    id bigint AUTO_INCREMENT primary key,
    book_id bigint unique comment '书ID',
    name varchar(128) default '' comment '名称',
    yw_cp_id bigint comment '阅文平台渠道ID',
    yw_book_id varchar(64) unique comment '阅文平台书ID',
    status tinyint(1) default 0 comment '状态: 0 pending 1 processing 2 processed 3 failed',
    book_updated_at int default 0 comment '书更新时间',
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment '阅文作品同步记录表';

create index idx_yw_store_book_id on yw_books(book_id);
create index idx_yw_books_yw_cp_id on yw_books(yw_cp_id);
create index idx_yw_books_yw_book_id on yw_books(yw_book_id);

create table yw_chapters(
    id bigint AUTO_INCREMENT primary key,
    book_id bigint comment '书ID',
    chapter_id bigint comment '章节ID',
    name varchar(128) default '' comment '章节名称',
    yw_cp_id bigint comment '阅文平台渠道ID',
    yw_book_id varchar(64) comment '阅文平台书ID',
    yw_chapter_id varchar(64) comment '阅文平台章节ID',
    status tinyint(1) default 0 comment '状态: 0 pending 1 processing 2 processed 3 failed',
    chapter_updated_at int default 0 comment '书更新时间',
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment '阅文作品章节同步记录表';

create index idx_yw_chapter_book_id on yw_chapters(book_id);
create index idx_yw_chapter_yw_cp_id on yw_chapters(yw_cp_id);
create index idx_yw_chapter_yw_book_id on yw_chapters(yw_book_id);
create index idx_yw_chapter_yw_chapter_id on yw_chapters(yw_chapter_id);
create index idx_yw_chapter_chapter_id on yw_chapters(chapter_id);
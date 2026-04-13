-- name: CreateComment :one
insert into comments (created_by, content, commented_on_, commented_on_uuid)
values (@created_by, @content, @commented_on_, @commented_on_uuid)
returning *;

-- name: FindAllCommentsByCommentedOnUuid :many
select *
from comments
where commented_on_ = @commented_on_ 
    and commented_on_uuid = @commented_on_uuid;

-- name: UpdateCommentByUuid :exec
update comments
set content = @content
where uuid = @uuid;

-- name: DeleteCommentByUuid :exec
delete from comments
where uuid = @uuid;

-- name: CreateComment :one
insert into comments (created_by, content, commented_on_, commented_on_uuid)
values (@created_by, @content, @commented_on_, @commented_on_uuid)
returning *;

-- name: FindAllCommentsByCommentedOnUuid :many
select c.*, u.name as user_name
from comments c
join users u on c.created_by = u.uuid
where c.commented_on_ = @commented_on_ 
    and c.commented_on_uuid = @commented_on_uuid;

-- name: UpdateCommentByUuid :exec
update comments
set content = @content
where uuid = @uuid;

-- name: DeleteCommentByUuid :exec
delete from comments
where uuid = @uuid;

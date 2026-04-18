-- name: CreateBoard :one
insert into boards (name, location_name, user_uuid)
values (@name, @location_name, @user_uuid)
returning *;

-- name: GetBoardById :one
select *
from boards
where uuid = @uuid;

-- name: GetAllBoards :many
select *
from boards
where user_uuid = @user_uuid;

-- name: UpdateBoardById :exec
update boards
set name = @name,
    details = @details,
    location_name = @location_name,
    starts_at = @starts_at,
    lasts_until = @lasts_until,
    thumbnail_url = @thumbnail_url
where uuid = @uuid;

-- name: DeleteBoardById :exec
delete from boards
where uuid = @uuid;

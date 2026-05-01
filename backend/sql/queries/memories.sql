-- name: CreateMemory :one
insert into memories (board_uuid, uploaded_by, image_url)
values (@board_uuid, @uploaded_by, @image_url)
returning *;

-- name: GetMemoryByUuid :one
select * from memories
where uuid = @uuid;

-- name: FindAllMemoriesByBoardUuid :many
select * from memories
where board_uuid = @board_uuid
order by created_at desc;

-- name: DeleteMemoryByUuid :exec
delete from memories
where uuid = @uuid;

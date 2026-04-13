-- name: CreateActivity :one
insert into activities (url, title, image_url, board_uuid)
values (@url, @title, @image_url, @board_uuid)
returning *;

-- name: GetActivityByUuid :one
select *
from activities
where uuid = @uuid;

-- name: FindAllActivitiesByBoardUuid :many
select *
from activities
where board_uuid = @board_uuid;

-- name: UpdateActivityByUuid :exec
update activities
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected
where uuid = @uuid;


-- name: DeleteActivityByUuid :exec
delete from activities
where uuid = @uuid;



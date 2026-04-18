-- name: CreateTransport :one
insert into transport (url, title, image_url, board_uuid, user_uuid)
values (@url, @title, @image_url, @board_uuid, @user_uuid)
returning *;

-- name: GetTransportByUuid :one
select *
from transport
where uuid = @uuid;

-- name: FindAllTransportByBoardUuid :many
select *
from transport
where board_uuid = @board_uuid;

-- name: UpdateTransportByUuid :exec
update transport
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected
where uuid = @uuid;


-- name: DeleteTransportByUuid :exec
delete from transport
where uuid = @uuid;

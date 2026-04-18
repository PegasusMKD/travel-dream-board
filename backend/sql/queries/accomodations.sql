-- name: CreateAccomodation :one
insert into accomodations (url, title, image_url, board_uuid, user_uuid)
values (@url, @title, @image_url, @board_uuid, @user_uuid)
returning *;

-- name: GetAccomodationByUuid :one
select *
from accomodations
where uuid = @uuid;

-- name: FindAllAccomodationsByBoardUuid :many
select *
from accomodations
where board_uuid = @board_uuid;

-- name: UpdateAccomodationByUuid :exec
update accomodations
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected
where uuid = @uuid;


-- name: DeleteAccomodationByUuid :exec
delete from accomodations
where uuid = @uuid;



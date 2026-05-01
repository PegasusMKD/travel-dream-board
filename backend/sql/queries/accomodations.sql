-- name: CreateAccomodation :one
insert into accomodations (
    url, title, image_url, board_uuid, user_uuid,
    price, currency, description
)
values (
    @url, @title, @image_url, @board_uuid, @user_uuid,
    @price, @currency, @description
)
returning *;

-- name: GetAccomodationByUuid :one
select a.*,
    coalesce(avg(v.rank)::float, 0)::float as avg_rating,
    count(v.uuid)::int as rating_count
from accomodations a
left join votes v on v.voted_on_uuid = a.uuid and v.voted_on_ = 'accomodation'
where a.uuid = @uuid
group by a.uuid;

-- name: FindAllAccomodationsByBoardUuid :many
select a.*,
    coalesce(avg(v.rank)::float, 0)::float as avg_rating,
    count(v.uuid)::int as rating_count
from accomodations a
left join votes v on v.voted_on_uuid = a.uuid and v.voted_on_ = 'accomodation'
where a.board_uuid = @board_uuid
group by a.uuid;

-- name: UpdateAccomodationByUuid :exec
update accomodations
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected,
    price = @price,
    currency = @currency,
    description = @description
where uuid = @uuid;


-- name: DeleteAccomodationByUuid :exec
delete from accomodations
where uuid = @uuid;

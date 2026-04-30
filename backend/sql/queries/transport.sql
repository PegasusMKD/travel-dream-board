-- name: CreateTransport :one
insert into transport (
    url,
    title,
    image_url,
    board_uuid,
    user_uuid,
    outbound_departing_location,
    outbound_arriving_location,
    outbound_departing_at,
    outbound_arriving_at,
    inbound_departing_location,
    inbound_arriving_location,
    inbound_departing_at,
    inbound_arriving_at,
    outbound_duration_minutes,
    inbound_duration_minutes
)
values (
    @url,
    @title,
    @image_url,
    @board_uuid,
    @user_uuid,
    @outbound_departing_location,
    @outbound_arriving_location,
    @outbound_departing_at,
    @outbound_arriving_at,
    @inbound_departing_location,
    @inbound_arriving_location,
    @inbound_departing_at,
    @inbound_arriving_at,
    @outbound_duration_minutes,
    @inbound_duration_minutes
)
returning *;

-- name: GetTransportByUuid :one
select t.*,
    coalesce(avg(v.rank)::float, 0)::float as avg_rating,
    count(v.uuid)::int as rating_count
from transport t
left join votes v on v.voted_on_uuid = t.uuid and v.voted_on_ = 'transport'
where t.uuid = @uuid
group by t.uuid;

-- name: FindAllTransportByBoardUuid :many
select t.*,
    coalesce(avg(v.rank)::float, 0)::float as avg_rating,
    count(v.uuid)::int as rating_count
from transport t
left join votes v on v.voted_on_uuid = t.uuid and v.voted_on_ = 'transport'
where t.board_uuid = @board_uuid
group by t.uuid;

-- name: UpdateTransportByUuid :exec
update transport
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected,
    outbound_departing_location = @outbound_departing_location,
    outbound_arriving_location = @outbound_arriving_location,
    outbound_departing_at = @outbound_departing_at,
    outbound_arriving_at = @outbound_arriving_at,
    inbound_departing_location = @inbound_departing_location,
    inbound_arriving_location = @inbound_arriving_location,
    inbound_departing_at = @inbound_departing_at,
    inbound_arriving_at = @inbound_arriving_at,
    outbound_duration_minutes = @outbound_duration_minutes,
    inbound_duration_minutes = @inbound_duration_minutes
where uuid = @uuid;


-- name: DeleteTransportByUuid :exec
delete from transport
where uuid = @uuid;

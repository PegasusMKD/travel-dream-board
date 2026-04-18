-- name: CreateScrapeAudit :one
insert into scrape_audit (url, host)
values (@url, @host)
returning *;

-- name: UpdateScrapeAuditByUuid :exec
update scrape_audit
set status = @status,
    title = @title,
    image_url = @image_url,
    description = @description,
    site_name = @site_name
where uuid = @uuid;

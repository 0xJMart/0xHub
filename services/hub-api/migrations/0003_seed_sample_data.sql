-- +goose Up

WITH category_homelab AS (
    INSERT INTO categories (id, name, slug, description, sort_order)
    VALUES (
        gen_random_uuid(),
        'Homelab',
        'homelab',
        'Self-hosted services and infrastructure projects.',
        10
    )
    RETURNING id
),
project_rows AS (
    SELECT
        gen_random_uuid() AS id,
        'Home Lab Inventory' AS title,
        'home-lab-inventory' AS slug,
        'Track hardware, firmware, and network assets in one place.' AS summary,
        'Inventory app that keeps homelab gear organized with lifecycle notes.' AS description,
        'active' AS status,
        (SELECT id FROM category_homelab) AS category_id
    UNION ALL
    SELECT
        gen_random_uuid(),
        'Media Server Stack',
        'media-server-stack',
        'Composable Plex, Jellyfin, and Arr stack with IaC deployment notes.',
        'A reproducible media stack with Terraform modules and compose profiles.',
        'active',
        (SELECT id FROM category_homelab)
),
inserted_projects AS (
    INSERT INTO projects (id, category_id, title, slug, summary, description, status)
    SELECT id, category_id, title, slug, summary, description, status
    FROM project_rows
    RETURNING id, slug
),
tag_rows AS (
    SELECT gen_random_uuid() AS id, 'automation' AS name, 'Automation' AS display_name, '#2563EB' AS color
    UNION ALL
    SELECT gen_random_uuid(), 'infra', 'Infrastructure', '#059669'
    UNION ALL
    SELECT gen_random_uuid(), 'media', 'Media', '#DB2777'
    UNION ALL
    SELECT gen_random_uuid(), 'monitoring', 'Monitoring', '#7C3AED'
),
inserted_tags AS (
    INSERT INTO tags (id, name, display_name, color)
    SELECT id, name, display_name, color
    FROM tag_rows
    RETURNING id, name
),
project_media AS (
    INSERT INTO media_assets (
        id,
        project_id,
        storage_path,
        original_filename,
        mime_type,
        size_bytes,
        description,
        status
    )
    VALUES (
        gen_random_uuid(),
        (SELECT id FROM inserted_projects WHERE slug = 'home-lab-inventory'),
        'projects/home-lab-inventory/dashboard.png',
        'dashboard.png',
        'image/png',
        524288,
        'Screenshot of the asset dashboard.',
        'available'
    )
    RETURNING project_id, id
),
tag_mappings AS (
    SELECT
        (SELECT id FROM inserted_projects WHERE slug = 'home-lab-inventory') AS project_id,
        (SELECT id FROM inserted_tags WHERE name = 'automation') AS tag_id
    UNION ALL
    SELECT
        (SELECT id FROM inserted_projects WHERE slug = 'home-lab-inventory'),
        (SELECT id FROM inserted_tags WHERE name = 'monitoring')
    UNION ALL
    SELECT
        (SELECT id FROM inserted_projects WHERE slug = 'media-server-stack'),
        (SELECT id FROM inserted_tags WHERE name = 'media')
    UNION ALL
    SELECT
        (SELECT id FROM inserted_projects WHERE slug = 'media-server-stack'),
        (SELECT id FROM inserted_tags WHERE name = 'infra')
),
ignored_insert_project_tags AS (
    INSERT INTO project_tags (project_id, tag_id)
    SELECT project_id, tag_id
    FROM tag_mappings
    ON CONFLICT DO NOTHING
    RETURNING project_id
),
ignored_update_projects AS (
    UPDATE projects
    SET hero_media_id = (SELECT id FROM project_media LIMIT 1)
    WHERE slug = 'home-lab-inventory'
    RETURNING id
)
INSERT INTO project_links (project_id, link_type, url, label)
VALUES
(
    (SELECT id FROM inserted_projects WHERE slug = 'home-lab-inventory'),
    'repo',
    'https://github.com/example/homelab-inventory',
    'Source Repository'
),
(
    (SELECT id FROM inserted_projects WHERE slug = 'home-lab-inventory'),
    'docs',
    'https://wiki.example.local/homelab-inventory',
    'Internal Docs'
),
(
    (SELECT id FROM inserted_projects WHERE slug = 'media-server-stack'),
    'repo',
    'https://github.com/example/media-server-stack',
    'Infrastructure as Code'
),
(
    (SELECT id FROM inserted_projects WHERE slug = 'media-server-stack'),
    'guide',
    'https://wiki.example.local/media-stack',
    'Setup Guide'
);

-- +goose Down

DELETE FROM project_links
WHERE project_id IN (
    SELECT id FROM projects WHERE slug IN ('home-lab-inventory', 'media-server-stack')
);

DELETE FROM project_tags
WHERE project_id IN (
    SELECT id FROM projects WHERE slug IN ('home-lab-inventory', 'media-server-stack')
);

DELETE FROM media_assets
WHERE project_id IN (
    SELECT id FROM projects WHERE slug IN ('home-lab-inventory', 'media-server-stack')
);

DELETE FROM projects WHERE slug IN ('home-lab-inventory', 'media-server-stack');
DELETE FROM tags WHERE name IN ('automation', 'infra', 'media', 'monitoring');
DELETE FROM categories WHERE slug = 'homelab';


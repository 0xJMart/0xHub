package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/0xHub/hub-api/internal/domain"
	"github.com/0xHub/hub-api/internal/ports"
)

var (
	_ ports.ProjectRepository = (*Store)(nil)
)

type projectRow struct {
	ID                uuid.UUID          `db:"id"`
	CategoryID        pgtype.UUID        `db:"category_id"`
	Title             string             `db:"title"`
	Slug              string             `db:"slug"`
	Summary           *string            `db:"summary"`
	Description       *string            `db:"description"`
	Status            string             `db:"status"`
	HeroMediaID       pgtype.UUID        `db:"hero_media_id"`
	CreatedAt         time.Time          `db:"created_at"`
	UpdatedAt         time.Time          `db:"updated_at"`
	ArchivedAt        pgtype.Timestamptz `db:"archived_at"`
	CategoryJoinID    pgtype.UUID        `db:"category_join_id"`
	CategoryName      *string            `db:"category_name"`
	CategorySlug      *string            `db:"category_slug"`
	CategoryDesc      *string            `db:"category_description"`
	CategorySort      pgtype.Int4        `db:"category_sort_order"`
	CategoryCreatedAt pgtype.Timestamptz `db:"category_created_at"`
	CategoryUpdatedAt pgtype.Timestamptz `db:"category_updated_at"`
	HeroID            pgtype.UUID        `db:"hero_id"`
	HeroProjectID     pgtype.UUID        `db:"hero_project_id"`
	HeroPath          *string            `db:"hero_storage_path"`
	HeroFilename      *string            `db:"hero_original_filename"`
	HeroMime          *string            `db:"hero_mime_type"`
	HeroSize          pgtype.Int8        `db:"hero_size_bytes"`
	HeroDesc          *string            `db:"hero_description"`
	HeroStatus        *string            `db:"hero_status"`
	HeroCreatedAt     pgtype.Timestamptz `db:"hero_created_at"`
	HeroUpdatedAt     pgtype.Timestamptz `db:"hero_updated_at"`
	HeroExpiresAt     pgtype.Timestamptz `db:"hero_expires_at"`
}

// ListProjects returns paginated project summaries.
func (s *Store) ListProjects(ctx context.Context, filter ports.ProjectFilter) ([]domain.Project, int, error) {
	baseQuery := `
SELECT
    p.id,
    p.category_id,
    p.title,
    p.slug,
    p.summary,
    p.description,
    p.status,
    p.hero_media_id,
    p.created_at,
    p.updated_at,
    p.archived_at,
    c.id AS category_join_id,
    c.name AS category_name,
    c.slug AS category_slug,
    c.description AS category_description,
    c.sort_order AS category_sort_order,
    c.created_at AS category_created_at,
    c.updated_at AS category_updated_at,
    hm.id AS hero_id,
    hm.project_id AS hero_project_id,
    hm.storage_path AS hero_storage_path,
    hm.original_filename AS hero_original_filename,
    hm.mime_type AS hero_mime_type,
    hm.size_bytes AS hero_size_bytes,
    hm.description AS hero_description,
    hm.status AS hero_status,
    hm.created_at AS hero_created_at,
    hm.updated_at AS hero_updated_at,
    hm.expires_at AS hero_expires_at
FROM projects p
LEFT JOIN categories c ON c.id = p.category_id
LEFT JOIN media_assets hm ON hm.id = p.hero_media_id
`

	whereParts := []string{"1=1"}
	args := []any{}
	position := 1

	if cat := strings.TrimSpace(filter.Category); cat != "" {
		whereParts = append(whereParts, fmt.Sprintf("c.slug = $%d", position))
		args = append(args, cat)
		position++
	}

	if tag := strings.TrimSpace(filter.Tag); tag != "" {
		whereParts = append(whereParts, fmt.Sprintf(`
EXISTS (
    SELECT 1
    FROM project_tags pt
    JOIN tags t ON t.id = pt.tag_id
    WHERE pt.project_id = p.id
      AND (LOWER(t.name) = $%[1]d OR LOWER(t.display_name) = $%[1]d)
)`, position))
		args = append(args, strings.ToLower(tag))
		position++
	}

	if search := strings.TrimSpace(filter.Search); search != "" {
		whereParts = append(whereParts, fmt.Sprintf("(p.title ILIKE $%[1]d OR p.summary ILIKE $%[1]d)", position))
		args = append(args, "%"+search+"%")
		position++
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	limitPos := position
	offsetPos := position + 1

	whereClause := strings.Join(whereParts, " AND ")
	query := fmt.Sprintf("%sWHERE %s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", baseQuery, whereClause, limitPos, offsetPos)

	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: list projects query: %w", err)
	}
	defer rows.Close()

	projectRows, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[projectRow])
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: list projects scan: %w", err)
	}

	if len(projectRows) == 0 {
		total, err := s.countProjects(ctx, whereClause, args[:len(args)-2])
		if err != nil {
			return nil, 0, err
		}
		return []domain.Project{}, total, nil
	}

	projects := make([]domain.Project, len(projectRows))
	projectPtrs := make([]*domain.Project, len(projectRows))
	ids := make([]uuid.UUID, len(projectRows))

	for i, row := range projectRows {
		projects[i] = mapProjectRow(row)
		projectPtrs[i] = &projects[i]
		ids[i] = projects[i].ID
	}

	if err := s.attachTags(ctx, ids, projectPtrs); err != nil {
		return nil, 0, err
	}

	total, err := s.countProjects(ctx, whereClause, args[:len(args)-2])
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (s *Store) countProjects(ctx context.Context, whereClause string, args []any) (int, error) {
	query := fmt.Sprintf(`
SELECT COUNT(*)
FROM projects p
LEFT JOIN categories c ON c.id = p.category_id
WHERE %s`, whereClause)
	var total int
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("postgres: count projects: %w", err)
	}
	return total, nil
}

func mapProjectRow(row projectRow) domain.Project {
	project := domain.Project{
		ID:          row.ID,
		CategoryID:  uuidPtr(row.CategoryID),
		Title:       row.Title,
		Slug:        row.Slug,
		Summary:     stringOrDefault(row.Summary, ""),
		Description: stringOrDefault(row.Description, ""),
		Status:      domain.ProjectStatus(row.Status),
		HeroMediaID: uuidPtr(row.HeroMediaID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		ArchivedAt:  timePtr(row.ArchivedAt),
		Tags:        []domain.Tag{},
		Links:       []domain.ProjectLink{},
		Media:       []domain.MediaAsset{},
	}

	if row.CategoryJoinID.Valid {
		project.Category = &domain.Category{
			ID:          uuidFromValue(row.CategoryJoinID),
			Name:        stringOrDefault(row.CategoryName, ""),
			Slug:        stringOrDefault(row.CategorySlug, ""),
			Description: stringOrDefault(row.CategoryDesc, ""),
			SortOrder:   intOrDefault(intPtr(row.CategorySort), 0),
			CreatedAt:   timeOrDefault(timePtr(row.CategoryCreatedAt), time.Time{}),
			UpdatedAt:   timeOrDefault(timePtr(row.CategoryUpdatedAt), time.Time{}),
		}
	}

	if row.HeroID.Valid {
		project.HeroMedia = &domain.MediaAsset{
			ID:               uuidFromValue(row.HeroID),
			ProjectID:        uuidPtr(row.HeroProjectID),
			StoragePath:      stringOrDefault(row.HeroPath, ""),
			OriginalFilename: stringOrDefault(row.HeroFilename, ""),
			MimeType:         stringOrDefault(row.HeroMime, ""),
			SizeBytes:        int64OrDefault(int64Ptr(row.HeroSize), 0),
			Description:      stringOrDefault(row.HeroDesc, ""),
			Status:           mediaStatusFromPtr(row.HeroStatus),
			CreatedAt:        timeOrDefault(timePtr(row.HeroCreatedAt), time.Time{}),
			UpdatedAt:        timeOrDefault(timePtr(row.HeroUpdatedAt), time.Time{}),
			ExpiresAt:        timePtr(row.HeroExpiresAt),
		}
	}

	return project
}

func (s *Store) attachTags(ctx context.Context, ids []uuid.UUID, projects []*domain.Project) error {
	const query = `
SELECT
    pt.project_id,
    t.id,
    t.name,
    t.display_name,
    t.color,
    t.created_at,
    t.updated_at
FROM project_tags pt
JOIN tags t ON t.id = pt.tag_id
WHERE pt.project_id = ANY($1::uuid[])
ORDER BY t.display_name ASC
`

	rows, err := s.pool.Query(ctx, query, ids)
	if err != nil {
		return fmt.Errorf("postgres: query project tags: %w", err)
	}
	defer rows.Close()

	type tagRow struct {
		ProjectID uuid.UUID `db:"project_id"`
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		Display   string    `db:"display_name"`
		Color     *string   `db:"color"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	tagRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[tagRow])
	if err != nil {
		return fmt.Errorf("postgres: scan project tags: %w", err)
	}

	tagMap := make(map[uuid.UUID][]domain.Tag, len(projects))
	for _, row := range tagRows {
		tagMap[row.ProjectID] = append(tagMap[row.ProjectID], domain.Tag{
			ID:          row.ID,
			Name:        row.Name,
			DisplayName: row.Display,
			Color:       stringOrDefault(row.Color, ""),
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	for _, project := range projects {
		if project == nil {
			continue
		}
		if tags, ok := tagMap[project.ID]; ok {
			project.Tags = tags
		}
	}

	return nil
}

// GetProjectBySlug returns a full project aggregate.
func (s *Store) GetProjectBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	project, err := s.fetchProject(ctx, "p.slug = $1", slug)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// GetProjectByID fetches a project aggregate by identifier.
func (s *Store) GetProjectByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	project, err := s.fetchProject(ctx, "p.id = $1", id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *Store) fetchProject(ctx context.Context, predicate string, arg any) (*domain.Project, error) {
	query := fmt.Sprintf(`
SELECT
    p.id,
    p.category_id,
    p.title,
    p.slug,
    p.summary,
    p.description,
    p.status,
    p.hero_media_id,
    p.created_at,
    p.updated_at,
    p.archived_at,
    c.id AS category_join_id,
    c.name AS category_name,
    c.slug AS category_slug,
    c.description AS category_description,
    c.sort_order AS category_sort_order,
    c.created_at AS category_created_at,
    c.updated_at AS category_updated_at,
    hm.id AS hero_id,
    hm.project_id AS hero_project_id,
    hm.storage_path AS hero_storage_path,
    hm.original_filename AS hero_original_filename,
    hm.mime_type AS hero_mime_type,
    hm.size_bytes AS hero_size_bytes,
    hm.description AS hero_description,
    hm.status AS hero_status,
    hm.created_at AS hero_created_at,
    hm.updated_at AS hero_updated_at,
    hm.expires_at AS hero_expires_at
FROM projects p
LEFT JOIN categories c ON c.id = p.category_id
LEFT JOIN media_assets hm ON hm.id = p.hero_media_id
WHERE %s
`, predicate)

	rows, err := s.pool.Query(ctx, query, arg)
	if err != nil {
		return nil, fmt.Errorf("postgres: fetch project query: %w", err)
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[projectRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("project")
		}
		return nil, fmt.Errorf("postgres: fetch project scan: %w", err)
	}

	project := mapProjectRow(record)

	if err := s.attachTags(ctx, []uuid.UUID{project.ID}, []*domain.Project{&project}); err != nil {
		return nil, err
	}

	if err := s.attachLinks(ctx, project.ID, &project); err != nil {
		return nil, err
	}

	if err := s.attachMedia(ctx, project.ID, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *Store) attachLinks(ctx context.Context, projectID uuid.UUID, project *domain.Project) error {
	const query = `
SELECT
    id,
    project_id,
    link_type,
    url,
    label,
    created_at,
    updated_at
FROM project_links
WHERE project_id = $1
ORDER BY created_at ASC
`

	rows, err := s.pool.Query(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("postgres: query project links: %w", err)
	}
	defer rows.Close()

	type linkRow struct {
		ID        uuid.UUID `db:"id"`
		ProjectID uuid.UUID `db:"project_id"`
		LinkType  string    `db:"link_type"`
		URL       string    `db:"url"`
		Label     *string   `db:"label"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	linkRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[linkRow])
	if err != nil {
		return fmt.Errorf("postgres: scan project links: %w", err)
	}

	project.Links = make([]domain.ProjectLink, 0, len(linkRows))
	for _, link := range linkRows {
		project.Links = append(project.Links, domain.ProjectLink{
			ID:        link.ID,
			ProjectID: link.ProjectID,
			LinkType:  link.LinkType,
			URL:       link.URL,
			Label:     stringOrDefault(link.Label, ""),
			CreatedAt: link.CreatedAt,
			UpdatedAt: link.UpdatedAt,
		})
	}

	return nil
}

func (s *Store) attachMedia(ctx context.Context, projectID uuid.UUID, project *domain.Project) error {
	const query = `
SELECT
    id,
    project_id,
    storage_path,
    original_filename,
    mime_type,
    size_bytes,
    description,
    status,
    created_at,
    updated_at,
    expires_at
FROM media_assets
WHERE project_id = $1
ORDER BY created_at ASC
`

	rows, err := s.pool.Query(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("postgres: query media assets: %w", err)
	}
	defer rows.Close()

	type mediaRow struct {
		ID               uuid.UUID          `db:"id"`
		ProjectID        pgtype.UUID        `db:"project_id"`
		StoragePath      string             `db:"storage_path"`
		OriginalFilename string             `db:"original_filename"`
		MimeType         *string            `db:"mime_type"`
		SizeBytes        pgtype.Int8        `db:"size_bytes"`
		Description      *string            `db:"description"`
		Status           string             `db:"status"`
		CreatedAt        time.Time          `db:"created_at"`
		UpdatedAt        time.Time          `db:"updated_at"`
		ExpiresAt        pgtype.Timestamptz `db:"expires_at"`
	}

	mediaRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[mediaRow])
	if err != nil {
		return fmt.Errorf("postgres: scan media assets: %w", err)
	}

	project.Media = make([]domain.MediaAsset, 0, len(mediaRows))
	for _, media := range mediaRows {
		project.Media = append(project.Media, domain.MediaAsset{
			ID:               media.ID,
			ProjectID:        uuidPtr(media.ProjectID),
			StoragePath:      media.StoragePath,
			OriginalFilename: media.OriginalFilename,
			MimeType:         stringOrDefault(media.MimeType, ""),
			SizeBytes:        int64OrDefault(int64Ptr(media.SizeBytes), 0),
			Description:      stringOrDefault(media.Description, ""),
			Status:           domain.MediaStatus(media.Status),
			CreatedAt:        media.CreatedAt,
			UpdatedAt:        media.UpdatedAt,
			ExpiresAt:        timePtr(media.ExpiresAt),
		})
	}

	return nil
}

// CreateProject inserts a project and its relations.
func (s *Store) CreateProject(ctx context.Context, project domain.Project) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres: begin tx: %w", err)
	}
	defer rollbackTx(ctx, tx)

	const insertProject = `
INSERT INTO projects (
    id,
    category_id,
    title,
    slug,
    summary,
    description,
    status,
    hero_media_id,
    archived_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
`

	if _, err := tx.Exec(ctx, insertProject,
		project.ID,
		project.CategoryID,
		project.Title,
		project.Slug,
		nullWhenEmpty(project.Summary),
		nullWhenEmpty(project.Description),
		project.Status,
		project.HeroMediaID,
		project.ArchivedAt,
	); err != nil {
		return fmt.Errorf("postgres: insert project: %w", err)
	}

	if err := upsertProjectTags(ctx, tx, project.ID, project.Tags); err != nil {
		return err
	}

	if err := replaceProjectLinks(ctx, tx, project.ID, project.Links); err != nil {
		return err
	}

	if err := replaceProjectMedia(ctx, tx, project.ID, project.Media); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateProject mutates a project and its relations.
func (s *Store) UpdateProject(ctx context.Context, project domain.Project) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres: begin tx: %w", err)
	}
	defer rollbackTx(ctx, tx)

	const updateProject = `
UPDATE projects
SET
    category_id = $2,
    title = $3,
    slug = $4,
    summary = $5,
    description = $6,
    status = $7,
    hero_media_id = $8,
    archived_at = $9,
    updated_at = NOW()
WHERE id = $1
`

	ct, err := tx.Exec(ctx, updateProject,
		project.ID,
		project.CategoryID,
		project.Title,
		project.Slug,
		nullWhenEmpty(project.Summary),
		nullWhenEmpty(project.Description),
		project.Status,
		project.HeroMediaID,
		project.ArchivedAt,
	)
	if err != nil {
		return fmt.Errorf("postgres: update project: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound("project")
	}

	if err := upsertProjectTags(ctx, tx, project.ID, project.Tags); err != nil {
		return err
	}

	if err := replaceProjectLinks(ctx, tx, project.ID, project.Links); err != nil {
		return err
	}

	if err := replaceProjectMedia(ctx, tx, project.ID, project.Media); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeleteProject removes a project.
func (s *Store) DeleteProject(ctx context.Context, id uuid.UUID) error {
	const deleteQuery = `DELETE FROM projects WHERE id = $1`
	ct, err := s.pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("postgres: delete project: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound("project")
	}
	return nil
}

func upsertProjectTags(ctx context.Context, tx pgx.Tx, projectID uuid.UUID, tags []domain.Tag) error {
	const deleteExisting = `DELETE FROM project_tags WHERE project_id = $1`
	if _, err := tx.Exec(ctx, deleteExisting, projectID); err != nil {
		return fmt.Errorf("postgres: delete project tags: %w", err)
	}

	if len(tags) == 0 {
		return nil
	}

	const insertTag = `
INSERT INTO project_tags (project_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

	for _, tag := range tags {
		if _, err := tx.Exec(ctx, insertTag, projectID, tag.ID); err != nil {
			return fmt.Errorf("postgres: insert project tag: %w", err)
		}
	}
	return nil
}

func replaceProjectLinks(ctx context.Context, tx pgx.Tx, projectID uuid.UUID, links []domain.ProjectLink) error {
	const deleteExisting = `DELETE FROM project_links WHERE project_id = $1`
	if _, err := tx.Exec(ctx, deleteExisting, projectID); err != nil {
		return fmt.Errorf("postgres: delete project links: %w", err)
	}

	if len(links) == 0 {
		return nil
	}

	const insertLink = `
INSERT INTO project_links (
    id,
    project_id,
    link_type,
    url,
    label,
    created_at,
    updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7)
`

	for _, link := range links {
		if _, err := tx.Exec(ctx, insertLink,
			link.ID,
			projectID,
			link.LinkType,
			link.URL,
			nullWhenEmpty(link.Label),
			link.CreatedAt,
			link.UpdatedAt,
		); err != nil {
			return fmt.Errorf("postgres: insert project link: %w", err)
		}
	}
	return nil
}

func replaceProjectMedia(ctx context.Context, tx pgx.Tx, projectID uuid.UUID, media []domain.MediaAsset) error {
	const deleteExisting = `DELETE FROM media_assets WHERE project_id = $1`
	if _, err := tx.Exec(ctx, deleteExisting, projectID); err != nil {
		return fmt.Errorf("postgres: delete media assets: %w", err)
	}

	if len(media) == 0 {
		return nil
	}

	const insertMedia = `
INSERT INTO media_assets (
    id,
    project_id,
    storage_path,
    original_filename,
    mime_type,
    size_bytes,
    description,
    status,
    created_at,
    updated_at,
    expires_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`

	for _, asset := range media {
		if _, err := tx.Exec(ctx, insertMedia,
			asset.ID,
			projectID,
			asset.StoragePath,
			asset.OriginalFilename,
			nullWhenEmpty(asset.MimeType),
			nullWhenZero(asset.SizeBytes),
			nullWhenEmpty(asset.Description),
			asset.Status,
			asset.CreatedAt,
			asset.UpdatedAt,
			asset.ExpiresAt,
		); err != nil {
			return fmt.Errorf("postgres: insert media asset: %w", err)
		}
	}
	return nil
}

func uuidPtr(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	id := uuid.UUID(value.Bytes)
	return &id
}

func uuidFromValue(value pgtype.UUID) uuid.UUID {
	if !value.Valid {
		return uuid.Nil
	}
	return uuid.UUID(value.Bytes)
}

func intPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	val := int(value.Int32)
	return &val
}

func intOrDefault(ptr *int, fallback int) int {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

func int64Ptr(value pgtype.Int8) *int64 {
	if !value.Valid {
		return nil
	}
	val := value.Int64
	return &val
}

func int64OrDefault(ptr *int64, fallback int64) int64 {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

func stringOrDefault(ptr *string, fallback string) string {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func timeOrDefault(ptr *time.Time, fallback time.Time) time.Time {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

func mediaStatusFromPtr(ptr *string) domain.MediaStatus {
	if ptr == nil || *ptr == "" {
		return domain.MediaStatusPending
	}
	return domain.MediaStatus(*ptr)
}

func nullWhenEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullWhenZero(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

// rollbackTx rolls back a transaction ignoring commit/rollback errors.
func rollbackTx(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

package schema

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type VideoDownloader struct {
	ID                  string    `db:"id" json:"id"`
	CrawlerID           *string   `db:"crawler_id" json:"crawler_id"` // self-referencing, can be nil
	VideoName           *string   `db:"video_name" json:"video_name"`
	VideoDownloadSource *string   `db:"video_download_source" json:"video_download_source"`
	VideoDownloadHeader *string   `db:"video_download_header" json:"video_download_header"`
	VideoDownloadPath   *string   `db:"video_download_path" json:"video_download_path"`
	Status              string    `db:"status" json:"status"`
	FileType            *string   `db:"file_type" json:"file_type"`
	FileHash            *string   `db:"file_hash" json:"file_hash"`
	FileSize            *int64    `db:"file_size" json:"file_size"`
	Doc                 string    `db:"doc" json:"doc"`
	IsDeleted           bool      `db:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

func (videoDownloader *VideoDownloader) Create(db *sqlx.DB) error {
	query := `
		INSERT INTO video_downloader (id, video_name, video_download_source, video_download_header, video_download_path, status, crawler_id, doc, is_deleted, file_type, file_hash, file_size)
		VALUES (:id, :video_name, :video_download_source, :video_download_header, :video_download_path, :status, :crawler_id, :doc, :is_deleted, :file_type, :file_hash, :file_size)
	`
	_, err := db.NamedExec(query, videoDownloader)
	return err
}

func (videoDownloader *VideoDownloader) Read(db *sqlx.DB, id string) error {
	query := `SELECT * FROM video_downloader WHERE id = $1`
	return db.Get(videoDownloader, query, id)
}

func (videoDownloader *VideoDownloader) Update(db *sqlx.DB) error {
	query := `
		UPDATE video_downloader SET
			video_name = :video_name, video_download_source = :video_download_source, video_download_header = :video_download_header,
			video_download_path = :video_download_path, status = :status, crawler_id = :crawler_id, doc = :doc, is_deleted = :is_deleted,
			file_type = :file_type, file_hash = :file_hash, file_size = :file_size
		WHERE id = :id
	`
	_, err := db.NamedExec(query, videoDownloader)
	return err
}

func (videoDownloader *VideoDownloader) Delete(db *sqlx.DB) error {
	query := `DELETE FROM video_downloader WHERE id = $1`
	_, err := db.Exec(query, videoDownloader.ID)
	return err
}

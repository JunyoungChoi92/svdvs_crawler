package schema

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type CrawlerTable struct {
	ID             string      `db:"id"`
	Name           string      `db:"name"`
	ParentUrl      string      `db:"parent_url"`
	CurrentUrl     string      `db:"current_url"`
	ExtractedLinks interface{} `db:"extracted_links"`
	Status         string      `db:"status"`
	Doc            string      `db:"doc"`
	IsDeleted      bool        `db:"is_deleted"`
	CreatedAt      interface{} `db:"created_at"`
	UpdatedAt      interface{} `db:"updated_at"`
	HTML           string      `db:"html"`
	ScreenshotPath string      `db:"screenshot_path"`
	SeedDomainID   string      `db:"seed_domain_id"`
}

func BatchInsertCrawlers(db *sqlx.DB, crawlers []*CrawlerTable) error {
	queryBase := `
		INSERT INTO crawler (id, name, parent_url, current_url, extracted_links, status, doc, is_deleted, html, screenshot_path, seed_domain_id)
		VALUES (:id, :name, :parent_url, :current_url, :extracted_links, :status, :doc, :is_deleted, :html, :screenshot_path, :seed_domain_id)
		ON CONFLICT (id) DO NOTHING 
	`

	_, err := db.NamedExec(queryBase, crawlers)
	if err != nil {
		return err
	}
	return nil
}

func ReadCrawler(db *sqlx.DB, crawlerID string) (*CrawlerTable, error) {
	var crawler CrawlerTable

	query := `SELECT * FROM crawler WHERE id = $1`
	err := db.Get(&crawler, query, crawlerID)

	if err != nil {
		return nil, err
	}

	return &crawler, nil
}

func SaveCrawledData(database *sqlx.DB, data CrawlerTable) error {
	// Save the extracted HTML content and screenshot in the PostgreSQL "scrap table".
	tx := database.MustBegin()
	insertSQL := `
		INSERT INTO crawler (
			id, name, parent_url, current_url, extracted_links, status, html, screenshot, seed_domain_id
		) VALUES (:id, :name, :parent_url, :current_url, :extracted_links, :status, :html, :screenshot, :seed_domain_id)
	`
	// if the argument "data" have all the keys in the insertSQL, it will be inserted into the database.

	tx.MustExec(insertSQL, data)

	return tx.Commit()
}

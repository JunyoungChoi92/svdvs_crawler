package schema

import "time"

type SeedDomainTable struct {
	Id               string    `db:"id" json:"id"`
	Host_name        string    `db:"host_name" json:"host_name"`
	Seed_domain      string    `db:"seed_domain" json:"seed_domain"`
	Seed_domain_path string    `db:"seed_domain_path" json:"seed_domain_path"`
	Is_connected     bool      `db:"is_connected" json:"is_connected"`
	Is_deleted       bool      `db:"is_deleted" json:"is_deleted"`
	Created_at       time.Time `db:"created_at" json:"created_at"`
	Updated_at       time.Time `db:"updated_at" json:"updated_at"`
	Doc              string    `db:"doc" json:"doc"`
}

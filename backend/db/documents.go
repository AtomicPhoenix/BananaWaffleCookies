package db

import (
	"context"
	"fmt"
	"time"
)

type Document struct {
	ID               int               `json:"id"`
	UserID           int               `json:"user_id"`
	Title            string            `json:"title"`
	DocumentType     string            `json:"document_type"`
	Tags             []string          `json:"tags"`
	IsArchived       bool              `json:"is_archived"`
	CurrentVersionID int               `json:"current_version_id"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Versions         []DocumentVersion `json:"versions"`
}

type DocumentVersion struct {
	ID            int       `json:"id"`
	DocumentID    int       `json:"document_id"`
	VersionNumber int       `json:"version_number"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileSizeBytes int64     `json:"file_size_bytes"`
	CreatedAt     time.Time `json:"created_at"`
}

func CreateDocument(doc Document) (int, error) {
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(context.Background())

	var docID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO documents (user_id, title, document_type, is_archived, tags)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		doc.UserID,
		doc.Title,
		doc.DocumentType,
		doc.IsArchived,
		doc.Tags,
	).Scan(&docID)

	if err != nil {
		return -1, fmt.Errorf("Failed to insert document for user_id=%d title=%s; Creating document failed: %v\n", doc.UserID, doc.Title, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return -1, fmt.Errorf("Failed to insert document for user_id=%d title=%s; Transaction Failed: %v\n", doc.UserID, doc.Title, err)
	}
	return docID, err
}

func CreateDocumentVersion(docID int, fileName, filePath string, fileSize int64) error {
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var versionID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO document_versions
		 (document_id, version_number, file_name, file_path, file_size_bytes)
		 VALUES (
			$1,
			COALESCE((SELECT MAX(version_number) + 1 FROM document_versions WHERE document_id = $1), 1),
			$2,
			$3,
			$4
		 )
		 RETURNING id`,
		docID,
		fileName,
		filePath,
		fileSize,
	).Scan(&versionID)

	if err != nil {
		return fmt.Errorf("Failed to create document version for document_id=%d; Insert version failed: %v\n",
			docID, err)
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		versionID,
		docID,
	)

	if err != nil {
		return fmt.Errorf("Failed to create document version for document_id=%d; Updating current version failed: %v\n",
			docID, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to create document version for document_id=%d; Transaction failed: %v\n",
			docID, err)
	}

	return nil
}

func DeleteDocument(user_id int, doc_id int) error {
	sql_query := `DELETE FROM documents WHERE id = $1 AND user_id = $2`

	_, err := DbConn.Exec(context.Background(), sql_query, doc_id, user_id)
	if err != nil {
		return fmt.Errorf("Failed to delete document for user_id=%d document_id=%d: %v\n", user_id, doc_id, err)
	}

	return nil
}

func UpdateDocument(userID, docID int, fileName, filePath string, fileSize int64) error {
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var tmp int

	err = tx.QueryRow(
		context.Background(),
		`SELECT id
		 FROM documents
		 WHERE id = $1 AND user_id = $2
		 FOR UPDATE`,
		docID,
		userID,
	).Scan(&tmp)

	if err != nil {
		return fmt.Errorf("Failed to update document for user_id=%d document_id=%d; User does not own document or document does not exist: %v\n",
			userID, docID, err)
	}

	var versionID int

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO document_versions
		 (document_id, version_number, file_name, file_path, file_size_bytes)
		 VALUES (
			$1,
			COALESCE((SELECT MAX(version_number) + 1 FROM document_versions WHERE document_id = $1), 1),
			$2,
			$3,
			$4
		 )
		 RETURNING id`,
		docID,
		fileName,
		filePath,
		fileSize,
	).Scan(&versionID)

	if err != nil {
		return fmt.Errorf("Failed to update document for user_id=%d document_id=%d; Could not insert new version into DB: %v\n",
			userID, docID, err)
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		versionID,
		docID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update document for user_id=%d document_id=%d; Could not update current version: %v\n",
			userID, docID, err)
	}

	return tx.Commit(context.Background())
}

func SetDocumentVersion(userID, docID, versionNumber int) error {
	// Begin transaction
	tx, err := DbConn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var versionID int

	err = tx.QueryRow(
		context.Background(),
		`SELECT dv.id
		 FROM document_versions dv
		 JOIN documents d ON d.id = dv.document_id
		 WHERE dv.document_id = $1
		   AND dv.version_number = $2
		   AND d.user_id = $3`,
		docID,
		versionNumber,
		userID,
	).Scan(&versionID)

	if err != nil {
		return fmt.Errorf("Failed to set document version for user_id=%d document_id=%d; Could not get version id: %v\n", userID, docID, err)
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE documents
		 SET current_version_id = $1, updated_at = NOW()
		 WHERE id = $2`,
		versionID,
		docID,
	)

	if err != nil {
		return fmt.Errorf("Failed to set document version for user_id=%d document_id=%d; Could not set version id: %v\n", userID, docID, err)
	}

	return tx.Commit(context.Background())
}

func GetDocument(docID int, userID int) (Document, error) {
	var doc Document

	err := DbConn.QueryRow(
		context.Background(),
		`SELECT 
			id, user_id, title, document_type, tags, is_archived, current_version_id, created_at, updated_at
		 FROM documents
		 WHERE id = $1 AND user_id = $2`,
		docID,
		userID,
	).Scan(
		&doc.ID,
		&doc.UserID,
		&doc.Title,
		&doc.DocumentType,
		&doc.Tags,
		&doc.IsArchived,
		&doc.CurrentVersionID,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if err != nil {
		return Document{}, fmt.Errorf("Failed to get document for user_id=%d, doc_id=%d; Querying document failed: %v\n", userID, docID, err)
	}

	rows, err := DbConn.Query(
		context.Background(),
		`SELECT id, document_id, version_number, file_name, file_path, file_size_bytes, created_at
		 FROM document_versions
		 WHERE document_id = $1
		 ORDER BY version_number DESC`,
		docID,
	)
	if err != nil {
		return Document{}, fmt.Errorf("Failed to get document for user_id=%d, doc_id=%d; Querying document versions failed: %v\n", userID, docID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var v DocumentVersion

		if err := rows.Scan(
			&v.ID,
			&v.DocumentID,
			&v.VersionNumber,
			&v.FileName,
			&v.FilePath,
			&v.FileSizeBytes,
			&v.CreatedAt,
		); err != nil {
			return Document{}, fmt.Errorf("Failed to get document for user_id=%d, doc_id=%d; Parsing document versions row failed: %v\n", userID, docID, err)
		}

		doc.Versions = append(doc.Versions, v)
	}

	if err = rows.Err(); err != nil {
		return Document{}, fmt.Errorf("Failed to get document for user_id=%d, doc_id=%d; Parsing document versions rows failed: %v\n", userID, docID, err)
	}

	return doc, nil
}

func GetAllDocuments(userID int) ([]Document, error) {
	rows, err := DbConn.Query(
		context.Background(),
		`SELECT id, user_id, title, document_type, tags, is_archived, current_version_id, created_at, updated_at
		 FROM documents
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Query failed: %v\n", userID, err)
	}
	defer rows.Close()

	var results []Document

	for rows.Next() {
		var d Document

		if err := rows.Scan(
			&d.ID,
			&d.UserID,
			&d.Title,
			&d.DocumentType,
			&d.Tags,
			&d.IsArchived,
			&d.CurrentVersionID,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while scanning document rows: %v\n", userID, err)
		}

		// Fetch all versions per doc
		vrows, err := DbConn.Query(
			context.Background(),
			`SELECT id, document_id, version_number, file_name, file_path, file_size_bytes, created_at
			 FROM document_versions
			 WHERE document_id = $1
			 ORDER BY version_number DESC`,
			d.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while querying document version rows: %v\n", userID, err)
		}

		for vrows.Next() {
			var v DocumentVersion
			if err := vrows.Scan(
				&v.ID,
				&v.DocumentID,
				&v.VersionNumber,
				&v.FileName,
				&v.FilePath,
				&v.FileSizeBytes,
				&v.CreatedAt,
			); err != nil {
				vrows.Close()
				return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while scanning document versions row: %v\n", userID, err)
			}
			d.Versions = append(d.Versions, v)
		}
		vrows.Close()

		if err = vrows.Err(); err != nil {
			return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while scanning document version rows: %v\n", userID, err)
		}

		results = append(results, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to get all documents for user_id=%d; Error while scanning document rows: %v\n", userID, err)
	}

	return results, rows.Err()
}

func AssertDocumentOwnership(docID, userID int) error {
	var id int
	return DbConn.QueryRow(
		context.Background(),
		`SELECT id FROM documents WHERE id=$1 AND user_id=$2`,
		docID,
		userID,
	).Scan(&id)
}

func GetNextVersionNumber(docID int) (int, error) {
	var v int
	err := DbConn.QueryRow(
		context.Background(),
		`SELECT COALESCE(MAX(version_number), 0) + 1
		 FROM document_versions
		 WHERE document_id = $1`,
		docID,
	).Scan(&v)

	return v, err
}

func GetLatestDocumentVersion(docID int) (DocumentVersion, error) {
	var v DocumentVersion

	// Attempt to get latest doc version using current_version_id
	err := DbConn.QueryRow(
		context.Background(),
		`SELECT dv.id, dv.document_id, dv.version_number, dv.file_name, dv.file_path, dv.file_size_bytes, dv.created_at
		 FROM document_versions dv
		 JOIN documents d ON d.current_version_id = dv.id
		 WHERE d.id = $1`,
		docID,
	).Scan(
		&v.ID,
		&v.DocumentID,
		&v.VersionNumber,
		&v.FileName,
		&v.FilePath,
		&v.FileSizeBytes,
		&v.CreatedAt,
	)

	if err == nil {
		return v, nil
	}

	// Fallback: Attempt to get latest doc version using highest version number
	err = DbConn.QueryRow(
		context.Background(),
		`SELECT id, document_id, version_number, file_name, file_path, file_size_bytes, created_at
		 FROM document_versions
		 WHERE document_id = $1
		 ORDER BY version_number DESC
		 LIMIT 1`,
		docID,
	).Scan(
		&v.ID,
		&v.DocumentID,
		&v.VersionNumber,
		&v.FileName,
		&v.FilePath,
		&v.FileSizeBytes,
		&v.CreatedAt,
	)

	if err != nil {
		return DocumentVersion{}, fmt.Errorf(
			"Failed to get latest document version for document_id=%d: %v",
			docID, err,
		)
	}

	return v, nil
}

func UpdateDocumentTitle(docID int, userID int, title string) error {
	query := `
		UPDATE documents
		SET title = $1
		WHERE id = $2 AND user_id = $3
	`

	_, err := DbConn.Exec(context.Background(), query, title, docID, userID)
	if err != nil {
		return fmt.Errorf("Failed to rename document for user_id=%d document_id=%d: %v\n", userID, docID, err)
	}

	return nil
}

func SetDocumentArchived(userID, docID int, archived bool) error {
	query := `
		UPDATE documents
		SET is_archived = $1,
		    updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`

	res, err := DbConn.Exec(context.Background(), query, archived, docID, userID)
	if err != nil {
		return fmt.Errorf(
			"Failed to update archive state for user_id=%d document_id=%d: %v",
			userID, docID, err,
		)
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("Failed to set document archival status: No rows updated")
	}

	return err
}

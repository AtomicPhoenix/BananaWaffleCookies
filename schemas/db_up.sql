-- =========================
-- UP SCRIPT
-- =========================

-- USERS
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- PROFILES (1:1 with users)
CREATE TABLE IF NOT EXISTS profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    phone TEXT,
    city TEXT,
    state TEXT,
    country TEXT,
    linkedin_url TEXT,
    portfolio_url TEXT,
    summary TEXT,
    completion_percent INT NOT NULL DEFAULT 0 CHECK (completion_percent BETWEEN 0 AND 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- JOBS
CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name TEXT NOT NULL,
    title TEXT NOT NULL,
    location_text TEXT,
    salary INT,
    status TEXT NOT NULL CHECK (
        status IN ('interested', 'applied', 'interview', 'offer', 'rejected', 'archived')
    ),
    deadline_date DATE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- JOB ACTIVITIES
CREATE TABLE IF NOT EXISTS job_activities (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    activity_type TEXT NOT NULL CHECK (
        activity_type IN (
            'created',
            'updated',
            'status_changed',
            'applied',
            'note_added',
            'document_linked',
            'document_unlinked'
        )
    ),
    activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    description TEXT,
    metadata JSONB
);

-- DOCUMENTS
CREATE TABLE IF NOT EXISTS documents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    document_type TEXT NOT NULL CHECK (
        document_type IN ('resume', 'cover_letter', 'other')
    ),
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- DOCUMENT VERSIONS
CREATE TABLE IF NOT EXISTS document_versions (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version_number INT NOT NULL,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    mime_type TEXT,
    file_size_bytes BIGINT,
    storage_provider TEXT,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_document_version UNIQUE (document_id, version_number)
);

-- JOB <-> DOCUMENT VERSION LINKS
CREATE TABLE IF NOT EXISTS job_document_links (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    document_version_id BIGINT NOT NULL REFERENCES document_versions(id) ON DELETE CASCADE,
    link_type TEXT NOT NULL CHECK (
        link_type IN ('resume', 'cover_letter', 'attachment', 'other')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_job_document_version_link UNIQUE (job_id, document_version_id, link_type)
);

-- INDEXES
CREATE INDEX IF NOT EXISTS idx_jobs_user_id
    ON jobs(user_id);

CREATE INDEX IF NOT EXISTS idx_jobs_user_status
    ON jobs(user_id, status);

CREATE INDEX IF NOT EXISTS idx_documents_user_id
    ON documents(user_id);

CREATE INDEX IF NOT EXISTS idx_document_versions_document_id
    ON document_versions(document_id);

CREATE INDEX IF NOT EXISTS idx_job_activities_job_id
    ON job_activities(job_id);

CREATE INDEX IF NOT EXISTS idx_job_activities_activity_at
    ON job_activities(job_id, activity_at DESC);

CREATE INDEX IF NOT EXISTS idx_job_document_links_job_id
    ON job_document_links(job_id);

CREATE INDEX IF NOT EXISTS idx_job_document_links_document_version_id
    ON job_document_links(document_version_id);

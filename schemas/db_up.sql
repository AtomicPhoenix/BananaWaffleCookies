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
    location TEXT, --added for preferences frontend
    city TEXT,
    state TEXT,
    country TEXT,
    headline TEXT,
    linkedin_url TEXT,
    portfolio_url TEXT,
    summary TEXT,
    preferred_city TEXT,
    preferred_state TEXT,
    preferred_role TEXT,
    preferred_salary_min INT,
    preferred_salary_max INT,
    work_mode TEXT, --added for frontend, was city&state+remote boolean

    completion_percent INT NOT NULL DEFAULT 0
        CHECK (completion_percent BETWEEN 0 AND 100),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (preferred_salary_min IS NULL OR preferred_salary_min >= 0),
    CHECK (preferred_salary_max IS NULL OR preferred_salary_max >= 0),
    CHECK (
        preferred_salary_min IS NULL OR
        preferred_salary_max IS NULL OR
        preferred_salary_max >= preferred_salary_min
    ),
    CHECK (
        preferred_state IS NULL OR preferred_state ~ '^[A-Z]{2}$'
    ),
    CHECK (
        preferred_city IS NULL OR preferred_state IS NOT NULL
    )
);

-- PROFILE EXPERIENCES
CREATE TABLE IF NOT EXISTS profile_experiences (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    experience_type TEXT NOT NULL,
    title TEXT NOT NULL,
    organization TEXT,
    location_text TEXT,
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

-- PROFILE EDUCATION
CREATE TABLE IF NOT EXISTS profile_education (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    institution TEXT NOT NULL,
    degree TEXT,
    field_of_study TEXT,
    start_date DATE,
    end_date DATE,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    honors TEXT NOT NULL DEFAULT '',
    gpa NUMERIC(3,2),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date),
    CHECK (gpa IS NULL OR (gpa >= 0 AND gpa <= 4.00))
);

-- PROFILE SKILLS
CREATE TABLE IF NOT EXISTS profile_skills (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_name TEXT NOT NULL,
    category TEXT,
    proficiency_label TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_profile_skill UNIQUE (user_id, skill_name)
);

-- PROJECTS
CREATE TABLE IF NOT EXISTS profile_projects (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    link TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name TEXT NOT NULL,
    title TEXT NOT NULL,
    location_text TEXT,
    posting_url TEXT,
    salary INT,
    status TEXT NOT NULL CHECK (
        status IN ('interested', 'applied', 'interview', 'offer', 'rejected', 'ghosted')
    ),
    deadline_date DATE,
    last_activity_at TIMESTAMPTZ,
    notes TEXT,
    description TEXT,
    company_notes TEXT,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
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
        'document_unlinked',
        'interview_scheduled',
        'interview_completed',
        'follow_up_created',
        'follow_up_completed',
        'outcome',
        'rejected',
        'ghosted'
    )
    ),
    activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    description TEXT
);

-- INTERVIEWS
CREATE TABLE IF NOT EXISTS interviews (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    round_type TEXT NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- FOLLOW UPS
CREATE TABLE IF NOT EXISTS follow_up_tasks (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    notes TEXT,
    due_at TIMESTAMPTZ,
    remind_at TIMESTAMPTZ,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (completed_at IS NULL OR is_completed = TRUE)
);


-- DOCUMENTS
CREATE TABLE IF NOT EXISTS documents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    document_type TEXT NOT NULL CHECK (
        document_type IN ('resume', 'cover_letter', 'other')
    ),
    tags TEXT[],
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    current_version_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS document_versions (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version_number INT NOT NULL,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size_bytes BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_document_version UNIQUE (document_id, version_number)
);

CREATE UNIQUE INDEX uq_doc_versions_docid_id
    ON document_versions(document_id, id);

ALTER TABLE documents
    ADD CONSTRAINT fk_doc_current_version_match
    FOREIGN KEY (id, current_version_id)
    REFERENCES document_versions(document_id, id);

-- JOB <-> DOCUMENT LINKS
CREATE TABLE IF NOT EXISTS document_links (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    link_type TEXT NOT NULL CHECK (
        link_type IN ('resume', 'cover_letter', 'attachment', 'other')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_job_document_link UNIQUE (job_id, document_id, link_type)
);

-- INDEXES
CREATE INDEX IF NOT EXISTS idx_profiles_user_id
    ON profiles(user_id);

CREATE INDEX IF NOT EXISTS idx_profile_experiences_user_id
    ON profile_experiences(user_id);

CREATE INDEX IF NOT EXISTS idx_profile_experiences_user_sort
    ON profile_experiences(user_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_profile_education_user_id
    ON profile_education(user_id);

CREATE INDEX IF NOT EXISTS idx_profile_education_user_sort
    ON profile_education(user_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_profile_skills_user_id
    ON profile_skills(user_id);

CREATE INDEX IF NOT EXISTS idx_profile_skills_user_sort
    ON profile_skills(user_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_profile_projects_user_sort
    ON profile_projects(user_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_jobs_user_id
    ON jobs(user_id);

CREATE INDEX IF NOT EXISTS idx_jobs_user_status
    ON jobs(user_id, status);

CREATE INDEX IF NOT EXISTS idx_documents_user_id
    ON documents(user_id);

CREATE INDEX IF NOT EXISTS idx_job_activities_job_id
    ON job_activities(job_id);

CREATE INDEX IF NOT EXISTS idx_job_activities_activity_at
    ON job_activities(job_id, activity_at DESC);

CREATE INDEX IF NOT EXISTS idx_interviews_job_id
    ON interviews(job_id);

CREATE INDEX IF NOT EXISTS idx_interviews_job_scheduled_at
    ON interviews(job_id, scheduled_at DESC);

CREATE INDEX IF NOT EXISTS idx_follow_up_tasks_job_id
    ON follow_up_tasks(job_id);

CREATE INDEX IF NOT EXISTS idx_follow_up_tasks_due_at
    ON follow_up_tasks(job_id, due_at);

CREATE INDEX IF NOT EXISTS idx_document_links_job_id
    ON document_links(job_id);

CREATE INDEX IF NOT EXISTS idx_doc_versions_doc_id_version
    ON document_versions(document_id, version_number DESC);

CREATE INDEX IF NOT EXISTS idx_documents_current_version_id
    ON documents(current_version_id);

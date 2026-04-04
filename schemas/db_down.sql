-- =========================
-- DOWN SCRIPT
-- =========================

-- Drop child tables first (reverse dependency order)

DROP TABLE IF EXISTS job_document_links;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS job_activities;
DROP TABLE IF EXISTS jobs;

DROP TABLE IF EXISTS profile_skills;
DROP TABLE IF EXISTS profile_education;
DROP TABLE IF EXISTS profile_experiences;

DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;
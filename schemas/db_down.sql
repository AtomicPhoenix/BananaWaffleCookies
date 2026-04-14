-- =========================
-- DOWN SCRIPT
-- =========================

-- Drop child tables first (reverse dependency order)

-- Job-related children
DROP TABLE IF EXISTS job_document_links;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;

DROP TABLE IF EXISTS interviews;
DROP TABLE IF EXISTS follow_up_tasks;

DROP TABLE IF EXISTS job_activities;
DROP TABLE IF EXISTS jobs;

-- Profile-related children
DROP TABLE IF EXISTS profile_skills;
DROP TABLE IF EXISTS profile_education;
DROP TABLE IF EXISTS profile_experiences;
DROP TABLE IF EXISTS profile_projects;

-- Parent tables
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;
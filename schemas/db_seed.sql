-- =========================
-- SEED DATA
-- =========================

-- USERS
INSERT INTO users (id, email, password_hash)
VALUES
    (1, 'demo@bwc.com', '$2a$12$rtHaGD5uJwlqtybSQO34fueIf/Fn64wgP.8ZkZmYTHYZorpa0z5mC')
ON CONFLICT (id) DO NOTHING;

-- PROFILES
INSERT INTO profiles (
    user_id,
    first_name,
    last_name,
    phone,
    city,
    state,
    country,
    headline,
    linkedin_url,
    portfolio_url,
    summary,
    preferred_role,
    preferred_salary_min,
    preferred_salary_max,
    preferred_city,
    preferred_state,
    completion_percent
)
VALUES
    (
        1,
        'Peter',
        'Griffin',
        '123-555-4567',
        'Quahog',
        'RI',
        'USA',
        'Family Guy',
        'https://linkedin.com/in/petergriffin',
        'https://petergriffin.dev',
        'Highly experienced problem solver with a proven track record in chaotic environments.',
        'Software Engineer',
        70000,
        100000,
        'Quahog',
        'RI',
        85
    )
ON CONFLICT (user_id) DO NOTHING;

-- PROFILE EXPERIENCES (at least 2)
INSERT INTO profile_experiences (
    user_id,
    experience_type,
    title,
    organization,
    location_text,
    start_date,
    end_date,
    is_current,
    description,
    sort_order
)
VALUES
    (
        1,
        'employment',
        'Safety Inspector',
        'Pawtucket Brewery',
        'Quahog, RI',
        '2018-01-01',
        '2018-01-02',
        FALSE,
        'Responsible for ensuring safety standards while occasionally causing large-scale workplace incidents.',
        1
    ),
    (
        1,
        'employment',
        'Freelancer',
        'Self-Employed',
        'Remote / RI',
        '2015-01-01',
        '2017-12-30',
        FALSE,
        'Handled a wide variety of unpredictable scenarios.',
        2
    );

-- PROFILE EDUCATION (at least 2)
INSERT INTO profile_education (
    user_id,
    institution,
    degree,
    field_of_study,
    start_date,
    end_date,
    is_current,
    honors,
    gpa,
    sort_order
)
VALUES
    (
        1,
        'Quahog Community College',
        'Associate of Applied Science',
        'Industrial Safety & Brewing Oversight',
        '2005-09-01',
        '2007-05-15',
        FALSE,
        NULL,
        2.75,
        1
    ),
    (
        1,
        'James Woods Regional High School',
        'High School Diploma',
        'General Studies',
        '1990-09-01',
        '1994-06-01',
        FALSE,
        'Voted Most Likely to Succeed',
        2.75,
        2
    );

-- PROFILE SKILLS (at least 2)
INSERT INTO profile_skills (
    user_id,
    skill_name,
    category,
    proficiency_label,
    sort_order
)
VALUES
    (1, 'Improvisation Under Pressure', 'Soft Skills', 'Expert', 1),
    (1, 'Systems Troubleshooting', 'Technical', 'Intermediate', 2),
    (1, 'Quality Assurance', 'Food & Beverage', 'Advanced', 3),
    (1, 'Linux', 'Systems', 'Intermediate', 4)
ON CONFLICT (user_id, skill_name) DO NOTHING;

-- PROFILE PROJECTS (optional, but useful)
INSERT INTO profile_projects (
    user_id,
    title,
    description,
    link,
    sort_order
)
VALUES
    (
        1,
        'Beverage Consumption Optimization Model',
        'Designed and executed a personal optimization strategy to maximize efficiency and enjoyment of consumption.',
        'https://example.com/beer-optimization',
        1
    ),
    (
        1,
        'Home Infrastructure Stress Testing',
        'Performed extensive real-world stress testing on residential structures.',
        'https://example.com/home-destruction',
        2
    );

-- JOBS (8 spanning multiple stages)
INSERT INTO jobs (
    id,
    user_id,
    company_name,
    title,
    location_text,
    posting_url,
    salary,
    status,
    deadline_date,
    last_activity_at,
    notes,
    description,
    is_archived
)
VALUES
    (
        1,
        1,
        'OpenAI',
        'Software Engineer Intern',
        'Remote',
        'https://example.com/openai-se-intern',
        95000,
        'interested',
        '2026-05-01',
        NOW() - INTERVAL '15 days',
        'Need to tailor resume toward backend experience.',
        'Internship focused on backend systems and developer tooling.',
        FALSE
    ),
    (
        2,
        1,
        'Google',
        'IT Support Specialist',
        'New York, NY',
        'https://example.com/google-it-support',
        80000,
        'applied',
        '2026-04-20',
        NOW() - INTERVAL '10 days',
        'Applied through referral portal.',
        'Support role involving endpoint deployment and troubleshooting.',
        FALSE
    ),
    (
        3,
        1,
        'Microsoft',
        'Security Analyst',
        'Redmond, WA',
        'https://example.com/microsoft-security-analyst',
        105000,
        'interview',
        '2026-04-28',
        NOW() - INTERVAL '4 days',
        'Technical screening completed.',
        'Entry-level analyst role focused on defensive security operations.',
        FALSE
    ),
    (
        4,
        1,
        'Amazon',
        'Cloud Support Associate',
        'Remote',
        'https://example.com/amazon-cloud-support',
        85000,
        'offer',
        '2026-04-18',
        NOW() - INTERVAL '2 days',
        'Received verbal offer, waiting on written details.',
        'Cloud support role working with customer infrastructure issues.',
        FALSE
    ),
    (
        5,
        1,
        'CrowdStrike',
        'SOC Analyst I',
        'Austin, TX',
        'https://example.com/crowdstrike-soc',
        90000,
        'rejected',
        '2026-04-10',
        NOW() - INTERVAL '8 days',
        'Good interview practice. Rejected after final round.',
        'Security operations center analyst position.',
        FALSE
    ),
    (
        6,
        1,
        'GitHub',
        'Junior Backend Engineer',
        'Remote',
        'https://example.com/github-backend',
        110000,
        'archived',
        '2026-03-30',
        NOW() - INTERVAL '30 days',
        'Archived because role no longer aligns with current goals.',
        'Backend engineering role in internal platform systems.',
        TRUE
    ),
    (
        7,
        1,
        'Datadog',
        'Technical Support Engineer',
        'Boston, MA',
        'https://example.com/datadog-support',
        88000,
        'applied',
        '2026-04-25',
        NOW() - INTERVAL '6 days',
        'Need follow-up if no response by next week.',
        'Technical customer support for monitoring and cloud tooling.',
        FALSE
    ),
    (
        8,
        1,
        'Palo Alto Networks',
        'Associate Security Engineer',
        'Santa Clara, CA',
        'https://example.com/paloalto-security',
        102000,
        'interview',
        '2026-04-30',
        NOW() - INTERVAL '1 day',
        'Behavioral interview scheduled.',
        'Associate-level security engineering role.',
        FALSE
    )
ON CONFLICT (id) DO NOTHING;

-- JOB ACTIVITIES (optional but very useful for timeline demos)
INSERT INTO job_activities (job_id, activity_type, activity_at, description, metadata)
VALUES
    (1, 'created', NOW() - INTERVAL '20 days', 'Job saved to tracker.', NULL),
    (1, 'updated', NOW() - INTERVAL '15 days', 'Updated notes and target resume.', NULL),

    (2, 'created', NOW() - INTERVAL '12 days', 'Job saved to tracker.', NULL),
    (2, 'applied', NOW() - INTERVAL '10 days', 'Applied through referral portal.', NULL),

    (3, 'created', NOW() - INTERVAL '14 days', 'Job saved to tracker.', NULL),
    (3, 'applied', NOW() - INTERVAL '11 days', 'Submitted application.', NULL),
    (3, 'status_changed', NOW() - INTERVAL '7 days', 'Status changed to interview.', '{"new_status":"interview"}'),

    (4, 'created', NOW() - INTERVAL '16 days', 'Job saved to tracker.', NULL),
    (4, 'applied', NOW() - INTERVAL '13 days', 'Applied on company website.', NULL),
    (4, 'status_changed', NOW() - INTERVAL '5 days', 'Status changed to offer.', '{"new_status":"offer"}'),
    (4, 'outcome', NOW() - INTERVAL '2 days', 'Received verbal offer.', NULL),

    (5, 'created', NOW() - INTERVAL '18 days', 'Job saved to tracker.', NULL),
    (5, 'applied', NOW() - INTERVAL '15 days', 'Application submitted.', NULL),
    (5, 'status_changed', NOW() - INTERVAL '6 days', 'Status changed to rejected.', '{"new_status":"rejected"}'),
    (5, 'outcome', NOW() - INTERVAL '6 days', 'Rejected after final round.', NULL),

    (6, 'created', NOW() - INTERVAL '40 days', 'Job saved to tracker.', NULL),
    (6, 'status_changed', NOW() - INTERVAL '30 days', 'Status changed to archived.', '{"new_status":"archived"}'),

    (7, 'created', NOW() - INTERVAL '8 days', 'Job saved to tracker.', NULL),
    (7, 'applied', NOW() - INTERVAL '6 days', 'Applied through company portal.', NULL),
    (7, 'follow_up_created', NOW() - INTERVAL '2 days', 'Created follow-up reminder.', NULL),

    (8, 'created', NOW() - INTERVAL '5 days', 'Job saved to tracker.', NULL),
    (8, 'applied', NOW() - INTERVAL '4 days', 'Applied on careers page.', NULL),
    (8, 'interview_scheduled', NOW() - INTERVAL '1 day', 'Behavioral interview scheduled.', NULL);

-- INTERVIEWS
INSERT INTO interviews (
    job_id,
    round_type,
    scheduled_at,
    completed_at,
    notes
)
VALUES
    (
        3,
        'technical_screen',
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '4 days',
        'Covered SQL, APIs, and troubleshooting scenarios.'
    ),
    (
        8,
        'behavioral',
        NOW() + INTERVAL '2 days',
        NULL,
        'Prepare STAR stories and security project examples.'
    );

-- FOLLOW-UP TASKS
INSERT INTO follow_up_tasks (
    job_id,
    title,
    notes,
    due_at,
    remind_at,
    is_completed,
    completed_at
)
VALUES
    (
        2,
        'Send follow-up email',
        'Follow up with recruiter if no update arrives.',
        NOW() + INTERVAL '2 days',
        NOW() + INTERVAL '1 day',
        FALSE,
        NULL
    ),
    (
        7,
        'Check application status',
        'Reach out if application remains pending.',
        NOW() + INTERVAL '3 days',
        NOW() + INTERVAL '2 days',
        FALSE,
        NULL
    );
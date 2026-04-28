package ai

import (
	"context"
	"fmt"
	"strings"

	"bananawafflecookies.com/m/v2/db"
	"bananawafflecookies.com/m/v2/settings"
	"google.golang.org/genai"
)

func GenerateResumeDraft(job db.Job, profile db.Profile) (string, error) {
	fullName := strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	// Load structured data
	experiences, _ := db.GetProfileExperiences(profile.UserID)
	skills, _ := db.GetProfileSkills(profile.UserID)
	education, _ := db.GetProfileEducation(profile.UserID)

	// Format experiences
	var expText strings.Builder
	for _, e := range experiences {
		expText.WriteString(fmt.Sprintf(
			"- %s at %s (%s): %s\n",
			e.Title,
			e.Organization,
			e.ExperienceType,
			e.Description,
		))
	}

	// Format skills
	var skillText strings.Builder
	for _, s := range skills {
		skillText.WriteString(fmt.Sprintf(
			"- %s (%s)\n",
			s.SkillName,
			s.ProficiencyLabel,
		))
	}

	// Format education
	var eduText strings.Builder
	for _, e := range education {
		eduText.WriteString(fmt.Sprintf(
			"- %s, %s in %s (GPA: %.2f)\n",
			e.Institution,
			e.Degree,
			e.FieldOfStudy,
			e.Gpa,
		))
	}

	// Build prompt
	query := fmt.Sprintf(`
Create a tailored, professional resume for the job below.

STRICT RULES:
- Use the provided data only; Do NOT invent facts
- Keep the resume concise and reviewer friendly
- Optimize for the job
- Your output should have the following format
	1. Professional Summary
	2. Education (omit if missing)
	3. Work Experience (omit if missing; tailor to job)
	4. Projects (omit if missing / not applicable to job)
	5. Skills (omit if missing / not applicable)

------------------------
CANDIDATE PROFILE
------------------------
Name: %s
Headline: %s
Location: %s, %s, %s
Phone: %s
LinkedIn: %s
Portfolio: %s
Summary: %s

Preferred Role: %s
Preferred Work Mode: %s
Preferred Salary: %d - %d

------------------------
EXPERIENCE
------------------------
%s

------------------------
SKILLS
------------------------
%s

------------------------
EDUCATION
------------------------
%s

------------------------
JOB
------------------------
Company: %s
Title: %s
Location: %s
Salary: %d
Status: %s
Description:
%s
`,
		fullName,
		profile.Headline,
		profile.City,
		profile.State,
		profile.Country,
		profile.Phone,
		profile.LinkedinURL,
		profile.PortfolioURL,
		profile.Summary,
		profile.PreferredRole,
		profile.WorkMode,
		profile.PreferredSalaryMin,
		profile.PreferredSalaryMax,
		expText.String(),
		skillText.String(),
		eduText.String(),
		job.CompanyName,
		job.Title,
		job.LocationText,
		job.Salary,
		job.Status,
		job.Description,
	)

	return queryModel(query)
}

func GenerateCoverLetter(job db.Job, profile db.Profile) (string, error) {
	fullName := strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	// Load structured data
	experiences, _ := db.GetProfileExperiences(profile.UserID)
	skills, _ := db.GetProfileSkills(profile.UserID)
	education, _ := db.GetProfileEducation(profile.UserID)

	// Format experiences
	var expText strings.Builder
	for _, e := range experiences {
		expText.WriteString(fmt.Sprintf(
			"- %s at %s (%s): %s\n",
			e.Title,
			e.Organization,
			e.ExperienceType,
			e.Description,
		))
	}

	// Format skills
	var skillText strings.Builder
	for _, s := range skills {
		skillText.WriteString(fmt.Sprintf(
			"- %s (%s)\n",
			s.SkillName,
			s.ProficiencyLabel,
		))
	}

	// Format education
	var eduText strings.Builder
	for _, e := range education {
		eduText.WriteString(fmt.Sprintf(
			"- %s, %s in %s (GPA: %.2f)\n",
			e.Institution,
			e.Degree,
			e.FieldOfStudy,
			e.Gpa,
		))
	}

	// Build prompt for cover letter
	query := fmt.Sprintf(`
Write a professional, tailored cover letter for the job below.

STRICT RULES:
- Use ONLY the provided information; do NOT invent facts
- Keep tone professional, confident, and natural
- Must be 3–5 short paragraphs
- Do NOT use bullet points
- Must include:
  1. Introduction (who you are + role you're applying for)
  2. Why this company / role
  3. Relevant experience & skills
  4. Closing + call to action
- Tailor content specifically to the job description

------------------------
CANDIDATE PROFILE
------------------------
Name: %s
Headline: %s
Location: %s, %s, %s
Phone: %s
LinkedIn: %s
Portfolio: %s
Summary: %s

Preferred Role: %s
Preferred Work Mode: %s

------------------------
EXPERIENCE
------------------------
%s

------------------------
SKILLS
------------------------
%s

------------------------
EDUCATION
------------------------
%s

------------------------
JOB
------------------------
Company: %s
Title: %s
Location: %s
Salary: %d
Status: %s
Description:
%s
`,
		fullName,
		profile.Headline,
		profile.City,
		profile.State,
		profile.Country,
		profile.Phone,
		profile.LinkedinURL,
		profile.PortfolioURL,
		profile.Summary,
		profile.PreferredRole,
		profile.WorkMode,
		expText.String(),
		skillText.String(),
		eduText.String(),
		job.CompanyName,
		job.Title,
		job.LocationText,
		job.Salary,
		job.Status,
		job.Description,
	)

	return queryModel(query)
}

func GenerateJobNotes(job db.Job, profile db.Profile) (string, error) {
	fullName := strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	experiences, _ := db.GetProfileExperiences(profile.UserID)
	skills, _ := db.GetProfileSkills(profile.UserID)

	var expText strings.Builder
	for _, e := range experiences {
		expText.WriteString(fmt.Sprintf("- %s at %s\n", e.Title, e.Organization))
	}

	var skillText strings.Builder
	for _, s := range skills {
		skillText.WriteString(fmt.Sprintf("- %s\n", s.SkillName))
	}

	userQuery := strings.TrimSpace(job.Notes)
	if userQuery == "" {
		userQuery = "General job and company analysis"
	}

	query := fmt.Sprintf(`
You are generating structured notes to help a candidate evaluate BOTH the job and the company.

IMPORTANT:
- The user provided a custom query/focus — prioritize answering it
- Include BOTH job insights AND company insights
- If company info is limited, infer cautiously from job description and say "Not specified" when needed

USER QUERY / FOCUS:
%s

STRICT RULES:
- Be concise
- Use bullet points only
- Do NOT invent specific facts about the company
- You MAY infer general patterns (e.g., startup vs enterprise) but label them clearly
- Max 250 words

OUTPUT FORMAT:

1. Direct Answer to User Query
- (Focus specifically on the user's question)

2. Job Insights
- (Key responsibilities, expectations, priorities)

3. Company Insights
- (What kind of company this appears to be: size, culture, industry hints, stability)
- (Any signals from job description: fast-paced, growth stage, etc.)

4. Fit Assessment
- (Candidate vs role + company alignment)

5. Preparation Tips
- (What to study for THIS company + role)

6. Potential Questions
- (Interview questions tailored to company + role)

7. Red Flags (Job or Company)
- (Compensation clarity, vague role, unrealistic expectations, etc.)

------------------------
CANDIDATE
------------------------
Name: %s
Headline: %s
Summary: %s

Experience:
%s

Skills:
%s

------------------------
JOB
------------------------
Company: %s
Title: %s
Location: %s
Salary: %d
Status: %s
Description:
%s
`,
		userQuery,
		fullName,
		profile.Headline,
		profile.Summary,
		expText.String(),
		skillText.String(),
		job.CompanyName,
		job.Title,
		job.LocationText,
		job.Salary,
		job.Status,
		job.Description,
	)

	return queryModel(query)
}

func queryModel(query string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		settings.Logger.Error("Failed to create ai client", "err", err)
		return "ERROR", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(query),
		nil,
	)
	if err != nil {
		settings.Logger.Error("Failed to query ai client", "err", err)
		return "ERROR", err
	}
	return result.Text(), nil
}

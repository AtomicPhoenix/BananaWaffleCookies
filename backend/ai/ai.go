package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"bananawafflecookies.com/m/v2/db"
	"google.golang.org/genai"
)

func GenerateResumeDraft(job db.Job, profile db.Profile) (string, error) {
	fullName := strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	query := fmt.Sprintf(`
Create a tailored, professional resume for the job below.

STRICT RULES:
- Use the provided data only; Do NOT invent facts
- Keep the resume concise and reviewer friendly
- Optimize for the job
- Your output should have the following fields in this order
	- Professional Summary
	- Education (omit if missing)
	- Work Experience (omit if missing; tailor to job)
	- Projects (omit if missing / not applicable to job)
	- Skills (omit if missing / not applicable)

------------------------
PROFILE INFO:
------------------------
Name: %s
Headline: %s
Location (City): %s
Location (State):  %s
Location (Countru): %s
Phone: %s
LinkedIn: %s
Portfolio: %s
Summary: %s

Preferred Role: %s
Preferred Work Mode: %s
Preferred Salary: %d - %d

------------------------
JOB INFO:
------------------------
Company: %s
Title: %s
Location: %s
Salary: %d
Status: %s
Description: %s
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
		job.CompanyName,
		job.Title,
		job.LocationText,
		job.Salary,
		job.Status,
		job.Description,
	)
	fmt.Printf("Query: %s\n---------------------------\n", query)
	return queryModel(query)
}

func queryModel(query string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create ai client: %v", err)
		return "ERROR", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(query),
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query ai: %v", err)
		return "ERROR", err
	}
	return result.Text(), err
}

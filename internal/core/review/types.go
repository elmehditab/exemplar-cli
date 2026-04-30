package review

type ReviewRequest struct {
	RepoPath string
}

type ReviewContext struct {
	RepositoryRoot string
	CurrentBranch  string
	ChangedFiles   []string
	Diff           string
	ParseDiff      ParseDiff
	Warnings       []string
}

type ReviewResult struct {
	RepositoryRoot string
	Message        string
	Status         string
	CurrentBranch  string
	ChangedFiles   []string
	ExecutedStages []string
	Diff           string
	Warnings       []string
	Findings       []Finding
}

type Finding struct {
	ID               string
	Title            string
	Description      string
	Severity         Severity
	Confidence       Confidence
	ConfidenceReason string
	FilePath         string
	Line             int
	Evidence         []Evidence
	Guidance         string
}

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

type Confidence string

const (
	ConfidenceLow    Confidence = "low"
	ConfidenceMedium Confidence = "medium"
	ConfidenceHigh   Confidence = "high"
)

type Evidence struct {
	Kind        string
	Reference   string
	Description string
}

type ParseDiff struct {
	Files []DiffFile
}

type DiffFile struct {
	Path               string
	OldPath            string
	Status             string
	Hunks              []DiffHunk
	AddedCount         int
	DeletedCount       int
	FileClassification string
}

type DiffHunk struct {
	Header string
	Lines  []DiffLine
}

type DiffLine struct {
	Content    string
	OldLineNum int
	NewLineNum int
	Kind       string
}

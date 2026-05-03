package review

type ReviewRequest struct {
	RepoPath string
}

type ReviewContext struct {
	RepositoryRoot string
	CurrentBranch  string
	ChangedFiles   []string
	Diff           string
	ParsedDiff     ParsedDiff
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
	ParsedDiff     ParsedDiff
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

type ParsedDiff struct {
	Files         []DiffFile
	ReviewTargets []ReviewTarget
	Stats         DiffStats
}

type DiffStats struct {
	FilesChanged int
	LinesAdded   int
	LinesDeleted int
	BinaryFiles  int
}

type DiffFile struct {
	OldPath      string
	NewPath      string
	Status       DiffFileStatus
	IsNew        bool
	IsDelete     bool
	IsRename     bool
	IsCopy       bool
	IsBinary     bool
	LinesAdded   int
	LinesDeleted int
	Hunks        []DiffHunk
}

type DiffFileStatus string

const (
	DiffFileStatusModified DiffFileStatus = "modified"
	DiffFileStatusAdded    DiffFileStatus = "added"
	DiffFileStatusDeleted  DiffFileStatus = "deleted"
	DiffFileStatusRenamed  DiffFileStatus = "renamed"
	DiffFileStatusCopied   DiffFileStatus = "copied"
	DiffFileStatusBinary   DiffFileStatus = "binary"
)

type DiffHunk struct {
	Header       string
	OldStart     int
	OldLineSpan  int
	NewStart     int
	NewLineSpan  int
	LinesAdded   int
	LinesDeleted int
	Lines        []DiffLine
}

type DiffLine struct {
	Operation     DiffLineOperation
	Content       string
	OldLineNumber int
	NewLineNumber int
	NoNewline     bool
}

type ReviewTarget struct {
	FilePath      string
	Line          int
	Content       string
	HunkHeader    string
	ContextBefore []DiffLine
	ContextAfter  []DiffLine
}

type DiffLineOperation string

const (
	DiffLineOperationContext DiffLineOperation = "context"
	DiffLineOperationAdded   DiffLineOperation = "added"
	DiffLineOperationDeleted DiffLineOperation = "deleted"
)

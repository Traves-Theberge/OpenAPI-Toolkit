package export

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// HTMLTemplateData contains all data needed for HTML report generation
type HTMLTemplateData struct {
	Timestamp   string
	SpecPath    string
	BaseURL     string
	TotalTests  int
	Passed      int
	Failed      int
	PassRate    float64
	Results     []HTMLResult
	HasVerbose  bool
	TotalTime   string
	AverageTime string
}

// HTMLResult represents a test result with additional display fields
type HTMLResult struct {
	Method       string
	Endpoint     string
	Status       string
	Message      string
	Duration     string
	RetryCount   int    // Number of retries performed
	RowClass     string // CSS class for row styling (success/failure)
	HasLog       bool
	RequestURL   string
	RequestBody  string
	ResponseBody string
	Timestamp    string
}

// htmlTemplate contains the complete HTML structure with embedded CSS
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OpenAPI Test Results - {{.Timestamp}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
            font-weight: 700;
        }
        .header p {
            font-size: 1.1rem;
            opacity: 0.95;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 30px 40px;
            background: #f8f9fa;
            border-bottom: 2px solid #e9ecef;
        }
        .stat-box {
            text-align: center;
            padding: 20px;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .stat-label {
            font-size: 0.9rem;
            color: #6c757d;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 8px;
        }
        .stat-value {
            font-size: 2rem;
            font-weight: bold;
            color: #212529;
        }
        .stat-value.pass {
            color: #28a745;
        }
        .stat-value.fail {
            color: #dc3545;
        }
        .stat-value.rate-good {
            color: #28a745;
        }
        .stat-value.rate-ok {
            color: #ffc107;
        }
        .stat-value.rate-bad {
            color: #dc3545;
        }
        .meta {
            padding: 20px 40px;
            background: #fff;
            border-bottom: 1px solid #e9ecef;
        }
        .meta-row {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #f1f3f5;
        }
        .meta-row:last-child {
            border-bottom: none;
        }
        .meta-label {
            font-weight: 600;
            color: #495057;
        }
        .meta-value {
            color: #6c757d;
            font-family: 'Courier New', monospace;
        }
        .results {
            padding: 40px;
        }
        .results h2 {
            font-size: 1.8rem;
            margin-bottom: 20px;
            color: #212529;
        }
        .results-table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .results-table thead {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .results-table th {
            padding: 15px;
            text-align: left;
            font-weight: 600;
            text-transform: uppercase;
            font-size: 0.85rem;
            letter-spacing: 0.5px;
        }
        .results-table td {
            padding: 12px 15px;
            border-bottom: 1px solid #e9ecef;
        }
        .results-table tbody tr:last-child td {
            border-bottom: none;
        }
        .results-table tbody tr.success {
            background: #f8fff9;
        }
        .results-table tbody tr.success:hover {
            background: #e8f5e9;
        }
        .results-table tbody tr.failure {
            background: #fff5f5;
        }
        .results-table tbody tr.failure:hover {
            background: #ffebee;
        }
        .method {
            display: inline-block;
            padding: 4px 10px;
            border-radius: 4px;
            font-weight: 600;
            font-size: 0.85rem;
            font-family: 'Courier New', monospace;
        }
        .method-get {
            background: #d1ecf1;
            color: #0c5460;
        }
        .method-post {
            background: #d4edda;
            color: #155724;
        }
        .method-put {
            background: #fff3cd;
            color: #856404;
        }
        .method-patch {
            background: #e7d4f5;
            color: #5a1e7a;
        }
        .method-delete {
            background: #f8d7da;
            color: #721c24;
        }
        .endpoint {
            font-family: 'Courier New', monospace;
            color: #495057;
            font-size: 0.95rem;
        }
        .status {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-weight: 600;
            font-size: 0.85rem;
        }
        .status-success {
            background: #28a745;
            color: white;
        }
        .status-error {
            background: #dc3545;
            color: white;
        }
        .status-warning {
            background: #ffc107;
            color: #333;
        }
        .message {
            color: #6c757d;
            font-size: 0.9rem;
        }
        .duration {
            font-family: 'Courier New', monospace;
            color: #495057;
            font-size: 0.9rem;
        }
        .footer {
            text-align: center;
            padding: 30px;
            background: #f8f9fa;
            color: #6c757d;
            font-size: 0.9rem;
        }
        .footer a {
            color: #667eea;
            text-decoration: none;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        @media print {
            body {
                background: white;
                padding: 0;
            }
            .container {
                box-shadow: none;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 OpenAPI Test Results</h1>
            <p>Generated on {{.Timestamp}}</p>
        </div>
        
        <div class="stats">
            <div class="stat-box">
                <div class="stat-label">Total Tests</div>
                <div class="stat-value">{{.TotalTests}}</div>
            </div>
            <div class="stat-box">
                <div class="stat-label">Passed</div>
                <div class="stat-value pass">{{.Passed}}</div>
            </div>
            <div class="stat-box">
                <div class="stat-label">Failed</div>
                <div class="stat-value fail">{{.Failed}}</div>
            </div>
            <div class="stat-box">
                <div class="stat-label">Pass Rate</div>
                <div class="stat-value {{if ge .PassRate 100.0}}rate-good{{else if ge .PassRate 50.0}}rate-ok{{else}}rate-bad{{end}}">
                    {{printf "%.1f" .PassRate}}%
                </div>
            </div>
            {{if .TotalTime}}
            <div class="stat-box">
                <div class="stat-label">Total Time</div>
                <div class="stat-value">{{.TotalTime}}</div>
            </div>
            <div class="stat-box">
                <div class="stat-label">Average Time</div>
                <div class="stat-value">{{.AverageTime}}</div>
            </div>
            {{end}}
        </div>
        
        <div class="meta">
            {{if .SpecPath}}
            <div class="meta-row">
                <span class="meta-label">Spec File:</span>
                <span class="meta-value">{{.SpecPath}}</span>
            </div>
            {{end}}
            {{if .BaseURL}}
            <div class="meta-row">
                <span class="meta-label">Base URL:</span>
                <span class="meta-value">{{.BaseURL}}</span>
            </div>
            {{end}}
        </div>
        
        <div class="results">
            <h2>📊 Test Results</h2>
            <table class="results-table">
                <thead>
                    <tr>
                        <th>Method</th>
                        <th>Endpoint</th>
                        <th>Status</th>
                        <th>Message</th>
                        <th>Duration</th>
                        <th>Retries</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Results}}
                    <tr class="{{.RowClass}}">
                        <td><span class="method method-{{.Method | lower}}">{{.Method}}</span></td>
                        <td><span class="endpoint">{{.Endpoint}}</span></td>
                        <td><span class="status {{if eq .RowClass "success"}}status-success{{else}}status-error{{end}}">{{.Status}}</span></td>
                        <td><span class="message">{{.Message}}</span></td>
                        <td><span class="duration">{{.Duration}}</span></td>
                        <td><span class="retry-count">{{.RetryCount}}</span></td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        
        <div class="footer">
            <p>Generated by <a href="https://github.com/Traves-Theberge/OpenAPI-Toolkit" target="_blank">OpenAPI TUI</a></p>
            <p>OpenAPI Testing Toolkit • Made with ❤️ for developers</p>
        </div>
    </div>
</body>
</html>`

// ExportResultsToHTML exports test results to a formatted HTML file
// Returns the filename and any error
func ExportResultsToHTML(results []models.TestResult, specPath, baseURL string) (string, error) {
	// Calculate statistics
	passed := 0
	failed := 0
	var totalDuration time.Duration

	for _, r := range results {
		totalDuration += r.Duration
		// Consider 2xx status codes as passed
		if strings.HasPrefix(r.Status, "2") {
			passed++
		} else {
			failed++
		}
	}

	// Calculate pass rate
	passRate := 0.0
	if len(results) > 0 {
		passRate = (float64(passed) / float64(len(results))) * 100
	}

	// Format timing statistics
	totalTime := formatDuration(totalDuration)
	averageTime := ""
	if len(results) > 0 {
		averageTime = formatDuration(totalDuration / time.Duration(len(results)))
	}

	// Convert results to HTML format
	htmlResults := make([]HTMLResult, len(results))
	for i, r := range results {
		rowClass := "success"
		if !strings.HasPrefix(r.Status, "2") {
			rowClass = "failure"
		}

		htmlResults[i] = HTMLResult{
			Method:     r.Method,
			Endpoint:   r.Endpoint,
			Status:     r.Status,
			Message:    r.Message,
			Duration:   formatDuration(r.Duration),
			RetryCount: r.RetryCount,
			RowClass:   rowClass,
		}

		// Add verbose log data if available
		if r.LogEntry != nil {
			htmlResults[i].HasLog = true
			htmlResults[i].RequestURL = r.LogEntry.RequestURL
			htmlResults[i].RequestBody = r.LogEntry.RequestBody
			htmlResults[i].ResponseBody = r.LogEntry.ResponseBody
			htmlResults[i].Timestamp = r.LogEntry.Timestamp.Format("15:04:05")
		}
	}

	// Prepare template data
	data := HTMLTemplateData{
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		SpecPath:    specPath,
		BaseURL:     baseURL,
		TotalTests:  len(results),
		Passed:      passed,
		Failed:      failed,
		PassRate:    passRate,
		Results:     htmlResults,
		HasVerbose:  false, // Can be enhanced later
		TotalTime:   totalTime,
		AverageTime: averageTime,
	}

	// Parse and execute template
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("openapi-test-results_%s.html", timestamp)

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return filename, nil
}

// formatDuration converts a duration to a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return "0µs"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

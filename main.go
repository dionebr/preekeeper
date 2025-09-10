package main

import (
        "bufio"
        "bytes"
        "crypto/tls"
        "fmt"
        "log"
        "os"
        "regexp"
        "strconv"
        "strings"
        "sync"
        "time"

        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/lipgloss"
        "github.com/spf13/cobra"
        "github.com/valyala/fasthttp"
)

// Config estrutura
type Config struct {
        URL         string
        Wordlist    string
        Threads     int
        Method      string
        StatusCodes string
        Extensions  string
        Headers     []string
        Delay       int
        Retries     int
        Timeout     int
        Recursion   bool
        MaxDepth    int
        FilterSize  string
        FilterLines string
        FilterRegex string
        NoTLS       bool
        UserAgent   string
        Cookies     string
        Proxy       string
        RateLimit   int
        Silent      bool
        Verbose     bool
        OutputFile  string
}

// Result estrutura
type Result struct {
        Path   string
        Status int
        Size   int
        Lines  int
}

// Stats estrutura
type Stats struct {
        ProcessedCount  int
        FoundCount      int
        RecursionCount  int
        RecursionActive bool
        CurrentPath     string
        RPS             float64
        Elapsed         string
}

// Job estrutura
type Job struct {
        URL   string
        Depth int
}

type scanState int

const (
        stateReady scanState = iota
        stateScanning
        stateCompleted
        statePaused
)

// Model principal do Bubble Tea
type Model struct {
        config         *Config
        state          scanState
        results        []Result
        stats          Stats
        terminalWidth  int
        terminalHeight int
        startTime      time.Time
        
        // Scanner internals
        jobs           chan Job
        mu             sync.Mutex
        progressMu     sync.Mutex
        stopChannel    chan bool
        workers        sync.WaitGroup
        producer       sync.WaitGroup
        
        // UI state
        scrollOffset   int
        showHelp       bool
        statusFilter   string
        
        // Wordlist
        wordlist       []string
}

// Estilos com paleta personalizada
var (
        // Paleta de cores personalizada
        LightGray    = lipgloss.Color("#cad4d9")  // Cinza claro
        BeigeGray    = lipgloss.Color("#b4b2a7")  // Cinza-bege  
        GrayBrown    = lipgloss.Color("#94866d")  // Marrom acinzentado
        WoodBrown    = lipgloss.Color("#746142")  // Marrom madeira
        DarkBrown    = lipgloss.Color("#42432e")  // Marrom escuro
        BlackBrown   = lipgloss.Color("#1d1f10")  // Preto
        
        // Status codes com nova paleta
        StatusOK        = lipgloss.NewStyle().Foreground(GrayBrown).Bold(true)     // Success 2xx
        StatusRedirect  = lipgloss.NewStyle().Foreground(WoodBrown).Bold(true)     // Redirect 3xx
        StatusClient    = lipgloss.NewStyle().Foreground(DarkBrown).Bold(true)     // Client Error 4xx
        StatusServer    = lipgloss.NewStyle().Foreground(BlackBrown).Bold(true)    // Server Error 5xx
        StatusNeutral   = lipgloss.NewStyle().Foreground(BeigeGray)                // Outros
        
        // UI Elements
        HeaderStyle     = lipgloss.NewStyle().Foreground(DarkBrown).Bold(true)
        BorderStyle     = lipgloss.NewStyle().Foreground(GrayBrown)
        InfoStyle       = lipgloss.NewStyle().Foreground(BeigeGray)
        SuccessStyle    = lipgloss.NewStyle().Foreground(GrayBrown).Bold(true)
        ErrorStyle      = lipgloss.NewStyle().Foreground(BlackBrown).Bold(true)
        ProgressStyle   = lipgloss.NewStyle().Foreground(WoodBrown).Bold(true)
        
        // Banner atualizado
        BannerStyle     = lipgloss.NewStyle().
                                                Foreground(DarkBrown).
                                                Bold(true).
                                                Border(lipgloss.RoundedBorder()).
                                                BorderForeground(WoodBrown).
                                                Padding(0, 1)
)

// Messages
type tickMsg time.Time
type resultMsg Result
type statsMsg Stats
type scanCompleteMsg struct{}

func NewFastHTTPClient(cfg *Config) *fasthttp.Client {
        return &fasthttp.Client{
                ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
                WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
                TLSConfig: &tls.Config{
                        InsecureSkipVerify: cfg.NoTLS,
                },
        }
}

func GetStatusColor(status int) lipgloss.Style {
        switch {
        case status >= 200 && status < 300:
                return StatusOK
        case status >= 300 && status < 400:
                return StatusRedirect
        case status >= 400 && status < 500:
                return StatusClient
        case status >= 500:
                return StatusServer
        default:
                return StatusNeutral
        }
}

func NewModel(cfg *Config) *Model {
        return &Model{
                config:      cfg,
                state:       stateReady,
                results:     []Result{},
                stopChannel: make(chan bool),
        }
}

func (m *Model) Init() tea.Cmd {
        return tea.Batch(
                tea.SetWindowTitle("Preekeeper Scanner 🐝"),
                tickCmd(),
        )
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                m.terminalWidth = msg.Width
                m.terminalHeight = msg.Height
                
        case tea.KeyMsg:
                return m.handleKeyPress(msg)
                
        case tickMsg:
                if m.state == stateScanning {
                        return m, tickCmd()
                }
                
        case resultMsg:
                m.mu.Lock()
                m.results = append(m.results, Result(msg))
                m.stats.FoundCount = len(m.results)
                m.mu.Unlock()
                
        case statsMsg:
                m.progressMu.Lock()
                m.stats = Stats(msg)
                m.progressMu.Unlock()
                
        case scanCompleteMsg:
                m.state = stateCompleted
                return m, nil
        }
        
        return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "q", "ctrl+c":
                if m.state == stateScanning {
                        close(m.stopChannel)
                        m.state = statePaused
                }
                return m, tea.Quit
                
        case "s":
                if m.state == stateReady {
                        return m, m.startScan()
                }
                
        case "p":
                if m.state == stateScanning {
                        m.state = statePaused
                        close(m.stopChannel)
                } else if m.state == statePaused {
                        m.state = stateScanning
                        m.stopChannel = make(chan bool)
                        return m, m.resumeScan()
                }
                
        case "r":
                if m.state == stateCompleted || m.state == statePaused {
                        m.resetScan()
                        return m, m.startScan()
                }
                
        case "h":
                m.showHelp = !m.showHelp
                
        case "up", "k":
                if m.scrollOffset > 0 {
                        m.scrollOffset--
                }
                
        case "down", "j":
                maxScroll := len(m.results) - (m.terminalHeight - 15)
                if maxScroll > 0 && m.scrollOffset < maxScroll {
                        m.scrollOffset++
                }
                
        case "1", "2", "3", "4", "5":
                statusMap := map[string]string{
                        "1": "2", "2": "3", "3": "4", "4": "5", "5": "",
                }
                m.statusFilter = statusMap[msg.String()]
        }
        
        return m, nil
}

func (m *Model) startScan() tea.Cmd {
        return tea.Sequence(
                func() tea.Msg {
                        m.state = stateScanning
                        m.startTime = time.Now()
                        m.loadWordlist()
                        m.initializeScanner()
                        go m.runScanner()
                        return tickMsg(time.Now())
                },
        )
}

func (m *Model) resumeScan() tea.Cmd {
        return tea.Sequence(
                func() tea.Msg {
                        go m.runScanner()
                        return tickMsg(time.Now())
                },
        )
}

func (m *Model) resetScan() {
        m.state = stateReady
        m.results = []Result{}
        m.stats = Stats{}
        m.scrollOffset = 0
        m.stopChannel = make(chan bool)
}

func (m *Model) loadWordlist() error {
        file, err := os.Open(m.config.Wordlist)
        if err != nil {
                return err
        }
        defer file.Close()
        
        m.wordlist = []string{}
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                m.wordlist = append(m.wordlist, scanner.Text())
        }
        return nil
}

func (m *Model) initializeScanner() {
        m.jobs = make(chan Job, m.config.Threads)
        m.stats = Stats{
                ProcessedCount: 0,
                FoundCount:     0,
                RecursionCount: 0,
        }
}

func (m *Model) runScanner() {
        // Parse status codes
        statusCodes := make(map[int]bool)
        for _, codeStr := range strings.Split(m.config.StatusCodes, ",") {
                code, _ := strconv.Atoi(codeStr)
                statusCodes[code] = true
        }
        
        // Parse filters
        var filterSize map[int]bool
        if m.config.FilterSize != "" {
                filterSize = make(map[int]bool)
                for _, s := range strings.Split(m.config.FilterSize, ",") {
                        size, _ := strconv.Atoi(s)
                        filterSize[size] = true
                }
        }
        
        var filterLines map[int]bool
        if m.config.FilterLines != "" {
                filterLines = make(map[int]bool)
                for _, l := range strings.Split(m.config.FilterLines, ",") {
                        lines, _ := strconv.Atoi(l)
                        filterLines[lines] = true
                }
        }
        
        var filterRegex *regexp.Regexp
        if m.config.FilterRegex != "" {
                filterRegex, _ = regexp.Compile(m.config.FilterRegex)
        }
        
        // Start job producer
        m.producer.Add(1)
        go m.produceJobs()
        
        // Start workers
        m.workers.Add(m.config.Threads)
        for i := 0; i < m.config.Threads; i++ {
                go m.worker(statusCodes, filterSize, filterLines, filterRegex)
        }
        
        // Close jobs channel when producer is done
        go func() {
                m.producer.Wait()
                close(m.jobs)
        }()
        
        // Wait for workers to finish
        m.workers.Wait()
}

func (m *Model) produceJobs() {
        defer m.producer.Done()
        
        var extensions []string
        if m.config.Extensions != "" {
                extensions = strings.Split(m.config.Extensions, ",")
        }
        
        for _, word := range m.wordlist {
                select {
                case <-m.stopChannel:
                        return
                default:
                }
                
                m.jobs <- Job{URL: word, Depth: 0}
                
                for _, ext := range extensions {
                        select {
                        case <-m.stopChannel:
                                return
                        default:
                        }
                        m.jobs <- Job{URL: word + ext, Depth: 0}
                }
        }
}

func (m *Model) worker(statusCodes map[int]bool, filterSize, filterLines map[int]bool, filterRegex *regexp.Regexp) {
        defer m.workers.Done()
        
        client := NewFastHTTPClient(m.config)
        req := fasthttp.AcquireRequest()
        resp := fasthttp.AcquireResponse()
        defer fasthttp.ReleaseRequest(req)
        defer fasthttp.ReleaseResponse(resp)
        
        req.Header.SetMethod(m.config.Method)
        req.Header.Set("User-Agent", "Preekeeper/1.0 🐝")
        
        for _, h := range m.config.Headers {
                parts := strings.SplitN(h, ":", 2)
                if len(parts) == 2 {
                        req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
                }
        }
        
        for job := range m.jobs {
                select {
                case <-m.stopChannel:
                        return
                default:
                }
                
                m.progressMu.Lock()
                m.stats.ProcessedCount++
                elapsed := time.Since(m.startTime).Seconds()
                if elapsed > 0 {
                        m.stats.RPS = float64(m.stats.ProcessedCount) / elapsed
                }
                m.stats.Elapsed = fmt.Sprintf("%02d:%02d:%02d",
                        int(time.Since(m.startTime).Hours()),
                        int(time.Since(m.startTime).Minutes())%60,
                        int(time.Since(m.startTime).Seconds())%60,
                )
                m.progressMu.Unlock()
                
                if m.config.Delay > 0 {
                        time.Sleep(time.Duration(m.config.Delay) * time.Millisecond)
                }
                
                var url string
                if strings.Contains(job.URL, "://") {
                        url = job.URL
                } else if strings.Contains(m.config.URL, "FUZZ") {
                        url = strings.Replace(m.config.URL, "FUZZ", job.URL, 1)
                } else {
                        url = fmt.Sprintf("%s/%s", strings.TrimRight(m.config.URL, "/"), job.URL)
                }
                
                req.SetRequestURI(url)
                
                m.progressMu.Lock()
                m.stats.CurrentPath = url
                m.progressMu.Unlock()
                
                var err error
                for i := 0; i <= m.config.Retries; i++ {
                        err = client.Do(req, resp)
                        if err == nil {
                                break
                        }
                        time.Sleep(50 * time.Millisecond)
                }
                
                if err == nil {
                        body := resp.Body()
                        bodySize := len(body)
                        lineCount := bytes.Count(body, []byte("\n"))
                        if bodySize > 0 {
                                lineCount++
                        }
                        
                        // Apply filters
                        if !((filterSize != nil && filterSize[bodySize]) ||
                                (filterLines != nil && filterLines[lineCount]) ||
                                (filterRegex != nil && filterRegex.Match(body))) {
                                
                                statusCode := resp.StatusCode()
                                if _, ok := statusCodes[statusCode]; ok {
                                        result := Result{
                                                Path:   url,
                                                Status: statusCode,
                                                Size:   bodySize,
                                                Lines:  lineCount,
                                        }
                                        
                                        m.mu.Lock()
                                        m.results = append(m.results, result)
                                        m.stats.FoundCount = len(m.results)
                                        m.mu.Unlock()
                                }
                        }
                }
        }
}

func (m *Model) View() string {
        if m.terminalWidth == 0 {
                return "Initializing..."
        }
        
        var b strings.Builder
        
        // Banner
        banner := BannerStyle.Render("Preekeeper Scanner 🐝 | Created by Dione Lima - Brazil")
        b.WriteString(HeaderStyle.Render(banner))
        b.WriteString("\n\n")
        
        // Help section
        if m.showHelp {
                return m.renderHelp()
        }
        
        // Configuration info
        b.WriteString(m.renderConfig())
        b.WriteString("\n")
        
        // Status and progress
        b.WriteString(m.renderProgress())
        b.WriteString("\n")
        
        // Results
        b.WriteString(m.renderResults())
        
        // Controls
        b.WriteString(m.renderControls())
        
        return b.String()
}

func (m *Model) renderConfig() string {
        width := m.terminalWidth - 4
        if width < 60 {
                width = 60
        }
        
        border := BorderStyle.Render(strings.Repeat("─", width))
        
        var b strings.Builder
        b.WriteString(border + "\n")
        
        configs := [][]string{
                {"Target", m.config.URL},
                {"Wordlist", m.config.Wordlist},
                {"Threads", fmt.Sprintf("%d", m.config.Threads)},
                {"Method", m.config.Method},
                {"Status Codes", m.config.StatusCodes},
        }
        
        if m.config.Recursion {
                configs = append(configs, []string{"Max Depth", fmt.Sprintf("%d", m.config.MaxDepth)})
        }
        
        for _, config := range configs {
                line := fmt.Sprintf("│ %-12s : %-*s │", config[0], width-20, config[1])
                b.WriteString(InfoStyle.Render(line) + "\n")
        }
        
        b.WriteString(border)
        return b.String()
}

func (m *Model) renderProgress() string {
        var b strings.Builder
        
        status := ""
        switch m.state {
        case stateReady:
                status = "Ready to scan"
        case stateScanning:
                status = "Scanning in progress..."
        case stateCompleted:
                status = "Scan completed"
        case statePaused:
                status = "Scan paused"
        }
        
        elapsed := m.stats.Elapsed
        if elapsed == "" {
                elapsed = "00:00:00"
        }
        
        statusLine := fmt.Sprintf("[%s] Elapsed: %s | Found: %d | RPS: %.2f | Processed: %d",
                status, elapsed, m.stats.FoundCount, m.stats.RPS, m.stats.ProcessedCount)
        
        b.WriteString(ProgressStyle.Render(statusLine) + "\n")
        
        if m.stats.CurrentPath != "" {
                currentLine := fmt.Sprintf("[>] Current: %s", m.stats.CurrentPath)
                b.WriteString(InfoStyle.Render(currentLine) + "\n")
        }
        
        if m.stats.RecursionActive {
                recursionLine := fmt.Sprintf("[↺] Recursion: %d additional directories", m.stats.RecursionCount)
                b.WriteString(InfoStyle.Render(recursionLine) + "\n")
        }
        
        return b.String()
}

func (m *Model) renderResults() string {
        if len(m.results) == 0 {
                return InfoStyle.Render("No results yet...\n")
        }
        
        var b strings.Builder
        b.WriteString(HeaderStyle.Render("Results:") + "\n")
        
        maxResults := m.terminalHeight - 15
        if maxResults < 5 {
                maxResults = 5
        }
        
        filteredResults := m.filterResults()
        start := m.scrollOffset
        end := start + maxResults
        
        if end > len(filteredResults) {
                end = len(filteredResults)
        }
        
        for i := start; i < end; i++ {
                result := filteredResults[i]
                statusColor := GetStatusColor(result.Status)
                
                line := fmt.Sprintf("  [%d] %s (Size: %d, Lines: %d)",
                        result.Status, result.Path, result.Size, result.Lines)
                
                b.WriteString(statusColor.Render(line) + "\n")
        }
        
        if len(filteredResults) > maxResults {
                scrollInfo := fmt.Sprintf("Showing %d-%d of %d results (↑↓ to scroll)",
                        start+1, end, len(filteredResults))
                b.WriteString(InfoStyle.Render(scrollInfo) + "\n")
        }
        
        return b.String()
}

func (m *Model) filterResults() []Result {
        if m.statusFilter == "" {
                return m.results
        }
        
        var filtered []Result
        for _, result := range m.results {
                statusPrefix := fmt.Sprintf("%d", result.Status)[0:1]
                if statusPrefix == m.statusFilter {
                        filtered = append(filtered, result)
                }
        }
        return filtered
}

func (m *Model) renderControls() string {
        controls := []string{
                "s: Start scan",
                "p: Pause/Resume",
                "r: Restart",
                "h: Help",
                "q: Quit",
        }
        
        if len(m.results) > 0 {
                controls = append(controls, "↑↓: Scroll", "1-5: Filter by status")
        }
        
        return InfoStyle.Render("\nControls: " + strings.Join(controls, " | "))
}

func (m *Model) renderHelp() string {
        help := `
Preekeeper Scanner Help

CONTROLS:
  s          - Start the scan
  p          - Pause/Resume the scan  
  r          - Restart the scan
  h          - Toggle this help
  q/Ctrl+C   - Quit the application
  ↑/k        - Scroll up in results
  ↓/j        - Scroll down in results

FILTERS:
  1          - Show only 2xx responses
  2          - Show only 3xx responses  
  3          - Show only 4xx responses
  4          - Show only 5xx responses
  5          - Show all responses

STATUS CODES:
  Green      - 2xx Success
  Yellow     - 3xx Redirection
  Purple     - 4xx Client Error
  Red        - 5xx Server Error

Press 'h' again to return to main view.
`
        return HeaderStyle.Render(help)
}

func tickCmd() tea.Cmd {
        return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
                return tickMsg(t)
        })
}

// Variáveis globais para flags
var (
        url          string
        wordlist     string
        threads      int
        method       string
        statusCodes  string
        extensions   string
        headers      []string
        delay        int
        retries      int
        timeout      int
        recursion    bool
        maxDepth     int
        filterSize   string
        filterLines  string
        filterRegex  string
        noTLS        bool
        silent       bool
        verbose      bool
        outputFile   string
        userAgent    string
        cookies      string
        proxy        string
        rateLimit    int
)

var rootCmd = &cobra.Command{
        Use:   "preekeeper",
        Short: "Preekeeper - Advanced Web Directory Scanner",
        Long: `Preekeeper Scanner 🐝
Advanced web directory brute-force tool with Bubble Tea TUI interface.
Created by Dione Lima - Brazil

A fast and feature-rich directory scanner similar to gobuster and dirb,
with a beautiful terminal user interface powered by Bubble Tea.`,
        Example: `  preekeeper -u http://example.com -w wordlist.txt
  preekeeper -u http://example.com -w wordlist.txt -t 50 -x .php,.html
  preekeeper -u http://example.com -w wordlist.txt -r -d 3
  preekeeper -u http://example.com/FUZZ -w wordlist.txt --mc 200,302`,
        Run: runScanner,
}

func init() {
        // URL flags
        rootCmd.Flags().StringVarP(&url, "url", "u", "", "Target URL (required)")
        rootCmd.MarkFlagRequired("url")
        
        // Wordlist flags
        rootCmd.Flags().StringVarP(&wordlist, "wordlist", "w", "wordlist.txt", "Wordlist file path")
        
        // Performance flags
        rootCmd.Flags().IntVarP(&threads, "threads", "t", 20, "Number of concurrent threads")
        rootCmd.Flags().IntVar(&delay, "delay", 0, "Delay between requests in milliseconds")
        rootCmd.Flags().IntVar(&timeout, "timeout", 10, "Request timeout in seconds")
        rootCmd.Flags().IntVar(&retries, "retries", 3, "Number of retries on request failure")
        rootCmd.Flags().IntVar(&rateLimit, "rate-limit", 0, "Rate limit requests per second (0 = unlimited)")
        
        // HTTP flags
        rootCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method")
        rootCmd.Flags().StringVarP(&userAgent, "user-agent", "a", "Preekeeper/1.0 🐝", "User agent string")
        rootCmd.Flags().StringSliceVarP(&headers, "headers", "H", []string{}, "Custom headers (can be used multiple times)")
        rootCmd.Flags().StringVar(&cookies, "cookies", "", "Cookies for requests")
        rootCmd.Flags().StringVar(&proxy, "proxy", "", "Proxy URL (http://host:port)")
        
        // Status and filtering flags
        rootCmd.Flags().StringVar(&statusCodes, "mc", "200,204,301,302,307,403,401,500", "Match status codes")
        rootCmd.Flags().StringVar(&filterSize, "fs", "", "Filter by response size (comma separated)")
        rootCmd.Flags().StringVar(&filterLines, "fl", "", "Filter by response lines (comma separated)")
        rootCmd.Flags().StringVar(&filterRegex, "fr", "", "Filter responses by regex pattern")
        
        // Extension and recursion flags
        rootCmd.Flags().StringVarP(&extensions, "extensions", "x", "", "File extensions (comma separated)")
        rootCmd.Flags().BoolVarP(&recursion, "recursive", "r", false, "Enable recursive scanning")
        rootCmd.Flags().IntVarP(&maxDepth, "depth", "d", 2, "Maximum recursion depth")
        
        // Security flags
        rootCmd.Flags().BoolVar(&noTLS, "no-tls-validation", false, "Skip TLS certificate validation")
        
        // Output flags
        rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silent mode (no banner)")
        rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
        rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for results")
}

func runScanner(cmd *cobra.Command, args []string) {
        // Validar URL
        if url == "" {
                fmt.Println(ErrorStyle.Render("Error: URL is required. Use -u flag."))
                os.Exit(1)
        }
        
        // Validar wordlist
        if _, err := os.Stat(wordlist); os.IsNotExist(err) {
                fmt.Println(ErrorStyle.Render(fmt.Sprintf("Error: Wordlist file '%s' not found", wordlist)))
                os.Exit(1)
        }
        
        // Criar configuração
        cfg := &Config{
                URL:         url,
                Wordlist:    wordlist,
                Threads:     threads,
                Method:      strings.ToUpper(method),
                StatusCodes: statusCodes,
                Extensions:  extensions,
                Headers:     headers,
                Delay:       delay,
                Retries:     retries,
                Timeout:     timeout,
                Recursion:   recursion,
                MaxDepth:    maxDepth,
                FilterSize:  filterSize,
                FilterLines: filterLines,
                FilterRegex: filterRegex,
                NoTLS:       noTLS,
                UserAgent:   userAgent,
                Cookies:     cookies,
                Proxy:       proxy,
                RateLimit:   rateLimit,
                Silent:      silent,
                Verbose:     verbose,
                OutputFile:  outputFile,
        }
        
        // Validações adicionais
        if cfg.Threads > 100 {
                fmt.Println(ErrorStyle.Render("Warning: High thread count (>100) may cause issues"))
        }
        
        if cfg.Delay < 0 {
                cfg.Delay = 0
        }
        
        // Criar modelo e iniciar TUI
        model := NewModel(cfg)
        
        // Configurar programa
        var opts []tea.ProgramOption
        if !silent {
                opts = append(opts, tea.WithAltScreen())
        }
        
        p := tea.NewProgram(model, opts...)
        
        if _, err := p.Run(); err != nil {
                log.Fatalf("Error running scanner: %v", err)
        }
}

func main() {
        if err := rootCmd.Execute(); err != nil {
                fmt.Println(ErrorStyle.Render(err.Error()))
                os.Exit(1)
        }
}
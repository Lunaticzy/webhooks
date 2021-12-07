package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	file, err := os.OpenFile("error.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(ioutil.Discard,
		"TRACE:",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO:",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING:",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR:",
		log.Ldate|log.Ltime|log.Lshortfile)
}

const PullRequestHookName = "merge_request_hooks"

type HookStruct struct {
	HookName    string      `json:"hook_name"`
	Password    string      `json:"password"`
	HookID      int         `json:"hook_id"`
	HookURL     string      `json:"hook_url"`
	Timestamp   string      `json:"timestamp"`
	Sign        string      `json:"sign"`
	PullRequest PullRequest `json:"pull_request"`
	Author      Author      `json:"author"`
	Sender      Sender      `json:"sender"`
	Enterprise  Enterprise  `json:"enterprise"`
}
type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Assignee struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Milestone struct {
	HtmlUrl        string      `json:"html_url"`
	ID             int         `json:"id"`
	Number         int         `json:"number"`
	Title          string      `json:"title"`
	Description    interface{} `json:"description"`
	OpenIssues     int         `json:"open_issues"`
	StartedIssues  int         `json:"started_issues"`
	ClosedIssues   int         `json:"closed_issues"`
	ApprovedIssues int         `json:"approved_issues"`
	State          string      `json:"state"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DueOn          interface{} `json:"due_on"`
}
type Owner struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Repo struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Path              string      `json:"path"`
	FullName          string      `json:"full_name"`
	Owner             Owner       `json:"owner"`
	Private           bool        `json:"private"`
	HtmlUrl           string      `json:"html_url"`
	URL               string      `json:"url"`
	Description       string      `json:"description"`
	Fork              bool        `json:"fork"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	PushedAt          time.Time   `json:"pushed_at"`
	GitURL            string      `json:"git_url"`
	SshUrl            string      `json:"ssh_url"`
	CloneURL          string      `json:"clone_url"`
	SvnURL            string      `json:"svn_url"`
	GitHttpUrl        string      `json:"git_http_url"`
	GitSshUrl         string      `json:"git_ssh_url"`
	GitSvnURL         string      `json:"git_svn_url"`
	Homepage          interface{} `json:"homepage"`
	StargazersCount   int         `json:"stargazers_count"`
	WatchersCount     int         `json:"watchers_count"`
	ForksCount        int         `json:"forks_count"`
	Language          string      `json:"language"`
	HasIssues         bool        `json:"has_issues"`
	HasWiki           bool        `json:"has_wiki"`
	HasPages          bool        `json:"has_pages"`
	License           interface{} `json:"license"`
	OpenIssuesCount   int         `json:"open_issues_count"`
	DefaultBranch     string      `json:"default_branch"`
	Namespace         string      `json:"namespace"`
	NameWithNamespace string      `json:"name_with_namespace"`
	PathWithNamespace string      `json:"path_with_namespace"`
}
type Head struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}
type Base struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}
type PullRequest struct {
	ID             int         `json:"id"`
	Number         int         `json:"number"`
	State          string      `json:"state"`
	HtmlUrl        string      `json:"html_url"`
	DiffURL        string      `json:"diff_url"`
	PatchURL       string      `json:"patch_url"`
	Title          string      `json:"title"`
	Body           interface{} `json:"body"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	ClosedAt       interface{} `json:"closed_at"`
	MergedAt       interface{} `json:"merged_at"`
	MergeCommitSha string      `json:"merge_commit_sha"`
	User           User        `json:"user"`
	Assignee       Assignee    `json:"assignee"`
	Tester         interface{} `json:"tester"`
	Milestone      Milestone   `json:"milestone"`
	Head           Head        `json:"head"`
	Base           Base        `json:"base"`
	Merged         bool        `json:"merged"`
	Mergeable      interface{} `json:"mergeable"`
	Comments       int         `json:"comments"`
	Commits        int         `json:"commits"`
	Additions      int         `json:"additions"`
	Deletions      int         `json:"deletions"`
	ChangedFiles   int         `json:"changed_files"`
}
type Author struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Sender struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Enterprise struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponseStruct(code int, msg string, data interface{}) *Response {
	return &Response{Code: code, Msg: msg, Data: data}
}

var config map[string]string
var mutex sync.Mutex

func init() {
	file, err := os.Open("config.yml")
	if err != nil {
		Error.Println("Failed to open config.yml resource")
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Warning.Println("Close config file failed")
		}
	}(file)

	configBuff, err := ioutil.ReadAll(file)
	if err != nil {
		Error.Println("Failed to read yaml config file")
		return
	}

	err = yaml.Unmarshal(configBuff, &config)
	if err != nil {
		Error.Println("Parse yaml config failed:", err)
		return
	}
	Info.Println("Load yaml config file success")
}

func hook(rw http.ResponseWriter, r *http.Request) {

	Trace.Printf("received request, method:%v, url:%v", r.Method, r.URL)

	if r.Method == "GET" {
		Error.Printf("IGNORE GET REQUEST")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	res := NewResponseStruct(200, "请求成功", nil)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println("Read body error:", err)
	}

	var hook HookStruct
	if err = json.Unmarshal(body, &hook); err != nil {
		Error.Println("Unmarshal body error:", err)
		return
	}

	if hook.HookName == PullRequestHookName {
		if hook.PullRequest.Merged {
			requestUrl := r.URL.String()
			if strings.Contains(requestUrl, "/hooks") {
				index := strings.LastIndex(requestUrl, "/")
				s := requestUrl[(index + 1):]
				shellPath := config[s]
				Info.Printf("Receive hook request, project name is %s", hook.PullRequest.Base.Repo.FullName)

				if strings.Trim(shellPath, " ") != "" {
					go func() {
						mutex.Lock()
						startTime := time.Now()
						command := exec.Command("/bin/bash", shellPath)
						err := command.Start()
						if err != nil {
							Error.Println("Start execute shell failed", err)
						}

						Info.Println("Start execute the shell, Process Pid:", command.Process.Pid)
						err = command.Wait()
						if err != nil {
							Error.Println("Execute shell error:", err)
						}
						endTime := time.Now()
						Info.Printf("Process execute success, PID:%v, cost:%v", command.ProcessState.Pid(), endTime.Sub(startTime))
						mutex.Unlock()
					}()
				}

			} else {
				Error.Println("Hook not running, not script to run")
				res = NewResponseStruct(200, "hook运行失败, 未配置脚本", nil)
			}

		} else {
			Error.Println("Hook not running, pull request not merge")
			res = NewResponseStruct(200, "hook运行失败, pr未合并", nil)
		}
	} else {
		Error.Printf("Hook not running, not support hook type: %s", hook.HookName)
		res = NewResponseStruct(200, "hook运行失败, 不支持的hook类型", nil)
	}

	err = json.NewEncoder(rw).Encode(&res)
	if err != nil {
		Error.Println("response error:", err)
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hook)

	server := &http.Server{
		Addr:    "0.0.0.0:6666",
		Handler: mux,
	}
	Info.Println("Start server listening at ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		Error.Println("start server failed:", err)
	}
}

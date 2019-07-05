package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/webhooks.v5/github"
)

var regexIssueKey = "\\[[A-Z]*\\-[0-9]+\\]"
var jiraClient *jira.Client

func InitJiraClient() {
	tp := jira.BasicAuthTransport{
		Username: "ramadhanm1998@gmail.com",
		Password: "icB26nXqVx90BRVTrxKKB68F",
	}

	jiraClient, _ = jira.NewClient(tp.Client(), "https://m-f-hafizh.atlassian.net/")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func handlers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hook, _ := github.New(github.Options.Secret("secret"))
	payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}

	//Regex issue key
	reg, _ := regexp.Compile(regexIssueKey)

	switch payload.(type) {
	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		// Do whatever you want from here...
		enc, err := json.MarshalIndent(release, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}
		fmt.Println("Release")
		fmt.Fprintf(w, string(enc))

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		// Do whatever you want from here...
		title := pullRequest.PullRequest.Title
		issueKey := strings.Replace(strings.Replace(reg.FindString(title), "[", "", -1), "]", "", -1)
		transitions, _, err := jiraClient.Issue.GetTransitions(issueKey)
		var transID string
		if pullRequest.Action == "open" {
			for _, transition := range transitions {
				if transition.To.StatusCategory.Name == "In Progress" {
					transID = transition.ID
				}
			}
		}
		res, err := jiraClient.Issue.DoTransition(issueKey, transID)
		// issue, _, err := jiraClient.Issue.Get(issueKey, nil)
		// fmt.Println("Pull Request", issue)
		enc, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}

		fmt.Println("Pull Request")
		fmt.Fprintf(w, string(enc))

	case github.PullRequestReviewPayload:
		pullReqReview := payload.(github.PullRequestReviewPayload)
		enc, err := json.MarshalIndent(pullReqReview, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}
		fmt.Println("Pull Request Review")
		fmt.Fprintf(w, string(enc))

	case github.RepositoryPayload:
		repositoryPayload := payload.(github.RepositoryPayload)
		enc, err := json.MarshalIndent(repositoryPayload, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}
		fmt.Println("Repository Payload")
		fmt.Fprintf(w, string(enc))

	case github.PushPayload:
		pushPayload := payload.(github.PushPayload)
		enc, err := json.MarshalIndent(pushPayload, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}

		fmt.Println("Push")
		fmt.Fprintf(w, string(enc))

	case github.MergedBy:
		merge := payload.(github.MergedBy)
		enc, err := json.MarshalIndent(merge, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}

		fmt.Println("Merge")
		fmt.Fprintf(w, string(enc))

	case github.CommitCommentPayload:
		commitComment := payload.(github.CommitCommentPayload)
		enc, err := json.MarshalIndent(commitComment, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}

		fmt.Println("Commit Comment")
		fmt.Fprintf(w, string(enc))
	}

}

func main() {
	port := ":" + os.Getenv("PORT")
	router := httprouter.New()
	fmt.Println("Running ...")
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.POST("/payload", handlers)
	log.Fatal(http.ListenAndServe(port, router))
}

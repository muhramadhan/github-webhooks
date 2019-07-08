package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
<<<<<<< HEAD
=======
	"os"
>>>>>>> 9df72e7d811d1d7b3eec0889431ee3f12a8761cb
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/webhooks.v5/github"
)

<<<<<<< HEAD
var regexIssueKey = "\\[[A-Z]*\\-[0-9]+\\]"
var regexIssueKeyBranch = "[A-Z]*\\-[0-9]+"
var jiraClient *jira.Client

func InitJiraClient() {
	tp := jira.BasicAuthTransport{
		Username: "ramadhanm1998@gmail.com",
		Password: "icB26nXqVx90BRVTrxKKB68F",
	}

	jiraClient, _ = jira.NewClient(tp.Client(), "https://m-f-hafizh.atlassian.net/")
}
=======
var regexProjectKey = "\\[[A-Z]*\\-[0-9]+\\]"
>>>>>>> 9df72e7d811d1d7b3eec0889431ee3f12a8761cb

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func handlers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
<<<<<<< HEAD

	hook, _ := github.New(github.Options.Secret("secret"))
	payload, err := hook.Parse(r, github.PullRequestEvent, github.CreateEvent)
=======
	tp := jira.BasicAuthTransport{
		Username: "hafizh203@gmail.com",
		Password: "nwXanAF4FVQVToP4OjDN9808",
	}
	client, _ := jira.NewClient(tp.Client(), "https://m-f-hafizh.atlassian.net/")

	hook, _ := github.New(github.Options.Secret("secret"))
	payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent, github.CreateEvent, github.PushEvent)
>>>>>>> 9df72e7d811d1d7b3eec0889431ee3f12a8761cb
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}
<<<<<<< HEAD
	//Regex issue key
	reg, _ := regexp.Compile(regexIssueKey)
=======
>>>>>>> 9df72e7d811d1d7b3eec0889431ee3f12a8761cb

	switch payload.(type) {
	case github.CreatePayload:
		// regBranch, _ := regexp.Compile(regexIssueKeyBranch)
		release := payload.(github.ReleasePayload)
		// Do whatever you want from here...
		enc, err := json.MarshalIndent(release, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest")
			return
		}
		fmt.Println("Release")
		fmt.Fprintf(w, string(enc))

<<<<<<< HEAD
	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		// Do whatever you want from here...
		title := pullRequest.PullRequest.Title
		issueKey := strings.Replace(strings.Replace(reg.FindString(title), "[", "", -1), "]", "", -1)
		if issueKey == "" {
			return
		}
		fmt.Println("jiraClient = ", jiraClient)
		transitions, _, err := jiraClient.Issue.GetTransitions(issueKey)
		if err != nil {
			fmt.Println(err)
		}
		var transID string
		if pullRequest.Action == "open" {
			for _, transition := range transitions {
				fmt.Println(transition.To.Name)
				if transition.To.Name == "In Review" {
					transID = transition.ID
				}
			}
		} else if pullRequest.Action == "closed" {
			if pullRequest.PullRequest.Merged {
				for _, transition := range transitions {
					if transition.To.Name == "Done" {
						transID = transition.ID
					}
				}
			} else {
				for _, transition := range transitions {
					if transition.To.Name == "In Progress" {
						transID = transition.ID
					}
				}
			}
		}

		if transID == "" {
			return
		}
		res, err := jiraClient.Issue.DoTransition(issueKey, transID)
		if err != nil {
			fmt.Fprint(w, "invalidRequest pertama: ", err)
			return
		}
		enc, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Fprint(w, "invalidRequest kedua: ", err)
			return
		}
		// issue, _, err := jiraClient.Issue.Get(issueKey, nil)
		// fmt.Println("Pull Request", issue)
=======
	case github.CreatePayload:
		createPayload := payload.(github.CreatePayload)
		branchName := createPayload.Ref
		fmt.Println(branchName)
		splitedName := strings.Split(branchName, "_")
		issueKey := splitedName[len(splitedName)-1]
		issue, _, _ := client.Issue.Get(issueKey, nil)
		if createPayload.RefType == "branch" {
			transitions, _, _ := client.Issue.GetTransitions(issueKey)
			for _, transition := range transitions {
				if transition.To.Name == "In Progress" {
					client.Issue.DoTransition(issue.ID, transition.ID)
					fmt.Println("Transition")
				}
			}
		}
		fmt.Println("New Branch")
>>>>>>> 9df72e7d811d1d7b3eec0889431ee3f12a8761cb

	case github.PullRequestPayload:
		fmt.Println("Pull Request")
		pullRequest := payload.(github.PullRequestPayload)
		reg, _ := regexp.Compile(regexProjectKey)
		title := pullRequest.PullRequest.Title
		issueKey := strings.Replace(strings.Replace(reg.FindString(title), "[", "", -1), "]", "", -1)
		issue, _, _ := client.Issue.Get(issueKey, nil)
		transitions, _, _ := client.Issue.GetTransitions(issueKey)
		for _, transition := range transitions {
			if transition.To.Name == "In Review" {
				client.Issue.DoTransition(issue.ID, transition.ID)
				fmt.Println("Transition")
			}
		}

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
	// port := ":" + os.Getenv("PORT")
	InitJiraClient()
	port := ":8080"
	router := httprouter.New()
	fmt.Println("Running ...")
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.POST("/payload", handlers)
	log.Fatal(http.ListenAndServe(port, router))
}

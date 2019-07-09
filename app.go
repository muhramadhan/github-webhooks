package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/webhooks.v5/github"
)

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

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func handlers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	hook, _ := github.New(github.Options.Secret("secret"))
	payload, err := hook.Parse(r, github.PullRequestEvent, github.CreateEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}
	//Regex issue key
	reg, _ := regexp.Compile(regexIssueKey)

	switch payload.(type) {

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		// Do whatever you want from here...
		title := pullRequest.PullRequest.Title
		issueKey := strings.Replace(strings.Replace(reg.FindString(title), "[", "", -1), "]", "", -1)
		if issueKey == "" {
			return
		}

		// Getting issue
		issue, _, err := jiraClient.Issue.Get(issueKey, nil)
		if err != nil {
			fmt.Println("Error : ", err)
		}

		transitions, _, err := jiraClient.Issue.GetTransitions(issueKey)
		if err != nil {
			fmt.Println("Error : ", err)
		}
		var transID string
		if pullRequest.Action == "opened" {
			for _, transition := range transitions {
				if transition.To.Name == "In Review" {
					transID = transition.ID
				}
			}
			// Getting the comment
			currUser, _, err := jiraClient.User.GetSelf()
			if err != nil {
				fmt.Println("Error : ", err)
			}
			var comment *jira.Comment
			if len(issue.Fields.Comments.Comments) != 0 {
				for _, commentLoop := range issue.Fields.Comments.Comments {
					if currUser.AccountID == commentLoop.Author.AccountID && strings.Contains(commentLoop.Body, "Pull Request") {
						comment = commentLoop
					}
				}
			}
			var BodyComment string
			if comment == nil {
				BodyComment = fmt.Sprintf("Pull Request:\n - %s\n", pullRequest.PullRequest.HTMLURL)
				newComment := jira.Comment{
					Body: BodyComment,
				}
				jiraClient.Issue.AddComment(issueKey, &newComment)
			} else {
				BodyComment = fmt.Sprintf("- %s\n", pullRequest.PullRequest.HTMLURL)
				updateComment := jira.Comment{
					ID:   comment.ID,
					Body: comment.Body + BodyComment,
				}
				jiraClient.Issue.UpdateComment(issueKey, &updateComment)
			}
		} else if pullRequest.Action == "closed" {

			// Getting the comment
			currUser, _, err := jiraClient.User.GetSelf()
			if err != nil {
				fmt.Println("Error : ", err)
			}

			var comment *jira.Comment
			if len(issue.Fields.Comments.Comments) != 0 {
				for _, commentLoop := range issue.Fields.Comments.Comments {
					if currUser.AccountID == commentLoop.Author.AccountID && strings.Contains(commentLoop.Body, "Pull Request") {
						comment = commentLoop
					}
				}
			}

			if pullRequest.PullRequest.Merged {
				for _, transition := range transitions {
					if transition.To.Name == "Done" {
						transID = transition.ID
					}
				}
				updateComment := jira.Comment{
					ID:   comment.ID,
					Body: comment.Body + " (Merged)\n",
				}
				jiraClient.Issue.UpdateComment(issueKey, &updateComment)
			} else {
				for _, transition := range transitions {
					if transition.To.Name == "In Progress" {
						transID = transition.ID
					}
				}
				updateComment := jira.Comment{
					ID:   comment.ID,
					Body: comment.Body + " (Closed)\n",
				}
				jiraClient.Issue.UpdateComment(issueKey, &updateComment)
			}
		}

		if transID == "" {
			return
		}
		_, err = jiraClient.Issue.DoTransition(issueKey, transID)
		if err != nil {
			fmt.Fprint(w, "invalidRequest pertama: ", err)
			return
		}
		// enc, err := json.MarshalIndent(res, "", "  ")
		// if err != nil {
		// 	fmt.Fprint(w, "invalidRequest kedua: ", err)
		// 	return
		// }
		// issue, _, err := jiraClient.Issue.Get(issueKey, nil)
		// fmt.Println("Pull Request", issue)
	case github.CreatePayload:
		createPayload := payload.(github.CreatePayload)
		branchName := createPayload.Ref
		fmt.Println(branchName)
		splitedName := strings.Split(branchName, "_")
		issueKey := splitedName[len(splitedName)-1]
		issue, _, _ := jiraClient.Issue.Get(issueKey, nil)
		if createPayload.RefType == "branch" {
			transitions, _, _ := jiraClient.Issue.GetTransitions(issueKey)
			BodyComment := fmt.Sprintf("Working Branch : %s", createPayload.Repository.HTMLURL+"/tree/"+createPayload.Ref)
			for _, transition := range transitions {
				if transition.To.Name == "In Progress" {
					jiraClient.Issue.DoTransition(issue.ID, transition.ID)
					fmt.Println("Transition")
					comment := jira.Comment{
						Body: BodyComment,
					}
					jiraClient.Issue.AddComment(issue.ID, &comment)
				}
			}
		}
		fmt.Println("New Branch")

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

package main

import (
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

var regexProjectKey = "\\[[A-Z]*\\-[0-9]+\\]"

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func handlers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tp := jira.BasicAuthTransport{
		Username: "hafizh203@gmail.com",
		Password: "nwXanAF4FVQVToP4OjDN9808",
	}
	client, _ := jira.NewClient(tp.Client(), "https://m-f-hafizh.atlassian.net/")

	hook, _ := github.New(github.Options.Secret("secret"))
	payload, err := hook.Parse(r, github.PullRequestEvent, github.CreateEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
		}
	}

	switch payload.(type) {

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
					comment := jira.Comment{
						Body: "Working Branch: " + createPayload.Repository.HTMLURL + "/tree/" + branchName,
					}
					client.Issue.AddComment(issue.ID, &comment)
				}
			}
		}
		fmt.Println("New Branch")

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

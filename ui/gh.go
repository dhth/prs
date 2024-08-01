package ui

import (
	"log"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	ghgql "github.com/cli/shurcooL-graphql"
)

func getPRDataFromQuery(ghClient *ghapi.GraphQLClient, queryStr string, prCount int) ([]pr, error) {
	var query prSearchQuery

	variables := map[string]interface{}{
		"query": ghgql.String(queryStr),
		"count": ghgql.Int(prCount),
	}
	err := ghClient.Query("PRQuery", &query, variables)
	if err != nil {
		return nil, err
	}
	var prs []pr
	for _, edge := range query.Search.Edges {
		if edge.Node.Type != "PullRequest" {
			continue
		}
		prs = append(prs, edge.Node.pr)
	}
	return prs, nil
}

func getPRMetadata(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prNumber int) (prDetails, error) {
	var query prDetailsQuery

	variables := map[string]interface{}{
		"repositoryOwner":    ghgql.String(repoOwner),
		"repositoryName":     ghgql.String(repoName),
		"pullRequestNumber":  ghgql.Int(prNumber),
		"latestReviewsCount": ghgql.Int(latestReviewsCount),
		"filesCount":         ghgql.Int(filesCount),
		"labelsCount":        ghgql.Int(labelsCount),
		"assigneesCount":     ghgql.Int(assigneesCount),
		"issuesCount":        ghgql.Int(issuesCount),
		"participantsCount":  ghgql.Int(participantsCount),
		"commentsCount":      ghgql.Int(commentsCount),
		"commitsCount":       ghgql.Int(commitsCount),
	}
	err := ghClient.Query("PRTL", &query, variables)
	if err != nil {
		log.Printf("error: %s\n", err)
		return prDetails{}, err
	}
	return query.RepositoryOwner.Repository.PullRequest, nil
}

func getViewerLoginData(ghClient *ghapi.GraphQLClient) (string, error) {
	var query userLoginQuery

	err := ghClient.Query("PullRequests", &query, nil)
	if err != nil {
		return "", err
	}
	return query.Viewer.Login, nil
}

func getPRTLData(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prNumber int, tlItemsCount int) ([]prTLItem, error) {
	var query prTLQuery

	variables := map[string]interface{}{
		"repositoryOwner":    ghgql.String(repoOwner),
		"repositoryName":     ghgql.String(repoName),
		"pullRequestNumber":  ghgql.Int(prNumber),
		"timelineItemsCount": ghgql.Int(tlItemsCount),
	}
	err := ghClient.Query("PRTL", &query, variables)
	if err != nil {
		return nil, err
	}
	return query.RepositoryOwner.Repository.PullRequest.TimelineItems.Nodes, nil
}

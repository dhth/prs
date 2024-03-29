package ui

import (
	ghapi "github.com/cli/go-gh/v2/pkg/api"
	ghgql "github.com/cli/shurcooL-graphql"
)

func GetPRs(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prCount int) ([]pr, error) {
	var query prsQuery

	variables := map[string]interface{}{
		"repositoryOwner":  ghgql.String(repoOwner),
		"repositoryName":   ghgql.String(repoName),
		"pullRequestCount": ghgql.Int(prCount),
	}
	err := ghClient.Query("PullRequests", &query, variables)
	if err != nil {
		return nil, err
	}
	return query.RepositoryOwner.Repository.PullRequests.Nodes, nil
}

func GetPRTL(ghClient *ghapi.GraphQLClient, repoOwner string, repoName string, prNumber int, tlItemsCount int, commentsCount int) ([]prTLItem, error) {
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

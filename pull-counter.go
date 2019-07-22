/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 *
 */
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Fetch all the PR count of a organization.
func ListPullRequests(org string, repo string, start time.Time) int {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "ACCESS_TOKEN"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// fmt.Printf("Time is: " + start.Format(time.RFC3339) + "\n")

	var count int = 0
	var page int = 1
	var perPage int = 100

	for {
		opt := &github.PullRequestListOptions{State: "all", Sort: "created", Direction: "desc", ListOptions: github.ListOptions{Page: page, PerPage: perPage}}
		prs, resp, err := client.PullRequests.List(context.Background(), org, repo, opt)
		if err != nil {
			fmt.Println(err)
			return 0
		}

		// for i, pr := range prs {
		// 	t := pr.CreatedAt
		// 	fmt.Printf(" %v. %v\n", i+1, t.Format(time.RFC3339))
		// }

		var countInCurrentPage int = 0

		for _, pr := range prs {
			check := pr.CreatedAt
			if check.Before(start) {
				break
			} else {
				countInCurrentPage++
				// fmt.Printf(" %v. %v : %v\n", countInCurrentPage, i+1, check.Format(time.RFC3339))
			}
		}

		count += countInCurrentPage

		if len(prs) > countInCurrentPage {
			break
		} else if resp.NextPage == 0 {
			break
		} else {
			page = resp.NextPage
			// fmt.Printf("Should pull page %d\n", page)
		}

		// fmt.Printf("PR list: %s\n", prs)
	}

	// fmt.Printf("PR count: %d\n", count)
	fmt.Printf("..")
	return count
}

func PrintPRStatsForRepo(org string, repo string) (int, int, int) {
	q3, _ := time.Parse(time.RFC3339, "2019-06-01T00:00:00Z")
	cq3 := ListPullRequests(org, repo, q3)

	q2, _ := time.Parse(time.RFC3339, "2019-04-01T00:00:00Z")
	cq2 := ListPullRequests(org, repo, q2)

	q1, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	cq1 := ListPullRequests(org, repo, q1)

	fmt.Printf("PR count for %s/%s: Q1=%d, Q2=%d, Q2+=%d\n", org, repo, cq1-cq2, cq2-cq3, cq3)

	return cq1 - cq2, cq2 - cq3, cq3
}

func PrintPRStatsForOrg(org string) (int, int, int) {
	var cq1 int = 0
	var cq2 int = 0
	var cq3 int = 0

	file, err := os.Open("/tmp/" + org + "-repo-list.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var repo string = scanner.Text()
		// fmt.Println(repo)
		q1, q2, q3 := PrintPRStatsForRepo(org, repo)

		cq1 += q1
		cq2 += q2
		cq3 += q3
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PR count for %s: Q1=%d, Q2=%d, Q2+=%d\n", org, cq1, cq2, cq3)
	return cq1, cq2, cq3
}

func main() {
	var cq1 int = 0
	var cq2 int = 0
	var cq3 int = 0

	var org string = "wso2"
	q1, q2, q3 := PrintPRStatsForOrg(org)

	cq1 += q1
	cq2 += q2
	cq3 += q3

	org = "wso2-extensions"
	eq1, eq2, eq3 := PrintPRStatsForOrg(org)

	cq1 += eq1
	cq2 += eq2
	cq3 += eq3

	fmt.Printf("Total PR count: Q1=%d, Q2=%d, Q2+=%d\n", cq1, cq2, cq3)
}

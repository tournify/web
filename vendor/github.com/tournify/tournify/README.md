# Tournify

[![GoDoc](https://godoc.org/github.com/tournify/tournify?status.svg)](https://godoc.org/github.com/tournify/tournify)
[![Go Report Card](https://goreportcard.com/badge/github.com/tournify/tournify)](https://goreportcard.com/report/github.com/tournify/tournify)
[![Build Status](https://api.travis-ci.org/tournify/tournify.svg?branch=master)](https://travis-ci.org/tournify/tournify)
[![Build status](https://ci.appveyor.com/api/projects/status/9s2ykpx3wdnf9eiw?svg=true)](https://ci.appveyor.com/project/markustenghamn/tournify)
[![CircleCI](https://circleci.com/gh/tournify/tournify.svg?style=svg)](https://circleci.com/gh/tournify/tournify)
[![codecov](https://codecov.io/gh/tournify/tournify/branch/master/graph/badge.svg)](https://codecov.io/gh/tournify/tournify)

This project aims to support the creation of any tournament.

Current features
 - Group tournament creation
 - Group tournament stats
 
Planned features
 - Elimination tournaments
 - Double elimination
 - Round robin

Example
=

To create a group tournament with 2 groups where all teams in each group meet one time simply do the following.

```go
package main

import (
	"fmt"
	"github.com/tournify/tournify"
	"math"
	"math/rand"
)

func main()  {
	teams := []tournify.Team{
		{ID:0},
		{ID:1},
		{ID:2},
		{ID:3},
		{ID:4},
		{ID:5},
		{ID:6},
		{ID:7},
		{ID:8},
		{ID:9},
		{ID:10},
		{ID:11},
		{ID:12},
		{ID:13},
		{ID:14},
		{ID:15},
		{ID:16},
	}

	teamInterfaces := make([]tournify.TeamInterface, len(teams))

	for i := range teams {
		teamInterfaces[i] = &teams[i]
	}

	// The CreateGroupTournamentFromTeams method takes a slice of teams along with the group count and meet count
	tournament := tournify.CreateGroupTournamentFromTeams(teamInterfaces, 2, 1)

	for _, game := range tournament.GetGames() {
		game.SetScore(randomEvenFloat(0.0, 10.0), randomEvenFloat(0.0, 10.0))
		err := tournament.SetGame(game)
		if err != nil {
			panic(err)
		}
	}

	// The print method gives us a string representing the current tournament
	fmt.Println(tournament.Print())
}

func randomEvenFloat(min float64, max float64) float64 {
	return math.RoundToEven(min + rand.Float64() * (max - min))
}
```

This will print something similar to the following output.

```text
TournamentType: Group

Groups
Group ID: 0
Team ID: 0
Team ID: 1
Team ID: 2
Team ID: 3
Team ID: 4
Team ID: 5
Team ID: 6
Team ID: 7

Group ID: 1
Team ID: 8
Team ID: 9
Team ID: 10
Team ID: 11
Team ID: 12
Team ID: 13
Team ID: 14
Team ID: 15
Team ID: 16


Games
Game ID: 0, HomeTeam: 0, AwayTeam: 7, HomeScore: 6.00, AwayScore: 9.00
Game ID: 1, HomeTeam: 1, AwayTeam: 6, HomeScore: 7.00, AwayScore: 4.00
Game ID: 2, HomeTeam: 2, AwayTeam: 5, HomeScore: 4.00, AwayScore: 7.00
Game ID: 3, HomeTeam: 3, AwayTeam: 4, HomeScore: 1.00, AwayScore: 2.00
Game ID: 4, HomeTeam: 0, AwayTeam: 6, HomeScore: 1.00, AwayScore: 3.00
Game ID: 5, HomeTeam: 2, AwayTeam: 5, HomeScore: 5.00, AwayScore: 8.00
Game ID: 6, HomeTeam: 3, AwayTeam: 4, HomeScore: 2.00, AwayScore: 4.00
Game ID: 7, HomeTeam: 7, AwayTeam: 1, HomeScore: 3.00, AwayScore: 5.00
Game ID: 8, HomeTeam: 0, AwayTeam: 5, HomeScore: 3.00, AwayScore: 3.00
Game ID: 9, HomeTeam: 3, AwayTeam: 4, HomeScore: 7.00, AwayScore: 2.00
Game ID: 10, HomeTeam: 7, AwayTeam: 1, HomeScore: 2.00, AwayScore: 4.00
Game ID: 11, HomeTeam: 6, AwayTeam: 2, HomeScore: 6.00, AwayScore: 9.00
Game ID: 12, HomeTeam: 0, AwayTeam: 4, HomeScore: 3.00, AwayScore: 3.00
Game ID: 13, HomeTeam: 7, AwayTeam: 1, HomeScore: 8.00, AwayScore: 2.00
Game ID: 14, HomeTeam: 6, AwayTeam: 2, HomeScore: 9.00, AwayScore: 7.00
Game ID: 15, HomeTeam: 5, AwayTeam: 3, HomeScore: 5.00, AwayScore: 0.00
Game ID: 16, HomeTeam: 0, AwayTeam: 1, HomeScore: 2.00, AwayScore: 6.00
Game ID: 17, HomeTeam: 6, AwayTeam: 2, HomeScore: 10.00, AwayScore: 1.00
Game ID: 18, HomeTeam: 5, AwayTeam: 3, HomeScore: 6.00, AwayScore: 1.00
Game ID: 19, HomeTeam: 4, AwayTeam: 7, HomeScore: 7.00, AwayScore: 3.00
Game ID: 20, HomeTeam: 0, AwayTeam: 2, HomeScore: 2.00, AwayScore: 5.00
Game ID: 21, HomeTeam: 5, AwayTeam: 3, HomeScore: 5.00, AwayScore: 3.00
Game ID: 22, HomeTeam: 4, AwayTeam: 7, HomeScore: 4.00, AwayScore: 5.00
Game ID: 23, HomeTeam: 1, AwayTeam: 6, HomeScore: 3.00, AwayScore: 3.00
Game ID: 24, HomeTeam: 0, AwayTeam: 3, HomeScore: 8.00, AwayScore: 4.00
Game ID: 25, HomeTeam: 4, AwayTeam: 7, HomeScore: 9.00, AwayScore: 3.00
Game ID: 26, HomeTeam: 1, AwayTeam: 6, HomeScore: 9.00, AwayScore: 1.00
Game ID: 27, HomeTeam: 2, AwayTeam: 5, HomeScore: 10.00, AwayScore: 1.00
Game ID: 29, HomeTeam: 9, AwayTeam: 16, HomeScore: 2.00, AwayScore: 7.00
Game ID: 30, HomeTeam: 10, AwayTeam: 15, HomeScore: 2.00, AwayScore: 3.00
Game ID: 31, HomeTeam: 11, AwayTeam: 14, HomeScore: 9.00, AwayScore: 7.00
Game ID: 32, HomeTeam: 12, AwayTeam: 13, HomeScore: 8.00, AwayScore: 7.00
Game ID: 33, HomeTeam: 8, AwayTeam: 16, HomeScore: 2.00, AwayScore: 4.00
Game ID: 34, HomeTeam: 10, AwayTeam: 15, HomeScore: 9.00, AwayScore: 7.00
Game ID: 35, HomeTeam: 11, AwayTeam: 14, HomeScore: 10.00, AwayScore: 9.00
Game ID: 36, HomeTeam: 12, AwayTeam: 13, HomeScore: 1.00, AwayScore: 5.00
Game ID: 38, HomeTeam: 8, AwayTeam: 15, HomeScore: 9.00, AwayScore: 10.00
Game ID: 39, HomeTeam: 11, AwayTeam: 14, HomeScore: 3.00, AwayScore: 7.00
Game ID: 40, HomeTeam: 12, AwayTeam: 13, HomeScore: 7.00, AwayScore: 6.00
Game ID: 42, HomeTeam: 16, AwayTeam: 10, HomeScore: 6.00, AwayScore: 6.00
Game ID: 43, HomeTeam: 8, AwayTeam: 14, HomeScore: 8.00, AwayScore: 4.00
Game ID: 44, HomeTeam: 12, AwayTeam: 13, HomeScore: 1.00, AwayScore: 10.00
Game ID: 46, HomeTeam: 16, AwayTeam: 10, HomeScore: 9.00, AwayScore: 3.00
Game ID: 47, HomeTeam: 15, AwayTeam: 11, HomeScore: 7.00, AwayScore: 6.00
Game ID: 48, HomeTeam: 8, AwayTeam: 13, HomeScore: 1.00, AwayScore: 7.00
Game ID: 50, HomeTeam: 16, AwayTeam: 10, HomeScore: 6.00, AwayScore: 4.00
Game ID: 51, HomeTeam: 15, AwayTeam: 11, HomeScore: 2.00, AwayScore: 5.00
Game ID: 52, HomeTeam: 14, AwayTeam: 12, HomeScore: 2.00, AwayScore: 2.00
Game ID: 53, HomeTeam: 8, AwayTeam: 9, HomeScore: 6.00, AwayScore: 1.00
Game ID: 54, HomeTeam: 16, AwayTeam: 10, HomeScore: 3.00, AwayScore: 4.00
Game ID: 55, HomeTeam: 15, AwayTeam: 11, HomeScore: 4.00, AwayScore: 6.00
Game ID: 56, HomeTeam: 14, AwayTeam: 12, HomeScore: 6.00, AwayScore: 6.00
Game ID: 58, HomeTeam: 8, AwayTeam: 10, HomeScore: 7.00, AwayScore: 8.00
Game ID: 59, HomeTeam: 15, AwayTeam: 11, HomeScore: 0.00, AwayScore: 7.00
Game ID: 60, HomeTeam: 14, AwayTeam: 12, HomeScore: 4.00, AwayScore: 5.00
Game ID: 62, HomeTeam: 9, AwayTeam: 16, HomeScore: 6.00, AwayScore: 4.00
Game ID: 63, HomeTeam: 8, AwayTeam: 11, HomeScore: 0.00, AwayScore: 0.00
Game ID: 64, HomeTeam: 14, AwayTeam: 12, HomeScore: 0.00, AwayScore: 9.00
Game ID: 66, HomeTeam: 9, AwayTeam: 16, HomeScore: 6.00, AwayScore: 6.00
Game ID: 67, HomeTeam: 10, AwayTeam: 15, HomeScore: 8.00, AwayScore: 9.00
Game ID: 68, HomeTeam: 8, AwayTeam: 12, HomeScore: 5.00, AwayScore: 6.00
Game ID: 70, HomeTeam: 9, AwayTeam: 16, HomeScore: 0.00, AwayScore: 8.00
Game ID: 71, HomeTeam: 10, AwayTeam: 15, HomeScore: 2.00, AwayScore: 6.00
Game ID: 72, HomeTeam: 11, AwayTeam: 14, HomeScore: 2.00, AwayScore: 2.00
```

You can also print statistics for the current tournament if we add the following code at the end of the main function.

```
	stats, err := tournify.GetGroupTournamentStats(tournament, 3, 0, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(tournify.PrintGroupTournamentStats(stats))
```

This will print the stats for the current tournament which will look similar to the following example

```
Group   Team    Played  Wins    Ties    Losses  +/-     Diff    Points
0       1       7       5       1       1       36/23   13      16
0       5       7       5       1       1       35/26   9       16
0       4       7       4       1       2       31/24   7       13
0       6       7       3       1       3       36/37   -1      10
0       2       7       3       0       4       41/43   -2      9
0       7       7       3       0       4       33/37   -4      9
0       0       7       1       2       4       25/33   -8      5
0       3       7       1       0       6       18/32   -14     3
1       16      9       5       2       2       53/33   20      17
1       11      9       5       2       2       48/38   10      17
1       12      9       5       2       2       45/45   0       17
1       15      9       5       0       4       48/54   -6      15
1       13      9       3       4       2       35/18   17      13
1       10      9       3       1       5       46/56   -10     10
1       -1      9       0       9       0       0/0     0       9
1       8       9       2       2       5       38/40   -2      8
1       9       9       1       5       3       15/31   -16     8
1       14      9       1       3       5       41/54   -13     6
```

Teams are ordered by the number of points. If a team has an equal number of points as another the diff value is used as a tie-breaker. If the diff value is also the same then the points scored against other teams is used.

Elimination tournaments are also possible. Here is an example of an elimination tournament with 8 teams:

```go
package main

import (
	"fmt"
	"github.com/tournify/tournify"
	"math"
	"math/rand"
	"time"
)

func main() {
	teams := []tournify.Team{
		{ID: 0},
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
		{ID: 6},
		{ID: 7},
	}

	teamInterfaces := make([]tournify.TeamInterface, len(teams))

	for i := range teams {
		teamInterfaces[i] = &teams[i]
	}

	tournament := tournify.CreateEliminationTournamentFromTeams(teamInterfaces)

	for _, game := range tournament.GetGames() {
		err := tournament.SetGameScore(game, randomEvenFloat(0.0, 100.0), randomEvenFloat(0.0, 100.0))
		if err != nil {
			panic(err)
		}
	}

	// The print method gives us a string representing the current tournament
	fmt.Println(tournament.Print())

	// Loop eliminated teams
	fmt.Println("Eliminated teams")
	for _, team := range tournament.GetEliminatedTeams() {
		fmt.Printf("Team ID: %d\n", team.GetID())
	}

	fmt.Println()

	// Loop remaining teams
	fmt.Println("Remaining teams")
	for _, team := range tournament.GetRemainingTeams() {
		fmt.Printf("Team ID: %d\n", team.GetID())
	}
}

func randomEvenFloat(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return math.RoundToEven(min + rand.Float64()*(max-min))
}
```

The output from the above code will be something like the following:

```bash
TournamentType: Elimination

Teams
Team ID: 0
Team ID: 1
Team ID: 2
Team ID: 3
Team ID: 4
Team ID: 5
Team ID: 6
Team ID: 7

Games
Game ID: 0, HomeTeam: 0, AwayTeam: 1, HomeScore: 85.00, AwayScore: 55.00
Game ID: 1, HomeTeam: 2, AwayTeam: 3, HomeScore: 14.00, AwayScore: 10.00
Game ID: 2, HomeTeam: 4, AwayTeam: 5, HomeScore: 53.00, AwayScore: 49.00
Game ID: 3, HomeTeam: 6, AwayTeam: 7, HomeScore: 95.00, AwayScore: 91.00
Game ID: 4, HomeTeam: 0, AwayTeam: 2, HomeScore: 34.00, AwayScore: 31.00
Game ID: 5, HomeTeam: 4, AwayTeam: 6, HomeScore: 0.00, AwayScore: 70.00
Game ID: 6, HomeTeam: 0, AwayTeam: 6, HomeScore: 13.00, AwayScore: 12.00

Eliminated teams
Team ID: 1
Team ID: 2
Team ID: 3
Team ID: 4
Team ID: 5
Team ID: 7

Remaining teams
Team ID: 0
Team ID: 6
```

Contributing
=

Please see [contributing](CONTRIBUTING.md).

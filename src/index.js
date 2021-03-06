// Collapse is needed for responsive navbar
import 'bootstrap/js/src/base-component'
import 'bootstrap/js/src/collapse'
import 'bootstrap/js/src/alert'

let createGroups = document.getElementById("create-groups");
let viewTournament = document.getElementById("view-tournament");
// If createGroups exists it means we are on the create tournament page
if (createGroups) {
    let tourType = document.getElementById("tourtype");
    let teamCount = document.getElementById("teamcount");
    let meetCount = document.getElementById("meetcount");
    let groupCount = document.getElementById("groupcount");
    let extraGroupFields = document.getElementById("extragroup")
    let advancedLink = document.getElementById("advlnk")
    let advancedArea = document.getElementById("advanced")
    let mixTeamsButton = document.getElementById("mix")
    let createButton = document.getElementById("create")
    let createForm = document.getElementById("create-tournament-form")
    tourType.addEventListener('change', renderCreateGroups);
    teamCount.addEventListener('change', renderCreateGroups);
    groupCount.addEventListener('change', renderCreateGroups);
    meetCount.addEventListener('change', renderCreateGroups);
    advancedLink.addEventListener('click', function (e) {
        e.preventDefault()
        if (advancedArea.style.display === "none") {
            advancedArea.style.display = "block"
        } else {
            advancedArea.style.display = "none"
        }
    })
    mixTeamsButton.addEventListener('click', function (e) {
        e.preventDefault()
        let teams = []
        document.querySelectorAll("input[name='team[]']").forEach(input => {
            teams.push(input.value)
        });
        teams = shuffle(teams)
        document.querySelectorAll("input[name='team[]']").forEach(input => {
            input.value = teams.pop()
        });
    })

    createForm.addEventListener('submit', function (e) {
        createButton.disabled = true
        createButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>\n' +
            'Loading...'
        return true
    })

    renderCreateGroups()

    function renderCreateGroups() {
        // Backup all current team names
        let teams = []
        document.querySelectorAll("input[name='team[]']").forEach(input => {
            teams.push(input.value)
        });
        let teamCount = document.getElementById("teamcount");
        // Begin by making the create groups div empty
        let teamCountInt = parseInt(teamCount.value, 10)
        createGroups.innerHTML = ""
        if (tourType.value === "0") {
            extraGroupFields.style.display = "block"
            if (teamCount.tagName === "SELECT") {
                let teamInput = document.createElement("input")
                teamInput.setAttribute("class", "form-control")
                teamInput.setAttribute("type", "number")
                teamInput.setAttribute("id", "teamcount")
                teamInput.setAttribute("name", "teamcount")
                teamInput.value = teamCount.value
                teamInput.addEventListener('change', renderCreateGroups);
                teamCount.replaceWith(teamInput)
            }
        } else {
            extraGroupFields.style.display = "none"
            if (teamCount.tagName === "INPUT") {
                let teamSelectOptions = ["4", "8", "16", "32", "64"];
                let teamSelect = document.createElement("select")
                for (let i = 0; i < teamSelectOptions.length; i++) {
                    let option = document.createElement("option");
                    option.value = teamSelectOptions[i];
                    option.text = teamSelectOptions[i];
                    if (teamSelectOptions[i] === "8") {
                        option.selected = true
                    }
                    teamSelect.appendChild(option);
                }
                teamSelect.setAttribute("class", "form-control")
                teamSelect.setAttribute("id", "teamcount")
                teamSelect.setAttribute("name", "teamcount")
                teamSelect.addEventListener('change', renderCreateGroups);
                teamCount.replaceWith(teamSelect)
            }
        }
        let maxTeams = 64
        if (teamCountInt > maxTeams) {
            teamCountInt = maxTeams
            teamCount.value = maxTeams
        }

        let groupCountInt = parseInt(groupCount.value, 10)

        if (tourType.value === "0") {
            if (teamCountInt < 2) {
                teamCountInt = 2
                teamCount.value = 2
            }
            let meetCountInt = parseInt(meetCount.value, 10)
            if (meetCountInt > 4) {
                meetCount.value = 4
            }
            if (meetCountInt < 1) {
                meetCount.value = 1
            }
            if (teamCountInt / groupCountInt < 2) {
                groupCountInt = 1
                groupCount.value = 1
            }
            if (groupCountInt < 1) {
                groupCountInt = 1
                groupCount.value = 1
            }
        } else {
            groupCountInt = 1
        }
        let tpg = teamCountInt / groupCountInt
        let group = 1
        let groups = []
        for (let i = 1; i <= teamCountInt; i++) {
            if (i === 1 || (i > (tpg * group) && teamCountInt !== i && tourType.value === "0")) {
                if (i >= (tpg * group) && teamCountInt !== i && tourType.value === "0") {
                    group++
                }
                groups.push(document.createElement("ul"))
                groups[group - 1].setAttribute("class", "list-unstyled")
            }
            let input = document.createElement("input")
            input.setAttribute("type", "text")
            input.setAttribute("name", "team[]")
            input.setAttribute("placeholder", teamLabel + " " + i)
            input.setAttribute("class", "form-control")
            let inputLi = document.createElement("li");
            inputLi.appendChild(input)
            groups[group - 1].appendChild(inputLi)
        }

        for (let x = 0; x < groups.length; x++) {
            let groupDiv = document.createElement("div")
            if (tourType.value === "0") {
                groupDiv.innerHTML = "<h2>" + groupLabel + " " + (x + 1) + "</h2>"
            } else {
                groupDiv.innerHTML = "<h2>" + teamsLabel + "</h2>"
            }
            groupDiv.appendChild(groups[x])
            createGroups.appendChild(groupDiv)
        }

        // We put all the teams back
        teams = teams.reverse()
        document.querySelectorAll("input[name='team[]']").forEach(input => {
            if (teams.length > 0) {
                input.value = teams.pop()
            }
        });
    }
} else if (viewTournament) {
    let statisticsBlock = document.getElementById("group-tournament-statistics")
    let eliminationGamesBlock = document.getElementById("elimination-games")
    document.querySelectorAll("input[name='home']").forEach(input => {
        input.addEventListener("change", updateGame)
    });
    document.querySelectorAll("input[name='away']").forEach(input => {
        input.addEventListener("change", updateGame)
    });

    function updateGame(e) {
        let data = {};
        let row = e.target.closest("tr")
        let inputs = row.getElementsByTagName('input');
        for (let z = 0; z < inputs.length; z++) {
            data[inputs[z].name] = inputs[z].value
        }
        fetch("/api/tournament/" + data["slug"] + "/game/" + data["id"], {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        }).then(res => {
            res.text().then(result => {
                if (data["tourType"] === "0") {
                    fetch("/api/tournament/" + data["slug"] + "/stats", {
                        method: "GET",
                    }).then(res => {
                        res.text().then(result => {
                            const obj = JSON.parse(result);
                            let newHtml = "<table class=\"table table-hover table-sm\">"
                            for (const [key, value] of Object.entries(obj)) {
                                newHtml += "<thead>"
                                newHtml += "<tr>"
                                newHtml += "<th scope=\"col\">" + groupLabel + " " + key + "</th>"
                                newHtml += "<th scope=\"col\">" + playedLabel + "</th>"
                                newHtml += "<th scope=\"col\">" + winsLabel + "</th>"
                                newHtml += "<th scope=\"col\">" + tiesLabel + "</th>"
                                newHtml += "<th scope=\"col\">" + lostLabel + "</th>"
                                newHtml += "<th scope=\"col\">+/-</th>"
                                newHtml += "<th scope=\"col\">" + diffLabel + "</th>"
                                newHtml += "<th scope=\"col\">" + pointsLabel + "</th>"
                                newHtml += "</tr>"
                                newHtml += "</thead>"
                                newHtml += "</tbody>"
                                for (const [statKey, statValue] of Object.entries(value.stats)) {
                                    newHtml += "<tr>"
                                    newHtml += "<th scope=\"row\">" + statValue.team.name + "</th>"
                                    newHtml += "<td>" + statValue.played + "</td>"
                                    newHtml += "<td>" + statValue.wins + "</td>"
                                    newHtml += "<td>" + statValue.ties + "</td>"
                                    newHtml += "<td>" + statValue.losses + "</td>"
                                    newHtml += "<td>" + statValue.points_for + "/" + statValue.points_against + "</td>"
                                    newHtml += "<td>" + (statValue.points_for - statValue.points_against) + "</td>"
                                    newHtml += "<td>" + statValue.points + "</td>"
                                    newHtml += "</tr>"
                                }
                                newHtml += "</tbody>"
                            }
                            newHtml += "</table>"
                            statisticsBlock.innerHTML = newHtml
                        })
                    });
                } else if (data["tourType"] === "1") {
                    fetch("/api/tournament/" + data["slug"] + "/games", {
                        method: "GET",
                    }).then(res => {
                        res.text().then(result => {
                            const obj = JSON.parse(result);
                            let newHtml = "<table class=\"table table-hover table-sm\" id=\"elimination-games\">"
                            for (const [key, value] of Object.entries(obj)) {
                                newHtml += "<thead>"
                                newHtml += "<tr>"
                                newHtml += "<th scope=\"col\">" + value.name + "</th>"
                                newHtml += "<th scope=\"col\" colSpan=\"2\">" + homeTeamLabel + "</th>"
                                newHtml += "<th scope=\"col\" colSpan=\"2\">" + awayTeamLabel + "</th>"
                                newHtml += "</tr>"
                                newHtml += "</thead>"
                                newHtml += "</tbody>"
                                for (const [gameKey, gameValue] of Object.entries(value.games)) {
                                    newHtml += "<tr>"
                                    newHtml += "<th scope=\"row\">"
                                    newHtml += gameValue.name
                                    newHtml += "<input type=\"hidden\" value=\"" + gameValue.ID + "\" name=\"id\">"
                                    newHtml += "<input type=\"hidden\" value=\"" + data["tourType"] + "\" name=\"tourType\">"
                                    newHtml += "<input type=\"hidden\" value=\"" + data["slug"] + "\" name=\"slug\">"
                                    newHtml += "</th>"
                                    newHtml += "<td>"
                                    if (gameValue.teams.length > 0) {
                                        newHtml += gameValue.teams[0].name
                                    }
                                    newHtml += "</td>"
                                    newHtml += "<td>"
                                    if (gameValue.teams.length > 0) {
                                        if (canEdit === "true") {
                                            if (gameValue.scores.length > 0) {
                                                newHtml += "<input class=\"small-input\" type=\"number\" name=\"home\" value=\"" + gameValue.scores[0].score + "\">"
                                            } else {
                                                newHtml += "<input class=\"small-input\" type=\"number\" name=\"home\" value=\"0\">"
                                            }
                                        } else {
                                            if (gameValue.scores.length > 0) {
                                                newHtml += gameValue.scores[0].points
                                            } else {
                                                newHtml += "0"
                                            }
                                        }
                                    } else {
                                        if (canEdit === "true") {
                                            newHtml += "<input class=\"small-input\" type=\"number\" name=\"home\" value=\"0\">"
                                        } else {
                                            newHtml += "0"
                                        }
                                    }
                                    newHtml += "</td>"
                                    newHtml += "<td>"
                                    if (gameValue.teams.length > 1) {
                                        newHtml += gameValue.teams[1].name
                                    }
                                    newHtml += "</td>"
                                    newHtml += "<td>"
                                    if (gameValue.teams.length > 1) {
                                        if (canEdit === "true") {
                                            if (gameValue.scores.length > 1) {
                                                newHtml += "<input class=\"small-input\" type=\"number\" name=\"away\" value=\"" + gameValue.scores[1].score + "\">"
                                            } else {
                                                newHtml += "<input class=\"small-input\" type=\"number\" name=\"away\" value=\"0\">"
                                            }
                                        } else {
                                            if (gameValue.scores.length > 1) {
                                                newHtml += gameValue.scores[1].points
                                            } else {
                                                newHtml += "0"
                                            }
                                        }
                                    } else {
                                        if (canEdit === "true") {
                                            newHtml += "<input class=\"small-input\" type=\"number\" name=\"away\" value=\"0\">"
                                        } else {
                                            newHtml += "0"
                                        }
                                    }
                                    newHtml += "</td>"
                                    newHtml += "</tr>"
                                }
                                newHtml += "</tbody>"
                            }
                            newHtml += "</table>"
                            eliminationGamesBlock.innerHTML = newHtml
                            document.querySelectorAll("input[name='home']").forEach(input => {
                                input.addEventListener("change", updateGame)
                            });
                            document.querySelectorAll("input[name='away']").forEach(input => {
                                input.addEventListener("change", updateGame)
                            });
                        })
                    });
                }
            })

        });
    }
}

// from https://stackoverflow.com/a/2450976/1260548
function shuffle(array) {
    let currentIndex = array.length, randomIndex;

    // While there remain elements to shuffle...
    while (currentIndex !== 0) {

        // Pick a remaining element...
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;

        // And swap it with the current element.
        [array[currentIndex], array[randomIndex]] = [
            array[randomIndex], array[currentIndex]];
    }

    return array;
}
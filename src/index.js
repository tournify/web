// Collapse is needed for responsive navbar
import 'bootstrap/js/src/base-component'
import 'bootstrap/js/src/collapse'
import 'bootstrap/js/src/alert'

let createGroups = document.getElementById("create-groups");
let viewTournament = document.getElementById("view-tournament");
// If createGroups exists it means we are on the create tournament page
if(createGroups){
    let tourType = document.getElementById("tourtype");
    let teamCount = document.getElementById("teamcount");
    let meetCount = document.getElementById("meetcount");
    let groupCount = document.getElementById("groupcount");
    let extraGroupFields = document.getElementById("extragroup")
    let advancedLink = document.getElementById("advlnk")
    let advancedArea = document.getElementById("advanced")
    let mixTeamsButton = document.getElementById("mix")
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
        document.querySelectorAll("input[name='team[]']").forEach( input => {
            teams.push(input.value)
        });
        teams = shuffle(teams)
        document.querySelectorAll("input[name='team[]']").forEach( input => {
            input.value = teams.pop()
        });
    })

    renderCreateGroups()

    function renderCreateGroups(){
        // Backup all current team names
        let teams = []
        document.querySelectorAll("input[name='team[]']").forEach( input => {
            teams.push(input.value)
        });
        // Begin by making the create groups div empty
        createGroups.innerHTML = ""
        if (tourType.value === "0") {
            extraGroupFields.style.display = "block"
        } else {
            extraGroupFields.style.display = "none"
        }
        let teamCountInt = parseInt(teamCount.value, 10)
        if (teamCountInt > 10) {
            teamCountInt = 10
            teamCount.value = 10
        }
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
        let groupCountInt = parseInt(groupCount.value, 10)
        if (groupCountInt < 1) {
            groupCountInt = 1
            groupCount.value = 1
        }
        let tpg = teamCountInt/groupCountInt
        let group = 1
        let groups = []
        for (let i = 1; i <= teamCountInt; i++) {
            if (i === 1 || (i > (tpg * group) && teamCountInt !== i && tourType.value === "0")) {
                if (i >= (tpg * group) && teamCountInt !== i && tourType.value === "0") {
                    group++
                }
                groups.push(document.createElement("ul"))
                groups[group-1].setAttribute("class", "list-unstyled")
            }
            let input = document.createElement("input")
            input.setAttribute("type", "text")
            input.setAttribute("name", "team[]")
            input.setAttribute("placeholder", "Team " + i)
            let inputLi = document.createElement("li");
            inputLi.appendChild(input)
            groups[group-1].appendChild(inputLi)
        }

        for (let x = 0; x < groups.length; x++) {
            let groupDiv = document.createElement("div")
            if (tourType.value === "0") {
                groupDiv.innerHTML = "<h2>Group " + (x + 1) + "</h2>"
            } else {
                groupDiv.innerHTML = "<h2>Teams</h2>"
            }
            groupDiv.appendChild(groups[x])
            createGroups.appendChild(groupDiv)
        }

        // We put all the teams back
        teams = teams.reverse()
        document.querySelectorAll("input[name='team[]']").forEach( input => {
            if (teams.length > 0) {
                input.value = teams.pop ()
            }
        });
    }
} else if (viewTournament) {
    let statisticsBlock = document.getElementById("group-tournament-statistics")
    document.querySelectorAll("input[name='home']").forEach( input => {
        input.addEventListener("change", updateGame)
    });
    document.querySelectorAll("input[name='away']").forEach( input => {
        input.addEventListener("change", updateGame)
    });

    function updateGame(e) {
        let data = {};
        let row = e.target.closest("tr")
        let inputs = row.getElementsByTagName( 'input' );
        for ( let z = 0; z < inputs.length; z++ ) {
            data[inputs[z].name] = inputs[z].value
        }
        fetch("/api/tournament/" + data["slug"] + "/game/" + data["id"], {
            method: "POST",
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        }).then(res => {
            res.text().then(result => {
                fetch("/api/tournament/" + data["slug"] + "/stats", {
                    method: "GET",
                }).then(res => {
                    res.text().then(result => {
                        const obj = JSON.parse(result);
                        let newHtml = "<table class=\"table table-hover table-sm\">"
                        for (const [key, value] of Object.entries(obj)) {
                            newHtml += "<thead>"
                            newHtml += "<tr>"
                            newHtml += "<th scope=\"col\">Group " + key + "</th>"
                            newHtml += "<th scope=\"col\">Team</th>"
                            newHtml += "<th scope=\"col\">Played</th>"
                            newHtml += "<th scope=\"col\">Wins</th>"
                            newHtml += "<th scope=\"col\">Ties</th>"
                            newHtml += "<th scope=\"col\">Lost</th>"
                            newHtml += "<th scope=\"col\">+/-</th>"
                            newHtml += "<th scope=\"col\">Diff</th>"
                            newHtml += "<th scope=\"col\">Points</th>"
                            newHtml += "</tr>"
                            newHtml += "</thead>"
                            newHtml += "</tbody>"
                            for (const [statKey, statValue] of Object.entries(value.stats)) {
                                newHtml += "<tr>"
                                newHtml += "<td></td>"
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
            })

        });
    }
}

// from https://stackoverflow.com/a/2450976/1260548
function shuffle(array) {
    let currentIndex = array.length,  randomIndex;

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
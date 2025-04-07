const alertURL = "http://localhost:8080/alerts";
const cycleURL = "http://localhost:8080/cycles"
const archonURL = "http://localhost:8080/archon"
const top5warframe = `http://localhost:8080/search`;
const newsURL = `http://localhost:8080/news`;
newsImageArray = [ ]


// TODO: Sort news in order of newest -> oldest
async function fetchAndDisplayNews() {
    try {
        const response = await fetch(newsURL);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        const news = await response.json();

        const newsList = document.getElementById("newsList");

        newsList.innerHTML = "";

        news.news.forEach(ev => {
            const rawTimeString = ev.asString
            const timeString = extractDaysAndHours(rawTimeString);
            const message = ev.message
            const link = ev.link
            const imageLink = ev.imageLink

            const listItem = document.createElement("li");
            listItem.classList.add("compact-li");
            const anchor = document.createElement("a");
            const timeAnchor = document.createElement("a");
            timeAnchor.textContent = timeString + ":  ";
            timeAnchor.style.cssText = 'font-size: 12px;'

            anchor.href = link;
            anchor.textContent = message;
            anchor.target = "_blank";

            listItem.appendChild(timeAnchor);
            listItem.appendChild(anchor);

            newsList.append(listItem);

            newsImageArray.push({key: message, value: imageLink});
        });
    }
    catch (error){
        console.error("Failed to fetch news:", error);
    }
}

// Rotates the news image every 5 seconds
function rotateNewsImage(newsImageArray) {
    const newsImage = document.getElementById("newsImage");
    // Load the first image NOTE: Not working
    newsImage.innerHTML = `<img class="news" src="${newsImageArray[0]}">`;

        let index = 0;
        setInterval(() => {
            const imageLink = newsImageArray[index].value;
            const keyText = newsImageArray[index].key;
            newsImage.innerHTML = `<img class="news" src="${imageLink}">`;
            const linkElement = [...document.querySelectorAll('a')].find(a => a.textContent === keyText);

            document.querySelectorAll('a').forEach(a => a.style.color = '');

            if (linkElement) {
                linkElement.style.color = 'white';
            }
            index = (index + 1) % newsImageArray.length;
        }, 5000);
     
}

async function fetchAndDisplayAlerts() {
    try {
        const response = await fetch(alertURL);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        const alerts = await response.json();

        const alertList = document.getElementById("alertList");
        alertList.innerHTML = "";

        if (alerts.length === 0) {
            const noItems = document.createElement("li");
            noItems.textContent = 'No alerts present';
            alertList.appendChild(noItems);
        }

        else{
        alerts.forEach(alert => {
            const type = alert.mission.typeKey || "Unknown type";
            const time = alert.eta || "Unknown time";
            const node = alert.mission.nodeKey || "Unknown node";
            const reward = alert.mission.reward.asString || "No reward";
            const minLevel = alert.mission.minEnemyLevel || "Unknown level";
            const maxLevel = alert.mission.maxEnemyLevel || "Unknown level";

            const listItem = document.createElement("li");
            listItem.textContent = `Type: ${type}, Time: ${time}, Node: ${node}, Reward: ${reward}, Level: ${minLevel}-${maxLevel}`;
            alertList.appendChild(listItem);
        });
    }
    } catch (error) {
        console.error("Failed to fetch alerts:", error);
    }
    
}

async function dayNightCycles() {
    try {
        const response = await fetch(cycleURL);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        const cycles = await response.json();

        const cycleList = document.getElementById("cycleList");
        cycleList.innerHTML = "";

        const state = cycles.state || "State not found";
        const timeString = cycles.shortString || "Not found";

        const listItem = document.createElement("li");
        listItem.textContent = `Cetus: ${state}, ${timeString}`;
        cycleList.appendChild(listItem);

    } catch (error) {
        console.error("Failed to fetch alerts:", error);
    }
}

async function archonHunt() {
    try {
        const response = await fetch(archonURL);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        const archonHunt = await response.json();

        const archonList = document.getElementById("archonList");
        const archonName = document.getElementById("archonName");

        archonList.innerHTML = "";

        const boss = archonHunt.boss || "Not found";
        const time = archonHunt.eta || "Not found";
        const timeWithoutSeconds = time.replace(/\s?\d+s$/, '');

        archonName.innerHTML = `<img src="IconNarmer.webp" style="width: 25px; height: 25px"> ${boss}, ${timeWithoutSeconds}`;

        archonHunt.missions.forEach(mission => {
            const missionItem = document.createElement("li");
            missionItem.textContent = `${mission.nodeKey || "Not found"}, ${mission.typeKey || "Not found"}`;
            archonList.appendChild(missionItem);
        });

    } catch (error) {
        console.error("Error fetching Archon Hunt data:", error);
    }
}

const checkbox = document.getElementById("archonCheckbox");

document.addEventListener("DOMContentLoaded", () => {
    const checkbox = document.getElementById("archonCheckbox");

    if (checkbox) {
        const isChecked = localStorage.getItem("archonCheckbox") === "true";
        checkbox.checked = isChecked;

        checkbox.addEventListener("change", () => {
            localStorage.setItem("archonCheckbox", checkbox.checked);
        });
    }
});

function getNextMondayAtMidnightUTC() {
    const now = new Date();
    const dayOfWeek = now.getUTCDay(); 
    const daysUntilNextMonday = (7 - dayOfWeek) % 7;
    const nextMonday = new Date(now);
    
    nextMonday.setUTCDate(now.getUTCDate() + daysUntilNextMonday); 
    nextMonday.setUTCHours(0, 0, 0, 0);
    
    return nextMonday;
}

// Resets archon hunt checkbox
function resetCheckBox(){
    localStorage.removeItem('archonCheckbox');
}

function setResetTimer() {
    const now = new Date();
    const nextMonday = getNextMondayAtMidnightUTC();
    
    const timeUntilNextMonday = nextMonday.getTime() - now.getTime();

    setTimeout(() => {
        resetCheckBox();
        setInterval(resetLocalStorage, 7 * 24 * 60 * 60 * 1000); // 7 days in milliseconds
    }, timeUntilNextMonday);
}

function restoreCheckBox()
{
    const checkbox = document.getElementById('archonCheckbox');
    const savedstate = localStorage.getItem('archonCheckbox');

    // Checkbox state is saved as a string
    if (savedstate === 'true')
    {
        checkbox.checked = true;
    }
}

async function fetchTop5(){
    const searchField = document.getElementById('searchField');
    const searchQuery = searchField.value.trim();

    if (!searchQuery) {
        alert("Please enter an in-game name.");
        return;
    }

    try {

        const response = await fetch(top5warframe,{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ query: searchQuery })
        });

        if (!response.ok) {
            throw new Error(`Error fetching data: ${response.statusText}`);
        }

        const data = await response.json();

        const resultsDiv = document.getElementById('top5Results');
        resultsDiv.innerHTML = "<h4>Top 5 Warframes:</h4>";
        
        const ul = document.createElement('ul');
        let listIndex = 1;
        data.forEach(warframe => {
            const li = document.createElement('li');
            li.textContent = `${listIndex}. ${warframe.cleanedName} - xp: ${warframe.xp}`;
        
            const details = document.createElement('details');
            const summary = document.createElement('summary');
            summary.textContent = 'Abilities';
        
            const passive = document.createElement('p');
            const passiveBold = document.createElement('strong');
            passiveBold.textContent = 'Passive: '
            const passiveName = document.createTextNode(warframe.passiveDescription);

            passive.appendChild(passiveBold);
            passive.appendChild(passiveName);

        
            details.appendChild(summary);
            details.appendChild(passive);
        
            warframe.abilities.forEach(ability => {
                const p = document.createElement('p');

                const abilityName = document.createElement('strong');
                abilityName.textContent = ability.abilityName + ': ';

                const abilityDescription = document.createTextNode(ability.description);

                p.appendChild(abilityName);
                p.appendChild(abilityDescription);

                details.appendChild(p);
            });
        
            li.appendChild(details);
            ul.appendChild(li);
            listIndex++;
        });
        resultsDiv.appendChild(ul);

    } catch (error) {
        console.error("Error:", error);
        alert("An error occurred while fetching data. Please try again.");
    }
    }
    

setResetTimer();


window.onload = function() {
    loadPageData();
};

function loadPage(page) {
    fetch(page)
        .then(response => response.text())
        .then(data => {
            const contentStart = data.indexOf('<body>');
            const contentEnd = data.indexOf('</body>');

            if (contentStart !== -1 && contentEnd !== -1) {
                const pageContent = data.slice(contentStart + 6, contentEnd);
                document.getElementById('content').innerHTML = pageContent;
            }

            if (page === 'index.html') {
                loadPageData();
            }
        })
        .catch(error => {
            document.getElementById('content').innerHTML = '<h1>Page not found!</h1>';
            console.error('Error loading page:', error);
        });
}

// Update data every minute
setInterval(() => {
    fetchAndDisplayAlerts();
    dayNightCycles();
    archonHunt();
}, 60000);

function loadPageData(){
    fetchAndDisplayAlerts();
    dayNightCycles();
    archonHunt();
    fetchAndDisplayNews();
    rotateNewsImage(newsImageArray);
}

function extractDaysAndHours(str) {
    const regex = /\[(.*?)\]/;

    const match = str.match(regex);

    return match ? match[1] : null;
}
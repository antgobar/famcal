function getCalMembers() {
    fetch('http://localhost:8080/members')
        .then(response => {
            if (!response.ok) {u
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            const members = document.getElementById('memberList');
            members.innerHTML = ''; 
            if (data.length === 0) {
                members.innerHTML = '<p>No members found.</p>';
                return;
            }

            data.forEach(member => {
                const article = document.createElement('article');
                article.classList.add('card'); 
                article.textContent = member.name;
                article.id = `member-${member.id}`;
                members.appendChild(article);
            });
        })
        .catch(error => console.error('Error fetching members:', error));
}

function addCalMemberOnSubmit() {
    const form = document.getElementById('addCalMemberForm');
    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const memberName = document.getElementById('name').value;
        if (!memberName.trim()) {
            alert('CalMember name cannot be empty.');
            return;
        }

        fetch('http://localhost:8080/members', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: memberName })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log('CalMember added:', data);
                document.getElementById('name').value = ''; // Reset input field
                getCalMembers(); // Refresh the member list
            })
            .catch(error => console.error('Error adding member:', error));
    });
}

function getCalendars() {
    fetch("http://localhost:8080/calendars")
        .then(response => response.json())
        .then(calendars => {
            const calendarsListDiv = document.getElementById("calendarsList");
            calendarsListDiv.innerHTML = "";

            calendars.forEach(calendar => {
                const calendarDiv = document.createElement("div");
                const calendarHTML = `
                    <article id="${calendar.id}" class="calendar-item">
                        <strong>${calendar.summary}</strong>: ${calendar.description || 'N/A'}</p>
                    </article>
                `;
                calendarDiv.innerHTML = calendarHTML;
                calendarsListDiv.appendChild(calendarDiv);
                calendarDiv.querySelector(".calendar-item").addEventListener("click", function () {
                    fetchEvents(this.id);
                });
            });

        })
        .catch(error => {
            console.error("Error fetching calendars:", error);
        });
}

function fetchEvents(calendarId) {
    eventsContainer = document.getElementById("eventsList");
    fetch(`http://localhost:8080/events?calendarId=${calendarId}&nEvents=10`)
        .then(response => response.json())
        .then(events => {
            eventsContainer.innerHTML = "";

            events.forEach(event => {
                const eventDiv = document.createElement("div");
                const eventDate = extractDateFromISO(event.start);
                const startTime = extractTimeFromISO(event.start);
                const endTime = extractTimeFromISO(event.end);
                const eventHTML = `
                    <article class="event-item">
                        <strong>${event.summary}</strong>: ${eventDate}
                        <p><strong>Start:</strong> ${startTime} <strong>End:</strong> ${endTime}</p>
                    </article>
                `;
                eventDiv.innerHTML = eventHTML;
                eventsContainer.appendChild(eventDiv);
            });
        })
        .catch(error => {
            console.error(`Error fetching events for calendar ${calendarId}:`, error);
        });
}

function extractDateFromISO(isoDateString) {
    const date = new Date(isoDateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Zero-padding for single digit months
    const day = String(date.getDate()).padStart(2, '0'); // Zero-padding for single digit days
    return `${year}-${month}-${day}`;
}

function extractTimeFromISO(isoDateString) {
    const date = new Date(isoDateString);
    const hours = String(date.getHours()).padStart(2, '0');  // Zero-padding for single digit hours
    const minutes = String(date.getMinutes()).padStart(2, '0'); // Zero-padding for single digit minutes
    return `${hours}:${minutes}`;
}

// Initialize
getCalendars();
getCalMembers();
addCalMemberOnSubmit();

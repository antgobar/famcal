function getCalMembers() {
    fetch(`${API_BASE_URL}/members`)
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

        fetch(`${API_BASE_URL}/members`, {
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
                document.getElementById('name').value = '';
                getCalMembers();
            })
            .catch(error => console.error('Error adding member:', error));
    });
}

function getCalendars() {
    fetch(`${API_BASE_URL}/calendars`)
        .then(response => response.json())
        .then(calendars => {
            const selectElement = document.getElementById("calendarSelect");
            selectElement.innerHTML = "";

            if (calendars.length > 0) {
                const defaultOption = document.createElement("option");
                defaultOption.selected = true;
                defaultOption.value = calendars[0].id;
                defaultOption.textContent = calendars[0].summary;
                defaultOption.setAttribute("data-calendarDescription", calendars[0].description || calendars[0].summary);
                selectElement.appendChild(defaultOption);
                fetchEvents(calendars[0].id);
                const descriptionElement = document.getElementById("calendarDescription");
                descriptionElement.textContent = calendars[0].description || calendars[0].summary;

                calendars.slice(1).forEach(calendar => {
                    const option = document.createElement("option");
                    option.value = encodeURIComponent(calendar.id);
                    option.textContent = calendar.summary;
                    option.setAttribute("data-calendarDescription", calendar.description || calendar.summary);
                    selectElement.appendChild(option);
                });

                selectElement.addEventListener("change", function () {
                    const selectedCalendarId = this.value;
                    if (selectedCalendarId) {
                        fetchEvents(selectedCalendarId);
                        const selectedOption = this.selectedOptions[0];
                        const calendarDescription = selectedOption.getAttribute("data-calendarDescription");
                        const descriptionElement = document.getElementById("calendarDescription");
                        descriptionElement.textContent = calendarDescription;
                    }
                });
            }

        })
        .catch(error => {
            console.error("Error fetching calendars:", error);
        });
}

function fetchEvents(calendarId) {
    eventsContainer = document.getElementById("eventsList");
    fetch(`${API_BASE_URL}/events?calendarId=${calendarId}&nEvents=10`)
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

function closeDialogById(dialogId) {
    const dialog = document.getElementById(dialogId);
    if (dialog) {
        dialog.close();
    } else {
        console.error(`Dialog with id ${dialogId} not found.`);
    }
}

// Initialize
getCalendars();
getCalMembers();
addCalMemberOnSubmit();

// Utility Functions
const fetchJSON = async (url, options = {}) => {
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`Error fetching from ${url}:`, error);
        throw error;
    }
};

const clearElement = (element) => {
    element.innerHTML = '';
};

// Fetch Members
const getCalMembers = async () => {
    const membersContainer = document.getElementById('memberList');
    try {
        const members = await fetchJSON(`${API_BASE_URL}/members`);
        clearElement(membersContainer);

        if (members.length === 0) {
            membersContainer.innerHTML = '<p>No members found.</p>';
            return;
        }

        members.forEach(({ id, name }) => {
            const article = document.createElement('article');
            article.classList.add('card');
            article.textContent = name;
            article.id = `member-${id}`;
            membersContainer.appendChild(article);
        });
    } catch (error) {
        console.error('Error fetching members:', error);
    }
};

// Add Member
const addCalMemberOnSubmit = () => {
    const form = document.getElementById('addCalMemberForm');
    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const memberName = form.name.value.trim();
        if (!memberName) {
            alert('Member name cannot be empty.');
            return;
        }

        try {
            await fetchJSON(`${API_BASE_URL}/members`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: memberName })
            });

            console.log('Member added:', memberName);
            form.reset();
            getCalMembers();
        } catch (error) {
            console.error('Error adding member:', error);
        }
    });
};

// Fetch Calendars
const getCalendars = async () => {
    const selectElement = document.getElementById('calendarSelect');
    const descriptionElement = document.getElementById('calendarDescription');
    clearElement(selectElement);

    try {
        const calendars = await fetchJSON(`${API_BASE_URL}/calendars`);
        if (calendars.length === 0) return;

        calendars.forEach(({ id, summary, description }, index) => {
            const option = document.createElement('option');
            option.value = encodeURIComponent(id);
            option.textContent = summary;
            option.dataset.description = description || summary;

            if (index === 0) {
                option.selected = true;
                descriptionElement.textContent = description || summary;
                fetchEvents(id);
            }

            selectElement.appendChild(option);
        });

        selectElement.addEventListener('change', ({ target }) => {
            const selectedOption = target.selectedOptions[0];
            const calendarId = target.value;
            descriptionElement.textContent = selectedOption.dataset.description;
            fetchEvents(calendarId);
        });
    } catch (error) {
        console.error('Error fetching calendars:', error);
    }
};

// Fetch Events
const fetchEvents = async (calendarId) => {
    const eventsContainer = document.getElementById('eventsList');
    clearElement(eventsContainer);

    try {
        const events = await fetchJSON(`${API_BASE_URL}/events?calendarId=${calendarId}&nEvents=10`);

        events.forEach(({ summary, start, end }) => {
            const eventHTML = `
                <article class="event-item">
                    <strong>${summary}</strong>: ${extractDateFromISO(start)}
                    <p><strong>Start:</strong> ${extractTimeFromISO(start)} 
                       <strong>End:</strong> ${extractTimeFromISO(end)}</p>
                </article>`;
            const eventDiv = document.createElement('div');
            eventDiv.innerHTML = eventHTML;
            eventsContainer.appendChild(eventDiv);
        });
    } catch (error) {
        console.error(`Error fetching events for calendar ${calendarId}:`, error);
    }
};

// ISO Date Helpers
const extractDateFromISO = (isoDate) => new Date(isoDate).toISOString().split('T')[0];
const extractTimeFromISO = (isoDate) => new Date(isoDate).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

// Dialog Controls
const setupDialogControls = () => {
    const dialog = document.getElementById('add-member-dialog');
    const openButton = document.getElementById('openMemberDialogButton');
    const closeButton = document.getElementById('closeMemberDialogButton');

    // Check if the buttons exist before adding event listeners
    if (openButton && closeButton) {
        openButton.addEventListener('click', () => dialog.showModal());
        closeButton.addEventListener('click', () => dialog.close());
    } else {
        console.error("Dialog control buttons not found.");
    }
};

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    getCalendars();
    getCalMembers();
    addCalMemberOnSubmit();
    setupDialogControls();
});

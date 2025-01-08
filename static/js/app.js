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
                article.classList.add('card'); // Use Pico's styling
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

document.getElementById("fetchCalendarsBtn").addEventListener("click", function() {
    fetch("http://localhost:8080/calendars")
        .then(response => response.json())
        .then(calendars => {
            const calendarsListDiv = document.getElementById("calendarsList");
            calendarsListDiv.innerHTML = "";  // Clear previous content

            calendars.forEach(calendar => {
                const calendarDiv = document.createElement("div");

                const calendarHTML = `
                    <article>
                    <h3>${calendar.summary}</h3>
                    <p><strong>ID:</strong> ${calendar.id}</p>
                    <p><strong>Description:</strong> ${calendar.description || 'N/A'}</p>
                    <p><strong>Time Zone:</strong> ${calendar.timeZone || 'N/A'}</p>
                    <p><strong>Location:</strong> ${calendar.location || 'N/A'}</p>
                    </article>
                `;
                calendarDiv.innerHTML = calendarHTML;
                calendarsListDiv.appendChild(calendarDiv);
            });
        })
        .catch(error => {
            console.error("Error fetching calendars:", error);
        });
});

// Initialize
getCalMembers();
addCalMemberOnSubmit();

let scoresData = {};

function loadScores(diveId) {
    console.log("Loading scores for dive:", diveId);
    const userId = document.querySelector('#scoreForm input[name="userId"]').value;
    const eventId = document.querySelector('#scoreForm input[name="eventId"]').value;

    if (!scoresData[diveId]) {
        scoresData[diveId] = { scores: [] };
        fetch(`/user/${userId}/event/${eventId}/dive/${diveId}/scores`)
            .then(function (response) {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Failed to fetch scores.');
                }
            })
            .then(function (data) {
                console.log("Fetched scores:", data);
                scoresData[diveId].scores = data.scores;
                updateScores(diveId);
            })
            .catch(function (error) {
                console.error('Error fetching scores:', error);
            });
    } else {
        updateScores(diveId);
    }
}

function updateScores(diveId) {
    scoresData[diveId].scores.forEach(function (scoreObject, index) {
        const score = scoreObject.Value;
        const scoreElement = document.querySelector(`[data-dive-id="${diveId}"] .score${index + 1}`);
        if (scoreElement) {
            scoreElement.textContent = score.toFixed(2);
        }
    });
}

document.addEventListener("DOMContentLoaded", function() {
    const scoreForm = document.getElementById("scoreForm");
    const scoreModal = document.getElementById("scoreModal");
    const errorAlert = document.getElementById("errorAlert");

    // Handle form submission
    scoreForm.addEventListener("submit", function(event) {
        event.preventDefault();
        console.log("Form submitted"); // Add this line

        // Get the form data
        const formData = new FormData(scoreForm);
        const diveId = parseInt(formData.get("diveId"), 10);
        const userId = parseInt(formData.get("userId"), 10);
        const eventId = parseInt(formData.get("eventId"), 10);

        // Get scores from the form data
        const scores = [];
        for (let i = 1; i <= 9; i++) {
            const score = formData.get("score" + i);
            if (score) {
                scores.push(parseFloat(score));
            }
        }

        // Send a request to the server to save the scores
        fetch(`/user/${userId}/event/${eventId}/scores`, {
            method: "POST",
            body: JSON.stringify({
                userId: userId,
                eventId: eventId,
                diveId: diveId,
                scores: scores
            }),
            headers: {
                "Content-Type": "application/json"
            }
        }).then(function(response) {
            if (response.ok) {
                return response.json(); // Return the response as a Promise
            } else {
                throw new Error("An error occurred while submitting the score.");
            }
        }).then(function(data) {
            // Close the modal and clear the form
            const modalInstance = bootstrap.Modal.getInstance(scoreModal);
            modalInstance.hide();
            scoreForm.reset();

            // Update the scores on the page
            scores.forEach(function(score, index) {
                const scoreElement = document.querySelector(`[data-dive-id="${diveId}"] [data-score-index="${index + 1}"]`);
                if (scoreElement) {
                    scoreElement.textContent = score.toFixed(2);
                } else {
                    const newScoreElement = document.createElement("div");
                    newScoreElement.textContent = score.toFixed(2);
                    newScoreElement.setAttribute("data-dive-id", diveId);
                    newScoreElement.setAttribute("data-score-index", index + 1);
                    scoreElement.parentNode.insertBefore(newScoreElement, scoreElement);
                }
            });

            // Hide the error alert
            errorAlert.style.display = "none";
        }).catch(function(error) {
            // Handle error
            // Display the error message
            errorAlert.textContent = error.message;
            errorAlert.style.display = "block";
        });
    });

    // Set the dive ID in the form when the Add Score button is clicked
    const addScoreButtons = document.querySelectorAll(".add-score");
    addScoreButtons.forEach(function(button) {
        button.addEventListener("click", function() {
            const diveId = button.getAttribute("data-dive-id");
            console.log("Add Score button clicked, diveId:", diveId); // Add this line
            scoreForm.querySelector("#diveId").value = diveId;
        });
    });

    // Load scores for all dives
    document.querySelectorAll('[data-dive-id]').forEach(function (element) {
        const diveId = element.getAttribute('data-dive-id');
        loadScores(diveId);
    });

});
